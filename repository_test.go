package gcs_proxy_test

import (
	"testing"
	"github.com/fsouza/fake-gcs-server/fakestorage"
	"github.com/springernature/gcs-proxy"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestShouldReturnAllTheObjectsAtRoot(t *testing.T) {
	objects := []fakestorage.Object{
		{
			BucketName: "bucket",
			Name:       "file1",
		},
		{
			BucketName: "bucket",
			Name:       "file2",
		},
		{
			BucketName: "bucket",
			Name:       "file3",
		},
	}

	server := fakestorage.NewServer(objects)
	client := server.Client()

	repo := gcs_proxy.NewRepository(client)
	path := "/"

	objs, err := repo.GetObjects(path)
	assert.NoError(t, err)
	assert.Contains(t, objs, gcs_proxy.Object{Name: objects[0].Name, Path: objects[0].Name})
	assert.Contains(t, objs, gcs_proxy.Object{Name: objects[1].Name, Path: objects[1].Name})
	assert.Contains(t, objs, gcs_proxy.Object{Name: objects[2].Name, Path: objects[2].Name})
}

func TestShouldReturnAllTheObjectsAtSomePath(t *testing.T) {
	path := "subPath"
	objects := []fakestorage.Object{
		{
			BucketName: "bucket",
			Name:       "file1",
		},
		{
			BucketName: "bucket",
			Name:       fmt.Sprintf("%s/file2", path),
		},
		{
			BucketName: "bucket",
			Name:       fmt.Sprintf("%s/file3", path),
		},
	}

	server := fakestorage.NewServer(objects)
	client := server.Client()

	repo := gcs_proxy.NewRepository(client)

	objs, err := repo.GetObjects(path)
	assert.NoError(t, err)
	assert.Contains(t, objs, gcs_proxy.Object{
		Name: "file2",
		Path: fmt.Sprintf("%s/file2", path),
	})
	assert.Contains(t, objs, gcs_proxy.Object{
		Name: "file3",
		Path: fmt.Sprintf("%s/file3", path),
	})
}

func TestShouldReturnAllTheObjectsAtSomeEvenDeeperPath(t *testing.T) {
	path := "subPath/subPathAgain"
	objects := []fakestorage.Object{
		{
			BucketName: "bucket",
			Name:       "file1",
		},
		{
			BucketName: "bucket",
			Name:       fmt.Sprintf("%s/file2", path),
		},
		{
			BucketName: "bucket",
			Name:       fmt.Sprintf("%s/file3", path),
		},
	}

	server := fakestorage.NewServer(objects)
	client := server.Client()

	repo := gcs_proxy.NewRepository(client)

	objs, err := repo.GetObjects(path)
	assert.NoError(t, err)
	assert.Contains(t, objs, gcs_proxy.Object{
		Name: "file2",
		Path: fmt.Sprintf("%s/file2", path),
	})
	assert.Contains(t, objs, gcs_proxy.Object{
		Name: "file3",
		Path: fmt.Sprintf("%s/file3", path),
	})
}

func TestGetObjectShouldReturnTheContentOfTheObject(t *testing.T) {
	path := "subPath/subPathAgain/file2"
	expected := []byte(("ImaFIIIIIIIIILE"))
	objects := []fakestorage.Object{
		{
			BucketName: "bucket",
			Name:       "file1",
		},
		{
			BucketName: "bucket",
			Name:       path,
			Content:    expected,
		},
		{
			BucketName: "bucket",
			Name:       fmt.Sprintf("%s/file3", path),
		},
	}

	server := fakestorage.NewServer(objects)
	client := server.Client()

	repo := gcs_proxy.NewRepository(client)

	object, err := repo.GetObject(path)
	assert.NoError(t, err)
	assert.Equal(t, expected, object)
}

func TestGetObjectShouldReturnErrorIfFileNotFound(t *testing.T) {
	path := "subPath/subPathAgain/file2"
	var objects []fakestorage.Object

	server := fakestorage.NewServer(objects)
	client := server.Client()

	repo := gcs_proxy.NewRepository(client)

	object, err := repo.GetObject(path)
	assert.Error(t, err)
	assert.Equal(t, []byte(nil), object)
}
