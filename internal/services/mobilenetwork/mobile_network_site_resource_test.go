package mobilenetwork_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/site"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MobileNetworkSiteResource struct{}

func TestAccMobileNetworkSite_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_site", "test")
	r := MobileNetworkSiteResource{}
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

func TestAccMobileNetworkSite_withNetworkFunctionIds(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_site", "test")
	r := MobileNetworkSiteResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withNetworkFunctionIds(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMobileNetworkSite_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_site", "test")
	r := MobileNetworkSiteResource{}
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

func TestAccMobileNetworkSite_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_site", "test")
	r := MobileNetworkSiteResource{}
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

func TestAccMobileNetworkSite_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_site", "test")
	r := MobileNetworkSiteResource{}
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

func (r MobileNetworkSiteResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := site.ParseSiteID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.MobileNetwork.SiteClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r MobileNetworkSiteResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}
resource "azurerm_mobile_network" "test" {
  name                = "acctest-mn-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[2]s"
  mobile_country_code = "001"
  mobile_network_code = "01"
}

`, data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkSiteResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_mobile_network_site" "test" {
  name              = "acctest-mns-%d"
  mobile_network_id = azurerm_mobile_network.test.id
  location          = "%s"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkSiteResource) withNetworkFunctionIds(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_mobile_network_packet_core_control_plane" "test" {
  name                = "acctest-mnpccp-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[2]s"
  sku                 = "EvaluationPackage"
  mobile_network_id   = azurerm_mobile_network.test.id

  control_plane_access_interface {
    name         = "default-interface"
    ipv4_address = "192.168.1.199"
    ipv4_gateway = "192.168.1.1"
    ipv4_subnet  = "192.168.1.0/25"
  }

  platform {
    type = "BaseVM"
  }
}

resource "azurerm_mobile_network_packet_core_data_plane" "test" {
  name                                        = "acctest-mnpcdp-%d"
  mobile_network_packet_core_control_plane_id = azurerm_mobile_network_packet_core_control_plane.test.id
  location                                    = "%s"

  user_plane_access_interface {
    name         = "default-interface"
    ipv4_address = "192.168.1.199"
    ipv4_gateway = "192.168.1.1"
    ipv4_subnet  = "192.168.1.0/25"
  }

}

resource "azurerm_mobile_network_site" "test" {
  name              = "acctest-mns-%d"
  mobile_network_id = azurerm_mobile_network.test.id
  location          = "%s"

  network_function_ids = [
    azurerm_mobile_network_packet_core_control_plane.id,
    azurerm_mobile_network_packet_core_data_plane.id
  ]
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkSiteResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_site" "import" {
  name              = azurerm_mobile_network_site.test.name
  mobile_network_id = azurerm_mobile_network.test.id
  location          = "%s"
}
`, config, data.Locations.Primary)
}

func (r MobileNetworkSiteResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_site" "test" {
  name              = "acctest-mns-%d"
  mobile_network_id = azurerm_mobile_network.test.id
  location          = "%s"

  tags = {
    key = "value"
  }

}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkSiteResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_site" "test" {
  name              = "acctest-mns-%d"
  mobile_network_id = azurerm_mobile_network.test.id
  location          = "%s"

  tags = {
    key = "value"
  }

}
`, template, data.RandomInteger, data.Locations.Primary)
}
