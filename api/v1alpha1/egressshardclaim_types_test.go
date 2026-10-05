package v1alpha1

import (
	"encoding/json"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newTestEgressShardClaim() *EgressShardClaim {
	return &EgressShardClaim{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "network.datumapis.com/v1alpha1",
			Kind:       "EgressShardClaim",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "web-0-eth1",
			Namespace: "tenant-a",
			Labels:    map[string]string{LabelEgressShardClaimNode: "node-a"},
		},
		Spec: EgressShardClaimSpec{
			Attachment: EgressShardClaimAttachmentRef{Name: "web-0-eth1"},
			VPC:        EgressShardClaimVPCRef{Name: "vpc-blue"},
			NodeName:   "node-a",
			Families:   []EgressAddressFamily{EgressAddressFamilyIPv6},
		},
		Status: EgressShardClaimStatus{
			ShardRef: &EgressShardClaimShardRef{Namespace: "galactic-system", Name: "node-a"},
			Conditions: []metav1.Condition{{
				Type:   "Ready",
				Status: metav1.ConditionTrue,
				Reason: EgressShardClaimReasonBound,
			}},
		},
	}
}

func TestEgressShardClaimDeepCopy(t *testing.T) {
	orig := newTestEgressShardClaim()
	dup := orig.DeepCopy()

	dup.Spec.Families[0] = EgressAddressFamilyIPv4
	dup.Spec.NodeName = "node-b"
	dup.Status.ShardRef.Name = "node-b"
	dup.Labels[LabelEgressShardClaimNode] = "node-b"

	if orig.Spec.Families[0] != EgressAddressFamilyIPv6 {
		t.Errorf("Families mutated: got %q", orig.Spec.Families[0])
	}
	if orig.Spec.NodeName != "node-a" {
		t.Errorf("NodeName mutated: got %q", orig.Spec.NodeName)
	}
	if orig.Status.ShardRef.Name != "node-a" {
		t.Errorf("ShardRef mutated: got %q", orig.Status.ShardRef.Name)
	}
	if orig.Labels[LabelEgressShardClaimNode] != "node-a" {
		t.Errorf("Labels mutated: got %q", orig.Labels[LabelEgressShardClaimNode])
	}
}

func TestEgressShardClaimDeepCopyNil(t *testing.T) {
	var c *EgressShardClaim
	if c.DeepCopy() != nil {
		t.Error("DeepCopy on nil pointer should return nil")
	}
}

func TestEgressShardClaimJSONRoundTrip(t *testing.T) {
	orig := newTestEgressShardClaim()

	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got EgressShardClaim
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Spec.Attachment.Name != "web-0-eth1" {
		t.Errorf("Attachment.Name = %q", got.Spec.Attachment.Name)
	}
	if got.Spec.VPC.Name != "vpc-blue" {
		t.Errorf("VPC.Name = %q", got.Spec.VPC.Name)
	}
	if got.Spec.NodeName != "node-a" {
		t.Errorf("NodeName = %q", got.Spec.NodeName)
	}
	if len(got.Spec.Families) != 1 || got.Spec.Families[0] != EgressAddressFamilyIPv6 {
		t.Errorf("Families = %v", got.Spec.Families)
	}
	if got.Status.ShardRef == nil || got.Status.ShardRef.Name != "node-a" {
		t.Errorf("ShardRef = %v", got.Status.ShardRef)
	}
}

func TestEgressShardClaimJSONFieldNames(t *testing.T) {
	data, err := json.Marshal(newTestEgressShardClaim())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	spec, ok := raw["spec"].(map[string]any)
	if !ok {
		t.Fatal("spec missing")
	}
	for _, key := range []string{"attachment", "vpc", "nodeName", "families"} {
		if _, present := spec[key]; !present {
			t.Errorf("spec.%s missing from wire form", key)
		}
	}
	status, ok := raw["status"].(map[string]any)
	if !ok {
		t.Fatal("status missing")
	}
	if _, present := status["shardRef"]; !present {
		t.Error("status.shardRef missing from wire form")
	}
}
