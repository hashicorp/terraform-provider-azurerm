// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package generators

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestInferIdentityProperties_CompositeResource(t *testing.T) {
	// Test against the actual nat_gateway_public_ip_association_resource.go
	filePath := filepath.Join(findProviderRoot(), "internal", "services", "network", "nat_gateway_public_ip_association_resource.go")
	inferred, err := InferIdentityProperties(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedProps := []string{"resource_id1", "resource_id2"}
	if !slices.Equal(inferred.Properties, expectedProps) {
		t.Fatalf("expected properties %v, got %v", expectedProps, inferred.Properties)
	}
	if inferred.HasSubscriptionID {
		t.Fatalf("expected HasSubscriptionID to be false for composite resource")
	}
	if inferred.IsVirtual {
		t.Fatalf("expected IsVirtual to be false")
	}
}

func TestInferIdentityProperties_GenerateCompositeIdentitySchema_CustomNames(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "example_resource.go")

	content := `package example

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceExample() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateCompositeIdentitySchema(nil, "first_id", "second_id"),
		},
	}
}
`
	if err := os.WriteFile(testFile, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	inferred, err := InferIdentityProperties(testFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedProps := []string{"first_id", "second_id"}
	if !slices.Equal(inferred.Properties, expectedProps) {
		t.Fatalf("expected properties %v, got %v", expectedProps, inferred.Properties)
	}
	if inferred.HasSubscriptionID {
		t.Fatalf("expected HasSubscriptionID to be false")
	}
}

func TestInferIdentityProperties_ScopeId(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "example_scope_resource.go")

	content := `package example

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceScopeExample() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.ScopeId{}),
		},
	}
}
`
	if err := os.WriteFile(testFile, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	inferred, err := InferIdentityProperties(testFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedProps := []string{"scope"}
	if !slices.Equal(inferred.Properties, expectedProps) {
		t.Fatalf("expected properties %v, got %v", expectedProps, inferred.Properties)
	}
	if inferred.HasSubscriptionID {
		t.Fatalf("expected HasSubscriptionID to be false for ScopeId")
	}
}

func TestInferIdentityProperties_ScopedBudget(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "example_scoped_budget_resource.go")

	content := `package example

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/consumption/2019-10-01/budgets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceScopedBudgetExample() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&budgets.ScopedBudgetId{}),
		},
	}
}
`
	if err := os.WriteFile(testFile, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	inferred, err := InferIdentityProperties(testFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedProps := []string{"scope", "name"}
	if !slices.Equal(inferred.Properties, expectedProps) {
		t.Fatalf("expected properties %v, got %v", expectedProps, inferred.Properties)
	}
	if inferred.HasSubscriptionID {
		t.Fatalf("expected HasSubscriptionID to be false for ScopedBudgetId")
	}
}

func TestResourceIdentityData_parseArgs_NoSubscriptionID(t *testing.T) {
	// Test 1: From GOFILE pointing to a composite ID resource
	filePath := filepath.Join(findProviderRoot(), "internal", "services", "network", "nat_gateway_public_ip_association_resource.go")
	t.Setenv("GOFILE", filePath)

	data := &resourceIdentityData{}
	args := []string{"-properties", "resource_id1:nat_gateway_id,resource_id2:public_ip_address_id"}
	errs := data.parseArgs(args)
	if len(errs) > 0 {
		t.Fatalf("unexpected parseArgs errors: %v", errs)
	}

	if !data.NoSubscriptionID {
		t.Fatalf("expected NoSubscriptionID to be true for composite resource without -no-subscription-id flag")
	}

	// Test 2: From -id without subscription_id
	t.Setenv("GOFILE", "")
	data2 := &resourceIdentityData{}
	args2 := []string{"-resource-name", "custom_res", "-id", "/providers/Microsoft.Foo/bars/{name}"}
	errs2 := data2.parseArgs(args2)
	if len(errs2) > 0 {
		t.Fatalf("unexpected parseArgs errors: %v", errs2)
	}

	if !data2.NoSubscriptionID {
		t.Fatalf("expected NoSubscriptionID to be true when -id does not contain subscription_id")
	}
}
