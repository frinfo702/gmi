package main

import (
	"flag"
	"fmt"
	"gmi/indexer"
	"gmi/searcher"
	"gmi/store"
	"gmi/ui"
	"os"
	"sort"
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
		fmt.Printf("%s Unknown command: %s\n", ui.Yellow("!"), ui.Red(command))
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(ui.Bold("Usage:"), "go_my_index <command> [arguments]")
	fmt.Println(ui.Bold("Commands:"))
	fmt.Println("  ", ui.Cyan("index"), "-dir <target_directory> [-out <index_file_path>]")
	fmt.Println("  ", ui.Cyan("search"), "-index <index_file_path> -q <query> [-mode <and|or>]")
}

func handleIndexCommand() {
	indexCmd := flag.NewFlagSet("index", flag.ExitOnError)
	targetDir := indexCmd.String("dir", "", "Directory to index (required)")
	indexPath := indexCmd.String("out", "myindex.idx", "Path to save/load the index file")
	indexCmd.Parse(os.Args[2:])

	if *targetDir == "" {
		fmt.Println(ui.Red("Error:"), "-dir flag is required for index command.")
		indexCmd.Usage()
		os.Exit(1)
	}

	fmt.Printf("%s Index command: targetDir='%s', indexPath='%s'\n", ui.Cyan("▶"), *targetDir, *indexPath)

	oldIdx, err := store.LoadIndex(*indexPath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("%s Error loading existing index: %v. A new index will be built.\n", ui.Yellow("Warning:"), err)
		}
		if oldIdx == nil {
			oldIdx = indexer.NewInvertedIndex()
		}
	}

	newIdx, buildErr := indexer.BuildIndex(*targetDir, oldIdx)
	if buildErr != nil {
		fmt.Printf("%s %v\n", ui.Red("Error building/updating index:"), buildErr)
		os.Exit(1)
	}

	saveErr := store.SaveIndex(newIdx, *indexPath)
	if saveErr != nil {
		fmt.Printf("%s %v\n", ui.Red("Error saving index:"), saveErr)
		os.Exit(1)
	}
	fmt.Println(ui.Green("Index built/updated and saved successfully."))
}

func handleSearchCommand() {
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	indexPath := searchCmd.String("index", "myindex.idx", "Path to the index file")
	query := searchCmd.String("q", "", "Search query (required)")
	mode := searchCmd.String("mode", "and", "Search mode: 'and' or 'or' (default: 'and')")
	searchCmd.Parse(os.Args[2:])

	if *query == "" {
		fmt.Println(ui.Red("Error:"), "-q flag is required for search command.")
		searchCmd.Usage()
		os.Exit(1)
	}
	normalizedMode := strings.ToLower(*mode)
	if normalizedMode != "and" && normalizedMode != "or" {
		fmt.Println(ui.Red("Error:"), "Invalid search mode. Must be 'and' or 'or'.")
		searchCmd.Usage()
		os.Exit(1)
	}

	fmt.Printf("%s Search command: indexPath='%s', query='%s', mode='%s'\n", ui.Cyan("▶"), *indexPath, *query, normalizedMode)
	idx, err := store.LoadIndex(*indexPath)
	if err != nil {
		fmt.Printf("%s %v\n", ui.Red("Error loading index for search:"), err)
		os.Exit(1)
	}
	if len(idx.Docs) == 0 && idx.NextDocID == 0 {
		fmt.Println(ui.Yellow("The index is empty or not found. Please build the index first using the 'index' command."))
		return
	}

	searchResults := searcher.Search(idx, *query, normalizedMode)

	if len(searchResults) == 0 {
		fmt.Println(ui.Yellow("No documents found matching your query."))
		return
	}

	fmt.Printf("%s Found %d document(s) matching query (mode: %s):\n", ui.Green("✔"), len(searchResults), normalizedMode)
	for i, res := range searchResults {
		fmt.Printf("%d. File: %s (DocID: %d, Score: %.4f)\n", i+1, res.Document.Path, res.Document.ID, res.Score)

		var termDetails []string
		foundTermsInDoc := []string{}
		for term := range res.QueryTermPositions {
			foundTermsInDoc = append(foundTermsInDoc, term)
		}
		sort.Strings(foundTermsInDoc)

		for _, term := range foundTermsInDoc {
			positions := res.QueryTermPositions[term]
			displayPositions := positions
			if len(displayPositions) > 3 {
				displayPositions = displayPositions[:3]
			}
			termDetails = append(termDetails, fmt.Sprintf("'%s' at %v", term, displayPositions))
		}
		if len(termDetails) > 0 {
			fmt.Printf("   Terms: %s (TotalWordsInDoc: %d)\n", strings.Join(termDetails, "; "), res.Document.TotalWords)
		} else {
			fmt.Printf("   (No specific term positions for this combined result, TotalWordsInDoc: %d)\n", res.Document.TotalWords)
		}

		if len(res.Snippets) > 0 {
			for _, snippet := range res.Snippets {
				fmt.Printf("   %s %s\n", ui.Cyan("Snippet:"), snippet)
			}
		} else {
			fmt.Println("   ", ui.Dim("Snippet: [Not available]"))
		}
		if i < len(searchResults)-1 {
			fmt.Println("   ---")
		}
	}
}
