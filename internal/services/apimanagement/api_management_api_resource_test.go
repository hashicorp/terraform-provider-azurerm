// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementApiResource struct{}

const (
	SkuNameConsumption = "Consumption_0"
	SkuNameDeveloper   = "Developer_1"
)

func TestAccApiManagementApi_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("api_type").HasValue("http"),
				check.That(data.ResourceName).Key("is_current").HasValue("true"),
				check.That(data.ResourceName).Key("is_online").HasValue("false"),
				check.That(data.ResourceName).Key("subscription_required").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApi_wordRevision(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.wordRevision(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.blankPath(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("api_type").HasValue("http"),
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.versionSet(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.oauth2Authorization(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.oauth2AuthorizationUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApi_openidAuthentication(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.openidAuthentication(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.openidAuthenticationUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApi_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccApiManagementApi_graphql(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.graphql(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("api_type").HasValue("graphql"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApi_soap(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.soap(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("api_type").HasValue("soap"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApi_websocket(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.websocket(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("api_type").HasValue("websocket"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApi_subscriptionRequired(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.subscriptionRequired(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.importSwagger(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("import"),
	})
}

func TestAccApiManagementApi_importOpenapi(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.importOpenapi(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("import"),
	})
}

func TestAccApiManagementApi_importOpenapiInvalid(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.importOpenapiInvalid(data),
			ExpectError: regexp.MustCompile("ValidationError"),
		},
	})
}

func TestAccApiManagementApi_updateImportOpenapiInvalid(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.importOpenapi(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("import"),
		{
			Config:      r.importOpenapiInvalid(data),
			ExpectError: regexp.MustCompile("ValidationError"),
		},
		{
			Config: r.importOpenapi(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("import"),
	})
}

func TestAccApiManagementApi_importSwaggerWithServiceUrl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.importSwaggerWithServiceUrl(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("import"),
		{
			Config: r.importSwaggerWithServiceUrlUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("import"),
	})
}

func TestAccApiManagementApi_importWsdl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.importWsdl(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("import"),
	})
}

func TestAccApiManagementApi_importWsdlWithSelector(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.importWsdlWithSelector(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("import"),
	})
}

func TestAccApiManagementApi_importUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.importWsdl(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("import"),
		{
			Config: r.importSwagger(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("import"),
	})
}

func TestAccApiManagementApi_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApi_completeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApi_cloneApi(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "clone")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cloneApi(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("source_api_id"),
	})
}

func TestAccApiManagementApi_createNewVersionFromExisting(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "version")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.createNewVersionFromExisting(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("source_api_id"),
	})
}

func TestAccApiManagementApi_createRevisionFromExisting(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "revision")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.createRevisionFromExisting(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("source_api_id"),
	})
}

func TestAccApiManagementApi_createRevisionFromExistingRevision(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "revision")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.createRevisionFromExistingRevision(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("source_api_id"),
	})
}

func TestAccApiManagementApi_contact(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "test")
	r := ApiManagementApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.contact(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ApiManagementApiResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := api.ParseApiID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.ApiClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
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
`, r.template(data, SkuNameConsumption), data.RandomInteger)
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
`, r.template(data, SkuNameConsumption), data.RandomInteger)
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
`, r.template(data, SkuNameConsumption), data.RandomInteger)
}

func (r ApiManagementApiResource) graphql(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  api_type            = "graphql"
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"
}
`, r.template(data, SkuNameConsumption), data.RandomInteger)
}

func (r ApiManagementApiResource) soap(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  api_type            = "soap"
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"
}
`, r.template(data, SkuNameConsumption), data.RandomInteger)
}

func (r ApiManagementApiResource) websocket(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  api_type            = "websocket"
  display_name        = "api1"
  path                = "api1"
  protocols           = ["wss"]
  revision            = "1"
  service_url         = "wss://example.com/foo/bar"
}
`, r.template(data, SkuNameDeveloper), data.RandomInteger)
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
`, r.template(data, SkuNameConsumption), data.RandomInteger)
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

func (r ApiManagementApiResource) importSwagger(data acceptance.TestData, ignoreImported bool) string {
	ignoreConfig := ""
	if ignoreImported {
		ignoreConfig = `lifecycle {
	ignore_changes = [description, display_name, contact, license]
}
`
	}

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
  %s
}
`, r.template(data, SkuNameConsumption), data.RandomInteger, ignoreConfig)
}

