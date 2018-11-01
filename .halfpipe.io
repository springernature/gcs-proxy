team: engineering-enablement
pipeline: gcs-proxy
tasks:
- type: run
  name: Test and build
  script: \go test ./... ; go build -o cf/server cmd/main.go
  docker:
    image: golang:1.11-stretch
