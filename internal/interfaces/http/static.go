package http

import (
	"embed"
	"github.com/gin-contrib/static"
	"io/fs"
	"log"
	"net/http"
	"path"
)

var Static embed.FS

type serverFileSystemType struct {
	http.FileSystem
}

func (f serverFileSystemType) Exists(prefix string, _path string) bool {
	_, err := f.Open(path.Join(prefix, _path))
	return err == nil
}

func mustFS(dir string) (serverFileSystem static.ServeFileSystem) {

	sub, err := fs.Sub(Static, dir)

	if err != nil {
		log.Println(err)
		return
	}

	serverFileSystem = serverFileSystemType{
		http.FS(sub),
	}

	return
}
