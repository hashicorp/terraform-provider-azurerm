package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSubnet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSubnet_basic_addressPrefixes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnet_basic_addressPrefixes(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSubnet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMSubnet_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_subnet"),
			},
		},
	})
}

func TestAccAzureRMSubnet_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
					testCheckAzureRMSubnetDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMSubnet_delegation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnet_delegation(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSubnet_delegationUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSubnet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSubnet_delegation(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSubnet_enforcePrivateLinkEndpointNetworkPolicies(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnet_enforcePrivateLinkEndpointNetworkPolicies(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSubnet_enforcePrivateLinkEndpointNetworkPolicies(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSubnet_enforcePrivateLinkEndpointNetworkPolicies(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSubnet_enforcePrivateLinkServiceNetworkPolicies(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnet_enforcePrivateLinkServiceNetworkPolicies(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSubnet_enforcePrivateLinkServiceNetworkPolicies(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSubnet_enforcePrivateLinkServiceNetworkPolicies(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSubnet_serviceEndpoints(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnet_serviceEndpoints(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSubnet_serviceEndpointsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				// remove them
				Config: testAccAzureRMSubnet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSubnet_serviceEndpoints(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSubnet_serviceEndpointPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnet_serviceEndpointPolicyBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSubnet_serviceEndpointPolicyUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSubnet_serviceEndpointPolicyBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMSubnet_updateAddressPrefix(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSubnet_updatedAddressPrefix(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMSubnetExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.SubnetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		log.Printf("[INFO] Checking Subnet addition.")

		name := rs.Primary.Attributes["name"]
		vnetName := rs.Primary.Attributes["virtual_network_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for subnet: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, vnetName, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Subnet %q (resource group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on subnetClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMSubnetDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.SubnetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		vnetName := rs.Primary.Attributes["virtual_network_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for subnet: %s", name)
		}

		future, err := client.Delete(ctx, resourceGroup, vnetName, name)
		if err != nil {
			if !response.WasNotFound(future.Response()) {
				return fmt.Errorf("Error deleting Subnet %q (Network %q / Resource Group %q): %+v", name, vnetName, resourceGroup, err)
			}
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for completion of Subnet %q (Network %q / Resource Group %q): %+v", name, vnetName, resourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMSubnetDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.SubnetsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_subnet" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		vnetName := rs.Primary.Attributes["virtual_network_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, vnetName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Subnet still exists:\n%#v", resp.SubnetPropertiesFormat)
			}
			return nil
		}
	}

	return nil
}

func testAccAzureRMSubnet_basic(data acceptance.TestData) string {
	template := testAccAzureRMSubnet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}
`, template)
}

func testAccAzureRMSubnet_delegation(data acceptance.TestData) string {
	template := testAccAzureRMSubnet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"

  delegation {
    name = "first"

    service_delegation {
      name = "Microsoft.ContainerInstance/containerGroups"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/action",
      ]
    }
  }
}
`, template)
}

func testAccAzureRMSubnet_delegationUpdated(data acceptance.TestData) string {
	template := testAccAzureRMSubnet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"

  delegation {
    name = "first"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }
}
`, template)
}

func testAccAzureRMSubnet_enforcePrivateLinkEndpointNetworkPolicies(data acceptance.TestData, enabled bool) string {
	template := testAccAzureRMSubnet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"

  enforce_private_link_endpoint_network_policies = %t
}
`, template, enabled)
}

func testAccAzureRMSubnet_enforcePrivateLinkServiceNetworkPolicies(data acceptance.TestData, enabled bool) string {
	template := testAccAzureRMSubnet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"

  enforce_private_link_service_network_policies = %t
}
`, template, enabled)
}

func testAccAzureRMSubnet_basic_addressPrefixes(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-n-%d"
  location = "%s"
}
resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16", "ace:cab:deca::/48"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefixes     = ["10.0.0.0/24", "ace:cab:deca:deed::/64"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMSubnet_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMSubnet_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "import" {
  name                 = azurerm_subnet.test.name
  resource_group_name  = azurerm_subnet.test.resource_group_name
  virtual_network_name = azurerm_subnet.test.virtual_network_name
  address_prefix       = azurerm_subnet.test.address_prefix
}
`, template)
}

func testAccAzureRMSubnet_serviceEndpoints(data acceptance.TestData) string {
	template := testAccAzureRMSubnet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.Sql"]
}
`, template)
}

func testAccAzureRMSubnet_serviceEndpointsUpdated(data acceptance.TestData) string {
	template := testAccAzureRMSubnet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.Sql", "Microsoft.Storage"]
}
`, template)
}

func testAccAzureRMSubnet_serviceEndpointPolicyBasic(data acceptance.TestData) string {
	template := testAccAzureRMSubnet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet_service_endpoint_storage_policy" "test" {
  name                = "acctestSEP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}
`, template, data.RandomInteger)
}

func testAccAzureRMSubnet_serviceEndpointPolicyUpdate(data acceptance.TestData) string {
	template := testAccAzureRMSubnet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet_service_endpoint_storage_policy" "test" {
  name                = "acctestSEP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_subnet" "test" {
  name                        = "internal"
  resource_group_name         = azurerm_resource_group.test.name
  virtual_network_name        = azurerm_virtual_network.test.name
  address_prefix              = "10.0.2.0/24"
  service_endpoints           = ["Microsoft.Sql"]
  service_endpoint_policy_ids = [azurerm_subnet_service_endpoint_storage_policy.test.id]
}
`, template, data.RandomInteger)
}

func testAccAzureRMSubnet_updatedAddressPrefix(data acceptance.TestData) string {
	template := testAccAzureRMSubnet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.3.0/24"
}
`, template)
}

func testAccAzureRMSubnet_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
