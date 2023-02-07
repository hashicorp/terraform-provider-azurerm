package mobilenetwork_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/slice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MobileNetworkSliceResource struct{}

func TestAccMobileNetworkSlice_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_slice", "test")
	// Limited regional availability for Mobile Network
	data.Locations.Primary = "eastus"

	r := MobileNetworkSliceResource{}
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

func TestAccMobileNetworkSlice_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_slice", "test")
	// Limited regional availability for Mobile Network
	data.Locations.Primary = "eastus"

	r := MobileNetworkSliceResource{}
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

func TestAccMobileNetworkSlice_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_slice", "test")
	// Limited regional availability for Mobile Network
	data.Locations.Primary = "eastus"

	r := MobileNetworkSliceResource{}
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

func TestAccMobileNetworkSlice_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_slice", "test")
	// Limited regional availability for Mobile Network
	data.Locations.Primary = "eastus"

	r := MobileNetworkSliceResource{}
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

func (r MobileNetworkSliceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := slice.ParseSliceID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.MobileNetwork.SliceClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r MobileNetworkSliceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
				%s

resource "azurerm_mobile_network_slice" "test" {
  name              = "acctest-mns-%d"
  mobile_network_id = azurerm_mobile_network.test.id
  location          = "%s"
  single_network_slice_selection_assistance_information {
    slice_service_type = 1
  }
}
`, MobileNetworkResource{}.basic(data), data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkSliceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_slice" "import" {
  name              = azurerm_mobile_network_slice.test.name
  mobile_network_id = azurerm_mobile_network_slice.test.mobile_network_id

  location = "%s"
  single_network_slice_selection_assistance_information {
    slice_service_type = 1
  }
}
`, r.basic(data), data.Locations.Primary)
}

func (r MobileNetworkSliceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_slice" "test" {
  name              = "acctest-mns-%d"
  mobile_network_id = azurerm_mobile_network.test.id
  location          = "%s"
  description       = "my favorite slice"
  single_network_slice_selection_assistance_information {
    slice_service_type = 1
  }
  tags = {
    key = "value"
  }

}
`, MobileNetworkResource{}.basic(data), data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkSliceResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mobile_network_slice" "test" {
  name              = "acctest-mns-%d"
  mobile_network_id = azurerm_mobile_network.test.id
  location          = "%s"
  description       = "my favorite slice2"
  single_network_slice_selection_assistance_information {
    slice_service_type = 1
  }

  tags = {
    key = "value"
  }

}
`, MobileNetworkResource{}.basic(data), data.RandomInteger, data.Locations.Primary)
}
