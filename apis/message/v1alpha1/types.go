package v1alpha1

import (
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Message struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec MessageSpec `json:"spec"`
}

type MessageSpec struct {
	Selector metav1.LabelSelector `json:"selector"`
	Reply    Reply                `json:"reply"`
}

type Reply struct {
	Text string `json:"text"`
	JSON any    `json:"json"`
}

func (r *Reply) YAMLString() string {
	out, err := yaml.Marshal(r)
	if err != nil {
		return ""
	}

	// yaml.Marshal Wrap the string in single quotes,
	// but we don't need it because after wrapping, characters such as newlines cannot be parsed
	return strings.ReplaceAll(string(out), "'", "")
}

func (r *Reply) String() string {
	if r.Text != "" {
		return r.Text
	}

	v, err := json.Marshal(r.JSON)
	if err != nil {
		return ""
	}
	return string(v)
}
