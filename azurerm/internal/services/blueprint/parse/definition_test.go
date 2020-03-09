package parse

import "testing"

func TestBlueprintDefinitionID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Error    bool
		Expected *BlueprintDefinitionId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "Blueprint ID in Management group",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprints/blueprint1",
			Expected: &BlueprintDefinitionId{
				ID:   "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprints/blueprint1",
				Name: "blueprint1",
				BlueprintDefinitionScopeId: BlueprintDefinitionScopeId{
					ScopeId:           "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
					Type:              AtManagementGroup,
					SubscriptionId:    "",
					ManagementGroupId: "00000000-0000-0000-0000-000000000000",
				},
			},
		},
		{
			Name:  "Blueprint ID in Management group but with wrong casing in management group part",
			Input: "/providers/microsoft.management/managementgroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprints/blueprint1",
			Expected: &BlueprintDefinitionId{
				ID:   "/providers/microsoft.management/managementgroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprints/blueprint1",
				Name: "blueprint1",
				BlueprintDefinitionScopeId: BlueprintDefinitionScopeId{
					ScopeId:           "/providers/microsoft.management/managementgroups/00000000-0000-0000-0000-000000000000",
					Type:              AtManagementGroup,
					SubscriptionId:    "",
					ManagementGroupId: "00000000-0000-0000-0000-000000000000",
				},
			},
		},
		{
			Name:  "Missing management group ID",
			Input: "/providers/Microsoft.Management/managementGroups/providers/Microsoft.Blueprint/blueprints/blueprint1",
			Error: true,
		},
		{
			Name:  "Missing blueprint name",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprints/",
			Error: true,
		},
		{
			Name:  "Blueprint ID in subscription",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprints/blueprint1",
			Expected: &BlueprintDefinitionId{
				ID:   "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprints/blueprint1",
				Name: "blueprint1",
				BlueprintDefinitionScopeId: BlueprintDefinitionScopeId{
					ScopeId:           "/subscriptions/00000000-0000-0000-0000-000000000000",
					Type:              AtSubscription,
					SubscriptionId:    "00000000-0000-0000-0000-000000000000",
					ManagementGroupId: "",
				},
			},
		},
		{
			Name:  "Missing blueprint name in subscription",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprints/",
			Error: true,
		},
		{
			Name:  "BlueprintID in resource group",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Blueprint/blueprints/blueprint1",
			Error: true,
		},
		{
			Name:  "Invalid resource ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/providers/Microsoft.Blueprint/blueprints/blueprint1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := BlueprintDefinitionID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %+v", err)
		}

		if actual.ID != v.Expected.ID {
			t.Fatalf("Expected %q but got %q", v.Expected.ID, actual.ID)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q", v.Expected.Name, actual.Name)
		}

		if actual.Type != v.Expected.Type {
			t.Fatalf("Expected type %q but got type %q", v.Expected.Type, actual.Type)
		}

		if actual.ScopeId != v.Expected.ScopeId {
			t.Fatalf("Expected %q but got %q", v.Expected.ScopeId, actual.ScopeId)
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q", v.Expected.SubscriptionId, actual.SubscriptionId)
		}

		if actual.ManagementGroupId != v.Expected.ManagementGroupId {
			t.Fatalf("Expected %q but got %q", v.Expected.ManagementGroupId, actual.ManagementGroupId)
		}
	}
}

func TestBlueprintDefinitionScopeID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Error    bool
		Expected *BlueprintDefinitionScopeId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "Management group ID",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
			Expected: &BlueprintDefinitionScopeId{
				ScopeId:           "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
				Type:              AtManagementGroup,
				SubscriptionId:    "",
				ManagementGroupId: "00000000-0000-0000-0000-000000000000",
			},
		},
		{
			Name:  "Management group ID with wrong casing",
			Input: "/providers/microsoft.management/managementgroups/00000000-0000-0000-0000-000000000000",
			Expected: &BlueprintDefinitionScopeId{
				ScopeId:           "/providers/microsoft.management/managementgroups/00000000-0000-0000-0000-000000000000",
				Type:              AtManagementGroup,
				SubscriptionId:    "",
				ManagementGroupId: "00000000-0000-0000-0000-000000000000",
			},
		},
		{
			Name:  "Management group ID but missing components",
			Input: "/providers/Microsoft.Management/managementGroups/",
			Error: true,
		},
		{
			Name:  "Subscription ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expected: &BlueprintDefinitionScopeId{
				ScopeId:           "/subscriptions/00000000-0000-0000-0000-000000000000",
				Type:              AtSubscription,
				SubscriptionId:    "00000000-0000-0000-0000-000000000000",
				ManagementGroupId: "",
			},
		},
		{
			Name:  "Resource Group ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1",
			Error: true,
		},
		{
			Name:  "Incomplete resource group ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Error: true,
		},
		{
			Name:  "Resource ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/vm1",
			Error: true,
		},
		{
			Name:  "Incomplete resource ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := BlueprintDefinitionScopeID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %+v", err)
		}

		if actual.Type != v.Expected.Type {
			t.Fatalf("Expected type %q but got type %q", v.Expected.Type, actual.Type)
		}

		if actual.ScopeId != v.Expected.ScopeId {
			t.Fatalf("Expected %q but got %q", v.Expected.ScopeId, actual.ScopeId)
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q", v.Expected.SubscriptionId, actual.SubscriptionId)
		}

		if actual.ManagementGroupId != v.Expected.ManagementGroupId {
			t.Fatalf("Expected %q but got %q", v.Expected.ManagementGroupId, actual.ManagementGroupId)
		}
	}
}

