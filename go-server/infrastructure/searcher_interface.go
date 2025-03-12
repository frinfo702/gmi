package infrastructure

import "github.com/frinfo702/fixer/domain/model"

type Searcher interface {
	Search(query model.SearchQuery) ([]model.SearchResult, error)
}
