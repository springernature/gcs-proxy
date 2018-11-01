package main

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/springernature/gcs-proxy"
	"log"
	"net/http"
	"os"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

func createGcsClient(jsonKey string) (client *storage.Client, err error) {
	if jsonKey == "" {
		log.Println("No jsonKey provided via env var, gonna try to use credz from the FS")
		return storage.NewClient(context.Background())
	}
	log.Println("jsonKey provided via env var, gonna try to use that")


	credz, err := google.CredentialsFromJSON(context.Background(), []byte(jsonKey), storage.ScopeReadOnly)
	if err != nil {
		return
	}
	return storage.NewClient(context.Background(), option.WithCredentials(credz))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
}


func main() {
	bucket := os.Getenv("BUCKET")
	if bucket == "" {
		log.Fatal("HEY! You need to specify the bucket you want to proxy via env var BUCKET")
	}
	gcsKey := os.Getenv("GCS_KEY")

	client, err := createGcsClient(gcsKey)
	if err != nil {
		log.Fatal(err)
	}
	or := gcs_proxy.NewRepository(bucket, client)
	server := gcs_proxy.NewServer(or)

	log.Println("Yay, gonna start serving stuff")
	http.HandleFunc("/", server.Handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/favicon.ico", faviconHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
