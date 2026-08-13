package v1alpha1

import (
	"encoding/json"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newTestEgressPolicy() *NetworkEgressPolicy {
	return &NetworkEgressPolicy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "network.datumapis.com/v1alpha1",
			Kind:       "NetworkEgressPolicy",
		},
		ObjectMeta: metav1.ObjectMeta{Name: "test-egress-policy", Namespace: "galactic-system"},
		Spec: NetworkEgressPolicySpec{
			VPCRef:           "vpc-a",
			VPCAttachmentRef: "vpcattachment-a",
		},
	}
}

// TestNetworkEgressPolicyDeepCopy verifies that DeepCopy produces an
// independent copy: mutations to the copy must not affect the original.
func TestNetworkEgressPolicyDeepCopy(t *testing.T) {
	orig := newTestEgressPolicy()
	dup := orig.DeepCopy()

	dup.Spec.VPCRef = "vpc-b"
	dup.Spec.VPCAttachmentRef = "vpcattachment-b"
	dup.Status.Conditions = append(dup.Status.Conditions, metav1.Condition{Type: ConditionTypeAccepted})

	if orig.Spec.VPCRef != "vpc-a" {
		t.Errorf("VPCRef mutated: got %q", orig.Spec.VPCRef)
	}
	if orig.Spec.VPCAttachmentRef != "vpcattachment-a" {
		t.Errorf("VPCAttachmentRef mutated: got %q", orig.Spec.VPCAttachmentRef)
	}
	if len(orig.Status.Conditions) != 0 {
		t.Errorf("Conditions mutated: got %v", orig.Status.Conditions)
	}
}

// TestNetworkEgressPolicyDeepCopyNil verifies DeepCopy on a nil pointer
// returns nil.
func TestNetworkEgressPolicyDeepCopyNil(t *testing.T) {
	var p *NetworkEgressPolicy
	if p.DeepCopy() != nil {
		t.Error("DeepCopy on nil pointer should return nil")
	}
}

// TestNetworkEgressPolicyJSONRoundTrip verifies that the struct serialises
// and deserialises through JSON without data loss.
func TestNetworkEgressPolicyJSONRoundTrip(t *testing.T) {
	orig := newTestEgressPolicy()
	orig.Status.Conditions = []metav1.Condition{
		{Type: ConditionTypeAccepted, Status: metav1.ConditionTrue, Reason: AcceptedReasonOwnershipVerified},
	}

	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var got NetworkEgressPolicy
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if got.Spec.VPCRef != orig.Spec.VPCRef {
		t.Errorf("VPCRef: got %q, want %q", got.Spec.VPCRef, orig.Spec.VPCRef)
	}
	if got.Spec.VPCAttachmentRef != orig.Spec.VPCAttachmentRef {
		t.Errorf("VPCAttachmentRef: got %q, want %q", got.Spec.VPCAttachmentRef, orig.Spec.VPCAttachmentRef)
	}
	if len(got.Status.Conditions) != 1 || got.Status.Conditions[0].Reason != AcceptedReasonOwnershipVerified {
		t.Errorf("Conditions: got %v", got.Status.Conditions)
	}
}

// TestNetworkEgressPolicyListDeepCopy verifies that
// NetworkEgressPolicyList.DeepCopy produces independent copies of each item.
func TestNetworkEgressPolicyListDeepCopy(t *testing.T) {
	list := &NetworkEgressPolicyList{
		Items: []NetworkEgressPolicy{*newTestEgressPolicy()},
	}
	copied := list.DeepCopy()
	copied.Items[0].Spec.VPCRef = "vpc-b"

	if list.Items[0].Spec.VPCRef != "vpc-a" {
		t.Errorf("original list item mutated via copy")
	}
}

