package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMVirtualHubConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualHubConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubConnectionExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMVirtualHubConnection_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_virtual_hub_connection"),
			},
		},
	})
}

func TestAccAzureRMVirtualHubConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubConnection_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualHubConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVirtualHubConnection_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVirtualHubConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualHubConnection_recreateWithSameConnectionName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVirtualHubConnection_enableInternetSecurity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMVirtualHubConnectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.HubVirtualNetworkConnectionClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Virtual Hub Connection not found: %s", resourceName)
		}

		id, err := parse.VirtualHubConnectionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Virtual Hub Connection %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Bad: Get on network.HubVirtualNetworkConnectionClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMVirtualHubConnectionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.HubVirtualNetworkConnectionClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_hub_connection" {
			continue
		}

		id, err := parse.VirtualHubConnectionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on network.HubVirtualNetworkConnectionClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMVirtualHubConnection_basic(data acceptance.TestData) string {
	template := testAccAzureRMVirtualHubConnection_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctestbasicvhubconn-%d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id
}
`, template, data.RandomInteger)
}

func testAccAzureRMVirtualHubConnection_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMVirtualHubConnection_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_connection" "import" {
  name                      = azurerm_virtual_hub_connection.test.name
  virtual_hub_id            = azurerm_virtual_hub_connection.test.virtual_hub_id
  remote_virtual_network_id = azurerm_virtual_hub_connection.test.remote_virtual_network_id
}
`, template)
}

func testAccAzureRMVirtualHubConnection_complete(data acceptance.TestData) string {
	template := testAccAzureRMVirtualHubConnection_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "test2" {
  name                = "acctestvirtnet2%d"
  address_space       = ["10.6.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_group" "test2" {
  name                = "acctestnsg2%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test2" {
  name               = "acctestsubnet2%d"
  virtual_network_id = azurerm_virtual_network.test2.id
  address_prefixes   = ["10.6.1.0/24"]
}

resource "azurerm_subnet_network_security_group_association" "test2" {
  subnet_id                 = azurerm_subnet.test2.id
  network_security_group_id = azurerm_network_security_group.test2.id
}

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctestvhubconn-%d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id
  internet_security_enabled = false
}

resource "azurerm_virtual_hub_connection" "test2" {
  name                      = "acctestvhubconn2-%d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test2.id
  internet_security_enabled = true
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualHubConnection_enableInternetSecurity(data acceptance.TestData) string {
	template := testAccAzureRMVirtualHubConnection_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctestbasicvhubconn-%d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id
  internet_security_enabled = true
}
`, template, data.RandomInteger)
}

func testAccAzureRMVirtualHubConnection_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vhub-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.5.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name               = "acctestsubnet%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefixes   = ["10.5.1.0/24"]
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctest-VHUB-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.2.0/24"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
