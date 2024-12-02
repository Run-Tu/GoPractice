// CUDA的options设置,如何知道我使用的是哪张卡？
package onnx

import (
	"fmt"
	"log"
	"testing"
	"time"

	ort "github.com/yalue/onnxruntime_go"
)

func TestOnnxGPU(t *testing.T) {
	// 初始化 ONNX Runtime环境
	err := ort.InitializeEnvironment()
	if err != nil {
		log.Fatalf("Failed to initialize ONNX Runtime: %v", err)
	}
	defer ort.DestroyEnvironment()

	// 获取模型的输入和输出信息
	inputs, outputs, err := ort.GetInputOutputInfo("../res/flag_embedding_model.onnx")
	if err != nil {
		log.Fatalf("获取模型输入输出信息失败: %v", err)
	}

	// 提取输入和输出的名称
	inputNames := make([]string, len(inputs))
	for i, input := range inputs {
		inputNames[i] = input.Name
	}
	outputNames := make([]string, len(outputs))
	for i, output := range outputs {
		outputNames[i] = output.Name
	}

	// 输入数据
	inputIds, attentionMask, tokenTypeIDs, err := GetOnnxInput(`这是一个测试案例`)
	if err != nil {
		log.Fatalf("获取模型input失败: %v", err)
	}
	inputIds64 := GetInt64Slice(inputIds)
	attentionMask64 := GetInt64Slice(attentionMask)
	tokenTypeIDs64 := GetInt64Slice(tokenTypeIDs)
	// convert to tensor
	inputIdsTensor, err := ort.NewTensor[int64]([]int64{1, int64(len(inputIds64))}, inputIds64)
	if err != nil {
		log.Fatalf("创建inputIds张量失败: %v", err)
	}
	defer inputIdsTensor.Destroy()

	attentionMaskTensor, err := ort.NewTensor[int64]([]int64{1, int64(len(attentionMask64))}, attentionMask64)
	if err != nil {
		log.Fatalf("创建attentionMask张量失败: %v", err)
	}
	defer attentionMaskTensor.Destroy()

	tokenTypeIDsTensor, err := ort.NewTensor[int64]([]int64{1, int64(len(tokenTypeIDs64))}, tokenTypeIDs64)
	if err != nil {
		log.Fatalf("创建tokenTypeIds张量失败: %v", err)
	}
	defer tokenTypeIDsTensor.Destroy()

	// 定义输出张量的形状（根据模型的实际输出形状进行调整）
	outputShape := ort.NewShape(1024) // 实际输出维度为1024
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		log.Fatalf("创建输出张量失败: %v", err)
	}
	defer outputTensor.Destroy()

	// 创建 SessionOptions
	options, err := ort.NewSessionOptions()
	if err != nil {
		log.Fatalf("Failed to create session options: %v", err)
	}
	// 创建CUDA Provider Options
	cudaOptions, err := ort.NewCUDAProviderOptions()
	if err != nil {
		log.Fatalf("Failed to append CUDA execution provider: %v", err)
	}
	defer cudaOptions.Destroy()
	// 设置 CUDA 执行提供程序
	err = options.AppendExecutionProviderCUDA(cudaOptions)
	if err != nil {
		log.Fatalf("Failed to append CUDA execution provider: %v", err)
	}
	// 创建会话
	session, err := ort.NewAdvancedSession(
		"../res/flag_embedding_model.onnx",
		inputNames,
		outputNames,
		[]ort.Value{inputIdsTensor, attentionMaskTensor, tokenTypeIDsTensor},
		[]ort.Value{outputTensor},
		options,
	)
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer session.Destroy()

	// 开始计时
	start := time.Now()

	// 运行推理
	// 只有一个err返回值可以设置在if作用域中,简化写法
	if err := session.Run(); err != nil {
		log.Fatalf("运行推理失败: %v", err)
	}

	results := outputTensor.GetData()
	fmt.Printf("推理结果: %v\n", results)
	fmt.Printf("推理结果形状: %v\n", len(results))

	duration := time.Since(start)
	fmt.Printf("推理时间: %v\n", duration)
}