// TestNetworkEgressPolicyFieldNames verifies the JSON keys for spec fields
// match the CRD schema ("vpcRef", "vpcAttachmentRef") and that no
// backend/port/VIP fields exist — egress enablement is on/off for a
// (vpcRef, vpcAttachmentRef) pair, not a per-flow rule.
func TestNetworkEgressPolicyFieldNames(t *testing.T) {
	orig := newTestEgressPolicy()

	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	spec, ok := m["spec"].(map[string]any)
	if !ok {
		t.Fatalf("spec not found or wrong type: %v", m["spec"])
	}
	if v, ok := spec["vpcRef"]; !ok || v != "vpc-a" {
		t.Errorf("expected spec.vpcRef=%q, got %v", "vpc-a", spec["vpcRef"])
	}
	if v, ok := spec["vpcAttachmentRef"]; !ok || v != "vpcattachment-a" {
		t.Errorf("expected spec.vpcAttachmentRef=%q, got %v", "vpcattachment-a", spec["vpcAttachmentRef"])
	}
	for _, unexpected := range []string{"vipAddresses", "protocol", "port", "backends"} {
		if _, ok := spec[unexpected]; ok {
			t.Errorf("unexpected field %q present in spec: %v", unexpected, spec)
		}
	}
}

// TestNetworkGatewayEgressAddressFieldName verifies NetworkGatewayStatus's
// new EgressAddress field round-trips under the JSON key "egressAddress"
// and stays independently settable from SRv6Address (a gateway node may
// have one, both, or neither populated).
func TestNetworkGatewayEgressAddressFieldName(t *testing.T) {
	gw := newTestGateway()
	gw.Status.EgressAddress = "2001:db8:ffff::1"

	data, err := json.Marshal(gw)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	status, ok := m["status"].(map[string]any)
	if !ok {
		t.Fatalf("status not found or wrong type: %v", m["status"])
	}
	if v, ok := status["egressAddress"]; !ok || v != "2001:db8:ffff::1" {
		t.Errorf("expected status.egressAddress=%q, got %v", "2001:db8:ffff::1", status["egressAddress"])
	}
	if v, ok := status["sRv6Address"]; !ok || v != "2001:db8:1::1" {
		t.Errorf("expected status.sRv6Address unaffected, got %v", v)
	}
}

// TestNetworkGatewayEgressSIDFieldName verifies NetworkGatewayStatus's new
// EgressSID field round-trips under the JSON key "egressSID" (design plan
// §3.1/§4.3) and stays independently settable from SRv6Address/
// EgressAddress.
func TestNetworkGatewayEgressSIDFieldName(t *testing.T) {
	gw := newTestGateway()
	gw.Status.EgressAddress = "2001:db8:ffff::1"
	gw.Status.EgressSID = "2001:db8:eeee::1"

	data, err := json.Marshal(gw)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	status, ok := m["status"].(map[string]any)
	if !ok {
		t.Fatalf("status not found or wrong type: %v", m["status"])
	}
	if v, ok := status["egressSID"]; !ok || v != "2001:db8:eeee::1" {
		t.Errorf("expected status.egressSID=%q, got %v", "2001:db8:eeee::1", status["egressSID"])
	}
	if v, ok := status["egressAddress"]; !ok || v != "2001:db8:ffff::1" {
		t.Errorf("expected status.egressAddress unaffected, got %v", v)
	}
}

// TestNetworkEgressPolicyAssignedGatewayNodeFieldName verifies the new
// AssignedGatewayNode field round-trips under the JSON key
// "assignedGatewayNode" (design plan §4.5).
func TestNetworkEgressPolicyAssignedGatewayNodeFieldName(t *testing.T) {
	orig := newTestEgressPolicy()
	orig.Status.AssignedGatewayNode = "gw-node-a"

	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	status, ok := m["status"].(map[string]any)
	if !ok {
		t.Fatalf("status not found or wrong type: %v", m["status"])
	}
	if v, ok := status["assignedGatewayNode"]; !ok || v != "gw-node-a" {
		t.Errorf("expected status.assignedGatewayNode=%q, got %v", "gw-node-a", status["assignedGatewayNode"])
	}
}

// TestNetworkEgressPolicyDeepCopyIncludesAssignedGatewayNode extends
// TestNetworkEgressPolicyDeepCopy to cover the new field: mutating the copy
// must not affect the original.
func TestNetworkEgressPolicyDeepCopyIncludesAssignedGatewayNode(t *testing.T) {
	orig := newTestEgressPolicy()
	orig.Status.AssignedGatewayNode = "gw-node-a"

	dup := orig.DeepCopy()
	dup.Status.AssignedGatewayNode = "gw-node-b"

	if orig.Status.AssignedGatewayNode != "gw-node-a" {
		t.Errorf("AssignedGatewayNode mutated: got %q", orig.Status.AssignedGatewayNode)
	}
}
