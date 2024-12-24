// gin 框架的request参数类型是定义在handler func()里面的嘛？
// binding:"request"的校验机制
// ShouldBindBodyWithJson()等类似方法的用法
// gin.H{}和c.JSON()用法区别？c.JSON()会自动返回不需要return
package controller

import (
	"math"
	"project01/onnx"
)

func CosineSimilarity(vec1, vec2 []float32) float64 {
	// 长度不一致｜向量的模为0都需要报错，下游进行try catch
	// go的try catch是不是一般不抓指定的错误类型，笼统的通过if err != nil 来判断
	if len(vec1) != len(vec2) {
		return 0
	}
	var dotProduct, magnitude1, magnitude2 float64
	for i := 0; i < len(vec1); i++ {
		dotProduct += float64(vec1[i] * vec2[i])
		magnitude1 += float64(vec1[i] * vec1[i])
		magnitude2 += float64(vec2[i] * vec2[i])
	}
	if magnitude1 == 0 || magnitude2 == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(magnitude1) * math.Sqrt(magnitude2))
}

func TextSimilarity(text1, text2 string) (float64, error) {
	embedding1, err := onnx.GetEmbedding(text1)
	if err != nil {
		return 0, err
	}
	embedding2, err := onnx.GetEmbedding(text2)
	if err != nil {
		return 0, err
	}

	return CosineSimilarity(embedding1, embedding2), nil
}
