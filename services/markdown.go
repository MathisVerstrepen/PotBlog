package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

var MetadataTagMap = map[string]string{
	"title":       "Title",
	"description": "Description",
	"date":        "Date",
	"tags":        "Tags",
	"author":      "Author",
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

		tag := strings.TrimSpace(mdTag[0])
		value := strings.TrimSpace(strings.Join(mdTag[1:], ":"))

		if metadataTag, ok := MetadataTagMap[tag]; ok {
			switch metadataTag {
			case "Title":
				metadata.Title = value
			case "Description":
				metadata.Description = value
			case "Date":
				metadata.Date = value
			case "Tags":
				metadata.Tags = []string{}
				for _, tag := range strings.Split(value, ",") {
					metadata.Tags = append(metadata.Tags, strings.TrimSpace(tag))
				}
			case "Author":
				metadata.Author = value
			}
		}
	}

	return metadata, nil
}

func MarkdownToHTML(md *string) (MarkdownHTML, error) {
	metadata, err := markdownToMetadata(md)
	if err != nil {
		return MarkdownHTML{}, err
	}

	return MarkdownHTML{
		RawHTML:  *md,
		Metadata: metadata,
	}, nil
}
