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
