// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestKubernetesClusterNodePoolSchemaAcceptsUbuntu2404(t *testing.T) {
	validate := resourceKubernetesClusterNodePoolSchema()["os_sku"].ValidateFunc
	if validate == nil {
		t.Fatal("node pool os_sku ValidateFunc is nil")
	}

	if warnings, errs := validate("Ubuntu2404", "os_sku"); len(warnings) != 0 || len(errs) != 0 {
		t.Fatalf("expected Ubuntu2404 to be accepted for node pool os_sku, got warnings=%v errs=%v", warnings, errs)
	}
}

func TestKubernetesClusterDefaultNodePoolSchemaAcceptsUbuntu2404(t *testing.T) {
	defaultNodePool := resourceKubernetesCluster().Schema["default_node_pool"]
	if defaultNodePool == nil {
		t.Fatal("default_node_pool schema is nil")
	}

	defaultNodePoolResource, ok := defaultNodePool.Elem.(*pluginsdk.Resource)
	if !ok {
		t.Fatalf("expected default_node_pool elem to be a resource, got %T", defaultNodePool.Elem)
	}

	validate := defaultNodePoolResource.Schema["os_sku"].ValidateFunc
	if validate == nil {
		t.Fatal("default node pool os_sku ValidateFunc is nil")
	}

	if warnings, errs := validate("Ubuntu2404", "os_sku"); len(warnings) != 0 || len(errs) != 0 {
		t.Fatalf("expected Ubuntu2404 to be accepted for default node pool os_sku, got warnings=%v errs=%v", warnings, errs)
	}
}
