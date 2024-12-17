package messenger

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/bonaysoft/gofbot/apis/message/v1alpha1"
	"github.com/bonaysoft/gofbot/pkg/storage"
)

type Manager interface {
	Match(params map[string]any) (*v1alpha1.Message, error)

	BuildReply(msg *v1alpha1.Message, params map[string]any) ([]byte, error)
}

type DefaultManager struct {
	store   storage.Manager
	funcMap template.FuncMap
}

func NewDefaultManager(store storage.Manager, funcMap template.FuncMap) *DefaultManager {
	return &DefaultManager{store: store, funcMap: funcMap}
}

func (d *DefaultManager) Match(params map[string]any) (*v1alpha1.Message, error) {
	messages, err := d.store.List(context.Background())
	if err != nil {
		return nil, fmt.Errorf("list messages: %w", err)
	}

	// 目前只支持value为string的kv，后续可以把复杂结构拍平成一维的，比如{a:{b:{c:1}}}=> {a.b.c:1}
	newParams := lo.MapEntries(lo.PickBy(params, func(key string, value any) bool {
		_, ok := value.(string)
		return ok
	}), func(key string, value any) (string, string) {
		return key, value.(string)
	})

	matcher := func(item v1alpha1.Message) bool {
		selector, err := metav1.LabelSelectorAsSelector(&item.Spec.Selector)
		if err != nil {
			return false
		}

		return selector.Matches(labels.Set(newParams))
	}
	message, ok := lo.Find(messages, matcher)
	if !ok {
		return nil, fmt.Errorf("find: not found any message")
	}

	return &message, nil
}

func (d *DefaultManager) BuildReply(msg *v1alpha1.Message, params map[string]any) ([]byte, error) {
	funcMap := sprig.TxtFuncMap()
	for k, f := range d.funcMap {
		funcMap[k] = f
	}

	buf := bytes.NewBufferString("")
	t := template.Must(template.New("msg").Funcs(funcMap).Parse(msg.Spec.Reply.YAMLString()))
	if err := t.Execute(buf, params); err != nil {
		return nil, fmt.Errorf("render message: %w", err)
	}

	var reply v1alpha1.Reply
	if err := yaml.Unmarshal(buf.Bytes(), &reply); err != nil {
		return nil, err
	}

	return []byte(reply.String()), nil
}
