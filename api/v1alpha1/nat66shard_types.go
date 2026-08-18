package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NAT66Shard marks a single node as a member of the sharded, stateful NAT66
// egress tier (galactic-nat66) — a component deliberately kept off the
// ingress load-balancer's own consistent-hash ring (see NetworkGateway):
// tenant egress traffic (backend -> arbitrary internet destination) is a
// different traffic pattern from ingress (fixed VIP, fixed backend pool)
// and needs its own placement ring, own per-flow state, and its own
// self-routing return path, entirely independent of any NetworkGateway node.
//
// Every shard owns a dedicated, BGP-advertised public IPv6 address
// (Status.ShardAddress) that a flow's allocated masquerade port lives
// within — so a reply is delivered to the correct shard by ordinary
// unicast SRv6/BGP routing alone, with no hashing or cross-shard lookup on
// the return path at all (the "any node can determine the owning shard from
// the tuple alone" property, satisfied by construction rather than by a
// replicated hash table).
//
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=nat66shard
// +kubebuilder:printcolumn:name="TARGET",type="string",JSONPath=".spec.targetRef.name"
// +kubebuilder:printcolumn:name="SHARD-ADDRESS",type="string",JSONPath=".status.shardAddress"
// +kubebuilder:printcolumn:name="SHARD-SID",type="string",JSONPath=".status.shardSID"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
type NAT66Shard struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NAT66ShardSpec   `json:"spec,omitempty"`
	Status NAT66ShardStatus `json:"status,omitempty"`
}

// NAT66ShardSpec defines the desired state of a NAT66Shard.
type NAT66ShardSpec struct {
	// TargetRef identifies the Node this shard executes on.
	// +kubebuilder:validation:Required
	TargetRef TargetRef `json:"targetRef"`
}

// NAT66ShardStatus defines the observed state of a NAT66Shard.
type NAT66ShardStatus struct {
	// ObservedGeneration is the .metadata.generation this status was computed from.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// ShardAddress is this shard's own dedicated, publicly-routable IPv6
	// address — every masquerade port this shard allocates lives within it,
	// so any node can route a reply to the correct shard using ordinary
	// unicast routing on this address alone, with no per-flow state lookup
	// anywhere but the owning shard itself. Operator-supplied per shard
	// today (no in-cluster derivation mechanism yet — the same gap
	// BGPRouter.Spec.SRv6Locator/NodeID assignment has today).
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == '' || isIP(self)",message="shardAddress must be a valid IPv6 address"
	ShardAddress string `json:"shardAddress,omitempty"`

	// ShardSID is this shard's own uSID locator — a real SRv6 uSID (unlike
	// ShardAddress, a plain routable address), advertised into BGP the same
	// way any other node-reachability route is (a /128 BGPAdvertisement, no
	// VRFID/Function) so every other node learns a kernel SEG6 route toward
	// it before installing a tenant VRF's default egress route against it.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == '' || isIP(self)",message="shardSID must be a valid IPv6 address"
	ShardSID string `json:"shardSID,omitempty"`

	// Conditions contains the standard conditions for this resource.
	//
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// NAT66ShardList is a list of NAT66Shard resources.
// +kubebuilder:object:root=true
type NAT66ShardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NAT66Shard `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NAT66Shard{}, &NAT66ShardList{})
}
