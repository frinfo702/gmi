package searcher

import "testing"

func Test_generateSnippet(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		docContent         string
		keywordToHighlight string
		positionsInDoc     []int
		want               string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateSnippet(tt.docContent, tt.keywordToHighlight, tt.positionsInDoc)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("generateSnippet() = %v, want %v", got, tt.want)
			}
		})
	}
}
