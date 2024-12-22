package services

import (
	"reflect"
	"testing"
)

func PointerTo[T ~string](s T) *T {
	return &s
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
				t.Errorf("ReadMarkdownFile() = %v, want %v", got, tt.want)
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
			got, err := extractMetadataBlock(tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractMetadataBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractMetadataBlock() = %v, want %v", got, tt.want)
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
			got, got1, err := parseMetadataLine(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMetadataLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseMetadataLine() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseMetadataLine() got1 = %v, want %v", got1, tt.want1)
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
			if got := parseTags(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTags() = %v, want %v", got, tt.want)
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
		want    Metadata
		wantErr bool
	}{
		{
			name: "givenMarkdown_WhenMarkdownToMetadata_ThenReturnMetadata",
			args: args{
				md: PointerTo(`---
title: This is an article !
description: desc.
date: 2024-08-13
tags: lorem, ipsum
author: Mathis Verstrepen
---`),
			},
			want: Metadata{
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
				md: PointerTo(`# This is an article !`),
			},
			want:    Metadata{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := markdownToMetadata(tt.args.md)
			if (err != nil) != tt.wantErr {
				t.Errorf("markdownToMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("markdownToMetadata() = %v, want %v", got, tt.want)
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarkdownToHTML(tt.args.md)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarkdownToHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarkdownToHTML() = %v, want %v", got, tt.want)
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
		want    string
		wantErr bool
	}{
		{
			name: "givenTitleH1Markdown_WhenMarkdownToRawHTML_ThenReturnHTML",
			args: args{
				md: PointerTo("# This is an article !"),
			},
			want:    "<h1>This is an article !</h1>",
			wantErr: false,
		}, {
			name: "givenTitleH2Markdown_WhenMarkdownToRawHTML_ThenReturnHTML",
			args: args{
				md: PointerTo("## This is an article !"),
			},
			want:    "<h2>This is an article !</h2>",
			wantErr: false,
		}, {
			name: "givenTitleParagraphMarkdown_WhenMarkdownToRawHTML_ThenReturnHTML",
			args: args{
				md: PointerTo("This is a paragraph."),
			},
			want:    "<p>This is a paragraph.</p>",
			wantErr: false,
		}, {
			name: "givenTitleQuoteMarkdown_WhenMarkdownToRawHTML_ThenReturnHTML",
			args: args{
				md: PointerTo("> This is a quote."),
			},
			want:    "<blockquote>This is a quote.</blockquote>",
			wantErr: false,
		}, {
			name: "givenTitleCodeMarkdown_WhenMarkdownToRawHTML_ThenReturnHTML",
			args: args{
				md: PointerTo("```python\nprint('Hello, World!')\n```"),
			},
			want:    "<div><p>python</p><pre><code>print('Hello, World!')\n</code></pre></div>",
			wantErr: false,
		}, {
			name: "givenTitleButtonMarkdown_WhenMarkdownToRawHTML_ThenReturnHTML",
			args: args{
				md: PointerTo("[button url='https://github.com/MathisVerstrepen' text='Github']"),
			},
			want:    `<a href="https://github.com/MathisVerstrepen" role="button">Github</a>`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := markdownToRawHTML(tt.args.md)
			if (err != nil) != tt.wantErr {
				t.Errorf("markdownToRawHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("markdownToRawHTML() = %v, want %v", got, tt.want)
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
			name: "givenCodeRow_WhenRowType_ThenReturnCode",
			args: args{
				row: "```python",
			},
			want: "code",
		}, {
			name: "givenButtonRow_WhenRowType_ThenReturnButton",
			args: args{
				row: "[button url='https://github.com/MathisVerstrepen' icon='github' text='Github']",
			},
			want: "button",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rowType(tt.args.row); got != tt.want {
				t.Errorf("rowType() = %v, want %v", got, tt.want)
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
			got, got1, got2 := extractButtonTags(tt.args.row)
			if got != tt.want {
				t.Errorf("extractButtonTags() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("extractButtonTags() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("extractButtonTags() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
