// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iotoperations_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iotoperations/parse"
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

func (IoTOperationsInstanceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_iotoperations_instance" "test" {
  name                = "example-iotoperations-instance"
  location            = "West US"
  resource_group_name = "%s"

  sku {
    name     = "S1"
    capacity = 1
  }
}
`, data.ResourceGroupName)
}

func (IoTOperationsInstanceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_iotoperations_instance" "test" {
  name                = "example-iotoperations-instance"
  location            = "West US"
  resource_group_name = "%s"

  sku {
    name     = "S1"
    capacity = 1
  }
}
`, data.ResourceGroupName)
}

func (IoTOperationsInstanceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_iotoperations_instance" "test" {
  name                = "example-iotoperations-instance"
  location            = "West US"
  resource_group_name = "%s"

  sku {
    name     = "S1"
    capacity = 1
  }

  tags = {
    environment = "testing"
    cost_center = "finance"
  }
}
`, data.ResourceGroupName)
}
