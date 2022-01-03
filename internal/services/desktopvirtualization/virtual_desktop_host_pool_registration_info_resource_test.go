package desktopvirtualization_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VirtualDesktopHostPoolRegistrationInfoResource struct {
}

func TestAccVirtualDesktopHostPoolRegInfo_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_host_pool_registration_info", "test")
	r := VirtualDesktopHostPoolRegistrationInfoResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("token").Exists(),
			),
			// a non-empty plan is expected as the expiration_date value is relative to execution so a continual change is expected in this case
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccVirtualDesktopHostPoolRegInfo_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_host_pool_registration_info", "test")
	r := VirtualDesktopHostPoolRegistrationInfoResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("token").Exists(),
			),
			// a non-empty plan is expected as the expiration_date value is relative to execution so a continual change is expected in this case
			ExpectNonEmptyPlan: true,
		},
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("token").Exists(),
			),
			// a non-empty plan is expected as the expiration_date value is relative to execution so a continual change is expected in this case
			ExpectNonEmptyPlan: true,
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("token").Exists(),
			),
			// a non-empty plan is expected as the expiration_date value is relative to execution so a continual change is expected in this case
			ExpectNonEmptyPlan: true,
		},
	})
}

func (VirtualDesktopHostPoolRegistrationInfoResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {

	id, err := parse.HostPoolRegistrationInfoID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DesktopVirtualization.HostPoolsClient.Get(ctx, id.ResourceGroup, id.HostPoolName)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}
	exists := resp.ID != nil && resp.HostPoolProperties != nil && resp.HostPoolProperties.RegistrationInfo != nil && len(*resp.HostPoolProperties.RegistrationInfo.Token) > 0

	return utils.Bool(exists), nil
}

func (VirtualDesktopHostPoolRegistrationInfoResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vdesktop-%d"
  location = "%s"
}

resource "azurerm_virtual_desktop_host_pool" "test" {
  name                 = "acctestHP%s"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  type                 = "Pooled"
  validate_environment = true
  load_balancer_type   = "BreadthFirst"

}

resource "azurerm_virtual_desktop_host_pool_registration_info" "test" {
  hostpool_id     = azurerm_virtual_desktop_host_pool.test.id
  expiration_date = timeadd(timestamp(), "48h")
}


`, data.RandomInteger, data.Locations.Secondary, data.RandomString)
}

func (VirtualDesktopHostPoolRegistrationInfoResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vdesktop-%d"
  location = "%s"
}

resource "azurerm_virtual_desktop_host_pool" "test" {
  name                 = "acctestHP%s"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  type                 = "Pooled"
  validate_environment = true
  load_balancer_type   = "BreadthFirst"

}

resource "azurerm_virtual_desktop_host_pool_registration_info" "test" {
  hostpool_id     = azurerm_virtual_desktop_host_pool.test.id
  expiration_date = timeadd(timestamp(), "72h")
}

`, data.RandomInteger, data.Locations.Secondary, data.RandomString)
}
