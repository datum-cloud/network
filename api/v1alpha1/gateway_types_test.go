package v1alpha1

import (
	"encoding/json"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newTestGateway() *NetworkGateway {
	return &NetworkGateway{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "network.datumapis.com/v1alpha1",
			Kind:       "NetworkGateway",
		},
		ObjectMeta: metav1.ObjectMeta{Name: "test-gateway"},
		Spec: NetworkGatewaySpec{
			TargetRef: TargetRef{Kind: "Node", Name: "gw-node-a"},
		},
		Status: NetworkGatewayStatus{
			SRv6Address: "2001:db8:1::1",
		},
	}
}

// TestNetworkGatewayDeepCopy verifies that DeepCopy produces an independent
// copy: mutations to the copy must not affect the original.
func TestNetworkGatewayDeepCopy(t *testing.T) {
	orig := newTestGateway()
	dup := orig.DeepCopy()

	dup.Spec.TargetRef.Name = "gw-node-b"
	dup.Status.SRv6Address = "2001:db8:1::2"
	dup.Status.Conditions = append(dup.Status.Conditions, metav1.Condition{Type: ConditionTypeReady})

	if orig.Spec.TargetRef.Name != "gw-node-a" {
		t.Errorf("TargetRef.Name mutated: got %q", orig.Spec.TargetRef.Name)
	}
	if orig.Status.SRv6Address != "2001:db8:1::1" {
		t.Errorf("SRv6Address mutated: got %q", orig.Status.SRv6Address)
	}
	if len(orig.Status.Conditions) != 0 {
		t.Errorf("Conditions mutated: got %v", orig.Status.Conditions)
	}
}

// TestNetworkGatewayDeepCopyNil verifies DeepCopy on a nil pointer returns nil.
func TestNetworkGatewayDeepCopyNil(t *testing.T) {
	var g *NetworkGateway
	if g.DeepCopy() != nil {
		t.Error("DeepCopy on nil pointer should return nil")
	}
}

// TestNetworkGatewayJSONRoundTrip verifies that the struct serialises and
// deserialises through JSON without data loss.
func TestNetworkGatewayJSONRoundTrip(t *testing.T) {
	orig := newTestGateway()
	orig.Status.Conditions = []metav1.Condition{
		{Type: ConditionTypeReady, Status: metav1.ConditionTrue, Reason: "Ready", Message: "ok"},
	}

	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var got NetworkGateway
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if got.Spec.TargetRef != orig.Spec.TargetRef {
		t.Errorf("TargetRef: got %+v, want %+v", got.Spec.TargetRef, orig.Spec.TargetRef)
	}
	if got.Status.SRv6Address != orig.Status.SRv6Address {
		t.Errorf("SRv6Address: got %q, want %q", got.Status.SRv6Address, orig.Status.SRv6Address)
	}
	if len(got.Status.Conditions) != 1 {
		t.Fatalf("Conditions len: got %d, want 1", len(got.Status.Conditions))
	}
}

// TestNetworkGatewayListDeepCopy verifies that NetworkGatewayList.DeepCopy
// produces independent copies of each item.
func TestNetworkGatewayListDeepCopy(t *testing.T) {
	list := &NetworkGatewayList{
		Items: []NetworkGateway{*newTestGateway()},
	}
	copied := list.DeepCopy()
	copied.Items[0].Spec.TargetRef.Name = "other-node"

	if list.Items[0].Spec.TargetRef.Name != "gw-node-a" {
		t.Errorf("original list item mutated via copy")
	}
}

// TestNetworkGatewayFieldNames verifies the JSON keys for spec/status fields
// match the CRD schema.
func TestNetworkGatewayFieldNames(t *testing.T) {
	orig := newTestGateway()

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
	if _, ok := spec["targetRef"]; !ok {
		t.Errorf("expected spec.targetRef field, got %v", spec)
	}

	status, ok := m["status"].(map[string]any)
	if !ok {
		t.Fatalf("status not found or wrong type: %v", m["status"])
	}
	if v, ok := status["sRv6Address"]; !ok || v != "2001:db8:1::1" {
		t.Errorf("expected status.sRv6Address=%q, got %v", "2001:db8:1::1", status["sRv6Address"])
	}
}
