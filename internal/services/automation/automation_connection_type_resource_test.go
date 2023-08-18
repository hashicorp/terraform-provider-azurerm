// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/connectiontype"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutomationConnectionTypeResource struct{}

func (a AutomationConnectionTypeResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := connectiontype.ParseConnectionTypeID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Automation.ConnectionType.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving Automation Connection Type %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (a AutomationConnectionTypeResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%[1]d"
  location = "%[2]s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (a AutomationConnectionTypeResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_automation_connection_type" "test" {
  name                    = "acctest-%[2]d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  is_global               = false
  field {
    name = "my_def"
    type = "string"
  }
}
`, a.template(data), data.RandomInteger)
}

func TestAccAutomationConnectionType_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.AutomationConnectionTypeResource{}.ResourceType(), "test")
	r := AutomationConnectionTypeResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("is_global").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}
