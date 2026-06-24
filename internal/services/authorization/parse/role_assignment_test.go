// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"testing"
)

func TestRoleAssignmentIDFormatter(t *testing.T) {
	testData := []struct {
		SubscriptionId           string
		ResourceGroup            string
		ResourceProvider         string
		ResourceScope            string
		ManagementGroup          string
		SubscriptionAlias        string
		IsSubscriptionLevel      bool
		IsSubscriptionAliasLevel bool
		Name                     string
		TenantId                 string
		Expected                 string
	}{
		{
			SubscriptionId:  "",
			ResourceGroup:   "",
			ResourceScope:   "",
			ManagementGroup: "",
			Name:            "23456781-2349-8764-5631-234567890121",
			TenantId:        "",
		},
		{
			SubscriptionId:  "12345678-1234-9876-4563-123456789012",
			ResourceGroup:   "group1",
			ResourceScope:   "",
			ManagementGroup: "managementGroup1",
			Name:            "23456781-2349-8764-5631-234567890121",
			TenantId:        "",
		},
		{
			SubscriptionId:  "12345678-1234-9876-4563-123456789012",
			ResourceGroup:   "",
			ResourceScope:   "",
			ManagementGroup: "managementGroup1",
			Name:            "23456781-2349-8764-5631-234567890121",
			TenantId:        "",
		},
		{
			SubscriptionId:  "12345678-1234-9876-4563-123456789012",
			ResourceGroup:   "",
			ResourceScope:   "",
			ManagementGroup: "",
			Name:            "23456781-2349-8764-5631-234567890121",
			TenantId:        "",
			Expected:        "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
		},
		{
			SubscriptionId:  "12345678-1234-9876-4563-123456789012",
			ResourceGroup:   "group1",
			ResourceScope:   "",
			ManagementGroup: "",
			Name:            "23456781-2349-8764-5631-234567890121",
			TenantId:        "",
			Expected:        "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
		},
		{
			SubscriptionId:  "",
			ResourceGroup:   "",
			ResourceScope:   "",
			ManagementGroup: "12345678-1234-9876-4563-123456789012",
			Name:            "23456781-2349-8764-5631-234567890121",
			TenantId:        "",
			Expected:        "/providers/Microsoft.Management/managementGroups/12345678-1234-9876-4563-123456789012/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
		},
		{
			SubscriptionId:  "",
			ResourceGroup:   "",
			ResourceScope:   "",
			ManagementGroup: "12345678-1234-9876-4563-123456789012",
			Name:            "23456781-2349-8764-5631-234567890121",
			TenantId:        "34567812-3456-7653-6742-345678901234",
			Expected:        "/providers/Microsoft.Management/managementGroups/12345678-1234-9876-4563-123456789012/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121|34567812-3456-7653-6742-345678901234",
		},
		{
			SubscriptionId:   "12345678-1234-9876-4563-123456789012",
			ResourceGroup:    "group1",
			ResourceProvider: "Microsoft.Storage",
			ResourceScope:    "storageAccounts/nameStorageAccount",
			ManagementGroup:  "",
			Name:             "23456781-2349-8764-5631-234567890121",
			TenantId:         "34567812-3456-7653-6742-345678901234",
			Expected:         "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/nameStorageAccount/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121|34567812-3456-7653-6742-345678901234",
		},
		{
			IsSubscriptionLevel: true,
			Name:                "23456781-2349-8764-5631-234567890121",
			TenantId:            "34567812-3456-7653-6742-345678901234",
			Expected:            "/providers/Microsoft.Subscription/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121|34567812-3456-7653-6742-345678901234",
		},
		{
			IsSubscriptionAliasLevel: true,
			Name:                     "23456781-2349-8764-5631-234567890121",
			TenantId:                 "34567812-3456-7653-6742-345678901234",
			SubscriptionAlias:        "my-awesome-sub",
			Expected:                 "/providers/Microsoft.Subscription/aliases/my-awesome-sub/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121|34567812-3456-7653-6742-345678901234",
		},
	}
	for _, v := range testData {
		t.Logf("testing %+v", v)
		actual, err := NewRoleAssignmentID(v.SubscriptionId, v.ResourceGroup, v.ResourceProvider, v.ResourceScope, v.ManagementGroup, v.Name, v.TenantId, v.SubscriptionAlias, v.IsSubscriptionLevel, v.IsSubscriptionAliasLevel)
		if err != nil {
			if v.Expected == "" {
				continue
			}
			t.Fatal(err)
		}
		actualId := actual.ID()
		if actualId != v.Expected {
			t.Fatalf("expected %q, got %q", v.Expected, actualId)
		}
	}
}

