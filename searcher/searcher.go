package searcher

import (
	"fmt"
	"gmi/indexer"
	"gmi/tokenizer"
	"sort"
)

// SearchResult は検索結果の1つのアイテムを表します。
type SearchResult struct {
	Document  indexer.Document
	Positions []int // クエリ単語の出現位置
}

// Search は指定されたインデックス内でクエリに一致するドキュメントを検索します。
// 現状は単一キーワード検索のみをサポートします。
func Search(idx *indexer.InvertedIndex, query string) []SearchResult {
	var results []SearchResult

	if idx == nil || idx.Index == nil || idx.Docs == nil {
		fmt.Println("Error: Index is not properly initialized.")
		return results
	}

	queryTokens := tokenizer.Tokenize(query)
	if len(queryTokens) == 0 {
		fmt.Println("Warning: Empty query after tokenization.")
		return results
	}

	searchTerm := queryTokens[0]
	fmt.Printf("Searching for term: '%s'\n", searchTerm)

	postings, found := idx.Index[searchTerm]
	if !found {
		fmt.Printf("Term '%s' not found in index.\n", searchTerm)
		return results
	}

	for _, posting := range postings {
		doc, docExists := idx.Docs[posting.DocID]
		if docExists {
			results = append(results, SearchResult{
				Document:  doc,
				Positions: posting.Positions,
			})
		} else {
			fmt.Printf("Warning: DocID %d found in posting for term '%s', but not in Docs map.\n", posting.DocID, searchTerm)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Document.ID < results[j].Document.ID
	})

	return results
}
