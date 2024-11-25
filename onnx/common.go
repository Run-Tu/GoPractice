package main

import (
	"log"

	"github.com/sugarme/tokenizer/pretrained"
)

func GetOnnxInput(sentence string) ([]int, []int, []int, error) {
	tk, err := pretrained.FromFile("../res/tokenizer.json")
	if err != nil {
		panic(err)
	}
	encoding, err := tk.EncodeSingle(sentence)
	if err != nil {
		log.Fatal(err)
		return nil, nil, nil, err
	}

	return encoding.Ids, encoding.AttentionMask, encoding.TypeIds, err
}

func convertToInt64Slice(data []int) []int64 {
	result := make([]int64, len(data))
	for i, v := range data {
		result[i] = int64(v)
	}
	return result
}
