package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"potblog/components"
	"potblog/handlers"
)

var (
	_, b, _, _ = runtime.Caller(0)
	Root       = filepath.Join(filepath.Dir(b), "..")
)

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
	rows := strings.Split(*md, "\n")

	var html string

	idx := 0
	for idx < len(rows) {
		row := rows[idx]
		row_type := rowType(row)

		switch row_type {
		case "title_h1":
			html += handlers.OfflineRender(components.TitleH1(row[2:]))
		case "title_h2":
			html += handlers.OfflineRender(components.TitleH2(row[3:]))
		case "paragraph":
			html += handlers.OfflineRender(components.Paragraph(row))
		case "quote":
			html += handlers.OfflineRender(components.Blockquote(row[2:]))
		case "code":
			language := strings.Trim(row, "`")
			fmt.Println(language)

			codeBlockMd := ""
			for _, codeRow := range rows[idx+1:] {
				idx++
				if strings.HasPrefix(codeRow, "```") {
					break
				}
				codeBlockMd += codeRow + "\n"
			}

			html += handlers.OfflineRender(components.CodeBlock(language, codeBlockMd))
		case "button":
			url, icon, text := extractButtonTags(row)

			fmt.Printf("url: %s, icon: %s, text: %s\n", url, icon, text)

			html += handlers.OfflineRender(components.Button(url, icon, text))
		}

		idx++
	}

	return html, nil
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
	return "paragraph"
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
