package searcher

import (
	"fmt"
	"gmi/indexer"
	"gmi/tokenizer"
	"sort"
	"strings"
)

// SearchResult は検索結果の1つのアイテムを表します。
type SearchResult struct {
	Document           indexer.Document
	QueryTermPositions map[string][]int // key: 検索クエリのトークン, value: そのトークンの出現位置リスト
}

// Searchは指定されたインデックス内でクエリに一致するドキュメントを検索します
func Search(idx *indexer.InvertedIndex, query string, mode string) []SearchResult {
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

	normalizedMode := strings.ToLower(mode)
	fmt.Printf("Searching for terms (%s): %v\n", normalizedMode, queryTokens)

	switch normalizedMode {
	case "or":
		finalResults = performOrSearch(idx, queryTokens)
	case "and":
		finalResults = performAndSearch(idx, queryTokens)
	default:
		fmt.Printf("Error: Unsupported search mode '%s'. Defaulting to AND search.\n", mode)
		finalResults = performAndSearch(idx, queryTokens) // 不明なモードの場合はANDとして扱う
	}

	sort.Slice(finalResults, func(i, j int) bool {
		return finalResults[i].Document.ID < finalResults[j].Document.ID
	})

	return finalResults
}

// performAndSearchはAND検索を実行します
func performAndSearch(idx *indexer.InvertedIndex, queryTokens []string) []SearchResult {
	var results []SearchResult
	postingLists := make(map[string][]indexer.Posting)
	shortestPostingListLength := -1
	var shortestToken string

	for _, token := range queryTokens {
		postings, found := idx.Index[token]
		if !found {
			fmt.Printf("Term '%s' not found in index. No results for AND search.\n", token)
			return results
		}
		postingLists[token] = postings
		if shortestPostingListLength == -1 || len(postings) < shortestPostingListLength {
			shortestPostingListLength = len(postings)
			shortestToken = token
		}
	}

	if shortestToken == "" {
		return results
	}

	candidateDocs := make(map[int]map[string][]int)
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
			fmt.Println("No documents match all query terms for AND search.")
			return results
		}
	}

	for docID, termPositionsMap := range candidateDocs {
		doc, docExists := idx.Docs[docID]
		if docExists {
			results = append(results, SearchResult{
				Document:           doc,
				QueryTermPositions: termPositionsMap,
			})
		}
	}
	return results
}

// performOrSearch はOR検索を実行します
func performOrSearch(idx *indexer.InvertedIndex, queryTokens []string) []SearchResult {
	var results []SearchResult
	// DocID -> {見つかったトークン -> positions}
	matchingDocs := make(map[int]map[string][]int)

	for _, token := range queryTokens {
		postings, found := idx.Index[token]
		if !found {
			fmt.Printf("Term '%s' not found in index, skipping for OR search.\n", token)
			continue
		}

		for _, p := range postings {
			if _, docExistsInMatch := matchingDocs[p.DocID]; !docExistsInMatch {
				matchingDocs[p.DocID] = make(map[string][]int)
			}
			matchingDocs[p.DocID][token] = p.Positions
		}
	}

	if len(matchingDocs) == 0 {
		fmt.Println("No documents match any query terms for OR search.")
		return results
	}

	for docID, termPositionsMap := range matchingDocs {
		doc, docExists := idx.Docs[docID]
		if docExists {
			results = append(results, SearchResult{
				Document:           doc,
				QueryTermPositions: termPositionsMap,
			})
		}
	}
	return results
}
