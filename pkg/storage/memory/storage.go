package memory

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/bonaysoft/gofbot/apis/message/v1alpha1"
)

type Storage struct {
	srcDir string

	store sync.Map
}

func NewStorage(srcDir string) *Storage {
	return &Storage{srcDir: srcDir}
}

func (s *Storage) Start(ctx context.Context) error {
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
				if !ok {
					return
				}
				// 检查是否是创建事件
				if event.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("Created file:", event.Name)
					if err := s.loadMessage(event.Name); err != nil {
						fmt.Println(err)
						return
					}
				}
			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}
		}
	}()

	return w.Add(s.srcDir)
}

func (s *Storage) List(ctx context.Context) ([]v1alpha1.Message, error) {
	messages := make([]v1alpha1.Message, 0)
	s.store.Range(func(key, value any) bool {
		messages = append(messages, value.(v1alpha1.Message))
		return true
	})
	return messages, nil
}

func (s *Storage) loadMessage(name string) error {
	yamlFile, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	var message v1alpha1.Message
	if err := yaml.Unmarshal(yamlFile, &message); err != nil {
		return err
	}

	message.SetUID(types.UID(uuid.NewSHA1(uuid.New(), yamlFile).String()))
	s.store.Store(message.UID, message)
	return nil
}

func (s *Storage) loadExistedMessages() error {
	return filepath.Walk(s.srcDir, func(filepath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if info.IsDir() {
			return nil
		} else if path.Ext(filepath) != ".yaml" && path.Ext(filepath) != ".yml" {
			return nil
		}

		return s.loadMessage(filepath)
	})
}
