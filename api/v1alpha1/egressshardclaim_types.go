package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EgressShardClaim is one attachment's standing request for internet egress
// on the node it landed on. It is written by the controller that owns the
// cell and read by the node.
//
// The claim is the contract between the two. The cell controller creates one
// when an attachment's network declares egress and the attachment has reported
// its node, and deletes it when the declaration is withdrawn or the attachment
// goes. The node's installer lists the claims naming it on every sweep and
// keeps each VRF's egress route in step: a VRF with a claim routes toward the
// node's shard, a VRF without one does not. That is what lets egress be turned
// on or off for a running workload without re-attaching it.
//
// The claim also records the binding. Status names the shard on the node, so
// the answer to "which shard does this attachment leave through" is readable,
// and a node without a usable shard produces a condition a consumer can see.
//
// There is one claim per attachment, owned by it, so an attachment that goes
// takes its claim with it. The claim names no selector, no address and no
// pool: the node is the binding, and the claim writes it down.
//
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=egressclaim
// +kubebuilder:printcolumn:name="ATTACHMENT",type="string",JSONPath=".spec.attachment.name"
// +kubebuilder:printcolumn:name="VPC",type="string",JSONPath=".spec.vpc.name"
// +kubebuilder:printcolumn:name="NODE",type="string",JSONPath=".spec.nodeName"
// +kubebuilder:printcolumn:name="SHARD",type="string",JSONPath=".status.shardRef.name"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="REASON",type="string",JSONPath=`.status.conditions[?(@.type=="Ready")].reason`
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
type EgressShardClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +required
	Spec EgressShardClaimSpec `json:"spec"`

	// +optional
	Status EgressShardClaimStatus `json:"status,omitempty"`
}

// Label keys carried by an EgressShardClaim so that the two readers who list
// claims can narrow the query server-side. A label restates a fact the spec or
// status already holds; the spec and status are what a reader trusts.
const (
	// LabelEgressShardClaimNode restates spec.nodeName. A node's installer
	// lists the claims carrying its own name to learn which of its VRFs
	// declare egress.
	LabelEgressShardClaimNode string = "network.datumapis.com/egress-node"

	// LabelEgressShardClaimShard names the shard the claim is bound to, the
	// value being the shard's name. A shard holds no list of the attachments
	// it serves, so "what does this shard serve" is answered by listing the
	// claims carrying this label.
	LabelEgressShardClaimShard string = "network.datumapis.com/egress-shard"
)

// EgressShardClaimSpec is the attachment, the VRF, the node and the families
// one claim stands for.
//
// The whole spec is immutable. An attachment that lands on a different node
// is a different request, so the claim is replaced rather than edited.
//
// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="spec is immutable; an attachment that moved nodes gets a new claim"
type EgressShardClaimSpec struct {
	// Attachment is the attachment this claim stands for. It is in the
	// claim's own namespace and owns the claim.
	// +required
	Attachment EgressShardClaimAttachmentRef `json:"attachment"`

	// VPC is the VPC the attachment is on. The node derives the VRF it
	// programs from this name, the same way it does at attach time, so the
	// node needs nothing else to find the routing table this claim governs.
	// +required
	VPC EgressShardClaimVPCRef `json:"vpc"`

	// NodeName is the node the attachment landed on, and therefore the node
	// whose shard serves it and whose installer acts on this claim.
	// +kubebuilder:validation:MinLength=1
	// +required
	NodeName string `json:"nodeName"`

	// Families are the destination address families the attachment's network
	// declared, so the shard on the node is one that translates them.
	// +listType=set
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=2
	// +required
	Families []EgressAddressFamily `json:"families"`
}

// EgressShardClaimAttachmentRef names the attachment a claim stands for.
type EgressShardClaimAttachmentRef struct {
	// Name of the attachment, in the claim's namespace.
	// +kubebuilder:validation:MinLength=1
	// +required
	Name string `json:"name"`
}

// EgressShardClaimVPCRef names the VPC an attachment is on.
type EgressShardClaimVPCRef struct {
	// Name of the VPC, in the claim's namespace.
	// +kubebuilder:validation:MinLength=1
	// +required
	Name string `json:"name"`
}

// EgressAddressFamily is a destination address family an egress shard
// translates toward.
//
// +kubebuilder:validation:Enum=IPv6;IPv4
type EgressAddressFamily string

const (
	// EgressAddressFamilyIPv6 is reached by NAT66 through the shard's IPv6
	// address.
	EgressAddressFamilyIPv6 EgressAddressFamily = "IPv6"

	// EgressAddressFamilyIPv4 is reached by NAT64 through the shard's IPv4
	// address.
	EgressAddressFamilyIPv4 EgressAddressFamily = "IPv4"
)

// EgressShardClaimShardRef names the shard a claim is bound to.
type EgressShardClaimShardRef struct {
	// Namespace of the EgressShard.
	// +kubebuilder:validation:MinLength=1
	// +required
	Namespace string `json:"namespace"`

	// Name of the EgressShard.
	// +kubebuilder:validation:MinLength=1
	// +required
	Name string `json:"name"`
}

// EgressShardClaimStatus is the shard a claim was bound to.
type EgressShardClaimStatus struct {
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// ShardRef is the shard on the attachment's node.
	//
	// Absent means the node holds no shard this claim can bind to, which is
	// what an attachment on a node an operator has not commissioned reads.
	// +optional
	ShardRef *EgressShardClaimShardRef `json:"shardRef,omitempty"`
}

// Reasons reported on an EgressShardClaim's Ready condition.
const (
	// EgressShardClaimReasonBound means this attachment egresses through the
	// shard status names.
	EgressShardClaimReasonBound string = "Bound"

	// EgressShardClaimReasonNoShardOnNode means no shard names the node the
	// attachment landed on.
	EgressShardClaimReasonNoShardOnNode string = "NoShardOnNode"

	// EgressShardClaimReasonShardNotReady means the shard on the node has not
	// reported the identifier a node routes toward.
	EgressShardClaimReasonShardNotReady string = "ShardNotReady"

	// EgressShardClaimReasonShardMismatch means the shard's spec and the
	// identity its process reported disagree, so which one the node runs is
	// unknown and nothing is bound to it.
	EgressShardClaimReasonShardMismatch string = "ShardMismatch"

	// EgressShardClaimReasonFamilyUnsupported means the shard on the node
	// translates none of a family the network declared.
	EgressShardClaimReasonFamilyUnsupported string = "FamilyUnsupported"

	// EgressShardClaimReasonShardMissing means the bound shard no longer
	// exists. The node's instances lost their egress with it.
	EgressShardClaimReasonShardMissing string = "ShardMissing"

	// EgressShardClaimReasonShardTerminating means the bound shard is being
	// deleted. The binding stands, and the shard is held until the claim goes.
	EgressShardClaimReasonShardTerminating string = "ShardTerminating"
)

// +kubebuilder:object:root=true

// EgressShardClaimList contains a list of EgressShardClaim.
type EgressShardClaimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EgressShardClaim `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EgressShardClaim{}, &EgressShardClaimList{})
}
