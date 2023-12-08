// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VirtualDesktopHostPoolRegistrationInfoResource struct{}

func TestAccVirtualDesktopHostPoolRegInfo_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_host_pool_registration_info", "test")
	r := VirtualDesktopHostPoolRegistrationInfoResource{}

	// Set the expiration times
	timeNow := time.Now().UTC()
	expirationTime := timeNow.AddDate(0, 0, 1).Format(time.RFC3339)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, expirationTime),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("token").Exists(),
			),
		},
	})
}

func TestAccVirtualDesktopHostPoolRegInfo_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_host_pool_registration_info", "test")
	r := VirtualDesktopHostPoolRegistrationInfoResource{}

	// Set the expiration times
	timeNow := time.Now().UTC()
	expirationTimeBasic := timeNow.AddDate(0, 0, 1).Format(time.RFC3339)
	expirationTimeComplete := timeNow.AddDate(0, 0, 2).Format(time.RFC3339)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, expirationTimeBasic),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("token").Exists(),
			),
		},
		{
			Config: r.complete(data, expirationTimeComplete),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("token").Exists(),
			),
		},
		{
			Config: r.basic(data, expirationTimeBasic),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("token").Exists(),
			),
		},
	})
}

func (VirtualDesktopHostPoolRegistrationInfoResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.HostPoolRegistrationInfoID(state.ID)
	if err != nil {
		return nil, err
	}

	hostPoolId := hostpool.NewHostPoolID(id.SubscriptionId, id.ResourceGroup, id.HostPoolName)

	resp, err := clients.DesktopVirtualization.HostPoolsClient.Get(ctx, hostPoolId)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", hostPoolId, err)
	}
	exists := false
	if model := resp.Model; model != nil {
		if info := model.Properties.RegistrationInfo; info != nil {
			exists = info.Token != nil && len(*info.Token) > 0
		}
	}
	return utils.Bool(exists), nil
}

func (VirtualDesktopHostPoolRegistrationInfoResource) basic(data acceptance.TestData, expirationDate string) string {
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
  expiration_date = "%s"
}


`, data.RandomInteger, data.Locations.Secondary, data.RandomString, expirationDate)
}

func (VirtualDesktopHostPoolRegistrationInfoResource) complete(data acceptance.TestData, expirationDate string) string {
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
  expiration_date = "%s"
}

`, data.RandomInteger, data.Locations.Secondary, data.RandomString, expirationDate)
}
