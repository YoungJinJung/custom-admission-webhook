package validation

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// labelValidator is a container for validating the labels of pods
type labelValidator struct {
	Logger logrus.FieldLogger
}

// validLabels is the list of valid label keys for pods
var validLabels = []string{
	"app",
	"tags.datadoghq.com/env",
}

// labelValidator implements the podValidator interface
var _ podValidator = (*labelValidator)(nil)

// Name returns the name of labelValidator
func (l labelValidator) Name() string {
	return "label_validator"
}

// validation is only valid if label exist in pod and key contains "validLabels".
func (l labelValidator) Validate(pod *corev1.Pod) (validation, error) {
	for _, k := range validLabels {
		found := false
		for labelKey := range pod.Labels {
			if strings.EqualFold(k, labelKey) {
				found = true
			}
		}
		if !found {
			validMsg := fmt.Sprintf("Label key must contain \"%s\" in Label.", k)
			v := validation{
				Valid:  false,
				Reason: validMsg,
			}
			return v, nil
		}
	}

	return validation{Valid: true, Reason: successValidationMsg}, nil
}
