// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/productpolicy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementProductPolicyResource struct{}

func TestAccApiManagementProductPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product_policy", "test")
	r := ApiManagementProductPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			ResourceName:            data.ResourceName,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"xml_link"},
		},
	})
}

func TestAccApiManagementProductPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product_policy", "test")
	r := ApiManagementProductPolicyResource{}

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

func TestAccApiManagementProductPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product_policy", "test")
	r := ApiManagementProductPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			ResourceName:            data.ResourceName,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"xml_link"},
		},
	})
}

func (ApiManagementProductPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := productpolicy.ParseProductID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.ProductPoliciesClient.Get(ctx, *id, productpolicy.GetOperationOptions{Format: pointer.To(productpolicy.PolicyExportFormatXml)})
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
}

func (ApiManagementProductPolicyResource) basic(data acceptance.TestData) string {
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
  sku_name            = "Consumption_0"
}

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product"
  api_management_name   = azurerm_api_management.test.name
  resource_group_name   = azurerm_resource_group.test.name
  display_name          = "Test Product"
  subscription_required = false
  published             = false
}

resource "azurerm_api_management_product_policy" "test" {
  product_id          = azurerm_api_management_product.test.product_id
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  xml_link            = "https://gist.githubusercontent.com/riordanp/ca22f8113afae0eb38cc12d718fd048d/raw/d6ac89a2f35a6881a7729f8cb4883179dc88eea1/example.xml"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ApiManagementProductPolicyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_product_policy" "import" {
  product_id          = azurerm_api_management_product_policy.test.product_id
  api_management_name = azurerm_api_management_product_policy.test.api_management_name
  resource_group_name = azurerm_api_management_product_policy.test.resource_group_name
  xml_link            = azurerm_api_management_product_policy.test.xml_link
}
`, r.basic(data))
}

func (ApiManagementProductPolicyResource) updated(data acceptance.TestData) string {
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
  sku_name            = "Consumption_0"
}

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product"
  api_management_name   = azurerm_api_management.test.name
  resource_group_name   = azurerm_resource_group.test.name
  display_name          = "Test Product"
  subscription_required = false
  published             = false
}

resource "azurerm_api_management_product_policy" "test" {
  product_id          = azurerm_api_management_product.test.product_id
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name

  xml_content = <<XML
<policies>
  <inbound>
    <set-variable name="abc" value="@(context.Request.Headers.GetValueOrDefault("X-Header-Name", ""))" />
    <find-and-replace from="xyz" to="abc" />
  </inbound>
</policies>
XML

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
