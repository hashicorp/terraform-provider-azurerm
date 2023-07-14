// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SharedImageVersionsDataSource struct{}

func TestAccDataSourceSharedImageVersions_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image_versions", "test")
	r := SharedImageVersionsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config:  SharedImageVersionResource{}.setup(data),
			Destroy: false,
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("images.0.id").Exists(),
				check.That(data.ResourceName).Key("images.#").HasValue("2"),
				check.That(data.ResourceName).Key("images.0.managed_image_id").Exists(),
				check.That(data.ResourceName).Key("images.0.target_region.#").HasValue("1"),
				check.That(data.ResourceName).Key("images.0.target_region.0.storage_account_type").HasValue("Standard_LRS"),
			),
		},
	})
}

func TestAccDataSourceSharedImageVersions_tagsFilterError(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image_versions", "test")
	r := SharedImageVersionsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config:  SharedImageVersionResource{}.setup(data),
			Destroy: false,
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config:      r.tagsFilterError(data),
			ExpectError: regexp.MustCompile("unable to find any images"),
		},
	})
}

func TestAccDataSourceSharedImageVersions_tagsFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image_versions", "test")
	r := SharedImageVersionsDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config:  SharedImageVersionResource{}.setup(data),
			Destroy: false,
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.tagsFilter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("images.#").HasValue("1"),
			),
		},
	})
}

func (r SharedImageVersionsDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image_versions" "test" {
  gallery_name        = azurerm_shared_image_version.test.gallery_name
  image_name          = azurerm_shared_image_version.test.image_name
  resource_group_name = azurerm_shared_image_version.test.resource_group_name
  depends_on          = [azurerm_shared_image_version.test]
}
`, r.template(data))
}

func (r SharedImageVersionsDataSource) tagsFilterError(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image_versions" "test" {
  gallery_name        = azurerm_shared_image_version.test.gallery_name
  image_name          = azurerm_shared_image_version.test.image_name
  resource_group_name = azurerm_shared_image_version.test.resource_group_name

  tags_filter = {
    environment = "error"
    cost-center = "Ops"
  }
}
`, r.template(data))
}

func (r SharedImageVersionsDataSource) tagsFilter(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image_versions" "test" {
  gallery_name        = azurerm_shared_image_version.test.gallery_name
  image_name          = azurerm_shared_image_version.test.image_name
  resource_group_name = azurerm_shared_image_version.test.resource_group_name
  depends_on          = [azurerm_shared_image_version.test]

  tags_filter = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, r.template(data))
}

func (SharedImageVersionsDataSource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image_gallery.test.name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  managed_image_id    = azurerm_image.test.id

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
    foo         = "bar"
  }
}

resource "azurerm_shared_image_version" "test2" {
  name                = "0.0.2"
  gallery_name        = azurerm_shared_image_gallery.test.name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  managed_image_id    = azurerm_image.test.id

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
  }
}
`, SharedImageVersionResource{}.provision(data))
}
