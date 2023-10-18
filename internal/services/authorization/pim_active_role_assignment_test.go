package authorization_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleassignmentscheduleinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PimActiveRoleAssignmentResource struct{}

// TODO: update the management policy configuration so that it can have no expiration
// Depends on new resource - azurerm_role_management_policy - https://github.com/hashicorp/terraform-provider-azurerm/pull/20496
// func TestAccPimActiveRoleAssignment_noExpiration(t *testing.T) {
// 	data := acceptance.BuildTestData(t, "azurerm_pim_active_role_assignment", "test")
// 	r := PimActiveRoleAssignmentResource{}

// 	data.ResourceTest(t, r, []acceptance.TestStep{
// 		{
// 			Config: r.noExpirationConfig(),
// 			Check: acceptance.ComposeTestCheckFunc(
// 				check.That(data.ResourceName).ExistsInAzure(r),
// 				check.That(data.ResourceName).Key("scope").Exists(),
// 			),
// 		},
// 		data.ImportStep("schedule.0.start_date_time"),
// 	})
// }

func TestAccPimActiveRoleAssignment_expirationByDurationHoursConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_pim_active_role_assignment", "test")
	r := PimActiveRoleAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.expirationByDurationHoursConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
			),
		},
		data.ImportStep("schedule.0.start_date_time"),
	})
}

func TestAccPimActiveRoleAssignment_expirationByDurationDaysConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_pim_active_role_assignment", "test")
	r := PimActiveRoleAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.expirationByDurationDaysConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				check.That(data.ResourceName).Key("schedule.0.expiration.0.duration_days").HasValue("8"),
			),
		},
		data.ImportStep("schedule.0.start_date_time"),
	})
}

func TestAccPimActiveRoleAssignment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_pim_active_role_assignment", "test")
	r := PimActiveRoleAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.update1(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
			),
		},
		data.ImportStep("schedule.0.start_date_time"),
		{
			Config: r.update2(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
			),
		},
		data.ImportStep("schedule.0.start_date_time"),
	})
}

func TestAccPimActiveRoleAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_pim_active_role_assignment", "test")

	r := PimActiveRoleAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.importTest(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
			),
		},
		data.RequiresImportErrorStep(func(data acceptance.TestData) string {
			return r.requiresImport()
		}),
	})
}

func TestAccPimActiveRoleAssignment_expirationByDateConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_pim_active_role_assignment", "test")
	r := PimActiveRoleAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.expirationByDateConfig(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
			),
		},
		data.ImportStep("schedule.0.start_date_time", "schedule.0.expiration.0.end_date_time"),
	})
}

func (r PimActiveRoleAssignmentResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.PimRoleAssignmentID(state.ID)
	if err != nil {
		return utils.Bool(false), err
	}

	filter := &roleassignmentscheduleinstances.ListForScopeOperationOptions{
		Filter: pointer.To(fmt.Sprintf("(principalId eq '%s' and roleDefinitionId eq '%s')", id.PrincipalId, id.RoleDefinitionId)),
	}

	items, err := client.Authorization.RoleAssignmentScheduleInstancesClient.ListForScopeComplete(ctx, id.ScopeID(), *filter)
	if err != nil {
		return nil, fmt.Errorf("listing role assignments on scope %s: %+v", id, err)
	}

	foundDirectAssignment := false

	for _, i := range items.Items {
		if *i.Properties.MemberType == roleassignmentscheduleinstances.MemberTypeDirect {
			foundDirectAssignment = true
			break
		}
	}

	return utils.Bool(foundDirectAssignment), nil
}

// func (PimActiveRoleAssignmentResource) noExpirationConfig() string {
// 	return `
// data "azurerm_subscription" "primary" {}

// data "azurerm_client_config" "test" {}

// data "azurerm_role_definition" "test" {
//   name = "Monitoring Data Reader"
// }

