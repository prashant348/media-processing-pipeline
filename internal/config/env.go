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
	OutputDir string
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
		OutputDir: os.Getenv("OUTPUT_DIR"),
	}
}