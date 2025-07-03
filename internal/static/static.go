package static

import (
	"net/http"
	"path/filepath"
	"runtime"
)

func Static() error {
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	webDir := filepath.Join(projectRoot, "web")
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	return nil
}
