package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// NewScribble returns a new Scribble instance with the given name
func NewScribble(name string) *Scribble {
	return &Scribble{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}
