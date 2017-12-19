package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func testCheckAzureRMApplicationGatewayExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %q", name)
		}

		ApplicationGatewayName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Gateway: %q", ApplicationGatewayName)
		}

		conn := testAccProvider.Meta().(*ArmClient).applicationGatewayClient

		resp, err := conn.Get(resourceGroup, ApplicationGatewayName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: App Gateway %q (resource group: %q) does not exist", ApplicationGatewayName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on ApplicationGatewayClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMApplicationGatewayDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).applicationGatewayClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_application_gateway" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("App Gateway still exists:\n%#v", resp.ApplicationGatewayPropertiesFormat)
	}

	return nil
}

func TestAccAzureRMApplicationGateway_basicPrivateHttp(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationGateway_basicPrivateHttp(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_basicPublicHttp(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationGateway_basicPublicHttp(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_basicHttpUpdate(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()
	location := testLocation()
	firstConfig := testAccAzureRMApplicationGateway_basicPrivateHttp(ri, location)
	secondConfig := testAccAzureRMApplicationGateway_basicPublicHttp(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: firstConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
				),
			},
			{
				Config: secondConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_basicPrivateHttpWafDetection(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationGateway_basicPrivateHttpWaf(ri, testLocation(), true, "Detection")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.firewall_mode", "Detection"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_type", "OWASP"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_version", "3.0"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_basicPrivateHttpWafPrevention(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationGateway_basicPrivateHttpWaf(ri, testLocation(), true, "Prevention")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.firewall_mode", "Prevention"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_type", "OWASP"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_version", "3.0"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_basicPrivateHttpWafUpdate(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()
	location := testLocation()
	firstConfig := testAccAzureRMApplicationGateway_basicPrivateHttp(ri, location)
	secondConfig := testAccAzureRMApplicationGateway_basicPrivateHttpWaf(ri, location, true, "Detection")
	thirdConfig := testAccAzureRMApplicationGateway_basicPrivateHttpWaf(ri, location, true, "Prevention")
	fourthConfig := testAccAzureRMApplicationGateway_basicPrivateHttpWaf(ri, location, false, "Prevention")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: firstConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.%", "0"),
				),
			},
			{
				Config: secondConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.firewall_mode", "Detection"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_type", "OWASP"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_version", "3.0"),
				),
			},
			{
				Config: thirdConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.firewall_mode", "Prevention"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_type", "OWASP"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_version", "3.0"),
				),
			},
			{
				Config: fourthConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.firewall_mode", "Prevention"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_type", "OWASP"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_version", "3.0"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_basicPublicHttpWafDetection(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationGateway_basicPublicHttpWaf(ri, testLocation(), true, "Detection")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.firewall_mode", "Detection"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_type", "OWASP"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_version", "3.0"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_basicPublicHttpWafPrevention(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationGateway_basicPublicHttpWaf(ri, testLocation(), true, "Prevention")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.firewall_mode", "Prevention"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_type", "OWASP"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_version", "3.0"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_basicPublicHttpWafUpdate(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()
	location := testLocation()
	firstConfig := testAccAzureRMApplicationGateway_basicPublicHttp(ri, location)
	secondConfig := testAccAzureRMApplicationGateway_basicPublicHttpWaf(ri, location, true, "Detection")
	thirdConfig := testAccAzureRMApplicationGateway_basicPublicHttpWaf(ri, location, true, "Prevention")
	fourthConfig := testAccAzureRMApplicationGateway_basicPublicHttpWaf(ri, location, false, "Prevention")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: firstConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.%", "0"),
				),
			},
			{
				Config: secondConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.firewall_mode", "Detection"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_type", "OWASP"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_version", "3.0"),
				),
			},
			{
				Config: thirdConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.firewall_mode", "Prevention"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_type", "OWASP"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_version", "3.0"),
				),
			},
			{
				Config: fourthConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.firewall_mode", "Prevention"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_type", "OWASP"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_version", "3.0"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_tags(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationGateway_tags(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Source", "AcceptanceTest"),
				),
			},
		},
	})
}

