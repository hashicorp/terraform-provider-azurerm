// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devcenters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DevCenterTestResource struct{}

func TestAccDevCenter(t *testing.T) {
	// Sequential testing because we can only provision a couple Dev Centers at a time.
	testCases := map[string]map[string]func(t *testing.T){
		"DevCenter": {
			"basic":          testAccDevCenter_basic,
			"requiresImport": testAccDevCenter_requiresImport,
			"update":         testAccDevCenter_update,
			"complete":       testAccDevCenter_complete,

			"basicDataSource": testAccDevCenterDataSource_basic,
		},

		"DevCenterAttachedNetwork": {
			"basic":          testAccDevCenterAttachedNetwork_basic,
			"requiresImport": testAccDevCenterAttachedNetwork_requiresImport,

			"basicDataSource": testAccDevCenterAttachedNetworkDataSource_basic,
		},

		"DevCenterCatalog": {
			"basic":  testAccDevCenterCatalogs_basic,
			"adoGit": testAccDevCenterCatalogs_adoGit,
			"update": testAccDevCenterCatalogs_update,

			"basicDataSource": testAccDevCenterCatalogDataSource_basic,
		},

		"DevCenterDevBoxDefinition": {
			"basic":          testAccDevCenterDevBoxDefinition_basic,
			"requiresImport": testAccDevCenterDevBoxDefinition_requiresImport,
			"update":         testAccDevCenterDevBoxDefinition_update,
			"complete":       testAccDevCenterDevBoxDefinition_complete,

			"basicDataSource": testAccDevCenterDevBoxDefinitionDataSource_basic,
		},

		"DevCenterEnvironmentType": {
			"basic":          testAccDevCenterEnvironmentType_basic,
			"requiresImport": testAccDevCenterEnvironmentType_requiresImport,
			"update":         testAccDevCenterEnvironmentType_update,
			"complete":       testAccDevCenterEnvironmentType_complete,

			"basicDataSource": testAccDevCenterEnvironmentTypeDataSource_basic,
		},

		"DevCenterGallery": {
			"basic":          testAccDevCenterGallery_basic,
			"requiresImport": testAccDevCenterGallery_requiresImport,
			"update":         testAccDevCenterGallery_update,
			"complete":       testAccDevCenterGallery_complete,

			"basicDataSource": testAccDevCenterGalleryDataSource_basic,
		},

		"DevCenterNetworkConnection": {
			"basic":          testAccDevCenterNetworkConnection_basic,
			"requiresImport": testAccDevCenterNetworkConnection_requiresImport,
			"update":         testAccDevCenterNetworkConnection_update,
			"complete":       testAccDevCenterNetworkConnection_complete,

			"basicDataSource": testAccDevCenterNetworkConnectionDataSource_basic,
		},

		"DevCenterProject": {
			"basic":          testAccDevCenterProject_basic,
			"requiresImport": testAccDevCenterProject_requiresImport,
			"update":         testAccDevCenterProject_update,
			"complete":       testAccDevCenterProject_complete,

			"basicDataSource": testAccDevCenterProjectDataSource_basic,
		},

		"DevCenterProjectEnvironmentType": {
			"basic":          testAccDevCenterProjectEnvironmentType_basic,
			"requiresImport": testAccDevCenterProjectEnvironmentType_requiresImport,
			"update":         testAccDevCenterProjectEnvironmentType_update,
			"complete":       testAccDevCenterProjectEnvironmentType_complete,

			"basicDataSource": testAccDevCenterProjectEnvironmentTypeDataSource_basic,
		},

		"DevCenterProjectPool": {
			"basic":          testAccDevCenterProjectPool_basic,
			"requiresImport": testAccDevCenterProjectPool_requiresImport,
			"update":         testAccDevCenterProjectPool_update,
			"complete":       testAccDevCenterProjectPool_complete,
			"managedNetwork": testAccDevCenterProjectPool_managedNetwork,

			"basicDataSource": testAccDevCenterProjectPoolDataSource_basic,
		},
	}

	for resource, tests := range testCases {
		t.Run(resource, func(t *testing.T) {
			for name, test := range tests {
				t.Run(name, func(t *testing.T) {
					test(t)
				})
			}
		})
	}
}

func testAccDevCenter_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center", "test")
	r := DevCenterTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccDevCenter_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center", "test")
	r := DevCenterTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func testAccDevCenter_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center", "test")
	r := DevCenterTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccDevCenter_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center", "test")
	r := DevCenterTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
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

func (r DevCenterTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := devcenters.ParseDevCenterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DevCenter.V20250201.DevCenters.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r DevCenterTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_dev_center" "test" {
  location            = azurerm_resource_group.test.location
  name                = "acctestdc-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
}
`, r.template(data))
}

func (r DevCenterTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dev_center" "import" {
  location            = azurerm_dev_center.test.location
  name                = azurerm_dev_center.test.name
  resource_group_name = azurerm_dev_center.test.resource_group_name
}
`, r.basic(data))
}

func (r DevCenterTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_dev_center" "test" {
  location                          = azurerm_resource_group.test.location
  name                              = "acctestdc-${var.random_string}"
  resource_group_name               = azurerm_resource_group.test.name
  project_catalog_item_sync_enabled = true
  tags = {
    environment = "terraform-acctests"
    some_key    = "some-value"
  }
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, r.template(data))
}

func (r DevCenterTestResource) template(data acceptance.TestData) string {
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

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-${var.random_integer}"
  location = var.primary_location
}


resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-${var.random_integer}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
