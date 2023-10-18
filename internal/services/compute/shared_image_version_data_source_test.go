// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

type SharedImageVersionDataSource struct{}

func TestAccDataSourceSharedImageVersion_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image_version", "test")
	r := SharedImageVersionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: SharedImageVersionResource{}.setup(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("managed_image_id").Exists(),
				check.That(data.ResourceName).Key("target_region.#").HasValue("1"),
				check.That(data.ResourceName).Key("target_region.0.storage_account_type").HasValue("Standard_LRS"),
			),
		},
	})
}

func TestAccDataSourceSharedImageVersion_latest(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image_version", "test")
	r := SharedImageVersionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: SharedImageVersionResource{}.setup(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.latest(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("0.0.2"),
				check.That(data.ResourceName).Key("managed_image_id").Exists(),
				check.That(data.ResourceName).Key("target_region.#").HasValue("1"),
				check.That(data.ResourceName).Key("target_region.0.storage_account_type").HasValue("Standard_LRS"),
			),
		},
	})
}

func TestAccDataSourceSharedImageVersion_excludeFromLatest(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image_version", "test")
	r := SharedImageVersionDataSource{}

	expectedVersion := "0.0.1"
	if !features.FourPointOhBeta() {
		// `ExcludeFromLatest` is not considered in 3.0 so `0.0.2` will still be the latest image
		expectedVersion = "0.0.2"
	}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: SharedImageVersionResource{}.setup(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.excludeFromLatest(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(expectedVersion),
				check.That(data.ResourceName).Key("managed_image_id").Exists(),
				check.That(data.ResourceName).Key("target_region.#").HasValue("1"),
				check.That(data.ResourceName).Key("target_region.0.storage_account_type").HasValue("Standard_LRS"),
			),
		},
	})
}

func TestAccDataSourceSharedImageVersion_recent(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image_version", "test")
	r := SharedImageVersionDataSource{}

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
			Config: r.recent(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("0.0.2"),
				check.That(data.ResourceName).Key("managed_image_id").Exists(),
				check.That(data.ResourceName).Key("target_region.#").HasValue("1"),
				check.That(data.ResourceName).Key("target_region.0.storage_account_type").HasValue("Standard_LRS"),
			),
		},
	})
}

func TestAccDataSourceSharedImageVersion_sortVersionsBySemver(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image_version", "test")
	r := SharedImageVersionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: SharedImageVersionResource{}.setup(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.sortVersionsBySemver(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("0.0.10"),
				check.That(data.ResourceName).Key("managed_image_id").Exists(),
				check.That(data.ResourceName).Key("target_region.#").HasValue("1"),
				check.That(data.ResourceName).Key("target_region.0.storage_account_type").HasValue("Standard_LRS"),
			),
		},
	})
}

func (SharedImageVersionDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image_version" "test" {
  name                = azurerm_shared_image_version.test.name
  gallery_name        = azurerm_shared_image_version.test.gallery_name
  image_name          = azurerm_shared_image_version.test.image_name
  resource_group_name = azurerm_shared_image_version.test.resource_group_name
}
`, SharedImageVersionResource{}.imageVersion(data))
}

func (SharedImageVersionDataSource) latest(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

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

data "azurerm_shared_image_version" "test" {
  name                = "latest"
  gallery_name        = azurerm_shared_image_version.test2.gallery_name
  image_name          = azurerm_shared_image_version.test2.image_name
  resource_group_name = azurerm_shared_image_version.test2.resource_group_name
}
`, SharedImageVersionResource{}.imageVersion(data))
}

func (SharedImageVersionDataSource) excludeFromLatest(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "test2" {
  name                = "0.0.2"
  gallery_name        = azurerm_shared_image_gallery.test.name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  managed_image_id    = azurerm_image.test.id
  exclude_from_latest = true

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
  }
}

data "azurerm_shared_image_version" "test" {
  name                = "latest"
  gallery_name        = azurerm_shared_image_version.test2.gallery_name
  image_name          = azurerm_shared_image_version.test2.image_name
  resource_group_name = azurerm_shared_image_version.test2.resource_group_name
}
`, SharedImageVersionResource{}.imageVersion(data))
}

func (SharedImageVersionDataSource) recent(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "test2" {
  name                = "0.0.3"
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

resource "azurerm_shared_image_version" "test3" {
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

  depends_on = [
    azurerm_shared_image_version.test2
  ]
}

data "azurerm_shared_image_version" "test" {
  name                = "recent"
  gallery_name        = azurerm_shared_image_version.test3.gallery_name
  image_name          = azurerm_shared_image_version.test3.image_name
  resource_group_name = azurerm_shared_image_version.test3.resource_group_name
}
`, SharedImageVersionResource{}.imageVersion(data))
}

func (SharedImageVersionDataSource) sortVersionsBySemver(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "test2" {
  name                = "0.0.9"
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

resource "azurerm_shared_image_version" "test3" {
  name                = "0.0.10"
  gallery_name        = azurerm_shared_image_gallery.test.name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  managed_image_id    = azurerm_image.test.id

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = 1
  }

  depends_on = [
    azurerm_shared_image_version.test2
  ]
}

data "azurerm_shared_image_version" "test" {
  name                = "latest"
  gallery_name        = azurerm_shared_image_version.test3.gallery_name
  image_name          = azurerm_shared_image_version.test3.image_name
  resource_group_name = azurerm_shared_image_version.test3.resource_group_name

  sort_versions_by_semver = true
}
`, SharedImageVersionResource{}.imageVersion(data))
}
