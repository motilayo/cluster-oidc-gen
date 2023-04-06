package utils

import (
	"context"
	"fmt"
	// "log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type BucketBasics struct {
	S3Client *s3.Client
}

const (
	AWS_REGION = types.BucketLocationConstraintCaCentral1
)

func CreateBucketBasics() BucketBasics {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err.Error())
	}

	bucketBasics := BucketBasics{S3Client: s3.NewFromConfig(cfg)}

	return bucketBasics
}

func (basics BucketBasics) CreateBucket(bucketName string) {
	createBucketInput := s3.CreateBucketInput{
		Bucket: &bucketName,
		ACL:    types.BucketCannedACLPublicRead,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: AWS_REGION,
		},
	}

	_, err := basics.S3Client.CreateBucket(context.TODO(), &createBucketInput)
	if err != nil {
		log.Printf("Could not create bucket %s", bucketName)

		var bucketAlreadyOwnedByYou *types.BucketAlreadyOwnedByYou
		var bucketAlreadyExists *types.BucketAlreadyExists

		if errors.As(err, &bucketAlreadyOwnedByYou) {
			log.Info("The bucket you tried to create already exists, and you own it")
		} else if errors.As(err, &bucketAlreadyExists) {
			log.Fatalln(bucketAlreadyExists.ErrorCode())
		} else {
			log.Fatalln(err.Error())
		}
	}

}

func (basics BucketBasics) UploadToS3(bucketName string, fileName string, key string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Could not open file %v to upload", fileName)
		log.Fatalln(err.Error())
	}
	defer file.Close()
	_, err = basics.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &key,
		Body:   file,
	})
	if err != nil {
		log.Printf("Could not upload file %v to %v", fileName, bucketName)
		log.Fatalln(err.Error())
	}
}

func GenerateS3URL(bucketName string) string {
	uri := fmt.Sprintf("https://%s.s3.%s.amazonaws.com", bucketName, AWS_REGION)
	return uri
}
