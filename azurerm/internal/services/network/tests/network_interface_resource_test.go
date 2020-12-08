package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMNetworkInterface_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkInterface_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
					testCheckAzureRMNetworkInterfaceDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMNetworkInterface_dnsServers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_dnsServers(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetworkInterface_dnsServersUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkInterface_enableAcceleratedNetworking(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				// Enabled
				Config: testAccAzureRMNetworkInterface_enableAcceleratedNetworking(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				// Disabled
				Config: testAccAzureRMNetworkInterface_enableAcceleratedNetworking(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				// Enabled
				Config: testAccAzureRMNetworkInterface_enableAcceleratedNetworking(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkInterface_enableIPForwarding(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				// Enabled
				Config: testAccAzureRMNetworkInterface_enableIPForwarding(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				// Disabled
				Config: testAccAzureRMNetworkInterface_enableIPForwarding(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				// Enabled
				Config: testAccAzureRMNetworkInterface_enableIPForwarding(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkInterface_internalDomainNameLabel(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_internalDomainNameLabel(data, "1"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetworkInterface_internalDomainNameLabel(data, "2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkInterface_ipv6(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_ipv6(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_configuration.0.private_ip_address_version", "IPv4"),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_configuration.1.private_ip_address_version", "IPv6"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkInterface_multipleIPConfigurations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_multipleIPConfigurations(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkInterface_multipleIPConfigurationsSecondaryAsPrimary(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_multipleIPConfigurationsSecondaryAsPrimary(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkInterface_publicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_publicIP(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetworkInterface_publicIPRemoved(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetworkInterface_publicIP(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkInterface_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMNetworkInterface_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_network_interface"),
			},
		},
	})
}

func TestAccAzureRMNetworkInterface_static(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_static(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkInterface_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_tags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetworkInterface_tagsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkInterface_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetworkInterface_multipleIPConfigurations(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkInterface_updateMultipleParameters(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_withMultipleParameters(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetworkInterface_updateMultipleParameters(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMNetworkInterfaceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.InterfacesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Network Interface %q (resource group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on ifaceClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMNetworkInterfaceDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.InterfacesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		future, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Error deleting Network Interface %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for the deletion of Network Interface %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMNetworkInterfaceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.InterfacesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_network_interface" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Network Interface still exists:\n%#v", resp.InterfacePropertiesFormat)
	}

	return nil
}

func testAccAzureRMNetworkInterface_basic(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterface_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_withMultipleParameters(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterface_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                    = "acctestni-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  enable_ip_forwarding    = true
  internal_dns_name_label = "acctestni-%s"

  dns_servers = [
    "10.0.0.5",
    "10.0.0.6"
  ]

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    env = "Test"
  }
}
`, template, data.RandomInteger, data.RandomString)
}

func testAccAzureRMNetworkInterface_updateMultipleParameters(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterface_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                    = "acctestni-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  enable_ip_forwarding    = true
  internal_dns_name_label = "acctestni-%s"

  dns_servers = [
    "10.0.0.5",
    "10.0.0.7"
  ]

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    env = "Test2"
  }
}
`, template, data.RandomInteger, data.RandomString)
}

func testAccAzureRMNetworkInterface_dnsServers(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterface_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  dns_servers = [
    "10.0.0.5",
    "10.0.0.6"
  ]

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_dnsServersUpdated(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterface_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  dns_servers = [
    "10.0.0.6",
    "10.0.0.5"
  ]

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_enableAcceleratedNetworking(data acceptance.TestData, enabled bool) string {
	template := testAccAzureRMNetworkInterface_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                          = "acctestni-%d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  enable_accelerated_networking = %t

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, template, data.RandomInteger, enabled)
}

func testAccAzureRMNetworkInterface_enableIPForwarding(data acceptance.TestData, enabled bool) string {
	template := testAccAzureRMNetworkInterface_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                 = "acctestni-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  enable_ip_forwarding = %t

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, template, data.RandomInteger, enabled)
}

func testAccAzureRMNetworkInterface_internalDomainNameLabel(data acceptance.TestData, suffix string) string {
	template := testAccAzureRMNetworkInterface_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                    = "acctestni-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  internal_dns_name_label = "acctestni-%s-%s"

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, template, data.RandomInteger, suffix, data.RandomString)
}

func testAccAzureRMNetworkInterface_ipv6(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterface_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    primary                       = true
  }

  ip_configuration {
    name                          = "secondary"
    private_ip_address_allocation = "Dynamic"
    private_ip_address_version    = "IPv6"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_multipleIPConfigurations(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterface_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    primary                       = true
  }

  ip_configuration {
    name                          = "secondary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_multipleIPConfigurationsSecondaryAsPrimary(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterface_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }

  ip_configuration {
    name                          = "secondary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    primary                       = true
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_publicIP(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterface_publicIPTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_publicIPRemoved(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterface_publicIPTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_publicIPTemplate(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterface_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}
`, template, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterface_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "import" {
  name                = azurerm_network_interface.test.name
  location            = azurerm_network_interface.test.location
  resource_group_name = azurerm_network_interface.test.resource_group_name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, template)
}

func testAccAzureRMNetworkInterface_static(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterface_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Static"
    private_ip_address            = "10.0.2.15"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_tags(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterface_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    Hello = "World"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_tagsUpdated(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterface_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    Hello     = "World"
    Elephants = "Five"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
