// go-my-index/indexer/indexer.go
package indexer

import (
	"fmt"
	"gmi/tokenizer"
	"gmi/ui"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// processFileResultはワーカーgoroutineからの処理結果を格納
type processedFileResult struct {
	filePath     string
	tokens       []string
	totalWords   int
	lastModified time.Time
	err          error
}

func BuildIndex(rootDirPath string, oldIdx *InvertedIndex) (*InvertedIndex, error) {
	fmt.Printf("%s Starting to build/update index for: %s\n", ui.Cyan("▶"), rootDirPath)

	newIdx := NewInvertedIndex()
	if oldIdx != nil && oldIdx.NextDocID > 0 {
		newIdx.NextDocID = oldIdx.NextDocID
	}

	currentFileSystemFiles := make(map[string]fs.FileInfo) // path -> FileInfo
	err := filepath.WalkDir(rootDirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("%s accessing path %q during WalkDir: %v\n", ui.Yellow("Warning:"), path, err)
			return err
		}
		lowerName := strings.ToLower(d.Name())
		if !d.IsDir() && (strings.HasSuffix(lowerName, ".txt") || strings.HasSuffix(lowerName, ".md")) {
			info, statErr := d.Info()
			if statErr != nil {
				fmt.Printf("%s getting FileInfo for %s: %v\n", ui.Yellow("Warning:"), path, statErr)
				return nil
			}
			currentFileSystemFiles[path] = info
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s walking the path %q to gather files: %w", "error", rootDirPath, err)
	}

	if len(currentFileSystemFiles) == 0 {
		fmt.Println(ui.Yellow("No files found in the target directory. Returning an empty index."))
		return NewInvertedIndex(), nil
	}
	fmt.Printf("%s Found %d files in current file system.\n", ui.Cyan("ℹ"), len(currentFileSystemFiles))

	var filesToProcess []string
	oldDocsByPath := make(map[string]Document)
	if oldIdx != nil {
		for _, doc := range oldIdx.Docs {
			oldDocsByPath[doc.Path] = doc
		}
	}

	for path, fileInfo := range currentFileSystemFiles {
		oldDoc, existsInOldIndex := oldDocsByPath[path]
		if existsInOldIndex && oldDoc.LastModified.Equal(fileInfo.ModTime()) {
			newIdx.Docs[oldDoc.ID] = oldDoc
			filesToProcess = append(filesToProcess, path)
		} else {
			if existsInOldIndex {
				fmt.Printf("%s File %s changed (OldTime: %s, NewTime: %s).\n", ui.Yellow("↺"), path, oldDoc.LastModified, fileInfo.ModTime())
			} else {
				fmt.Printf("%s New file %s found.\n", ui.Green("+"), path)
			}
			filesToProcess = append(filesToProcess, path)
		}
	}

	if oldIdx != nil {
		for path := range oldDocsByPath {
			if _, existsInCurrentFS := currentFileSystemFiles[path]; !existsInCurrentFS {
				fmt.Printf("%s File %s was deleted.\n", ui.Yellow("-"), path)
			}
		}
	}

	if len(filesToProcess) == 0 {
		fmt.Println("No files to process (all files unchanged or directory empty). Returning old index (or new if old was nil).")
		if oldIdx != nil && len(currentFileSystemFiles) > 0 {
			return oldIdx, nil
		}
		return newIdx, nil
	}
	fmt.Printf("%s %d files will be (re)processed.\n", ui.Cyan("▶"), len(filesToProcess))

	numWorkers := runtime.NumCPU()
	if numWorkers > len(filesToProcess) {
		numWorkers = len(filesToProcess)
	}

	jobs := make(chan string, len(filesToProcess))
	results := make(chan processedFileResult, len(filesToProcess))
	var wg sync.WaitGroup

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for filePath := range jobs {
				fileInfo := currentFileSystemFiles[filePath]
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
				results <- processedFileResult{filePath: filePath, tokens: tokens, totalWords: validTokensCount, lastModified: fileInfo.ModTime(), err: nil}
			}
		}(w)
	}

	for _, fp := range filesToProcess {
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
				fmt.Printf("%s processing file %s: %v\n", ui.Yellow("Warning:"), result.filePath, result.err)
				continue
			}

			var docID int
			oldDoc, pathExistedInOld := oldDocsByPath[result.filePath]
			if pathExistedInOld {
				docID = oldDoc.ID
				newIdx.Docs[docID] = Document{ID: docID, Path: result.filePath, TotalWords: result.totalWords, LastModified: result.lastModified}
			} else {
				docID = newIdx.NextDocID
				newIdx.Docs[docID] = Document{ID: docID, Path: result.filePath, TotalWords: result.totalWords, LastModified: result.lastModified}
				newIdx.NextDocID++
			}

			addTokensToInvertedIndex(newIdx, docID, result.tokens)
		}
	}()

	wg.Wait()
	close(results)
	resultWg.Wait()

	fmt.Println(ui.Green("Index update process completed."))
	return newIdx, nil
}

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
