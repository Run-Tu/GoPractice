// 包之间的引用关系,common和inference都属于onnx包下，inference应该可以直接用common中的方法吧
// inference中的outputShape放在yaml里是否合适？如果合适请从yaml文件中读取outputShape
// inference.go和session.go如何共用一个Session？需要提前执行InitONNX()方法
package onnx

import (
	"fmt"
	"sync"

	ort "github.com/yalue/onnxruntime_go"
)

var sessionMutex sync.Mutex

func GetEmbedding(text string) ([]float32, error) {
	InitONNX()

	// tokenizer获取模型输入
	inputIds, attentionMask, tokenTypeIDs, err := GetOnnxInput(text)
	if err != nil {
		return nil, fmt.Errorf("获取模型输入失败: %v", err)
	}

	inputIds64 := GetInt64Slice(inputIds)
	attentionMask64 := GetInt64Slice(attentionMask)
	tokenTypeIDs64 := GetInt64Slice(tokenTypeIDs)

	inputIdsTensor, err := ort.NewTensor[int64]([]int64{1, int64(len(inputIds64))}, inputIds64)
	if err != nil {
		return nil, fmt.Errorf("failed to create inputIdsTensor: %v", err)
	}
	defer inputIdsTensor.Destroy()
	attentionMaskTensor, err := ort.NewTensor[int64]([]int64{1, int64(len(attentionMask64))}, attentionMask64)
	if err != nil {
		return nil, fmt.Errorf("failed to create attentionMaskTensor: %v", err)
	}
	defer attentionMaskTensor.Destroy()
	tokenTypeTensor, err := ort.NewTensor[int64]([]int64{1, int64(len(tokenTypeIDs64))}, tokenTypeIDs64)
	if err != nil {
		return nil, fmt.Errorf("failed to create tokenTypeTensor: %v", err)
	}
	defer tokenTypeTensor.Destroy()

	// output
	outputShape := []int64{1, 1024}
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		return nil, fmt.Errorf("failde to create outputTensor: %v", err)
	}
	defer outputTensor.Destroy()

	// convert tensor to onnxruntime.Value
	inputValues := []ort.Value{inputIdsTensor, attentionMaskTensor, tokenTypeTensor}
	outputValues := []ort.Value{outputTensor}

	// 线程安全(ort的Session是否自动支持线程安全？)
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	if err := Session.Run(inputValues, outputValues); err != nil {
		return nil, fmt.Errorf("运行推理失败: %v", err)
	}
	if len(outputValues) == 0 {
		return nil, fmt.Errorf("获取输出数据失败")
	}

	return outputTensor.GetData(), nil
}
