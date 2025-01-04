package services

import (
	"io"
	"os"
	"path/filepath"
	"potblog/infrastructure"
	"reflect"
	"strings"
	"testing"
)

func pointerTo[T ~string](s T) *T {
	return &s
}

func getTestFileData(filename string) string {
	filepath := filepath.Join(Root, "tests/data", filename)
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

func Test_ReadMarkdownFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "givenAFile_WhenReadMarkdownFile_ThenReturnFileContent",
			args: args{
				filename: "./tests/data/TestReadMarkdownFile.md",
			},
			want: "# Test",
		}, {
			name: "givenAnInvalidFile_WhenReadMarkdownFile_ThenReturnEmptyString",
			args: args{
				filename: "./tests/data/InvalidFile.md",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadMarkdownFile(tt.args.filename); got != tt.want {
				t.Errorf("ReadMarkdownFile() = %v\nWant %v", got, tt.want)
			}
		})
	}
}

func TestMarkdownToHTML(t *testing.T) {
	type args struct {
		md *string
	}
	tests := []struct {
		name    string
		args    args
		want    MarkdownHTML
		wantErr bool
	}{
		{
			name: "givenMarkdown_WhenMarkdownToHTML_ThenReturnHTML",
			args: args{
				md: pointerTo(getTestFileData("TestMarkdownToHTML.md")),
			},
			want: MarkdownHTML{
				RawHTML: getTestFileData("TestMarkdownToHTMLResult.html"),
				Metadata: infrastructure.Metadata{
					Title:       "This is an article !",
					Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec nec odio vitae nunc.",
					Date:        "2024-08-13",
					Tags:        []string{"lorem", "ipsum"},
					Author:      "Mathis Verstrepen",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertMarkdownToHTML(tt.args.md)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarkdownToHTML() error = %v \nWantErr %v", err, tt.wantErr)
				return
			}

			gotNormalized := strings.ReplaceAll(got.RawHTML, "\n", "")
			wantNormalized := strings.ReplaceAll(tt.want.RawHTML, "\n", "")

			if !reflect.DeepEqual(gotNormalized, wantNormalized) {
				t.Errorf("MarkdownToHTML() = %v \nWant %v", got.RawHTML, tt.want.RawHTML)
			}
			if !reflect.DeepEqual(got.Metadata, tt.want.Metadata) {
				t.Errorf("MarkdownToHTML() = %v \nWant %v", got.Metadata, tt.want.Metadata)
			}
		})
	}
}

func Test_markdownToMetadata(t *testing.T) {
	type args struct {
		md *string
	}
	tests := []struct {
		name    string
		args    args
		want    infrastructure.Metadata
		wantErr bool
	}{
		{
			name: "givenMarkdown_WhenMarkdownToMetadata_ThenReturnMetadata",
			args: args{
				md: pointerTo(`---
title: This is an article !
description: desc.
date: 2024-08-13
tags: lorem, ipsum
author: Mathis Verstrepen
---`),
			},
			want: infrastructure.Metadata{
				Title:       "This is an article !",
				Description: "desc.",
				Date:        "2024-08-13",
				Tags:        []string{"lorem", "ipsum"},
				Author:      "Mathis Verstrepen",
			},
			wantErr: false,
		}, {
			name: "givenInvalidMarkdown_WhenMarkdownToMetadata_ThenReturnEmptyMetadata",
			args: args{
				md: pointerTo(`# This is an article !`),
			},
			want:    infrastructure.Metadata{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractMetadataFromMarkdown(tt.args.md)
			if (err != nil) != tt.wantErr {
				t.Errorf("markdownToMetadata() error = %v\nWantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("markdownToMetadata() = %v\nWant %v", got, tt.want)
			}
		})
	}
}

func Test_extractMetadataBlock(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "givenMarkdown_WhenExtractMetadataBlock_ThenReturnMetadataBlock",
			args: args{
				content: `---
title: This is an article !
description: desc.
---`,
			},
			want: `title: This is an article !
description: desc.`,
			wantErr: false,
		}, {
			name: "givenInvalidMarkdown_WhenExtractMetadataBlock_ThenReturnError",
			args: args{
				content: `# This is an article !`,
			},
			want:    "",
			wantErr: true,
		}, {
			name: "givenInvalidMarkdown_WhenExtractMetadataBlock_ThenReturnError",
			args: args{
				content: `---
title: This is an article !
`,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := retrieveMetadataSection(tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractMetadataBlock() error = %v\nWantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractMetadataBlock() = %v\nWant %v", got, tt.want)
			}
		})
	}
}

func Test_parseMetadataLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "givenMetadataLine_WhenParseMetadataLine_ThenReturnTagAndValue",
			args: args{
				line: "title: This is an article !",
			},
			want:    "title",
			want1:   "This is an article !",
			wantErr: false,
		}, {
			name: "givenInvalidMetadataLine_WhenParseMetadataLine_ThenReturnError",
			args: args{
				line: "title This is an article !",
			},
			want:    "",
			want1:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := processMetadataEntry(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMetadataLine() error = %v\nWantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseMetadataLine() got = %v\nWant %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseMetadataLine() got1 = %v\nWant %v", got1, tt.want1)
			}
		})
	}
}

