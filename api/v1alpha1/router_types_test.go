package v1alpha1

import (
	"encoding/json"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newTestRouter(roles ...RouterRole) *BGPRouter {
	return &BGPRouter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "network.datumapis.com/v1alpha1",
			Kind:       "BGPRouter",
		},
		ObjectMeta: metav1.ObjectMeta{Name: "test-router"},
		Spec: BGPRouterSpec{
			TargetRef: TargetRef{Kind: "Node", Name: "node-1"},
			LocalASN:  65000,
			RouterID:  "10.0.0.1",
			Roles:     roles,
		},
	}
}

// TestBGPRouterDeepCopyRoles verifies that mutating Roles on the copy does not
// affect the original.
func TestBGPRouterDeepCopyRoles(t *testing.T) {
	orig := newTestRouter(RouterRoleTransit, RouterRoleFabric)
	dup := orig.DeepCopy()

	dup.Spec.Roles[0] = RouterRoleTenant

	if orig.Spec.Roles[0] != RouterRoleTransit {
		t.Errorf("Roles[0] mutated: got %q, want %q", orig.Spec.Roles[0], RouterRoleTransit)
	}
	if orig.Spec.Roles[1] != RouterRoleFabric {
		t.Errorf("Roles[1] mutated: got %q, want %q", orig.Spec.Roles[1], RouterRoleFabric)
	}
}

// TestBGPRouterDeepCopyNilRoles verifies that a router with no roles deep-copies
// without allocating a non-nil slice.
func TestBGPRouterDeepCopyNilRoles(t *testing.T) {
	orig := newTestRouter()
	dup := orig.DeepCopy()

	if dup.Spec.Roles != nil {
		t.Errorf("expected nil Roles on copy, got %v", dup.Spec.Roles)
	}
}

// TestBGPRouterJSONRoundTripRoles verifies that Roles survives JSON marshal/unmarshal.
func TestBGPRouterJSONRoundTripRoles(t *testing.T) {
	cases := []struct {
		name  string
		roles []RouterRole
	}{
		{"no roles", nil},
		{"single transit", []RouterRole{RouterRoleTransit}},
		{"all roles", []RouterRole{RouterRoleTransit, RouterRoleFabric, RouterRoleTenant}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			orig := newTestRouter(tc.roles...)

			data, err := json.Marshal(orig)
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}

			var got BGPRouter
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("Unmarshal: %v", err)
			}

			if len(got.Spec.Roles) != len(tc.roles) {
				t.Fatalf("Roles len: got %d, want %d", len(got.Spec.Roles), len(tc.roles))
			}
			for i, r := range tc.roles {
				if got.Spec.Roles[i] != r {
					t.Errorf("Roles[%d]: got %q, want %q", i, got.Spec.Roles[i], r)
				}
			}
		})
	}
}

// TestBGPRouterJSONRolesFieldName verifies the JSON key is "roles".
func TestBGPRouterJSONRolesFieldName(t *testing.T) {
	orig := newTestRouter(RouterRoleFabric)
	data, err := json.Marshal(orig.Spec)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	raw, ok := m["roles"]
	if !ok {
		t.Fatal("expected JSON key \"roles\" not found")
	}
	roles, ok := raw.([]any)
	if !ok || len(roles) != 1 || roles[0] != "fabric" {
		t.Errorf("unexpected roles value: %v", raw)
	}
}

// TestBGPRouterRolesOmitEmpty verifies that a router with no roles omits the
// "roles" key from JSON output.
func TestBGPRouterRolesOmitEmpty(t *testing.T) {
	orig := newTestRouter()
	data, err := json.Marshal(orig.Spec)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if _, ok := m["roles"]; ok {
		t.Error("expected \"roles\" key to be absent when empty")
	}
}

// TestBGPRouterLargeLocalASN verifies that 4-byte ASNs (values above signed int32 max)
// survive JSON round-trip correctly.
func TestBGPRouterLargeLocalASN(t *testing.T) {
	// Max 4-byte ASN — the boundary of the uint32 range.
	const maxASN = int64(4294967295)

	router := &BGPRouter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "network.datumapis.com/v1alpha1",
			Kind:       "BGPRouter",
		},
		ObjectMeta: metav1.ObjectMeta{Name: "large-asn-router"},
		Spec: BGPRouterSpec{
			TargetRef: TargetRef{Kind: "Node", Name: "node-1"},
			LocalASN:  maxASN,
			RouterID:  "10.0.0.1",
			Roles:     []RouterRole{RouterRoleTransit},
			AddressFamilies: []AddressFamily{
				{AFI: AFIIPv4, SAFI: SAFIUnicast},
			},
		},
	}

	data, err := json.Marshal(router)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var got BGPRouter
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if got.Spec.LocalASN != maxASN {
		t.Errorf("LocalASN after round-trip: got %d, want %d", got.Spec.LocalASN, maxASN)
	}
}

