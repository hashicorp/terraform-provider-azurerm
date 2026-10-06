// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package managedidentity_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func testAccFederatedIdentityCredential_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_federated_identity_credential", "test")
	r := FederatedIdentityCredentialResource{}
	resourceName := fmt.Sprintf("acctest-0-%d", data.RandomInteger)
	userAssignedIdentityName := fmt.Sprintf("acctestuai-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestrg-%d", data.RandomInteger)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInitWithTestName(context.Background(), t.Name(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.listQuery(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_federated_identity_credential.list", 3),
					querycheck.ExpectIdentity("azurerm_federated_identity_credential.list", map[string]knownvalue.Check{
						"name":                        knownvalue.StringExact(resourceName),
						"resource_group_name":         knownvalue.StringExact(resourceGroupName),
						"subscription_id":             knownvalue.StringExact(data.Subscriptions.Primary),
						"user_assigned_identity_name": knownvalue.StringExact(userAssignedIdentityName),
					}),
				},
			},
		},
	})
}

func (r FederatedIdentityCredentialResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_federated_identity_credential" "test" {
  count = 3

  name                      = "acctest-${count.index}-%[1]d"
  user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  audience                  = ["api://AzureADTokenExchange"]
  issuer                    = "https://token.actions.githubusercontent.com"
  subject                   = "repo:example/repo:environment:test-${count.index}"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r FederatedIdentityCredentialResource) listQuery(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_federated_identity_credential" "list" {
  provider = azurerm
  config {
    user_assigned_identity_id = "/subscriptions/%[1]s/resourceGroups/acctestrg-%[2]d/providers/Microsoft.ManagedIdentity/userAssignedIdentities/acctestuai-%[2]d"
  }
}
`, data.Subscriptions.Primary, data.RandomInteger)
}
