package gcs_proxy

type Object interface {
	GetName() string
	GetPath() string
}

type File struct {
	Path string
	Name string
}

func (f File) GetName() string {
	return f.Name
}

func (f File) GetPath() string {
	return f.Path
}

type Directory struct {
	Path string
	Name string
}

func (d Directory) GetName() string {
	return d.Name
}

func (d Directory) GetPath() string {
	return d.Path
}

