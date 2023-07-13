// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ApiManagementProductTagResource struct{}

func TestAccApiManagementProductTag_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product_tag", "test")
	r := ApiManagementProductTagResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("api_management_product_id").Exists(),
				check.That(data.ResourceName).Key("api_management_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementProductTag_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product_tag", "test")
	r := ApiManagementProductTagResource{}

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

func (ApiManagementProductTagResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ProductTagID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.TagClient.GetByProduct(ctx, id.ResourceGroup, id.ServiceName, id.ProductName, id.TagName)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ApiManagementProductTagResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Consumption_0"
}

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product"
  api_management_name   = azurerm_api_management.test.name
  resource_group_name   = azurerm_resource_group.test.name
  display_name          = "Test Product"
  subscription_required = false
  published             = false
}

resource "azurerm_api_management_tag" "test" {
  api_management_id = azurerm_api_management.test.id
  name              = "acctestTag-%d"
}

resource "azurerm_api_management_product_tag" "test" {
  api_management_product_id = azurerm_api_management_product.test.product_id
  api_management_name       = azurerm_api_management.test.name
  resource_group_name       = azurerm_resource_group.test.name
  name                      = azurerm_api_management_tag.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementProductTagResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_product_tag" "import" {
  api_management_product_id = azurerm_api_management_product_tag.test.api_management_product_id
  api_management_name       = azurerm_api_management_product_tag.test.api_management_name
  resource_group_name       = azurerm_api_management_product_tag.test.resource_group_name
  name                      = azurerm_api_management_product_tag.test.name
}
`, r.basic(data))
}
