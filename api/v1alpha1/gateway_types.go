package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NetworkGateway defines an XDP ingress NAT+LB gateway engine instance bound
// to a single dedicated gateway-role node. Exactly one NetworkGateway exists
// per gateway node (spec.targetRef.name is the Kubernetes node name),
// mirroring the BGPRouter node-scoped root object pattern. NetworkRule
// resources are assigned to a NetworkGateway via status.primaryNode.
//
// There is no tunnel overlay in this design (an earlier Geneve-based
// approach was superseded before this type shipped): the gateway's XDP
// program does Full-NAT (DNAT the VIP to a backend Pod's address, SNAT the
// client's source to status.sRv6Address) and pushes an SRv6 uSID outer
// header addressed to the backend's worker node directly, so return traffic
// (addressed to status.sRv6Address) arrives back at this same gateway node
// over the ordinary SRv6 fabric — no compute-node encap agent, no tunnel
// endpoint to publish. status.sRv6Address is advertised into BGP the same
// way any workload prefix is (a BGPAdvertisement naming it, /128, Argument
// 0 — the value PR #740 reserves and forbids registering into any tenant
// VRF, guaranteeing it never collides with a real tenant's Argument), so
// every other node learns a real kernel SEG6 route to it for free through
// the existing EVPN pipeline.
//
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=netgw
// +kubebuilder:printcolumn:name="TARGET",type="string",JSONPath=".spec.targetRef.name"
// +kubebuilder:printcolumn:name="SRV6-ADDRESS",type="string",JSONPath=".status.sRv6Address"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
type NetworkGateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetworkGatewaySpec   `json:"spec,omitempty"`
	Status NetworkGatewayStatus `json:"status,omitempty"`
}

// NetworkGatewaySpec defines the desired state of a NetworkGateway.
type NetworkGatewaySpec struct {
	// TargetRef identifies the Node this gateway engine executes on.
	// +kubebuilder:validation:Required
	TargetRef TargetRef `json:"targetRef"`
}

// NetworkGatewayStatus defines the observed state of a NetworkGateway.
type NetworkGatewayStatus struct {
	// ObservedGeneration is the .metadata.generation this status was computed from.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// SRv6Address is this gateway node's own SRv6-reachable IPv6 address,
	// used as the Full-NAT SNAT source for every ingress flow this node
	// translates. Backend Pods' replies are naturally routed back to it
	// over the ordinary SRv6 fabric (the same mechanism that routes any
	// other node's traffic), where this node's XDP program decapsulates
	// and un-NATs them using its own conn_table — there is no separate
	// tunnel endpoint or overlay device to publish. Populated by the
	// engine once it has computed the address (a uFMT 48+16 uSID over this
	// node's own BGPRouter locator/node-ID, at the reserved Argument 0)
	// and advertised it into BGP.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == '' || isIP(self)",message="sRv6Address must be a valid IPv6 address"
	SRv6Address string `json:"sRv6Address,omitempty"`

	// Conditions contains the standard conditions for this resource.
	//
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// NetworkGatewayList is a list of NetworkGateway resources.
// +kubebuilder:object:root=true
type NetworkGatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NetworkGateway `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NetworkGateway{}, &NetworkGatewayList{})
}
