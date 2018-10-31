package main

import (
	"net/http"
	"github.com/springernature/gcs-proxy"
	"log"
)

func main() {
	or := gcs_proxy.Repository{}
	server := gcs_proxy.NewServer(or)
	http.HandleFunc("/", server.Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
