package validation

import (
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"

)

func TestValidatePod(t *testing.T) {
	v := NewValidator(logger())

	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "lgudax",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "lgudax",
				Image: "busybox",
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

	val, err := v.ValidatePod(pod)
	assert.Nil(t, err)
	assert.True(t, val.Valid)
}

func logger() *logrus.Entry {
	mute := logrus.StandardLogger()
	mute.Out = io.Discard
	return mute.WithField("logger", "test")
}
