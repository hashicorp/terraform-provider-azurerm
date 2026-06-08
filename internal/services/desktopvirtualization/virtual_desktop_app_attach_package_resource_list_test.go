// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccVirtualDesktopAppAttachPackage_list_basic(t *testing.T) {
	r := VirtualDesktopAppAttachPackageResource{}
	listResourceAddress := "azurerm_virtual_desktop_app_attach_package.list"

	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_app_attach_package", "test")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 3),
				},
			},
			{
				Query:  true,
				Config: r.basicQueryByResourceGroupName(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
				},
			},
		},
	})
}

func (r VirtualDesktopAppAttachPackageResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_desktop_app_attach_package" "test" {
  count = 3

  name                = "acctest-msix-${count.index}-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  display_name        = "XmlNotepad${count.index}"
  host_pool_ids = [
    azurerm_virtual_desktop_host_pool.test.id
  ]
  msix_package_name     = "43906ChrisLovett.XmlNotepad_2.9.0.21_neutral__hndwmj480pefj"
  storage_share_file_id = azurerm_storage_share_file.test6.id

  depends_on = [
    azurerm_virtual_machine_extension.test0,
    azurerm_virtual_machine_extension.test1,
    azurerm_virtual_machine_extension.test2
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualDesktopAppAttachPackageResource) basicQuery() string {
	return `
list "azurerm_virtual_desktop_app_attach_package" "list" {
  provider = azurerm
  config {}
}
`
}

func (r VirtualDesktopAppAttachPackageResource) basicQueryByResourceGroupName() string {
	return `
list "azurerm_virtual_desktop_app_attach_package" "list" {
  provider = azurerm
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`
}
