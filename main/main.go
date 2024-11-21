package main

import (
	"fmt"
	utils "project01/package"
	"sync"
)

// 定义WaitGroup确保所有的goroutine都执行完毕
var wg sync.WaitGroup

func import_test() {
	var i, result, f_result = 1, 1, "2"
	result, f_result = utils.Add(i, 1)
	fmt.Println("result", result, "f_result", f_result)
}

func goroutine_test() {
	numbers := []int{1, 2, 3, 4, 5, 6}
	ch := make(chan int, len(numbers))

	for _, n := range numbers {
		wg.Add(1)
		// square()方法
		// go允许在func通过go func()的方式定义匿名函数
		go func(num int) {
			defer wg.Done()
			ch <- num * num
		}(n)
	}

	wg.Wait() // 等待所有goroutine完成
	close(ch) // 关闭channel

	for result := range ch {
		fmt.Println(result)
	}
}

func main() {
	// Part1 import
	import_test()
	// Part2 goroutine
	goroutine_test()

}