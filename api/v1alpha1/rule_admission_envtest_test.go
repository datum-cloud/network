// Copyright 2026 Datum Cloud, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package v1alpha1_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/yaml"
)

func TestNetworkRuleBackendSelectionAdmission(t *testing.T) {
	if os.Getenv("KUBEBUILDER_ASSETS") == "" {
		t.Skip("set KUBEBUILDER_ASSETS for isolated API admission tests")
	}
	apiRoot := os.Getenv("GALACTIC_TEST_NETWORK_API")
	if apiRoot == "" {
		apiRoot = "../.."
	}
	data, err := os.ReadFile(filepath.Join(apiRoot, "config", "crd", "network.datumapis.com_networkrules.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	crd := &apiextensionsv1.CustomResourceDefinition{}
	if err := yaml.Unmarshal(data, crd); err != nil {
		t.Fatal(err)
	}
	environment := &envtest.Environment{CRDs: []*apiextensionsv1.CustomResourceDefinition{crd}}
	cfg, err := environment.Start()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := environment.Stop(); err != nil {
			t.Errorf("stop API: %v", err)
		}
	})
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	c, err := client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if err := c.Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "rule-admission"}}); err != nil {
		t.Fatal(err)
	}

	for _, tt := range []struct {
		name      string
		selector  map[string]any
		addresses any
		wantError bool
	}{
		{name: "neither-selection", wantError: true},
		{name: "empty-selector", selector: map[string]any{}, wantError: true},
		{name: "selector-only", selector: map[string]any{"matchLabels": map[string]any{"app": "web"}}},
		{name: "addresses-only", addresses: []any{"fd00:10::1"}},
		{name: "explicit-empty-addresses", addresses: []any{}},
		{name: "empty-selector-with-addresses", selector: map[string]any{}, addresses: []any{"fd00:10::1"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			spec := map[string]any{
				"vpcRef": "vpc-a", "vipAddresses": []any{"2001:db8::10"},
				"protocol": "tcp", "port": int64(443), "backendPort": int64(8443),
			}
			if tt.selector != nil {
				spec["backendSelector"] = tt.selector
			}
			if tt.addresses != nil {
				spec["backendAddresses"] = tt.addresses
			}
			rule := &unstructured.Unstructured{Object: map[string]any{
				"apiVersion": "network.datumapis.com/v1alpha1", "kind": "NetworkRule",
				"metadata": map[string]any{"name": tt.name, "namespace": "rule-admission"},
				"spec":     spec,
			}}
			err := c.Create(ctx, rule)
			if (err != nil) != tt.wantError {
				t.Fatalf("create error=%v, wantError=%t", err, tt.wantError)
			}
			if tt.wantError && !apierrors.IsInvalid(err) {
				t.Fatalf("expected schema rejection, got %v", err)
			}
		})
	}
}
