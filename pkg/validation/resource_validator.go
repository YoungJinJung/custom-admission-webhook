package validation

import (
	"fmt"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// resourceValidator is a container for validating the name of pods
type resourceValidator struct {
	Logger logrus.FieldLogger
}

// resourceValidator implements the podValidator interface
var _ podValidator = (*resourceValidator)(nil)

// Name returns the name of resourceValidator
func (n resourceValidator) Name() string {
	return "resource_validator"
}

// Validate inspects the name of a given pod and returns validation.
// The returned validation is only valid if the pod name does not contain some
// bad string.
func (n resourceValidator) Validate(pod *corev1.Pod) (validation, error) {
	for _, container := range pod.Spec.Containers {
		cpuSetMsg, cpuNonZeroMsg  := validateResource(container.Resources.Limits, "limit", corev1.ResourceCPU)
		if cpuNonZeroMsg != "" {
			v := validation{
				Valid:  false,
				Reason: cpuNonZeroMsg,
			}
			return v, nil
		}
		if cpuSetMsg != "" {
			v := validation{
				Valid:  false,
				Reason: cpuSetMsg,
			}
			return v, nil
		}
	}

	return validation{Valid: true, Reason: "valid resource"}, nil
}


func validateResource(resList corev1.ResourceList, resourceName string, name corev1.ResourceName) (string, string){
	if !isResourceSet(resList, name) {
		msg := fmt.Sprintf("'%s' resource %s must be specified.", name, resourceName)
		return msg, ""
	}
	if !isResourceNonZero(resList, name) {
		msg := fmt.Sprintf("'%s' resource %s must be a nonzero value.", name, resourceName)
		return "", msg
	}
	return "", ""
}

func isResourceSet(resList corev1.ResourceList, name corev1.ResourceName) bool {
	var missing = resList == nil
	if !missing {
		if _, ok := resList[name]; !ok {
			missing = true
		}
	}
	return !missing
}

func isResourceNonZero(resList corev1.ResourceList, name corev1.ResourceName) bool {
	if resList == nil {
		return true
	}
	if r, ok := resList[name]; ok {
		return !r.IsZero()
	} else {
		return true
	}
}