package main

import (
	"fmt"
	"testing"
)

// TestBasicTypes 基本类型
func TestBasicTypes(t *testing.T) {
	// 整型
	var a int = 10
	var b int = 3
	var sum int = a + b
	expectedSum := 14
	if sum != expectedSum {
		t.Errorf("Expected sum to be %d but got %d", expectedSum, sum)
	}

	// 浮点型
	var x float64 = 3.14
	var y float32 = 2.71
	result := x * float64(y)
	expectedResult := 8.6080 // go中的小数默认是float64类型,所以result和expectedResult都要转成float64类型才能比较
	if result != expectedResult {
		t.Errorf("Expected result to be %f but got %f", expectedResult, result)
	}

	// 布尔型
	var flag bool = true
	if !flag {
		t.Error("Expected flag to be true but got false")
	}

	// 字符串
	var s1 string = "Hello"
	var s2 string = "World"
	var s3 = s1 + " " + s2
	expectedString := "Hello World"
	if s3 != expectedString {
		t.Errorf("Expected string to be %s but got %s", expectedString, s3)
	}
}

/*
	数组
	切片
	map
	Struct
*/
// TestArray 测试数组类型
func TestArray(t *testing.T) {
	// 定义数组
	var arr [3]int = [3]int{1, 2, 3}
	expectedLength := 3
	if len(arr) != expectedLength {
		t.Errorf("Expected array length to be %d but got %d", expectedLength, len(arr))
	}

	// 测试数组内容
	for i, v := range arr {
		fmt.Println("i, v pair is", i, v)
		expectedValue := i + 1
		if v != expectedValue {
			t.Errorf("Expected value at index %d to be %d but got %d", i, expectedValue, v)
		}
	}
}

// TestSlice 测试切片类型
func TestSlice(t *testing.T) {
	//定义切片
	var slice []int = []int{1, 2, 3}
	slice = append(slice, 4, 5)
	expectedSlice := []int{1, 2, 3, 4, 5}
	//测试切片内容
	for i, v := range slice {
		if v != expectedSlice[i] {
			t.Errorf("Expected value at index %d to be %d but got %d", i, expectedSlice[i], v)
		}
	}
	//测试切片长度和容量
	expectedLength := 5
	expectedCapacity := 6
	fmt.Println("slice real length is", len(slice))
	fmt.Println("slice real capacity is", cap(slice))
	if len(slice) != expectedLength {
		t.Errorf("Expected length to be %d but got %d", expectedLength, len(slice))
	}
	if cap(slice) != expectedCapacity {
		t.Errorf("Expected capacity to be %d but got %d", expectedCapacity, cap(slice))
	}
}

// TestStruct 测试结构体类型
func TestStruct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	// 定义并初始化结构体
	p := Person{Name: "Alice", Age: 25}
	expectedName := "Alice"
	expectedAge := 25
	if p.Name != expectedName {
		t.Errorf("Expected Name to be %s but got %s", expectedName, p.Name)
	}
	if p.Age != expectedAge {
		t.Errorf("Expected Age to be %d but got %d", expectedAge, p.Age)
	}
}

// TestMap 测试映射类型
func TestMap(t *testing.T) {
	m := map[string]int{
		"age":    30,
		"height": 180,
	}
	if m["age"] != 30 {
		t.Errorf("Expected 'age' to be %d but got %d", 30, m["age"])
	}
	// 删除键
	delete(m, "height")
	if _, ok := m["height"]; ok {
		t.Error("Expected 'height' to be deleted but it's still present")
	}
}