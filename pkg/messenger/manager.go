package messenger

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
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
	newParams := flattenMap(params, 5)
	messages, err := d.store.List(context.Background())
	if err != nil {
		return nil, fmt.Errorf("list messages: %w", err)
	}

	matcher := func(item v1alpha1.Message) bool {
		selector, err := NewInternalSelector(&item.Spec.Selector)
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

// flattenMap converts a nested map into a flat map with dot notation.
// Example: {a:{b:{c:1}}} => {a.b.c:"1"}
// maxDepth limits the recursion depth to prevent stack overflow.
// Non-object and non-string values are converted to strings.
func flattenMap(m map[string]any, maxDepth int) map[string]string {
	result := make(map[string]string)

	var flatten func(prefix string, value any, depth int)
	flatten = func(prefix string, value any, depth int) {
		if depth > maxDepth {
			// Convert to string if max depth reached
			result[prefix] = fmt.Sprint(value)
			return
		}

		if nestedMap, ok := value.(map[string]any); ok {
			for k, v := range nestedMap {
				newKey := k
				if prefix != "" {
					newKey = prefix + "." + k
				}
				flatten(newKey, v, depth+1)
			}
		} else if _, isArray := value.([]any); isArray {
			// Ignore array values
			return
		} else if strVal, ok := value.(string); ok {
			result[prefix] = strVal
		} else {
			// Convert non-string, non-map values to string
			result[prefix] = fmt.Sprint(value)
		}
	}

	flatten("", m, 0)

	// Remove the empty prefix key if it exists (shouldn't in normal usage)
	delete(result, "")

	return result
}
