package validation

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// Validator is a container for mutation
type Validator struct {
	Logger *logrus.Entry
}

// NewValidator returns an initialised instance of Validator
func NewValidator(logger *logrus.Entry) *Validator {
	return &Validator{Logger: logger}
}

// podValidators is an interface used to group functions mutating pods
type podValidator interface {
	Validate(*corev1.Pod) (validation, error)
	Name() string
}

type validation struct {
	Valid  bool
	Reason string
}

// ValidatePod returns true if a pod is valid
func (v *Validator) ValidatePod(pod *corev1.Pod) (validation, error) {
	// list of all validations to be applied to the pod
	validations := []podValidator{
		resourceValidator{v.Logger},
	}

	// apply all validations
	for _, v := range validations {
		var err error
		validPod, err := v.Validate(pod)
		if err != nil {
			return validation{Valid: false, Reason: err.Error()}, err
		}
		if !validPod.Valid {
			return validation{Valid: false, Reason: validPod.Reason}, err
		}
	}

	return validation{Valid: true, Reason: "valid pod"}, nil
}
