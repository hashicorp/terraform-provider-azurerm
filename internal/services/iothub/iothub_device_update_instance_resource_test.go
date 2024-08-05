// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceupdate/2022-10-01/deviceupdates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type IotHubDeviceUpdateInstanceResource struct{}

func TestAccIotHubDeviceUpdateInstance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_device_update_instance", "test")
	r := IotHubDeviceUpdateInstanceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("diagnostic_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubDeviceUpdateInstance_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_device_update_instance", "test")
	r := IotHubDeviceUpdateInstanceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("diagnostic_storage_account.0.connection_string"),
	})
}

func TestAccIotHubDeviceUpdateInstance_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_device_update_instance", "test")
	r := IotHubDeviceUpdateInstanceResource{}
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

func TestAccIotHubDeviceUpdateInstance_diagnosticStorageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_device_update_instance", "test")
	r := IotHubDeviceUpdateInstanceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.diagnosticStorageAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("diagnostic_storage_account.0.connection_string"),
		{
			Config: r.diagnosticStorageAccountUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("diagnostic_storage_account.0.connection_string"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubDeviceUpdateInstance_diagnosticEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_device_update_instance", "test")
	r := IotHubDeviceUpdateInstanceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.diagnosticEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.diagnosticDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubDeviceUpdateInstance_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_device_update_instance", "test")
	r := IotHubDeviceUpdateInstanceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.tagsUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r IotHubDeviceUpdateInstanceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := deviceupdates.ParseInstanceID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.IoTHub.DeviceUpdatesClient
	resp, err := client.InstancesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r IotHubDeviceUpdateInstanceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[2]d"
  location = "%[1]s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }
}

resource "azurerm_iothub_device_update_account" "test" {
  name                = "acc-dua-%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func (r IotHubDeviceUpdateInstanceResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_device_update_instance" "test" {
  name                     = "acc-dui-%s"
  device_update_account_id = azurerm_iothub_device_update_account.test.id
  iothub_id                = azurerm_iothub.test.id
}
`, template, data.RandomString)
}

func (r IotHubDeviceUpdateInstanceResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_iothub_device_update_instance" "test" {
  name                     = "acc-dui-%[2]s"
  device_update_account_id = azurerm_iothub_device_update_account.test.id
  iothub_id                = azurerm_iothub.test.id
  diagnostic_enabled       = true

  diagnostic_storage_account {
    connection_string = azurerm_storage_account.test.primary_connection_string
    id                = azurerm_storage_account.test.id
  }

  tags = {
    environment = "AccTest"
  }
}
`, template, data.RandomString)
}

func (r IotHubDeviceUpdateInstanceResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_device_update_instance" "import" {
  name                     = azurerm_iothub_device_update_instance.test.name
  device_update_account_id = azurerm_iothub_device_update_instance.test.device_update_account_id
  iothub_id                = azurerm_iothub_device_update_instance.test.iothub_id
}
`, config)
}

func (r IotHubDeviceUpdateInstanceResource) diagnosticStorageAccount(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_iothub_device_update_instance" "test" {
  name                     = "acc-dui-%[2]s"
  device_update_account_id = azurerm_iothub_device_update_account.test.id
  iothub_id                = azurerm_iothub.test.id


  diagnostic_storage_account {
    connection_string = azurerm_storage_account.test.primary_connection_string
    id                = azurerm_storage_account.test.id
  }
}
`, template, data.RandomString)
}

func (r IotHubDeviceUpdateInstanceResource) diagnosticStorageAccountUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctestsa2%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_iothub_device_update_instance" "test" {
  name                     = "acc-dui-%[2]s"
  device_update_account_id = azurerm_iothub_device_update_account.test.id
  iothub_id                = azurerm_iothub.test.id

  diagnostic_storage_account {
    connection_string = azurerm_storage_account.test2.primary_connection_string
    id                = azurerm_storage_account.test2.id
  }
}
`, template, data.RandomString)
}

func (r IotHubDeviceUpdateInstanceResource) diagnosticEnabled(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_device_update_instance" "test" {
  name                     = "acc-dui-%s"
  device_update_account_id = azurerm_iothub_device_update_account.test.id
  iothub_id                = azurerm_iothub.test.id

  diagnostic_enabled = true
}
`, template, data.RandomString)
}

func (r IotHubDeviceUpdateInstanceResource) diagnosticDisabled(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_device_update_instance" "test" {
  name                     = "acc-dui-%s"
  device_update_account_id = azurerm_iothub_device_update_account.test.id
  iothub_id                = azurerm_iothub.test.id

  diagnostic_enabled = false
}
`, template, data.RandomString)
}

func (r IotHubDeviceUpdateInstanceResource) tags(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_device_update_instance" "test" {
  name                     = "acc-dui-%s"
  device_update_account_id = azurerm_iothub_device_update_account.test.id
  iothub_id                = azurerm_iothub.test.id

  tags = {
    environment = "AccTest"
  }
}
`, template, data.RandomString)
}

func (r IotHubDeviceUpdateInstanceResource) tagsUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_device_update_instance" "test" {
  name                     = "acc-dui-%s"
  device_update_account_id = azurerm_iothub_device_update_account.test.id
  iothub_id                = azurerm_iothub.test.id

  tags = {
    environment = "AccTest2"
    purpose     = "Testing"
  }
}
`, template, data.RandomString)
}
