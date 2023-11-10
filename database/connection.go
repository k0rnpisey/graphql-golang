package database

import (
	"fmt"
	"github.com/go-gorm/caches/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"regexp"
)

func ConnectWith(url string) (db *gorm.DB, err error) {
	fmt.Printf("Connecting go %s\n\n", maskDBUrl(url))
	level := logger.Error
	if os.Getenv("ENABLE_QUERY_LOG") == "1" {
		level = logger.Info
	}
	db, err = gorm.Open(postgres.Open(url), &gorm.Config{
		PrepareStmt: true,
		Logger:      logger.Default.LogMode(level),
	})

	cachesPlugin := &caches.Caches{Conf: &caches.Config{
		Cacher: &Cacher{},
	}}
	_ = db.Use(cachesPlugin)

	if err != nil {
		fmt.Printf("Gorm connection open failed. Reason: [%T] %+v\n", err, err)
		return nil, err
	}

	// set timezone to Asia/Phnom_Penh
	db.Exec("SET TIME ZONE 'Asia/Phnom_Penh';")
	return db, nil
}

func maskDBUrl(url string) string {
	masked := regexp.MustCompile(`//.*@`).ReplaceAllString(url, "//***:***@")
	masked = regexp.MustCompile(`([?&])user=[^&]*`).ReplaceAllString(masked, "${1}user=***")
	masked = regexp.MustCompile(`([?&])password=[^&]*`).ReplaceAllString(masked, "${1}password=***")
	return masked
}

const envDBUrl = "DB_URL"

var engine *gorm.DB

// Connect connects to the database using GORM
func Connect() (conn *gorm.DB, err error) {
	url := os.Getenv(envDBUrl)
	if url == "" {
		err = fmt.Errorf("$%s isn't given", envDBUrl)
		return
	}
	conn, err = ConnectWith(url)
	engine = conn
	return conn, err
}

func GetConnection() *gorm.DB {
	return engine
}