func Test_parseTags(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "givenTags_WhenParseTags_ThenReturnTags",
			args: args{
				value: " lorem ,  ipsum ",
			},
			want: []string{"lorem", "ipsum"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractTags(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTags() = %v\nWant %v", got, tt.want)
			}
		})
	}
}

func Test_markdownToRawHTML(t *testing.T) {
	type args struct {
		md *string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "givenTitleH1Markdown_WhenMarkdownToRawHTML_ThenReturnHTML",
			args: args{
				md: pointerTo("# This is an article !"),
			},
			want:    []string{`<h1 class="article-title-h1">This is an article !</h1>`},
			wantErr: false,
		}, {
			name: "givenTitleH2Markdown_WhenMarkdownToRawHTML_ThenReturnHTML",
			args: args{
				md: pointerTo("## This is an article !"),
			},
			want:    []string{`<h2 class="article-title-h2">This is an article !</h2>`},
			wantErr: false,
		}, {
			name: "givenTitleParagraphMarkdown_WhenMarkdownToRawHTML_ThenReturnHTML",
			args: args{
				md: pointerTo("This is a paragraph."),
			},
			want:    []string{`<p class="article-paragraph">This is a paragraph.</p>`},
			wantErr: false,
		}, {
			name: "givenTitleQuoteMarkdown_WhenMarkdownToRawHTML_ThenReturnHTML",
			args: args{
				md: pointerTo("> This is a quote."),
			},
			want:    []string{`<blockquote class="article-blockquote standard"><hr>`, `This is a quote.</blockquote>`},
			wantErr: false,
		}, {
			name: "givenTitleQuoteMarkdown_WhenMarkdownToRawHTML_ThenReturnHTML",
			args: args{
				md: pointerTo("> [!WARNING] This is a quote."),
			},
			want:    []string{`<blockquote class="article-blockquote warning"><hr>`, `This is a quote.</blockquote>`},
			wantErr: false,
		}, {
			name: "givenTitleQuoteMarkdown_WhenMarkdownToRawHTML_ThenReturnHTML",
			args: args{
				md: pointerTo("> [!IMPORTANT] This is a quote."),
			},
			want:    []string{`<blockquote class="article-blockquote important"><hr>`, `This is a quote.</blockquote>`},
			wantErr: false,
		}, {
			name: "givenTitleCodeMarkdown_WhenMarkdownToRawHTML_ThenReturnHTML",
			args: args{
				md: pointerTo("```python\nprint('Hello, World!')\n```"),
			},
			want:    []string{"<div class=\"article-codeblock\">", "<button class=\"article-copy-button\"", "</button><p class=\"article-language\">python</p><pre class=\"article-code\"><code", "print('Hello, World!')\n</code></pre></div>"},
			wantErr: false,
		}, {
			name: "givenTitleButtonMarkdown_WhenMarkdownToRawHTML_ThenReturnHTML",
			args: args{
				md: pointerTo("[button url='https://github.com/MathisVerstrepen' text='Github']"),
			},
			want:    []string{`<a href="https://github.com/MathisVerstrepen" role="button" class="article-button" target="_blank">Github</a>`},
			wantErr: false,
		}, {
			name: "givenTitleButtonMarkdown_WhenMarkdownToRawHTML_ThenReturnHTML",
			args: args{
				md: pointerTo("![This is an image](https://superimage.com)"),
			},
			want:    []string{`<figure class="article-image"><img src="https://superimage.com" alt="This is an image"><figcaption>This is an image</figcaption></figure>`},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertMarkdownToHTML(tt.args.md)
			if (err != nil) != tt.wantErr {
				t.Errorf("markdownToRawHTML() error = %v\nWantErr %v", err, tt.wantErr)
				return
			}
			for i, html := range tt.want {
				if !strings.Contains(got, html) {
					t.Errorf("markdownToRawHTML() = %v\nWant %v", got, tt.want[i])
				}
			}
		})
	}
}

func Test_skipMetadataBlock(t *testing.T) {
	type args struct {
		content *string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "givenMarkdown_WhenSkipMetadataBlock_ThenReturnMarkdownWithoutMetadata",
			args: args{
				content: pointerTo(`---
title: This is an article !
description: desc.
---

# This is an article !`),
			},
			want: "# This is an article !",
		}, {
			name: "givenMarkdownWithoutMetadata_WhenSkipMetadataBlock_ThenReturnMarkdown",
			args: args{
				content: pointerTo("# This is an article !"),
			},
			want: "# This is an article !",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeMetadataBlock(tt.args.content); got != tt.want {
				t.Errorf("skipMetadataBlock() = %v \nWant %v", got, tt.want)
			}
		})
	}
}

