// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_api_definition_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apidefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiCenterApiDefinitionResource struct{}

func TestAccApiCenterApiDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_api_definition", "test")
	r := ApiCenterApiDefinitionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("api_version_id").MatchesOtherKey(check.That("azurerm_api_center_api_version.test").Key("id")),
				check.That(data.ResourceName).Key("display_name").MatchesOtherKey(check.That(data.ResourceName).Key("name")),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiCenterApiDefinition_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_api_definition", "test")
	r := ApiCenterApiDefinitionResource{}

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

func TestAccApiCenterApiDefinition_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_api_definition", "test")
	r := ApiCenterApiDefinitionResource{}

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

func TestAccApiCenterApiDefinition_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_api_definition", "test")
	r := ApiCenterApiDefinitionResource{}

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

func (ApiCenterApiDefinitionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := apidefinitions.ParseDefinitionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiCenter.ApiDefinitionsClient.Get(ctx, *id)
	if response.WasNotFound(resp.HttpResponse) {
		return pointer.To(false), nil
	}
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (ApiCenterApiDefinitionResource) template(data acceptance.TestData) string {
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

resource "azurerm_api_center_api" "test" {
  name                  = "acctest-api-%[1]d"
  service_id = azurerm_api_center_service.test.id
  display_name          = "Acceptance Test API"
  kind                  = "rest"
  terms_of_service_url  = "https://example.com/terms"
}

resource "azurerm_api_center_api_version" "test" {
  name              = "acctest-version-%[1]d"
  api_id = azurerm_api_center_api.test.id
  lifecycle_stage   = "testing"
  display_name      = "Acceptance Test API Version"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ApiCenterApiDefinitionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_center_api_definition" "test" {
  name                      = "acctest%d"
  api_version_id = azurerm_api_center_api_version.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r ApiCenterApiDefinitionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_center_api_definition" "test" {
  name                      = "acctest%d"
  api_version_id = azurerm_api_center_api_version.test.id
  display_name              = "Complete API Definition"
  description               = "An API definition for acceptance testing."
}
`, r.template(data), data.RandomInteger)
}

func (r ApiCenterApiDefinitionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_center_api_definition" "import" {
  name                      = azurerm_api_center_api_definition.test.name
  api_version_id = azurerm_api_center_api_definition.test.api_version_id
  display_name              = azurerm_api_center_api_definition.test.display_name
}
`, r.basic(data))
}
