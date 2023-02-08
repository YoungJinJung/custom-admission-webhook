package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestNameValidatorValidate(t *testing.T) {
	t.Run("Success Case", func(t *testing.T) {
		pod := &corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name: "lgudax",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{
					Name:  "lgudax",
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU: resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("500Mi"),
						},
						Limits: corev1.ResourceList{
							corev1.ResourceCPU: resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("500Mi"),
						},
					},
				}},
			},
		}

		v, err := resourceValidator{logger()}.Validate(pod)
		assert.Nil(t, err)
		assert.True(t, v.Valid)
	})

	t.Run("Failed Case#1 Missing resource", func(t *testing.T) {
		pod := &corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name: "lgudax",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{
					Name:  "lgudax",
				}},
			},
		}

		v, err := resourceValidator{logger()}.Validate(pod)
		assert.Nil(t, err)
		assert.False(t, v.Valid)
	})

	t.Run("Failed Case#1 zero assign resource", func(t *testing.T) {
		pod := &corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name: "lgudax",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{
					Name:  "lgudax",
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU: resource.MustParse("0m"),
							corev1.ResourceMemory: resource.MustParse("0Mi"),
						},
						Limits: corev1.ResourceList{
							corev1.ResourceCPU: resource.MustParse("0m"),
							corev1.ResourceMemory: resource.MustParse("500Mi"),
						},
					},
				}},
			},
		}

		v, err := resourceValidator{logger()}.Validate(pod)
		assert.Nil(t, err)
		assert.False(t, v.Valid)
	})
}
