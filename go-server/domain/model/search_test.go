package model

import "testing"

func TestSearchResult(t *testing.T) {
	result := SearchResult{
		FilePath:   "file.go",
		LineNumber: 41,
		LineText:   "Hello, World",
	}
	if result.FilePath != "file.go" || result.LineNumber != 41 || result.LineText != "Hello, World" {
		t.Errorf("SearchResult not assigned properly: %+v", result)
	}
}
