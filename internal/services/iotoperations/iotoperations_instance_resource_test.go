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

// IotOperationsInstanceResource is a test harness for azurerm_iotoperations_instance acceptance tests.
type IotOperationsInstanceResource struct{}

func TestAccIotOperationsInstance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_instance", "test")
	r := IotOperationsInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctest-ioi-%s", data.RandomString)),
				check.That(data.ResourceName).Key("schema_registry_ref").Exists(),
				check.That(data.ResourceName).Key("extended_location.0.type").HasValue("CustomLocation"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotOperationsInstance_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_instance", "test")
	r := IotOperationsInstanceResource{}

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

func TestAccIotOperationsInstance_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_instance", "test")
	r := IotOperationsInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctest-ioi-%s", data.RandomString)),
				check.That(data.ResourceName).Key("schema_registry_ref").Exists(),
				check.That(data.ResourceName).Key("description").HasValue("This is a test IoT Operations instance for terraform acceptance test"),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotOperationsInstance_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_instance", "test")
	r := IotOperationsInstanceResource{}

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
				check.That(data.ResourceName).Key("description").HasValue("This is a test IoT Operations instance for terraform acceptance test"),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func (r IotOperationsInstanceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r IotOperationsInstanceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iotops-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r IotOperationsInstanceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_instance" "test" {
  name                = "acctest-ioi-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  schema_registry_ref = "/subscriptions/%s/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.DeviceRegistry/schemaRegistries/acctest-registry-%s"

  extended_location {
    name = "acctest-custom-location-%s"
    type = "CustomLocation"
  }
}
`, r.template(data), data.RandomString, data.Client().SubscriptionID, data.RandomString, data.RandomString)
}

func (r IotOperationsInstanceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_instance" "import" {
  name                = azurerm_iotoperations_instance.test.name
  resource_group_name = azurerm_iotoperations_instance.test.resource_group_name
  location            = azurerm_iotoperations_instance.test.location
  schema_registry_ref = azurerm_iotoperations_instance.test.schema_registry_ref

  extended_location {
    name = azurerm_iotoperations_instance.test.extended_location[0].name
    type = azurerm_iotoperations_instance.test.extended_location[0].type
  }
}
`, r.basic(data))
}

func (r IotOperationsInstanceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_instance" "test" {
  name                = "acctest-ioi-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  schema_registry_ref = "/subscriptions/%s/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.DeviceRegistry/schemaRegistries/acctest-registry-%s"
  description         = "This is a test IoT Operations instance for terraform acceptance test"

  extended_location {
    name = "acctest-custom-location-%s"
    type = "CustomLocation"
  }

  tags = {
    environment = "testing"
    cost_center = "finance"
  }
}
`, r.template(data), data.RandomString, data.Client().SubscriptionID, data.RandomString, data.RandomString)
}
