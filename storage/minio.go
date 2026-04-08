package storage

import (
	"log"
	"media_processing_pipeline/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func InitMinIO(env *config.Env) {

	client, err := minio.New(env.MinioEndpoint, &minio.Options{
		Creds: credentials.NewStaticV4(env.MinioRootUser, env.MinioRootPass, ""),
		Secure: false,
	})

	if err != nil {
		log.Fatal(err)
	}

	MinioClient = client
}