func TestPublishedBlueprintID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Error    bool
		Expected *PublishedBlueprintId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "Published Blueprint ID in Management group",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprints/blueprint1/versions/v1",
			Expected: &PublishedBlueprintId{
				Version: "v1",
				BlueprintDefinitionId: BlueprintDefinitionId{
					ID:   "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprints/blueprint1",
					Name: "blueprint1",
					BlueprintDefinitionScopeId: BlueprintDefinitionScopeId{
						ScopeId:           "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
						Type:              AtManagementGroup,
						SubscriptionId:    "",
						ManagementGroupId: "00000000-0000-0000-0000-000000000000",
					},
				},
			},
		},
		{
			Name:  "Published Blueprint ID in Management group but with wrong casing in management group part",
			Input: "/providers/microsoft.management/managementgroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprints/blueprint1/versions/v1",
			Expected: &PublishedBlueprintId{
				Version: "v1",
				BlueprintDefinitionId: BlueprintDefinitionId{
					ID:   "/providers/microsoft.management/managementgroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprints/blueprint1",
					Name: "blueprint1",
					BlueprintDefinitionScopeId: BlueprintDefinitionScopeId{
						ScopeId:           "/providers/microsoft.management/managementgroups/00000000-0000-0000-0000-000000000000",
						Type:              AtManagementGroup,
						SubscriptionId:    "",
						ManagementGroupId: "00000000-0000-0000-0000-000000000000",
					},
				},
			},
		},
		{
			Name:  "Missing version",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprints/blueprint1/versions/",
			Error: true,
		},
		{
			Name:  "Unpublished Blueprint ID",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprints/blueprint1",
			Error: true,
		},
		{
			Name:  "Missing management group ID",
			Input: "/providers/Microsoft.Management/managementGroups/providers/Microsoft.Blueprint/blueprints/blueprint1",
			Error: true,
		},
		{
			Name:  "Missing blueprint name",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprints/",
			Error: true,
		},
		{
			Name:  "Blueprint ID in subscription",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprints/blueprint1/versions/v1",
			Expected: &PublishedBlueprintId{
				Version: "v1",
				BlueprintDefinitionId: BlueprintDefinitionId{
					ID:   "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprints/blueprint1",
					Name: "blueprint1",
					BlueprintDefinitionScopeId: BlueprintDefinitionScopeId{
						ScopeId:           "/subscriptions/00000000-0000-0000-0000-000000000000",
						Type:              AtSubscription,
						SubscriptionId:    "00000000-0000-0000-0000-000000000000",
						ManagementGroupId: "",
					},
				},
			},
		},
		{
			Name:  "Missing version in subscription",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprints/blueprint1/versions/",
			Error: true,
		},
		{
			Name:  "Published BlueprintID in resource group",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Blueprint/blueprints/blueprint1/versions/v1",
			Error: true,
		},
		{
			Name:  "Invalid resource ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/providers/Microsoft.Blueprint/blueprints/blueprint1/versions/v1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := PublishedBlueprintID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %+v", err)
		}

		if actual.Version != v.Expected.Version {
			t.Fatalf("Expected %q but got %q", v.Expected.Version, actual.Version)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q", v.Expected.Name, actual.Name)
		}

		if actual.ID != v.Expected.ID {
			t.Fatalf("Expected %q but got %q", v.Expected.ID, actual.Name)
		}

		if actual.Type != v.Expected.Type {
			t.Fatalf("Expected type %q but got type %q", v.Expected.Type, actual.Type)
		}

		if actual.ScopeId != v.Expected.ScopeId {
			t.Fatalf("Expected %q but got %q", v.Expected.ScopeId, actual.ScopeId)
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q", v.Expected.SubscriptionId, actual.SubscriptionId)
		}

		if actual.ManagementGroupId != v.Expected.ManagementGroupId {
			t.Fatalf("Expected %q but got %q", v.Expected.ManagementGroupId, actual.ManagementGroupId)
		}
	}
}

func TestManagementGroupID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Error    bool
		Expected *ManagementGroupId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "Missing management group segment",
			Input: "/providers/Microsoft.Management/resourceGroups/group1",
			Error: true,
		},
		{
			Name:  "Missing right provider",
			Input: "/managementGroups/00000000-0000-0000-0000-000000000000",
			Error: true,
		},
		{
			Name:  "Subscription ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000",
			Error: true,
		},
		{
			Name:  "Resource Group ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1",
			Error: true,
		},
		{
			Name:  "Resource ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/vm1",
			Error: true,
		},
		{
			Name:  "Resource ID-like",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/virtualMachines/vm1",
			Error: true,
		},
		{
			Name:  "Management Group ID",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
			Expected: &ManagementGroupId{
				GroupId: "00000000-0000-0000-0000-000000000000",
			},
		},
		{
			Name:  "Management Group ID with wrong casing in provider",
			Input: "/providers/microsoft.management/managementGroups/00000000-0000-0000-0000-000000000000",
			Expected: &ManagementGroupId{
				GroupId: "00000000-0000-0000-0000-000000000000",
			},
		},
		{
			Name:  "Management Group ID with wrong casing in provider and managementGroup segment",
			Input: "/providers/microsoft.management/managementgroups/00000000-0000-0000-0000-000000000000",
			Expected: &ManagementGroupId{
				GroupId: "00000000-0000-0000-0000-000000000000",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ManagementGroupID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %+v", err)
		}

		if actual.GroupId != v.Expected.GroupId {
			t.Fatalf("Expected %q but got %q", v.Expected.GroupId, actual.GroupId)
		}
	}
}
