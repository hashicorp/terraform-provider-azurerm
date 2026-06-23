// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package managedidentity_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccUserAssignedIdentity_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_user_assigned_identity", "test")
	r := UserAssignedIdentityTestResource{}
	resourceName := fmt.Sprintf("acctestuai-0-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestrg-%d", data.RandomInteger)

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
				// Wait for the subscription-level User Assigned Identity listing API to become eventually consistent after resource creation.
				PreConfig: func() { time.Sleep(2 * time.Minute) },
				Query:     true,
				Config:    r.subscriptionListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_user_assigned_identity.list", 3),
					querycheck.ExpectIdentity("azurerm_user_assigned_identity.list", map[string]knownvalue.Check{
						"name":                knownvalue.StringExact(resourceName),
						"resource_group_name": knownvalue.StringExact(resourceGroupName),
						"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
					}),
				},
			},
			{
				Query:  true,
				Config: r.resourceGroupListQuery(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_user_assigned_identity.list", 3),
					querycheck.ExpectIdentity("azurerm_user_assigned_identity.list", map[string]knownvalue.Check{
						"name":                knownvalue.StringExact(resourceName),
						"resource_group_name": knownvalue.StringExact(resourceGroupName),
						"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
					}),
				},
			},
		},
	})
}

func (UserAssignedIdentityTestResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  count = 3

  name                = "acctestuai-${count.index}-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (UserAssignedIdentityTestResource) subscriptionListQuery() string {
	return `
list "azurerm_user_assigned_identity" "list" {
  provider = azurerm
}
`
}

func (UserAssignedIdentityTestResource) resourceGroupListQuery(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_user_assigned_identity" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctestrg-%[1]d"
  }
}
`, data.RandomInteger)
}
