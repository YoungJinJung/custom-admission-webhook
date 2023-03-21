package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestLabelValidator(t *testing.T) {
	t.Run("Success Case", func(t *testing.T) {
		pod := &corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name: "test-pod",
				Labels: map[string]string{
					"tags.datadoghq.com/env": "test-app",
					"app":                    "test-app",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{
					Name: "test-container",
				}},
			},
		}

		v, err := labelValidator{Logger: logger()}.Validate(pod)
		assert.Nil(t, err)
		assert.True(t, v.Valid)
	})

	t.Run("Failed Case: Label key is missing", func(t *testing.T) {
		pod := &corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name: "test-pod",
				Labels: map[string]string{
					"otherKey": "otherValue",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{
					Name: "test-container",
				}},
			},
		}

		v, err := labelValidator{Logger: logger()}.Validate(pod)
		assert.Nil(t, err)
		assert.False(t, v.Valid)
		assert.Equal(t, "Label key must contain \"app\" in Label.", v.Reason)
	})
}
