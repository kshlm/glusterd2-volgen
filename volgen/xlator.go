package volgen

import (
	"os"
	"path/filepath"
)

var xlatorMap map[string]*Node

func LoadXlators(path string) error {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != xlatorExt {
			return nil
		}
		n, err := NodeFromFile(path)
		if err != nil {
			return nil
		}
		xlatorMap[n.ID] = n
		return nil
	})
	return err
}
