// go-my-index/indexer/indexer.go
package indexer

import (
	"fmt"
	"gmi/tokenizer"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// processFileResultはワーカーgoroutineからの処理結果を格納
type processedFileResult struct {
	filePath   string
	tokens     []string
	totalWords int
	err        error
}

func BuildIndex(rootDirPath string, idx *InvertedIndex) error {
	// ... (関数の前半部分は変更なし) ...
	fmt.Printf("Starting to build index for directory (concurrently): %s\n", rootDirPath)

	var filePaths []string
	err := filepath.WalkDir(rootDirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %q during WalkDir: %v\n", path, err)
			return err
		}
		lowerName := strings.ToLower(d.Name())
		if !d.IsDir() && (strings.HasSuffix(lowerName, ".txt") || strings.HasSuffix(lowerName, ".md")) {
			filePaths = append(filePaths, path)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking the path %q to gather files: %w", rootDirPath, err)
	}

	if len(filePaths) == 0 {
		fmt.Println("No files found to index.")
		return nil
	}
	fmt.Printf("Found %d files to process.\n", len(filePaths))

	numWorkers := runtime.NumCPU()
	if numWorkers > len(filePaths) {
		numWorkers = len(filePaths)
	}
	fmt.Printf("Using %d worker goroutines.\n", numWorkers)

	jobs := make(chan string, len(filePaths))
	results := make(chan processedFileResult, len(filePaths))

	var wg sync.WaitGroup

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for filePath := range jobs {
				content, err := os.ReadFile(filePath)
				if err != nil {
					results <- processedFileResult{filePath: filePath, err: fmt.Errorf("worker %d error reading file %q: %w", workerID, filePath, err)}
					continue
				}
				tokens := tokenizer.Tokenize(string(content))
				validTokensCount := 0
				for _, t := range tokens {
					if t != "" {
						validTokensCount++
					}
				}
				results <- processedFileResult{filePath: filePath, tokens: tokens, totalWords: validTokensCount, err: nil} 
			}
		}(w)
	}

	for _, fp := range filePaths {
		jobs <- fp
	}
	close(jobs)

	var resultWg sync.WaitGroup
	resultWg.Add(1)

	go func() {
		defer resultWg.Done()
		processedCount := 0
		for result := range results {
			processedCount++
			if result.err != nil {
				fmt.Printf("Error processing file %s: %v\n", result.filePath, result.err)
				continue
			}

			var docID int
			existingDocID := -1
			for id, doc := range idx.Docs {
				if doc.Path == result.filePath {
					existingDocID = id
					break
				}
			}

			if existingDocID != -1 {
				docID = existingDocID
				currentDoc := idx.Docs[docID]
				currentDoc.TotalWords = result.totalWords
				idx.Docs[docID] = currentDoc
			} else {
				docID = idx.NextDocID
				idx.Docs[docID] = Document{ID: docID, Path: result.filePath, TotalWords: result.totalWords} 
				idx.NextDocID++
			}
			addTokensToInvertedIndex(idx, docID, result.tokens) 
		}
	}()

	wg.Wait()
	close(results)
	resultWg.Wait()

	fmt.Println("Index building process completed.")
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
		frequency := len(positions) 
		postingsList := idx.Index[token]
		foundPostingForDoc := false
		for i, p := range postingsList {
			if p.DocID == docID {
				postingsList[i].Positions = positions
				postingsList[i].Frequency = frequency 
				foundPostingForDoc = true
				break
			}
		}

		if !foundPostingForDoc {
			newPosting := Posting{DocID: docID, Positions: positions, Frequency: frequency} 
			idx.Index[token] = append(postingsList, newPosting)
		}
	}
}
