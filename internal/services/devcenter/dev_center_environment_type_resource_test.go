// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/environmenttypes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DevCenterEnvironmentTypeTestResource struct{}

func TestAccDevCenterEnvironmentType_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_environment_type", "test")
	r := DevCenterEnvironmentTypeTestResource{}

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

func TestAccDevCenterEnvironmentType_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_environment_type", "test")
	r := DevCenterEnvironmentTypeTestResource{}

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

func TestAccDevCenterEnvironmentType_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_environment_type", "test")
	r := DevCenterEnvironmentTypeTestResource{}

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

func TestAccDevCenterEnvironmentType_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_environment_type", "test")
	r := DevCenterEnvironmentTypeTestResource{}

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
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r DevCenterEnvironmentTypeTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := environmenttypes.ParseDevCenterEnvironmentTypeID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DevCenter.V20250201.EnvironmentTypes.EnvironmentTypesGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r DevCenterEnvironmentTypeTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_dev_center_environment_type" "test" {
  name          = "acctest-dcet-%d"
  dev_center_id = azurerm_dev_center.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r DevCenterEnvironmentTypeTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dev_center_environment_type" "import" {
  name          = azurerm_dev_center_environment_type.test.name
  dev_center_id = azurerm_dev_center_environment_type.test.dev_center_id
}
`, r.basic(data))
}

func (r DevCenterEnvironmentTypeTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_dev_center_environment_type" "test" {
  name          = "acctest-dcet-%d"
  dev_center_id = azurerm_dev_center.test.id

  tags = {
    Env = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r DevCenterEnvironmentTypeTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_dev_center_environment_type" "test" {
  name          = "acctest-dcet-%d"
  dev_center_id = azurerm_dev_center.test.id

  tags = {
    Env = "Test2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r DevCenterEnvironmentTypeTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-dcet-%d"
  location = "%s"
}

resource "azurerm_dev_center" "test" {
  name                = "acctest-dc-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
