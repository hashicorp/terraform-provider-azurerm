// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apicenter_test

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

func TestAccApiCenterEnvironment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_environment", "test")
	r := ApiCenterEnvironmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
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

func (r ApiCenterEnvironmentResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_center_environment" "test" {
  name                  = "test"
  api_center_service_id = azurerm_api_center_service.test.id
  title                 = "testtitle"
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
  api_center_service_id  = azurerm_api_center_service.test.id
  title                  = "testtitle"
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
  api_center_service_id  = azurerm_api_center_service.test.id
  title                  = "testtitle"
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
  api_center_service_id = azurerm_api_center_environment.test.api_center_service_id
  title                 = azurerm_api_center_environment.test.title
  environment_type      = azurerm_api_center_environment.test.environment_type
}
`, template)
}
