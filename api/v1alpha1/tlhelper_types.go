/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TlhelperSpec defines the desired state of Tlhelper
type TlhelperSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The image of the Ingress Controller.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Image Image `json:"image"`

	// Replicas indicate the replicas to mantain
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=1
	Replicas *int32 `json:"replicas"`

	// The type of the Service for the Ingress Controller. Valid Service types are: ClusterIP and LoadBalancer.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=ClusterIP;LoadBalancer
	// +kubebuilder:default:=ClusterIP
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ServiceType string `json:"serviceType"`

	// Stdout log format. Valid formats are: text and json
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=text;json
	// +kubebuilder:default:=json
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Format string `json:"format"`

	// Log level for V logs.
	// Valit level formats: info, warn, debug
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=info;warn;debug
	// +kubebuilder:default:=info
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	LogLevel string `json:"logLevel"`

	// Initial values of the TLhelper ConfigMap.
	// +kubebuilder:validation:Optional
	// +nullable
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ConfigMapData map[string]string `json:"configMapData,omitempty"`

	// DB is remote or localhost
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	Remote bool `json:"remote"`
}

// Image defines the Repository, Tag and ImagePullPolicy of the Ingress Controller Image.
type Image struct {
	// The repository of the image.
	Repository string `json:"repository"`
	// The tag (version) of the image.
	Tag string `json:"tag"`
	// The ImagePullPolicy of the image.
	// +kubebuilder:validation:Enum=Never;Always;IfNotPresent
	// +kubebuilder:default:=Always
	PullPolicy string `json:"pullPolicy"`
}

// TlhelperStatus defines the observed state of Tlhelper
type TlhelperStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Pods are the name of the Pods hosting the App
	Pods []string `json:"pods"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Tlhelper is the Schema for the tlhelpers API
type Tlhelper struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TlhelperSpec   `json:"spec,omitempty"`
	Status TlhelperStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TlhelperList contains a list of Tlhelper
type TlhelperList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Tlhelper `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Tlhelper{}, &TlhelperList{})
}
