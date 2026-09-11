package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EgressShard marks a single node as a member of the sharded, stateful egress
// translation tier (galactic-nat) — a component deliberately kept off the
// ingress load-balancer's own consistent-hash ring (see NetworkGateway):
// tenant egress traffic (backend -> arbitrary internet destination) is a
// different traffic pattern from ingress (fixed VIP, fixed backend pool)
// and needs its own placement ring, own per-flow state, and its own
// self-routing return path, entirely independent of any NetworkGateway node.
//
// A shard serves one or both address families. NAT66 (IPv6 -> IPv6) and NAT64
// (IPv6 -> IPv4, RFC 6146) are the same function — stateful egress PAT with a
// VRF-scoped session table — over different families, so one shard object
// describes both rather than there being a second, near-duplicate kind.
// Status.ShardAddressIPv6 and Status.ShardAddressIPv4 are each set only for
// the family this shard actually translates; a shard serving only NAT66
// leaves the IPv4 field empty and behaves exactly as it did before NAT64
// existed.
//
// Every shard owns a dedicated, publicly-routable address per family it
// serves, and a flow's allocated masquerade port lives within it — so a reply
// is delivered to the correct shard by ordinary unicast routing alone, with
// no hashing or cross-shard lookup on the return path at all (the "any node
// can determine the owning shard from the tuple alone" property, satisfied by
// construction rather than by a replicated hash table).
//
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=egressshard
// +kubebuilder:printcolumn:name="TARGET",type="string",JSONPath=".spec.targetRef.name"
// +kubebuilder:printcolumn:name="SHARD-SID",type="string",JSONPath=".status.shardSID"
// +kubebuilder:printcolumn:name="IPV6",type="string",JSONPath=".status.shardAddressIPv6"
// +kubebuilder:printcolumn:name="IPV4",type="string",JSONPath=".status.shardAddressIPv4"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
type EgressShard struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EgressShardSpec   `json:"spec,omitempty"`
	Status EgressShardStatus `json:"status,omitempty"`
}

// EgressShardSpec defines the desired state of an EgressShard.
type EgressShardSpec struct {
	// TargetRef identifies the Node this shard executes on.
	// +kubebuilder:validation:Required
	TargetRef TargetRef `json:"targetRef"`
}

// EgressShardStatus defines the observed state of an EgressShard.
//
// Every field here is echoed from what the shard's datapath process was
// actually started with, not derived: the shard publishes what it is running,
// so a status that disagrees with an operator's intent is a visible
// misconfiguration rather than a silently reconciled one.
type EgressShardStatus struct {
	// ObservedGeneration is the .metadata.generation this status was computed from.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// ShardSID is this shard's own uSID locator — a real SRv6 uSID (unlike
	// the ShardAddress fields, which are plain routable addresses),
	// advertised into BGP the same way any other node-reachability route is
	// (a /128 BGPAdvertisement, no VRFID/Function) so every other node learns
	// a kernel SEG6 route toward it before installing a tenant VRF's egress
	// route against it. One SID serves both families: which translation a
	// packet gets is decided from the inner destination, not from a second
	// SID.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == '' || (isIP(self) && ip(self).family() == 6)",message="shardSID must be a valid IPv6 address"
	ShardSID string `json:"shardSID,omitempty"`

	// ShardAddressIPv6 is this shard's own dedicated, publicly-routable IPv6
	// address — every NAT66 masquerade port this shard allocates lives within
	// it, so any node can route a reply to the correct shard using ordinary
	// unicast routing on this address alone, with no per-flow state lookup
	// anywhere but the owning shard itself. Operator-supplied per shard today
	// (no in-cluster derivation mechanism yet — the same gap
	// BGPRouter.Spec.SRv6Locator/NodeID assignment has today).
	//
	// Empty means this shard does not perform IPv6-to-IPv6 translation.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == '' || (isIP(self) && ip(self).family() == 6)",message="shardAddressIPv6 must be a valid IPv6 address"
	ShardAddressIPv6 string `json:"shardAddressIPv6,omitempty"`

	// ShardAddressIPv4 is this shard's own dedicated, publicly-routable IPv4
	// address — every NAT64 masquerade port this shard allocates lives within
	// it, and it is the source an IPv4-only destination sees. Unlike
	// ShardAddressIPv6, reachability for this address is not established by a
	// BGPAdvertisement into the EVPN fabric: an IPv4 reply arrives from the
	// internet, so the address must be attracted to this node by the underlay
	// or upstream announcement instead. Publishing it here is what makes that
	// operator prerequisite checkable.
	//
	// Empty means this shard does not perform NAT64.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == '' || (isIP(self) && ip(self).family() == 4)",message="shardAddressIPv4 must be a valid IPv4 address"
	ShardAddressIPv4 string `json:"shardAddressIPv4,omitempty"`

	// NAT64Prefix is the IPv6 prefix whose synthesized addresses this shard
	// translates to IPv4 — one Datum-operated Network-Specific Prefix, shared
	// fabric-wide, never per-tenant. It is echoed here, rather than only
	// existing as process configuration, because it is the single fact DNS64
	// synthesis has to agree with: a shard translating for a different prefix
	// than the resolver synthesizes into is otherwise a silent blackhole.
	//
	// Empty whenever ShardAddressIPv4 is empty.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == '' || isCIDR(self)",message="nat64Prefix must be a valid CIDR"
	NAT64Prefix string `json:"nat64Prefix,omitempty"`

	// Conditions contains the standard conditions for this resource.
	//
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// EgressShardList is a list of EgressShard resources.
// +kubebuilder:object:root=true
type EgressShardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EgressShard `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EgressShard{}, &EgressShardList{})
}
