package usecase

type SearchInteractor struct {
	searcher Searcher
}

func NewSearchInteractor(searcher Searcher) *SearchInteractor {
	return &SearchInteractor{searcher: searcher}
}
