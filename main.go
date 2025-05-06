package main

import (
	"flag"
	"fmt"
	"gmi/indexer"
	"gmi/searcher"
	"gmi/store"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "index":
		handleIndexCommand()
	case "search":
		handleSearchCommand()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: go_my_index <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  index -dir <target_directory> [-out <index_file_path>]")
	fmt.Println("  search -index <index_file_path> -q <query> [-mode <and|or>]")
}

func handleIndexCommand() {
	indexCmd := flag.NewFlagSet("index", flag.ExitOnError)
	targetDir := indexCmd.String("dir", "", "Directory to index (required)")
	indexPath := indexCmd.String("out", "myindex.idx", "Path to save/load the index file")
	indexCmd.Parse(os.Args[2:])

	if *targetDir == "" {
		fmt.Println("Error: -dir flag is required for index command.")
		indexCmd.Usage()
		os.Exit(1)
	}

	fmt.Printf("Index command: targetDir='%s', indexPath='%s'\n", *targetDir, *indexPath)

	oldIdx, err := store.LoadIndex(*indexPath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("Warning: Error loading existing index: %v. A new index will be built.\n", err)
		}
		if oldIdx == nil {
			oldIdx = indexer.NewInvertedIndex()
		}
	}

	newIdx, buildErr := indexer.BuildIndex(*targetDir, oldIdx)
	if buildErr != nil {
		fmt.Printf("Error building/updating index: %v\n", buildErr)
		os.Exit(1)
	}

	saveErr := store.SaveIndex(newIdx, *indexPath)
	if saveErr != nil {
		fmt.Printf("Error saving index: %v\n", saveErr)
		os.Exit(1)
	}
	fmt.Println("Index built/updated and saved successfully.")
}

func handleSearchCommand() {
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	indexPath := searchCmd.String("index", "myindex.idx", "Path to the index file")
	query := searchCmd.String("q", "", "Search query (required)")
	mode := searchCmd.String("mode", "and", "Search mode: 'and' or 'or' (default: 'and')")
	searchCmd.Parse(os.Args[2:])

	if *query == "" {
		fmt.Println("Error: -q flag is required for search command.")
		searchCmd.Usage()
		os.Exit(1)
	}
	normalizedMode := strings.ToLower(*mode)
	if normalizedMode != "and" && normalizedMode != "or" {
		fmt.Println("Error: Invalid search mode. Must be 'and' or 'or'.")
		searchCmd.Usage()
		os.Exit(1)
	}

	fmt.Printf("Search command: indexPath='%s', query='%s', mode='%s'\n", *indexPath, *query, normalizedMode)
	idx, err := store.LoadIndex(*indexPath)
	if err != nil {
		fmt.Printf("Error loading index for search: %v\n", err)
		os.Exit(1)
	}
	if len(idx.Docs) == 0 && idx.NextDocID == 0 {
		fmt.Println("The index is empty or not found. Please build the index first using the 'index' command.")
		return
	}

	searchResults := searcher.Search(idx, *query, normalizedMode)

	if len(searchResults) == 0 {
		fmt.Println("No documents found matching your query.")
		return
	}

	fmt.Printf("Found %d document(s) matching query (mode: %s):\n", len(searchResults), normalizedMode)
	for i, res := range searchResults {
		fmt.Printf("%d. File: %s (DocID: %d, Score: %.4f)\n", i+1, res.Document.Path, res.Document.ID, res.Score)

		var termDetails []string
		for term, positions := range res.QueryTermPositions {
			displayPositions := positions
			if len(displayPositions) > 3 {
				displayPositions = displayPositions[:3]
			}
			termDetails = append(termDetails, fmt.Sprintf("'%s' at %v", term, displayPositions))
		}
		fmt.Printf("   Terms found: %s (TotalWordsInDoc: %d)\n", strings.Join(termDetails, "; "), res.Document.TotalWords)
	}
}
