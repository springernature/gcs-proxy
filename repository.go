package gcs_proxy

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/iterator"
	"strings"
)

type ObjectRepository interface {
	GetObjects(path string) (objects []Object, err error)
	GetObject(path string) []byte
	IsFile(path string) bool
}

type Object struct {
	Path string
	Name string
}

type repo struct {
	client *storage.Client
}

func NewRepository(client *storage.Client) ObjectRepository {
	return repo{
		client: client,
	}
}

func (repo) GetObject(path string) []byte {
	panic("implement me")
}

func (repo) IsFile(path string) bool {
	panic("implement me")
}

func (r repo) GetObjects(path string) (objects []Object, err error) {
	if path == "/" {
		path = ""
	} else if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	bh := r.client.Bucket("bucket")
	oi := bh.Objects(context.Background(), &storage.Query{
		Delimiter: "/",
		Prefix:    path,
	})

	for {
		attrs, iteratorErr := oi.Next()
		if iteratorErr == iterator.Done {
			break
		}
		if iteratorErr != nil {
			err = iteratorErr
			return
		}

		numOfSlashes := len(strings.Split(attrs.Name, "/"))
		name := strings.Split(attrs.Name, "/")[numOfSlashes-1]
		objects = append(objects, Object{
			Name: name,
			Path: attrs.Name,
		})
	}

	return
}
