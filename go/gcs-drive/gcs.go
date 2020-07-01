package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/storage"
)

func GetGCS() *storage.Client {
	var gcsPath string
	if os.Getenv("ENV") != "" {
		gcsPath = "/root/"
	}
	gcsPath += "files/gcs.json"
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", gcsPath)
	ctx := context.Background()
	gcsClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	return gcsClient
}

func GetGCSObject(bucket, object string) *storage.ObjectHandle {
	client := GetGCS()
	return client.Bucket(bucket).Object(object)
}
