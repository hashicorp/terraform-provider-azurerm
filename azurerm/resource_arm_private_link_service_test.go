package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPrivateLinkService_basic(t *testing.T) {
	resourceName := "azurerm_private_link_service.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateLinkServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateLinkService_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctestPLS-%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_frontend_ip_configuration_ids.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMPrivateLinkService_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_private_link_service.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateLinkServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateLinkService_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMPrivateLinkService_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_private_link_service"),
			},
		},
	})
}

func TestAccAzureRMPrivateLinkService_update(t *testing.T) {
	resourceName := "azurerm_private_link_service.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateLinkServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateLinkService_basicIp(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_frontend_ip_configuration_ids.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMPrivateLinkService_update(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.0.primary", "true"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_frontend_ip_configuration_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMPrivateLinkService_basicIp(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_frontend_ip_configuration_ids.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMPrivateLinkService_move(t *testing.T) {
	resourceName := "azurerm_private_link_service.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateLinkServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateLinkService_moveSetup(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.0.private_ip_address", "10.5.2.17"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMPrivateLinkService_moveAdd(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.0.private_ip_address", "10.5.2.17"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.1.private_ip_address", "10.5.2.18"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.2.private_ip_address", "10.5.2.19"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.3.private_ip_address", "10.5.2.20"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMPrivateLinkService_moveChangeOne(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.0.private_ip_address", "10.5.2.17"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.1.private_ip_address", "10.5.2.18"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.2.private_ip_address", "10.5.2.19"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.3.private_ip_address", "10.5.2.21"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMPrivateLinkService_moveChangeTwo(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.0.private_ip_address", "10.5.2.17"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.1.private_ip_address", "10.5.2.20"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.2.private_ip_address", "10.5.2.19"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.3.private_ip_address", "10.5.2.21"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMPrivateLinkService_moveChangeThree(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.0.private_ip_address", "10.5.2.17"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.1.private_ip_address", "10.5.2.20"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.2.private_ip_address", "10.5.2.19"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.3.private_ip_address", "10.5.2.18"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMPrivateLinkService_enableProxyProtocol(t *testing.T) {
	resourceName := "azurerm_private_link_service.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateLinkServiceDestroy,
		Steps: []resource.TestStep{
			{
				// Enable
				Config: testAccAzureRMPrivateLinkService_enableProxyProtocol(ri, location, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Disable
				Config: testAccAzureRMPrivateLinkService_enableProxyProtocol(ri, location, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Enable
				Config: testAccAzureRMPrivateLinkService_enableProxyProtocol(ri, location, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMPrivateLinkService_complete(t *testing.T) {
	resourceName := "azurerm_private_link_service.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateLinkServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateLinkService_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auto_approval_subscription_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "visibility_subscription_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.0.private_ip_address", "10.5.1.40"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.0.private_ip_address_version", "IPv4"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.1.private_ip_address", "10.5.1.41"),
					resource.TestCheckResourceAttr(resourceName, "nat_ip_configuration.1.private_ip_address_version", "IPv4"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_frontend_ip_configuration_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMPrivateLinkServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Private Link Service not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.PrivateLinkServiceClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMPrivateLinkService_basic(rInt int, location string) string {
	template := testAccAzureRMPrivateLinkService_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                                  = "acctestsnet-basic-%d"
  resource_group_name                   = azurerm_resource_group.test.name
  virtual_network_name                  = azurerm_virtual_network.test.name
  address_prefix                        = "10.5.4.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name

  nat_ip_configuration {
    name                         = "primaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    primary                      = true
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}
`, template, rInt, rInt, rInt)
}

func testAccAzureRMPrivateLinkService_basicIp(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                                  = "acctestsnet-update-%d"
  resource_group_name                   = azurerm_resource_group.test.name
  virtual_network_name                  = azurerm_virtual_network.test.name
  address_prefix                        = "10.5.3.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name

  nat_ip_configuration {
    name                         = "primaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.3.30"
    private_ip_address_version   = "IPv4"
    primary                      = true
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}
`, testAccAzureRMPrivateLinkService_template(rInt, location), rInt, rInt, rInt)
}

func testAccAzureRMPrivateLinkService_requiresImport(rInt int, location string) string {
	template := testAccAzureRMPrivateLinkService_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_private_link_service" "import" {
  name                           = azurerm_private_link_service.test.name
  location                       = azurerm_private_link_service.test.location
  resource_group_name            = azurerm_private_link_service.test.name

  nat_ip_configuration {
    name                         = "primaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    primary                      = true
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}
`, template, rInt)
}

func testAccAzureRMPrivateLinkService_enableProxyProtocol(rInt int, location string, enabled bool) string {
	template := testAccAzureRMPrivateLinkService_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                                  = "acctestsnet-basic-%d"
  resource_group_name                   = azurerm_resource_group.test.name
  virtual_network_name                  = azurerm_virtual_network.test.name
  address_prefix                        = "10.5.4.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  enable_proxy_protocol          = %t

  nat_ip_configuration {
    name                         = "primaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    primary                      = true
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}
`, template, rInt, rInt, enabled, rInt)
}

func testAccAzureRMPrivateLinkService_update(rInt int, location string) string {
	template := testAccAzureRMPrivateLinkService_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                                  = "acctestsnet-update-%d"
  resource_group_name                   = azurerm_resource_group.test.name
  virtual_network_name                  = azurerm_virtual_network.test.name
  address_prefix                        = "10.5.3.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                         = "primaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.3.30"
    private_ip_address_version   = "IPv4"
    primary                      = true
  }

  nat_ip_configuration {
    name                         = "secondaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.3.22"
    private_ip_address_version   = "IPv4"
    primary                      = false
  }

  nat_ip_configuration {
    name                         = "thirdaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.3.23"
    private_ip_address_version   = "IPv4"
    primary                      = false
  }

  nat_ip_configuration {
    name                         = "fourtharyIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.3.24"
    private_ip_address_version   = "IPv4"
    primary                      = false
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, template, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMPrivateLinkService_moveSetup(rInt int, location string) string {
	template := testAccAzureRMPrivateLinkService_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                                  = "acctestsnet-move-%d"
  resource_group_name                   = azurerm_resource_group.test.name
  virtual_network_name                  = azurerm_virtual_network.test.name
  address_prefix                        = "10.5.2.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                         = "primaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.2.17"
    private_ip_address_version   = "IPv4"
    primary                      = true
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, template, rInt, rInt, rInt)
}

func testAccAzureRMPrivateLinkService_moveAdd(rInt int, location string) string {
	template := testAccAzureRMPrivateLinkService_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                                  = "acctestsnet-move-%d"
  resource_group_name                   = azurerm_resource_group.test.name
  virtual_network_name                  = azurerm_virtual_network.test.name
  address_prefix                        = "10.5.2.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                         = "primaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.2.17"
    private_ip_address_version   = "IPv4"
    primary                      = true
  }

  nat_ip_configuration {
    name                         = "secondaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.2.18"
    private_ip_address_version   = "IPv4"
    primary                      = false
  }

  nat_ip_configuration {
    name                         = "thirdaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.2.19"
    private_ip_address_version   = "IPv4"
    primary                      = false
  }

  nat_ip_configuration {
    name                         = "fourtharyIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.2.20"
    private_ip_address_version   = "IPv4"
    primary                      = false
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, template, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMPrivateLinkService_moveChangeOne(rInt int, location string) string {
	template := testAccAzureRMPrivateLinkService_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                                  = "acctestsnet-move-%d"
  resource_group_name                   = azurerm_resource_group.test.name
  virtual_network_name                  = azurerm_virtual_network.test.name
  address_prefix                        = "10.5.2.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                         = "primaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.2.17"
    private_ip_address_version   = "IPv4"
    primary                      = true
  }

  nat_ip_configuration {
    name                         = "secondaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.2.18"
    private_ip_address_version   = "IPv4"
    primary                      = false
  }

  nat_ip_configuration {
    name                         = "thirdaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.2.19"
    private_ip_address_version   = "IPv4"
    primary                      = false
  }

  nat_ip_configuration {
    name                         = "fourtharyIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.2.21"
    private_ip_address_version   = "IPv4"
    primary                      = false
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, template, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMPrivateLinkService_moveChangeTwo(rInt int, location string) string {
	template := testAccAzureRMPrivateLinkService_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                                  = "acctestsnet-move-%d"
  resource_group_name                   = azurerm_resource_group.test.name
  virtual_network_name                  = azurerm_virtual_network.test.name
  address_prefix                        = "10.5.2.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                         = "primaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.2.17"
    private_ip_address_version   = "IPv4"
    primary                      = true
  }

  nat_ip_configuration {
    name                         = "secondaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.2.20"
    private_ip_address_version   = "IPv4"
    primary                      = false
  }

  nat_ip_configuration {
    name                         = "thirdaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.2.19"
    private_ip_address_version   = "IPv4"
    primary                      = false
  }

  nat_ip_configuration {
    name                         = "fourtharyIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.2.21"
    private_ip_address_version   = "IPv4"
    primary                      = false
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, template, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMPrivateLinkService_moveChangeThree(rInt int, location string) string {
	template := testAccAzureRMPrivateLinkService_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                                  = "acctestsnet-move-%d"
  resource_group_name                   = azurerm_resource_group.test.name
  virtual_network_name                  = azurerm_virtual_network.test.name
  address_prefix                        = "10.5.2.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                         = "primaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.2.17"
    private_ip_address_version   = "IPv4"
    primary                      = true
  }

  nat_ip_configuration {
    name                         = "secondaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.2.20"
    private_ip_address_version   = "IPv4"
    primary                      = false
  }

  nat_ip_configuration {
    name                         = "thirdaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.2.19"
    private_ip_address_version   = "IPv4"
    primary                      = false
  }

  nat_ip_configuration {
    name                         = "fourtharyIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.2.18"
    private_ip_address_version   = "IPv4"
    primary                      = false
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, template, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMPrivateLinkService_complete(rInt int, location string) string {
	template := testAccAzureRMPrivateLinkService_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                                  = "acctestsnet-complete-%d"
  resource_group_name                   = azurerm_resource_group.test.name
  virtual_network_name                  = azurerm_virtual_network.test.name
  address_prefix                        = "10.5.1.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                         = "primaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.1.40"
    private_ip_address_version   = "IPv4"
    primary                      = true
  }

  nat_ip_configuration {
    name                         = "secondaryIpConfiguration-%d"
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.1.41"
    private_ip_address_version   = "IPv4"
    primary                      = false
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, template, rInt, rInt, rInt, rInt)
}

func testAccAzureRMPrivateLinkService_template(rInt int, location string) string {
	return fmt.Sprintf(`
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
`, rInt, location, rInt, rInt, rInt)
}
