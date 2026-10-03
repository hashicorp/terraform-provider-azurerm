// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package pluginsdk

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestSnakeCase(t *testing.T) {
	cases := []struct {
		Input  string
		Output string
	}{
		{
			Input:  "ResourceGroupName",
			Output: "resource_group_name",
		},
		{
			Input:  "ServerGroupsv2Name",
			Output: "server_groups_v2_name",
		},
		{
			Input:  "ServerGroupsV2Name",
			Output: "server_groups_v2_name",
		},
		{
			Input:  "ServerGroupsV42Name",
			Output: "server_groups_v42_name",
		},
		{
			Input:  "ServerGroupsV2",
			Output: "server_groups_v2",
		},
		{
			Input:  "Terraform404",
			Output: "terraform_404",
		},
		{
			Input:  "1TerraformProvider",
			Output: "1_terraform_provider",
		},
		{
			Input:  "TFProvider",
			Output: "tf_provider",
		},
		{
			Input:  "Peer2Peer",
			Output: "peer_2_peer",
		},
		{
			Input:  "V2_Terraform",
			Output: "v2_terraform",
		},
		{
			Input:  "vV2_Terraform",
			Output: "v_v2_terraform",
		},
		{
			Input:  "VV2_Terraform",
			Output: "v_v2_terraform",
		},
	}

	failures := make([]string, 0)
	for _, tc := range cases {
		if v := ToSnakeCase(tc.Input); v != tc.Output {
			failures = append(failures, fmt.Sprintf("expected %s, got %s", tc.Output, v))
		}
	}
	if len(failures) > 0 {
		t.Fatal(strings.Join(failures, "\n"))
	}
}

func TestSegmentTypeSupported_Composite(t *testing.T) {
	if !SegmentTypeSupported(resourceids.ResourceIDSegmentType) {
		t.Fatalf("expected ResourceIDSegmentType to be supported")
	}
}

func TestGenerateIdentitySchema_Composite(t *testing.T) {
	rgId := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")
	subId := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")
	compositeId := commonids.NewCompositeResourceID(&rgId, &subId)

	schemaFn := GenerateIdentitySchema(compositeId)
	idSchema := schemaFn()

	if len(idSchema) != 2 {
		t.Fatalf("expected 2 schema items, got %d", len(idSchema))
	}

	for _, k := range []string{"resource_id1", "resource_id2"} {
		s, ok := idSchema[k]
		if !ok {
			t.Fatalf("expected schema to have %q", k)
		}
		if !s.RequiredForImport {
			t.Fatalf("expected %q to have RequiredForImport=true", k)
		}
		if s.Type != schema.TypeString {
			t.Fatalf("expected %q to be TypeString, got %v", k, s.Type)
		}
	}
}

func TestGenerateCompositeIdentitySchema(t *testing.T) {
	rgId := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")
	subId := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")
	compositeId := commonids.NewCompositeResourceID(&rgId, &subId)

	schemaFn := GenerateCompositeIdentitySchema(compositeId, "resource_group_id", "subscription_id")
	idSchema := schemaFn()

	if len(idSchema) != 2 {
		t.Fatalf("expected 2 schema items, got %d", len(idSchema))
	}

	for _, k := range []string{"resource_group_id", "subscription_id"} {
		s, ok := idSchema[k]
		if !ok {
			t.Fatalf("expected schema to have %q", k)
		}
		if !s.RequiredForImport {
			t.Fatalf("expected %q to have RequiredForImport=true", k)
		}
	}
}

func TestValidateResourceId_Composite(t *testing.T) {
	rgId := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")
	subId := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")
	compositeId := commonids.NewCompositeResourceID(&rgId, &subId)

	validId := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group|/subscriptions/12345678-1234-9876-4563-123456789012"
	if err := ValidateResourceId(compositeId, validId, false); err != nil {
		t.Fatalf("expected valid ID to pass, got: %s", err)
	}

	invalidMissingPipe := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group"
	if err := ValidateResourceId(compositeId, invalidMissingPipe, false); err == nil {
		t.Fatalf("expected missing pipe ID to fail validation")
	}

	invalidSubId := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group|not-a-sub-id"
	if err := ValidateResourceId(compositeId, invalidSubId, false); err == nil {
		t.Fatalf("expected invalid component ID to fail validation")
	}
}

