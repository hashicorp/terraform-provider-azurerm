package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ApiManagementPolicyResource struct {
}

func TestAccApiManagementPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_policy", "test")
	r := ApiManagementPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("xml_link"),
	})
}

func TestAccApiManagementPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_policy", "test")
	r := ApiManagementPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("xml_link"),
		{
			Config: r.customPolicy(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("xml_link"),
	})
}

func TestAccApiManagementPolicy_customPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_policy", "test")
	r := ApiManagementPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.customPolicy(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("xml_link"),
	})
}

func (t ApiManagementPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.PolicyID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName

	resp, err := clients.ApiManagement.ServiceClient.Get(ctx, resourceGroup, serviceName)
	if err != nil {
		return nil, fmt.Errorf("reading ApiManagement Policy (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ApiManagementPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
}

resource "azurerm_api_management_policy" "test" {
  api_management_id = azurerm_api_management.test.id
  xml_link          = "https://raw.githubusercontent.com/terraform-providers/terraform-provider-azurerm/master/azurerm/internal/services/apimanagement/tests/testdata/api_management_policy_test.xml"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (ApiManagementPolicyResource) customPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
}

resource "azurerm_api_management_policy" "test" {
  api_management_id = azurerm_api_management.test.id

  xml_content = <<XML
<policies>
  <inbound>
    <set-variable name="abc" value="@(context.Request.Headers.GetValueOrDefault("X-Header-Name", ""))" />
    <find-and-replace from="xyz" to="foo" />
  </inbound>
</policies>
XML

}
`, data.RandomInteger, data.Locations.Primary)
}
