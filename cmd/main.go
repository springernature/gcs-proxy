package main

import (
	"net/http"
	"github.com/springernature/gcs-proxy"
	"log"
	"os"
	"cloud.google.com/go/storage"
	"context"
)

func main() {
	bucket := os.Getenv("BUCKET")
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	or := gcs_proxy.NewRepository(bucket, client)
	server := gcs_proxy.NewServer(or)
	http.HandleFunc("/", server.Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
