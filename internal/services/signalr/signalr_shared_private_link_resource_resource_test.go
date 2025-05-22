// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package signalr_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SignalrSharedPrivateLinkResource struct{}

func TestAccSignalrSharedPrivateLinkResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_shared_private_link_resource", "test")
	r := SignalrSharedPrivateLinkResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccSignalrSharedPrivateLinkResource_basicSites(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_shared_private_link_resource", "test")
	r := SignalrSharedPrivateLinkResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicSites(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccSignalrSharedPrivateLinkResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_shared_private_link_resource", "test")
	r := SignalrSharedPrivateLinkResource{}

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

func (r SignalrSharedPrivateLinkResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := signalr.ParseSharedPrivateLinkResourceIDInsensitively(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := clients.SignalR.SignalRClient.SharedPrivateLinkResourcesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r SignalrSharedPrivateLinkResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_key_vault" "test" {
  name                       = "vault%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id
    certificate_permissions = [
      "ManageContacts",
    ]
    key_permissions = [
      "Create",
    ]
    secret_permissions = [
      "Set",
    ]
  }
}

resource "azurerm_signalr_shared_private_link_resource" "test" {
  name               = "acctest-%d"
  signalr_service_id = azurerm_signalr_service.test.id
  sub_resource_name  = "vault"
  target_resource_id = azurerm_key_vault.test.id
  request_message    = "please approve"
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r SignalrSharedPrivateLinkResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_signalr_shared_private_link_resource" "import" {
  name               = azurerm_signalr_shared_private_link_resource.test.name
  signalr_service_id = azurerm_signalr_shared_private_link_resource.test.signalr_service_id
  sub_resource_name  = azurerm_signalr_shared_private_link_resource.test.sub_resource_name
  target_resource_id = azurerm_signalr_shared_private_link_resource.test.target_resource_id
}
`, config)
}

func (r SignalrSharedPrivateLinkResource) basicSites(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Windows"
  sku_name            = "S1"

}
resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}
}

resource "azurerm_signalr_shared_private_link_resource" "test" {
  name               = "acctest-%d"
  signalr_service_id = azurerm_signalr_service.test.id
  sub_resource_name  = "sites"
  target_resource_id = azurerm_windows_web_app.test.id
  request_message    = "please approve"
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (SignalrSharedPrivateLinkResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-signalr-%d"
  location = "%s"
}

resource "azurerm_signalr_service" "test" {
  name                = "acctestSignalr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Standard_S1"
    capacity = 1
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
