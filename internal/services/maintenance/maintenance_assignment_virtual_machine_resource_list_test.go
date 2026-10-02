// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package maintenance_test

import (
	"context"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccMaintenanceAssignmentVirtualMachine_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment_virtual_machine", "list")
	r := MaintenanceAssignmentVirtualMachineResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.linkMultipleMaintenanceAssignmentsToOneVM(data),
			},
			{
				Query:  true,
				Config: r.basicListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_maintenance_assignment_virtual_machine.list", 2),
					querycheck.ExpectIdentity(
						"azurerm_maintenance_assignment_virtual_machine.list",
						map[string]knownvalue.Check{
							"scope": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"name":  knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
						},
					),
				},
			},
		},
	})
}

func (MaintenanceAssignmentVirtualMachineResource) basicListQuery() string {
	return `
list "azurerm_maintenance_assignment_virtual_machine" "list" {
  provider = azurerm
  config {
    virtual_machine_id = azurerm_linux_virtual_machine.test.id
  }
}
`
}
