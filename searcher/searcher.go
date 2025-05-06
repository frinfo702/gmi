package searcher

import (
	"fmt"
	"gmi/indexer"
	"gmi/tokenizer"
	"math"
	"sort"
	"strings"
)

// SearchResult は検索結果の1つのアイテムを表します。
type SearchResult struct {
	Document           indexer.Document
	QueryTermPositions map[string][]int // key: 検索クエリのトークン, value: そのトークンの出現位置リスト
	Score              float64          // TF-IDFスコア
}

// calculateIDF calculates the Inverse Document Frequency for a term.
func calculateIDF(totalDocuments int, docsContainingTerm int) float64 {
	if docsContainingTerm == 0 {
		return 0
	}
	// IDF(t) = log(N / df_t) -- N:総文書数, df_t:単語tを含む文書数
	return math.Log(float64(totalDocuments) / float64(docsContainingTerm))
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

	totalDocsInIndex := len(idx.Docs)
	idfScores := make(map[string]float64)
	for _, token := range queryTokens {
		if _, exists := idfScores[token]; !exists {
			postingsForToken, foundInIndex := idx.Index[token]
			if foundInIndex {
				idfScores[token] = calculateIDF(totalDocsInIndex, len(postingsForToken))
			} else {
				idfScores[token] = 0
			}
		}
	}

	// DocID -> SearchResult (スコア計算途中)
	intermediateResults := make(map[int]map[string]indexer.Posting)

	switch normalizedMode {
	case "or":
		for _, token := range queryTokens {
			postings, found := idx.Index[token]
			if !found {
				continue
			}
			for _, p := range postings {
				if _, docExists := intermediateResults[p.DocID]; !docExists {
					intermediateResults[p.DocID] = make(map[string]indexer.Posting)
				}
				intermediateResults[p.DocID][token] = p
			}
		}
	case "and":
		postingLists := make(map[string][]indexer.Posting)
		shortestPostingListLength := -1
		var shortestToken string

		for _, token := range queryTokens {
			postings, found := idx.Index[token]
			if !found {
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

		initialCandidates := make(map[int]map[string]indexer.Posting)
		for _, p := range postingLists[shortestToken] {
			initialCandidates[p.DocID] = map[string]indexer.Posting{
				shortestToken: p,
			}
		}

		for token, postings := range postingLists {
			if token == shortestToken {
				continue
			}
			nextCandidates := make(map[int]map[string]indexer.Posting)
			for _, p := range postings {
				if existingData, ok := initialCandidates[p.DocID]; ok {
					for t, post := range existingData {
						if _, exists := nextCandidates[p.DocID]; !exists {
							nextCandidates[p.DocID] = make(map[string]indexer.Posting)
						}
						nextCandidates[p.DocID][t] = post
					}
					nextCandidates[p.DocID][token] = p
				}
			}
			initialCandidates = nextCandidates
			if len(initialCandidates) == 0 {
				return finalResults
			}
		}
		intermediateResults = initialCandidates
	default:
		fmt.Printf("Error: Unsupported search mode '%s'.\n", mode)
		return finalResults
	}

	if len(intermediateResults) == 0 {
		return finalResults
	}

	for docID, termPostingMap := range intermediateResults {
		doc, docExists := idx.Docs[docID]
		if !docExists {
			continue
		}

		currentDocScore := 0.0
		queryTermPositions := make(map[string][]int)

		for queryToken, posting := range termPostingMap {
			tf := float64(posting.Frequency)
			idf := idfScores[queryToken]
			currentDocScore += tf * idf
			queryTermPositions[queryToken] = posting.Positions
		}

		finalResults = append(finalResults, SearchResult{
			Document:           doc,
			QueryTermPositions: queryTermPositions,
			Score:              currentDocScore,
		})
	}

	sort.Slice(finalResults, func(i, j int) bool {
		return finalResults[i].Score > finalResults[j].Score
	})

	return finalResults
}
