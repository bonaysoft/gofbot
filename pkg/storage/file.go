package storage

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/bonaysoft/gofbot/apis/message/v1alpha1"
	"github.com/bonaysoft/gofbot/pkg/errors"
)

type FileStorage struct {
	storeDir string

	cache sync.Map
}

func NewStorage() (Manager, error) {
	return &FileStorage{storeDir: viper.GetString("storage-file-location")}, nil
}

func (s *FileStorage) Start(ctx context.Context) error {
	if err := s.loadExistedMessages(); err != nil {
		return err
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-w.Events:
				if !ok || strings.HasSuffix(event.Name, "~") {
					return
				}

				if event.Op.Has(fsnotify.Create) {
					slog.Info("Create file", "filename", event.Name)
					if err := s.cacheStaticMessageFile(event.Name); err != nil {
						slog.Error("cacheStaticMessageFile failed", errors.With(err))
					}
				} else if event.Op.Has(fsnotify.Remove) || event.Op.Has(fsnotify.Rename) {
					slog.Info("Remove file", "filename", event.Name)
					s.cache.Delete(event.Name)
				}
			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				slog.Error("storage watch", errors.With(err))
			}
		}
	}()

	return w.Add(s.storeDir)
}

func (s *FileStorage) List(ctx context.Context) ([]v1alpha1.Message, error) {
	messages := make([]v1alpha1.Message, 0)
	s.cache.Range(func(key, value any) bool {
		filename := key.(string)
		message, _ := s.readFile2Message(filename) // ignore err because already validated
		messages = append(messages, *message)
		return true
	})

	return messages, nil
}

func (s *FileStorage) readFile2Message(name string) (*v1alpha1.Message, error) {
	yamlFile, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	var message v1alpha1.Message
	if err := yaml.Unmarshal(yamlFile, &message); err != nil {
		return nil, err
	}
	return &message, nil
}

func (s *FileStorage) loadExistedMessages() error {
	return filepath.Walk(s.storeDir, func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if info.IsDir() {
			return nil
		} else if path.Ext(filename) != ".yaml" && path.Ext(filename) != ".yml" {
			return nil
		}

		return s.cacheStaticMessageFile(filename)
	})
}

func (s *FileStorage) cacheStaticMessageFile(filename string) error {
	if _, err := s.readFile2Message(filename); err != nil {
		return fmt.Errorf("readFile2Message: %s - %s\n", filename, err)
	}

	s.cache.Store(filename, 1)
	return nil
}
