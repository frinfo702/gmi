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
type processdFileResult struct {
	filePath string
	tokens   []string
	err      error
}

func BuildIndex(rootDirPath string, idx *InvertedIndex) error {
	fmt.Printf("Starting to build index for directory: %s\n", rootDirPath)

	var filePaths []string
	err := filepath.WalkDir(rootDirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %q: %v\n", path, err)
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
		fmt.Println("No file found to index")
		return nil
	}
	fmt.Printf("Found %d files to process.\n", len(filePaths))

	numWorkers := runtime.NumCPU()
	if numWorkers > len(filePaths) {
		numWorkers = len(filePaths)
	}
	fmt.Printf("Using %d workers goroutines.\n", numWorkers)

	jobs := make(chan string, len(filePaths))
	results := make(chan processdFileResult, len(filePaths))

	var wg sync.WaitGroup

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			fmt.Printf("Workers %d started\n", workerID)
			for filePath := range jobs {
				content, err := os.ReadFile(filePath)
				if err != nil {
					results <- processdFileResult{filePath: filePath, err: fmt.Errorf("worker %d error reading file %q: %w", workerID, filePath, err)}
					continue
				}
				tokens := tokenizer.Tokenize(string(content))
				results <- processdFileResult{filePath: filePath, tokens: tokens, err: nil}
			}
			fmt.Printf("Worker %d finishd", workerID)
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
				fmt.Print("Error processing file %s: %v\n", result.filePath, result.err)
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
				// fmt.Printf("  File %s (DocID %d) will be updated.\n", result.filePath, docID)
			} else {
				docID = idx.NextDocID
				idx.Docs[docID] = Document{ID: docID, Path: result.filePath}
				idx.NextDocID++
			}
			addTokensToInvertedIndex(idx, docID, result.tokens)
			if processedCount%100 == 0 { // 100ファイル処理するごとに進捗表示
				fmt.Printf("Collected results for %d/%d files...\n", processedCount, len(filePaths))
			}
		}
		fmt.Println("All results collected.")
	}()

	wg.Wait()
	close(results)

	resultWg.Wait()

	fmt.Println("Index building process (concurrent file processing and index construction) completed.")
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
