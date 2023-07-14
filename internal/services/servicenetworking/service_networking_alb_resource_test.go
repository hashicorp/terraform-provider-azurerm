package servicenetworking_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-05-01-preview/trafficcontrollerinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AlbResource struct{}

func (r AlbResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := trafficcontrollerinterface.ParseTrafficControllerID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("while parsing resource ID: %+v", err)
	}

	resp, err := clients.ServiceNetworking.ServiceNetworkingClient.TrafficControllerInterface.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("checking for presence of existing %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestAccServiceNetworkingALB_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_alb", "test")

	// for preview only, remove before merge
	data.Locations.Primary = "northeurope"
	r := AlbResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("configuration_endpoint.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceNetworkingALB_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_alb", "test")

	// for preview only, remove before merge
	data.Locations.Primary = "northeurope"
	r := AlbResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("configuration_endpoint.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceNetworkingALB_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_alb", "test")

	// for preview only, remove before merge
	data.Locations.Primary = "northeurope"
	r := AlbResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("configuration_endpoint.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("configuration_endpoint.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func (r AlbResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acct-sn-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r AlbResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

resource "azurerm_alb" "test" {
  name                = "acct-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r AlbResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

resource "azurerm_alb" "test" {
  name                = "acct-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
  tags = {
    key = "value"
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}
