// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iotoperations_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/instance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type IoTOperationsInstanceResource struct{}

func TestAccIoTOperationsInstance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_instance", "test")
	r := IoTOperationsInstanceResource{}

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

func TestAccIoTOperationsInstance_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_instance", "test")
	r := IoTOperationsInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_iotoperations_instance"),
		},
	})
}

func TestAccIoTOperationsInstance_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_instance", "test")
	r := IoTOperationsInstanceResource{}

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

func TestAccIoTOperationsInstance_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_instance", "test")
	r := IoTOperationsInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (r IoTOperationsInstanceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := instance.ParseInstanceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.IoTOperations.InstanceClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (IoTOperationsInstanceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iotops-%d"
  location = "%s"
}

resource "azurerm_iotoperations_instance" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location           = azurerm_resource_group.test.location
  schema_registry_ref = "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-rg/providers/Microsoft.DeviceRegistry/schemaRegistries/example-registry"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IoTOperationsInstanceResource) requiresImport(data acceptance.TestData) string {
	template := IoTOperationsInstanceResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_instance" "import" {
  name                = azurerm_iotoperations_instance.test.name
  resource_group_name = azurerm_iotoperations_instance.test.resource_group_name
  location           = azurerm_iotoperations_instance.test.location
  schema_registry_ref = azurerm_iotoperations_instance.test.schema_registry_ref
}
`, template)
}

func (IoTOperationsInstanceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iotops-%d"
  location = "%s"
}

resource "azurerm_iotoperations_instance" "test" {
  name                     = "acctest-%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                = azurerm_resource_group.test.location
  schema_registry_ref     = "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-rg/providers/Microsoft.DeviceRegistry/schemaRegistries/example-registry"
  description             = "This is a test IoT Operations instance for terraform acceptance test"
  version                 = "1.0.0"
  extended_location_name  = "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-rg/providers/Microsoft.ExtendedLocation/customLocations/example-location"
  extended_location_type  = "CustomLocation"

  tags = {
    environment = "testing"
    cost_center = "finance"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
