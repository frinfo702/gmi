package indexer

// Document は検索対象のドキュメントを表します。
type Document struct {
	ID   int    // ドキュメントの一意なID
	Path string // ドキュメントのファイルパス
}

// Posting は転置インデックスのポスティングリストの要素です。
type Posting struct {
	DocID     int   // ドキュメントID
	Positions []int // 単語の出現位置 (ドキュメント内のトークンindex)
}

// InvertedIndex は転置インデックス全体を表します。
type InvertedIndex struct {
	Index     map[string][]Posting
	Docs      map[int]Document // ドキュメントIDからドキュメント情報へのマップ
	NextDocID int              // 次に割り当てるドキュメントID
}

// NewInvertedIndex は新しいInvertedIndexのインスタンスを作成します。
func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		Index:     make(map[string][]Posting),
		Docs:      make(map[int]Document),
		NextDocID: 0,
	}
}
