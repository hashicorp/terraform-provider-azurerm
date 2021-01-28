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

func TestDeviceSecurityGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_security_group", "test")
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

func TestDeviceSecurityGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_security_group", "test")
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

func TestDeviceSecurityGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_security_group", "test")
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

func TestDeviceSecurityGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_device_security_group", "test")
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
	id, err := parse.DeviceSecurityGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SecurityCenter.DeviceSecurityGroupsClient.Get(ctx, id.TargetResourceID, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Device Security Group %q: %+v", id.ID(), err)
	}

	return utils.Bool(resp.DeviceSecurityGroupProperties != nil), nil
}

func (r DeviceSecurityGroupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_device_security_group" "test" {
  name               = "acctest-DSG-%d"
  target_resource_id = azurerm_iothub.test.id

  depends_on = [azurerm_iot_security_solution.test]
}
`, r.template(data), data.RandomInteger)
}

func (r DeviceSecurityGroupResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_device_security_group" "import" {
  name               = azurerm_device_security_group.test.name
  target_resource_id = azurerm_device_security_group.test.target_resource_id
}
`, r.basic(data))
}

func (r DeviceSecurityGroupResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_device_security_group" "test" {
  name               = "acctest-DSG-%d"
  target_resource_id = azurerm_iothub.test.id

  allow_list_rule {
    type   = "LocalUserNotAllowed"
    values = ["user1"]
  }

  allow_list_rule {
    type   = "ProcessNotAllowed"
    values = ["ssh"]
  }

  allow_list_rule {
    type   = "ConnectionToIpNotAllowed"
    values = ["10.0.0.0/24"]
  }

  time_window_rule {
    type             = "ActiveConnectionsNotInAllowedRange"
    min_threshold    = 0
    max_threshold    = 30
    time_window_size = "PT5M"
  }

  time_window_rule {
    type             = "AmqpC2DMessagesNotInAllowedRange"
    min_threshold    = 0
    max_threshold    = 30
    time_window_size = "PT5M"
  }

  time_window_rule {
    type             = "MqttC2DMessagesNotInAllowedRange"
    min_threshold    = 0
    max_threshold    = 30
    time_window_size = "PT5M"
  }

  time_window_rule {
    type             = "HttpC2DMessagesNotInAllowedRange"
    min_threshold    = 0
    max_threshold    = 30
    time_window_size = "PT5M"
  }

  time_window_rule {
    type             = "AmqpC2DRejectedMessagesNotInAllowedRange"
    min_threshold    = 0
    max_threshold    = 30
    time_window_size = "PT5M"
  }

  time_window_rule {
    type             = "MqttC2DRejectedMessagesNotInAllowedRange"
    min_threshold    = 0
    max_threshold    = 30
    time_window_size = "PT5M"
  }

  time_window_rule {
    type             = "HttpC2DRejectedMessagesNotInAllowedRange"
    min_threshold    = 0
    max_threshold    = 30
    time_window_size = "PT5M"
  }


  time_window_rule {
    type             = "AmqpD2CMessagesNotInAllowedRange"
    min_threshold    = 0
    max_threshold    = 30
    time_window_size = "PT5M"
  }

  time_window_rule {
    type             = "MqttD2CMessagesNotInAllowedRange"
    min_threshold    = 0
    max_threshold    = 30
    time_window_size = "PT5M"
  }

  time_window_rule {
    type             = "HttpD2CMessagesNotInAllowedRange"
    min_threshold    = 0
    max_threshold    = 30
    time_window_size = "PT5M"
  }

  time_window_rule {
    type             = "DirectMethodInvokesNotInAllowedRange"
    min_threshold    = 0
    max_threshold    = 30
    time_window_size = "PT5M"
  }

  time_window_rule {
    type             = "FailedLocalLoginsNotInAllowedRange"
    min_threshold    = 0
    max_threshold    = 30
    time_window_size = "PT5M"
  }

  time_window_rule {
    type             = "FileUploadsNotInAllowedRange"
    min_threshold    = 0
    max_threshold    = 30
    time_window_size = "PT5M"
  }

  time_window_rule {
    type             = "QueuePurgesNotInAllowedRange"
    min_threshold    = 0
    max_threshold    = 30
    time_window_size = "PT5M"
  }

  time_window_rule {
    type             = "TwinUpdatesNotInAllowedRange"
    min_threshold    = 0
    max_threshold    = 30
    time_window_size = "PT5M"
  }

  time_window_rule {
    type             = "UnauthorizedOperationsNotInAllowedRange"
    min_threshold    = 0
    max_threshold    = 30
    time_window_size = "PT5M"
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