func Test_rowType(t *testing.T) {
	type args struct {
		row string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "givenTitleH1Row_WhenRowType_ThenReturnTitleH1",
			args: args{
				row: "# This is an article !",
			},
			want: "title_h1",
		}, {
			name: "givenTitleH2Row_WhenRowType_ThenReturnTitleH2",
			args: args{
				row: "## This is an article !",
			},
			want: "title_h2",
		}, {
			name: "givenParagraphRow_WhenRowType_ThenReturnParagraph",
			args: args{
				row: "This is a paragraph.",
			},
			want: "paragraph",
		}, {
			name: "givenQuoteRow_WhenRowType_ThenReturnQuote",
			args: args{
				row: "> This is a quote.",
			},
			want: "quote",
		}, {
			name: "givenQuoteRow_WhenRowType_ThenReturnQuote",
			args: args{
				row: "> [!WARNING] This is a quote.",
			},
			want: "quote-warning",
		}, {
			name: "givenQuoteRow_WhenRowType_ThenReturnQuote",
			args: args{
				row: "> [!IMPORTANT] This is a quote.",
			},
			want: "quote-important",
		}, {
			name: "givenCodeRow_WhenRowType_ThenReturnCode",
			args: args{
				row: "```python",
			},
			want: "code",
		}, {
			name: "givenQuoteRow_WhenRowType_ThenReturnQuote",
			args: args{
				row: "![This is an image](https://superimage.com)",
			},
			want: "image",
		}, {
			name: "givenButtonRow_WhenRowType_ThenReturnButton",
			args: args{
				row: "[button url='https://github.com/MathisVerstrepen' icon='github' text='Github']",
			},
			want: "button",
		}, {
			name: "givenEmptyRow_WhenRowType_ThenReturnEmpty",
			args: args{
				row: "",
			},
			want: "empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rowType(tt.args.row); got != tt.want {
				t.Errorf("rowType() = %v\nWant %v", got, tt.want)
			}
		})
	}
}

func Test_generateHashFromCodeBlock(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "givenCodeBlock_WhenGenerateHashFromCodeBlock_ThenReturnHash",
			args: args{
				code: "print('Hello, World!')\n",
			},
			want: "XEkuiIsx",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hashFromCodeBlock(tt.args.code); got != tt.want {
				t.Errorf("generateHashFromCodeBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractButtonTags(t *testing.T) {
	type args struct {
		row string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
		want2 string
	}{
		{
			name: "givenButtonRow_WhenExtractButtonTags_ThenReturnButtonTags",
			args: args{
				row: "[button url='https://github.com/MathisVerstrepen' icon='github' text='Github']",
			},
			want:  "https://github.com/MathisVerstrepen",
			want1: "github",
			want2: "Github",
		}, {
			name: "givenInvalidButtonRow_WhenExtractButtonTags_ThenReturnEmptyButtonTags",
			args: args{
				row: "[button]",
			},
			want:  "",
			want1: "",
			want2: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := extractButtonProperties(tt.args.row)
			if got != tt.want {
				t.Errorf("extractButtonTags() got = %v\nWant %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("extractButtonTags() got1 = %v\nWant %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("extractButtonTags() got2 = %v\nWant %v", got2, tt.want2)
			}
		})
	}
}

func Test_extractImageTags(t *testing.T) {
	type args struct {
		row string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "givenImageRow_WhenExtractImageTags_ThenReturnImageTags",
			args: args{
				row: "![This is an image](https://superimage.com)",
			},
			want:  "This is an image",
			want1: "https://superimage.com",
		}, {
			name: "givenInvalidImageRow_WhenExtractImageTags_ThenReturnEmptyImageTags",
			args: args{
				row: "![This is an image]",
			},
			want:  "",
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := extractImageDetails(tt.args.row)
			if got != tt.want {
				t.Errorf("extractImageTags() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("extractImageTags() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_boldify(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "givenTextWithNoAsterix_WhenBoldify_ThenReturnText",
			args: args{
				text: "This is a paragraph.",
			},
			want: "This is a paragraph.",
		}, {
			name: "givenTextWithAsterix_WhenBoldify_ThenReturnBoldText",
			args: args{
				text: "**This is a paragraph.**",
			},
			want: "<b>This is a paragraph.</b>",
		}, {
			name: "givenTextWithMultipleAsterix_WhenBoldify_ThenReturnBoldText",
			args: args{
				text: "**This is a paragraph.** **This is a paragraph.**",
			},
			want: "<b>This is a paragraph.</b> <b>This is a paragraph.</b>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := boldify(tt.args.text); got != tt.want {
				t.Errorf("boldify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_linkify(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "givenTextWithNoLink_WhenLinkify_ThenReturnText",
			args: args{
				text: "This is a paragraph.",
			},
			want: "This is a paragraph.",
		}, {
			name: "givenTextWithLink_WhenLinkify_ThenReturnLink",
			args: args{
				text: `[This is a link](https://superlink.com)`,
			},
			want: `<a href="https://superlink.com" class="article-external-link" target="_blank">This is a link🡵</a>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := linkify(tt.args.text); got != tt.want {
				t.Errorf("linkify() = %v, want %v", got, tt.want)
			}
		})
	}
}
