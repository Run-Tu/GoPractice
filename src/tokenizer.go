package main

import (
	"fmt"
	"log"
	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretrained"
)

func main() {
    // Download and cache pretrained tokenizer. In this case `bert-base-uncased` from Huggingface
    // can be any model with `tokenizer.json` available. E.g. `tiiuae/falcon-7b`
	configFile, err := tokenizer.CachedPath("flag-embedding", "tokenizer.json")
	if err != nil {
		panic(err)
	}

	tk, err := pretrained.FromFile(configFile)
	if err != nil {
		panic(err)
	}

	sentence := `这是一个测试用例`
	en, err := tk.EncodeSingle(sentence)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tokens: %q\n", en.Tokens)
	fmt.Printf("offsets: %v\n", en.Offsets)

	// Output
	// tokens: ["这" "是" "一" "个" "测" "试" "用" "例"]
	// offsets: [[0 3] [3 6] [6 9] [9 12] [12 15] [15 18] [18 21] [21 24]]
}