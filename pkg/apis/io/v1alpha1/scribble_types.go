package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ScribbleList is a list of Scribble structs
type ScribbleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Scribble `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Scribble defines the main structure for the custom-resource
type Scribble struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ScribbleSpec `json:"spec"`
}

// ScribbleSpec defines the structure of the Scribble JSON object from the CR
type ScribbleSpec struct {
	Name string `json:"name"`
}

func init() {
	SchemeBuilder.Register(&Scribble{}, &ScribbleList{})
}
