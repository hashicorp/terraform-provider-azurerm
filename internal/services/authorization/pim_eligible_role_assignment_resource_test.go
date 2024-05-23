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
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleeligibilityschedules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PimEligibleRoleAssignmentResource struct{}

func TestAccPimEligibleRoleAssignment_noExpiration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_pim_eligible_role_assignment", "test")
	r := PimEligibleRoleAssignmentResource{}

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

func TestAccPimEligibleRoleAssignment_expirationByDurationHours(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_pim_eligible_role_assignment", "test")
	r := PimEligibleRoleAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.expirationByDurationHours(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				check.That(data.ResourceName).Key("schedule.0.expiration.0.duration_hours").HasValue("8"),
			),
		},
		data.ImportStep("schedule.0.start_date_time"),
	})
}

func TestAccPimEligibleRoleAssignment_expirationByDurationDays(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_pim_eligible_role_assignment", "test")
	r := PimEligibleRoleAssignmentResource{}

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

func TestAccPimEligibleRoleAssignment_pending(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_pim_eligible_role_assignment", "test")
	r := PimEligibleRoleAssignmentResource{}

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

func TestAccPimEligibleRoleAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_pim_eligible_role_assignment", "test")

	r := PimEligibleRoleAssignmentResource{}

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

func TestAccPimEligibleRoleAssignment_expirationByDate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_pim_eligible_role_assignment", "test")
	r := PimEligibleRoleAssignmentResource{}

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

func (r PimEligibleRoleAssignmentResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.PimRoleAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}

	scopeId, err := commonids.ParseScopeID(id.Scope)
	if err != nil {
		return nil, err
	}

	schedulesResult, err := client.Authorization.RoleEligibilitySchedulesClient.ListForScopeComplete(ctx, *scopeId, roleeligibilityschedules.ListForScopeOperationOptions{
		Filter: pointer.To(fmt.Sprintf("(principalId eq '%s')", id.PrincipalId)),
	})
	if err != nil {
		return nil, fmt.Errorf("listing role eligiblity schedules for %s: %+v", scopeId, err)
	}

	for _, schedule := range schedulesResult.Items {
		if props := schedule.Properties; props != nil {
			if props.RoleDefinitionId != nil && strings.EqualFold(*props.RoleDefinitionId, id.RoleDefinitionId) &&
				props.Scope != nil && strings.EqualFold(*props.Scope, scopeId.ID()) &&
				props.PrincipalId != nil && strings.EqualFold(*props.PrincipalId, id.PrincipalId) &&
				props.MemberType != nil && *props.MemberType == roleeligibilityschedules.MemberTypeDirect {
				return utils.Bool(true), nil
			}
		}
	}

	return utils.Bool(false), nil
}

func (PimEligibleRoleAssignmentResource) template(data acceptance.TestData) string {
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

func (r PimEligibleRoleAssignmentResource) noExpiration(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "Disk Backup Reader"
}

%[1]s

resource "azurerm_pim_eligible_role_assignment" "test" {
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

func (r PimEligibleRoleAssignmentResource) expirationByDurationHours(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "ContainerApp Reader"
}

%[1]s

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[3]s"
}

resource "time_static" "test" {}

resource "azurerm_pim_eligible_role_assignment" "test" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"
  principal_id       = azuread_user.test.object_id

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
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r PimEligibleRoleAssignmentResource) expirationByDurationDays(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "Disk Backup Reader"
}

%[1]s

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[3]s"
}
resource "azurerm_virtual_network" "test" {
  name                = "acctestVNET1-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "time_static" "test" {}

resource "azurerm_pim_eligible_role_assignment" "test" {
  scope              = azurerm_virtual_network.test.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"
  principal_id       = azuread_user.test.object_id

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
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r PimEligibleRoleAssignmentResource) expirationByDate(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "Workbook Reader"
}

%[1]s

resource "time_static" "test" {}
resource "time_offset" "test" {
  offset_days = 7
}

resource "azurerm_pim_eligible_role_assignment" "test" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"
  principal_id       = azuread_group.test.object_id

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
`, r.template(data))
}

func (r PimEligibleRoleAssignmentResource) pending(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "Key Vault Contributor"
}

%[1]s

resource "time_offset" "test" {
  offset_days = 1
}

resource "azurerm_pim_eligible_role_assignment" "test" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"
  principal_id       = azuread_user.test.object_id

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
`, r.template(data))
}

func (r PimEligibleRoleAssignmentResource) importSetup(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "AcrPull"
}

%[1]s

resource "time_static" "test" {}

resource "azurerm_pim_eligible_role_assignment" "test" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"
  principal_id       = azuread_user.test.object_id

  schedule {
    start_date_time = time_static.test.rfc3339
    expiration {
      duration_hours = 3
    }
  }

  justification = "Expiration Duration Set"

  ticket {
    number = "1"
    system = "example ticket system"
  }
}
`, r.template(data))
}

func (r PimEligibleRoleAssignmentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_pim_eligible_role_assignment" "import" {
  scope              = azurerm_pim_eligible_role_assignment.test.scope
  role_definition_id = azurerm_pim_eligible_role_assignment.test.role_definition_id
  principal_id       = azurerm_pim_eligible_role_assignment.test.principal_id
}
`, r.importSetup(data))
}
