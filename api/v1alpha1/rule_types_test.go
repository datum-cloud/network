package v1alpha1

import (
	"encoding/json"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newTestRule() *NetworkRule {
	return &NetworkRule{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "network.datumapis.com/v1alpha1",
			Kind:       "NetworkRule",
		},
		ObjectMeta: metav1.ObjectMeta{Name: "test-rule", Namespace: "galactic-system"},
		Spec: NetworkRuleSpec{
			VPCRef:           "vpc-a",
			VPCAttachmentRef: "vpcattachment-a",
			VIPAddresses:     []string{"2001:db8:1::10"},
			Protocol:         NetworkRuleProtocolTCP,
			Port:             443,
			Backends: []NetworkRuleBackend{
				{Address: "fd00:10:1::1", Port: 8443},
			},
		},
	}
}

// TestNetworkRuleDeepCopy verifies that DeepCopy produces an independent
// copy: mutations to slices in the copy must not affect the original.
func TestNetworkRuleDeepCopy(t *testing.T) {
	orig := newTestRule()
	dup := orig.DeepCopy()

	dup.Spec.VIPAddresses[0] = "2001:db8:1::20"
	dup.Spec.Backends[0].Address = "fd00:10:1::2"

	if orig.Spec.VIPAddresses[0] != "2001:db8:1::10" {
		t.Errorf("VIPAddresses[0] mutated: got %q", orig.Spec.VIPAddresses[0])
	}
	if orig.Spec.Backends[0].Address != "fd00:10:1::1" {
		t.Errorf("Backends[0].Address mutated: got %q", orig.Spec.Backends[0].Address)
	}
}

// TestNetworkRuleDeepCopyNil verifies DeepCopy on a nil pointer returns nil.
func TestNetworkRuleDeepCopyNil(t *testing.T) {
	var r *NetworkRule
	if r.DeepCopy() != nil {
		t.Error("DeepCopy on nil pointer should return nil")
	}
}

// TestNetworkRuleJSONRoundTrip verifies that the struct serialises and
// deserialises through JSON without data loss.
func TestNetworkRuleJSONRoundTrip(t *testing.T) {
	orig := newTestRule()
	orig.Spec.Backends = append(orig.Spec.Backends, NetworkRuleBackend{Address: "fd00:10:1::3", Port: 8444})
	orig.Status.Conditions = []metav1.Condition{
		{Type: ConditionTypeAccepted, Status: metav1.ConditionTrue, Reason: AcceptedReasonOwnershipVerified},
	}

	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var got NetworkRule
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if got.Spec.VPCRef != orig.Spec.VPCRef {
		t.Errorf("VPCRef: got %q, want %q", got.Spec.VPCRef, orig.Spec.VPCRef)
	}
	if got.Spec.VPCAttachmentRef != orig.Spec.VPCAttachmentRef {
		t.Errorf("VPCAttachmentRef: got %q, want %q", got.Spec.VPCAttachmentRef, orig.Spec.VPCAttachmentRef)
	}
	if len(got.Spec.Backends) != 2 {
		t.Errorf("Backends len: got %d, want 2", len(got.Spec.Backends))
	}
	if len(got.Status.Conditions) != 1 || got.Status.Conditions[0].Reason != AcceptedReasonOwnershipVerified {
		t.Errorf("Conditions: got %v", got.Status.Conditions)
	}
}

// TestNetworkRuleListDeepCopy verifies that NetworkRuleList.DeepCopy
// produces independent copies of each item.
func TestNetworkRuleListDeepCopy(t *testing.T) {
	list := &NetworkRuleList{
		Items: []NetworkRule{*newTestRule()},
	}
	copied := list.DeepCopy()
	copied.Items[0].Spec.VPCRef = "vpc-b"

	if list.Items[0].Spec.VPCRef != "vpc-a" {
		t.Errorf("original list item mutated via copy")
	}
}

// TestNetworkRuleStatusHasNoPrimaryNode is a regression test: the earlier
// active-passive Full-NAT design assigned a single primaryNode per rule.
// This design's anycast/DSR model has every NetworkGateway serve every rule
// identically, so no such field should exist any more.
func TestNetworkRuleStatusHasNoPrimaryNode(t *testing.T) {
	orig := newTestRule()
	orig.Status.Conditions = []metav1.Condition{{Type: ConditionTypeAccepted}}

	data, err := json.Marshal(orig.Status)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if _, ok := m["primaryNode"]; ok {
		t.Errorf("unexpected primaryNode field present in status: %v", m)
	}
}

// TestNetworkRuleBackendFieldNames verifies the JSON keys for
// NetworkRuleBackend match the CRD schema ("address", "port").
func TestNetworkRuleBackendFieldNames(t *testing.T) {
	b := NetworkRuleBackend{Address: "fd00:10:1::1", Port: 8443}
	data, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if v, ok := m["address"]; !ok || v != "fd00:10:1::1" {
		t.Errorf("expected JSON key \"address\"=%q, got %v", "fd00:10:1::1", m)
	}
	if v, ok := m["port"]; !ok || v != float64(8443) {
		t.Errorf("expected JSON key \"port\"=8443, got %v", m)
	}
}

// TestNetworkRuleProtocolValues is a regression test pinning the accepted
// NetworkRuleProtocol enum values.
func TestNetworkRuleProtocolValues(t *testing.T) {
	if NetworkRuleProtocolTCP != "tcp" {
		t.Errorf("NetworkRuleProtocolTCP: got %q, want %q", NetworkRuleProtocolTCP, "tcp")
	}
	if NetworkRuleProtocolUDP != "udp" {
		t.Errorf("NetworkRuleProtocolUDP: got %q, want %q", NetworkRuleProtocolUDP, "udp")
	}
}
