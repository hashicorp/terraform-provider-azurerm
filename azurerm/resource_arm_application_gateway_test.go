package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApplicationGateway_basic_base(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					testCheckAzureRMApplicationGatewaySslCertificateAssigned(resourceName, "ssl-1"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_basic_changeSslCert(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					testCheckAzureRMApplicationGatewaySslCertificateAssigned(resourceName, "ssl-1"),
				),
				Destroy: false,
			},
			{
				Config: testAccAzureRMApplicationGateway_basic_changeSslCert(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					testCheckAzureRMApplicationGatewaySslCertificateAssigned(resourceName, "ssl-2"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_basic_authCert(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_basic_authCert(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					testCheckAzureRMApplicationGatewaySslCertificateAssigned(resourceName, "ssl-1"),
					testCheckAzureRMApplicationGatewayAuthenticationCertificateAssigned(resourceName, "auth-1"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_basic_changeAuthCert(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_basic_authCert(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					testCheckAzureRMApplicationGatewaySslCertificateAssigned(resourceName, "ssl-1"),
					testCheckAzureRMApplicationGatewayAuthenticationCertificateAssigned(resourceName, "auth-1"),
				),
				Destroy: false,
			},
			{
				Config: testAccAzureRMApplicationGateway_basic_changeAuthCert(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					testCheckAzureRMApplicationGatewaySslCertificateAssigned(resourceName, "ssl-1"),
					testCheckAzureRMApplicationGatewayAuthenticationCertificateAssigned(resourceName, "auth-2"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_waf(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := acctest.RandInt()

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	gwID := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestrg-%d/providers/Microsoft.Network/applicationGateways/acctestgw-%d",
		subscriptionID, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_waf(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					testCheckAzureRMApplicationGatewaySslCertificateAssigned(resourceName, "ssl-1"),
					resource.TestCheckResourceAttr(resourceName, "id", gwID),
				),
			},
		},
	})
}

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

func testCheckAzureRMApplicationGatewaySslCertificateAssigned(name string, certName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
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

		var certId *string

		for _, cert := range *resp.SslCertificates {
			if *cert.Name == certName {
				certId = cert.ID
			}
		}

		if certId == nil {
			return fmt.Errorf("Bad: SSL certificate not found: %q", certName)
		}

		for _, listener := range *resp.HTTPListeners {
			if listener.SslCertificate != nil && *listener.SslCertificate.ID == *certId {
				return nil
			}
		}

		return fmt.Errorf("Bad: SSL certificate not assigned to a listener: %s", certName)
	}
}

func testCheckAzureRMApplicationGatewayAuthenticationCertificateAssigned(name string, certName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
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

		var certId *string

		for _, cert := range *resp.AuthenticationCertificates {
			if *cert.Name == certName {
				certId = cert.ID
			}
		}

		if certId == nil {
			return fmt.Errorf("Bad: Authentication certificate not found: %q", certName)
		}

		for _, backendHttpSettings := range *resp.BackendHTTPSettingsCollection {
			if backendHttpSettings.AuthenticationCertificates != nil {
				for _, authCert := range *backendHttpSettings.AuthenticationCertificates {
					if *authCert.ID == *certId {
						return nil
					}
				}
			}
		}

		return fmt.Errorf("Bad: Authentication certificate not assigned: %q", certName)
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
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return fmt.Errorf("App Gateway still exists:\n%#v", resp.ApplicationGatewayPropertiesFormat)
	}

	return nil
}

func testAccAzureRMApplicationGateway_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "subnet-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_public_ip" "test" {
  name                         = "acctest-pubip-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "dynamic"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestgw-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "Standard_Medium"
    tier     = "Standard"
    capacity = 1
  }

  disabled_ssl_protocols = [
    "TLSv1_0",
  ]

  gateway_ip_configuration {
    # id = computed
    name      = "gw-ip-config1"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_ip_configuration {
    # id = computed
    name                 = "ip-config-public"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  frontend_ip_configuration {
    # id = computed
    name      = "ip-config-private"
    subnet_id = "${azurerm_subnet.test.id}"

    # private_ip_address = computed
    private_ip_address_allocation = "Dynamic"
  }

  frontend_port {
    # id = computed
    name = "port-8080"
    port = 8080
  }

  backend_address_pool {
    # id = computed
    name = "pool-1"

    fqdn_list = [
      "terraform.io",
    ]
  }

  backend_http_settings {
    # id = computed
    name                  = "backend-http-1"
    port                  = 8010
    protocol              = "Https"
    cookie_based_affinity = "Enabled"
    request_timeout       = 30

    # probe_id = computed
    probe_name = "probe-1"
  }

  http_listener {
    # id = computed
    name = "listener-1"

    # frontend_ip_configuration_id = computed
    frontend_ip_configuration_name = "ip-config-public"

    # frontend_ip_port_id = computed
    frontend_port_name = "port-8080"
    protocol           = "Http"
  }

  http_listener {
    name                           = "listener-2"
    frontend_ip_configuration_name = "ip-config-public"
    frontend_port_name             = "port-8080"
    protocol                       = "Https"

    # ssl_certificate_id = computed
    ssl_certificate_name = "ssl-1"
    host_name            = "terraform.io"
    require_sni          = true
  }

  probe {
    # id = computed
    name                = "probe-1"
    protocol            = "Https"
    path                = "/test"
    host                = "azure.com"
    timeout             = 120
    interval            = 300
    unhealthy_threshold = 8
  }

  url_path_map {
    # id = computed
    name                               = "path-map-1"
    default_backend_address_pool_name  = "pool-1"
    default_backend_http_settings_name = "backend-http-1"

    path_rule {
      # id = computed
      name                       = "path-rule-1"
      backend_address_pool_name  = "pool-1"
      backend_http_settings_name = "backend-http-1"

      paths = [
        "/test",
      ]
    }
  }

  request_routing_rule {
    # id = computed
    name      = "rule-basic-1"
    rule_type = "Basic"

    # http_listener_id = computed
    http_listener_name = "listener-1"

    # backend_address_pool_id = computed
    backend_address_pool_name = "pool-1"

    # backend_http_settings_id = computed
    backend_http_settings_name = "backend-http-1"
  }

  request_routing_rule {
    # id = computed
    name              = "rule-path-1"
    rule_type         = "PathBasedRouting"
    url_path_map_name = "path-map-1"

    # http_listener_id = computed
    http_listener_name = "listener-2"
  }

  ssl_certificate {
    # id = computed
    name     = "ssl-1"
    data     = "${file("resource_arm_application_gateway_test.pfx")}"
    password = "terraform"
  }

  tags {
    environment = "tf01"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMApplicationGateway_basic_changeSslCert(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "subnet-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_public_ip" "test" {
  name                         = "acctest-pubip-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "dynamic"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestgw-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "Standard_Medium"
    tier     = "Standard"
    capacity = 1
  }

  disabled_ssl_protocols = [
    "TLSv1_0",
  ]

  gateway_ip_configuration {
    # id = computed
    name      = "gw-ip-config1"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_ip_configuration {
    # id = computed
    name                 = "ip-config-public"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  frontend_ip_configuration {
    # id = computed
    name      = "ip-config-private"
    subnet_id = "${azurerm_subnet.test.id}"

    # private_ip_address = computed
    private_ip_address_allocation = "Dynamic"
  }

  frontend_port {
    # id = computed
    name = "port-8080"
    port = 8080
  }

  backend_address_pool {
    # id = computed
    name = "pool-1"

    fqdn_list = [
      "terraform.io",
    ]
  }

  backend_http_settings {
    # id = computed
    name                  = "backend-http-1"
    port                  = 8010
    protocol              = "Https"
    cookie_based_affinity = "Enabled"
    request_timeout       = 30

    # probe_id = computed
    probe_name = "probe-1"
  }

  http_listener {
    # id = computed
    name = "listener-1"

    # frontend_ip_configuration_id = computed
    frontend_ip_configuration_name = "ip-config-public"

    # frontend_ip_port_id = computed
    frontend_port_name = "port-8080"
    protocol           = "Http"
  }

  http_listener {
    name                           = "listener-2"
    frontend_ip_configuration_name = "ip-config-public"
    frontend_port_name             = "port-8080"
    protocol                       = "Https"

    # ssl_certificate_id = computed
    ssl_certificate_name = "ssl-2"
    host_name            = "terraform.io"
    require_sni          = true
  }

  probe {
    # id = computed
    name                = "probe-1"
    protocol            = "Https"
    path                = "/test"
    host                = "azure.com"
    timeout             = 120
    interval            = 300
    unhealthy_threshold = 8
  }

  url_path_map {
    # id = computed
    name                               = "path-map-1"
    default_backend_address_pool_name  = "pool-1"
    default_backend_http_settings_name = "backend-http-1"

    path_rule {
      # id = computed
      name                       = "path-rule-1"
      backend_address_pool_name  = "pool-1"
      backend_http_settings_name = "backend-http-1"

      paths = [
        "/test",
      ]
    }
  }

  request_routing_rule {
    # id = computed
    name      = "rule-basic-1"
    rule_type = "Basic"

    # http_listener_id = computed
    http_listener_name = "listener-1"

    # backend_address_pool_id = computed
    backend_address_pool_name = "pool-1"

    # backend_http_settings_id = computed
    backend_http_settings_name = "backend-http-1"
  }

  request_routing_rule {
    # id = computed
    name              = "rule-path-1"
    rule_type         = "PathBasedRouting"
    url_path_map_name = "path-map-1"

    # http_listener_id = computed
    http_listener_name = "listener-2"
  }

  ssl_certificate {
    # id = computed
    name     = "ssl-2"
    data     = "${file("resource_arm_application_gateway_test.pfx")}"
    password = "terraform"
  }

  tags {
    environment = "tf01"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMApplicationGateway_basic_authCert(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "subnet-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_public_ip" "test" {
  name                         = "acctest-pubip-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "dynamic"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestgw-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "Standard_Medium"
    tier     = "Standard"
    capacity = 1
  }

  disabled_ssl_protocols = [
    "TLSv1_0",
  ]

  gateway_ip_configuration {
    # id = computed
    name      = "gw-ip-config1"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_ip_configuration {
    # id = computed
    name                 = "ip-config-public"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  frontend_ip_configuration {
    # id = computed
    name      = "ip-config-private"
    subnet_id = "${azurerm_subnet.test.id}"

    # private_ip_address = computed
    private_ip_address_allocation = "Dynamic"
  }

  frontend_port {
    # id = computed
    name = "port-8080"
    port = 8080
  }

  backend_address_pool {
    # id = computed
    name = "pool-1"

    fqdn_list = [
      "terraform.io",
    ]
  }

  backend_http_settings {
    # id = computed
    name                  = "backend-http-1"
    port                  = 8010
    protocol              = "Http"
    cookie_based_affinity = "Enabled"
    request_timeout       = 30
  }

  backend_http_settings {
    # id = computed
    name                  = "backend-http-2"
    port                  = 8011
    protocol              = "Https"
    cookie_based_affinity = "Enabled"
    request_timeout       = 30

    authentication_certificate {
      name = "auth-1"
    }

    # probe_id = computed
    probe_name = "probe-1"
  }

  http_listener {
    # id = computed
    name = "listener-1"

    # frontend_ip_configuration_id = computed
    frontend_ip_configuration_name = "ip-config-public"

    # frontend_ip_port_id = computed
    frontend_port_name = "port-8080"
    protocol           = "Http"
  }

  http_listener {
    name                           = "listener-2"
    frontend_ip_configuration_name = "ip-config-public"
    frontend_port_name             = "port-8080"
    protocol                       = "Https"

    # ssl_certificate_id = computed
    ssl_certificate_name = "ssl-1"
    host_name            = "terraform.io"
    require_sni          = true
  }

  probe {
    # id = computed
    name                = "probe-1"
    protocol            = "Https"
    path                = "/test"
    host                = "azure.com"
    timeout             = 120
    interval            = 300
    unhealthy_threshold = 8
  }

  url_path_map {
    # id = computed
    name                               = "path-map-1"
    default_backend_address_pool_name  = "pool-1"
    default_backend_http_settings_name = "backend-http-1"

    path_rule {
      # id = computed
      name                       = "path-rule-1"
      backend_address_pool_name  = "pool-1"
      backend_http_settings_name = "backend-http-1"

      paths = [
        "/test",
      ]
    }
  }

  request_routing_rule {
    # id = computed
    name      = "rule-basic-1"
    rule_type = "Basic"

    # http_listener_id = computed
    http_listener_name = "listener-1"

    # backend_address_pool_id = computed
    backend_address_pool_name = "pool-1"

    # backend_http_settings_id = computed
    backend_http_settings_name = "backend-http-1"
  }

  request_routing_rule {
    # id = computed
    name              = "rule-path-1"
    rule_type         = "PathBasedRouting"
    url_path_map_name = "path-map-1"

    # http_listener_id = computed
    http_listener_name = "listener-2"
  }

  authentication_certificate {
    name = "auth-1"
    data = "${file("resource_arm_application_gateway_test.cer")}"
  }

  ssl_certificate {
    # id = computed
    name     = "ssl-1"
    data     = "${file("resource_arm_application_gateway_test.pfx")}"
    password = "terraform"
  }

  tags {
    environment = "tf01"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMApplicationGateway_basic_changeAuthCert(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "subnet-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_public_ip" "test" {
  name                         = "acctest-pubip-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "dynamic"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestgw-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "Standard_Medium"
    tier     = "Standard"
    capacity = 1
  }

  disabled_ssl_protocols = [
    "TLSv1_0",
  ]

  gateway_ip_configuration {
    # id = computed
    name      = "gw-ip-config1"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_ip_configuration {
    # id = computed
    name                 = "ip-config-public"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  frontend_ip_configuration {
    # id = computed
    name      = "ip-config-private"
    subnet_id = "${azurerm_subnet.test.id}"

    # private_ip_address = computed
    private_ip_address_allocation = "Dynamic"
  }

  frontend_port {
    # id = computed
    name = "port-8080"
    port = 8080
  }

  backend_address_pool {
    # id = computed
    name = "pool-1"

    fqdn_list = [
      "terraform.io",
    ]
  }

  backend_http_settings {
    # id = computed
    name                  = "backend-http-1"
    port                  = 8010
    protocol              = "Http"
    cookie_based_affinity = "Enabled"
    request_timeout       = 30
  }

  backend_http_settings {
    # id = computed
    name                  = "backend-http-2"
    port                  = 8011
    protocol              = "Https"
    cookie_based_affinity = "Enabled"
    request_timeout       = 30

    authentication_certificate {
      name = "auth-2"
    }

    # probe_id = computed
    probe_name = "probe-1"
  }

  http_listener {
    # id = computed
    name = "listener-1"

    # frontend_ip_configuration_id = computed
    frontend_ip_configuration_name = "ip-config-public"

    # frontend_ip_port_id = computed
    frontend_port_name = "port-8080"
    protocol           = "Http"
  }

  http_listener {
    name                           = "listener-2"
    frontend_ip_configuration_name = "ip-config-public"
    frontend_port_name             = "port-8080"
    protocol                       = "Https"

    # ssl_certificate_id = computed
    ssl_certificate_name = "ssl-1"
    host_name            = "terraform.io"
    require_sni          = true
  }

  probe {
    # id = computed
    name                = "probe-1"
    protocol            = "Https"
    path                = "/test"
    host                = "azure.com"
    timeout             = 120
    interval            = 300
    unhealthy_threshold = 8
  }

  url_path_map {
    # id = computed
    name                               = "path-map-1"
    default_backend_address_pool_name  = "pool-1"
    default_backend_http_settings_name = "backend-http-1"

    path_rule {
      # id = computed
      name                       = "path-rule-1"
      backend_address_pool_name  = "pool-1"
      backend_http_settings_name = "backend-http-1"

      paths = [
        "/test",
      ]
    }
  }

  request_routing_rule {
    # id = computed
    name      = "rule-basic-1"
    rule_type = "Basic"

    # http_listener_id = computed
    http_listener_name = "listener-1"

    # backend_address_pool_id = computed
    backend_address_pool_name = "pool-1"

    # backend_http_settings_id = computed
    backend_http_settings_name = "backend-http-1"
  }

  request_routing_rule {
    # id = computed
    name              = "rule-path-1"
    rule_type         = "PathBasedRouting"
    url_path_map_name = "path-map-1"

    # http_listener_id = computed
    http_listener_name = "listener-2"
  }

  authentication_certificate {
    name = "auth-2"
    data = "${file("resource_arm_application_gateway_test.cer")}"
  }

  ssl_certificate {
    # id = computed
    name     = "ssl-1"
    data     = "${file("resource_arm_application_gateway_test.pfx")}"
    password = "terraform"
  }

  tags {
    environment = "tf01"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMApplicationGateway_waf(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "subnet-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_public_ip" "test" {
  name                         = "acctest-pubip-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "dynamic"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestgw-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "WAF_Medium"
    tier     = "WAF"
    capacity = 1
  }

  disabled_ssl_protocols = [
    "TLSv1_0",
  ]

  waf_configuration {
    enabled = "true"
    firewall_mode = "Detection"
    rule_set_type = "OWASP"
    rule_set_version = "3.0"
  }

  gateway_ip_configuration {
    # id = computed
    name      = "gw-ip-config1"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_ip_configuration {
    # id = computed
    name                 = "ip-config-public"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  frontend_ip_configuration {
    # id = computed
    name      = "ip-config-private"
    subnet_id = "${azurerm_subnet.test.id}"

    # private_ip_address = computed
    private_ip_address_allocation = "Dynamic"
  }

  frontend_port {
    # id = computed
    name = "port-8080"
    port = 8080
  }

  backend_address_pool {
    # id = computed
    name = "pool-1"

    fqdn_list = [
      "terraform.io",
    ]
  }

  backend_http_settings {
    # id = computed
    name                  = "backend-http-1"
    port                  = 8010
    protocol              = "Https"
    cookie_based_affinity = "Enabled"
    request_timeout       = 30

    # probe_id = computed
    probe_name = "probe-1"
  }

  http_listener {
    # id = computed
    name = "listener-1"

    # frontend_ip_configuration_id = computed
    frontend_ip_configuration_name = "ip-config-public"

    # frontend_ip_port_id = computed
    frontend_port_name = "port-8080"
    protocol           = "Http"
  }

  http_listener {
    name                           = "listener-2"
    frontend_ip_configuration_name = "ip-config-public"
    frontend_port_name             = "port-8080"
    protocol                       = "Https"

    # ssl_certificate_id = computed
    ssl_certificate_name = "ssl-1"
    host_name            = "terraform.io"
    require_sni          = true
  }

  probe {
    # id = computed
    name                = "probe-1"
    protocol            = "Https"
    path                = "/test"
    host                = "azure.com"
    timeout             = 120
    interval            = 300
    unhealthy_threshold = 8
  }

  url_path_map {
    # id = computed
    name                               = "path-map-1"
    default_backend_address_pool_name  = "pool-1"
    default_backend_http_settings_name = "backend-http-1"

    path_rule {
      # id = computed
      name                       = "path-rule-1"
      backend_address_pool_name  = "pool-1"
      backend_http_settings_name = "backend-http-1"

      paths = [
        "/test",
      ]
    }
  }

  request_routing_rule {
    # id = computed
    name      = "rule-basic-1"
    rule_type = "Basic"

    # http_listener_id = computed
    http_listener_name = "listener-1"

    # backend_address_pool_id = computed
    backend_address_pool_name = "pool-1"

    # backend_http_settings_id = computed
    backend_http_settings_name = "backend-http-1"
  }

  request_routing_rule {
    # id = computed
    name              = "rule-path-1"
    rule_type         = "PathBasedRouting"
    url_path_map_name = "path-map-1"

    # http_listener_id = computed
    http_listener_name = "listener-2"
  }

  ssl_certificate {
    # id = computed
    name     = "ssl-1"
    data     = "${file("resource_arm_application_gateway_test.pfx")}"
    password = "terraform"
  }

  tags {
    environment = "tf01"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt)
}
