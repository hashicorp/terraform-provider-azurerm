// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mysql_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servers"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccMySqlFlexibleServer_rename(t *testing.T) {
	if os.Getenv("ARM_MYSQL_RENAME_ENABLED") != "true" {
		t.Skip("requires explicit ARM_MYSQL_RENAME_ENABLED=true and a rename-enabled test subscription/region")
	}
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MysqlFlexibleServerResource{}
	renameStep := func(suffix string) acceptance.TestStep {
		name := fmt.Sprintf("acctest-fs-%d%s", data.RandomInteger, suffix)
		id := servers.NewFlexibleServerID(data.Subscriptions.Primary, fmt.Sprintf("acctestRG-mysql-%d", data.RandomInteger), name)
		return acceptance.TestStep{
			Config: r.rename(data, suffix),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").HasValue(id.ID()),
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				acceptance.TestCheckOutput("server_id", id.ID()),
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectIdentityValue(data.ResourceName, tfjsonpath.New("name"), knownvalue.StringExact(name)),
			},
		}
	}
	data.ResourceTest(t, r, []acceptance.TestStep{
		renameStep(""),
		renameStep("-renamed"),
		{
			Config:   r.rename(data, "-renamed"),
			PlanOnly: true,
		},
		data.ImportStep("administrator_password", "rename_enabled"),
		renameStep("-again"),
		data.ImportStep("administrator_password", "rename_enabled"),
	})
}

func (r MysqlFlexibleServerResource) rename(data acceptance.TestData, suffix string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d%s"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "_admin_Terraform_892123456789312"
  administrator_password = "QAZwsx123"
  sku_name               = "B_Standard_B1ms"
  rename_enabled         = true
}

output "server_id" {
  value = azurerm_mysql_flexible_server.test.id
}
`, r.template(data), data.RandomInteger, suffix)
}
