package messenger

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// RequirementOperator defines supported operators for requirements
type RequirementOperator string

const (
	// OperatorPrefix matches if field value starts with the specified prefix
	OperatorPrefix RequirementOperator = "Prefix"
	// OperatorSuffix matches if field value ends with the specified suffix
	OperatorSuffix RequirementOperator = "Suffix"
	// OperatorContains matches if field value contains the specified substring
	OperatorContains RequirementOperator = "Contains"
)

// CustomRequirement represents a custom matching requirement
type CustomRequirement struct {
	Key      string
	Operator RequirementOperator
	Value    string
}

// internalSelector extends the k8s selector with custom operators
type internalSelector struct {
	standardSelector   labels.Selector
	customRequirements []CustomRequirement
}

// NewInternalSelector creates a selector that supports both standard k8s selectors
// and custom operators
func NewInternalSelector(labelSelector *metav1.LabelSelector) (*internalSelector, error) {
	// Create a copy of the labelSelector without custom operations
	standardSelector := &metav1.LabelSelector{
		MatchLabels: labelSelector.MatchLabels,
	}

	// Extract custom requirements from matchExpressions
	var customReqs []CustomRequirement
	if labelSelector.MatchExpressions != nil {
		for _, expr := range labelSelector.MatchExpressions {
			// Check for custom operators by convention - operators starting with uppercase
			if len(expr.Values) > 0 && isCustomOperator(string(expr.Operator)) {
				customReqs = append(customReqs, CustomRequirement{
					Key:      expr.Key,
					Operator: RequirementOperator(expr.Operator),
					Value:    expr.Values[0], // Use first value
				})
			} else {
				// Add to standard expressions
				if standardSelector.MatchExpressions == nil {
					standardSelector.MatchExpressions = []metav1.LabelSelectorRequirement{}
				}
				standardSelector.MatchExpressions = append(standardSelector.MatchExpressions, expr)
			}
		}
	}

	// Process standard k8s selector
	selector, err := metav1.LabelSelectorAsSelector(standardSelector)
	if err != nil {
		return nil, err
	}

	return &internalSelector{
		standardSelector:   selector,
		customRequirements: customReqs,
	}, nil
}

// Matches returns true if this selector matches the given set of labels.
func (s *internalSelector) Matches(labelSet labels.Set) bool {
	// First use standard k8s selector to filter out most labels
	if !s.standardSelector.Matches(labelSet) {
		return false
	}

	// If there are custom requirements, check them
	for _, r := range s.customRequirements {
		value, ok := labelSet[r.Key]
		if !ok {
			return false
		}

		switch r.Operator {
		case OperatorPrefix:
			if !strings.HasPrefix(value, r.Value) {
				return false
			}
		case OperatorSuffix:
			if !strings.HasSuffix(value, r.Value) {
				return false
			}
		case OperatorContains:
			if !strings.Contains(value, r.Value) {
				return false
			}
		default:
			return false
		}
	}

	return true
}

// isCustomOperator determines if an operator string represents a custom operator
func isCustomOperator(op string) bool {
	switch RequirementOperator(op) {
	case OperatorPrefix, OperatorSuffix, OperatorContains:
		return true
	default:
		return false
	}
}
