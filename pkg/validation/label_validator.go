package validation

import (
	"fmt"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// labelValidator is a container for validating the labels of pods
type labelValidator struct {
	Logger logrus.FieldLogger
}

// labelValidator implements the podValidator interface
var _ podValidator = (*labelValidator)(nil)

// Name returns the name of labelValidator
func (n labelValidator) Name() string {
	return "label_validator"
}


// Validate inspects the labels of a given pod and returns validation.
// Validation is only valid if the pod labels exist
func (n labelValidator) Validate(pod *corev1.Pod) (validation, error) {
	if pod.Labels == nil {
		v := validation{
			Valid:  false,
			Reason: fmt.Sprintln("pod labels isn't exist"),
		}
		return v, nil
	}

	return validation{Valid: true, Reason: "valid name"}, nil
}
