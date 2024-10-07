// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleassignmentschedules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PimActiveRoleAssignmentResource struct{}

func TestAccPimActiveRoleAssignment_noExpiration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_pim_active_role_assignment", "test")
	r := PimActiveRoleAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.noExpiration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
			),
		},
		data.ImportStep("schedule.0.start_date_time"),
	})
}

func TestAccPimActiveRoleAssignment_expirationByDurationHours(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_pim_active_role_assignment", "test")
	r := PimActiveRoleAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.expirationByDurationHours(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
			),
		},
		data.ImportStep("schedule.0.start_date_time"),
	})
}

func TestAccPimActiveRoleAssignment_expirationByDurationDays(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_pim_active_role_assignment", "test")
	r := PimActiveRoleAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.expirationByDurationDays(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				check.That(data.ResourceName).Key("schedule.0.expiration.0.duration_days").HasValue("8"),
			),
		},
		data.ImportStep("schedule.0.start_date_time"),
	})
}

func TestAccPimActiveRoleAssignment_pending(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_pim_active_role_assignment", "test")
	r := PimActiveRoleAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.pending(data),
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
			Config: r.importSetup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccPimActiveRoleAssignment_expirationByDate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_pim_active_role_assignment", "test")
	r := PimActiveRoleAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.expirationByDate(data),
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
		return nil, err
	}

	scopeId, err := commonids.ParseScopeID(id.Scope)
	if err != nil {
		return nil, err
	}

	schedulesResult, err := client.Authorization.RoleAssignmentSchedulesClient.ListForScopeComplete(ctx, *scopeId, roleassignmentschedules.ListForScopeOperationOptions{
		Filter: pointer.To(fmt.Sprintf("(principalId eq '%s')", id.PrincipalId)),
	})
	if err != nil {
		return nil, fmt.Errorf("listing role assignment schedules for %s: %+v", scopeId, err)
	}

	for _, schedule := range schedulesResult.Items {
		if props := schedule.Properties; props != nil {
			if props.RoleDefinitionId != nil && strings.EqualFold(*props.RoleDefinitionId, id.RoleDefinitionId) &&
				props.Scope != nil && strings.EqualFold(*props.Scope, scopeId.ID()) &&
				props.PrincipalId != nil && strings.EqualFold(*props.PrincipalId, id.PrincipalId) &&
				props.MemberType != nil && *props.MemberType == roleassignmentschedules.MemberTypeDirect {
				return utils.Bool(true), nil
			}
		}
	}

	return utils.Bool(false), nil
}

func (PimActiveRoleAssignmentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azuread_domains" "test" {
  only_initial = true
}

resource "azuread_user" "test" {
  user_principal_name = "acctestUser-%[1]d1@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestUser-%[1]d1"
  password            = "p@$$Wd%[2]s"
}

resource "azuread_group" "test" {
  display_name     = "acctest-group-%[1]d"
  security_enabled = true
}

`, data.RandomInteger, data.RandomString)
}

func (r PimActiveRoleAssignmentResource) noExpiration(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "Monitoring Data Reader"
}

%[1]s

resource "azurerm_pim_active_role_assignment" "test" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"
  principal_id       = azuread_user.test.object_id

  justification = "No Expiration"

  ticket {
    number = "1"
    system = "example ticket system"
  }
}
`, r.template(data))
}

func (PimActiveRoleAssignmentResource) expirationByDurationHours(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "ContainerApp Reader"
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-%[1]d"
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

func (PimActiveRoleAssignmentResource) expirationByDurationDays(data acceptance.TestData) string {
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

func (PimActiveRoleAssignmentResource) importSetup(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "AcrPull"
}

resource "azuread_application_registration" "test" {
  display_name = "acctest-%[1]d"
}

resource "azuread_service_principal" "test" {
  client_id = azuread_application_registration.test.client_id
}

resource "time_static" "test" {}

resource "azurerm_pim_active_role_assignment" "test" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"
  principal_id       = azuread_service_principal.test.object_id

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
`, data.RandomInteger)
}

func (r PimActiveRoleAssignmentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_pim_active_role_assignment" "import" {
  scope              = azurerm_pim_active_role_assignment.test.scope
  role_definition_id = azurerm_pim_active_role_assignment.test.role_definition_id
  principal_id       = azurerm_pim_active_role_assignment.test.principal_id
}
`, r.importSetup(data))
}

func (PimActiveRoleAssignmentResource) expirationByDate(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "Workbook Contributor"
}

resource "azuread_application_registration" "test" {
  display_name = "acctest-%[1]d"
}

resource "azuread_service_principal" "test" {
  client_id = azuread_application_registration.test.client_id
}

resource "time_static" "test" {}
resource "time_offset" "test" {
  offset_days = 7
}

resource "azurerm_pim_active_role_assignment" "test" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"
  principal_id       = azuread_service_principal.test.object_id

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
`, data.RandomInteger)
}

func (PimActiveRoleAssignmentResource) pending(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "Key Vault Reader"
}

resource "azuread_application_registration" "test" {
  display_name = "acctest-%[1]d"
}

resource "azuread_service_principal" "test" {
  client_id = azuread_application_registration.test.client_id
}

resource "time_offset" "test" {
  offset_days = 1
}

resource "azurerm_pim_active_role_assignment" "test" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"
  principal_id       = azuread_service_principal.test.object_id

  schedule {
    start_date_time = time_offset.test.rfc3339
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
`, data.RandomInteger)
}
