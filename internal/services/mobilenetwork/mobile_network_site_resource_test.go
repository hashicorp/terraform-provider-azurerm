package mobilenetwork_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/site"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MobileNetworkSiteResource struct{}

func TestAccMobileNetworkSite_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_site", "test")

	// Limited regional availability for Mobile Network
	data.Locations.Primary = "eastus"

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

func TestAccMobileNetworkSite_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_site", "test")

	// Limited regional availability for Mobile Network
	data.Locations.Primary = "eastus"

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

	// Limited regional availability for Mobile Network
	data.Locations.Primary = "eastus"

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

	// Limited regional availability for Mobile Network
	data.Locations.Primary = "eastus"

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

func (r MobileNetworkSiteResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
				%s

resource "azurerm_mobile_network_site" "test" {
  name              = "acctest-mns-%d"
  mobile_network_id = azurerm_mobile_network.test.id
  location          = "%s"
}
`, MobileNetworkResource{}.basic(data), data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkSiteResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_site" "import" {
  name              = azurerm_mobile_network_site.test.name
  mobile_network_id = azurerm_mobile_network.test.id
  location          = "%s"
}
`, r.basic(data), data.Locations.Primary)
}

func (r MobileNetworkSiteResource) complete(data acceptance.TestData) string {
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
`, MobileNetworkResource{}.basic(data), data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkSiteResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_site" "test" {
  name              = "acctest-mns-%d"
  mobile_network_id = azurerm_mobile_network.test.id
  location          = "%s"

  tags = {
    key = "update"
  }

}
`, MobileNetworkResource{}.basic(data), data.RandomInteger, data.Locations.Primary)
}
