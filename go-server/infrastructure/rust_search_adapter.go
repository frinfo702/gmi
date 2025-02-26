package infrastructure

// RustSearchAdapterはRust製検索エンジンを呼び出すアダプタ
type RustSearchAdapter struct {
	binaryPath string
}

func NewRustSearchAdapter(binaryPath string) *RustSearchAdapter {
	return &RustSearchAdapter{binaryPath: binaryPath}
}
