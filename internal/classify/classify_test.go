package classify

import (
	"sort"
	"testing"
)

func TestDetectCategory(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{"plain text", "some random note about Go", "note"},
		{"http link", "https://golang.org/doc/effective_go", "link"},
		{"http with path", "http://example.com/path?q=1", "link"},
		{"link with extra text", "https://github.com/foo/bar\nsome description", "link"},
		{"image png", "screenshot.png", "image"},
		{"image jpg", "photo.jpg", "image"},
		{"image webp", "banner.webp", "image"},
		{"pdf file", "report.pdf", "file"},
		{"docx file", "notes.docx", "file"},
		{"zip file", "archive.zip", "file"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectCategory(tt.content)
			if got != tt.want {
				t.Errorf("detectCategory(%q) = %q, want %q", tt.content, got, tt.want)
			}
		})
	}
}

func TestExtractTags(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		content string
		want    []string
	}{
		{
			"go concurrency",
			"Go并发模式",
			"goroutine channel mutex",
			[]string{"concurrency", "go"},
		},
		{
			"docker",
			"Docker部署",
			"dockerfile compose container",
			[]string{"docker"},
		},
		{
			"python api",
			"Flask API",
			"python flask rest api endpoint",
			[]string{"api", "python"},
		},
		{
			"no match",
			"random title",
			"nothing special here",
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTags(tt.title, tt.content)
			sort.Strings(got)
			sort.Strings(tt.want)

			if len(got) == 0 && len(tt.want) == 0 {
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("extractTags() = %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("extractTags() = %v, want %v", got, tt.want)
					return
				}
			}
		})
	}
}

func TestClassify(t *testing.T) {
	r := Classify("Go goroutine patterns", "goroutine channel select")
	if r.Category != "note" {
		t.Errorf("Category = %q, want note", r.Category)
	}
	if len(r.Tags) == 0 {
		t.Error("expected tags, got none")
	}

	r = Classify("bookmark", "https://go.dev")
	if r.Category != "link" {
		t.Errorf("Category = %q, want link", r.Category)
	}
}

func TestFormatTags(t *testing.T) {
	got := FormatTags([]string{"go", "docker", "api"})
	want := "go,docker,api"
	if got != want {
		t.Errorf("FormatTags() = %q, want %q", got, want)
	}
}
