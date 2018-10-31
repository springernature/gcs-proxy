package gcs_proxy

type ObjectRepository interface {
	GetObjects(path string) []string
}
