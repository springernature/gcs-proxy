package gcs_proxy

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/iterator"
	"strings"
	"io/ioutil"
)

type ObjectRepository interface {
	GetObjects(path string) (objects []Object, err error)
	GetObject(path string) (objectContent []byte, err error)
	IsFile(path string) (isFile bool, err error)
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

func (r repo) GetObject(path string) (objectContent []byte, err error) {
	bh := r.client.Bucket("bucket")
	oh := bh.Object(path)
	reader, err := oh.NewReader(context.Background())
	if err != nil {
		return
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}

func (r repo) IsFile(path string) (bool, error) {
	bh := r.client.Bucket("bucket")
	oh := bh.Object(path)
	_, err := oh.Attrs(context.Background())
	if err == storage.ErrObjectNotExist {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
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
