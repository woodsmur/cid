/*


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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CustomImageDeploySpec defines the desired state of CustomImageDeploy
type CustomImageDeploySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Image is the docker image with version info of CustomImageDeploy.
	Image string `json:"image,omitempty"`

	// Size is the number of pods to run
	Size int32 `json:"size"`

	// Port is the port of container
	Port int32 `json:"port"`
}

// CustomImageDeployStatus defines the observed state of CustomImageDeploy
type CustomImageDeployStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// CustomImageDeploy is the Schema for the customimagedeploys API
type CustomImageDeploy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CustomImageDeploySpec   `json:"spec,omitempty"`
	Status CustomImageDeployStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CustomImageDeployList contains a list of CustomImageDeploy
type CustomImageDeployList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CustomImageDeploy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CustomImageDeploy{}, &CustomImageDeployList{})
}
