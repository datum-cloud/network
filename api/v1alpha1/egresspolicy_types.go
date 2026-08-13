package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NetworkEgressPolicy enables internet egress for a single tenant
// VPC/VPCAttachment, served by the shared hyperconverged gateway engine's
// masquerade (SNAT/PAT) datapath. Unlike NetworkRule, it carries no
// VIP/backend/port: egress is on or off for a (vpcRef, vpcAttachmentRef)
// pair, existence-implies-enabled, not a per-flow rule — because the
// destination of an egress flow is an arbitrary internet address, not a
// pre-configured backend list.
//
// It is namespaced (deployed to galactic-system) and tenant-writable; like
// NetworkRule, vpcRef/vpcAttachmentRef are opaque string identifiers because
// the VPC API is owned by a separate companion operator, not this repo. An
// admission webhook (implemented by the consuming controller) must verify
// the requester is authorized for vpcRef/vpcAttachmentRef before a policy is
// accepted — see the Accepted condition.
//
// Presence of an accepted NetworkEgressPolicy resolves only *enablement*
// (should this tenant reach the egress datapath at all) — a routing-layer
// decision (does the tenant's VRF have a default route toward the shared
// egress_sid locator), not a per-packet datapath lookup. *Isolation*
// (preventing two tenants with colliding ULA source addresses from
// colliding in the egress connection table) is a separate, datapath-level
// concern resolved by tagging each flow with the tenant/VRF identifier
// carried in the egress_sid locator's own Argument bits, not by anything in
// this spec.
//
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=netegress
// +kubebuilder:printcolumn:name="VPC",type="string",JSONPath=".spec.vpcRef"
// +kubebuilder:printcolumn:name="VPC-ATTACHMENT",type="string",JSONPath=".spec.vpcAttachmentRef"
// +kubebuilder:printcolumn:name="ASSIGNED-NODE",type="string",JSONPath=".status.assignedGatewayNode"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
type NetworkEgressPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetworkEgressPolicySpec   `json:"spec,omitempty"`
	Status NetworkEgressPolicyStatus `json:"status,omitempty"`
}

// NetworkEgressPolicySpec defines the desired egress-enablement state for a
// tenant VPC/VPCAttachment.
type NetworkEgressPolicySpec struct {
	// VPCRef is the opaque identifier of the target VPC this policy applies
	// to. This repo does not own the VPC API and does not validate the
	// identifier beyond non-emptiness; the admission webhook of the
	// consuming controller is responsible for verifying the requester is
	// authorized for this VPC before the policy is accepted.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	VPCRef string `json:"vpcRef"`

	// VPCAttachmentRef is the opaque identifier of the target
	// VPCAttachment this policy applies to. Like VPCRef, this is an opaque
	// string reference validated by the admission webhook, not by this API.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	VPCAttachmentRef string `json:"vpcAttachmentRef"`
}

// NetworkEgressPolicyStatus defines the observed state of a
// NetworkEgressPolicy.
type NetworkEgressPolicyStatus struct {
	// ObservedGeneration is the .metadata.generation this status was computed from.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// AssignedGatewayNode is the name of the NetworkGateway-backed gateway
	// node this policy's tenant should route egress traffic through,
	// mirroring NetworkRule's own status.primaryNode field and computed
	// the same way: assigned_node = hash(vpcRef) % <gateway node count>
	// (design plan §4.5 — a tenant's egress node and its primary ingress
	// node are the same node, by design, so both fields are computed by
	// the identical AssignPrimaryNode function). The controller consuming
	// this CRD sets this field exactly once, at creation.
	//
	// This value must never be silently recomputed by a reconciler once
	// set, for the exact same reason NetworkRuleStatus.PrimaryNode's own
	// doc comment gives: recomputing it on a later reconcile can flip
	// which node a tenant's egress traffic routes through and cause an
	// avoidable traffic flap; a reconciler that observes a stale or
	// removed node here must surface that via a condition instead of
	// overwriting the value.
	// +optional
	AssignedGatewayNode string `json:"assignedGatewayNode,omitempty"`

	// Conditions contains the standard conditions for this resource,
	// including Accepted (see AcceptedReasonOwnershipVerified /
	// AcceptedReasonOwnershipDenied in rule_types.go, reused as-is here).
	//
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// NetworkEgressPolicyList is a list of NetworkEgressPolicy resources.
// +kubebuilder:object:root=true
type NetworkEgressPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NetworkEgressPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NetworkEgressPolicy{}, &NetworkEgressPolicyList{})
}
