package utils

import "strconv"

func Add(num1 int, num2 int) (int, string) {
	return num1 + num2, strconv.Itoa(num1 + num2)
}