func TestRoleAssignmentID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *RoleAssignmentId
	}{
		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Error: true,
		},

		{
			// just subscription
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},

		{
			// missing ResourceGroup value
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},

		{
			// missing Management Group value
			Input: "/providers/Microsoft.Management/managementGroups/",
			Error: true,
		},

		{
			// missing Role Assignment value at Subscription Scope
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Authorization/roleAssignments/",
			Error: true,
		},

		{
			// missing Role Assignment value at Management Group scope
			Input: "/providers/Microsoft.Management/managementGroups/managementGroup1/providers/Microsoft.Authorization/roleAssignments/",
			Error: true,
		},

		{
			// valid at subscriptions scope
			Input: "/providers/Microsoft.Subscription/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
			Expected: &RoleAssignmentId{
				IsSubscriptionLevel: true,
				Name:                "23456781-2349-8764-5631-234567890121",
			},
		},

		{
			// valid at subscriptions aliases scope
			Input: "/providers/Microsoft.Subscription/aliases/my-awesome-sub/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
			Expected: &RoleAssignmentId{
				IsSubscriptionAliasLevel: true,
				Name:                     "23456781-2349-8764-5631-234567890121",
				SubscriptionAlias:        "my-awesome-sub",
			},
		},

		{
			// valid at subscription scope
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
			Expected: &RoleAssignmentId{
				SubscriptionID:  "12345678-1234-9876-4563-123456789012",
				ResourceGroup:   "",
				ResourceScope:   "",
				ManagementGroup: "",
				Name:            "23456781-2349-8764-5631-234567890121",
			},
		},

		{
			// valid at resource group scope
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
			Expected: &RoleAssignmentId{
				SubscriptionID:  "12345678-1234-9876-4563-123456789012",
				ResourceGroup:   "group1",
				ManagementGroup: "",
				Name:            "23456781-2349-8764-5631-234567890121",
			},
		},

		{
			// valid at management group scope
			Input: "/providers/Microsoft.Management/managementGroups/managementGroup1/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
			Expected: &RoleAssignmentId{
				SubscriptionID:  "",
				ResourceGroup:   "",
				ManagementGroup: "managementGroup1",
				Name:            "23456781-2349-8764-5631-234567890121",
			},
		},
		{
			Input: "/providers/Microsoft.Management/managementGroups/managementGroup1/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121|34567812-3456-7653-6742-345678901234",
			Expected: &RoleAssignmentId{
				SubscriptionID:  "",
				ResourceGroup:   "",
				ManagementGroup: "managementGroup1",
				Name:            "23456781-2349-8764-5631-234567890121",
				TenantId:        "34567812-3456-7653-6742-345678901234",
			},
		},
		{
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/nameStorageAccount/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121|34567812-3456-7653-6742-345678901234",
			Expected: &RoleAssignmentId{
				SubscriptionID:   "12345678-1234-9876-4563-123456789012",
				ResourceGroup:    "group1",
				ResourceProvider: "Microsoft.Storage",
				ResourceScope:    "storageAccounts/nameStorageAccount",
				ManagementGroup:  "",
				Name:             "23456781-2349-8764-5631-234567890121",
				TenantId:         "34567812-3456-7653-6742-345678901234",
			},
		},
		{
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.AppPlatform/Spring/spring1/apps/app1/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121|34567812-3456-7653-6742-345678901234",
			Expected: &RoleAssignmentId{
				SubscriptionID:   "12345678-1234-9876-4563-123456789012",
				ResourceGroup:    "group1",
				ResourceProvider: "Microsoft.AppPlatform",
				ResourceScope:    "Spring/spring1/apps/app1",
				ManagementGroup:  "",
				Name:             "23456781-2349-8764-5631-234567890121",
				TenantId:         "34567812-3456-7653-6742-345678901234",
			},
		},
		{
			Input: "/providers/Microsoft.Capacity/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
			Expected: &RoleAssignmentId{
				SubscriptionID:   "",
				ResourceGroup:    "",
				ResourceProvider: "Microsoft.Capacity",
				ResourceScope:    "",
				ManagementGroup:  "",
				Name:             "23456781-2349-8764-5631-234567890121",
				TenantId:         "34567812-3456-7653-6742-345678901234",
			},
		},
		{
			Input: "/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
			Expected: &RoleAssignmentId{
				SubscriptionID:   "",
				ResourceGroup:    "",
				ResourceProvider: "",
				ResourceScope:    "",
				ManagementGroup:  "",
				Name:             "23456781-2349-8764-5631-234567890121",
				TenantId:         "34567812-3456-7653-6742-345678901234",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := RoleAssignmentID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("expected a value but got an error: %+v", err)
		}

		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Role Assignment Name", v.Expected.Name, actual.Name)
		}

		if actual.SubscriptionID != v.Expected.SubscriptionID {
			t.Fatalf("Expected %q but got %q for Role Assignment Subscription ID", v.Expected.SubscriptionID, actual.SubscriptionID)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Role Assignment Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.ResourceProvider != v.Expected.ResourceProvider {
			t.Fatalf("Expected %q but got %q for Role Assignment Resource Provider", v.Expected.ResourceProvider, actual.ResourceProvider)
		}

		if actual.ResourceScope != v.Expected.ResourceScope {
			t.Fatalf("Expected %q but got %q for Role Assignment Resource Scope", v.Expected.ResourceScope, actual.ResourceScope)
		}

		if actual.ManagementGroup != v.Expected.ManagementGroup {
			t.Fatalf("Expected %q but got %q for Role Assignment Management Group", v.Expected.ManagementGroup, actual.ManagementGroup)
		}

		if actual.IsSubscriptionLevel != v.Expected.IsSubscriptionLevel {
			t.Fatalf("Expected %v but got %v for Role Assignment SubscriptionLevel flag", v.Expected.IsSubscriptionLevel, actual.IsSubscriptionLevel)
		}

		if actual.IsSubscriptionAliasLevel != v.Expected.IsSubscriptionAliasLevel {
			t.Fatalf("Expected %v but got %v for Role Assignment SubscriptionAliasLevel flag", v.Expected.IsSubscriptionAliasLevel, actual.IsSubscriptionAliasLevel)
		}

		if actual.SubscriptionAlias != v.Expected.SubscriptionAlias {
			t.Fatalf("Expected %q but got %q for Role Assignment SubscriptionAlias", v.Expected.SubscriptionAlias, actual.SubscriptionAlias)
		}
	}
}

