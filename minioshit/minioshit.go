package minioshit

import (
	"context"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

/*
docker run --rm -d \
   -p 9000:9000 \
   -p 9011:9011 \
   --name minio \
   -v ~/minio/data:/data \
   -e "MINIO_ROOT_USER=ROOTNAME" \
   -e "MINIO_ROOT_PASSWORD=CHANGEME123" \
   quay.io/minio/minio server /data --console-address ":9011"
*/

const (
	// endpoint := "play.min.io"
	// accessKeyID := "Q3AM3UQ867SPQQA43P2F"
	// secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
	// useSSL := true
	endpoint        = "127.0.0.1:9000"
	accessKeyID     = "A9FY9Qgi1zbo2eYeH8V3"
	secretAccessKey = "yYD5SVUxEeN5wr2w1u9y1FuDzVeE7uzQKhQdk0Lk"
	useSSL          = false

	// Make a new bucket called testbucket.
	bucketName = "pic-meiqia-upload"
	location   = "us-east-1"
)

var (
	ctx         = context.Background()
	minioClient *minio.Client
)

func init() {
	// Initialize minio client object.
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	minioClient = client
}

func Upload() {
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	// Upload the test file
	// Change the value of filePath if the file is in another location
	objectName := "superman.png"
	filePath := "/Users/ranya/Downloads/superman.png"
	contentType := "application/octet-stream"

	// Upload the test file with FPutObject
	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %+v\n", objectName, info)

	// reqParams := make(url.Values)
	// reqParams.Set("response-content-disposition", "attachment; filename="+filePath)
	// u, err := minioClient.PresignedGetObject(ctx, bucketName, objectName, time.Minute*1, reqParams)
	// if err != nil {
	// 	log.Fatalf("PresignedGetObject failed: %v", err)
	// }
	// log.Println(u.String())
}

func ExpireRead() {
	objectName := "superman.png"
	// filePath := "/Users/ranya/Downloads/superman.png"
	// reqParams := make(url.Values)
	// reqParams.Set("response-content-disposition", "attachment; filename="+filePath)
	u, err := minioClient.PresignedGetObject(ctx, bucketName, objectName, time.Second*30, nil)
	if err != nil {
		log.Fatalf("PresignedGetObject failed: %v", err)
	}
	log.Println(u.String())
}

func Delete() {
	objectName := "superman.png"
	err := minioClient.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Fatalf("RemoveObject failed: %v", err)
	}
}

func Main() {
	Upload()
	ExpireRead()
	// Delete()
}
