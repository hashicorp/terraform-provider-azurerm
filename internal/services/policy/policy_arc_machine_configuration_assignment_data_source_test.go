// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package policy_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type PolicyArcMachineConfigurationAssignmentDataSource struct{}

func TestAccPolicyArcMachineConfigurationAssignmentDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_arc_machine_configuration_assignment", "test")
	r := PolicyArcMachineConfigurationAssignmentDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("compliance_status").Exists(),
			),
		},
	})
}

func (r PolicyArcMachineConfigurationAssignmentDataSource) template(data acceptance.TestData) string {
	rgid := data.RandomInteger
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-arcgc-%d"
  location = "%s"
}

resource "azurerm_arc_machine" "test" {
  name                = "acctestrg-arcgc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  kind                = "SCVMM"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    foo = "bar"
  }
}
`, rgid, data.Locations.Primary, rgid)
}

func (r PolicyArcMachineConfigurationAssignmentDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_arc_machine_configuration_assignment" "test" {
  name       = "AzureWindowsBaseline"
  location   = azurerm_arc_machine.test.location
  machine_id = azurerm_arc_machine.test.id

  configuration {
    assignment_type = "ApplyAndMonitor"
    version         = "1.*"

    parameter {
      name  = "Minimum Password Length;ExpectedValue"
      value = "16"
    }
    parameter {
      name  = "Minimum Password Age;ExpectedValue"
      value = "0"
    }
    parameter {
      name  = "Maximum Password Age;ExpectedValue"
      value = "30,45"
    }
    parameter {
      name  = "Enforce Password History;ExpectedValue"
      value = "10"
    }
    parameter {
      name  = "Password Must Meet Complexity Requirements;ExpectedValue"
      value = "1"
    }
  }
}

data "azurerm_policy_arc_machine_configuration_assignment" "test" {
  name                = azurerm_policy_arc_machine_configuration_assignment.test.name
  resource_group_name = azurerm_resource_group.test.name
  machine_name        = azurerm_arc_machine.test.name
}
`, r.template(data))
}
