package main

import (
	"fmt"
	utils "project01/package"
)

func main() {
	i, result := 1, 1
	var f_result string
	result, f_result = utils.Add(i, 1)
	fmt.Println("result", result, "f_result", f_result)
}
