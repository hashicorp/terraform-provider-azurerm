package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApplicationGateway_basic(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_Small"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "2"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.#", "0"),
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

func TestAccAzureRMApplicationGateway_http2(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_http2(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_http2", "true"),
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

func TestAccAzureRMApplicationGateway_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMApplicationGateway_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_application_gateway"),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_authCertificate(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_authCertificate(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "authentication_certificate.0.name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// since these are read from the existing state
					"authentication_certificate.0.data",
				},
			},
			{
				Config: testAccAzureRMApplicationGateway_authCertificateUpdated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "authentication_certificate.0.name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// since these are read from the existing state
					"authentication_certificate.0.data",
				},
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_pathBasedRouting(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_pathBasedRouting(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
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

func TestAccAzureRMApplicationGateway_customErrorConfigurations(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_customErrorConfigurations(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
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

func TestAccAzureRMApplicationGateway_probes(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_probes(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
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

func TestAccAzureRMApplicationGateway_probesPickHostNameFromBackendHTTPSettings(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_probesPickHostNameFromBackendHTTPSettings(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "probe.0.pick_host_name_from_backend_http_settings", "true"),
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

func TestAccAzureRMApplicationGateway_settingsPickHostNameFromBackendAddress(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_settingsPickHostNameFromBackendAddress(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "backend_http_settings.0.pick_host_name_from_backend_address", "true"),
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

func TestAccAzureRMApplicationGateway_sslCertificate(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_sslCertificate(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// since these are read from the existing state
					"ssl_certificate.0.data",
					"ssl_certificate.0.password",
				},
			},
			{
				Config: testAccAzureRMApplicationGateway_sslCertificateUpdated(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// since these are read from the existing state
					"ssl_certificate.0.data",
					"ssl_certificate.0.password",
				},
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_webApplicationFirewall(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_webApplicationFirewall(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "WAF_Medium"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "WAF"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.firewall_mode", "Detection"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_type", "OWASP"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_version", "3.0"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.file_upload_limit_mb", "100"),
				),
			},
		},
	})
}

func testCheckAzureRMApplicationGatewayExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		gatewayName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Application Gateway: %q", gatewayName)
		}

		client := testAccProvider.Meta().(*ArmClient).applicationGatewayClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, gatewayName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Application Gateway %q (resource group: %q) does not exist", gatewayName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on applicationGatewayClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMApplicationGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).applicationGatewayClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_application_gateway" {
			continue
		}

		gatewayName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, gatewayName)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Application Gateway still exists:\n%#v", resp.ApplicationGatewayPropertiesFormat)
	}

	return nil
}

func testAccAzureRMApplicationGateway_basic(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "${local.request_routing_rule_name}"
    rule_type                  = "Basic"
    http_listener_name         = "${local.listener_name}"
    backend_address_pool_name  = "${local.backend_address_pool_name}"
    backend_http_settings_name = "${local.http_setting_name}"
  }
}
`, template, rInt)
}

func testAccAzureRMApplicationGateway_http2(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  enable_http2        = true

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "${local.request_routing_rule_name}"
    rule_type                  = "Basic"
    http_listener_name         = "${local.listener_name}"
    backend_address_pool_name  = "${local.backend_address_pool_name}"
    backend_http_settings_name = "${local.http_setting_name}"
  }
}
`, template, rInt)
}

func testAccAzureRMApplicationGateway_requiresImport(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_application_gateway" "import" {
  name                = "${azurerm_application_gateway.test.name}"
  resource_group_name = "${azurerm_application_gateway.test.resource_group_name}"
  location            = "${azurerm_application_gateway.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "${local.request_routing_rule_name}"
    rule_type                  = "Basic"
    http_listener_name         = "${local.listener_name}"
    backend_address_pool_name  = "${local.backend_address_pool_name}"
    backend_http_settings_name = "${local.http_setting_name}"
  }
}
`, template)
}

func testAccAzureRMApplicationGateway_authCertificate(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  auth_cert_name                 = "${azurerm_virtual_network.test.name}-auth"
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 443
    protocol              = "Https"
    request_timeout       = 1

    authentication_certificate {
      name = "${local.auth_cert_name}"
    }
  }

  authentication_certificate {
    name = "${local.auth_cert_name}"
    data = "${file("testdata/application_gateway_test.cer")}"
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "${local.request_routing_rule_name}"
    rule_type                  = "Basic"
    http_listener_name         = "${local.listener_name}"
    backend_address_pool_name  = "${local.backend_address_pool_name}"
    backend_http_settings_name = "${local.http_setting_name}"
  }
}
`, template, rInt)
}

