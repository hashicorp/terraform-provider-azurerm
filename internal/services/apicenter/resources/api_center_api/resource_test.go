// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_api_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiCenterApiResource struct{}

func TestAccApiCenterApi_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_api", "test")
	r := ApiCenterApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("api_center_service_id").MatchesOtherKey(check.That("azurerm_api_center_service.test").Key("id")),
				check.That(data.ResourceName).Key("display_name").MatchesOtherKey(check.That(data.ResourceName).Key("name")),
				check.That(data.ResourceName).Key("lifecycle_stage").HasValue("design"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiCenterApi_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_api", "test")
	r := ApiCenterApiResource{}

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

func TestAccApiCenterApi_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_api", "test")
	r := ApiCenterApiResource{}

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

func TestAccApiCenterApi_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_api", "test")
	r := ApiCenterApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ApiCenterApiResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := apis.ParseApiID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiCenter.ApisClient.Get(ctx, *id)
	if response.WasNotFound(resp.HttpResponse) {
		return pointer.To(false), nil
	}
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (ApiCenterApiResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-apicenter-%[1]d"
  location = "%[2]s"
}

resource "azurerm_api_center_service" "test" {
  name                = "acctestApiCSvc%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ApiCenterApiResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_center_api" "test" {
  name                  = "acctest%d"
  api_center_service_id = azurerm_api_center_service.test.id
  kind                  = "rest"
}
`, r.template(data), data.RandomInteger)
}

func (r ApiCenterApiResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_center_api" "test" {
  name                  = "acctest%d"
  api_center_service_id = azurerm_api_center_service.test.id
  display_name          = "Acceptance Test GraphQL API"
  kind                  = "graphql"
  description           = "An API for acceptance testing."
  summary               = "Acceptance test API summary"
  terms_of_service_url  = "https://example.com/terms"

  contacts {
    name  = "API Support"
    email = "support@example.com"
    url   = "https://example.com/support"
  }

  external_documentation {
    title       = "API Documentation"
    description = "Documentation for the acceptance test API."
    url         = "https://example.com/docs"
  }

  license {
    name       = "MIT License"
    url        = "https://opensource.org/licenses/MIT"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApiCenterApiResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_center_api" "import" {
  name                  = azurerm_api_center_api.test.name
  api_center_service_id = azurerm_api_center_api.test.api_center_service_id
  display_name          = azurerm_api_center_api.test.display_name
  kind                  = azurerm_api_center_api.test.kind
}
`, r.basic(data))
}
