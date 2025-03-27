package infrastructure

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/frinfo702/fixer/domain/model"
)

// RustSearchAdapterは、Rust 製検索エンジンを呼び出すためのアダプタ
type RustSearchAdapter struct {
	binaryPath string
}

func NewRustSearchAdapter(binaryPath string) *RustSearchAdapter {
	log.Printf("RustSearchAdapter initialized with binary path: %s", binaryPath)
	return &RustSearchAdapter{binaryPath: binaryPath}
}

func (r *RustSearchAdapter) Search(query model.SearchQuery) ([]model.SearchResult, error) {
	// 検索パスを絶対パスに変換
	searchPath := query.Path
	if searchPath != "" && searchPath != "." {
		absPath, err := filepath.Abs(searchPath)
		if err != nil {
			log.Printf("Warning: Failed to convert path to absolute: %v", err)
		} else {
			log.Printf("Converting path '%s' to absolute: '%s'", searchPath, absPath)
			searchPath = absPath
		}
	}

	// 引数リストを構築
	args := []string{query.Query} // 最初にクエリを追加

	// 検索対象ディレクトリを追加 (指定されていれば)
	if searchPath != "" && searchPath != "." {
		args = append(args, "--dir", searchPath)
	}

	// ファジー検索フラグを追加
	if query.Fuzzy {
		args = append(args, "--fuzzy")
	}

	log.Printf("Executing Rust search: %s %v", r.binaryPath, args)
	cmd := exec.Command(r.binaryPath, args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		errMsg := strings.TrimSpace(errBuf.String())
		if errMsg == "" {
			errMsg = err.Error()
		}
		log.Printf("Rust search error: %s", errMsg)
		return nil, fmt.Errorf("rust search error: %s", errMsg)
	}

	output := outBuf.String()
	log.Printf("Rust search output received: %d bytes", len(output))
	var results []model.SearchResult
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "検索結果") || strings.HasPrefix(line, "スコア:") {
			continue
		}

		// 新しい出力形式に対応
		// 形式: "スコア: XX, ファイルパス:行番号: テキスト"
		if strings.HasPrefix(line, "スコア:") {
			parts := strings.SplitN(line, ",", 2)
			if len(parts) < 2 {
				continue
			}

			// "ファイルパス:行番号: テキスト" の部分を取得
			fileParts := strings.SplitN(strings.TrimSpace(parts[1]), ":", 3)
			if len(fileParts) < 3 {
				continue
			}

			var lineNum int
			fmt.Sscanf(fileParts[1], "%d", &lineNum)
			results = append(results, model.SearchResult{
				FilePath:   fileParts[0],
				LineNumber: lineNum,
				LineText:   strings.TrimSpace(fileParts[2]),
			})
			continue
		}

		// 旧形式のサポートも残す（互換性のため）
		parts := strings.SplitN(line, ":", 3)
		if len(parts) < 3 {
			continue
		}
		var lineNum int
		fmt.Sscanf(parts[1], "%d", &lineNum)
		results = append(results, model.SearchResult{
			FilePath:   parts[0],
			LineNumber: lineNum,
			LineText:   strings.TrimSpace(parts[2]),
		})
	}

	// 検索結果を関連性でソート
	// まず、クエリが行のテキストに含まれるかを確認し、含まれる場合は優先度を高くする
	sort.Slice(results, func(i, j int) bool {
		// クエリが含まれているかをチェック
		queryInI := strings.Contains(strings.ToLower(results[i].LineText), strings.ToLower(query.Query))
		queryInJ := strings.Contains(strings.ToLower(results[j].LineText), strings.ToLower(query.Query))

		// どちらも含まれている場合、文字列の長さが短い方が良い（より具体的なマッチ）
		if queryInI && queryInJ {
			return len(results[i].LineText) < len(results[j].LineText)
		}

		// クエリを含むものを優先
		return queryInI && !queryInJ
	})

	log.Printf("Processed search results: %d items", len(results))
	return results, nil
}
