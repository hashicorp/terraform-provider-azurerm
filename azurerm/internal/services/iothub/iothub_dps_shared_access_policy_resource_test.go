package iothub_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type IotHubDpsSharedAccessPolicyResource struct {
}

func TestAccIotHubDpsSharedAccessPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_dps_shared_access_policy", "test")
	r := IotHubDpsSharedAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctest"),
				check.That(data.ResourceName).Key("enrollment_read").HasValue("false"),
				check.That(data.ResourceName).Key("enrollment_write").HasValue("false"),
				check.That(data.ResourceName).Key("registration_read").HasValue("false"),
				check.That(data.ResourceName).Key("registration_write").HasValue("false"),
				check.That(data.ResourceName).Key("service_config").HasValue("true"),
				check.That(data.ResourceName).Key("primary_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_key").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
	})
}

func TestAccIotHubDpsSharedAccessPolicy_writeWithoutRead(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_dps_shared_access_policy", "test")
	r := IotHubDpsSharedAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.writeWithoutRead(data),
			ExpectError: regexp.MustCompile("If `registration_write` is set to true, `registration_read` must also be set to true"),
		},
	})
}

func TestAccIotHubDpsSharedAccessPolicy_enrollmentReadWithoutRegistration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_dps_shared_access_policy", "test")
	r := IotHubDpsSharedAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.enrollmentReadWithoutRegistration(data),
			ExpectError: regexp.MustCompile("If `enrollment_read` is set to true, `registration_read` must also be set to true"),
		},
	})
}

func TestAccIotHubDpsSharedAccessPolicy_enrollmentWriteWithoutOthers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_dps_shared_access_policy", "test")
	r := IotHubDpsSharedAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.enrollmentWriteWithoutOthers(data),
			ExpectError: regexp.MustCompile("If `enrollment_write` is set to true, `enrollment_read`, `registration_read`, and `registration_write` must also be set to true"),
		},
	})
}

func (IotHubDpsSharedAccessPolicyResource) basic(data acceptance.TestData) string {
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

resource "azurerm_iothub_dps_shared_access_policy" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_dps_name     = azurerm_iothub_dps.test.name
  name                = "acctest"
  service_config      = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IotHubDpsSharedAccessPolicyResource) writeWithoutRead(data acceptance.TestData) string {
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

resource "azurerm_iothub_dps_shared_access_policy" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_dps_name     = azurerm_iothub_dps.test.name
  name                = "acctest"
  registration_write  = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IotHubDpsSharedAccessPolicyResource) enrollmentReadWithoutRegistration(data acceptance.TestData) string {
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

resource "azurerm_iothub_dps_shared_access_policy" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_dps_name     = azurerm_iothub_dps.test.name
  name                = "acctest"
  enrollment_read     = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IotHubDpsSharedAccessPolicyResource) enrollmentWriteWithoutOthers(data acceptance.TestData) string {
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

resource "azurerm_iothub_dps_shared_access_policy" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_dps_name     = azurerm_iothub_dps.test.name
  name                = "acctest"
  enrollment_write    = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (t IotHubDpsSharedAccessPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	iothubDpsName := id.Path["provisioningServices"]
	keyName := id.Path["keys"]

	accessPolicy, err := clients.IoTHub.DPSResourceClient.ListKeysForKeyName(ctx, iothubDpsName, keyName, resourceGroup)
	if err != nil {
		return nil, fmt.Errorf("loading Shared Access Policy (%s): %+v", id, err)
	}

	return utils.Bool(accessPolicy.PrimaryKey != nil), nil
}
