package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DeviceSecurityGroupResource struct {
}

func TestIotSecurityDeviceGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_security_device_group", "test")
	r := DeviceSecurityGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestIotSecurityDeviceGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_security_device_group", "test")
	r := DeviceSecurityGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestIotSecurityDeviceGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_security_device_group", "test")
	r := DeviceSecurityGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestIotSecurityDeviceGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iot_security_device_group", "test")
	r := DeviceSecurityGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (DeviceSecurityGroupResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.IotSecurityDeviceGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SecurityCenter.DeviceSecurityGroupsClient.Get(ctx, id.IotHubID, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Iot Security Device Group %q: %+v", id.ID(), err)
	}

	return utils.Bool(resp.DeviceSecurityGroupProperties != nil), nil
}

func (r DeviceSecurityGroupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iot_security_device_group" "test" {
  name      = "acctest-ISDG-%d"
  iothub_id = azurerm_iothub.test.id

  depends_on = [azurerm_iot_security_solution.test]
}
`, r.template(data), data.RandomInteger)
}

func (r DeviceSecurityGroupResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iot_security_device_group" "import" {
  name      = azurerm_iot_security_device_group.test.name
  iothub_id = azurerm_iot_security_device_group.test.iothub_id
}
`, r.basic(data))
}

func (r DeviceSecurityGroupResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iot_security_device_group" "test" {
  name      = "acctest-ISDG-%d"
  iothub_id = azurerm_iothub.test.id

  allow_rule {
    connection_to_ip_not_allowed = ["10.0.0.0/24"]
    local_user_not_allowed       = ["user1"]
    process_not_allowed          = ["ssh"]
  }

  range_rule {
    type     = "ActiveConnectionsNotInAllowedRange"
    min      = 0
    max      = 30
    duration = "PT5M"
  }

  range_rule {
    type     = "AmqpC2DMessagesNotInAllowedRange"
    min      = 0
    max      = 30
    duration = "PT5M"
  }

  range_rule {
    type     = "MqttC2DMessagesNotInAllowedRange"
    min      = 0
    max      = 30
    duration = "PT5M"
  }

  range_rule {
    type     = "HttpC2DMessagesNotInAllowedRange"
    min      = 0
    max      = 30
    duration = "PT5M"
  }

  range_rule {
    type     = "AmqpC2DRejectedMessagesNotInAllowedRange"
    min      = 0
    max      = 30
    duration = "PT5M"
  }

  range_rule {
    type     = "MqttC2DRejectedMessagesNotInAllowedRange"
    min      = 0
    max      = 30
    duration = "PT5M"
  }

  range_rule {
    type     = "HttpC2DRejectedMessagesNotInAllowedRange"
    min      = 0
    max      = 30
    duration = "PT5M"
  }


  range_rule {
    type     = "AmqpD2CMessagesNotInAllowedRange"
    min      = 0
    max      = 30
    duration = "PT5M"
  }

  range_rule {
    type     = "MqttD2CMessagesNotInAllowedRange"
    min      = 0
    max      = 30
    duration = "PT5M"
  }

  range_rule {
    type     = "HttpD2CMessagesNotInAllowedRange"
    min      = 0
    max      = 30
    duration = "PT5M"
  }

  range_rule {
    type     = "DirectMethodInvokesNotInAllowedRange"
    min      = 0
    max      = 30
    duration = "PT5M"
  }

  range_rule {
    type     = "FailedLocalLoginsNotInAllowedRange"
    min      = 0
    max      = 30
    duration = "PT5M"
  }

  range_rule {
    type     = "FileUploadsNotInAllowedRange"
    min      = 0
    max      = 30
    duration = "PT5M"
  }

  range_rule {
    type     = "QueuePurgesNotInAllowedRange"
    min      = 0
    max      = 30
    duration = "PT5M"
  }

  range_rule {
    type     = "TwinUpdatesNotInAllowedRange"
    min      = 0
    max      = 30
    duration = "PT5M"
  }

  range_rule {
    type     = "UnauthorizedOperationsNotInAllowedRange"
    min      = 0
    max      = 30
    duration = "PT5M"
  }

  depends_on = [azurerm_iot_security_solution.test]
}
`, r.template(data), data.RandomInteger)
}

func (DeviceSecurityGroupResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-security-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku {
    name     = "S1"
    capacity = "1"
  }
}

resource "azurerm_iot_security_solution" "test" {
  name                = "acctest-Iot-Security-Solution-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  display_name        = "Iot Security Solution"
  iothub_ids          = [azurerm_iothub.test.id]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
