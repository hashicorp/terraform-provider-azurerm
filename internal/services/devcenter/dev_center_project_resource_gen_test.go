package devcenter_test

// NOTE: this file is generated - manual changes will be overwritten.
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.
import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/projects"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DevCenterProjectTestResource struct{}

func TestAccDevCenterProject_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_project", "test")
	r := DevCenterProjectTestResource{}

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

func TestAccDevCenterProject_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_project", "test")
	r := DevCenterProjectTestResource{}

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

func TestAccDevCenterProject_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_project", "test")
	r := DevCenterProjectTestResource{}

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

func TestAccDevCenterProject_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_project", "test")
	r := DevCenterProjectTestResource{}

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
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r DevCenterProjectTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := projects.ParseProjectID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DevCenter.V20230401.Projects.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r DevCenterProjectTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_dev_center_project" "test" {
  dev_center_id       = azurerm_dev_center.test.id
  location            = azurerm_resource_group.test.location
  name                = "acctestdcp-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
}
`, r.template(data))
}

func (r DevCenterProjectTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dev_center_project" "import" {
  dev_center_id       = azurerm_dev_center_project.test.dev_center_id
  location            = azurerm_dev_center_project.test.location
  name                = azurerm_dev_center_project.test.name
  resource_group_name = azurerm_dev_center_project.test.resource_group_name
}
`, r.basic(data))
}

func (r DevCenterProjectTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_dev_center_project" "test" {
  dev_center_id              = azurerm_dev_center.test.id
  location                   = azurerm_resource_group.test.location
  name                       = "acctestdcp-${var.random_string}"
  resource_group_name        = azurerm_resource_group.test.name
  description                = "Description for the Dev Center Project"
  maximum_dev_boxes_per_user = 21
  tags = {
    environment = "terraform-acctests"
    some_key    = "some-value"
  }
}
`, r.template(data))
}

func (r DevCenterProjectTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "primary_location" {
  default = %q
}
variable "random_integer" {
  default = %d
}
variable "random_string" {
  default = %q
}

resource "azurerm_dev_center" "test" {
  name                = "acctestdc-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}


resource "azurerm_resource_group" "test" {
  name     = "acctestrg-${var.random_integer}"
  location = var.primary_location
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
