package infrastructure

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/frinfo702/rustysearch/domain/model"
)

// RustSearchAdapter は、Rust 製検索エンジンを呼び出すためのアダプタ
type RustSearchAdapter struct {
	binaryPath string
}

func NewRustSearchAdapter(binaryPath string) *RustSearchAdapter {
	return &RustSearchAdapter{binaryPath: binaryPath}
}

func (r *RustSearchAdapter) Search(query model.SearchQuery) ([]model.SearchResult, error) {
	// 引数の構築：パス、検索クエリ、オプションで --fuzzy を付与
	args := []string{query.Path, query.Query}
	if query.Fuzzy {
		args = append(args, "--fuzzy")
	}

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
		return nil, fmt.Errorf("Rust search error: %s", errMsg)
	}

	output := outBuf.String()
	var results []model.SearchResult
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
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
	return results, nil
}
