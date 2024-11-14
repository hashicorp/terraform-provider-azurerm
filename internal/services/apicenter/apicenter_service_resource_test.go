// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apicenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ApiCenterServiceResource struct{}

func TestAccApicenterService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_apicenter_service", "test")
	r := ApiCenterServiceResource{}

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

func TestAccDatabricksAccessConnector_identityUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_access_connector", "test")
	r := ApiCenterServiceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identityUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabricksAccessConnector_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_access_connector", "test")
	r := ApiCenterServiceResource{}

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

func (ApiCenterServiceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := services.ParseServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiCenter.ServicesClient.Get(ctx, *id)

	if response.WasNotFound(resp.HttpResponse) {
		return pointer.To(false), nil
	}

	if err != nil {
		return nil, fmt.Errorf("making Read request on Databricks %s: %+v", id.ID(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (ApiCenterServiceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-apicenter-%d"
  location = "%s"
}
`, data.RandomInteger, "eastus") // Only available in a few select regions for now
}

func (r ApiCenterServiceResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_apicenter_service" "test" {
  name                = "acctestApiCSvc%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func (r ApiCenterServiceResource) identityUserAssigned(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestAPICUAI-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_apicenter_service" "test" {
  name                = "acctestApiCSvc%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}
`, template, data.RandomInteger)
}

func (ApiCenterServiceResource) requiresImport(data acceptance.TestData) string {
	template := ApiCenterServiceResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_apicenter_service" "import" {
  name                = azurerm_apicenter_service.test.name
  resource_group_name = azurerm_apicenter_service.test.resource_group_name
  location            = azurerm_apicenter_service.test.location
}
`, template)
}
