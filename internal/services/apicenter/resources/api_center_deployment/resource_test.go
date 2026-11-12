// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_deployment_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/deployments"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiCenterDeploymentResource struct{}

func TestAccApiCenterDeployment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_deployment", "test")
	r := ApiCenterDeploymentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("api_center_api_definition_id").MatchesOtherKey(check.That("azurerm_api_center_api_definition.test").Key("id")),
				check.That(data.ResourceName).Key("display_name").MatchesOtherKey(check.That(data.ResourceName).Key("name")),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiCenterDeployment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_deployment", "test")
	r := ApiCenterDeploymentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("api_center_api_definition_id").MatchesOtherKey(check.That("azurerm_api_center_api_definition.test").Key("id")),
				check.That(data.ResourceName).Key("environment_id").MatchesOtherKey(check.That("azurerm_api_center_environment.test").Key("id")),
				check.That(data.ResourceName).Key("description").HasValue("An API deployment for acceptance testing."),
				check.That(data.ResourceName).Key("display_name").HasValue("Acceptance Test Deployment"),
				check.That(data.ResourceName).Key("state").HasValue("active"),
				check.That(data.ResourceName).Key("server_runtime_uri.#").HasValue("2"),
				resource.TestCheckTypeSetElemAttr(data.ResourceName, "server_runtime_uri.*", "https://example.com/api"),
				resource.TestCheckTypeSetElemAttr(data.ResourceName, "server_runtime_uri.*", "https://example.com/api-secondary"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiCenterDeployment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_deployment", "test")
	r := ApiCenterDeploymentResource{}

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

func TestAccApiCenterDeployment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_deployment", "test")
	r := ApiCenterDeploymentResource{}

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
				check.That(data.ResourceName).Key("api_center_api_definition_id").MatchesOtherKey(check.That("azurerm_api_center_api_definition.test").Key("id")),
				check.That(data.ResourceName).Key("environment_id").MatchesOtherKey(check.That("azurerm_api_center_environment.test").Key("id")),
				check.That(data.ResourceName).Key("description").HasValue("An API deployment for acceptance testing."),
				check.That(data.ResourceName).Key("display_name").HasValue("Acceptance Test Deployment"),
				check.That(data.ResourceName).Key("state").HasValue("active"),
				check.That(data.ResourceName).Key("server_runtime_uri.#").HasValue("2"),
				resource.TestCheckTypeSetElemAttr(data.ResourceName, "server_runtime_uri.*", "https://example.com/api"),
				resource.TestCheckTypeSetElemAttr(data.ResourceName, "server_runtime_uri.*", "https://example.com/api-secondary"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("server_runtime_uri.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func (ApiCenterDeploymentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := deployments.ParseDeploymentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiCenter.DeploymentsClient.Get(ctx, *id)
	if response.WasNotFound(resp.HttpResponse) {
		return pointer.To(false), nil
	}
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (ApiCenterDeploymentResource) template(data acceptance.TestData) string {
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
  api_center_service_id = azurerm_api_center_service.test.id
  display_name          = "Acceptance Test API"
  kind                  = "rest"
  terms_of_service_url  = "https://example.com/terms"
}

resource "azurerm_api_center_api_version" "test" {
  name              = "acctest-version-%[1]d"
  api_center_api_id = azurerm_api_center_api.test.id
  lifecycle_stage   = "testing"
}

resource "azurerm_api_center_api_definition" "test" {
  name                      = "acctest-definition-%[1]d"
  api_center_api_version_id = azurerm_api_center_api_version.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ApiCenterDeploymentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_center_deployment" "test" {
  name                         = "acctest%d"
  api_center_api_definition_id = azurerm_api_center_api_definition.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r ApiCenterDeploymentResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_api_center_environment" "test" {
  name                  = "acctest-environment-%[2]d"
  api_center_service_id = azurerm_api_center_service.test.id
  display_name          = "Acceptance Test Environment"
  environment_type      = "testing"
}

resource "azurerm_api_center_deployment" "test" {
  name                         = "acctest%[2]d"
  api_center_api_definition_id = azurerm_api_center_api_definition.test.id
  description                  = "An API deployment for acceptance testing."
  environment_id               = azurerm_api_center_environment.test.id
  state                        = "active"
  display_name                 = "Acceptance Test Deployment"

  server_runtime_uri = [
    "https://example.com/api",
    "https://example.com/api-secondary",
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r ApiCenterDeploymentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_center_deployment" "import" {
  name                         = azurerm_api_center_deployment.test.name
  api_center_api_definition_id = azurerm_api_center_deployment.test.api_center_api_definition_id
}
`, r.basic(data))
}
