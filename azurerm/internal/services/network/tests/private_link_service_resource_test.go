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

func TestAccAzureRMPrivateLinkService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_link_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateLinkServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateLinkService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctestPLS-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "load_balancer_frontend_ip_configuration_ids.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPrivateLinkService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_link_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateLinkServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateLinkService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMPrivateLinkService_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_private_link_service"),
			},
		},
	})
}

func TestAccAzureRMPrivateLinkService_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_link_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateLinkServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateLinkService_basicIp(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "load_balancer_frontend_ip_configuration_ids.#", "1"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPrivateLinkService_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.#", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.0.primary", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "load_balancer_frontend_ip_configuration_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.env", "test"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPrivateLinkService_basicIp(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "load_balancer_frontend_ip_configuration_ids.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPrivateLinkService_move(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_link_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateLinkServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateLinkService_moveSetup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.0.private_ip_address", "10.5.2.17"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPrivateLinkService_moveAdd(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.#", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.0.private_ip_address", "10.5.2.17"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.1.private_ip_address", "10.5.2.18"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.2.private_ip_address", "10.5.2.19"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.3.private_ip_address", "10.5.2.20"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPrivateLinkService_moveChangeOne(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.#", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.0.private_ip_address", "10.5.2.17"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.1.private_ip_address", "10.5.2.18"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.2.private_ip_address", "10.5.2.19"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.3.private_ip_address", "10.5.2.21"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPrivateLinkService_moveChangeTwo(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.#", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.0.private_ip_address", "10.5.2.17"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.1.private_ip_address", "10.5.2.20"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.2.private_ip_address", "10.5.2.19"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.3.private_ip_address", "10.5.2.21"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPrivateLinkService_moveChangeThree(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.#", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.0.private_ip_address", "10.5.2.17"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.1.private_ip_address", "10.5.2.20"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.2.private_ip_address", "10.5.2.19"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.3.private_ip_address", "10.5.2.18"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPrivateLinkService_enableProxyProtocol(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_link_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateLinkServiceDestroy,
		Steps: []resource.TestStep{
			{
				// Enable
				Config: testAccAzureRMPrivateLinkService_enableProxyProtocol(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				// Disable
				Config: testAccAzureRMPrivateLinkService_enableProxyProtocol(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				// Enable
				Config: testAccAzureRMPrivateLinkService_enableProxyProtocol(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPrivateLinkService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_link_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateLinkServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateLinkService_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_approval_subscription_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "visibility_subscription_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.0.private_ip_address", "10.5.1.40"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.0.private_ip_address_version", "IPv4"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.1.private_ip_address", "10.5.1.41"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.1.private_ip_address_version", "IPv4"),
					resource.TestCheckResourceAttr(data.ResourceName, "load_balancer_frontend_ip_configuration_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.env", "test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMPrivateLinkServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.PrivateLinkServiceClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Private Link Service not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Private Link Service %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on network.PrivateLinkServiceClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMPrivateLinkServiceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.PrivateLinkServiceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_private_link_service" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name, ""); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on network.PrivateLinkServiceClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMPrivateLinkService_basic(data acceptance.TestData) string {
	template := testAccAzureRMPrivateLinkService_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name               = "acctestsnet-basic-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.5.4.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                = "acctestPLS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  nat_ip_configuration {
    name      = "primaryIpConfiguration-%d"
    subnet_id = azurerm_subnet.test.id
    primary   = true
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateLinkService_basicIp(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name               = "acctestsnet-update-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.5.3.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                = "acctestPLS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  nat_ip_configuration {
    name                       = "primaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.3.30"
    private_ip_address_version = "IPv4"
    primary                    = true
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}
`, testAccAzureRMPrivateLinkService_template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateLinkService_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMPrivateLinkService_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_private_link_service" "import" {
  name                = azurerm_private_link_service.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  nat_ip_configuration {
    name      = "primaryIpConfiguration-%d"
    subnet_id = azurerm_subnet.test.id
    primary   = true
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}
`, template, data.RandomInteger)
}

func testAccAzureRMPrivateLinkService_enableProxyProtocol(data acceptance.TestData, enabled bool) string {
	template := testAccAzureRMPrivateLinkService_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name               = "acctestsnet-basic-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.5.4.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                  = "acctestPLS-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  enable_proxy_protocol = %t

  nat_ip_configuration {
    name      = "primaryIpConfiguration-%d"
    subnet_id = azurerm_subnet.test.id
    primary   = true
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}
`, template, data.RandomInteger, data.RandomInteger, enabled, data.RandomInteger)
}

func testAccAzureRMPrivateLinkService_update(data acceptance.TestData) string {
	template := testAccAzureRMPrivateLinkService_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name               = "acctestsnet-update-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.5.3.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                       = "primaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.3.30"
    private_ip_address_version = "IPv4"
    primary                    = true
  }

  nat_ip_configuration {
    name                       = "secondaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.3.22"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  nat_ip_configuration {
    name                       = "thirdaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.3.23"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  nat_ip_configuration {
    name                       = "fourtharyIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.3.24"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateLinkService_moveSetup(data acceptance.TestData) string {
	template := testAccAzureRMPrivateLinkService_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name               = "acctestsnet-move-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.5.2.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                       = "primaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.17"
    private_ip_address_version = "IPv4"
    primary                    = true
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateLinkService_moveAdd(data acceptance.TestData) string {
	template := testAccAzureRMPrivateLinkService_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name               = "acctestsnet-move-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.5.2.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                       = "primaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.17"
    private_ip_address_version = "IPv4"
    primary                    = true
  }

  nat_ip_configuration {
    name                       = "secondaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.18"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  nat_ip_configuration {
    name                       = "thirdaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.19"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  nat_ip_configuration {
    name                       = "fourtharyIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.20"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateLinkService_moveChangeOne(data acceptance.TestData) string {
	template := testAccAzureRMPrivateLinkService_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name               = "acctestsnet-move-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.5.2.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                       = "primaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.17"
    private_ip_address_version = "IPv4"
    primary                    = true
  }

  nat_ip_configuration {
    name                       = "secondaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.18"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  nat_ip_configuration {
    name                       = "thirdaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.19"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  nat_ip_configuration {
    name                       = "fourtharyIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.21"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateLinkService_moveChangeTwo(data acceptance.TestData) string {
	template := testAccAzureRMPrivateLinkService_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name               = "acctestsnet-move-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.5.2.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                       = "primaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.17"
    private_ip_address_version = "IPv4"
    primary                    = true
  }

  nat_ip_configuration {
    name                       = "secondaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.20"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  nat_ip_configuration {
    name                       = "thirdaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.19"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  nat_ip_configuration {
    name                       = "fourtharyIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.21"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateLinkService_moveChangeThree(data acceptance.TestData) string {
	template := testAccAzureRMPrivateLinkService_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name               = "acctestsnet-move-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.5.2.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                       = "primaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.17"
    private_ip_address_version = "IPv4"
    primary                    = true
  }

  nat_ip_configuration {
    name                       = "secondaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.20"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  nat_ip_configuration {
    name                       = "thirdaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.19"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  nat_ip_configuration {
    name                       = "fourtharyIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.18"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateLinkService_complete(data acceptance.TestData) string {
	template := testAccAzureRMPrivateLinkService_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name               = "acctestsnet-complete-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.5.1.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                       = "primaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.1.40"
    private_ip_address_version = "IPv4"
    primary                    = true
  }

  nat_ip_configuration {
    name                       = "secondaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.1.41"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateLinkService_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-privatelinkservice-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  sku                 = "Standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  sku                 = "Standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = azurerm_public_ip.test.name
    public_ip_address_id = azurerm_public_ip.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
