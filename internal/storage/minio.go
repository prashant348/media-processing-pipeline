package storage

import (
	"context"
	"io"
	"log"
	"media_processing_pipeline/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ObjectStore interface {
	PutObject(
		ctx context.Context,
		bucketName string,
		objectName string,
		reader io.Reader,
		size int64,
		opts minio.PutObjectOptions,
	) (info minio.UploadInfo, err error)
	
	GetObject(
		ctx context.Context, 
		bucketName string, 
		objectName string, 
		opts minio.GetObjectOptions,
	) (*minio.Object, error)
}

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