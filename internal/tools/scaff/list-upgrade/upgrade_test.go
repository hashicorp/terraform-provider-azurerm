package list_upgrade

import (
	"go/format"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpgrade_TypedPlain_AddsIdentityAndFlatten(t *testing.T) {
	r, err := Analyze(filepath.Join("testdata", "typed_plain.go"))
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}

	out, changed, err := r.Upgrade(UpgradeOptions{
		AddIdentity:    true,
		ExtractFlatten: true,
		ReadModel:      "AzureMonitorWorkspaceResource",
		ResourceName:   "monitor_workspace",
	})
	if err != nil {
		t.Fatalf("Upgrade: %v", err)
	}
	if !changed {
		t.Fatalf("expected changed to be true")
	}

	// The result must be valid, gofmt-able Go.
	formatted, err := format.Source(out)
	if err != nil {
		t.Fatalf("result is not valid Go: %v\n\n%s", err, out)
	}
	got := string(formatted)

	mustContain(t, got, "var (")
	mustContain(t, got, "_ sdk.ResourceWithIdentity = WorkspaceResource{}")
	mustContain(t, got, "_ sdk.ResourceWithUpdate   = WorkspaceResource{}")
	mustContain(t, got, "func (r WorkspaceResource) Identity() resourceids.ResourceId {")
	mustContain(t, got, "return &azuremonitorworkspaces.AccountId{}")
	mustContain(t, got, "func (r WorkspaceResource) flatten(metadata sdk.ResourceMetaData, id *azuremonitorworkspaces.AccountId, model *azuremonitorworkspaces.AzureMonitorWorkspaceResource) error {")
	mustContain(t, got, "return r.flatten(metadata, id, resp.Model)")
	mustContain(t, got, "//go:generate go run ../../tools/generator-tests resourceidentity -resource-name monitor_workspace -properties \"name,resource_group_name\"")

	// The flatten body must set identity before encoding, and the model guard
	// must be collapsed to use the parameter.
	mustContain(t, got, "pluginsdk.SetResourceIdentityData(metadata.ResourceData, id)")
	mustContain(t, got, "if model != nil {")
	if strings.Contains(got, "if model := resp.Model; model != nil") {
		t.Errorf("expected the resp.Model guard to be collapsed in flatten")
	}

	// Create must set identity after SetID.
	mustContain(t, got, "return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)")
}

func TestUpgrade_TypedUpgraded_NoChanges(t *testing.T) {
	r, err := Analyze(filepath.Join("testdata", "typed_upgraded.go"))
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	_, changed, err := r.Upgrade(UpgradeOptions{
		AddIdentity:    true,
		ExtractFlatten: true,
		ReadModel:      "StorageMover",
		ResourceName:   "storage_mover",
	})
	if err != nil {
		t.Fatalf("Upgrade: %v", err)
	}
	if changed {
		t.Errorf("expected no changes for an already-upgraded resource")
	}
}

func TestUpgrade_Untyped_Unsupported(t *testing.T) {
	r, err := Analyze(filepath.Join("testdata", "untyped_plain.go"))
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if _, _, err := r.Upgrade(UpgradeOptions{AddIdentity: true}); err == nil {
		t.Errorf("expected an error upgrading an untyped resource")
	}
}

func mustContain(t *testing.T, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Errorf("expected output to contain:\n%s", needle)
	}
}