// resource "azurerm_pim_active_role_assignment" "test" {
//   scope              = data.azurerm_subscription.primary.id
//   role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"
//   principal_id       = data.azurerm_client_config.test.object_id

//   justification = "No Expiration"

//   ticket {
//     number = "1"
//     system = "example ticket system"
//   }
// }
// `
// }

func (PimActiveRoleAssignmentResource) expirationByDurationHoursConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "ContainerApp Reader"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "time_static" "test" {}

resource "azurerm_pim_active_role_assignment" "test" {
  scope              = azurerm_resource_group.test.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"
  principal_id       = data.azurerm_client_config.test.object_id

  schedule {
    start_date_time = time_static.test.rfc3339
    expiration {
      duration_hours = 8
    }
  }

  justification = "Expiration Duration Set"

  ticket {
    number = "1"
    system = "example ticket system"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (PimActiveRoleAssignmentResource) expirationByDurationDaysConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "Disk Backup Reader"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
resource "azurerm_virtual_network" "test" {
  name                = "acctestVNET1-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "time_static" "test" {}

resource "azurerm_pim_active_role_assignment" "test" {
  scope              = azurerm_virtual_network.test.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"
  principal_id       = data.azurerm_client_config.test.object_id

  schedule {
    start_date_time = time_static.test.rfc3339
    expiration {
      duration_days = 8
    }
  }

  justification = "Expiration Duration Set"

  ticket {
    number = "1"
    system = "example ticket system"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (PimActiveRoleAssignmentResource) importTest() string {
	return `
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "AcrPull"
}

resource "time_static" "test" {}

resource "azurerm_pim_active_role_assignment" "test" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"
  principal_id       = data.azurerm_client_config.test.object_id

  schedule {
    start_date_time = time_static.test.rfc3339
    expiration {
      duration_hours = 8
    }
  }

  justification = "Expiration Duration Set"

  ticket {
    number = "1"
    system = "example ticket system"
  }
}
`
}

func (r PimActiveRoleAssignmentResource) requiresImport() string {
	return fmt.Sprintf(`
%s

resource "azurerm_pim_active_role_assignment" "import" {
  scope              = azurerm_pim_active_role_assignment.test.scope
  role_definition_id = azurerm_pim_active_role_assignment.test.role_definition_id
  principal_id       = azurerm_pim_active_role_assignment.test.principal_id
}
`, r.importTest())
}

func (PimActiveRoleAssignmentResource) expirationByDateConfig() string {
	return `
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "Workbook Reader"
}

resource "time_static" "test" {}
resource "time_offset" "test" {
  offset_days = 7
}

resource "azurerm_pim_active_role_assignment" "test" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"
  principal_id       = data.azurerm_client_config.test.object_id

  schedule {
    start_date_time = time_static.test.rfc3339
    expiration {
      end_date_time = time_offset.test.rfc3339
    }
  }

  justification = "Expiration End Date Set"

  ticket {
    number = "1"
    system = "example ticket system"
  }
}
`
}

func (PimActiveRoleAssignmentResource) update1() string {
	return `
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "Billing Reader"
}

resource "time_static" "test" {}

resource "azurerm_pim_active_role_assignment" "test" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"
  principal_id       = data.azurerm_client_config.test.object_id

  schedule {
    start_date_time = time_static.test.rfc3339
    expiration {
      duration_hours = 8
    }
  }

  justification = "Expiration Duration Set"

  ticket {
    number = "1"
    system = "example ticket system"
  }
}
`
}

func (PimActiveRoleAssignmentResource) update2() string {
	return `
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "Billing Reader"
}

resource "time_static" "test" {}

resource "azurerm_pim_active_role_assignment" "test" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"
  principal_id       = data.azurerm_client_config.test.object_id

  schedule {
    start_date_time = time_static.test.rfc3339
    expiration {
      duration_hours = 8
    }
  }

  justification = "Expiration Duration Set"

  ticket {
    number = "1"
    system = "example ticket system"
  }
}
`
}
