package services

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/a-h/templ"

	"potblog/components"
	"potblog/infrastructure"
)

var (
	Root = os.Getenv("ROOT_DIR")
)

func renderTemplate(t templ.Component) string {
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

type MarkdownHTML struct {
	RawHTML  string
	Metadata infrastructure.Metadata
}

func ConvertMarkdownToHTML(md *string) (MarkdownHTML, error) {
	metadata, err := extractMetadataFromMarkdown(md)
	if err != nil {
		return MarkdownHTML{}, err
	}

	html, err := convertMarkdownToHTML(md)
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

func extractMetadataFromMarkdown(md *string) (infrastructure.Metadata, error) {
	metadataSection, err := retrieveMetadataSection(*md)
	if err != nil {
		return infrastructure.Metadata{}, err
	}

	metadataLines := strings.Split(metadataSection, "\n")

	metadata := infrastructure.Metadata{}
	for _, line := range metadataLines {
		mdTag := strings.Split(line, ":")
		if len(mdTag) < 2 {
			continue
		}

		tag, value, err := processMetadataEntry(line)
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
				metadata.Tags = extractTags(value)
			case "Author":
				metadata.Author = value
			}
		}
	}

	return metadata, nil
}

func retrieveMetadataSection(content string) (string, error) {
	if !strings.HasPrefix(content, "---") {
		return "", fmt.Errorf("invalid metadata format: missing opening separator")
	}

	indexSecondSeparator := strings.Index(content[3:], "---")
	if indexSecondSeparator == -1 {
		return "", fmt.Errorf("invalid metadata format: missing closing separator")
	}

	return content[4 : indexSecondSeparator+2], nil
}

func processMetadataEntry(line string) (string, string, error) {
	parts := strings.Split(line, ":")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid metadata line format")
	}

	tag := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(strings.Join(parts[1:], ":"))
	return tag, value, nil
}

func extractTags(value string) []string {
	var tags []string
	for _, tag := range strings.Split(value, ",") {
		tags = append(tags, strings.TrimSpace(tag))
	}
	return tags
}

func convertMarkdownToHTML(md *string) (string, error) {
	rows := strings.Split(removeMetadataBlock(md), "\n")

	var html strings.Builder

	idx := 0
	for idx < len(rows) {
		row := rows[idx]
		row_type := rowType(row)

		switch row_type {
		case "title_h1":
			html.WriteString(renderTemplate(components.TitleH1(row[2:])))
		case "title_h2":
			html.WriteString(renderTemplate(components.TitleH2(row[3:])))
		case "paragraph":
			html.WriteString(renderTemplate(components.Paragraph(row)))
		case "quote-warning":
			html.WriteString(renderTemplate(components.Blockquote(row[13:], "warning")))
		case "quote-important":
			html.WriteString(renderTemplate(components.Blockquote(row[15:], "important")))
		case "quote":
			html.WriteString(renderTemplate(components.Blockquote(row[2:], "standard")))
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

			codeHash := hashFromCodeBlock(codeBlockMd)
			html.WriteString(renderTemplate(components.CodeBlock(language, codeBlockMd, codeHash)))
		case "button":
			url, icon, text := extractButtonProperties(row)
			html.WriteString(renderTemplate(components.Button(url, icon, text)))
		case "image":
			caption, url := extractImageDetails(row)
			html.WriteString(renderTemplate(components.Image(url, caption)))
		case "empty":
			html.WriteString("\n\n")
		}

		idx++
	}

	htmlStr := html.String()
	htmlStr = linkify(htmlStr)
	htmlStr = boldify(htmlStr)

	return htmlStr, nil
}

func removeMetadataBlock(content *string) string {
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
	if strings.HasPrefix(row, "> [!WARNING] ") {
		return "quote-warning"
	}
	if strings.HasPrefix(row, "> [!IMPORTANT] ") {
		return "quote-important"
	}
	if strings.HasPrefix(row, "> ") {
		return "quote"
	}
	if strings.HasPrefix(row, "```") {
		return "code"
	}
	if strings.HasPrefix(row, "![") {
		return "image"
	}
	if strings.HasPrefix(row, "[button") {
		return "button"
	}
	if row == "" {
		return "empty"
	}
	return "paragraph"
}

func hashFromCodeBlock(code string) string {
	hasher := sha256.New()
	hasher.Write([]byte(code))
	hash := hasher.Sum(nil)

	encoded := base64.URLEncoding.EncodeToString(hash)

	return encoded[:8]
}

func extractButtonProperties(row string) (string, string, string) {
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

func extractImageDetails(row string) (string, string) {
	re := regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)
	matches := re.FindAllStringSubmatch(row, -1)

	if len(matches) == 0 {
		return "", ""
	}

	return matches[0][1], matches[0][2]
}

func boldify(text string) string {
	starCount := strings.Count(text, "**")
	if starCount%2 != 0 {
		return text
	}

	for i := 0; i < starCount/2; i++ {
		text = strings.Replace(text, "**", "<b>", 1)
		text = strings.Replace(text, "**", "</b>", 1)
	}

	return text
}

func linkify(text string) string {
	re := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	matches := re.FindAllStringSubmatchIndex(text, -1)

	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		linkText := text[match[2]:match[3]]
		linkURL := text[match[4]:match[5]]
		replacement := renderTemplate(components.ExternalLink(linkURL, linkText))
		text = text[:match[0]] + replacement + text[match[1]:]
	}

	return text
}
