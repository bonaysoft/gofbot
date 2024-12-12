package v1alpha1

import (
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
