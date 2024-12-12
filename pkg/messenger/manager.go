package messenger

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/samber/lo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/bonaysoft/gofbot/apis/message/v1alpha1"
	"github.com/bonaysoft/gofbot/pkg/storage"
)

type Manager interface {
	Match(params map[string]any) (*v1alpha1.Message, bool)

	BuildReply(msg *v1alpha1.Message, params map[string]any) ([]byte, error)
}

type DefaultManager struct {
	store storage.Manager
}

func NewDefaultManager(store storage.Manager) *DefaultManager {
	return &DefaultManager{store: store}
}

func (d *DefaultManager) Match(params map[string]any) (*v1alpha1.Message, bool) {
	messages, err := d.store.List(context.Background())
	if err != nil {
		return nil, false
	}

	// 目前只支持value为string的kv，后续可以把复杂结构拍平成一维的，比如{a:{b:{c:1}}}=> {a.b.c:1}
	newParams := lo.MapEntries(lo.PickBy(params, func(key string, value any) bool {
		_, ok := value.(string)
		return ok
	}), func(key string, value any) (string, string) {
		return key, value.(string)
	})

	message, ok := lo.Find(messages, func(item v1alpha1.Message) bool {
		selector, err := metav1.LabelSelectorAsSelector(&item.Spec.Selector)
		if err != nil {
			return false
		}

		return selector.Matches(labels.Set(newParams))
	})
	return &message, ok
}

func (d *DefaultManager) BuildReply(msg *v1alpha1.Message, params map[string]any) ([]byte, error) {
	msgTemplate := msg.Spec.Reply.Text
	if msg.Spec.Reply.JSON != nil {
		data, err := json.Marshal(msg.Spec.Reply.JSON)
		if err != nil {
			return nil, err
		}

		msgTemplate = string(data)
	}

	funcMap := sprig.TxtFuncMap()
	// for k, f := range lark.FuncMap() {
	// 	funcMap[k] = f
	// }

	buf := bytes.NewBufferString("")
	t := template.Must(template.New("msg").Funcs(funcMap).Parse(msgTemplate))
	if err := t.Execute(buf, params); err != nil {
		return nil, err
	}

	newMsg := buf.String()
	if strconv.CanBackquote(newMsg) {
		return buf.Bytes(), nil
	}

	result := strconv.Quote(newMsg)
	return []byte(result)[1 : len(result)-1], nil
}
