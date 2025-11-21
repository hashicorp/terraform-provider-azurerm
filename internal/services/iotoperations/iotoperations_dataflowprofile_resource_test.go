// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iotoperations_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/dataflowprofile"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// IotOperationsDataflowProfileResource is a test harness for azurerm_iotoperations_dataflow_profile acceptance tests.
type IotOperationsDataflowProfileResource struct{}

func TestAccIotOperationsDataflowProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_dataflow_profile", "test")
	r := IotOperationsDataflowProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctest-dfp-%s", data.RandomString)),
				check.That(data.ResourceName).Key("properties.0.instance_count").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotOperationsDataflowProfile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_dataflow_profile", "test")
	r := IotOperationsDataflowProfileResource{}

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

func TestAccIotOperationsDataflowProfile_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_dataflow_profile", "test")
	r := IotOperationsDataflowProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctest-dfp-%s", data.RandomString)),
				check.That(data.ResourceName).Key("properties.0.instance_count").HasValue("3"),
				check.That(data.ResourceName).Key("properties.0.diagnostics.0.logs.0.level").HasValue("Info"),
				check.That(data.ResourceName).Key("properties.0.diagnostics.0.metrics.0.prometheus_port").HasValue("9090"),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotOperationsDataflowProfile_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_dataflow_profile", "test")
	r := IotOperationsDataflowProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("properties.0.instance_count").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("properties.0.instance_count").HasValue("3"),
				check.That(data.ResourceName).Key("properties.0.diagnostics.0.logs.0.level").HasValue("Info"),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("properties.0.instance_count").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

// Exists implements the acceptance existence check using the generated SDK ID parser.
func (IotOperationsDataflowProfileResource) Exists(ctx context.Context, c *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := dataflowprofile.ParseDataflowProfileID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing ID %q: %w", state.ID, err)
	}

	resp, err := c.IoTOperations.DataflowProfileClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id.ID(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r IotOperationsDataflowProfileResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iotops-%d"
  location = "%s"
}

resource "azurerm_iotoperations_instance" "test" {
  name                = "acctest-instance-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  
  extended_location {
    name = "acctest-custom-location-%s"
    type = "CustomLocation"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func (r IotOperationsDataflowProfileResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_dataflow_profile" "test" {
  name                = "acctest-dfp-%s"
  resource_group_name = azurerm_resource_group.test.name
  instance_name       = azurerm_iotoperations_instance.test.name
  location            = azurerm_resource_group.test.location

  properties {
    instance_count = 1
  }

  extended_location {
    name = azurerm_iotoperations_instance.test.extended_location[0].name
    type = azurerm_iotoperations_instance.test.extended_location[0].type
  }
}
`, r.template(data), data.RandomString)
}

func (r IotOperationsDataflowProfileResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_dataflow_profile" "import" {
  name                = azurerm_iotoperations_dataflow_profile.test.name
  resource_group_name = azurerm_iotoperations_dataflow_profile.test.resource_group_name
  instance_name       = azurerm_iotoperations_dataflow_profile.test.instance_name
  location            = azurerm_iotoperations_dataflow_profile.test.location

  properties {
    instance_count = azurerm_iotoperations_dataflow_profile.test.properties[0].instance_count
  }

  extended_location {
    name = azurerm_iotoperations_dataflow_profile.test.extended_location[0].name
    type = azurerm_iotoperations_dataflow_profile.test.extended_location[0].type
  }
}
`, r.basic(data))
}

func (r IotOperationsDataflowProfileResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_dataflow_profile" "test" {
  name                = "acctest-dfp-%s"
  resource_group_name = azurerm_resource_group.test.name
  instance_name       = azurerm_iotoperations_instance.test.name
  location            = azurerm_resource_group.test.location

  properties {
    instance_count = 3

    diagnostics {
      logs {
        level = "Info"
      }
      metrics {
        prometheus_port = 9090
      }
    }
  }

  extended_location {
    name = azurerm_iotoperations_instance.test.extended_location[0].name
    type = azurerm_iotoperations_instance.test.extended_location[0].type
  }

  tags = {
    environment = "testing"
    purpose     = "dataflow-profile-acceptance-test"
  }
}
`, r.template(data), data.RandomString)
}
