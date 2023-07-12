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

type SpringCloudDevToolPortalResource struct{}

func TestAccSpringCloudDevToolPortal_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_dev_tool_portal", "test")
	r := SpringCloudDevToolPortalResource{}
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

func TestAccSpringCloudDevToolPortal_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_dev_tool_portal", "test")
	r := SpringCloudDevToolPortalResource{}
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

func TestAccSpringCloudDevToolPortal_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_dev_tool_portal", "test")
	r := SpringCloudDevToolPortalResource{}
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
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

func TestAccSpringCloudDevToolPortal_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_dev_tool_portal", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	r := SpringCloudDevToolPortalResource{}
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

func (r SpringCloudDevToolPortalResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SpringCloudDevToolPortalID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.AppPlatform.DevToolPortalClient.Get(ctx, id.ResourceGroup, id.SpringName, id.DevToolPortalName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r SpringCloudDevToolPortalResource) template(data acceptance.TestData) string {
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
`, data.Locations.Primary, data.RandomInteger)
}

func (r SpringCloudDevToolPortalResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_dev_tool_portal" "test" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
}
`, template)
}

func (r SpringCloudDevToolPortalResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_dev_tool_portal" "import" {
  name                    = azurerm_spring_cloud_dev_tool_portal.test.name
  spring_cloud_service_id = azurerm_spring_cloud_dev_tool_portal.test.spring_cloud_service_id
}
`, config)
}

func (r SpringCloudDevToolPortalResource) complete(data acceptance.TestData, clientId, clientSecret string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

resource "azurerm_spring_cloud_dev_tool_portal" "test" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id

  public_network_access_enabled = true

  sso {
    client_id     = "%s"
    client_secret = "%s"
    metadata_url  = "https://login.microsoftonline.com/${data.azurerm_client_config.current.tenant_id}/v2.0/.well-known/openid-configuration"
    scope         = ["openid", "profile", "email"]
  }

  application_accelerator_enabled = true
  application_live_view_enabled   = true
}
`, template, clientId, clientSecret)
}
