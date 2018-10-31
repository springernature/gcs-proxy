package gcs_proxy

type ObjectRepository interface {
	GetObjects(path string) []string
}

type Repository struct {
}

func (Repository) GetObjects(path string) []string {
	return []string{
		"a",
		"b",
		"c",
	}
}
