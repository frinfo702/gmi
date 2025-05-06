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

// BuildIndex は指定されたディレクトリを走査し、転置インデックスを構築します。
// (ToDo 6 の中核部分。今回はファイル走査、読み込み、トークナイズまでを行い、
//
//	インデックス構造体への格納は次のステップで行います。)
func BuildIndex(rootDirPath string, idx *InvertedIndex) error {
	fmt.Printf("Starting to build index for directory: %s\n", rootDirPath)

	err := filepath.WalkDir(rootDirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %q: %v\n", path, err)
			return err // エラーが発生したら探索を中断しない場合は `nil` を返す
		}
		if !d.IsDir() && (strings.HasSuffix(d.Name(), ".txt") || strings.HasSuffix(d.Name(), ".md")) {
			fmt.Printf("Processing file: %s\n", path)

			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("Error reading file %q: %v\n", path, err)
				return nil // このファイルはスキップして次に進む
			}

			// ドキュメントをインデックスに追加
			docID := idx.NextDocID
			idx.Docs[docID] = Document{ID: docID, Path: path}
			idx.NextDocID++

			tokens := tokenizer.Tokenize(string(content))
			fmt.Printf("  Tokens (%s): %v\n", path, tokens)

			// ここで tokens と docID を使って転置インデックス (idx.Index) を更新する (ToDo 6 の残り)
			// (この部分は次のステップで実装します)
			addTokensToInvertedIndex(idx, docID, tokens)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking the path %q: %w", rootDirPath, err)
	}

	fmt.Println("Index building process (file scan, tokenize) completed.")
	return nil
}

// addTokensToInvertedIndex はドキュメントのトークンリストを転置インデックスに追加します。
// (ToDo 6 の中核部分、この関数を次のステップで詳細に実装します)
func addTokensToInvertedIndex(idx *InvertedIndex, docID int, tokens []string) {
	tokenPositions := make(map[string][]int)
	for i, token := range tokens {
		tokenPositions[token] = append(tokenPositions[token], i)
	}

	for token, positions := range tokenPositions {
		// 既存のPostingリストを取得、または新規作成
		postings := idx.Index[token]

		// このドキュメントに関するPostingが存在するか確認
		found := false
		for i, p := range postings {
			if p.DocID == docID {
				// 通常、同じドキュメントを2回インデックスすることはないはずだが、念のため
				// もし既存ならポジションを追加することもできるが、ここでは単純に上書き（またはエラー）
				// 今回のロジックでは1ドキュメント1回処理なので、ここには到達しない想定
				postings[i].Positions = append(postings[i].Positions, positions...)
				found = true
				break
			}
		}

		if !found {
			newPosting := Posting{DocID: docID, Positions: positions}
			idx.Index[token] = append(postings, newPosting)
		}
	}
}
