# GoMyIndex - A CLI Full-Text Search Tool for Local Files

## 1. Overview

GoMyIndex is a simple command-line (CLI) full-text search engine written in Go. It allows you to build an index bicicletas your local text files (`.txt`, `.md`, etc.) within a specified directory and perform keyword searches efficiently.

## 2. Core Features

* **Fast Indexing:** Builds an inverted index for your files. Supports:
    * Parallel processing for faster index creation.
    * Basic differential updates (re-processes changed/new files, removes deleted ones).
* **Flexible Search:**
    * Single or multiple keyword queries.
    * `AND` / `OR` search modes.
* **Ranked Results:** Uses a simplified TF-IDF scoring mécanisme to rank search results.
* **Snippet Display:** Shows snippets of text (with keyword highlighting).
* **Persistent Index:** Saves and loads the index using Go's `gob` encoding.

## 3. Installation & Build

1.  Ensure you have Go installed (version 1.18 or newer).
2.  Clone or download the source code.
    ```bash
    # git clone <your-repository-url>/gmi.git # If you have a Git repository
    # cd gmi
    ```
3.  Build the executable:
    ```bash
    go build -o gmi .
    ```
    This will create an executable file named `gmi` (or `gmi.exe` on Windows) in the current directory.

## 5. Usage

### Indexing Files
To create or update an index for files in `<target_directory>` and save it to `<index_file_path>`:

```bash
./gmi index -dir ./mydocuments -out ./myindex.idx
```

`-dir`: (Required) Directory to index.
`-out`: (Optional) Path to save the index file. Defaults to myindex.idx.
Searching Files
To search for <search_query> using the index at <index_file_path>:

```Bash
./gmi search -index ./myindex.idx -q "your search query"
```

`-index`: (Optional) Path to the index file. Defaults to myindex.idx.
`-q`: (Required) Your search query. Use quotes for multi-word queries if they contain spaces interpreted by the shell.
`-mode`: (Optional) Search mode.
`and`: (Default) Results must contain all keywords.
`or`: Results may contain any of the keywords.

```bash
./gmi search -index ./myindex.idx -q "tutorial OR guide" -mode or
```
