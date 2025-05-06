// go-my-index/indexer/indexer.go
package indexer

import (
	"fmt"
	"gmi/tokenizer"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func BuildIndex(rootDirPath string, idx *InvertedIndex) error {
	fmt.Printf("Starting to build index for directory: %s\n", rootDirPath)

	err := filepath.WalkDir(rootDirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}
		// 対象ファイルの判定を改善 (大文字・小文字を区別しない)
		lowerName := strings.ToLower(d.Name())
		if !d.IsDir() && (strings.HasSuffix(lowerName, ".txt") || strings.HasSuffix(lowerName, ".md")) {
			fmt.Printf("Processing file: %s\n", path)

			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("Error reading file %q: %v\n", path, err)
				return nil // このファイルはスキップ
			}

			// 既にこのパスがインデックスされているかチェック (重複インデックス防止)
			// より厳密には、ファイルの内容が同じかどうかで判断すべきだが、ここではパスで簡易的に判断
			var existingDocID = -1
			for id, doc := range idx.Docs {
				if doc.Path == path {
					existingDocID = id
					break
				}
			}

			docID := idx.NextDocID
			if existingDocID != -1 {
				docID = existingDocID
				fmt.Printf("  File %s already partially indexed as DocID %d. Will update.\n", path, docID)
			} else {
				idx.Docs[docID] = Document{ID: docID, Path: path}
				idx.NextDocID++
			}

			tokens := tokenizer.Tokenize(string(content))
			// fmt.Printf("  Tokens (%s): %v\n", path, tokens) // ログが多いのでコメントアウト

			addTokensToInvertedIndex(idx, docID, tokens)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking the path %q: %w", rootDirPath, err)
	}

	fmt.Println("Index building process (file scan, tokenize, index construction) completed.")
	return nil
}

// addTokensToInvertedIndex はドキュメントのトークンリストを転置インデックスに追加します。
func addTokensToInvertedIndex(idx *InvertedIndex, docID int, tokens []string) {
	tokenPositionsInDoc := make(map[string][]int)
	for i, token := range tokens {
		if token == "" {
			continue
		}
		tokenPositionsInDoc[token] = append(tokenPositionsInDoc[token], i)
	}

	for token, positions := range tokenPositionsInDoc {
		postingsList := idx.Index[token]

		foundPostingForDoc := false
		for i, p := range postingsList {
			if p.DocID == docID {
				postingsList[i].Positions = positions
				foundPostingForDoc = true
				break
			}
		}

		if !foundPostingForDoc {
			newPosting := Posting{DocID: docID, Positions: positions}
			idx.Index[token] = append(postingsList, newPosting)
		}
	}
}
