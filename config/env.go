package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Env struct {
	MinioEndpoint string
	MinioRootUser string
	MinioRootPass string
	MinioBucketName string
}

func InitEnv(filenames ...string) error {
	if err := godotenv.Load(filenames...); err != nil {
		return err
	}

	return nil
} 

func LoadEnv() *Env {
	return &Env{
		MinioEndpoint: os.Getenv("MINIO_ENDPOINT"),
		MinioRootUser: os.Getenv("MINIO_ROOT_USER"),
		MinioRootPass: os.Getenv("MINIO_ROOT_PASSWORD"),
		MinioBucketName: os.Getenv("MINIO_BUCKET_NAME"),
	}
}


// func LoadEnv() *Env {
// 	return &Env{
// 		MinioEndpoint: getEnv("MINIO_ENDPOINT"),
// 		MinioRootUser: getEnv("MINIO_ROOT_USER"),
// 		MinioRootPass: getEnv("MINIO_ROOT_PASSWORD"),
// 		MinioBucketName: getEnv("MINIO_BUCKET_NAME"),
// 	}
// }

// func getEnv(key string) string {
// 	if val, ok := os.LookupEnv(key); ok {
// 		return val
// 	}
// 	return ""
// }

