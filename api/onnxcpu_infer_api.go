// log.Fatalf()是否有更复杂的动态log组件？
package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"net/http"
	onnx "project01/onnx"
	"sync"

	"github.com/gin-gonic/gin"
	ort "github.com/yalue/onnxruntime_go"
)

var (
	db     *sql.DB
	dbOnce sync.Once
)

var (
	ortAdvancedSession *ort.AdvancedSession
	inputNames         []string
	outputNames        []string
)

func main() {
	// Initialize ONNX Runtime environment
	err := ort.InitializeEnvironment()
	if err != nil {
		log.Fatalf("Failed to initialize ONNX Runtime: %v", err)
	}
	defer ort.DestroyEnvironment()

	// Initialize ONNX Runtime Session
	modelPath := "../res/flag_embedding_model.onnx"
	inputs, outputs, err := ort.GetInputOutputInfo(modelPath)
	if err != nil {
		log.Fatalf("Failed to get input output info from onnx file: %v", err)
	}
	// 提取inputNames和outputNames
	inputNames = make([]string, len(inputs))
	for i, input := range inputs {
		inputNames[i] = input.Name
	}
	outputNames = make([]string, len(outputs))
	for i, output := range outputs {
		outputNames[i] = output.Name
	}
	//Create ONNX Session
	options, err := ort.NewSessionOptions()
	if err != nil {
		log.Fatalf("Failed to create NewSessionOptions()")
	}
	cudaOptions, err := ort.NewCUDAProviderOptions()
	if err != nil {
		log.Fatalf("Failed to create NewCUDAProviderOptions:%v", err)
	}
	err = options.AppendExecutionProviderCUDA(cudaOptions)
	if err != nil {
		log.Fatalf("Failed to append CUDA execution provider: %v", err)
	}
	session, err := ort.NewAdvancedSession(modelPath, inputNames, outputNames, nil, nil, options)
	if err != nil {
		log.Fatalf("创建 ONNX 高级会话失败: %v", err)
	}
	defer session.Destroy()

	inputNames, outputNames, err := ort.GetInputOutputNames(modelPath)
	if err != nil {
		log.Fatalf("Failed to get input/output names: %v", err)
	}
	ortSession, err = ort.NewAdvancedSession(modelPath, inputNames, outputNames, nil, nil, nil)
	if err != nil {
		log.Fatalf("Failed to create ONNX AdvancedSession: %v", err)
	}
	defer ortSession.Destroy()

	// Initialize database connection
	initDB()

	// Set up Gin router
	router := gin.Default()
	// Define the API endpoint
	router.POST("/similarity", func(c *gin.Context) {
		var request struct {
			Text1 string `json:"text1" binding:"required"`
			Text2 string `json:"text2" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		similarity, err := calculateSimilarity(request.Text1, request.Text2)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		// Write to MYSQL database
		if err := saveToDatabase(request.Text1, request.Text2, similarity); err != nil {
			log.Printf("Failed to save to databse: %v", err)
		}

		c.JSON(http.StatusOK, gin.H{"similarity": similarity})
	})

	router.Run(":11260")
}

func initDB() {
	dbOnce.Do(func() {
		var err error
		db, err = sql.Open("mysql", "username:password@tcp(localhost:3306)/dbname")
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		// Create table if not exists
		createTableQuery := `
		CREATE TABLE IF NOT EXISTS similarity_results (
			id INT AUTO_INCREMENT PRIMARY KEY,
			text1 TEXT NOT NULL,
			text2 TEXT NOT NULL,
			similarity DOUBLE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
		`
		_, err = db.Exec(createTableQuery)
		if err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}
	})
}

func calculateSimilarity(text1, text2 string) (float64, error) {
	modelPath := "../res/flag_embedding_model.onnx" // TODO:要替换成config
	inputs, outputs, err := ort.GetInputOutputInfo(modelPath)
	if err != nil {
		log.Fatalf("Failed to GetInputOutputInfo from onnx file: %v", err)
	}
	// Get inputNames & outputNames
	inputNames := make([]string, len(inputs))
	for i, input := range inputs {
		inputNames[i] = input.Name
	}
	outputNames := make([]string, len(outputs))
	for i, output := range outputs {
		outputNames[i] = output.Name
	}

	embedding1, err := getEmbedding(text1, modelPath, inputNames, outputNames)
	if err != nil {
		return 0, fmt.Errorf("Failed to get embedding for text1: %v", err)
	}
	embedding2, err := getEmbedding(text2, modelPath, inputNames, outputNames)
	if err != nil {
		return 0, fmt.Errorf("Failed to get embedding for text2: %v", err)
	}

	similarity := cosineSimilarity(embedding1, embedding2)

	return similarity, nil
}

func getEmbedding(text, modelPath string, inputNames, outputNames []string) ([]float32, error) {
	// Prepare input data
	inputIds, attentionMask, tokenTypeIDs, err := onnx.GetOnnxInput(text)
	if err != nil {
		log.Fatalf("Failed to get model input: %v", err)
	}
	inputIds64 := onnx.GetInt64Slice(inputIds)
	attentionMask64 := onnx.GetInt64Slice(attentionMask)
	tokenTypeIDs64 := onnx.GetInt64Slice(tokenTypeIDs)

	// Convert to tensors
	inputIdsTensor, err := ort.NewTensor[int64]([]int64{1, int64(len(inputIds64))}, inputIds64)
	if err != nil {
		return nil, fmt.Errorf("Failed to create inputIds tensor: %v", err)
	}
	defer inputIdsTensor.Destroy()

	attentionMaskTensor, err := ort.NewTensor[int64]([]int64{1, int64(len(attentionMask64))}, attentionMask64)
	if err != nil {
		return nil, fmt.Errorf("Failed to create attentionMask tensor: %v", err)
	}
	defer attentionMaskTensor.Destroy()

	tokenTypeIDsTensor, err := ort.NewTensor[int64]([]int64{1, int64(len(tokenTypeIDs64))}, tokenTypeIDs64)
	if err != nil {
		return nil, fmt.Errorf("Failed to create tokenTypeIDs tensor: %v", err)
	}
	defer tokenTypeIDsTensor.Destroy()

	// Define output tensor shape (adjust according to your model)
	outputShape := ort.NewShape(1, 1024)
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		log.Fatalf("Failed to create ouput tensor: %v", err)
	}
	defer outputTensor.Destroy()

	// Create session
	session, err := ort.NewAdvancedSession(
		modelPath,
		inputNames,
		outputNames,
		[]ort.Value{inputIdsTensor, attentionMaskTensor, tokenTypeIDsTensor},
		[]ort.Value{outputTensor},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create session: %v", err)
	}
	defer session.Destroy()

	// Run inference
	if err := session.Run(); err != nil {
		return nil, fmt.Errorf("Failed to run inference: %v", err)
	}

	results := outputTensor.GetData()
	return results, nil
}

func cosineSimilarity(vec1, vec2 []float32) float64 {
	if len(vec1) != len(vec2) {
		return 0
	}
	var dotProduct float64
	var magnitude1 float64
	var magnitude2 float64

	for i := 0; i < len(vec1); i++ {
		dotProduct += float64(vec1[i] * vec2[i])
		magnitude1 += float64(vec1[i] * vec1[i])
		magnitude2 += float64(vec2[i] * vec2[i])
	}

	if magnitude1 == 0 || magnitude2 == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(magnitude1) + math.Sqrt(magnitude2))
}

func saveToDatabase(text1, text2 string, similarity float64) error {
	query := `
		INSERT INTO similarity_results (text1, text2, similarity)
		VALUES (?,?,?)
	`
	_, err := db.Exec(query, text1, text2, similarity)

	return err
}
