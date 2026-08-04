package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Env struct {
	MinioEndpoint   string
	MinioRootUser   string
	MinioRootPass   string
	MinioBucketName string
	OutputDir       string
	BaseUrl         string
	ClientBaseUrl   string
}

func InitEnv(filenames ...string) error {
	if err := godotenv.Load(filenames...); err != nil {
		if len(filenames) == 0 && os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return nil
}

func LoadEnv() *Env {
	return &Env{
		MinioEndpoint:   os.Getenv("MINIO_ENDPOINT"),
		MinioRootUser:   os.Getenv("MINIO_ROOT_USER"),
		MinioRootPass:   os.Getenv("MINIO_ROOT_PASSWORD"),
		MinioBucketName: os.Getenv("MINIO_BUCKET_NAME"),
		OutputDir:       os.Getenv("OUTPUT_DIR"),
		BaseUrl:         os.Getenv("BASE_URL"),
		ClientBaseUrl:   os.Getenv("CLIENT_BASE_URL"),
	}
}
