// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_environment_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiCenterEnvironmentResource struct{}

func TestAccApiCenterEnvironment_customProperties(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_environment", "test")
	r := ApiCenterEnvironmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customProperties(data, "catalog-team", "medium"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("custom_properties").HasValue(`{"governance":{"criticality":"medium","owning_team":"catalog-team"}}`),
			),
		},
		data.ImportStep(),
		{
			Config: r.customProperties(data, "platform-team", "high"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("custom_properties").HasValue(`{"governance":{"criticality":"high","owning_team":"platform-team"}}`),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiCenterEnvironment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_environment", "test")
	r := ApiCenterEnvironmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").MatchesOtherKey(check.That(data.ResourceName).Key("name")),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiCenterEnvironment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_environment", "test")
	r := ApiCenterEnvironmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiCenterEnvironment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_environment", "test")
	r := ApiCenterEnvironmentResource{}

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

func (ApiCenterEnvironmentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := environments.ParseEnvironmentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiCenter.EnvironmentsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (ApiCenterEnvironmentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-apicenter-%d"
  location = "%s"
}

resource "azurerm_api_center_service" "test" {
  name                = "acctestApiCSvc%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}

`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ApiCenterEnvironmentResource) customProperties(data acceptance.TestData, team, criticality string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_center_metadata_schema" "test" {
  name       = "governance"
  service_id = azurerm_api_center_service.test.id

  metadata_assignment {
    entity   = "environment"
    required = true
  }

  schema = jsonencode({
    type = "object"
    properties = {
      owning_team = { type = "string" }
      criticality = {
        type = "string"
        oneOf = [
          {
            const       = "low"
            description = "Low criticality"
          },
          {
            const       = "medium"
            description = "Medium criticality"
          },
          {
            const       = "high"
            description = "High criticality"
          },
        ]
      }
    }
  })
}

resource "azurerm_api_center_environment" "test" {
  name             = "test"
  service_id       = azurerm_api_center_service.test.id
  environment_type = "testing"

  custom_properties = jsonencode({
    (azurerm_api_center_metadata_schema.test.name) = {
      owning_team = %q
      criticality = %q
    }
  })
}
`, r.template(data), team, criticality)
}

func (r ApiCenterEnvironmentResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_center_environment" "test" {
  name                  = "test"
  service_id = azurerm_api_center_service.test.id
  environment_type      = "testing"
  description           = "testing environment"
}`, template)
}

func (r ApiCenterEnvironmentResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_center_environment" "test" {
  name                   = "test"
  service_id  = azurerm_api_center_service.test.id
  display_name           = "testtitle"
  environment_type       = "testing"
  description            = "testing environment"
  development_portal_uri = "https://developer.com"
  instructions           = "Use this wonderful API to CRUD brilliant data."
  server_type            = "Azure API Management"
  management_portal_uri  = "https://azure-apim-mgmt-portal.azure.com"
}
`, template)
}

func (r ApiCenterEnvironmentResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_center_environment" "test" {
  name                   = "test"
  service_id  = azurerm_api_center_service.test.id
  display_name           = "testtitle"
  environment_type       = "testing"
  description            = "testing environment 2"
  development_portal_uri = "https://developer2.com"
  instructions           = "Use this wonderful API2 to CRUD brilliant data."
  server_type            = "Apigee API Management"
  management_portal_uri  = "https://azure-apim-mgmt-portal2.azure.com"
}
`, template)
}

func (ApiCenterEnvironmentResource) requiresImport(data acceptance.TestData) string {
	template := ApiCenterEnvironmentResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_center_environment" "import" {
  name                  = azurerm_api_center_environment.test.name
  service_id = azurerm_api_center_environment.test.service_id
  display_name          = azurerm_api_center_environment.test.display_name
  environment_type      = azurerm_api_center_environment.test.environment_type
}
`, template)
}
