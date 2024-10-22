// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseRoleDefinitionID(t *testing.T) {
	testData := []struct {
		Input       string
		Expected    *RoleDefinitionID
		ShouldParse bool
	}{
		{
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/12345678-1234-1234-1234-1234567890ab|/subscriptions/00000000-0000-0000-0000-000000000000",
			Expected: &RoleDefinitionID{
				ResourceID: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/12345678-1234-1234-1234-1234567890ab",
				Scope:      "/subscriptions/00000000-0000-0000-0000-000000000000",
				RoleID:     "12345678-1234-1234-1234-1234567890ab",
			},
			ShouldParse: true,
		},
		{
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/12345678-1234-1234-1234-1234567890ab|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/valid.resource()-_group",
			Expected: &RoleDefinitionID{
				ResourceID: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/12345678-1234-1234-1234-1234567890ab",
				Scope:      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/valid.resource()-_group",
				RoleID:     "12345678-1234-1234-1234-1234567890ab",
			},
			ShouldParse: true,
		},
		{
			Input: "/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/providers/Microsoft.Authorization/roleDefinitions/12345678-1234-1234-1234-1234567890ab|/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup/providers/Microsoft.Compute/virtualMachines/myVM",
			Expected: &RoleDefinitionID{
				ResourceID: "/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/providers/Microsoft.Authorization/roleDefinitions/12345678-1234-1234-1234-1234567890ab",
				Scope:      "/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup/providers/Microsoft.Compute/virtualMachines/myVM",
				RoleID:     "12345678-1234-1234-1234-1234567890ab",
			},
			ShouldParse: true,
		},
		{
			Input: "/providers/Microsoft.Authorization/rOlEdEfiNiTiOns/AB65ee35-Dc39-43aC-86a0-aCCaaDf4abcd|/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup/providers/Microsoft.Compute/virtualMachines/myVM",
			Expected: &RoleDefinitionID{
				ResourceID: "/providers/Microsoft.Authorization/rOlEdEfiNiTiOns/AB65ee35-Dc39-43aC-86a0-aCCaaDf4abcd",
				Scope:      "/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup/providers/Microsoft.Compute/virtualMachines/myVM",
				RoleID:     "AB65ee35-Dc39-43aC-86a0-aCCaaDf4abcd",
			},
			ShouldParse: true,
		},
		{
			Input:       "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/00000000-0000-0000-0000-000000000000|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/invalid.resource()-_group.",
			Expected:    nil,
			ShouldParse: false,
		},
		{
			Input: "/providers/Microsoft.Authorization/roleDefinitions/ab65ee35-dc39-43ac-86a0-accaadf4abcd|/providers/Microsoft.Management/managementGroups/Some-Management-Group",
			Expected: &RoleDefinitionID{
				ResourceID: "/providers/Microsoft.Authorization/roleDefinitions/ab65ee35-dc39-43ac-86a0-accaadf4abcd",
				Scope:      "/providers/Microsoft.Management/managementGroups/Some-Management-Group",
				RoleID:     "ab65ee35-dc39-43ac-86a0-accaadf4abcd",
			},
			ShouldParse: true,
		},
		{
			Input: "/pRoVidErs/MiCrOsOft.AuThORiZaTiOn/RoLeDeFinItiOns/ab65ee35-DC39-43AC-86a0-accaadf4abcd|/PrOviDeRs/MiCrOsOfT.MaNaGeMeNt/maNAGementGroups/Some-Management-Group",
			Expected: &RoleDefinitionID{
				ResourceID: "/pRoVidErs/MiCrOsOft.AuThORiZaTiOn/RoLeDeFinItiOns/ab65ee35-DC39-43AC-86a0-accaadf4abcd",
				Scope:      "/PrOviDeRs/MiCrOsOfT.MaNaGeMeNt/maNAGementGroups/Some-Management-Group",
				RoleID:     "ab65ee35-DC39-43AC-86a0-accaadf4abcd",
			},
			ShouldParse: true,
		},
		{
			Input: "/providers/Microsoft.Authorization/roleDefinitions/ab65ee35-dc39-43ac-86a0-accaadf4abcd|/providers/Microsoft.Management/managementGroups/Some-Management-Group/subscriptions/12345678-1234-1234-1234-1234567890ab",
			Expected: &RoleDefinitionID{
				ResourceID: "/providers/Microsoft.Authorization/roleDefinitions/ab65ee35-dc39-43ac-86a0-accaadf4abcd",
				Scope:      "/providers/Microsoft.Management/managementGroups/Some-Management-Group/subscriptions/12345678-1234-1234-1234-1234567890ab",
				RoleID:     "ab65ee35-dc39-43ac-86a0-accaadf4abcd",
			},
			ShouldParse: true,
		},
		{
			Input: "/providers/Microsoft.Authorization/roleDefinitions/ab65ee35-dc39-43ac-86a0-accaadf4abcd|/providers/Microsoft.Management/managementGroups/Some-Management-Group/subscriptions/12345678-1234-1234-1234-1234567890ab/resourceGroups/myGroup/providers/Microsoft.Compute/virtualMachines/myVM",
			Expected: &RoleDefinitionID{
				ResourceID: "/providers/Microsoft.Authorization/roleDefinitions/ab65ee35-dc39-43ac-86a0-accaadf4abcd",
				Scope:      "/providers/Microsoft.Management/managementGroups/Some-Management-Group/subscriptions/12345678-1234-1234-1234-1234567890ab/resourceGroups/myGroup/providers/Microsoft.Compute/virtualMachines/myVM",
				RoleID:     "ab65ee35-dc39-43ac-86a0-accaadf4abcd",
			},
			ShouldParse: true,
		},
		{
			Input: "/providers/Microsoft.Authorization/roleDefinitions/ab65ee35-dc39-43ac-86a0-accaadf4abcd|/",
			Expected: &RoleDefinitionID{
				ResourceID: "/providers/Microsoft.Authorization/roleDefinitions/ab65ee35-dc39-43ac-86a0-accaadf4abcd",
				Scope:      "/",
				RoleID:     "ab65ee35-dc39-43ac-86a0-accaadf4abcd",
			},
			ShouldParse: true,
		},
		{
			Input: "/providers/Microsoft.Authorization/rOlEdEfiNiTiOns/AB65ee35-Dc39-43aC-86a0-aCCaaDf4abcd|/",
			Expected: &RoleDefinitionID{
				ResourceID: "/providers/Microsoft.Authorization/rOlEdEfiNiTiOns/AB65ee35-Dc39-43aC-86a0-aCCaaDf4abcd",
				Scope:      "/",
				RoleID:     "AB65ee35-Dc39-43aC-86a0-aCCaaDf4abcd",
			},
			ShouldParse: true,
		},
		{
			Input:       "12345678-1234-1234-1234-1234567890ab|/providers/Microsoft.Management/managementGroups/Some-Management-Group",
			Expected:    nil,
			ShouldParse: false,
		},
		{
			Input:       "12345678-1234-inva-lid1-|/providers/Microsoft.Management/managementGroups/Some-Management-Group",
			Expected:    nil,
			ShouldParse: false,
		},
		{
			Input:       "12345678-nonsense",
			Expected:    nil,
			ShouldParse: false,
		},
		{
			Input:       "/providers/Microsoft.Authorization/roleDefinitions/ab65ee35-INVALID-86a0-accaadf4abcd|/",
			Expected:    nil,
			ShouldParse: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] parsing %q, ShouldParse %+v", v.Input, v.ShouldParse)

		roleDefinitionId, err := RoleDefinitionId(v.Input)

		if v.ShouldParse {
			switch {
			case err != nil:
				t.Fatalf("expected %q parse success, got failure: %q", v.Input, err)
			case !cmp.Equal(v.Expected, roleDefinitionId):
				t.Fatalf("parse succeeded but expected %+v, got %+v", v.Expected, roleDefinitionId)
			default:
				t.Logf("[DEBUG] parse succeeded as expected, got:\n%+v", roleDefinitionId)
			}
		} else {
			if err == nil {
				t.Fatalf("expected %q parse failure, got success: %+v", v.Input, roleDefinitionId)
			}
			t.Logf("[DEBUG] parse failed as expected, got %q", err)
		}
	}
}
