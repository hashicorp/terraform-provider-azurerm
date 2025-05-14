// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type EventGridPartnerConfigurationsTestResource struct{}

func TestAccEventGridPartnerConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_partner_configuration", "test")
	r := EventGridPartnerConfigurationsTestResource{}

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

func TestAccEventGridPartnerConfiguration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_partner_configuration", "test")
	r := EventGridPartnerConfigurationsTestResource{}

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

func TestAccEventGridPartnerConfiguration_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_partner_configuration", "test")
	r := EventGridPartnerConfigurationsTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridPartnerConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_partner_configuration", "test")
	r := EventGridPartnerConfigurationsTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (EventGridPartnerConfigurationsTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseResourceGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.EventGrid.PartnerConfigurations.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (EventGridPartnerConfigurationsTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventgrid_partner_configuration" "test" {
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r EventGridPartnerConfigurationsTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_eventgrid_partner_configuration" "import" {
  resource_group_name = azurerm_resource_group.test.name
}
`, r.basic(data))
}

func (EventGridPartnerConfigurationsTestResource) complete(data acceptance.TestData) string {
	expiryTime := time.Now().In(time.UTC).Add(7 * 24 * time.Hour).Format(time.RFC3339)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventgrid_partner_configuration" "test" {
  resource_group_name                     = azurerm_resource_group.test.name
  default_maximum_expiration_time_in_days = 180

  partner_authorization {
    partner_registration_id              = "804a11ca-ce9b-4158-8e94-3c8dc7a072ec"
    partner_name                         = "Auth0"
    authorization_expiration_time_in_utc = "%s"
  }

  tags = {
    "foo" = "bar"
  }
}
`, data.RandomInteger, data.Locations.Primary, expiryTime)
}
