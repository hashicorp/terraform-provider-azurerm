// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apicenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ApiCenterEnvironmentResource struct{}

func TestAccApicenterEnvironment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_apicenter_environment", "test")
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

func TestAccApicenterEnvironment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_apicenter_environment", "test")
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

	if response.WasNotFound(resp.HttpResponse) {
		return pointer.To(false), nil
	}

	if err != nil {
		return nil, fmt.Errorf("making Read request on ApiCenter Environment %s: %+v", id.ID(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (ApiCenterEnvironmentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-apicenter-%d"
  location = "%s"
}

resource "azurerm_apicenter_service" "test" {
  name                = "acctestApiCSvc%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}

`, data.RandomInteger, "eastus", data.RandomInteger) // Only available in a few select regions for now
}

func (r ApiCenterEnvironmentResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_apicenter_environment" "test" {
  name             = "test"
  service_id       = azurerm_apicenter_service.test.id
  identification   = "testid"
  environment_type = "testing"
  description      = "testing environment"
}
`, template)
}

func (ApiCenterEnvironmentResource) requiresImport(data acceptance.TestData) string {
	template := ApiCenterEnvironmentResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_apicenter_environment" "import" {
  name             = azurerm_apicenter_environment.test.name
  service_id       = azurerm_apicenter_environment.test.service_id
  identification   = azurerm_apicenter_environment.test.identification
  environment_type = azurerm_apicenter_environment.test.environment_type
}
`, template)
}