func testAccAzureRMApplicationGateway_authCertificateUpdated(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  auth_cert_name                 = "${azurerm_virtual_network.test.name}-auth2"
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    port                  = 443
    protocol              = "Https"
    request_timeout       = 1

    authentication_certificate {
      name = "${local.auth_cert_name}"
    }
  }

  authentication_certificate {
    name = "${local.auth_cert_name}"
    data = "${file("testdata/application_gateway_test_2.crt")}"
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "${local.request_routing_rule_name}"
    rule_type                  = "Basic"
    http_listener_name         = "${local.listener_name}"
    backend_address_pool_name  = "${local.backend_address_pool_name}"
    backend_http_settings_name = "${local.http_setting_name}"
  }
}
`, template, rInt)
}

func testAccAzureRMApplicationGateway_pathBasedRouting(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  path_rule_name                 = "${azurerm_virtual_network.test.name}-pathrule1"
  url_path_map_name              = "${azurerm_virtual_network.test.name}-urlpath1"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name               = "${local.request_routing_rule_name}"
    rule_type          = "PathBasedRouting"
    url_path_map_name  = "${local.url_path_map_name}"
    http_listener_name = "${local.listener_name}"
  }

  url_path_map {
    name                               = "${local.url_path_map_name}"
    default_backend_address_pool_name  = "${local.backend_address_pool_name}"
    default_backend_http_settings_name = "${local.http_setting_name}"

    path_rule {
      name                       = "${local.path_rule_name}"
      backend_address_pool_name  = "${local.backend_address_pool_name}"
      backend_http_settings_name = "${local.http_setting_name}"

      paths = [
        "/test",
      ]
    }
  }
}
`, template, rInt)
}

func testAccAzureRMApplicationGateway_probes(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  probe1_name                    = "${azurerm_virtual_network.test.name}-probe1"
  probe2_name                    = "${azurerm_virtual_network.test.name}-probe2"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    probe_name            = "${local.probe1_name}"
    protocol              = "Http"
    request_timeout       = 1
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}-2"
    cookie_based_affinity = "Disabled"
    port                  = 8080
    probe_name            = "${local.probe2_name}"
    protocol              = "Http"
    request_timeout       = 1
  }

  probe {
    name                = "${local.probe1_name}"
    protocol            = "Http"
    path                = "/test"
    host                = "azure.com"
    timeout             = 120
    interval            = 300
    unhealthy_threshold = 8
  }

  probe {
    name                = "${local.probe2_name}"
    protocol            = "Http"
    path                = "/other"
    host                = "azure.com"
    timeout             = 120
    interval            = 300
    unhealthy_threshold = 8
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "${local.request_routing_rule_name}"
    rule_type                  = "Basic"
    http_listener_name         = "${local.listener_name}"
    backend_address_pool_name  = "${local.backend_address_pool_name}"
    backend_http_settings_name = "${local.http_setting_name}"
  }
}
`, template, rInt)
}

func testAccAzureRMApplicationGateway_probesPickHostNameFromBackendHTTPSettings(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  probe_name                     = "${azurerm_virtual_network.test.name}-probe"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                                = "${local.http_setting_name}"
    cookie_based_affinity               = "Disabled"
    pick_host_name_from_backend_address = true
    port                                = 80
    probe_name                          = "${local.probe_name}"
    protocol                            = "Http"
    request_timeout                     = 1
  }

  probe {
    name                                      = "${local.probe_name}"
    protocol                                  = "Http"
    path                                      = "/test"
    timeout                                   = 120
    interval                                  = 300
    unhealthy_threshold                       = 8
    pick_host_name_from_backend_http_settings = true
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "${local.request_routing_rule_name}"
    rule_type                  = "Basic"
    http_listener_name         = "${local.listener_name}"
    backend_address_pool_name  = "${local.backend_address_pool_name}"
    backend_http_settings_name = "${local.http_setting_name}"
  }
}
`, template, rInt)
}

