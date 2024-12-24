// `gorm:"type:text;not null"`结构体的这种用法讲讲，还有哪些类似的情况？
// DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})这个方法是什么意思？dsn的作用，以及fmt.Sprintf()的作用
package db

import (
	"fmt"
	"log"
	"sync"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	dbOnce sync.Once
)

// InitDB initializes the database connection
func InitDB() {
	dbOnce.Do(func() {
		// Get DSN from config file
		user := viper.GetString("database.user")
		password := viper.GetString("database.password")
		host := viper.GetString("database.host")
		port := viper.GetInt("database.port")
		name := viper.GetString("database.name")
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, name)

		var err error
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		// Auto migrate the database schema
		if err := DB.AutoMigrate(&SimilarityResult{}); err != nil {
			log.Fatalf("Failed to migrate database schema: %v", err)
		}

	})
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

type SimilarityResult struct {
	ID         uint   `gorm:"primaryKey"`
	Text1      string `gorm:"type:text;not null"`
	Text2      string `gorm:"type:text;not null"`
	Similarity float64
	CreatedAt  int64 `gorm:"autoCreateTime"`
}