func TestSetResourceIdentityData_Composite(t *testing.T) {
	rgId := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")
	subId := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")
	compositeId := commonids.NewCompositeResourceID(&rgId, &subId)

	r := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"dummy": {Type: schema.TypeString, Optional: true},
		},
		Identity: &schema.ResourceIdentity{
			SchemaFunc: GenerateIdentitySchema(compositeId),
		},
	}

	d := r.TestResourceData()
	if err := SetResourceIdentityData(d, compositeId); err != nil {
		t.Fatalf("expected SetResourceIdentityData to succeed, got: %s", err)
	}

	identity, err := d.Identity()
	if err != nil {
		t.Fatalf("getting identity: %s", err)
	}

	if v := identity.Get("resource_id1"); v != rgId.ID() {
		t.Fatalf("expected resource_id1 to be %q, got %q", rgId.ID(), v)
	}
	if v := identity.Get("resource_id2"); v != subId.ID() {
		t.Fatalf("expected resource_id2 to be %q, got %q", subId.ID(), v)
	}
}

func TestSetCompositeResourceIdentityData(t *testing.T) {
	rgId := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")
	subId := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")
	compositeId := commonids.NewCompositeResourceID(&rgId, &subId)

	r := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"dummy": {Type: schema.TypeString, Optional: true},
		},
		Identity: &schema.ResourceIdentity{
			SchemaFunc: GenerateCompositeIdentitySchema(compositeId, "rg_id", "sub_id"),
		},
	}

	d := r.TestResourceData()
	if err := SetCompositeResourceIdentityData(d, compositeId, "rg_id", "sub_id"); err != nil {
		t.Fatalf("expected SetCompositeResourceIdentityData to succeed, got: %s", err)
	}

	identity, err := d.Identity()
	if err != nil {
		t.Fatalf("getting identity: %s", err)
	}

	if v := identity.Get("rg_id"); v != rgId.ID() {
		t.Fatalf("expected rg_id to be %q, got %q", rgId.ID(), v)
	}
	if v := identity.Get("sub_id"); v != subId.ID() {
		t.Fatalf("expected sub_id to be %q, got %q", subId.ID(), v)
	}
}

func TestValidateResourceIdentityData_Composite(t *testing.T) {
	rgId := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")
	subId := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")
	compositeId := commonids.NewCompositeResourceID(&rgId, &subId)

	idSchema := GenerateIdentitySchema(compositeId)()
	rawIdentity := map[string]string{
		"resource_id1": rgId.ID(),
		"resource_id2": subId.ID(),
	}

	d := schema.TestResourceDataWithIdentityRaw(t, map[string]*schema.Schema{}, idSchema, rawIdentity)
	if err := ValidateResourceIdentityData(d, compositeId); err != nil {
		t.Fatalf("expected ValidateResourceIdentityData to succeed, got: %s", err)
	}

	expectedId := compositeId.ID()
	if d.Id() != expectedId {
		t.Fatalf("expected d.Id() to be %q, got %q", expectedId, d.Id())
	}
}

func TestValidateCompositeResourceIdentityData(t *testing.T) {
	rgId := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")
	subId := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")
	compositeId := commonids.NewCompositeResourceID(&rgId, &subId)

	idSchema := GenerateCompositeIdentitySchema(compositeId, "rg_id", "sub_id")()
	rawIdentity := map[string]string{
		"rg_id":  rgId.ID(),
		"sub_id": subId.ID(),
	}

	d := schema.TestResourceDataWithIdentityRaw(t, map[string]*schema.Schema{}, idSchema, rawIdentity)
	if err := ValidateCompositeResourceIdentityData(d, compositeId, "rg_id", "sub_id"); err != nil {
		t.Fatalf("expected ValidateCompositeResourceIdentityData to succeed, got: %s", err)
	}

	expectedId := compositeId.ID()
	if d.Id() != expectedId {
		t.Fatalf("expected d.Id() to be %q, got %q", expectedId, d.Id())
	}
}

