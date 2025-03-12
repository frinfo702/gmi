package usecase

import (
	"github.com/frinfo702/fixer/domain/model"
	"github.com/frinfo702/fixer/infrastructure"
)

type SearchInteractor struct {
	searcher infrastructure.Searcher
}

func NewSearchInteractor(searcher infrastructure.Searcher) *SearchInteractor {
	return &SearchInteractor{searcher: searcher}
}

// Executeは指定された検索クエリに基づき検索を実行
func (s *SearchInteractor) Execute(query model.SearchQuery) ([]model.SearchResult, error) {
	return s.searcher.Search(query)
}
