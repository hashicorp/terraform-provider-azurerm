// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package maintenance_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccMaintenanceAssignmentArcMachine_listByArcMachineID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment_arc_machine", "test")
	r := MaintenanceAssignmentArcMachineResource{}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acceptance.PreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
			},
			{
				Query: true,
				Config: `
list "azurerm_maintenance_assignment_arc_machine" "test" {
  provider         = azurerm
  include_resource = true
  config {
    arc_machine_id = azurerm_arc_machine.test.id
  }
}
`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_maintenance_assignment_arc_machine.test", 1),
					querycheck.ExpectResourceKnownValues("azurerm_maintenance_assignment_arc_machine.test", nil, []querycheck.KnownValueCheck{
						{
							Path:       tfjsonpath.New("maintenance_configuration_id"),
							KnownValue: knownvalue.StringRegexp(regexp.MustCompile(fmt.Sprintf("(?i)^/subscriptions/%s/resourceGroups/acctestRG-maint-%d/providers/Microsoft.Maintenance/maintenanceConfigurations/acctest-mc%d$", data.Subscriptions.Primary, data.RandomInteger, data.RandomInteger))),
						},
					}),
					querycheck.ExpectIdentity("azurerm_maintenance_assignment_arc_machine.test", map[string]knownvalue.Check{
						"scope": knownvalue.StringRegexp(regexp.MustCompile(fmt.Sprintf("(?i)^/subscriptions/%s/resourceGroups/acctestRG-maint-%d/providers/Microsoft.HybridCompute/machines/acctest-arc-%d$", data.Subscriptions.Primary, data.RandomInteger, data.RandomInteger))),
						"name":  knownvalue.StringRegexp(regexp.MustCompile(fmt.Sprintf("(?i)^acctest-mc%d$", data.RandomInteger))),
					}),
				},
			},
		},
	})
}
