package main

import (
	"encoding/json"
	"fmt"
)

type StructData struct {
	S_id   int
	S_name string
	S_xx   interface{}
}

func BuildData() (
	struct_data StructData,
	map_data map[string]interface{},
	slice_data []map[string]interface{},
	array_data [3]StructData) {

	// build struct_data
	struct_data = StructData{
		S_id:   1,
		S_name: "asd",
		S_xx:   [...]int{1, 2, 3},
	}
	// build map_data
	A_map_column := make(map[string]string)
	B_map_column := make(map[string]string)
	A_map_column["A"] = "abc"
	B_map_column["C"] = "cba"
	propA, propB := BuildProperty(123, A_map_column, "abc",
		1234, B_map_column, "123")
	map_data = make(map[string]interface{})
	map_data["property_slices"] = []propertyA{propA}
	map_data["property_arr"] = [...]propertyB{propB}
	// build slice_data
	slice_map_data := make(map[string]interface{})
	slice_map_data["A"] = map_data
	slice_data = append(slice_data, slice_map_data)
	// build array_data
	array_data = [3]StructData{struct_data, struct_data, struct_data}

	return struct_data, map_data, slice_data, array_data
}

func serialize_data() ([]byte, []byte, []byte, []byte) {
	var (
		struct_data StructData
		map_data    map[string]interface{}
		slice_data  []map[string]interface{}
		array_data  [3]StructData
	)

	struct_data, map_data, slice_data, array_data = BuildData()

	serialize_struct_data, err := json.Marshal(struct_data)
	if err != nil {
		fmt.Println("序列化失败", err)
	}

	serialize_map_data, err := json.Marshal(map_data)
	if err != nil {
		fmt.Println("serialize_map_data,序列化失败", err)
	}

	serialize_slice_data, err := json.Marshal(slice_data)
	if err != nil {
		fmt.Println("序列化失败", err)
	}

	serialize_array_data, err := json.Marshal(array_data)
	if err != nil {
		fmt.Println("序列化失败", err)
	}

	return serialize_map_data, serialize_struct_data, serialize_slice_data, serialize_array_data
}

func main() {
	serialize_map_data, serialize_struct_data, serialize_slice_data, serialize_array_data := serialize_data()
	fmt.Printf("序列化后的结果是：\n")
	fmt.Printf("Map Data: %s\n", serialize_map_data)
	fmt.Printf("Struct Data: %s\n", serialize_struct_data)
	fmt.Printf("Slice Data: %s\n", serialize_slice_data)
	fmt.Printf("Array Data: %s\n", serialize_array_data)
}
