package store

import (
	"encoding/gob"
	"fmt"
	"gmi/indexer"
	"os"
)

// SaveIndex は転置インデックスを指定されたファイルパスに保存します。
func SaveIndex(idx *indexer.InvertedIndex, filePath string) error {
	fmt.Printf("TODO: Implement SaveIndex to %s\n", filePath)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create index file %s: %w", filePath, err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(idx); err != nil {
		return fmt.Errorf("failed to encode index to file %s: %w", filePath, err)
	}
	fmt.Printf("Index saved to %s\n", filePath)
	return nil
}

// LoadIndex は指定されたファイルパスから転置インデックスを読み込みます。
func LoadIndex(filePath string) (*indexer.InvertedIndex, error) {
	fmt.Printf("TODO: Implement LoadIndex from %s\n", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Index file %s not found, creating new index.\n", filePath)
			return indexer.NewInvertedIndex(), nil
		}
		return nil, fmt.Errorf("failed to open index file %s: %w", filePath, err)
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	var idx indexer.InvertedIndex
	if err := decoder.Decode(&idx); err != nil {
		return nil, fmt.Errorf("failed to decode index from file %s: %w", filePath, err)
	}
	fmt.Printf("Index loaded from %s. NextDocID: %d, Index size: %d tokens, Docs: %d\n",
		filePath, idx.NextDocID, len(idx.Index), len(idx.Docs))

	// gobでデコードした際、mapがnilになる場合があるので初期化しておく
	if idx.Index == nil {
		idx.Index = make(map[string][]indexer.Posting)
	}
	if idx.Docs == nil {
		idx.Docs = make(map[int]indexer.Document)
	}

	return &idx, nil
}