func testAccAzureRMApplicationGateway_basicPrivateHttp(rInt int, location string) string {
	return fmt.Sprintf(`
variable "frontend_ip_configuration_name" {
  default = "appGatewayFrontendIP"
}

variable "frontend_port_name" {
  default = "appGatewayFrontendPort"
}

variable "http_listener_name" {
  default = "appGatewayHttpListener"
}

variable "backend_address_pool_name" {
  default = "appGatewayBackendPool"
}

variable "backend_http_settings_name" {
  default = "appGatewayBackendHttpSettings"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "internal"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_ip_configuration {
    name      = "${var.frontend_ip_configuration_name}"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${var.frontend_port_name}"
    port = 80
  }

  backend_address_pool {
    name = "${var.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${var.backend_http_settings_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
  }

  http_listener {
    name                           = "${var.http_listener_name}"
    frontend_ip_configuration_name = "${var.frontend_ip_configuration_name}"
    frontend_port_name             = "${var.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "rule1"
    rule_type                  = "Basic"
    http_listener_name         = "${var.http_listener_name}"
    backend_address_pool_name  = "${var.backend_address_pool_name}"
    backend_http_settings_name = "${var.backend_http_settings_name}"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMApplicationGateway_tags(rInt int, location string) string {
	return fmt.Sprintf(`
variable "frontend_ip_configuration_name" {
  default = "appGatewayFrontendIP"
}

variable "frontend_port_name" {
  default = "appGatewayFrontendPort"
}

variable "http_listener_name" {
  default = "appGatewayHttpListener"
}

variable "backend_address_pool_name" {
  default = "appGatewayBackendPool"
}

variable "backend_http_settings_name" {
  default = "appGatewayBackendHttpSettings"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "internal"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_ip_configuration {
    name      = "${var.frontend_ip_configuration_name}"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${var.frontend_port_name}"
    port = 80
  }

  backend_address_pool {
    name = "${var.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${var.backend_http_settings_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
  }

  http_listener {
    name                           = "${var.http_listener_name}"
    frontend_ip_configuration_name = "${var.frontend_ip_configuration_name}"
    frontend_port_name             = "${var.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "rule1"
    rule_type                  = "Basic"
    http_listener_name         = "${var.http_listener_name}"
    backend_address_pool_name  = "${var.backend_address_pool_name}"
    backend_http_settings_name = "${var.backend_http_settings_name}"
  }

  tags {
    Source = "AcceptanceTest"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMApplicationGateway_basicPrivateHttpWaf(rInt int, location string, firewallEnabled bool, firewallMode string) string {
	return fmt.Sprintf(`
variable "frontend_ip_configuration_name" {
  default = "appGatewayFrontendIP"
}

variable "frontend_port_name" {
  default = "appGatewayFrontendPort"
}

variable "http_listener_name" {
  default = "appGatewayHttpListener"
}

variable "backend_address_pool_name" {
  default = "appGatewayBackendPool"
}

variable "backend_http_settings_name" {
  default = "appGatewayBackendHttpSettings"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "WAF_Medium"
    tier     = "WAF"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "internal"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_ip_configuration {
    name      = "${var.frontend_ip_configuration_name}"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${var.frontend_port_name}"
    port = 80
  }

  backend_address_pool {
    name = "${var.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${var.backend_http_settings_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
  }

  http_listener {
    name                           = "${var.http_listener_name}"
    frontend_ip_configuration_name = "${var.frontend_ip_configuration_name}"
    frontend_port_name             = "${var.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "rule1"
    rule_type                  = "Basic"
    http_listener_name         = "${var.http_listener_name}"
    backend_address_pool_name  = "${var.backend_address_pool_name}"
    backend_http_settings_name = "${var.backend_http_settings_name}"
  }

  waf_configuration {
    enabled          = %t
    firewall_mode    = "%s"
    rule_set_type    = "OWASP"
    rule_set_version = "3.0"
  }
}
`, rInt, location, rInt, rInt, firewallEnabled, firewallMode)
}

func testAccAzureRMApplicationGateway_basicPublicHttp(rInt int, location string) string {
	return fmt.Sprintf(`
variable "frontend_ip_configuration_name" {
  default = "appGatewayFrontendIP"
}

variable "frontend_port_name" {
  default = "appGatewayFrontendPort"
}

variable "http_listener_name" {
  default = "appGatewayHttpListener"
}

variable "backend_address_pool_name" {
  default = "appGatewayBackendPool"
}

variable "backend_http_settings_name" {
  default = "appGatewayBackendHttpSettings"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_public_ip" "test" {
  name                         = "acctestpip-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "dynamic"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "internal"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_ip_configuration {
    name                 = "${var.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  frontend_port {
    name = "${var.frontend_port_name}"
    port = 80
  }

  backend_address_pool {
    name = "${var.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${var.backend_http_settings_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
  }

  http_listener {
    name                           = "${var.http_listener_name}"
    frontend_ip_configuration_name = "${var.frontend_ip_configuration_name}"
    frontend_port_name             = "${var.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "rule1"
    rule_type                  = "Basic"
    http_listener_name         = "${var.http_listener_name}"
    backend_address_pool_name  = "${var.backend_address_pool_name}"
    backend_http_settings_name = "${var.backend_http_settings_name}"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMApplicationGateway_basicPublicHttpWaf(rInt int, location string, firewallEnabled bool, firewallMode string) string {
	return fmt.Sprintf(`
variable "frontend_ip_configuration_name" {
  default = "appGatewayFrontendIP"
}

variable "frontend_port_name" {
  default = "appGatewayFrontendPort"
}

variable "http_listener_name" {
  default = "appGatewayHttpListener"
}

variable "backend_address_pool_name" {
  default = "appGatewayBackendPool"
}

variable "backend_http_settings_name" {
  default = "appGatewayBackendHttpSettings"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_public_ip" "test" {
  name                         = "acctestpip-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "dynamic"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "WAF_Medium"
    tier     = "WAF"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "internal"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_ip_configuration {
    name                 = "${var.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  frontend_port {
    name = "${var.frontend_port_name}"
    port = 80
  }

  backend_address_pool {
    name = "${var.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${var.backend_http_settings_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
  }

  http_listener {
    name                           = "${var.http_listener_name}"
    frontend_ip_configuration_name = "${var.frontend_ip_configuration_name}"
    frontend_port_name             = "${var.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "rule1"
    rule_type                  = "Basic"
    http_listener_name         = "${var.http_listener_name}"
    backend_address_pool_name  = "${var.backend_address_pool_name}"
    backend_http_settings_name = "${var.backend_http_settings_name}"
  }

  waf_configuration {
    enabled          = %t
    firewall_mode    = "%s"
    rule_set_type    = "OWASP"
    rule_set_version = "3.0"
  }
}
`, rInt, location, rInt, rInt, rInt, firewallEnabled, firewallMode)
}

func TestAccAzureRMApplicationGateway_privateSpecifiedIP(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationGateway_basicPrivateSpecifiedIP(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "frontend_ip_configuration.0.private_ip_address_allocation", "Static"),
					resource.TestCheckResourceAttr(resourceName, "frontend_ip_configuration.0.private_ip_address", "10.254.0.15"),
				),
			},
		},
	})
}

func testAccAzureRMApplicationGateway_basicPrivateSpecifiedIP(rInt int, location string) string {
	return fmt.Sprintf(`
variable "frontend_ip_configuration_name" {
  default = "appGatewayFrontendIP"
}

variable "frontend_port_name" {
  default = "appGatewayFrontendPort"
}

variable "http_listener_name" {
  default = "appGatewayHttpListener"
}

variable "backend_address_pool_name" {
  default = "appGatewayBackendPool"
}

variable "backend_http_settings_name" {
  default = "appGatewayBackendHttpSettings"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "internal"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_ip_configuration {
    name                          = "${var.frontend_ip_configuration_name}"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Static"
    private_ip_address            = "10.254.0.15"
  }

  frontend_port {
    name = "${var.frontend_port_name}"
    port = 80
  }

  backend_address_pool {
    name = "${var.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${var.backend_http_settings_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
  }

  http_listener {
    name                           = "${var.http_listener_name}"
    frontend_ip_configuration_name = "${var.frontend_ip_configuration_name}"
    frontend_port_name             = "${var.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "rule1"
    rule_type                  = "Basic"
    http_listener_name         = "${var.http_listener_name}"
    backend_address_pool_name  = "${var.backend_address_pool_name}"
    backend_http_settings_name = "${var.backend_http_settings_name}"
  }
}
`, rInt, location, rInt, rInt)
}

func TestAccAzureRMApplicationGateway_basicPrivateMultipleInstances(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationGateway_basicPrivateMultipleInstances(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "5"),
				),
			},
		},
	})
}

func testAccAzureRMApplicationGateway_basicPrivateMultipleInstances(rInt int, location string) string {
	return fmt.Sprintf(`
variable "frontend_ip_configuration_name" {
  default = "appGatewayFrontendIP"
}

variable "frontend_port_name" {
  default = "appGatewayFrontendPort"
}

variable "http_listener_name" {
  default = "appGatewayHttpListener"
}

variable "backend_address_pool_name" {
  default = "appGatewayBackendPool"
}

variable "backend_http_settings_name" {
  default = "appGatewayBackendHttpSettings"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 5
  }

  gateway_ip_configuration {
    name      = "internal"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_ip_configuration {
    name      = "${var.frontend_ip_configuration_name}"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${var.frontend_port_name}"
    port = 80
  }

  backend_address_pool {
    name = "${var.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${var.backend_http_settings_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
  }

  http_listener {
    name                           = "${var.http_listener_name}"
    frontend_ip_configuration_name = "${var.frontend_ip_configuration_name}"
    frontend_port_name             = "${var.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "rule1"
    rule_type                  = "Basic"
    http_listener_name         = "${var.http_listener_name}"
    backend_address_pool_name  = "${var.backend_address_pool_name}"
    backend_http_settings_name = "${var.backend_http_settings_name}"
  }
}
`, rInt, location, rInt, rInt)
}

func TestAccAzureRMApplicationGateway_multipleFrontEndPorts(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationGateway_multipleFrontEndPorts(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "frontend_port.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "frontend_port.0.port", "80"),
					resource.TestCheckResourceAttr(resourceName, "frontend_port.1.port", "3000"),
				),
			},
		},
	})
}

func testAccAzureRMApplicationGateway_multipleFrontEndPorts(rInt int, location string) string {
	return fmt.Sprintf(`
variable "frontend_ip_configuration_name" {
  default = "appGatewayFrontendIP"
}

variable "frontend_port_name" {
  default = "appGatewayFrontendPort"
}

variable "frontend_port_name2" {
  default = "appGatewayFrontendPort2"
}

variable "http_listener_name" {
  default = "appGatewayHttpListener"
}

variable "backend_address_pool_name" {
  default = "appGatewayBackendPool"
}

variable "backend_http_settings_name" {
  default = "appGatewayBackendHttpSettings"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "internal"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_ip_configuration {
    name      = "${var.frontend_ip_configuration_name}"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${var.frontend_port_name}"
    port = 80
  }

  frontend_port {
    name = "${var.frontend_port_name2}"
    port = 3000
  }

  backend_address_pool {
    name = "${var.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${var.backend_http_settings_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
  }

  http_listener {
    name                           = "${var.http_listener_name}"
    frontend_ip_configuration_name = "${var.frontend_ip_configuration_name}"
    frontend_port_name             = "${var.frontend_port_name}"
    protocol                       = "Http"
  }

  http_listener {
    name                           = "${var.http_listener_name}2"
    frontend_ip_configuration_name = "${var.frontend_ip_configuration_name}"
    frontend_port_name             = "${var.frontend_port_name2}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "rule1"
    rule_type                  = "Basic"
    http_listener_name         = "${var.http_listener_name}"
    backend_address_pool_name  = "${var.backend_address_pool_name}"
    backend_http_settings_name = "${var.backend_http_settings_name}"
  }

  request_routing_rule {
    name                       = "rule2"
    rule_type                  = "Basic"
    http_listener_name         = "${var.http_listener_name}2"
    backend_address_pool_name  = "${var.backend_address_pool_name}"
    backend_http_settings_name = "${var.backend_http_settings_name}"
  }
}
`, rInt, location, rInt, rInt)
}

func TestAccAzureRMApplicationGateway_multipleListeners(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationGateway_multipleListeners(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "frontend_port.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "frontend_port.0.port", "80"),
					resource.TestCheckResourceAttr(resourceName, "frontend_port.1.port", "3000"),
				),
			},
		},
	})
}

func testAccAzureRMApplicationGateway_multipleListeners(rInt int, location string) string {
	return fmt.Sprintf(`
variable "frontend_ip_configuration_name" {
  default = "appGatewayFrontendIP"
}

variable "frontend_port_name" {
  default = "appGatewayFrontendPort"
}

variable "frontend_port_name2" {
  default = "appGatewayFrontendPort2"
}

variable "http_listener_name" {
  default = "appGatewayHttpListener"
}

variable "http_listener_name2" {
  default = "appGatewayHttpListener2"
}

variable "backend_address_pool_name" {
  default = "appGatewayBackendPool"
}

variable "backend_http_settings_name" {
  default = "appGatewayBackendHttpSettings"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "internal"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_ip_configuration {
    name      = "${var.frontend_ip_configuration_name}"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${var.frontend_port_name}"
    port = 80
  }

  frontend_port {
    name = "${var.frontend_port_name2}"
    port = 3000
  }

  backend_address_pool {
    name = "${var.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${var.backend_http_settings_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
  }

  http_listener {
    name                           = "${var.http_listener_name}"
    frontend_ip_configuration_name = "${var.frontend_ip_configuration_name}"
    frontend_port_name             = "${var.frontend_port_name}"
    protocol                       = "Http"
  }

  http_listener {
    name                           = "${var.http_listener_name2}"
    frontend_ip_configuration_name = "${var.frontend_ip_configuration_name}"
    frontend_port_name             = "${var.frontend_port_name2}"
    protocol                       = "Http"
    host_name                      = "www.terraform.io"
    require_sni                    = false
  }

  request_routing_rule {
    name                       = "rule1"
    rule_type                  = "Basic"
    http_listener_name         = "${var.http_listener_name}"
    backend_address_pool_name  = "${var.backend_address_pool_name}"
    backend_http_settings_name = "${var.backend_http_settings_name}"
  }

  request_routing_rule {
    name                       = "rule2"
    rule_type                  = "Basic"
    http_listener_name         = "${var.http_listener_name2}"
    backend_address_pool_name  = "${var.backend_address_pool_name}"
    backend_http_settings_name = "${var.backend_http_settings_name}"
  }
}
`, rInt, location, rInt, rInt)
}

func TestAccAzureRMApplicationGateway_pathBasedRouting(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationGateway_pathBasedRouting(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "url_path_map.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "url_path_map.0.path_rule.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "url_path_map.0.path_rule.0.name", "test1"),
					resource.TestCheckResourceAttr(resourceName, "url_path_map.0.path_rule.0.paths.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "url_path_map.0.path_rule.0.paths.0", "/foo/*"),
					resource.TestCheckResourceAttr(resourceName, "url_path_map.0.path_rule.1.name", "test2"),
					resource.TestCheckResourceAttr(resourceName, "url_path_map.0.path_rule.1.paths.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "url_path_map.0.path_rule.1.paths.0", "/bar"),
				),
			},
		},
	})
}

func testAccAzureRMApplicationGateway_pathBasedRouting(rInt int, location string) string {
	return fmt.Sprintf(`
variable "frontend_ip_configuration_name" {
  default = "appGatewayFrontendIP"
}

variable "frontend_port_name" {
  default = "appGatewayFrontendPort"
}

variable "http_listener_name" {
  default = "appGatewayHttpListener"
}

variable "backend_address_pool_name" {
  default = "appGatewayBackendPool"
}

variable "backend_http_settings_name" {
  default = "appGatewayBackendHttpSettings"
}

variable "url_path_map_name" {
  default = "urlPathMap1"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "internal"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_ip_configuration {
    name      = "${var.frontend_ip_configuration_name}"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${var.frontend_port_name}"
    port = 80
  }

  backend_address_pool {
    name = "${var.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${var.backend_http_settings_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
  }

  http_listener {
    name                           = "${var.http_listener_name}"
    frontend_ip_configuration_name = "${var.frontend_ip_configuration_name}"
    frontend_port_name             = "${var.frontend_port_name}"
    protocol                       = "Http"
  }

  url_path_map {
    name                               = "${var.url_path_map_name}"
    default_backend_address_pool_name  = "${var.backend_address_pool_name}"
    default_backend_http_settings_name = "${var.backend_http_settings_name}"

    path_rule {
      name                       = "test1"
      paths                      = ["/foo/*"]
      backend_address_pool_name  = "${var.backend_address_pool_name}"
      backend_http_settings_name = "${var.backend_http_settings_name}"
    }

    path_rule {
      name                       = "test2"
      paths                      = ["/bar"]
      backend_address_pool_name  = "${var.backend_address_pool_name}"
      backend_http_settings_name = "${var.backend_http_settings_name}"
    }
  }

  request_routing_rule {
    name               = "rule1"
    rule_type          = "PathBasedRouting"
    http_listener_name = "${var.http_listener_name}"
    url_path_map_name  = "${var.url_path_map_name}"
  }
}
`, rInt, location, rInt, rInt)
}

func TestAccAzureRMApplicationGateway_multipleFEIPConfigurations(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()

	// a Standard Application Gateway can only support 1 x Public and 1 x Private Frontend IP Configurations
	config := testAccAzureRMApplicationGateway_multipleFEIPConfigurations(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "frontend_ip_configuration.0.subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "frontend_ip_configuration.1.public_ip_address_id"),
				),
			},
		},
	})
}

func testAccAzureRMApplicationGateway_multipleFEIPConfigurations(rInt int, location string) string {
	return fmt.Sprintf(`
variable "frontend_ip_configuration_name_private" {
  default = "appGatewayFrontendIPPrivate"
}

variable "frontend_ip_configuration_name_public" {
	default = "appGatewayFrontendIPPublic"
}

variable "frontend_port_name" {
  default = "appGatewayFrontendPort"
}

variable "http_listener_name_private" {
  default = "appGatewayHttpListenerPrivate"
}

variable "http_listener_name_public" {
  default = "appGatewayHttpListenerPublic"
}

variable "backend_address_pool_name" {
  default = "appGatewayBackendPool"
}

variable "backend_http_settings_name" {
  default = "appGatewayBackendHttpSettings"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_public_ip" "test" {
  name                         = "acctestpip-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "dynamic"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "internal"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_ip_configuration {
    name      = "${var.frontend_ip_configuration_name_private}"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_ip_configuration {
    name                 = "${var.frontend_ip_configuration_name_public}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  frontend_port {
    name = "${var.frontend_port_name}"
    port = 80
  }

  backend_address_pool {
    name = "${var.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${var.backend_http_settings_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
  }

  http_listener {
    name                           = "${var.http_listener_name_private}"
    frontend_ip_configuration_name = "${var.frontend_ip_configuration_name_private}"
    frontend_port_name             = "${var.frontend_port_name}"
    protocol                       = "Http"
  }

  http_listener {
    name                           = "${var.http_listener_name_public}"
    frontend_ip_configuration_name = "${var.frontend_ip_configuration_name_public}"
    frontend_port_name             = "${var.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "rule1"
    rule_type                  = "Basic"
    http_listener_name         = "${var.http_listener_name_private}"
    backend_address_pool_name  = "${var.backend_address_pool_name}"
    backend_http_settings_name = "${var.backend_http_settings_name}"
  }

  request_routing_rule {
    name                       = "rule2"
    rule_type                  = "Basic"
    http_listener_name         = "${var.http_listener_name_public}"
    backend_address_pool_name  = "${var.backend_address_pool_name}"
    backend_http_settings_name = "${var.backend_http_settings_name}"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

// Multiple Probes?
func TestAccAzureRMApplicationGateway_probe(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMApplicationGateway_probe(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "probe.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "probe.0.name", "probe1"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_probes(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMApplicationGateway_probes(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "probe.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "probe.0.name", "probe1"),
					resource.TestCheckResourceAttr(resourceName, "probe.1.name", "probe2"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_probesUpdate(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMApplicationGateway_probe(ri, location)
	updatedConfig := testAccAzureRMApplicationGateway_probes(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "probe.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "probe.0.name", "FirstProbe"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "probe.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "probe.0.name", "FirstProbe"),
					resource.TestCheckResourceAttr(resourceName, "probe.1.name", "SecondProbe"),
				),
			},
		},
	})
}

func testAccAzureRMApplicationGateway_probe(rInt int, location string) string {
	return fmt.Sprintf(`
variable "frontend_ip_configuration_name" {
  default = "appGatewayFrontendIP"
}

variable "frontend_port_name" {
  default = "appGatewayFrontendPort"
}

variable "http_listener_name" {
  default = "appGatewayHttpListener"
}

variable "backend_address_pool_name" {
  default = "appGatewayBackendPool"
}

variable "backend_http_settings_name" {
  default = "appGatewayBackendHttpSettings1"
}

variable "probe_name" {
  default = "probe1"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "internal"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_ip_configuration {
    name      = "${var.frontend_ip_configuration_name}"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${var.frontend_port_name}"
    port = 80
  }

  backend_address_pool {
    name = "${var.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${var.backend_http_settings_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    probe_name            = "${var.probe_name}"
  }

  http_listener {
    name                           = "${var.http_listener_name}"
    frontend_ip_configuration_name = "${var.frontend_ip_configuration_name}"
    frontend_port_name             = "${var.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "rule1"
    rule_type                  = "Basic"
    http_listener_name         = "${var.http_listener_name}"
    backend_address_pool_name  = "${var.backend_address_pool_name}"
    backend_http_settings_name = "${var.backend_http_settings_name}"
  }

  probe {
    name                = "${var.probe_name}"
    protocol            = "Http"
    host                = "www.terraform.io"
    path                = "/"
    interval            = 30
    timeout             = 10
    unhealthy_threshold = 3
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMApplicationGateway_probes(rInt int, location string) string {
	return fmt.Sprintf(`
variable "frontend_ip_configuration_name" {
  default = "appGatewayFrontendIP"
}

variable "frontend_port_name" {
  default = "appGatewayFrontendPort"
}

variable "http_listener_name" {
  default = "appGatewayHttpListener"
}

variable "backend_address_pool_name" {
  default = "appGatewayBackendPool"
}

variable "backend_http_settings_name_first" {
  default = "appGatewayBackendHttpSettings1"
}

variable "backend_http_settings_name_second" {
  default = "appGatewayBackendHttpSettings2"
}

variable "probe_name_first" {
  default = "probe1"
}

variable "probe_name_second" {
  default = "probe2"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "internal"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_ip_configuration {
    name      = "${var.frontend_ip_configuration_name}"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${var.frontend_port_name}"
    port = 80
  }

  backend_address_pool {
    name = "${var.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${var.backend_http_settings_name_first}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    probe_name            = "${var.probe_name_first}"
  }

  backend_http_settings {
    name                  = "${var.backend_http_settings_name_second}"
    cookie_based_affinity = "Disabled"
    port                  = 3000
    protocol              = "Http"
    probe_name            = "${var.probe_name_second}"
  }

  http_listener {
    name                           = "${var.http_listener_name_first}"
    frontend_ip_configuration_name = "${var.frontend_ip_configuration_name}"
    frontend_port_name             = "${var.frontend_port_name}"
    protocol                       = "Http"
  }

  http_listener {
    name                           = "${var.http_listener_name_second}"
    frontend_ip_configuration_name = "${var.frontend_ip_configuration_name}"
    frontend_port_name             = "${var.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "rule1"
    rule_type                  = "Basic"
    http_listener_name         = "${var.http_listener_name}"
    backend_address_pool_name  = "${var.backend_address_pool_name}"
    backend_http_settings_name = "${var.backend_http_settings_name_first}"
  }

  request_routing_rule {
    name                       = "rule2"
    rule_type                  = "Basic"
    http_listener_name         = "${var.http_listener_name}"
    backend_address_pool_name  = "${var.backend_address_pool_name}"
    backend_http_settings_name = "${var.backend_http_settings_name_second}"
  }

  probe {
    name                = "${var.probe_name_first}"
    protocol            = "Http"
    host                = "www.terraform.io"
    path                = "/"
    interval            = 30
    timeout             = 10
    unhealthy_threshold = 3
  }

  probe {
    name                = "${var.probe_name_second}"
    protocol            = "Http"
    host                = "www.terraform.io"
    path                = "/"
    interval            = 30
    timeout             = 10
    unhealthy_threshold = 3
  }
}
`, rInt, location, rInt, rInt)
}

func TestAccAzureRMApplicationGateway_cookieBasedAffinity(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationGateway_cookieBasedAffinity(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "backend_http_settings.0.cookie_based_affinity", "Enabled"),
				),
			},
		},
	})
}

func testAccAzureRMApplicationGateway_cookieBasedAffinity(rInt int, location string) string {
	return fmt.Sprintf(`
variable "frontend_ip_configuration_name" {
  default = "appGatewayFrontendIP"
}

variable "frontend_port_name" {
  default = "appGatewayFrontendPort"
}

variable "http_listener_name" {
  default = "appGatewayHttpListener"
}

variable "backend_address_pool_name" {
  default = "appGatewayBackendPool"
}

variable "backend_http_settings_name" {
  default = "appGatewayBackendHttpSettings"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "internal"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_ip_configuration {
    name      = "${var.frontend_ip_configuration_name}"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${var.frontend_port_name}"
    port = 80
  }

  backend_address_pool {
    name = "${var.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${var.backend_http_settings_name}"
    cookie_based_affinity = "Enabled"
    port                  = 80
    protocol              = "Http"
  }

  http_listener {
    name                           = "${var.http_listener_name}"
    frontend_ip_configuration_name = "${var.frontend_ip_configuration_name}"
    frontend_port_name             = "${var.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "rule1"
    rule_type                  = "Basic"
    http_listener_name         = "${var.http_listener_name}"
    backend_address_pool_name  = "${var.backend_address_pool_name}"
    backend_http_settings_name = "${var.backend_http_settings_name}"
  }
}
`, rInt, location, rInt, rInt)
}

func TestAccAzureRMApplicationGateway_targetDomainName(t *testing.T) {
	// Pointing to a Domain Name
	t.Skip("Not Yet Implemented")
}
func TestAccAzureRMApplicationGateway_targetIP(t *testing.T) {
	// Pointing to an IP Address
	t.Skip("Not Yet Implemented")
}
func TestAccAzureRMApplicationGateway_targetVM(t *testing.T) {
	// Pointing to a VM?
	t.Skip("Not Yet Implemented")
}

func TestAccAzureRMApplicationGateway_sslCert(t *testing.T) {
	t.Skip("Not Yet Implemented")
}
func TestAccAzureRMApplicationGateway_sslCertUpdate(t *testing.T) {
	t.Skip("Not Yet Implemented")
}
func TestAccAzureRMApplicationGateway_sslCertWithDisabledSSLCyphers(t *testing.T) {
	t.Skip("Not Yet Implemented")
}
func TestAccAzureRMApplicationGateway_sslCertWithProbe(t *testing.T) {
	t.Skip("Not Yet Implemented")
}
func TestAccAzureRMApplicationGateway_sslCertOffloading(t *testing.T) {
	t.Skip("Not Yet Implemented")
}
func TestAccAzureRMApplicationGateway_sslCertWithAuthenticationCertificates(t *testing.T) {
	t.Skip("Not Yet Implemented")
}

// - Import tests for all of the above too
// Redirect Configurations
