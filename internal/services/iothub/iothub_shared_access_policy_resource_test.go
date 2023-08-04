// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type IoTHubSharedAccessPolicyResource struct{}

func TestAccIotHubSharedAccessPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_shared_access_policy", "test")
	r := IoTHubSharedAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.writeWithoutRead(data),
			ExpectError: regexp.MustCompile("If `registry_write` is set to true, `registry_read` must also be set to true"),
		},
	})
}

func TestAccIotHubSharedAccessPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_shared_access_policy", "test")
	r := IoTHubSharedAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
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

func (t IoTHubSharedAccessPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SharedAccessPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	accessPolicy, err := clients.IoTHub.ResourceClient.GetKeysForKeyName(ctx, id.ResourceGroup, id.IotHubName, id.IotHubKeyName)
	if err != nil {
		return nil, fmt.Errorf("loading IotHub Shared Access Policy %q: %+v", id, err)
	}

	return utils.Bool(accessPolicy.PrimaryKey != nil), nil
}
