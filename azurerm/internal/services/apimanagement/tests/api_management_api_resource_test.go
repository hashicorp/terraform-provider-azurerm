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

func TestAccAzureRMApiManagementApi_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApi_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "soap_pass_through", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "is_current", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "is_online", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "subscription_required", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementApi_wordRevision(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApi_wordRevision(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "revision", "one-point-oh"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementApi_blankPath(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApi_blankPath(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "soap_pass_through", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "is_current", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "is_online", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "path", ""),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementApi_version(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApi_versionSet(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "v1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementApi_oauth2Authorization(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApi_oauth2Authorization(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementApi_openidAuthentication(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApi_openidAuthentication(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementApi_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApi_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementApi_requiresImport),
		},
	})
}

func TestAccAzureRMApiManagementApi_soapPassthrough(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApi_soapPassthrough(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementApi_subscriptionRequired(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApi_subscriptionRequired(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "subscription_required", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementApi_importSwagger(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApi_importSwagger(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(data.ResourceName),
				),
			},
			{
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"import",
				},
			},
		},
	})
}

func TestAccAzureRMApiManagementApi_importWsdl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApi_importWsdl(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(data.ResourceName),
				),
			},
			{
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"import",
				},
			},
		},
	})
}

func TestAccAzureRMApiManagementApi_importUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApi_importWsdl(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(data.ResourceName),
				),
			},
			{
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"import",
				},
			},
			{
				Config: testAccAzureRMApiManagementApi_importSwagger(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(data.ResourceName),
				),
			},
			{
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"import",
				},
			},
		},
	})
}

func TestAccAzureRMApiManagementApi_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApi_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMApiManagementApiDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ApiClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_api" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		revision := rs.Primary.Attributes["revision"]
		apiId := fmt.Sprintf("%s;rev=%s", name, revision)

		resp, err := conn.Get(ctx, resourceGroup, serviceName, apiId)
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

func testCheckAzureRMApiManagementApiExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ApiClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		revision := rs.Primary.Attributes["revision"]

		apiId := fmt.Sprintf("%s;rev=%s", name, revision)
		resp, err := conn.Get(ctx, resourceGroup, serviceName, apiId)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: API %q Revision %q (API Management Service %q / Resource Group: %q) does not exist", name, revision, serviceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on apiManagementClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementApi_basic(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApi_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"
}
`, template, data.RandomInteger)
}

func testAccAzureRMApiManagementApi_blankPath(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApi_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api1"
  path                = ""
  protocols           = ["https"]
  revision            = "1"
}
`, template, data.RandomInteger)
}

func testAccAzureRMApiManagementApi_wordRevision(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApi_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "one-point-oh"
}
`, template, data.RandomInteger)
}

func testAccAzureRMApiManagementApi_soapPassthrough(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApi_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"
  soap_pass_through   = true
}
`, template, data.RandomInteger)
}

func testAccAzureRMApiManagementApi_subscriptionRequired(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApi_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "test" {
  name                  = "acctestapi-%d"
  resource_group_name   = azurerm_resource_group.test.name
  api_management_name   = azurerm_api_management.test.name
  display_name          = "api1"
  path                  = "api1"
  protocols             = ["https"]
  revision              = "1"
  subscription_required = false
}
`, template, data.RandomInteger)
}

func testAccAzureRMApiManagementApi_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApi_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "import" {
  name                = azurerm_api_management_api.test.name
  resource_group_name = azurerm_api_management_api.test.resource_group_name
  api_management_name = azurerm_api_management_api.test.api_management_name
  display_name        = azurerm_api_management_api.test.display_name
  path                = azurerm_api_management_api.test.path
  protocols           = azurerm_api_management_api.test.protocols
  revision            = azurerm_api_management_api.test.revision
}
`, template)
}

func testAccAzureRMApiManagementApi_importSwagger(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApi_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"

  import {
    content_value  = file("testdata/api_management_api_swagger.json")
    content_format = "swagger-json"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMApiManagementApi_importWsdl(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApi_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"

  import {
    content_value  = file("testdata/api_management_api_wsdl.xml")
    content_format = "wsdl"

    wsdl_selector {
      service_name  = "Calculator"
      endpoint_name = "CalculatorHttpsSoap11Endpoint"
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMApiManagementApi_complete(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApi_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "Butter Parser"
  path                = "butter-parser"
  protocols           = ["https", "http"]
  revision            = "3"
  description         = "What is my purpose? You parse butter."
  service_url         = "https://example.com/foo/bar"

  subscription_key_parameter_names {
    header = "X-Butter-Robot-API-Key"
    query  = "location"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMApiManagementApi_versionSet(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApi_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_version_set" "test" {
  name                = "acctestAMAVS-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "Butter Parser"
  versioning_scheme   = "Segment"
}

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"
  version             = "v1"
  version_set_id      = azurerm_api_management_api_version_set.test.id
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementApi_oauth2Authorization(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApi_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_authorization_server" "test" {
  name                         = "acctestauthsrv-%d"
  resource_group_name          = azurerm_resource_group.test.name
  api_management_name          = azurerm_api_management.test.name
  display_name                 = "Test Group"
  authorization_endpoint       = "https://azacctest.hashicorptest.com/client/authorize"
  client_id                    = "42424242-4242-4242-4242-424242424242"
  client_registration_endpoint = "https://azacctest.hashicorptest.com/client/register"

  grant_types = [
    "implicit",
  ]

  authorization_methods = [
    "GET",
  ]
}

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"
  oauth2_authorization {
    authorization_server_name = azurerm_api_management_authorization_server.test.name
    scope                     = "acctest"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementApi_openidAuthentication(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApi_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_openid_connect_provider" "test" {
  name                = "acctest-%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  client_id           = "00001111-2222-3333-%d"
  client_secret       = "%d-cwdavsxbacsaxZX-%d"
  display_name        = "Initial Name"
  metadata_endpoint   = "https://azacctest.hashicorptest.com/example/foo"
}

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"
  openid_authentication {
    openid_provider_name = azurerm_api_management_openid_connect_provider.test.name
    bearer_token_sending_methods = [
      "authorizationHeader",
      "query",
    ]
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementApi_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
