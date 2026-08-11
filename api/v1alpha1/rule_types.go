package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NetworkRuleProtocol is the transport protocol matched by a NetworkRule's
// ingress VIP.
//
// +kubebuilder:validation:Enum=tcp;udp
type NetworkRuleProtocol string

const (
	// NetworkRuleProtocolTCP matches TCP traffic.
	NetworkRuleProtocolTCP NetworkRuleProtocol = "tcp"

	// NetworkRuleProtocolUDP matches UDP traffic.
	NetworkRuleProtocolUDP NetworkRuleProtocol = "udp"
)

// Accepted reasons — used as Accepted.Reason on NetworkRule, set by the
// admission webhook that verifies the requester is authorized for the
// vpcRef/vpcAttachmentRef named in the rule.
const (
	// AcceptedReasonOwnershipVerified indicates admission verified the
	// requester is authorized for the target VPC/VPCAttachment.
	AcceptedReasonOwnershipVerified string = "OwnershipVerified"

	// AcceptedReasonOwnershipDenied indicates admission rejected the rule
	// because the requester is not authorized for the target
	// VPC/VPCAttachment named in vpcRef/vpcAttachmentRef.
	AcceptedReasonOwnershipDenied string = "OwnershipDenied"
)

// NetworkRuleBackend is a single backend endpoint that ingress traffic
// matching a NetworkRule's VIP addresses is load-balanced to.
type NetworkRuleBackend struct {
	// Address is the backend's IPv4 or IPv6 address.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=45
	// +kubebuilder:validation:XValidation:rule="isIP(self)",message="address must be a valid IPv4 or IPv6 address"
	Address string `json:"address"`

	// Port is the backend's destination port.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port"`
}

// NetworkRule defines ingress load-balancing and NAT for a single tenant
// VPC/VPCAttachment, served by the shared hyperconverged gateway engine.
// It is namespaced (deployed to galactic-system) and tenant-writable; the
// vpcRef/vpcAttachmentRef fields are opaque string identifiers because the
// VPC API is owned by a separate companion operator, not this repo. An
// admission webhook (implemented by the consuming controller) must verify
// the requester is authorized for vpcRef/vpcAttachmentRef before a rule is
// accepted — see the Accepted condition.
//
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=netrule
// +kubebuilder:printcolumn:name="VPC",type="string",JSONPath=".spec.vpcRef"
// +kubebuilder:printcolumn:name="PROTOCOL",type="string",JSONPath=".spec.protocol"
// +kubebuilder:printcolumn:name="PORT",type="integer",JSONPath=".spec.port"
// +kubebuilder:printcolumn:name="PRIMARY-NODE",type="string",JSONPath=".status.primaryNode"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
type NetworkRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetworkRuleSpec   `json:"spec,omitempty"`
	Status NetworkRuleStatus `json:"status,omitempty"`
}

// NetworkRuleSpec defines the desired ingress load-balancing state for a
// tenant VPC/VPCAttachment.
type NetworkRuleSpec struct {
	// VPCRef is the opaque identifier of the target VPC this rule applies
	// to. This repo does not own the VPC API and does not validate the
	// identifier beyond non-emptiness; the admission webhook of the
	// consuming controller is responsible for verifying the requester is
	// authorized for this VPC before the rule is accepted.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	VPCRef string `json:"vpcRef"`

	// VPCAttachmentRef is the opaque identifier of the target
	// VPCAttachment this rule applies to. Like VPCRef, this is an opaque
	// string reference validated by the admission webhook, not by this API.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	VPCAttachmentRef string `json:"vpcAttachmentRef"`

	// VIPAddresses is the list of ingress VIP addresses (IPv4 and/or IPv6)
	// this rule provisions on the assigned gateway node(s).
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=8
	// +kubebuilder:validation:items:MaxLength=45
	// +kubebuilder:validation:XValidation:rule="self.all(v, isIP(v))",message="vipAddresses must all be valid IPv4 or IPv6 addresses"
	// +listType=set
	VIPAddresses []string `json:"vipAddresses"`

	// Protocol is the transport protocol matched by VIPAddresses/Port.
	// +kubebuilder:validation:Required
	Protocol NetworkRuleProtocol `json:"protocol"`

	// Port is the ingress port on VIPAddresses that this rule load-balances.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port"`

	// Backends is the list of backend address:port targets that ingress
	// traffic matching VIPAddresses/Protocol/Port is load-balanced to.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=64
	Backends []NetworkRuleBackend `json:"backends"`
}

// NetworkRuleStatus defines the observed state of a NetworkRule.
type NetworkRuleStatus struct {
	// ObservedGeneration is the .metadata.generation this status was computed from.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// PrimaryNode is the name of the NetworkGateway-backed gateway node
	// assigned to advertise this rule's VIPAddresses at the preferred BGP
	// local-preference, per the active-active model: primary_node =
	// hash(vpcRef) % <gateway node count>. The controller consuming this
	// CRD sets this field exactly once, at creation.
	//
	// This value must never be silently recomputed by a reconciler once
	// set. Recomputing it on a later reconcile can flip which node is
	// primary for a live VIP and cause an avoidable traffic flap; a
	// reconciler that observes a stale or removed node here must surface
	// that via a condition instead of overwriting the value.
	// +optional
	PrimaryNode string `json:"primaryNode,omitempty"`

	// Conditions contains the standard conditions for this resource.
	//
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// NetworkRuleList is a list of NetworkRule resources.
// +kubebuilder:object:root=true
type NetworkRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NetworkRule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NetworkRule{}, &NetworkRuleList{})
}
