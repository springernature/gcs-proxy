package gcs_proxy

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/iterator"
	"io"
	"strings"
)

type ObjectRepository interface {
	GetObjects(path string) (objects []Object, err error)
	GetObject(path string) (objectContent []byte, err error)
	IsFile(path string) (isFile bool, err error)
}

type repo struct {
	bucket string
	client *storage.Client
}

func NewRepository(bucket string, client *storage.Client) ObjectRepository {
	return repo{
		bucket: bucket,
		client: client,
	}
}

func (r repo) GetObject(path string) (objectContent []byte, err error) {
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	bh := r.client.Bucket(r.bucket)
	oh := bh.Object(path)
	reader, err := oh.NewReader(context.Background())
	if err != nil {
		return
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

func (r repo) IsFile(path string) (bool, error) {
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	bh := r.client.Bucket(r.bucket)
	oh := bh.Object(path)
	oa, err := oh.Attrs(context.Background())
	if err == storage.ErrObjectNotExist {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return oa.Name != "" && oa.Prefix == "", nil
}

func (r repo) GetObjects(objectPath string) (objects []Object, err error) {
	if objectPath == "/" {
		objectPath = ""
	} else if !strings.HasSuffix(objectPath, "/") {
		objectPath += "/"
	}
	if strings.HasPrefix(objectPath, "/") {
		objectPath = objectPath[1:]
	}

	bh := r.client.Bucket(r.bucket)
	oi := bh.Objects(context.Background(), &storage.Query{
		Delimiter: "/",
		Prefix:    objectPath,
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

		if isDirectory(attrs) {
			objects = append(objects, createDirectoryObject(attrs))
		} else {
			objects = append(objects, createFileObject(attrs))
		}
	}

	return
}

func isDirectory(attrs *storage.ObjectAttrs) bool {
	return attrs.Name == "" && attrs.Prefix != ""
}

func createDirectoryObject(attrs *storage.ObjectAttrs) Object {
	splitPrefix := strings.Split(attrs.Prefix, "/")
	return Directory{
		Name: splitPrefix[len(splitPrefix)-2],
		Path: attrs.Prefix,
	}
}

func createFileObject(attrs *storage.ObjectAttrs) Object {
	numOfSlashes := len(strings.Split(attrs.Name, "/"))
	return File{
		Name: strings.Split(attrs.Name, "/")[numOfSlashes-1],
		Path: attrs.Name,
	}
}
