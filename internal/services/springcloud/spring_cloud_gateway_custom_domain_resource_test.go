// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SpringCloudGatewayCustomDomainResource struct{}

func TestAccSpringCloudGatewayCustomDomain_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_gateway_custom_domain", "test")
	r := SpringCloudGatewayCustomDomainResource{}
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

func TestAccSpringCloudGatewayCustomDomain_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_gateway_custom_domain", "test")
	r := SpringCloudGatewayCustomDomainResource{}
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

func (r SpringCloudGatewayCustomDomainResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SpringCloudGatewayCustomDomainID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.AppPlatform.GatewayCustomDomainClient.Get(ctx, id.ResourceGroup, id.SpringName, id.GatewayName, id.DomainName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r SpringCloudGatewayCustomDomainResource) template(data acceptance.TestData) string {
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

func (r SpringCloudGatewayCustomDomainResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_gateway_custom_domain" "test" {
  name                    = "${azurerm_spring_cloud_service.test.name}.azuremicroservices.io"
  spring_cloud_gateway_id = azurerm_spring_cloud_gateway.test.id
}
`, template)
}

func (r SpringCloudGatewayCustomDomainResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_gateway_custom_domain" "import" {
  name                    = azurerm_spring_cloud_gateway_custom_domain.test.name
  spring_cloud_gateway_id = azurerm_spring_cloud_gateway_custom_domain.test.spring_cloud_gateway_id
}
`, config)
}
