package iothub_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type IotHubDPSResource struct {
}

func TestAccIotHubDPS_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_dps", "test")
	r := IotHubDPSResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allocation_policy").Exists(),
				check.That(data.ResourceName).Key("device_provisioning_host_name").Exists(),
				check.That(data.ResourceName).Key("id_scope").Exists(),
				check.That(data.ResourceName).Key("service_operations_host_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubDPS_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_dps", "test")
	r := IotHubDPSResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_iothub_dps"),
		},
	})
}

func TestAccIotHubDPS_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_dps", "test")
	r := IotHubDPSResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubDPS_linkedHubs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_dps", "test")
	r := IotHubDPSResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.linkedHubs(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.linkedHubsUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t IotHubDPSResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["provisioningServices"]
	// the name path can use the ProvisioningServices in older iterations
	if name == "" {
		name = id.Path["ProvisioningServices"]
	}

	resp, err := clients.IoTHub.DPSResourceClient.Get(ctx, name, resourceGroup)
	if err != nil {
		return nil, fmt.Errorf("reading IotHuB DPS (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (IotHubDPSResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub_dps" "test" {
  name                = "acctestIoTDPS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r IotHubDPSResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_dps" "import" {
  name                = azurerm_iothub_dps.test.name
  resource_group_name = azurerm_iothub_dps.test.resource_group_name
  location            = azurerm_iothub_dps.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }
}
`, r.basic(data))
}

func (IotHubDPSResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub_dps" "test" {
  name                = "acctestIoTDPS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IotHubDPSResource) linkedHubs(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub_dps" "test" {
  name                = "acctestIoTDPS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  linked_hub {
    connection_string       = "HostName=test.azure-devices.net;SharedAccessKeyName=iothubowner;SharedAccessKey=booo"
    location                = azurerm_resource_group.test.location
    allocation_weight       = 15
    apply_allocation_policy = true
  }

  linked_hub {
    connection_string = "HostName=test2.azure-devices.net;SharedAccessKeyName=iothubowner2;SharedAccessKey=key2"
    location          = azurerm_resource_group.test.location
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IotHubDPSResource) linkedHubsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub_dps" "test" {
  name                = "acctestIoTDPS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  linked_hub {
    connection_string = "HostName=test.azure-devices.net;SharedAccessKeyName=iothubowner;SharedAccessKey=booo"
    location          = azurerm_resource_group.test.location
    allocation_weight = 150
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
