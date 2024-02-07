// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type IotHubDPSResource struct{}

func TestAccIotHubDPS_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_dps", "test")
	r := IotHubDPSResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allocation_policy").Exists(),
				check.That(data.ResourceName).Key("device_provisioning_host_name").Exists(),
				check.That(data.ResourceName).Key("id_scope").Exists(),
				check.That(data.ResourceName).Key("service_operations_host_name").Exists(),
				check.That(data.ResourceName).Key("data_residency_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubDPS_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_dps", "test")
	r := IotHubDPSResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubDPS_dataResidencyEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_dps", "test")
	r := IotHubDPSResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataResidencyEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubDPS_linkedHubs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_dps", "test")
	r := IotHubDPSResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linkedHubs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allocation_policy").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.linkedHubsUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubDPS_ipFilterRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_dps", "test")
	r := IotHubDPSResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipFilterRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allocation_policy").Exists(),
				check.That(data.ResourceName).Key("device_provisioning_host_name").Exists(),
				check.That(data.ResourceName).Key("id_scope").Exists(),
				check.That(data.ResourceName).Key("service_operations_host_name").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("ip_filter_rule.0.name").HasValue("test"),
				check.That(data.ResourceName).Key("ip_filter_rule.0.ip_mask").HasValue("10.0.0.0/31"),
				check.That(data.ResourceName).Key("ip_filter_rule.0.action").HasValue("Accept"),
				check.That(data.ResourceName).Key("ip_filter_rule.0.target").HasValue("all"),
				check.That(data.ResourceName).Key("ip_filter_rule.1.name").HasValue("test2"),
				check.That(data.ResourceName).Key("ip_filter_rule.1.ip_mask").HasValue("10.0.2.0/31"),
				check.That(data.ResourceName).Key("ip_filter_rule.1.action").HasValue("Accept"),
				check.That(data.ResourceName).Key("ip_filter_rule.1.target").HasValue("serviceApi"),
			),
		},
		{
			Config: r.ipFilterRulesUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_filter_rule.0.name").HasValue("test"),
				check.That(data.ResourceName).Key("ip_filter_rule.0.ip_mask").HasValue("10.0.0.0/31"),
				check.That(data.ResourceName).Key("ip_filter_rule.0.action").HasValue("Reject"),
				check.That(data.ResourceName).Key("ip_filter_rule.0.target").HasValue("all"),
				check.That(data.ResourceName).Key("ip_filter_rule.1.name").HasValue("test2"),
				check.That(data.ResourceName).Key("ip_filter_rule.1.ip_mask").HasValue("10.0.2.0/31"),
				check.That(data.ResourceName).Key("ip_filter_rule.1.action").HasValue("Reject"),
				check.That(data.ResourceName).Key("ip_filter_rule.1.target").HasValue("deviceApi"),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func (t IotHubDPSResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseProvisioningServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.IoTHub.DPSResourceClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
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
  name                          = "acctestIoTDPS-%d"
  resource_group_name           = azurerm_resource_group.test.name
  location                      = azurerm_resource_group.test.location
  public_network_access_enabled = false

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

func (IotHubDPSResource) dataResidencyEnabled(data acceptance.TestData) string {
	// Data residency has limited region support
	data.Locations.Primary = "brazilsouth"
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

  data_residency_enabled = true

  sku {
    name     = "S1"
    capacity = "1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r IotHubDPSResource) linkedHubs(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_iothub_dps" "test" {
  name                = "acctestIoTDPS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  linked_hub {
    connection_string       = azurerm_iothub_shared_access_policy.test.primary_connection_string
    location                = azurerm_resource_group.test.location
    allocation_weight       = 15
    apply_allocation_policy = true
  }

  linked_hub {
    connection_string = azurerm_iothub_shared_access_policy.test2.primary_connection_string
    location          = azurerm_resource_group.test.location
  }
}
`, r.linkedHubsDependencies(data), data.RandomInteger)
}

func (r IotHubDPSResource) linkedHubsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_iothub_dps" "test" {
  name                = "acctestIoTDPS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  linked_hub {
    connection_string       = azurerm_iothub_shared_access_policy.test.primary_connection_string
    location                = azurerm_resource_group.test.location
    allocation_weight       = 150
    apply_allocation_policy = true
  }
}
`, r.linkedHubsDependencies(data), data.RandomInteger)
}

func (IotHubDPSResource) linkedHubsDependencies(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_shared_access_policy" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_name         = azurerm_iothub.test.name
  name                = "acctest"

  registry_read   = true
  registry_write  = true
  service_connect = true
}

resource "azurerm_iothub" "test2" {
  name                = "acctestIoTHub2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_shared_access_policy" "test2" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_name         = azurerm_iothub.test2.name
  name                = "acctest2"

  registry_read   = true
  registry_write  = true
  service_connect = true
}
`, data.Locations.Primary, data.RandomInteger)
}

func (IotHubDPSResource) ipFilterRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub_dps" "test" {
  name                          = "acctestIoTDPS-%d"
  resource_group_name           = azurerm_resource_group.test.name
  location                      = azurerm_resource_group.test.location
  public_network_access_enabled = false

  ip_filter_rule {
    name    = "test"
    ip_mask = "10.0.0.0/31"
    action  = "Accept"
    target  = "all"
  }

  ip_filter_rule {
    name    = "test2"
    ip_mask = "10.0.2.0/31"
    action  = "Accept"
    target  = "serviceApi"
  }

  ip_filter_rule {
    name    = "test3"
    ip_mask = "10.0.3.0/31"
    action  = "Accept"
  }

  sku {
    name     = "S1"
    capacity = "1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IotHubDPSResource) ipFilterRulesUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub_dps" "test" {
  name                          = "acctestIoTDPS-%d"
  resource_group_name           = azurerm_resource_group.test.name
  location                      = azurerm_resource_group.test.location
  public_network_access_enabled = false

  ip_filter_rule {
    name    = "test"
    ip_mask = "10.0.0.0/31"
    action  = "Reject"
    target  = "all"
  }

  ip_filter_rule {
    name    = "test2"
    ip_mask = "10.0.2.0/31"
    action  = "Reject"
    target  = "deviceApi"
  }

  ip_filter_rule {
    name    = "test3"
    ip_mask = "10.0.3.0/31"
    action  = "Accept"
  }

  sku {
    name     = "S1"
    capacity = "1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
