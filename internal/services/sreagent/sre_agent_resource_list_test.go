// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sreagent_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-01-01/agents"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccSreAgent_list_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sre_agent", "test")
	r := SreAgentResource{}
	listResourceAddress := "azurerm_sre_agent.list"
	name := fmt.Sprintf("acctest-sre-0-%d", data.RandomIntOfLength(8))
	resourceGroupName := fmt.Sprintf("acctest-rg-sre-%d", data.RandomInteger)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				ResourceName:      data.ResourceName + "[0]",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      data.ResourceName + "[1]",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Query:  true,
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 2),
					querycheck.ExpectIdentity(listResourceAddress, map[string]knownvalue.Check{
						"name":                knownvalue.StringExact(name),
						"resource_group_name": knownvalue.StringExact(resourceGroupName),
						"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
					}),
				},
			},
			{
				Query:  true,
				Config: r.queryBySubscription(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 2),
					querycheck.ExpectIdentity(listResourceAddress, map[string]knownvalue.Check{
						"name":                knownvalue.StringExact(name),
						"resource_group_name": knownvalue.StringExact(resourceGroupName),
						"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
					}),
				},
			},
			{
				Query:  true,
				Config: r.queryByResourceGroup(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 2),
					querycheck.ExpectResourceKnownValues(
						listResourceAddress,
						queryfilter.ByDisplayName(knownvalue.StringExact(name)),
						[]querycheck.KnownValueCheck{
							{
								Path:       tfjsonpath.New("id"),
								KnownValue: knownvalue.StringExact(agents.NewAgentID(data.Subscriptions.Primary, resourceGroupName, name).ID()),
							},
							{
								Path:       tfjsonpath.New("name"),
								KnownValue: knownvalue.StringExact(name),
							},
							{
								Path:       tfjsonpath.New("resource_group_name"),
								KnownValue: knownvalue.StringExact(resourceGroupName),
							},
							{
								Path:       tfjsonpath.New("location"),
								KnownValue: knownvalue.StringExact(location.Normalize(data.Locations.Primary)),
							},
						},
					),
				},
			},
		},
	})
}

func TestAccSreAgent_list_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sre_agent", "test")
	r := SreAgentResource{}
	listResourceAddress := "azurerm_sre_agent.list"
	name := fmt.Sprintf("acctest-sre-%d", data.RandomIntOfLength(8))
	resourceGroupName := fmt.Sprintf("acctest-rg-sre-%d", data.RandomInteger)
	identityPrefix := commonids.NewResourceGroupID(data.Subscriptions.Primary, resourceGroupName).ID() + "/providers/Microsoft.ManagedIdentity/userAssignedIdentities/"
	actionsIdentity := identityPrefix + fmt.Sprintf("acctest-sre-actions-%d", data.RandomIntOfLength(8))
	resourcesIdentity := identityPrefix + fmt.Sprintf("acctest-sre-resources-%d", data.RandomIntOfLength(8))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.complete(data),
			},
			{
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Query:  true,
				Config: r.queryByResourceGroup(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 1),
					querycheck.ExpectIdentity(listResourceAddress, map[string]knownvalue.Check{
						"name":                knownvalue.StringExact(name),
						"resource_group_name": knownvalue.StringExact(resourceGroupName),
						"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
					}),
					querycheck.ExpectResourceKnownValues(
						listResourceAddress,
						queryfilter.ByDisplayName(knownvalue.StringExact(name)),
						[]querycheck.KnownValueCheck{
							{
								Path:       tfjsonpath.New("identity").AtSliceIndex(0).AtMapKey("type"),
								KnownValue: knownvalue.StringExact("UserAssigned"),
							},
							{
								Path: tfjsonpath.New("identity").AtSliceIndex(0).AtMapKey("identity_ids"),
								KnownValue: knownvalue.SetExact([]knownvalue.Check{
									knownvalue.StringExact(actionsIdentity),
									knownvalue.StringExact(resourcesIdentity),
								}),
							},
							{
								Path:       tfjsonpath.New("action_configuration").AtSliceIndex(0).AtMapKey("identity_id"),
								KnownValue: knownvalue.StringExact(actionsIdentity),
							},
							{
								Path:       tfjsonpath.New("action_configuration").AtSliceIndex(0).AtMapKey("mode"),
								KnownValue: knownvalue.StringExact("Review"),
							},
							{
								Path:       tfjsonpath.New("action_configuration").AtSliceIndex(0).AtMapKey("access_level"),
								KnownValue: knownvalue.StringExact("Low"),
							},
							{
								Path:       tfjsonpath.New("resources_configuration").AtSliceIndex(0).AtMapKey("identity_id"),
								KnownValue: knownvalue.StringExact(resourcesIdentity),
							},
							{
								Path: tfjsonpath.New("resources_configuration").AtSliceIndex(0).AtMapKey("resource_ids"),
								KnownValue: knownvalue.SetExact([]knownvalue.Check{
									knownvalue.StringExact(commonids.NewResourceGroupID(data.Subscriptions.Primary, fmt.Sprintf("acctest-rg-sre-managed-0-%d", data.RandomInteger)).ID()),
								}),
							},
							{
								Path:       tfjsonpath.New("tags").AtMapKey("stage"),
								KnownValue: knownvalue.StringExact("initial"),
							},
						},
					),
				},
			},
		},
	})
}

func (r SreAgentResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sre_agent" "test" {
  count = 2

  name                = "acctest-sre-${count.index}-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, r.template(data), data.RandomIntOfLength(8))
}

func (SreAgentResource) basicQuery() string {
	return `
list "azurerm_sre_agent" "list" {
  provider = azurerm
  config {}
}
`
}

func (SreAgentResource) queryBySubscription(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_sre_agent" "list" {
  provider = azurerm
  config {
    subscription_id = "%s"
  }
}
`, data.Subscriptions.Primary)
}

func (SreAgentResource) queryByResourceGroup() string {
	return `
list "azurerm_sre_agent" "list" {
  provider         = azurerm
  include_resource = true
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`
}
