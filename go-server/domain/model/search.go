package model

// SearchResultは1件の検索結果を表す
type SearchResult struct {
	FilePath   string
	LineNumber int
	LineText   string
}

// SearchQueryは検索条件を保持
type SearchQuery struct {
	Path  string
	Query string
	Fuzzy bool // Fuzzy means if obvious search is turned on.
}

// Searcherインターフェースは検索機能の抽象化をする
type Searcher interface {
	Search(query SearchQuery) ([]SearchResult, error)
}
