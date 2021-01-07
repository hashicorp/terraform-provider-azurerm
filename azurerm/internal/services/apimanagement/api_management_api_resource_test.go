package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ApiManagementApiResource struct {
}

func TestAccApiManagementApi_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("soap_pass_through").HasValue("false"),
				check.That(data.ResourceName).Key("is_current").HasValue("true"),
				check.That(data.ResourceName).Key("is_online").HasValue("false"),
				check.That(data.ResourceName).Key("subscription_required").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApi_wordRevision(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.wordRevision(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("revision").HasValue("one-point-oh"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApi_blankPath(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.blankPath(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("soap_pass_through").HasValue("false"),
				check.That(data.ResourceName).Key("is_current").HasValue("true"),
				check.That(data.ResourceName).Key("is_online").HasValue("false"),
				check.That(data.ResourceName).Key("path").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApi_version(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.versionSet(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("version").HasValue("v1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApi_oauth2Authorization(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.oauth2Authorization(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApi_openidAuthentication(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.openidAuthentication(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApi_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccApiManagementApi_soapPassthrough(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.soapPassthrough(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApi_subscriptionRequired(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.subscriptionRequired(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subscription_required").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApi_importSwagger(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.importSwagger(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
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
	})
}

func TestAccApiManagementApi_importWsdl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.importWsdl(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
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
	})
}

func TestAccApiManagementApi_importUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.importWsdl(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
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
			Config: r.importSwagger(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
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
	})
}

func TestAccApiManagementApi_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t ApiManagementApiResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	apiid := id.Path["apis"]

	resp, err := clients.ApiManagement.ApiClient.Get(ctx, resourceGroup, serviceName, apiid)
	if err != nil {
		return nil, fmt.Errorf("reading ApiManagementApi (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r ApiManagementApiResource) basic(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementApiResource) blankPath(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementApiResource) wordRevision(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementApiResource) soapPassthrough(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementApiResource) subscriptionRequired(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementApiResource) requiresImport(data acceptance.TestData) string {
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
`, r.basic(data))
}

func (r ApiManagementApiResource) importSwagger(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementApiResource) importWsdl(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementApiResource) complete(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementApiResource) versionSet(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementApiResource) oauth2Authorization(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementApiResource) openidAuthentication(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (ApiManagementApiResource) template(data acceptance.TestData) string {
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
