// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/datanetwork"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MobileNetworkDataNetworkResource struct{}

func TestAccMobileNetworkDataNetwork_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_data_network", "test")
	r := MobileNetworkDataNetworkResource{}
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

func TestAccMobileNetworkDataNetwork_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_data_network", "test")
	r := MobileNetworkDataNetworkResource{}
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

func TestAccMobileNetworkDataNetwork_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_data_network", "test")
	r := MobileNetworkDataNetworkResource{}
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

func TestAccMobileNetworkDataNetwork_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_data_network", "test")
	r := MobileNetworkDataNetworkResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r MobileNetworkDataNetworkResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := datanetwork.ParseDataNetworkID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.MobileNetwork.DataNetworkClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r MobileNetworkDataNetworkResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_mobile_network_data_network" "test" {
  name              = "acctest-mndn-%d"
  mobile_network_id = azurerm_mobile_network.test.id
  location          = azurerm_resource_group.test.location
}
`, template, data.RandomInteger)
}

func (r MobileNetworkDataNetworkResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mobile_network_data_network" "import" {
  name              = azurerm_mobile_network_data_network.test.name
  mobile_network_id = azurerm_mobile_network_data_network.test.mobile_network_id
  location          = azurerm_mobile_network_data_network.test.location
}
`, template)
}

func (r MobileNetworkDataNetworkResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_mobile_network_data_network" "test" {
  name              = "acctest-mndn-%d"
  mobile_network_id = azurerm_mobile_network.test.id
  location          = azurerm_resource_group.test.location
  description       = "my favourite data network"
  tags = {
    key = "value"
  }
}
`, template, data.RandomInteger)
}

func (r MobileNetworkDataNetworkResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_mobile_network_data_network" "test" {
  name              = "acctest-mndn-%d"
  mobile_network_id = azurerm_mobile_network.test.id
  location          = azurerm_resource_group.test.location
  description       = "my favourite data network 2"
  tags = {
    key = "updated"
  }

}
`, template, data.RandomInteger)
}

func (r MobileNetworkDataNetworkResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-mn-%[1]d"
  location = %[2]q
}

resource "azurerm_mobile_network" "test" {
  name                = "acctest-mn-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  mobile_country_code = "001"
  mobile_network_code = "01"
}
`, data.RandomInteger, data.Locations.Primary)
}
