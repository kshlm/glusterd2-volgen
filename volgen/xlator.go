package volgen

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/apex/log"
)

var xlatorMap = make(map[string]*Node)

var (
	ErrXlatorNotFound = errors.New("xlator not found")
)

func LoadXlators(path string) error {
	log.WithField("path", path).Debug("loading xlator nodefiles")

	var count = 0

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		log.WithField("file", path).Trace("found file")
		if filepath.Ext(path) != xlatorExt {
			log.WithField("file", path).Trace("file not xlator nodefile, skipping")
			return nil
		}
		n, err := NodeFromFile(path)
		if err != nil {
			return nil
		}
		xlatorMap[filepath.Base(path)] = n
		log.WithFields(log.Fields{
			"id":   n.ID,
			"dump": n,
		}).Debug("loaded xlator")
		count += 1
		return nil
	})
	log.WithFields(log.Fields{
		"path":  path,
		"count": count,
	}).Info("loaded xlators from path")
	return err
}

func FindXlator(id string) (*Node, error) {
	x, ok := xlatorMap[id]
	if !ok {
		return nil, ErrXlatorNotFound
	}
	return x, nil
}
