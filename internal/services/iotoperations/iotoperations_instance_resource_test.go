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

	return pluginsdk.Bool(resp.Model != nil), nil
}

func (IoTOperationsInstanceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"
}

resource "azurerm_iotoperations_instance" "test" {
  instanceName      = "example-iotoperations-instance"
  resourceGroupName = azurerm_resource_group.test.name
  subscriptionId    = data.azurerm_client_config.current.subscription_id
  api-version       = "2024-11-01"
  location          = azurerm_resource_group.test.location
}

data "azurerm_client_config" "current" {}
`, data.ResourceGroupName, data.Locations.Primary)
}

func (IoTOperationsInstanceResource) requiresImport(data acceptance.TestData) string {
	template := IoTOperationsInstanceResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_instance" "import" {
  instanceName      = azurerm_iotoperations_instance.test.instanceName
  resourceGroupName = azurerm_iotoperations_instance.test.resourceGroupName
  subscriptionId    = azurerm_iotoperations_instance.test.subscriptionId
  api-version       = azurerm_iotoperations_instance.test.api-version
  location          = azurerm_iotoperations_instance.test.location
}
`, template)
}

func (IoTOperationsInstanceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"
}

resource "azurerm_iotoperations_instance" "test" {
  instanceName      = "example-iotoperations-instance"
  resourceGroupName = azurerm_resource_group.test.name
  subscriptionId    = data.azurerm_client_config.current.subscription_id
  api-version       = "2024-11-01"
  location          = azurerm_resource_group.test.location
  description       = "Test IoT Operations instance"
  version           = "1.0.0"

  tags = {
    environment = "testing"
    cost_center = "finance"
  }
}

data "azurerm_client_config" "current" {}
`, data.ResourceGroupName, data.Locations.Primary)
}