func (r ApiManagementApiResource) importOpenapi(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "current"

  import {
    content_value  = file("testdata/api_management_api_openapi.yaml")
    content_format = "openapi"
  }

  lifecycle {
    ignore_changes = [description]
  }
}
`, r.template(data, SkuNameConsumption), data.RandomInteger)
}

func (r ApiManagementApiResource) importOpenapiInvalid(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "current"

  import {
    content_value  = file("testdata/api_management_api_openapi_invalid.yaml")
    content_format = "openapi"
  }
}
`, r.template(data, SkuNameConsumption), data.RandomInteger)
}

func (r ApiManagementApiResource) importSwaggerWithServiceUrl(data acceptance.TestData) string {
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
  description         = "import swagger"
  service_url         = "https://example.com/foo/bar"

  import {
    content_value  = file("testdata/api_management_api_swagger.json")
    content_format = "swagger-json"
  }
}
`, r.template(data, SkuNameConsumption), data.RandomInteger)
}

func (r ApiManagementApiResource) importSwaggerWithServiceUrlUpdate(data acceptance.TestData) string {
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
  description         = "import swagger update"
  service_url         = "https://example.com/foo/bar"

  import {
    content_value  = file("testdata/api_management_api_swagger.json")
    content_format = "swagger-json"
  }
}
`, r.template(data, SkuNameConsumption), data.RandomInteger)
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
  }
}
`, r.template(data, SkuNameConsumption), data.RandomInteger)
}

func (r ApiManagementApiResource) importWsdlWithSelector(data acceptance.TestData) string {
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
    content_value  = file("testdata/api_management_api_wsdl_multiple.xml")
    content_format = "wsdl"

    wsdl_selector {
      service_name  = "Calculator"
      endpoint_name = "CalculatorHttpsSoap11Endpoint"
    }
  }
}
`, r.template(data, SkuNameConsumption), data.RandomInteger)
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

  contact {
    email = "test@test.com"
    name  = "test"
    url   = "https://example:8080"
  }

  license {
    name = "test-license"
    url  = "https://example:8080/license"
  }

  terms_of_service_url = "https://example:8080/service"
}
`, r.template(data, SkuNameConsumption), data.RandomInteger)
}

func (r ApiManagementApiResource) completeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "Butter Parser Update"
  path                = "butter-parser-update"
  protocols           = ["https"]
  revision            = "3"
  description         = "What is my purpose? You parse butter."
  service_url         = "https://example.com/foo/bar/update"

  subscription_key_parameter_names {
    header = "X-Butter-Robot-API-Key"
    query  = "location-update"
  }

  contact {
    email = "test-update@test.com"
    name  = "test-update"
    url   = "https://example-update:8080"
  }

  license {
    name = "test-license-update"
    url  = "https://example:8080/license-update"
  }

  terms_of_service_url = "https://example:8080/service-update"
}
`, r.template(data, SkuNameConsumption), data.RandomInteger)
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
`, r.template(data, SkuNameConsumption), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementApiResource) oauth2Authorization(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_authorization_server" "test" {
  name                         = "acctestauthsrv-%d"
  resource_group_name          = azurerm_resource_group.test.name
  api_management_name          = azurerm_api_management.test.name
  display_name                 = "Test Group"
  authorization_endpoint       = "https://azacceptance.hashicorptest.com/client/authorize"
  client_id                    = "42424242-4242-4242-4242-424242424242"
  client_registration_endpoint = "https://azacceptance.hashicorptest.com/client/register"

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
`, r.template(data, SkuNameConsumption), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementApiResource) oauth2AuthorizationUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_authorization_server" "test" {
  name                         = "acctestauthsrv-%d"
  resource_group_name          = azurerm_resource_group.test.name
  api_management_name          = azurerm_api_management.test.name
  display_name                 = "Test Group"
  authorization_endpoint       = "https://azacceptance.hashicorptest.com/client/authorize"
  client_id                    = "42424242-4242-4242-4242-424242424242"
  client_registration_endpoint = "https://azacceptance.hashicorptest.com/client/register"

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
  display_name        = "api1update"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"
  oauth2_authorization {
    authorization_server_name = azurerm_api_management_authorization_server.test.name
    scope                     = "acctest"
  }
}
`, r.template(data, SkuNameConsumption), data.RandomInteger, data.RandomInteger)
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
  metadata_endpoint   = "https://azacceptance.hashicorptest.com/example/foo"
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
`, r.template(data, SkuNameConsumption), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementApiResource) openidAuthenticationUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_openid_connect_provider" "test" {
  name                = "acctest-%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  client_id           = "00001111-2222-3333-%d"
  client_secret       = "%d-cwdavsxbacsaxZX-%d"
  display_name        = "Initial Name"
  metadata_endpoint   = "https://azacceptance.hashicorptest.com/example/foo"
}

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api1update"
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
`, r.template(data, SkuNameConsumption), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementApiResource) cloneApi(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "clone" {
  name                = "acctestClone-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api1_clone"
  revision            = "1"
  source_api_id       = azurerm_api_management_api.test.id
  description         = "Copy of Existing Echo Api including Operations."
}
`, r.basic(data), data.RandomInteger)
}