func TestRoleAssignmentName(t *testing.T) {
	testData := []struct {
		Scope            string
		PrincipalId      string
		RoleDefinitionId string
		Expected         string
		Error            bool
	}{
		{
			Scope:            "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/myRg",
			PrincipalId:      "22222222-2222-2222-2222-222222222222",
			RoleDefinitionId: "/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Authorization/roleDefinitions/33333333-3333-3333-3333-333333333333",
			Expected:         "eec47c62-a9a9-54dd-b8d8-f9e2839fab35",
		},
		{
			Scope:            "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/myRg",
			PrincipalId:      "22222222-2222-2222-2222-222222222222",
			RoleDefinitionId: "/providers/Microsoft.Authorization/roleDefinitions/33333333-3333-3333-3333-333333333333",
			Expected:         "eec47c62-a9a9-54dd-b8d8-f9e2839fab35",
		},
		{
			Scope:            "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/myRg",
			PrincipalId:      "22222222-2222-2222-2222-222222222222",
			RoleDefinitionId: "/providers/Microsoft.Management/managementGroups/myMg/providers/Microsoft.Authorization/roleDefinitions/33333333-3333-3333-3333-333333333333",
			Expected:         "eec47c62-a9a9-54dd-b8d8-f9e2839fab35",
		},
		{
			Scope:            "/SUBSCRIPTIONS/11111111-1111-1111-1111-111111111111/RESOURCEGROUPS/MYRG",
			PrincipalId:      "22222222-2222-2222-2222-222222222222",
			RoleDefinitionId: "/SUBSCRIPTIONS/11111111-1111-1111-1111-111111111111/PROVIDERS/MICROSOFT.AUTHORIZATION/ROLEDEFINITIONS/33333333-3333-3333-3333-333333333333",
			Expected:         "eec47c62-a9a9-54dd-b8d8-f9e2839fab35",
		},
		{
			Scope:            "/",
			PrincipalId:      "22222222-2222-2222-2222-222222222222",
			RoleDefinitionId: "/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Authorization/roleDefinitions/33333333-3333-3333-3333-333333333333",
			Expected:         "c095807d-5244-5404-8bce-e0e762e70618",
		},
		{
			Scope:            "/providers/Microsoft.Management/managementGroups/myMg",
			PrincipalId:      "22222222-2222-2222-2222-222222222222",
			RoleDefinitionId: "/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Authorization/roleDefinitions/33333333-3333-3333-3333-333333333333",
			Expected:         "6aab5ec9-d102-539e-8b65-b053c060e64f",
		},
		{
			Scope:            "/subscriptions/11111111-1111-1111-1111-111111111111",
			PrincipalId:      "99999999-9999-9999-9999-999999999999",
			RoleDefinitionId: "/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Authorization/roleDefinitions/33333333-3333-3333-3333-333333333333",
			Expected:         "54e40269-cc30-57bd-8031-5ca4efc75fd5",
		},
		{
			Scope:            "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/myRg",
			PrincipalId:      "22222222-2222-2222-2222-222222222222",
			RoleDefinitionId: "/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Authorization/roleDefinitions/not-a-uuid",
			Error:            true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing Scope=%q PrincipalId=%q RoleDefinitionId=%q", v.Scope, v.PrincipalId, v.RoleDefinitionId)

		actual, err := RoleAssignmentName(v.Scope, v.PrincipalId, v.RoleDefinitionId)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("expected a value but got an error: %+v", err)
		}

		if v.Error {
			t.Fatalf("expected an error but got none")
		}

		if actual != v.Expected {
			t.Fatalf("expected %q, got %q", v.Expected, actual)
		}
	}
}
