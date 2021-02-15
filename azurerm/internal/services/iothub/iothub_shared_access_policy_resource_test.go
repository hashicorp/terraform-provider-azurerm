package iothub_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type IoTHubSharedAccessPolicyResource struct {
}

func TestAccIotHubSharedAccessPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_shared_access_policy", "test")
	r := IoTHubSharedAccessPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("registry_read").HasValue("true"),
				check.That(data.ResourceName).Key("registry_write").HasValue("true"),
				check.That(data.ResourceName).Key("service_connect").HasValue("false"),
				check.That(data.ResourceName).Key("device_connect").HasValue("false"),
				check.That(data.ResourceName).Key("name").HasValue("acctest"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubSharedAccessPolicy_writeWithoutRead(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_shared_access_policy", "test")
	r := IoTHubSharedAccessPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.writeWithoutRead(data),
			ExpectError: regexp.MustCompile("If `registry_write` is set to true, `registry_read` must also be set to true"),
		},
	})
}

func TestAccIotHubSharedAccessPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_shared_access_policy", "test")
	r := IoTHubSharedAccessPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_iothub_shared_access_policy"),
		},
	})
}

func (IoTHubSharedAccessPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
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

  registry_read  = true
  registry_write = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r IoTHubSharedAccessPolicyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_shared_access_policy" "import" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_name         = azurerm_iothub.test.name
  name                = "acctest"

  registry_read  = true
  registry_write = true
}
`, r.basic(data))
}

func (IoTHubSharedAccessPolicyResource) writeWithoutRead(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
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

  registry_write = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (t IoTHubSharedAccessPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	iothubName := id.Path["IotHubs"]
	keyName := id.Path["IotHubKeys"]

	accessPolicy, err := clients.IoTHub.ResourceClient.GetKeysForKeyName(ctx, resourceGroup, iothubName, keyName)
	if err != nil {
		return nil, fmt.Errorf("loading IotHub Shared Access Policy %q: %+v", id, err)
	}

	return utils.Bool(accessPolicy.PrimaryKey != nil), nil
}
