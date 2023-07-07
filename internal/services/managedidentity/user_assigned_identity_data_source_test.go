// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedidentity_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type UserAssignedIdentityDataSource struct{}

func TestAccDataSourceAzureRMUserAssignedIdentity_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_user_assigned_identity", "test")
	d := UserAssignedIdentityDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctest%s-uai", data.RandomString)),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(fmt.Sprintf("acctestRG-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("location").HasValue(azure.NormalizeLocation(data.Locations.Primary)),
				check.That(data.ResourceName).Key("principal_id").IsUUID(),
				check.That(data.ResourceName).Key("client_id").IsUUID(),
				check.That(data.ResourceName).Key("tenant_id").IsUUID(),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("principal_id").MatchesOtherKey(
					check.That("azurerm_user_assigned_identity.test").Key("principal_id"),
				),
				check.That(data.ResourceName).Key("client_id").MatchesOtherKey(
					check.That("azurerm_user_assigned_identity.test").Key("client_id"),
				),
				check.That(data.ResourceName).Key("tenant_id").MatchesOtherKey(
					check.That("azurerm_user_assigned_identity.test").Key("tenant_id"),
				),
			),
		},
	})
}

func (d UserAssignedIdentityDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%s-uai"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    "foo" = "bar"
  }
}

data "azurerm_user_assigned_identity" "test" {
  name                = azurerm_user_assigned_identity.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
