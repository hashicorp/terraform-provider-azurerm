// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SpringCloudAPIPortalResource struct{}

func TestAccSpringCloudAPIPortal_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_api_portal", "test")
	r := SpringCloudAPIPortalResource{}
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

func TestAccSpringCloudAPIPortal_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_api_portal", "test")
	r := SpringCloudAPIPortalResource{}
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

func TestAccSpringCloudAPIPortal_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_api_portal", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	r := SpringCloudAPIPortalResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, clientId, clientSecret),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sso.0.client_id", "sso.0.client_secret"),
	})
}

func TestAccSpringCloudAPIPortal_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_api_portal", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	r := SpringCloudAPIPortalResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, clientId, clientSecret),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sso.0.client_id", "sso.0.client_secret"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SpringCloudAPIPortalResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SpringCloudAPIPortalID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.AppPlatform.APIPortalClient.Get(ctx, id.ResourceGroup, id.SpringName, id.ApiPortalName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r SpringCloudAPIPortalResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%[2]d"
  location = "%[1]s"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "E0"
}

resource "azurerm_spring_cloud_gateway" "test" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
}
`, data.Locations.Primary, data.RandomInteger)
}

func (r SpringCloudAPIPortalResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_api_portal" "test" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
}
`, template)
}

func (r SpringCloudAPIPortalResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_api_portal" "import" {
  name                    = azurerm_spring_cloud_api_portal.test.name
  spring_cloud_service_id = azurerm_spring_cloud_api_portal.test.spring_cloud_service_id
}
`, config)
}

func (r SpringCloudAPIPortalResource) complete(data acceptance.TestData, clientId, clientSecret string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
data "azurerm_client_config" "current" {
}

resource "azurerm_spring_cloud_api_portal" "test" {
  name                          = "default"
  spring_cloud_service_id       = azurerm_spring_cloud_service.test.id
  gateway_ids                   = [azurerm_spring_cloud_gateway.test.id]
  https_only_enabled            = false
  public_network_access_enabled = false
  instance_count                = 1

  sso {
    client_id     = "%s"
    client_secret = "%s"
    issuer_uri    = "https://login.microsoftonline.com/${data.azurerm_client_config.current.tenant_id}/v2.0"
    scope         = ["read"]
  }
}
`, template, clientId, clientSecret)
}
