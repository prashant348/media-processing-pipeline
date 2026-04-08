package config

import "os"

type Env struct {
	MinioEndpoint string
	MinioRootUser string
	MinioRootPass string
	MinioBucketName string
}

func LoadEnv() *Env {
	return &Env{
		MinioEndpoint: getEnv("MINIO_ENDPOINT"),
		MinioRootUser: getEnv("MINIO_ROOT_USER"),
		MinioRootPass: getEnv("MINIO_ROOT_PASSWORD"),
		MinioBucketName: getEnv("MINIO_BUCKET_NAME"),
	}
}

func getEnv(key string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return ""
}

