package entity

import "io/fs"

// File represents a file to be written to disk.
type File struct {
	Path    string      // relative path from project root
	Content []byte      // content to write
	Mode    fs.FileMode // file permissions
}

// Dir represents a directory to create.
type Dir struct {
	Path string
	Mode fs.FileMode
}
