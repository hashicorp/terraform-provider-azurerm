// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package systemcentervirtualmachinemanager_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccSystemCenterVirtualMachineManagerServerListSequential(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests because only one System Center Virtual Machine Manager Server can be onboarded at a time on a given Custom Location

	if os.Getenv("ARM_TEST_CUSTOM_LOCATION_ID") == "" || os.Getenv("ARM_TEST_FQDN") == "" || os.Getenv("ARM_TEST_USERNAME") == "" || os.Getenv("ARM_TEST_PASSWORD") == "" {
		t.Skip("Skipping as one of `ARM_TEST_CUSTOM_LOCATION_ID`, `ARM_TEST_FQDN`, `ARM_TEST_USERNAME`, `ARM_TEST_PASSWORD` was not specified")
	}

	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"scvmmServerList": {
			"basic": testAccSystemCenterVirtualMachineManagerServer_list_basic,
		},
	})
}

func testAccSystemCenterVirtualMachineManagerServer_list_basic(t *testing.T) {

	r := SystemCenterVirtualMachineManagerServerResource{}
	listResourceAddress := "azurerm_system_center_virtual_machine_manager_server.list"

	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_server", "testlist")

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
					querycheck.ExpectLengthAtLeast(listResourceAddress, 2),
				},
			},
			{
				Query:  true,
				Config: r.basicQueryByResourceGroupName(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 2),
				},
			},
		},
	})
}

func (r SystemCenterVirtualMachineManagerServerResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_system_center_virtual_machine_manager_server" "test" {
  count = 2

  name                = "acctest-scvmmms-${count.index}-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = "%s"
  fqdn                = "%s"
  username            = "%s"
  password            = "%s"
}
`, r.template(data), data.RandomInteger, os.Getenv("ARM_TEST_CUSTOM_LOCATION_ID"), os.Getenv("ARM_TEST_FQDN"), os.Getenv("ARM_TEST_USERNAME"), os.Getenv("ARM_TEST_PASSWORD"))
}

func (r SystemCenterVirtualMachineManagerServerResource) basicQuery() string {
	return `
list "azurerm_system_center_virtual_machine_manager_server" "list" {
  provider = azurerm
  config {}
}
`
}

func (r SystemCenterVirtualMachineManagerServerResource) basicQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_system_center_virtual_machine_manager_server" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctestrg-scvmmms-%d"
  }
}
`, data.RandomInteger)
}