func (r ApiManagementApiResource) createNewVersionFromExisting(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_version_set" "test" {
  name                = "acctestAMAVS-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "Butter Parser"
  versioning_scheme   = "Segment"
}

resource "azurerm_api_management_api" "version" {
  name                = "acctestVersion-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api_version"
  revision            = "1"
  source_api_id       = azurerm_api_management_api.test.id
  version             = "v1"
  version_set_id      = azurerm_api_management_api_version_set.test.id
  version_description = "Create Echo API into a new Version using Existing Version Set and Copy all Operations."
}
`, r.basic(data), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementApiResource) createRevisionFromExisting(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "revision" {
  name                 = "acctestRevision-%d"
  resource_group_name  = azurerm_resource_group.test.name
  api_management_name  = azurerm_api_management.test.name
  revision             = "18"
  source_api_id        = azurerm_api_management_api.test.id
  revision_description = "Creating a Revision of an existing API"
}
`, r.basic(data), data.RandomInteger)
}

func (r ApiManagementApiResource) createRevisionFromExistingRevision(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "revision" {
  name                 = "acctestRevision-%d"
  resource_group_name  = azurerm_resource_group.test.name
  api_management_name  = azurerm_api_management.test.name
  revision             = "18"
  description          = "What is my purpose? You parse butter."
  source_api_id        = "${azurerm_api_management_api.test.id};rev=3"
  revision_description = "Creating a Revision of an existing API"
  contact {
    email = "test@test.com"
    name  = "test"
    url   = "https://example:8080"
  }

  license {
    name = "test-license"
    url  = "https://example:8080/license"
  }

  terms_of_service_url = "https://example:8080/service"
}
`, r.complete(data), data.RandomInteger)
}

func (r ApiManagementApiResource) contact(data acceptance.TestData) string {
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

  contact {
    email = "test@test.com"
    name  = "test"
    url   = "https://example:8080"
  }

  license {
    name = "test-license"
    url  = "https://example:8080/license"
  }

  terms_of_service_url = "https://example:8080/service"
}
`, r.template(data, SkuNameConsumption), data.RandomInteger)
}

func (ApiManagementApiResource) template(data acceptance.TestData, skuName string) string {
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

  sku_name = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, skuName)
}
