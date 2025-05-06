package searcher

import (
	"fmt"
	"gmi/indexer"
	"gmi/tokenizer"
	"math"
	"os"
	"regexp"
	"sort"
	"strings"
)

// SearchResult は検索結果の1つのアイテムを表します。
type SearchResult struct {
	Document           indexer.Document
	QueryTermPositions map[string][]int // key: 検索クエリのトークン, value: そのトークンの出現位置リスト
	Score              float64          // TF-IDFスコア
	Snippets           []string         // キーワード周辺のスニペット
}

const (
	snippetContextWords = 5 // スニペットでキーワードの前後に表示する単語数
	maxSnippetsPerDoc   = 2 // 1ドキュメントあたり表示するスニペットの最大数
)

func generateSnippet(docContent string, keywordToHighlight string, positionsInDoc []int, contextWords int) string {
	if len(positionsInDoc) == 0 {
		return ""
	}

	words := strings.Fields(docContent) // strings.Fieldsは空白文字で分割
	if len(words) == 0 {
		return ""
	}

	re, err := regexp.Compile(`(?i)\b` + regexp.QuoteMeta(keywordToHighlight) + `\b`)
	if err != nil {
		return "[Error compiling regex for snippet]"
	}

	matches := re.FindAllStringIndex(docContent, -1) // 全てのマッチ位置(バイトオフセット)
	if len(matches) == 0 {
		return "[Keyword not found in content for snippet]" // 理論上ここには来ないはず
	}

	firstMatchStart := matches[0][0]
	firstMatchEnd := matches[0][1]
	snippetWindowChars := 40
	startOffset := firstMatchStart - snippetWindowChars
	if startOffset < 0 {
		startOffset = 0
	}
	endOffset := firstMatchEnd + snippetWindowChars
	if endOffset > len(docContent) {
		endOffset = len(docContent)
	}
	rawSnippet := docContent[startOffset:endOffset]
	highlightedSnippet := re.ReplaceAllString(rawSnippet, `**$0**`)

	prefix := ""
	if startOffset > 0 {
		prefix = "... "
	}
	suffix := ""
	if endOffset < len(docContent) {
		suffix = " ..."
	}

	return prefix + highlightedSnippet + suffix
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
	uniqueQueryTokens := make(map[string]bool)
	for _, token := range queryTokens {
		uniqueQueryTokens[token] = true
	}
	for token := range uniqueQueryTokens {
		postingsForToken, foundInIndex := idx.Index[token]
		if foundInIndex {
			idfScores[token] = calculateIDF(totalDocsInIndex, len(postingsForToken))
		} else {
			idfScores[token] = 0
		}
	}

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
		validQueryTokensForAND := []string{}

		for _, token := range queryTokens {
			postings, found := idx.Index[token]
			if !found {
				fmt.Printf("Term '%s' not found in index. Cannot satisfy AND condition.\n", token)
				return finalResults
			}
			postingLists[token] = postings
			validQueryTokensForAND = append(validQueryTokensForAND, token)
			if shortestPostingListLength == -1 || len(postings) < shortestPostingListLength {
				shortestPostingListLength = len(postings)
				shortestToken = token
			}
		}
		if shortestToken == "" {
			return finalResults
		} // 有効なトークンが一つもなかった

		currentCandidates := make(map[int]map[string]indexer.Posting)
		for _, p := range postingLists[shortestToken] {
			currentCandidates[p.DocID] = map[string]indexer.Posting{
				shortestToken: p,
			}
		}

		for _, token := range validQueryTokensForAND {
			if token == shortestToken {
				continue
			}
			nextCandidates := make(map[int]map[string]indexer.Posting)
			for _, p := range postingLists[token] {
				if existingData, ok := currentCandidates[p.DocID]; ok {
					// 既存のデータに現在のトークンのPostingを追加
					newData := make(map[string]indexer.Posting)
					for t, post := range existingData {
						newData[t] = post
					}
					newData[token] = p
					nextCandidates[p.DocID] = newData
				}
			}
			currentCandidates = nextCandidates
			if len(currentCandidates) == 0 {
				return finalResults
			}
		}
		intermediateResults = currentCandidates
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
		queryTermPositionsForThisDoc := make(map[string][]int)

		for queryToken, posting := range termPostingMap {
			tf := float64(posting.Frequency)
			idf := idfScores[queryToken]
			currentDocScore += tf * idf
			queryTermPositionsForThisDoc[queryToken] = posting.Positions
		}

		// スニペット生成
		var snippets []string
		docContentBytes, err := os.ReadFile(doc.Path)
		if err != nil {
			fmt.Printf("Warning: Could not read file %s to generate snippet: %v\n", doc.Path, err)
			snippets = append(snippets, "[Could not load content for snippet]")
		} else {
			docContent := string(docContentBytes)
			generatedSnippetsCount := 0
			for term := range termPostingMap {
				if generatedSnippetsCount >= maxSnippetsPerDoc {
					break
				}
				snippet := generateSnippet(docContent, term, termPostingMap[term].Positions, snippetContextWords)
				if snippet != "" {
					snippets = append(snippets, snippet)
					generatedSnippetsCount++
				}
			}
			if len(snippets) == 0 {
				limit := 100
				if len(docContent) < limit {
					limit = len(docContent)
				}
				snippets = append(snippets, strings.TrimSpace(docContent[:limit])+"...")
			}
		}

		finalResults = append(finalResults, SearchResult{
			Document:           doc,
			QueryTermPositions: queryTermPositionsForThisDoc,
			Score:              currentDocScore,
			Snippets:           snippets,
		})
	}

	sort.Slice(finalResults, func(i, j int) bool {
		return finalResults[i].Score > finalResults[j].Score
	})

	return finalResults
}