func TestImporterValidatingIdentity_Composite(t *testing.T) {
	rgId := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")
	subId := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")
	compositeId := commonids.NewCompositeResourceID(&rgId, &subId)

	importer := ImporterValidatingIdentity(compositeId)
	if importer == nil || importer.StateContext == nil {
		t.Fatalf("expected non-nil importer with StateContext")
	}
}

func TestImporterValidatingCompositeIdentity(t *testing.T) {
	rgId := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")
	subId := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")
	compositeId := commonids.NewCompositeResourceID(&rgId, &subId)

	importer := ImporterValidatingCompositeIdentity(compositeId, "rg_id", "sub_id")
	if importer == nil || importer.StateContext == nil {
		t.Fatalf("expected non-nil importer with StateContext")
	}
}

func TestSegmentTypeSupported_Scope(t *testing.T) {
	if !SegmentTypeSupported(resourceids.ScopeSegmentType) {
		t.Fatalf("expected ScopeSegmentType to be supported")
	}
}

func TestGenerateIdentitySchema_Scoped(t *testing.T) {
	scopedId := commonids.NewChaosStudioTargetID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "target1")
	schemaFn := GenerateIdentitySchema(&scopedId)
	idSchema := schemaFn()

	if len(idSchema) != 2 {
		t.Fatalf("expected 2 schema items (scope, name), got %d: %+v", len(idSchema), idSchema)
	}

	for _, k := range []string{"scope", "name"} {
		s, ok := idSchema[k]
		if !ok {
			t.Fatalf("expected schema to have %q", k)
		}
		if !s.RequiredForImport {
			t.Fatalf("expected %q to have RequiredForImport=true", k)
		}
		if s.Type != schema.TypeString {
			t.Fatalf("expected %q to be TypeString, got %v", k, s.Type)
		}
	}
}

func TestSetResourceIdentityData_Scoped(t *testing.T) {
	scopedId := commonids.NewChaosStudioTargetID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "target1")

	r := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"dummy": {Type: schema.TypeString, Optional: true},
		},
		Identity: &schema.ResourceIdentity{
			SchemaFunc: GenerateIdentitySchema(&scopedId),
		},
	}

	d := r.TestResourceData()
	if err := SetResourceIdentityData(d, &scopedId); err != nil {
		t.Fatalf("expected SetResourceIdentityData to succeed, got: %s", err)
	}

	identity, err := d.Identity()
	if err != nil {
		t.Fatalf("getting identity: %s", err)
	}

	if v := identity.Get("scope"); v != scopedId.Scope {
		t.Fatalf("expected scope to be %q, got %q", scopedId.Scope, v)
	}
	if v := identity.Get("name"); v != scopedId.TargetName {
		t.Fatalf("expected name to be %q, got %q", scopedId.TargetName, v)
	}
}

func TestValidateResourceIdentityData_Scoped(t *testing.T) {
	scopedId := commonids.NewChaosStudioTargetID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "target1")

	idSchema := GenerateIdentitySchema(&scopedId)()
	rawIdentity := map[string]string{
		"scope": scopedId.Scope,
		"name":  scopedId.TargetName,
	}

	d := schema.TestResourceDataWithIdentityRaw(t, map[string]*schema.Schema{}, idSchema, rawIdentity)
	if err := ValidateResourceIdentityData(d, &scopedId); err != nil {
		t.Fatalf("expected ValidateResourceIdentityData to succeed, got: %s", err)
	}

	expectedId := scopedId.ID()
	if d.Id() != expectedId {
		t.Fatalf("expected d.Id() to be %q, got %q", expectedId, d.Id())
	}
}
