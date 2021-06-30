package azure_test

import (
	"reflect"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func TestParseAzureResourceID(t *testing.T) {
	testCases := []struct {
		id                 string
		expectedResourceID *azure.ResourceID
		expectError        bool
	}{
		{
			// Missing "resourceGroups".
			"/subscriptions/00000000-0000-0000-0000-000000000000//myResourceGroup/",
			nil,
			true,
		},
		{
			// Empty resource group ID.
			"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups//",
			nil,
			true,
		},
		{
			"random",
			nil,
			true,
		},
		{
			"/subscriptions/6d74bdd2-9f84-11e5-9bd9-7831c1c4c038",
			&azure.ResourceID{
				SubscriptionID: "6d74bdd2-9f84-11e5-9bd9-7831c1c4c038",
				ResourceGroup:  "",
				Provider:       "",
				Path:           map[string]string{},
			},
			false,
		},
		{
			"subscriptions/6d74bdd2-9f84-11e5-9bd9-7831c1c4c038",
			nil,
			true,
		},
		{
			"/subscriptions/6d74bdd2-9f84-11e5-9bd9-7831c1c4c038/resourceGroups/testGroup1",
			&azure.ResourceID{
				SubscriptionID: "6d74bdd2-9f84-11e5-9bd9-7831c1c4c038",
				ResourceGroup:  "testGroup1",
				Provider:       "",
				Path:           map[string]string{},
			},
			false,
		},
		{
			"/subscriptions/6d74bdd2-9f84-11e5-9bd9-7831c1c4c038/resourceGroups/testGroup1/providers/Microsoft.Network",
			&azure.ResourceID{
				SubscriptionID: "6d74bdd2-9f84-11e5-9bd9-7831c1c4c038",
				ResourceGroup:  "testGroup1",
				Provider:       "Microsoft.Network",
				Path:           map[string]string{},
			},
			false,
		},
		{
			// Missing leading /
			"subscriptions/6d74bdd2-9f84-11e5-9bd9-7831c1c4c038/resourceGroups/testGroup1/providers/Microsoft.Network/virtualNetworks/virtualNetwork1/",
			nil,
			true,
		},
		{
			"/subscriptions/6d74bdd2-9f84-11e5-9bd9-7831c1c4c038/resourceGroups/testGroup1/providers/Microsoft.Network/virtualNetworks/virtualNetwork1",
			&azure.ResourceID{
				SubscriptionID: "6d74bdd2-9f84-11e5-9bd9-7831c1c4c038",
				ResourceGroup:  "testGroup1",
				Provider:       "Microsoft.Network",
				Path: map[string]string{
					"virtualNetworks": "virtualNetwork1",
				},
			},
			false,
		},
		{
			"/subscriptions/6d74bdd2-9f84-11e5-9bd9-7831c1c4c038/resourceGroups/testGroup1/providers/Microsoft.Network/virtualNetworks/virtualNetwork1?api-version=2006-01-02-preview",
			&azure.ResourceID{
				SubscriptionID: "6d74bdd2-9f84-11e5-9bd9-7831c1c4c038",
				ResourceGroup:  "testGroup1",
				Provider:       "Microsoft.Network",
				Path: map[string]string{
					"virtualNetworks": "virtualNetwork1",
				},
			},
			false,
		},
		{
			"/subscriptions/6d74bdd2-9f84-11e5-9bd9-7831c1c4c038/resourceGroups/testGroup1/providers/Microsoft.Network/virtualNetworks/virtualNetwork1/subnets/publicInstances1?api-version=2006-01-02-preview",
			&azure.ResourceID{
				SubscriptionID: "6d74bdd2-9f84-11e5-9bd9-7831c1c4c038",
				ResourceGroup:  "testGroup1",
				Provider:       "Microsoft.Network",
				Path: map[string]string{
					"virtualNetworks": "virtualNetwork1",
					"subnets":         "publicInstances1",
				},
			},
			false,
		},
		{
			"/subscriptions/34ca515c-4629-458e-bf7c-738d77e0d0ea/resourcegroups/acceptanceTestResourceGroup1/providers/Microsoft.Cdn/profiles/acceptanceTestCdnProfile1",
			&azure.ResourceID{
				SubscriptionID: "34ca515c-4629-458e-bf7c-738d77e0d0ea",
				ResourceGroup:  "acceptanceTestResourceGroup1",
				Provider:       "Microsoft.Cdn",
				Path: map[string]string{
					"profiles": "acceptanceTestCdnProfile1",
				},
			},
			false,
		},
		{
			"/subscriptions/34ca515c-4629-458e-bf7c-738d77e0d0ea/resourceGroups/testGroup1/providers/Microsoft.ServiceBus/namespaces/testNamespace1/topics/testTopic1/subscriptions/testSubscription1",
			&azure.ResourceID{
				SubscriptionID: "34ca515c-4629-458e-bf7c-738d77e0d0ea",
				ResourceGroup:  "testGroup1",
				Provider:       "Microsoft.ServiceBus",
				Path: map[string]string{
					"namespaces":    "testNamespace1",
					"topics":        "testTopic1",
					"subscriptions": "testSubscription1",
				},
			},
			false,
		},
		{
			"/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/example-resources/providers/Microsoft.ApiManagement/service/service1/subscriptions/22222222-2222-2222-2222-222222222222",
			&azure.ResourceID{
				SubscriptionID: "11111111-1111-1111-1111-111111111111",
				ResourceGroup:  "example-resources",
				Provider:       "Microsoft.ApiManagement",
				Path: map[string]string{
					"service":       "service1",
					"subscriptions": "22222222-2222-2222-2222-222222222222",
				},
			},
			false,
		},
		{
			"/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/example-resources/providers/Microsoft.Storage/storageAccounts/nameStorageAccount/providers/Microsoft.Authorization/roleAssignments/22222222-2222-2222-2222-222222222222",
			&azure.ResourceID{
				SubscriptionID:    "11111111-1111-1111-1111-111111111111",
				ResourceGroup:     "example-resources",
				Provider:          "Microsoft.Storage",
				SecondaryProvider: "Microsoft.Authorization",
				Path: map[string]string{
					"storageAccounts": "nameStorageAccount",
					"roleAssignments": "22222222-2222-2222-2222-222222222222",
				},
			},
			false,
		},
		{
			// missing resource group
			"/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.ApiManagement/service/service1/subscriptions/22222222-2222-2222-2222-222222222222",
			&azure.ResourceID{
				SubscriptionID: "11111111-1111-1111-1111-111111111111",
				Provider:       "Microsoft.ApiManagement",
				Path: map[string]string{
					"service":       "service1",
					"subscriptions": "22222222-2222-2222-2222-222222222222",
				},
			},
			false,
		},
	}

	for _, test := range testCases {
		t.Logf("[DEBUG] Testing %q", test.id)
		parsed, err := azure.ParseAzureResourceID(test.id)
		if test.expectError && err != nil {
			continue
		}
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		if !reflect.DeepEqual(test.expectedResourceID, parsed) {
			t.Fatalf("Unexpected resource ID:\nExpected: %+v\nGot:      %+v\n", test.expectedResourceID, parsed)
		}
	}
}

func TestParseAzureResourceIDWithoutSubscription(t *testing.T) {
	testCases := []struct {
		id                 string
		expectedResourceID *azure.ResourceID
		expectError        bool
	}{
		{
			id:          "",
			expectError: true,
		},
		{
			id:          "/providers/Microsoft.Billing/billingAccounts//enrollmentAccounts/123456",
			expectError: true,
		},
		{
			id:          "/providers/Microsoft.Billing/billingAccounts/12345678/enrollmentAccounts",
			expectError: true,
		},
		{
			id: "/providers/Microsoft.Billing/billingAccounts/12345678/enrollmentAccounts/123456",
			expectedResourceID: &azure.ResourceID{
				Provider: "Microsoft.Billing",
				Path: map[string]string{
					"billingAccounts":    "12345678",
					"enrollmentAccounts": "123456",
				},
			},
		},
		{
			id:          "/providers/Microsoft.Management/managementGroups/",
			expectError: true,
		},
		{
			id:          "providers/Microsoft.Management/managementGroups/testManagementGroup",
			expectError: true,
		},
		{
			id:          "/Microsoft.Management/managementGroups/testManagementGroup",
			expectError: true,
		},
		{
			id: "/providers/Microsoft.Management/managementGroups/testManagementGroup",
			expectedResourceID: &azure.ResourceID{
				Provider: "Microsoft.Management",
				Path: map[string]string{
					"managementGroups": "testManagementGroup",
				},
			},
		},
		{
			id: "/providers/microsoft.management/managementGroups/testManagementGroup",
			expectedResourceID: &azure.ResourceID{
				Provider: "microsoft.management",
				Path: map[string]string{
					"managementGroups": "testManagementGroup",
				},
			},
		},
	}

	for _, test := range testCases {
		t.Logf("[DEBUG] Testing %q", test.id)
		parsed, err := azure.ParseAzureResourceIDWithoutSubscription(test.id)
		if test.expectError && err != nil {
			continue
		}
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		if !reflect.DeepEqual(test.expectedResourceID, parsed) {
			t.Fatalf("Unexpected resource ID:\nExpected: %+v\nGot:      %+v\n", test.expectedResourceID, parsed)
		}
	}
}
