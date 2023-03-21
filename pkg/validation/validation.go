package validation

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

const (
	successValidationMsg = "Pod Validation"
)

type Validator struct {
	Logger *logrus.Entry
}

func NewValidator(logger *logrus.Entry) *Validator {
	return &Validator{Logger: logger}
}

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
	validations := []podValidator{
		resourceValidator{v.Logger},
		labelValidator{v.Logger},
	}

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

	return validation{Valid: true, Reason: successValidationMsg}, nil
}
