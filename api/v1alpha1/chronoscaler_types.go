package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ChronoScalerSpec defines the desired state of ChronoScaler
type ChronoScalerSpec struct {
	// Define the Start hour (24h Local format) to scale the deployment to replicas specified in Replicas
	Start string `json:"start"`
	// Define the End hour (24h Local format) to scale the deployment to replicas specified in DefaultReplicas
	End string `json:"end"`
	// Define the number of Replicas to scale to
	Replicas int32 `json:"replicas"`
	// Define the number of Replicas to scale back after EndTime
	DefaultReplicas int32 `json:"defaultReplicas"`
	// List of Deployments to scale
	Deployments []NamespacedName `json:"deployments"`
}

type NamespacedName struct {
	// Define the Name of deployment to Scale
	Name string `json:"name"`
	//Define the Namespace in which the deployment is present
	Namespace string `json:"namespace"`
}

// ChronoScalerStatus defines the observed state of ChronoScaler
type ChronoScalerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ChronoScaler is the Schema for the chronoscalers API
type ChronoScaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChronoScalerSpec   `json:"spec,omitempty"`
	Status ChronoScalerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ChronoScalerList contains a list of ChronoScaler
type ChronoScalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChronoScaler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChronoScaler{}, &ChronoScalerList{})
}