func testAccAzureRMApplicationGateway_settingsPickHostNameFromBackendAddress(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                                = "${local.http_setting_name}"
    cookie_based_affinity               = "Disabled"
    pick_host_name_from_backend_address = true
    port                                = 80
    protocol                            = "Http"
    request_timeout                     = 1
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "${local.request_routing_rule_name}"
    rule_type                  = "Basic"
    http_listener_name         = "${local.listener_name}"
    backend_address_pool_name  = "${local.backend_address_pool_name}"
    backend_http_settings_name = "${local.http_setting_name}"
  }
}
`, template, rInt)
}

func testAccAzureRMApplicationGateway_sslCertificate(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  ssl_certificate_name           = "${azurerm_virtual_network.test.name}-ssl1"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 443
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Https"
    ssl_certificate_name           = "${local.ssl_certificate_name}"
  }

  request_routing_rule {
    name                       = "${local.request_routing_rule_name}"
    rule_type                  = "Basic"
    http_listener_name         = "${local.listener_name}"
    backend_address_pool_name  = "${local.backend_address_pool_name}"
    backend_http_settings_name = "${local.http_setting_name}"
  }

  ssl_certificate {
    name     = "${local.ssl_certificate_name}"
    data     = "${file("testdata/application_gateway_test.pfx")}"
    password = "terraform"
  }
}
`, template, rInt)
}

func testAccAzureRMApplicationGateway_sslCertificateUpdated(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  ssl_certificate_name           = "${azurerm_virtual_network.test.name}-ssl2"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 443
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Https"
    ssl_certificate_name           = "${local.ssl_certificate_name}"
  }

  request_routing_rule {
    name                       = "${local.request_routing_rule_name}"
    rule_type                  = "Basic"
    http_listener_name         = "${local.listener_name}"
    backend_address_pool_name  = "${local.backend_address_pool_name}"
    backend_http_settings_name = "${local.http_setting_name}"
  }

  ssl_certificate {
    name     = "${local.ssl_certificate_name}"
    data     = "${file("testdata/application_gateway_test_2.pfx")}"
    password = "hello-world"
  }
}
`, template, rInt)
}

func testAccAzureRMApplicationGateway_webApplicationFirewall(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
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

  disabled_ssl_protocols = [
    "TLSv1_0",
  ]

  waf_configuration {
    enabled          = true
    firewall_mode    = "Detection"
    rule_set_type    = "OWASP"
    rule_set_version = "3.0"
    file_upload_limit_mb = 100
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "${local.request_routing_rule_name}"
    rule_type                  = "Basic"
    http_listener_name         = "${local.listener_name}"
    backend_address_pool_name  = "${local.backend_address_pool_name}"
    backend_http_settings_name = "${local.http_setting_name}"
  }
}
`, template, rInt)
}

func testAccAzureRMApplicationGateway_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "subnet-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.0.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-pubip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Dynamic"
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMApplicationGateway_customErrorConfigurations(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  path_rule_name                 = "${azurerm_virtual_network.test.name}-pathrule1"
  url_path_map_name              = "${azurerm_virtual_network.test.name}-urlpath1"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Http"
    custom_error_configuration {
      status_code           = "HTTPStatus403"
      custom_error_page_url = "http://azure.com/error403_listener.html"
    }
    custom_error_configuration {
      status_code           = "HTTPStatus502"
      custom_error_page_url = "http://azure.com/error502_listener.html"
    }
  }

  custom_error_configuration {
    status_code           = "HTTPStatus403"
    custom_error_page_url = "http://azure.com/error.html"
  }

  custom_error_configuration {
    status_code           = "HTTPStatus502"
    custom_error_page_url = "http://azure.com/error.html"
  }
}
`, template, rInt)
}
