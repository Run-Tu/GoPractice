// viper的使用逻辑，是怎么根据GetString将配置文件映射成参数的？
package onnx

import (
	"log"
	"sync"

	"github.com/spf13/viper"
	ort "github.com/yalue/onnxruntime_go"
)

var (
	Session     *ort.DynamicAdvancedSession
	sessionOnce sync.Once
)

func InitONNX() {
	sessionOnce.Do(func() {
		if err := ort.InitializeEnvironment(); err != nil {
			log.Fatalf("初始化ONNX Runtime失败: %v", err)
		}
		modelPath := viper.GetString("onnx.modelPath")

		// 提取输入输出名称
		inputs, outputs, err := ort.GetInputOutputInfo(modelPath)
		if err != nil {
			log.Fatalf("获取模型输入输出信息失败: %v", err)
		}
		inputNames := make([]string, len(inputs))
		for i, input := range inputs {
			inputNames[i] = input.Name
		}
		outputNames := make([]string, len(outputs))
		for i, output := range outputs {
			outputNames[i] = output.Name
		}

		// 创建会话选项
		options, err := ort.NewSessionOptions()
		if err != nil {
			log.Fatalf("创建会话选项失败: %v", err)
		}
		defer options.Destroy()

		// 判断是否使用GPU
		if viper.GetBool("onnx.use_gpu") {
			cudaOptions, err := ort.NewCUDAProviderOptions()
			if err != nil {
				log.Fatalf("创建CUDA选项失败: %v ", err)
			}
			defer cudaOptions.Destroy()

			if err := options.AppendExecutionProviderCUDA(cudaOptions); err != nil {
				log.Fatalf("添加CUDA会话执行失败:%v", err)
			}

		}

		// 创建ONNX会话
		Session, err = ort.NewDynamicAdvancedSession(modelPath, inputNames, outputNames, options)
		if err != nil {
			log.Fatalf("创建ONNX会话失败: %v", err)
		}
	})
}

func DestroyONNX() {
	if Session != nil {
		Session.Destroy()
	}
	ort.DestroyEnvironment()
	log.Println("ONNX资源已释放")
}
