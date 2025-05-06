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
	fmt.Println("  search -index <index_file_path> -q <query>")
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

	idx, err := store.LoadIndex(*indexPath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("Error loading existing index: %v\n", err)
		fmt.Println("Attempting to create a new index.")
		idx = indexer.NewInvertedIndex()
	}

	err = indexer.BuildIndex(*targetDir, idx)
	if err != nil {
		fmt.Printf("Error building index: %v\n", err)
		os.Exit(1)
	}

	err = store.SaveIndex(idx, *indexPath)
	if err != nil {
		fmt.Printf("Error saving index: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Index built and saved successfully.")
}

func handleSearchCommand() {
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	indexPath := searchCmd.String("index", "myindex.idx", "Path to the index file")
	query := searchCmd.String("q", "", "Search query (required)")
	searchCmd.Parse(os.Args[2:])

	if *query == "" {
		fmt.Println("Error: -q flag is required for search command.")
		searchCmd.Usage()
		os.Exit(1)
	}

	fmt.Printf("Search command: indexPath='%s', query='%s'\n", *indexPath, *query)
	idx, err := store.LoadIndex(*indexPath)
	if err != nil {
		fmt.Printf("Error loading index for search: %v\n", err)
		os.Exit(1)
	}
	if len(idx.Docs) == 0 && idx.NextDocID == 0 {
		fmt.Println("The index is empty or not found. Please build the index first using the 'index' command.")
		return
	}

	searchResults := searcher.Search(idx, *query)

	if len(searchResults) == 0 {
		fmt.Println("No documents found matching your query.")
		return
	}

	fmt.Printf("Found %d document(s) matching all terms:\n", len(searchResults))
	for i, res := range searchResults {
		fmt.Printf("%d. File: %s (DocID: %d)\n", i+1, res.Document.Path, res.Document.ID)

		var termDetails []string
		for term, positions := range res.QueryTermPositions {
			displayPositions := positions
			if len(displayPositions) > 3 {
				displayPositions = displayPositions[:3]
			}
			termDetails = append(termDetails, fmt.Sprintf("'%s' at %v", term, displayPositions))
		}
		fmt.Printf("   Terms found: %s\n", strings.Join(termDetails, "; "))
	}
}
