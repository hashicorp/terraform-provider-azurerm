package mssql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MsSqlVirtualNetworkRuleResource struct{}

func TestAccMsSqlVirtualNetworkRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_network_rule", "test")
	r := MsSqlVirtualNetworkRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlVirtualNetworkRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_network_rule", "test")
	r := MsSqlVirtualNetworkRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlVirtualNetworkRule_ignoreMissingServiceEndpoint(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_network_rule", "test")
	r := MsSqlVirtualNetworkRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ignoreMissingServiceEndpoint(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccMsSqlVirtualNetworkRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_network_rule", "test")
	r := MsSqlVirtualNetworkRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r MsSqlVirtualNetworkRuleResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.VirtualNetworkRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQL.VirtualNetworkRulesClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving MSSQL %s: %+v", id.String(), err)
	}

	return utils.Bool(true), nil
}

func (MsSqlVirtualNetworkRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%[1]d"
  address_space       = ["10.7.28.0/23"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test1" {
  name                 = "subnet1%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.7.28.0/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_subnet" "test2" {
  name                 = "subnet2%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.7.28.128/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_subnet" "test3" {
  name                 = "subnet3%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.7.29.0/25"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadmin"
  administrator_login_password = "P@55W0rD!!%[3]s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r MsSqlVirtualNetworkRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_virtual_network_rule" "test" {
  name      = "acctestsqlvnetrule%[2]d"
  server_id = azurerm_mssql_server.test.id
  subnet_id = azurerm_subnet.test1.id
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlVirtualNetworkRuleResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_virtual_network_rule" "test" {
  name      = "acctestsqlvnetrule%[2]d"
  server_id = azurerm_mssql_server.test.id
  subnet_id = azurerm_subnet.test2.id
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlVirtualNetworkRuleResource) ignoreMissingServiceEndpoint(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_virtual_network_rule" "test" {
  name      = "acctestsqlvnetrule%[2]d"
  server_id = azurerm_mssql_server.test.id
  subnet_id = azurerm_subnet.test3.id

  ignore_missing_vnet_service_endpoint = true
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlVirtualNetworkRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_virtual_network_rule" "import" {
  name      = azurerm_mssql_virtual_network_rule.test.name
  server_id = azurerm_mssql_virtual_network_rule.test.server_id
  subnet_id = azurerm_mssql_virtual_network_rule.test.subnet_id
}
`, r.basic(data))
}
