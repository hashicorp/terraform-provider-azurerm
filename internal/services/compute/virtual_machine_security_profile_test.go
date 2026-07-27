// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestFlattenVirtualMachineSecurityProfileConfiguration(t *testing.T) {
	profile := &virtualmachines.SecurityProfile{
		EncryptionAtHost: pointer.To(false),
	}

	for _, testCase := range securityProfileConfigurationTestCases(t) {
		t.Run(testCase.name, func(t *testing.T) {
			got := flattenVirtualMachineSecurityProfile(profile, testCase.data)
			if len(got) != testCase.wantLen {
				t.Fatalf("expected %d security profiles, got %d", testCase.wantLen, len(got))
			}
			if len(got) != 0 && got[0].(map[string]interface{})["host_encryption_enabled"] != false {
				t.Fatal("expected disabled host encryption to remain in the security profile")
			}
		})
	}
}

func TestFlattenVirtualMachineScaleSetSecurityProfileConfiguration(t *testing.T) {
	profile := &virtualmachinescalesets.SecurityProfile{
		EncryptionAtHost: pointer.To(false),
	}

	for _, testCase := range securityProfileConfigurationTestCases(t) {
		t.Run(testCase.name, func(t *testing.T) {
			got := flattenVirtualMachineScaleSetSecurityProfile(profile, testCase.data)
			if len(got) != testCase.wantLen {
				t.Fatalf("expected %d security profiles, got %d", testCase.wantLen, len(got))
			}
			if len(got) != 0 && got[0].(map[string]interface{})["host_encryption_enabled"] != false {
				t.Fatal("expected disabled host encryption to remain in the security profile")
			}
		})
	}
}

func securityProfileConfigurationTestCases(t *testing.T) []struct {
	name    string
	data    *schema.ResourceData
	wantLen int
} {
	t.Helper()

	profileSchema := map[string]*schema.Schema{
		"security_profile": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"host_encryption_enabled": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},
	}
	profileType := cty.Object(map[string]cty.Type{
		"host_encryption_enabled": cty.Bool,
	})
	configuredData := func(profile cty.Value) *schema.ResourceData {
		return (&schema.Resource{Schema: profileSchema}).Data(&terraform.InstanceState{
			RawConfig: cty.ObjectVal(map[string]cty.Value{
				"security_profile": profile,
			}),
		})
	}

	return []struct {
		name    string
		data    *schema.ResourceData
		wantLen int
	}{
		{
			name: "omitted block",
			data: configuredData(cty.NullVal(cty.List(profileType))),
		},
		{
			name: "explicitly disabled host encryption",
			data: configuredData(cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"host_encryption_enabled": cty.False,
				}),
			})),
			wantLen: 1,
		},
		{
			name:    "import without configuration",
			data:    (&schema.Resource{Schema: profileSchema}).Data(nil),
			wantLen: 1,
		},
	}
}
