package devcenter_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/devcenters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"testing"
)

type DevCenterResource struct{}

func TestAccDevCenter_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center", "test")
	r := DevCenterResource{}

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

func (r DevCenterResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := devcenters.ParseDevCenterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.DevCenter.DevCenterClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (r DevCenterResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_dev_center" "test" {
  name = "testdevcenter-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  identity {
	type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func (r DevCenterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRg-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
