package api

import (
	"bytes"
	"net/http"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

func (s *Server) handleReadme(w http.ResponseWriter, r *http.Request) {
	md, err := os.ReadFile("README.md")
	if err != nil {
		writeError(w, http.StatusNotFound, "README.md not found")
		return
	}

	parts := strings.SplitN(string(md), "\n---\n", 3)
	var enMd, zhMd string
	if len(parts) >= 3 {
		enMd = strings.TrimSpace(parts[1])
		zhMd = strings.TrimSpace(parts[2])
	} else {
		enMd = string(md)
		zhMd = string(md)
	}

	enHTML, err := renderMarkdown(enMd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to render markdown")
		return
	}
	zhHTML, err := renderMarkdown(zhMd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to render markdown")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"en": enHTML,
		"zh": zhHTML,
	})
}

func renderMarkdown(src string) (string, error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.Table),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(src), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
