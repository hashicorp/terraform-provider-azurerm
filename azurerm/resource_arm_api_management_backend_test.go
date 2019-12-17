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

func TestAccAzureRMApiManagementBackend_basic(t *testing.T) {
	resourceName := "azurerm_api_management_backend.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementBackendDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementBackend_basic(ri, "basic", location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementBackendExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "http"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://acctest"),
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

func TestAccAzureRMApiManagementBackend_allProperties(t *testing.T) {
	resourceName := "azurerm_api_management_backend.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementBackendDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementBackend_allProperties(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementBackendExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "http"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://acctest"),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "resource_id", "https://resourceid"),
					resource.TestCheckResourceAttr(resourceName, "title", "title"),
					resource.TestCheckResourceAttr(resourceName, "credentials.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.authorization.0.parameter", "parameter"),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.authorization.0.scheme", "scheme"),
					resource.TestCheckResourceAttrSet(resourceName, "credentials.0.certificate.0"),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.header.header1", "header1value1,header1value2"),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.header.header2", "header2value1,header2value2"),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.query.query1", "query1value1,query1value2"),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.query.query2", "query2value1,query2value2"),
					resource.TestCheckResourceAttr(resourceName, "proxy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "proxy.0.url", "http://192.168.1.1:8080"),
					resource.TestCheckResourceAttr(resourceName, "proxy.0.username", "username"),
					resource.TestCheckResourceAttr(resourceName, "proxy.0.password", "password"),
					resource.TestCheckResourceAttr(resourceName, "tls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tls.0.validate_certificate_chain", "false"),
					resource.TestCheckResourceAttr(resourceName, "tls.0.validate_certificate_name", "true"),
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

func TestAccAzureRMApiManagementBackend_credentialsNoCertificate(t *testing.T) {
	resourceName := "azurerm_api_management_backend.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementBackendDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementBackend_credentialsNoCertificate(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementBackendExists(resourceName),
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

func TestAccAzureRMApiManagementBackend_update(t *testing.T) {
	resourceName := "azurerm_api_management_backend.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementBackendDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementBackend_basic(ri, "update", location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementBackendExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "http"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://acctest"),
				),
			},
			{
				Config: testAccAzureRMApiManagementBackend_update(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementBackendExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "soap"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://updatedacctest"),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "resource_id", "https://resourceid"),
					resource.TestCheckResourceAttr(resourceName, "proxy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "proxy.0.url", "http://192.168.1.1:8080"),
					resource.TestCheckResourceAttr(resourceName, "proxy.0.username", "username"),
					resource.TestCheckResourceAttr(resourceName, "proxy.0.password", "password"),
					resource.TestCheckResourceAttr(resourceName, "tls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tls.0.validate_certificate_chain", "false"),
					resource.TestCheckResourceAttr(resourceName, "tls.0.validate_certificate_name", "true"),
				),
			},
			{
				Config: testAccAzureRMApiManagementBackend_basic(ri, "update", location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementBackendExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "http"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://acctest"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "resource_id", ""),
					resource.TestCheckResourceAttr(resourceName, "proxy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tls.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMApiManagementBackend_serviceFabric(t *testing.T) {
	resourceName := "azurerm_api_management_backend.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementBackendDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementBackend_serviceFabric(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementBackendExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "http"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://acctest"),
					resource.TestCheckResourceAttr(resourceName, "service_fabric_cluster.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "service_fabric_cluster.0.client_certificate_thumbprint"),
					resource.TestCheckResourceAttr(resourceName, "service_fabric_cluster.0.max_partition_resolution_retries", "5"),
					resource.TestCheckResourceAttr(resourceName, "service_fabric_cluster.0.management_endpoints.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "service_fabric_cluster.0.server_certificate_thumbprints.#", "2"),
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

func TestAccAzureRMApiManagementBackend_disappears(t *testing.T) {
	resourceName := "azurerm_api_management_backend.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementBackendDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementBackend_basic(ri, "disappears", location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementBackendExists(resourceName),
					testCheckAzureRMApiManagementBackendDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMApiManagementBackend_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_api_management_backend.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementBackendDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementBackend_basic(ri, "import", location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementBackendExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMApiManagementBackend_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_api_management_backend"),
			},
		},
	})
}

func testCheckAzureRMApiManagementBackendDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.BackendClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_backend" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := conn.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return nil
	}

	return nil
}

func testCheckAzureRMApiManagementBackendExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.BackendClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := conn.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Backend %q (API Management Service %q / Resource Group: %q) does not exist", name, serviceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on BackendClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMApiManagementBackendDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for backend: %s", name)
		}

		conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.BackendClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := conn.Delete(ctx, resourceGroup, serviceName, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp) {
				return nil
			}
			return fmt.Errorf("Bad: Delete on BackendClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementBackend_basic(rInt int, testName string, location string) string {
	template := testAccAzureRMApiManagementBackend_template(rInt, testName, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_backend" "test" {
  name                = "acctestapi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  protocol            = "http"
  url                 = "https://acctest"
}
`, template, rInt)
}

func testAccAzureRMApiManagementBackend_update(rInt int, location string) string {
	template := testAccAzureRMApiManagementBackend_template(rInt, "update", location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_backend" "test" {
  name                = "acctestapi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  protocol            = "soap"
  url                 = "https://updatedacctest"
  description         = "description"
  resource_id         = "https://resourceid"
  proxy {
    url      = "http://192.168.1.1:8080"
    username = "username"
    password = "password"
  }
  tls {
    validate_certificate_chain = false
    validate_certificate_name  = true
  }
}
`, template, rInt)
}

func testAccAzureRMApiManagementBackend_allProperties(rInt int, location string) string {
	template := testAccAzureRMApiManagementBackend_template(rInt, "all", location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_certificate" "test" {
  name                = "example-cert"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  data                = "${filebase64("testdata/keyvaultcert.pfx")}"
  password            = ""
}

resource "azurerm_api_management_backend" "test" {
  name                = "acctestapi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  protocol            = "http"
  url                 = "https://acctest"
  description         = "description"
  resource_id         = "https://resourceid"
  title               = "title"
  credentials {
    authorization {
      parameter = "parameter"
      scheme    = "scheme"
    }
    certificate = [
      "${azurerm_api_management_certificate.test.thumbprint}",
    ]
    header = {
      header1 = "header1value1,header1value2"
      header2 = "header2value1,header2value2"
    }
    query = {
      query1 = "query1value1,query1value2"
      query2 = "query2value1,query2value2"
    }
  }
  proxy {
    url      = "http://192.168.1.1:8080"
    username = "username"
    password = "password"
  }
  tls {
    validate_certificate_chain = false
    validate_certificate_name  = true
  }
}
`, template, rInt)
}

func testAccAzureRMApiManagementBackend_serviceFabric(rInt int, location string) string {
	template := testAccAzureRMApiManagementBackend_template(rInt, "sf", location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_certificate" "test" {
  name                = "example-cert"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  data                = "${filebase64("testdata/keyvaultcert.pfx")}"
  password            = ""
}

resource "azurerm_api_management_backend" "test" {
  name                = "acctestapi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  protocol            = "http"
  url                 = "https://acctest"
  service_fabric_cluster {
    client_certificate_thumbprint = "${azurerm_api_management_certificate.test.thumbprint}"
    management_endpoints = [
      "https://acctestsf.com",
    ]
    max_partition_resolution_retries = 5
    server_certificate_thumbprints = [
      "thumb1",
      "thumb2",
    ]
  }
}
`, template, rInt)
}

func testAccAzureRMApiManagementBackend_requiresImport(rInt int, location string) string {
	template := testAccAzureRMApiManagementBackend_basic(rInt, "requiresimport", location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_backend" "import" {
  name                = "${azurerm_api_management_backend.test.name}"
  resource_group_name = "${azurerm_api_management_backend.test.resource_group_name}"
  api_management_name = "${azurerm_api_management_backend.test.api_management_name}"
  protocol            = "${azurerm_api_management_backend.test.protocol}"
  url                 = "${azurerm_api_management_backend.test.url}"
}
`, template)
}

func testAccAzureRMApiManagementBackend_template(rInt int, testName string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d-%s"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d-%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku {
    name     = "Developer"
    capacity = 1
  }
}
`, rInt, testName, location, rInt, testName)
}

func testAccAzureRMApiManagementBackend_credentialsNoCertificate(rInt int, location string) string {
	template := testAccAzureRMApiManagementBackend_template(rInt, "all", location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_backend" "test" {
  name                = "acctestapi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  protocol            = "http"
  url                 = "https://acctest"
  description         = "description"
  resource_id         = "https://resourceid"
  title               = "title"
  credentials {
    authorization {
      parameter = "parameter"
      scheme    = "scheme"
    }
    header = {
      header1 = "header1value1,header1value2"
      header2 = "header2value1,header2value2"
    }
    query = {
      query1 = "query1value1,query1value2"
      query2 = "query2value1,query2value2"
    }
  }
  proxy {
    url      = "http://192.168.1.1:8080"
    username = "username"
    password = "password"
  }
  tls {
    validate_certificate_chain = false
    validate_certificate_name  = true
  }
}
`, template, rInt)
}
