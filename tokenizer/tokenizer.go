// go-my-index/tokenizer/tokenizer.go
package tokenizer

import (
	"regexp"
	"strings"
)

var (
	// 正規表現で単語として認識するパターン (英数字の連続)
	// より高度にするならUnicodeの文字クラスなどを考慮
	wordRegex = regexp.MustCompile(`[a-zA-Z0-9]+`)
)

// Tokenize は与えられたテキストを単語のリストに分割し、正規化します。
// 正規化処理として、小文字化を行います。
func Tokenize(text string) []string {
	words := wordRegex.FindAllString(text, -1)
	var tokens []string
	for _, word := range words {
		if word != "" {
			tokens = append(tokens, strings.ToLower(word))
		}
	}
	return tokens
}
