package services

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/a-h/templ"

	"potblog/components"
)

var (
	_, b, _, _ = runtime.Caller(0)
	Root       = filepath.Join(filepath.Dir(b), "..")
)

func offlineRender(t templ.Component) string {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(context.Background(), buf); err != nil {
		return ""
	}

	return buf.String()
}

func ReadMarkdownFile(relative_filepath string) string {
	filepath := filepath.Join(Root, relative_filepath)
	file, err := os.Open(filepath)
	if err != nil {
		return ""
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return ""
	}

	return string(data)
}

type Metadata struct {
	Title       string
	Description string
	Date        string
	Tags        []string
	Author      string
}

type MarkdownHTML struct {
	RawHTML  string
	Metadata Metadata
}

func MarkdownToHTML(md *string) (MarkdownHTML, error) {
	metadata, err := markdownToMetadata(md)
	if err != nil {
		return MarkdownHTML{}, err
	}

	html, err := markdownToRawHTML(md)
	if err != nil {
		return MarkdownHTML{}, err
	}

	return MarkdownHTML{
		RawHTML:  html,
		Metadata: metadata,
	}, nil
}

var MetadataTagMap = map[string]string{
	"title":       "Title",
	"description": "Description",
	"date":        "Date",
	"tags":        "Tags",
	"author":      "Author",
}

func markdownToMetadata(md *string) (Metadata, error) {
	metadataBlock, err := extractMetadataBlock(*md)
	if err != nil {
		return Metadata{}, err
	}

	metadataLines := strings.Split(metadataBlock, "\n")

	metadata := Metadata{}
	for _, line := range metadataLines {
		mdTag := strings.Split(line, ":")
		if len(mdTag) < 2 {
			continue
		}

		tag, value, err := parseMetadataLine(line)
		if err != nil {
			continue
		}

		if metadataTag, ok := MetadataTagMap[tag]; ok {
			switch metadataTag {
			case "Title":
				metadata.Title = value
			case "Description":
				metadata.Description = value
			case "Date":
				metadata.Date = value
			case "Tags":
				metadata.Tags = parseTags(value)
			case "Author":
				metadata.Author = value
			}
		}
	}

	return metadata, nil
}

func extractMetadataBlock(content string) (string, error) {
	if !strings.HasPrefix(content, "---") {
		return "", fmt.Errorf("invalid metadata format: missing opening separator")
	}

	indexSecondSeparator := strings.Index(content[3:], "---")
	if indexSecondSeparator == -1 {
		return "", fmt.Errorf("invalid metadata format: missing closing separator")
	}

	return content[4 : indexSecondSeparator+2], nil
}

func parseMetadataLine(line string) (string, string, error) {
	parts := strings.Split(line, ":")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid metadata line format")
	}

	tag := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(strings.Join(parts[1:], ":"))
	return tag, value, nil
}

func parseTags(value string) []string {
	var tags []string
	for _, tag := range strings.Split(value, ",") {
		tags = append(tags, strings.TrimSpace(tag))
	}
	return tags
}

func markdownToRawHTML(md *string) (string, error) {
	rows := strings.Split(skipMetadataBlock(md), "\n")

	var html strings.Builder

	idx := 0
	for idx < len(rows) {
		row := rows[idx]
		row_type := rowType(row)

		switch row_type {
		case "title_h1":
			html.WriteString(offlineRender(components.TitleH1(row[2:])))
		case "title_h2":
			html.WriteString(offlineRender(components.TitleH2(row[3:])))
		case "paragraph":
			html.WriteString(offlineRender(components.Paragraph(row)))
		case "quote":
			html.WriteString(offlineRender(components.Blockquote(row[2:])))
		case "code":
			language := strings.Trim(row, "`")

			codeBlockMd := ""
			for _, codeRow := range rows[idx+1:] {
				idx++
				if strings.HasPrefix(codeRow, "```") {
					break
				}
				codeBlockMd += codeRow + "\n"
			}

			codeHash := generateHashFromCodeBlock(codeBlockMd)
			html.WriteString(offlineRender(components.CodeBlock(language, codeBlockMd, codeHash)))
		case "button":
			url, icon, text := extractButtonTags(row)
			html.WriteString(offlineRender(components.Button(url, icon, text)))
		case "empty":
			html.WriteString("\n\n")
		}

		idx++
	}

	return boldify(html.String()), nil
}

func skipMetadataBlock(content *string) string {
	indexSecondSeparator := strings.Index((*content)[3:], "---")
	if indexSecondSeparator == -1 {
		return *content
	}

	return strings.Trim((*content)[indexSecondSeparator+6:], "\n")
}

func rowType(row string) string {
	if strings.HasPrefix(row, "# ") {
		return "title_h1"
	}
	if strings.HasPrefix(row, "## ") {
		return "title_h2"
	}
	if strings.HasPrefix(row, "> ") {
		return "quote"
	}
	if strings.HasPrefix(row, "```") {
		return "code"
	}
	if strings.HasPrefix(row, "[button") {
		return "button"
	}
	if row == "" {
		return "empty"
	}
	return "paragraph"
}

func generateHashFromCodeBlock(code string) string {
	hasher := sha256.New()
	hasher.Write([]byte(code))
	hash := hasher.Sum(nil)

	encoded := base64.URLEncoding.EncodeToString(hash)

	return encoded[:8]
}

func extractButtonTags(row string) (string, string, string) {
	innerData := row[7 : len(row)-1]

	tags := strings.Split(innerData, " ")
	var tagMap = make(map[string]string)
	for _, tag := range tags {
		tagData := strings.Split(tag, "=")
		if len(tagData) < 2 {
			continue
		}
		tagMap[tagData[0]] = strings.Trim(strings.Trim(tagData[1], "'"), "\"")
	}

	url := tagMap["url"]
	icon := tagMap["icon"]
	text := tagMap["text"]

	return url, icon, text
}

func boldify(row string) string {
	starCount := strings.Count(row, "**")
	if starCount%2 != 0 {
		return row
	}

	for i := 0; i < starCount/2; i++ {
		row = strings.Replace(row, "**", "<b>", 1)
		row = strings.Replace(row, "**", "</b>", 1)
	}

	return row
}
