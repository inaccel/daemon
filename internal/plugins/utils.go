package plugins

import (
	"net"
	"os"
	"path/filepath"
)

func listen(path string) (net.Listener, error) {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return nil, err
	}

	if err := os.RemoveAll(path); err != nil {
		return nil, err
	}

	return net.Listen("unix", path)
}
