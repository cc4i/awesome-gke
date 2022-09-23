/*
Copyright 2022.

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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TrackerTopSpec defines the desired state of TrackerTop
type TrackerTopSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Specified a namespace to provison resources
	Where string `json:"where"`

	// Tracker will be placed as per topology
	Trackers []Tracker `json:"trackers"`

	// Topology is to define relations between Trackers
	Graph []Topology `json:"graph,omitempty"`

	// Shared Redis for kv store
	Redis ThirdParty `json:"redis"`
}

// ServingType describe how to expose Tracker service, basically it's same as normal Service
// +kubebuilder:validation:Enum=ClusterIP;LoadBalancer;NodePort
type ServingType string

const (
	ClusterIP    ServingType = "ClusterIP"
	LoadBalancer ServingType = "LoadBalancer"
	NodePort     ServingType = "NodePort"
)

type ThirdParty struct {
	Name string `json:"name"`
	// Dependent container image for Tracker
	Image    string `json:"image"`
	Host     string `json:"host"`
	Port     int32  `json:"port"`
	Protocol string `json:"protocol,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

type Tracker struct {
	// Service name for Tracker
	Name string `json:"name"`
	// Dependent container image for Tracker
	Image string `json:"image"`
	// Verison of Tracker
	Version string `json:"version"`
	// Replicas of Tracker
	Replicas *int32 `json:"replicas"`
	// Expose URI, eg: http://host:port/path
	ServingUri string `json:"servingUri,omitempty"`
	// Service protocol
	ServingProtocol string `json:"servingProtocol,omitempty"` //eg: HTTP, HTTPS, TCP, GRPC
	// Service Type
	ServingType ServingType `json:"servingType,omitempty"` //eg: ClusterIP, LoadBalancer, NodePort
	// Connection to Redis
	RedisConn string `json:"redisConn,omitempty"`
	// Where to host service
	HostedCloud string `json:"hostedCloud,omitempty"`
}

type Topology struct {
	Name       string   `json:"name"`
	Upstream   string   `json:"upstream,omitempty"`
	Downstream []string `json:"downstream,omitempty"`
	// Condition in header to determine which downstream would be called
	ServingCondition map[string]string `json:"servingCondition,omitempty"`
}

// TrackerTopStatus defines the observed state of TrackerTop
type TrackerTopStatus struct {
	Active           []v1.ObjectReference `json:"active,omitempty"`
	LastScheduleTime *metav1.Time         `json:"lastScheduleTime,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TrackerTop is the Schema for the trackertops API
type TrackerTop struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TrackerTopSpec   `json:"spec,omitempty"`
	Status TrackerTopStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TrackerTopList contains a list of TrackerTop
type TrackerTopList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TrackerTop `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TrackerTop{}, &TrackerTopList{})
}
