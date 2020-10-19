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

func TestAccAzureRMApiManagementBackend_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_backend", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementBackendDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementBackend_basic(data, "basic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementBackendExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "protocol", "http"),
					resource.TestCheckResourceAttr(data.ResourceName, "url", "https://acctest"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementBackend_allProperties(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_backend", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementBackendDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementBackend_allProperties(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementBackendExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "protocol", "http"),
					resource.TestCheckResourceAttr(data.ResourceName, "url", "https://acctest"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "description"),
					resource.TestCheckResourceAttr(data.ResourceName, "resource_id", "https://resourceid"),
					resource.TestCheckResourceAttr(data.ResourceName, "title", "title"),
					resource.TestCheckResourceAttr(data.ResourceName, "credentials.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "credentials.0.authorization.0.parameter", "parameter"),
					resource.TestCheckResourceAttr(data.ResourceName, "credentials.0.authorization.0.scheme", "scheme"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "credentials.0.certificate.0"),
					resource.TestCheckResourceAttr(data.ResourceName, "credentials.0.header.header1", "header1value1,header1value2"),
					resource.TestCheckResourceAttr(data.ResourceName, "credentials.0.header.header2", "header2value1,header2value2"),
					resource.TestCheckResourceAttr(data.ResourceName, "credentials.0.query.query1", "query1value1,query1value2"),
					resource.TestCheckResourceAttr(data.ResourceName, "credentials.0.query.query2", "query2value1,query2value2"),
					resource.TestCheckResourceAttr(data.ResourceName, "proxy.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "proxy.0.url", "http://192.168.1.1:8080"),
					resource.TestCheckResourceAttr(data.ResourceName, "proxy.0.username", "username"),
					resource.TestCheckResourceAttr(data.ResourceName, "proxy.0.password", "password"),
					resource.TestCheckResourceAttr(data.ResourceName, "tls.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tls.0.validate_certificate_chain", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "tls.0.validate_certificate_name", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementBackend_credentialsNoCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_backend", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementBackendDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementBackend_credentialsNoCertificate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementBackendExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementBackend_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_backend", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementBackendDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementBackend_basic(data, "update"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementBackendExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "protocol", "http"),
					resource.TestCheckResourceAttr(data.ResourceName, "url", "https://acctest"),
				),
			},
			{
				Config: testAccAzureRMApiManagementBackend_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementBackendExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "protocol", "soap"),
					resource.TestCheckResourceAttr(data.ResourceName, "url", "https://updatedacctest"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "description"),
					resource.TestCheckResourceAttr(data.ResourceName, "resource_id", "https://resourceid"),
					resource.TestCheckResourceAttr(data.ResourceName, "proxy.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "proxy.0.url", "http://192.168.1.1:8080"),
					resource.TestCheckResourceAttr(data.ResourceName, "proxy.0.username", "username"),
					resource.TestCheckResourceAttr(data.ResourceName, "proxy.0.password", "password"),
					resource.TestCheckResourceAttr(data.ResourceName, "tls.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tls.0.validate_certificate_chain", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "tls.0.validate_certificate_name", "true"),
				),
			},
			{
				Config: testAccAzureRMApiManagementBackend_basic(data, "update"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementBackendExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "protocol", "http"),
					resource.TestCheckResourceAttr(data.ResourceName, "url", "https://acctest"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "resource_id", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "proxy.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "tls.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMApiManagementBackend_serviceFabric(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_backend", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementBackendDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementBackend_serviceFabric(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementBackendExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "service_fabric_cluster.0.client_certificate_thumbprint"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementBackend_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_backend", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementBackendDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementBackend_basic(data, "disappears"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementBackendExists(data.ResourceName),
					testCheckAzureRMApiManagementBackendDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMApiManagementBackend_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_backend", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementBackendDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementBackend_basic(data, "import"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementBackendExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementBackend_requiresImport),
		},
	})
}

func testCheckAzureRMApiManagementBackendDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.BackendClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_backend" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

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
		conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.BackendClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

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
		conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.BackendClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMApiManagementBackend_basic(data acceptance.TestData, testName string) string {
	template := testAccAzureRMApiManagementBackend_template(data, testName)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_backend" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  protocol            = "http"
  url                 = "https://acctest"
}
`, template, data.RandomInteger)
}

func testAccAzureRMApiManagementBackend_update(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementBackend_template(data, "update")
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_backend" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
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
`, template, data.RandomInteger)
}

func testAccAzureRMApiManagementBackend_allProperties(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementBackend_template(data, "all")
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_certificate" "test" {
  name                = "example-cert"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  data                = filebase64("testdata/keyvaultcert.pfx")
  password            = ""
}

resource "azurerm_api_management_backend" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
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
      azurerm_api_management_certificate.test.thumbprint,
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
`, template, data.RandomInteger)
}

func testAccAzureRMApiManagementBackend_serviceFabric(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementBackend_template(data, "sf")
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_certificate" "test" {
  name                = "example-cert"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  data                = filebase64("testdata/keyvaultcert.pfx")
  password            = ""
}

resource "azurerm_api_management_backend" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  protocol            = "http"
  url                 = "fabric:/mytestapp/acctest"
  service_fabric_cluster {
    client_certificate_thumbprint = azurerm_api_management_certificate.test.thumbprint
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
`, template, data.RandomInteger)
}

func testAccAzureRMApiManagementBackend_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementBackend_basic(data, "requiresimport")
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_backend" "import" {
  name                = azurerm_api_management_backend.test.name
  resource_group_name = azurerm_api_management_backend.test.resource_group_name
  api_management_name = azurerm_api_management_backend.test.api_management_name
  protocol            = azurerm_api_management_backend.test.protocol
  url                 = azurerm_api_management_backend.test.url
}
`, template)
}

func testAccAzureRMApiManagementBackend_template(data acceptance.TestData, testName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d-%s"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
}
`, data.RandomInteger, testName, data.Locations.Primary, data.RandomInteger, testName)
}

func testAccAzureRMApiManagementBackend_credentialsNoCertificate(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementBackend_template(data, "all")
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_backend" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
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
`, template, data.RandomInteger)
}
