package gcs_proxy_test

import (
	"testing"
	"github.com/fsouza/fake-gcs-server/fakestorage"
	"github.com/springernature/gcs-proxy"
	"github.com/stretchr/testify/assert"
	"fmt"
)

var bucket = "myBucket"
func TestShouldReturnAllTheFoldersAtRoot(t *testing.T) {
	objects := []fakestorage.Object{
		{
			BucketName: bucket,
			Name:       "folder1/",
		},
		{
			BucketName: bucket,
			Name:       "folder2/",
		},
		{
			BucketName: bucket,
			Name:       "folder3/",
		},
	}

	server := fakestorage.NewServer(objects)
	client := server.Client()

	repo := gcs_proxy.NewRepository(bucket, client)
	path := "/"

	objs, err := repo.GetObjects(path)
	assert.NoError(t, err)
	assert.Contains(t, objs, gcs_proxy.Directory{Name: objects[0].Name, Path: objects[0].Name})
	assert.Contains(t, objs, gcs_proxy.Directory{Name: objects[1].Name, Path: objects[1].Name})
	assert.Contains(t, objs, gcs_proxy.Directory{Name: objects[2].Name, Path: objects[2].Name})
}

func TestShouldReturnAllTheObjectsAtSomePath(t *testing.T) {
	path := "subPath"
	objects := []fakestorage.Object{
		{
			BucketName: bucket,
			Name:       "file1",
		},
		{
			BucketName: bucket,
			Name:       fmt.Sprintf("%s/file2", path),
		},
		{
			BucketName: bucket,
			Name:       fmt.Sprintf("%s/file3", path),
		},
		{
			BucketName: bucket,
			Name:       fmt.Sprintf("%s/subDirectory/file4", path),
		},
	}

	server := fakestorage.NewServer(objects)
	client := server.Client()

	repo := gcs_proxy.NewRepository(bucket, client)

	objs, err := repo.GetObjects(path)
	assert.NoError(t, err)
	assert.Contains(t, objs, gcs_proxy.File{
		Name: "file2",
		Path: fmt.Sprintf("%s/file2", path),
	})
	assert.Contains(t, objs, gcs_proxy.File{
		Name: "file3",
		Path: fmt.Sprintf("%s/file3", path),
	})
	assert.Contains(t, objs, gcs_proxy.Directory{
		Name: "subDirectory/",
		Path: fmt.Sprintf("%s/subDirectory/", path),
	})
}

func TestShouldReturnAllTheObjectsAtSomeEvenDeeperPath(t *testing.T) {
	path := "subPath/subPathAgain"
	objects := []fakestorage.Object{
		{
			BucketName: bucket,
			Name:       "file1",
		},
		{
			BucketName: bucket,
			Name:       fmt.Sprintf("%s/file2", path),
		},
		{
			BucketName: bucket,
			Name:       fmt.Sprintf("%s/file3", path),
		},
	}

	server := fakestorage.NewServer(objects)
	client := server.Client()

	repo := gcs_proxy.NewRepository(bucket, client)

	objs, err := repo.GetObjects(path)
	assert.NoError(t, err)
	assert.Contains(t, objs, gcs_proxy.File{
		Name: "file2",
		Path: fmt.Sprintf("%s/file2", path),
	})
	assert.Contains(t, objs, gcs_proxy.File{
		Name: "file3",
		Path: fmt.Sprintf("%s/file3", path),
	})
}

func TestGetObjectShouldReturnTheContentOfTheObject(t *testing.T) {
	path := "subPath/subPathAgain/file2"
	expected := []byte(("ImaFIIIIIIIIILE"))
	objects := []fakestorage.Object{
		{
			BucketName: bucket,
			Name:       "file1",
		},
		{
			BucketName: bucket,
			Name:       path,
			Content:    expected,
		},
		{
			BucketName: bucket,
			Name:       fmt.Sprintf("%s/file3", path),
		},
	}

	server := fakestorage.NewServer(objects)
	client := server.Client()

	repo := gcs_proxy.NewRepository(bucket, client)

	object, err := repo.GetObject(path)
	assert.NoError(t, err)
	assert.Equal(t, expected, object)
}

func TestGetObjectShouldReturnErrorIfFileNotFound(t *testing.T) {
	path := "subPath/subPathAgain/file2"
	var objects []fakestorage.Object

	server := fakestorage.NewServer(objects)
	client := server.Client()

	repo := gcs_proxy.NewRepository(bucket, client)

	object, err := repo.GetObject(path)
	assert.Error(t, err)
	assert.Equal(t, []byte(nil), object)
}

func TestIsFileShouldReturnsFalseIffFileNotFound(t *testing.T) {
	path := "subPath/subPathAgain/file2"
	var objects []fakestorage.Object

	server := fakestorage.NewServer(objects)
	client := server.Client()

	repo := gcs_proxy.NewRepository(bucket, client)

	isFile, err := repo.IsFile(path)
	assert.NoError(t, err)
	assert.False(t, isFile)
}

func TestIsFileShouldReturnTrueIfFileIsFound(t *testing.T) {
	path := "subPath/subPathAgain/file2"
	objects := []fakestorage.Object{
		{BucketName: bucket, Name: "file1"},
		{BucketName: bucket, Name: path},
		{BucketName: bucket, Name: "some/random/path/file3"},
	}

	server := fakestorage.NewServer(objects)
	client := server.Client()

	repo := gcs_proxy.NewRepository(bucket, client)

	isFile, err := repo.IsFile(path)
	assert.NoError(t, err)
	assert.True(t, isFile)
}
