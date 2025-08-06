package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(&Cat{}, &CatList{})
}

// Cat is a cat.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type Cat struct {
	metav1.TypeMeta `json:",inline"`
	// metadata contains the standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// spec defines the desired state of the Cat.
	// +required
	Spec CatSpec `json:"spec"`
	// status defines the observed state of the Cat.
	// +optional
	Status *CatStatus `json:"status,omitempty"`
}

// CatList contains a list of Cat.
// +kubebuilder:object:root=true
type CatList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cat `json:"items"`
}

// CatSpec defines the desired state of Cat.
type CatSpec struct {
}

// CatStatus defines the observed state of Cat.
type CatStatus struct {
	// sleepy represents if the cat is sleepy.
	// +required
	// +kubebuilder:example=false
	Sleepy bool `json:"sleepy"`
}
