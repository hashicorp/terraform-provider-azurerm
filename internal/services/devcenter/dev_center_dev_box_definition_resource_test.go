// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devboxdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DevCenterDevBoxDefinitionTestResource struct{}

func TestAccDevCenterDevBoxDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_dev_box_definition", "test")
	r := DevCenterDevBoxDefinitionTestResource{}

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

func TestAccDevCenterDevBoxDefinition_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_dev_box_definition", "test")
	r := DevCenterDevBoxDefinitionTestResource{}

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

func TestAccDevCenterDevBoxDefinition_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_dev_box_definition", "test")
	r := DevCenterDevBoxDefinitionTestResource{}

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

func TestAccDevCenterDevBoxDefinition_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_dev_box_definition", "test")
	r := DevCenterDevBoxDefinitionTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
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
	})
}

func (r DevCenterDevBoxDefinitionTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := devboxdefinitions.ParseDevCenterDevBoxDefinitionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DevCenter.V20250201.DevBoxDefinitions.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r DevCenterDevBoxDefinitionTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_dev_center_dev_box_definition" "test" {
  name               = "acctest-dcet-%d"
  location           = azurerm_resource_group.test.location
  dev_center_id      = azurerm_dev_center.test.id
  image_reference_id = "${azurerm_dev_center.test.id}/galleries/default/images/microsoftvisualstudio_visualstudioplustools_vs-2022-ent-general-win10-m365-gen2"
  sku_name           = "general_i_8c32gb256ssd_v2"
}
`, r.template(data), data.RandomInteger)
}

func (r DevCenterDevBoxDefinitionTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dev_center_dev_box_definition" "import" {
  name               = azurerm_dev_center_dev_box_definition.test.name
  location           = azurerm_dev_center_dev_box_definition.test.location
  dev_center_id      = azurerm_dev_center_dev_box_definition.test.dev_center_id
  image_reference_id = azurerm_dev_center_dev_box_definition.test.image_reference_id
  sku_name           = "general_i_8c32gb256ssd_v2"
}
`, r.basic(data))
}

func (r DevCenterDevBoxDefinitionTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_dev_center_dev_box_definition" "test" {
  name               = "acctest-dcet-%d"
  location           = azurerm_resource_group.test.location
  dev_center_id      = azurerm_dev_center.test.id
  image_reference_id = "${azurerm_dev_center.test.id}/galleries/default/images/microsoftvisualstudio_visualstudioplustools_vs-2022-ent-general-win10-m365-gen2"
  sku_name           = "general_i_8c32gb256ssd_v2"

  tags = {
    Env = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r DevCenterDevBoxDefinitionTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_dev_center_dev_box_definition" "test" {
  name               = "acctest-dcet-%d"
  location           = azurerm_resource_group.test.location
  dev_center_id      = azurerm_dev_center.test.id
  image_reference_id = "${azurerm_dev_center.test.id}/galleries/default/images/microsoftvisualstudio_visualstudioplustools_vs-2022-ent-general-win11-m365-gen2"
  sku_name           = "general_i_8c32gb512ssd_v2"

  tags = {
    Env = "Test2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r DevCenterDevBoxDefinitionTestResource) template(data acceptance.TestData) string {
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
