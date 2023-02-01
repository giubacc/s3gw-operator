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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BucketSpec defines the desired state of Bucket
type BucketSpec struct {
	// Important: Run "make" to regenerate code after modifying this file

	Name          string `json:"name,omitempty"`
	ObjectLocking bool   `json:"objectLocking"`
}

type BStatus string

const (
	NotCreated BStatus = "NotCreated"
	Created    BStatus = "Created"
	Error      BStatus = "Error"
)

// BucketStatus defines the observed state of Bucket
type BucketStatus struct {
	// Important: Run "make" to regenerate code after modifying this file

	Status BStatus `json:"status,omitempty"`
}

// Bucket is the Schema for the buckets API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`,description="Bucket Status"
// +kubebuilder:printcolumn:name="Object-Locking",type=string,JSONPath=`.spec.objectLocking`,description="Object Locking"
// +kubebuilder:printcolumn:name="Creation-Time",type=string,JSONPath=`.metadata.creationTimestamp`,description="Creation-Time"
type Bucket struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BucketSpec   `json:"spec,omitempty"`
	Status BucketStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BucketList contains a list of Bucket
type BucketList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Bucket `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Bucket{}, &BucketList{})
}
