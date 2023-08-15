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

type ImagesDataSource struct{}

func TestAccDataSourceAzureRMImages_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_images", "test")
	r := ImagesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: ImageResource{}.setupUnmanagedDisks(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("images.#").HasValue("2"),
				check.That(data.ResourceName).Key("images.0.os_disk.0.os_type").HasValue("Linux"),
			),
		},
	})
}

func TestAccDataSourceAzureRMImages_tagsFilterError(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_images", "test")
	r := ImagesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: ImageResource{}.setupUnmanagedDisks(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(ImageResource{}.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(ImageResource{}.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config:      r.tagsFilterError(data),
			ExpectError: regexp.MustCompile("no images were found that match the specified tags"),
		},
	})
}

func TestAccDataSourceAzureRMImages_tagsFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_images", "test")
	r := ImagesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: ImageResource{}.setupUnmanagedDisks(data),
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

func (r ImagesDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_images" "test" {
  resource_group_name = azurerm_image.test.resource_group_name
}
`, r.template(data))
}

func (r ImagesDataSource) tagsFilterError(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_images" "test" {
  resource_group_name = azurerm_image.test.resource_group_name
  tags_filter = {
    environment = "error"
    cost-center = "Ops"
  }
}
`, r.template(data))
}

func (r ImagesDataSource) tagsFilter(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_images" "test" {
  resource_group_name = azurerm_image.test.resource_group_name
  tags_filter = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, r.template(data))
}

func (ImagesDataSource) template(data acceptance.TestData) string {
	template := ImageResource{}.setupUnmanagedDisks(data)
	return fmt.Sprintf(`
%s

resource "azurerm_image" "test" {
  name                = "accteste"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  os_disk {
    os_type  = "Linux"
    os_state = "Generalized"
    blob_uri = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    size_gb  = 30
    caching  = "None"
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
    foo         = "bar"
  }
}

resource "azurerm_image" "test2" {
  name                = "accteste2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  os_disk {
    os_type  = "Linux"
    os_state = "Generalized"
    blob_uri = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    size_gb  = 30
    caching  = "None"
  }
}
`, template)
}
