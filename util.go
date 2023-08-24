package logf

import (
	"fmt"
	"os"
	"path"
)

func FileWriter(pathfile string) (*os.File, error) {
	p := path.Dir(pathfile)
	info, err := os.Stat(p)
	if (err != nil && os.IsNotExist(err)) || !info.IsDir() {
		if err := os.MkdirAll(p, 0755); err != nil {
			return nil, fmt.Errorf("make log path [%s]: %w", p, err)
		}
		err = nil
	}
	if err != nil {
		return nil, err
	}
	return os.OpenFile(pathfile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
}
