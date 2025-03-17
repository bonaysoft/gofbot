package messenger

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func TestIsCustomOperator(t *testing.T) {
	tests := []struct {
		operator string
		expected bool
	}{
		{"In", false},
		{"NotIn", false},
		{"Exists", false},
		{"DoesNotExist", false},
		{"Prefix", true},
		{"Suffix", true},
		{"Contains", true},
	}

	for _, test := range tests {
		if result := isCustomOperator(test.operator); result != test.expected {
			t.Errorf("isCustomOperator(%s) = %v; expected %v", test.operator, result, test.expected)
		}
	}
}

func TestInternalSelector_Matches(t *testing.T) {
	tests := []struct {
		name           string
		labelSelector  *metav1.LabelSelector
		labelSet       labels.Set
		expectedResult bool
	}{
		{
			name: "only standard - match",
			labelSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "test"},
			},
			labelSet:       labels.Set{"app": "test"},
			expectedResult: true,
		},
		{
			name: "only standard - no match",
			labelSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "test"},
			},
			labelSet:       labels.Set{"app": "not-test"},
			expectedResult: false,
		},
		{
			name: "only custom - match",
			labelSelector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "version",
						Operator: metav1.LabelSelectorOperator(OperatorPrefix),
						Values:   []string{"v1"},
					},
				},
			},
			labelSet:       labels.Set{"version": "v1.0.0"},
			expectedResult: true,
		},
		{
			name: "only custom - no match",
			labelSelector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "version",
						Operator: metav1.LabelSelectorOperator(OperatorPrefix),
						Values:   []string{"v1"},
					},
				},
			},
			labelSet:       labels.Set{"version": "v2.0.0"},
			expectedResult: false,
		},
		{
			name: "standard and custom - match",
			labelSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "test"},
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "version",
						Operator: metav1.LabelSelectorOperator(OperatorPrefix),
						Values:   []string{"v1"},
					},
				},
			},
			labelSet:       labels.Set{"app": "test", "version": "v1.0.0"},
			expectedResult: true,
		},
		{
			name: "standard and custom - no match",
			labelSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "test"},
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "version",
						Operator: metav1.LabelSelectorOperator(OperatorPrefix),
						Values:   []string{"v1"},
					},
				},
			},
			labelSet:       labels.Set{"app": "test", "version": "v2.0.0"},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selector, err := NewInternalSelector(tt.labelSelector)
			if err != nil {
				t.Fatalf("failed to create selector: %v", err)
			}
			result := selector.Matches(tt.labelSet)
			if result != tt.expectedResult {
				t.Errorf("expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}