// TestBGPRouterLocalASNAboveSignedInt32Max verifies that values > 2^31-1
// (the signed int32 ceiling) are handled correctly.
func TestBGPRouterLocalASNAboveSignedInt32Max(t *testing.T) {
	// 2^31 = 2147483648 — the first value beyond signed int32 range.
	const aboveSignedMax int64 = 2147483648

	router := &BGPRouter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "network.datumapis.com/v1alpha1",
			Kind:       "BGPRouter",
		},
		ObjectMeta: metav1.ObjectMeta{Name: "above-int32-router"},
		Spec: BGPRouterSpec{
			TargetRef: TargetRef{Kind: "Node", Name: "node-1"},
			LocalASN:  aboveSignedMax,
			RouterID:  "10.0.0.1",
			Roles:     []RouterRole{RouterRoleTransit},
			AddressFamilies: []AddressFamily{
				{AFI: AFIIPv4, SAFI: SAFIUnicast},
			},
		},
	}

	data, err := json.Marshal(router)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var got BGPRouter
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if got.Spec.LocalASN != aboveSignedMax {
		t.Errorf("LocalASN after round-trip: got %d, want %d", got.Spec.LocalASN, aboveSignedMax)
	}
}

// TestBGPRouterDeepCopyLocalASN verifies that mutating LocalASN on the copy
// does not affect the original.
func TestBGPRouterDeepCopyLocalASN(t *testing.T) {
	orig := newTestRouter()
	orig.Spec.LocalASN = 654321
	dup := orig.DeepCopy()

	dup.Spec.LocalASN = 999999
	if orig.Spec.LocalASN != 654321 {
		t.Errorf("LocalASN mutated: got %d, want 654321", orig.Spec.LocalASN)
	}
}

// TestBGPRouterJSONRoundTripSRv6Locator verifies that SRv6Locator survives
// JSON marshal/unmarshal.
func TestBGPRouterJSONRoundTripSRv6Locator(t *testing.T) {
	orig := newTestRouter(RouterRoleFabric)
	orig.Spec.SRv6Locator = "2001:db8:ff01::/48"

	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var got BGPRouter
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if got.Spec.SRv6Locator != orig.Spec.SRv6Locator {
		t.Errorf("SRv6Locator: got %q, want %q", got.Spec.SRv6Locator, orig.Spec.SRv6Locator)
	}
}

// TestBGPRouterSRv6LocatorFieldName verifies the JSON key is "srv6Locator".
func TestBGPRouterSRv6LocatorFieldName(t *testing.T) {
	orig := newTestRouter(RouterRoleFabric)
	orig.Spec.SRv6Locator = "2001:db8:ff01::/48"

	data, err := json.Marshal(orig.Spec)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if raw, ok := m["srv6Locator"]; !ok || raw != "2001:db8:ff01::/48" {
		t.Errorf("unexpected srv6Locator value: %v", raw)
	}
}

// TestBGPRouterSRv6LocatorOmitEmpty verifies that a router without a
// SRv6Locator omits the "srv6Locator" key from JSON output.
func TestBGPRouterSRv6LocatorOmitEmpty(t *testing.T) {
	orig := newTestRouter(RouterRoleFabric)

	data, err := json.Marshal(orig.Spec)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if _, ok := m["srv6Locator"]; ok {
		t.Error("expected \"srv6Locator\" key to be absent when empty")
	}
}

// TestBGPRouterJSONRoundTripNodeID verifies that NodeID survives JSON
// marshal/unmarshal.
func TestBGPRouterJSONRoundTripNodeID(t *testing.T) {
	orig := newTestRouter(RouterRoleFabric)
	orig.Spec.NodeID = 42

	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var got BGPRouter
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if got.Spec.NodeID != orig.Spec.NodeID {
		t.Errorf("NodeID: got %d, want %d", got.Spec.NodeID, orig.Spec.NodeID)
	}
}

// TestBGPRouterNodeIDFieldName verifies the JSON key is "nodeID".
func TestBGPRouterNodeIDFieldName(t *testing.T) {
	orig := newTestRouter(RouterRoleFabric)
	orig.Spec.NodeID = 42

	data, err := json.Marshal(orig.Spec)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if raw, ok := m["nodeID"]; !ok || raw != float64(42) {
		t.Errorf("unexpected nodeID value: %v", raw)
	}
}

// TestBGPRouterNodeIDOmitEmpty verifies that a router without a NodeID omits
// the "nodeID" key from JSON output.
func TestBGPRouterNodeIDOmitEmpty(t *testing.T) {
	orig := newTestRouter(RouterRoleFabric)

	data, err := json.Marshal(orig.Spec)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if _, ok := m["nodeID"]; ok {
		t.Error("expected \"nodeID\" key to be absent when empty")
	}
}
