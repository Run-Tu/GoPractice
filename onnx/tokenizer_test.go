package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/sugarme/tokenizer/pretrained"
)

func TestTokenizer(t *testing.T) {
	// Download and cache pretrained tokenizer. In this case `bert-base-uncased` from Huggingface
	// can be any model with `tokenizer.json` available. E.g. `tiiuae/falcon-7b`
	tk, err := pretrained.FromFile("../res/tokenizer.json")
	if err != nil {
		panic(err)
	}

	sentence := `这是一个测试用例`
	encoding, err := tk.EncodeSingle(sentence)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tokens: %q\n", encoding.Tokens)
	fmt.Printf("offsets: %v\n", encoding.Offsets)

	// 获取input_ids
	inputIDs := encoding.Ids
	fmt.Println("Input IDs:", inputIDs)
	// 获取attention_mask
	attentionMask := encoding.AttentionMask
	fmt.Println("Attention Mask:", attentionMask)
	// 获取token_type_ids
	tokenTypeIDs := encoding.TypeIds
	fmt.Println("Token Type IDs:", tokenTypeIDs)
}
