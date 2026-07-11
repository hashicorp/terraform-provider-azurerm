// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package list_upgrade

import (
	"path/filepath"
	"testing"
)

func TestAnalyze_TypedPlain(t *testing.T) {
	r, err := Analyze(filepath.Join("testdata", "typed_plain.go"))
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}

	assertEqual(t, "kind", r.Kind.String(), "typed")
	assertEqual(t, "package", r.Package, "testdata")
	assertEqual(t, "struct", r.StructName, "WorkspaceResource")
	assertEqual(t, "model", r.ModelStruct, "WorkspaceResourceModel")
	assertEqual(t, "terraform type", r.TerraformType, "azurerm_monitor_workspace")
	assertEqual(t, "sdk package", r.SDKPackage, "azuremonitorworkspaces")
	assertEqual(t, "sdk import", r.SDKImportPath, "github.com/hashicorp/go-azure-sdk/resource-manager/monitor/2023-04-03/azuremonitorworkspaces")
	assertEqual(t, "id type", r.IDTypeName, "AccountId")
	assertEqual(t, "id parse", r.IDParseFunc, "ParseAccountID")
	assertEqual(t, "id validate", r.IDValidateFunc, "ValidateAccountID")
	assertEqual(t, "service", r.ServiceName, "Monitor")
	assertEqual(t, "client field", r.ClientField, "WorkspacesClient")

	if !r.HasUpdate {
		t.Errorf("expected HasUpdate to be true")
	}
	if r.HasIdentity {
		t.Errorf("expected HasIdentity to be false")
	}
	if r.HasFlatten {
		t.Errorf("expected HasFlatten to be false")
	}
}

func TestAnalyze_TypedUpgraded(t *testing.T) {
	r, err := Analyze(filepath.Join("testdata", "typed_upgraded.go"))
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}

	assertEqual(t, "kind", r.Kind.String(), "typed")
	assertEqual(t, "struct", r.StructName, "StorageMoverResource")
	assertEqual(t, "id type", r.IDTypeName, "StorageMoverId")
	assertEqual(t, "sdk package", r.SDKPackage, "storagemovers")

	if !r.HasIdentity {
		t.Errorf("expected HasIdentity to be true")
	}
	if !r.HasFlatten {
		t.Errorf("expected HasFlatten to be true")
	}
	if !r.HasUpdate {
		t.Errorf("expected HasUpdate to be true")
	}
}

func TestAnalyze_Untyped(t *testing.T) {
	r, err := Analyze(filepath.Join("testdata", "untyped_plain.go"))
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	assertEqual(t, "kind", r.Kind.String(), "untyped")
}

func assertEqual(t *testing.T, label, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %q, want %q", label, got, want)
	}
}
