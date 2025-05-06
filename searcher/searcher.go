package searcher

import (
	"fmt"
	"gmi/indexer"
	"gmi/tokenizer"
	"sort"
)

// SearchResult は検索結果の1つのアイテムを表します。
type SearchResult struct {
	Document           indexer.Document
	QueryTermPositions map[string][]int // key: 検索クエリのトークン, value: そのトークンの出現位置リスト
}

// Searchは指定されたインデックス内でクエリに一致するドキュメントを検索します
func Search(idx *indexer.InvertedIndex, query string) []SearchResult {
	var finalResults []SearchResult

	if idx == nil || idx.Index == nil || idx.Docs == nil {
		fmt.Println("Error: Index is not properly initialized.")
		return finalResults
	}

	queryTokens := tokenizer.Tokenize(query)
	if len(queryTokens) == 0 {
		fmt.Println("Warning: Empty query after tokenization.")
		return finalResults
	}

	fmt.Printf("Searching for terms (AND): %v\n", queryTokens)

	postingLists := make(map[string][]indexer.Posting)
	shortestPostingListLength := -1
	var shortestToken string

	for _, token := range queryTokens {
		postings, found := idx.Index[token]
		if !found {
			fmt.Printf("Term '%s' not found in index. No results for AND search.\n", token)
			return finalResults
		}
		postingLists[token] = postings
		if shortestPostingListLength == -1 || len(postings) < shortestPostingListLength {
			shortestPostingListLength = len(postings)
			shortestToken = token
		}
	}

	if shortestToken == "" {
		return finalResults
	}

	candidateDocs := make(map[int]map[string][]int) // DocID -> {token -> positions}

	for _, p := range postingLists[shortestToken] {
		candidateDocs[p.DocID] = map[string][]int{
			shortestToken: p.Positions,
		}
	}

	for token, postings := range postingLists {
		if token == shortestToken {
			continue
		}

		currentMatchingDocs := make(map[int]map[string][]int)
		for _, p := range postings {
			if existingMatchData, ok := candidateDocs[p.DocID]; ok {
				existingMatchData[token] = p.Positions
				currentMatchingDocs[p.DocID] = existingMatchData
			}
		}
		candidateDocs = currentMatchingDocs
		if len(candidateDocs) == 0 {
			fmt.Println("No documents match all query terms.")
			return finalResults
		}
	}

	for docID, termPositionsMap := range candidateDocs {
		doc, docExists := idx.Docs[docID]
		if docExists {
			finalResults = append(finalResults, SearchResult{
				Document:           doc,
				QueryTermPositions: termPositionsMap,
			})
		}
	}

	sort.Slice(finalResults, func(i, j int) bool {
		return finalResults[i].Document.ID < finalResults[j].Document.ID
	})

	return finalResults
}
