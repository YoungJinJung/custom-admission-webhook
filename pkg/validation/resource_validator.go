package validation

import (
	"fmt"
	"strings"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)
const (
	empty = ""
	limits = "limits"
	reqeusts = "reqeusts"
	delimiter = "."
)
// resourceValidator is a container for validating the resources of pods
type resourceValidator struct {
	Logger logrus.FieldLogger
}

// resourceValidator implements the podValidator interface
var _ podValidator = (*resourceValidator)(nil)

// Name returns the name of resourceValidator
func (n resourceValidator) Name() string {
	return "resource_validator"
}

var resources = [4]corev1.ResourceName{
	corev1.ResourceRequestsCPU,
	corev1.ResourceLimitsCPU,
	corev1.ResourceRequestsMemory,
	corev1.ResourceLimitsMemory,
}

// Validate inspects the resources(requests, limits) of a given pod and returns validation.
// validation is only valid if Resourse exist pod and values is not zero.
func (n resourceValidator) Validate(pod *corev1.Pod) (validation, error) {
	var validMsg string
	var resourceList corev1.ResourceList
	var resourceName corev1.ResourceName

	for _, container := range pod.Spec.Containers {
		for index, resource := range resources {
			arr := strings.Split(resource.String(), delimiter)
			if (strings.Compare(arr[0], limits) == 0){
				resourceList = container.Resources.Limits
			} else {
				resourceList = container.Resources.Requests
			}

			if index < 2 {
				resourceName = corev1.ResourceCPU
			} else {
				resourceName = corev1.ResourceMemory
			}
			validMsg = validateResource(resourceList, reqeusts, resourceName)
			if validMsg != empty {
				v := validation{
					Valid:  false,
					Reason: validMsg,
				}
				return v, nil
			}
		}		
	}

	return validation{Valid: true, Reason: successValidationMsg}, nil
}


func validateResource(resList corev1.ResourceList, resourceName string, name corev1.ResourceName) string{
	if !isEmpty(resList, name) {
		msg := fmt.Sprintf("'%s' resource %s must be specified.", name, resourceName)
		return msg
	}
	if !isNonZero(resList, name) {
		msg := fmt.Sprintf("'%s' resource %s must be a nonzero value.", name, resourceName)
		return msg
	}
	return empty
}

func isEmpty(resList corev1.ResourceList, name corev1.ResourceName) bool {
	var missing = resList == nil
	if !missing {
		if _, ok := resList[name]; !ok {
			missing = true
		}
	}
	return !missing
}

func isNonZero(resList corev1.ResourceList, name corev1.ResourceName) bool {
	if resList == nil {
		return true
	}
	if r, ok := resList[name]; ok {
		return !r.IsZero()
	} else {
		return true
	}
}