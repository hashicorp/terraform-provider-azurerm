// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/rolemanagementpolicies"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
)

type RoleManagementPolicyResource struct{}

func TestRoleManagementPolicy_resourceGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_management_policy", "test")
	r := RoleManagementPolicyResource{}

	// Ignore the dangling resource post-test as the policy remains while the group is in a pending deletion state
	data.ResourceTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.resourceGroupCreate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("active_assignment_rules.0.expire_after").HasValue("P180D"),
				check.That(data.ResourceName).Key("eligible_assignment_rules.0.expiration_required").HasValue("false"),
				check.That(data.ResourceName).Key("notification_rules.0.eligible_assignments.0.approver_notifications.0.notification_level").HasValue("Critical"),
			),
		},
		data.ImportStep(),
		{
			Config: r.resourceGroupUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("active_assignment_rules.0.expire_after").HasValue("P365D"),
				check.That(data.ResourceName).Key("eligible_assignment_rules.0.expiration_required").HasValue("false"),
				check.That(data.ResourceName).Key("activation_rules.0.approval_stage.0.primary_approver.0.type").HasValue("Group"),
				check.That(data.ResourceName).Key("notification_rules.0.eligible_assignments.0.approver_notifications.0.notification_level").HasValue("Critical"),
			),
		},
	})
}

func TestRoleManagementPolicy_managementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_management_policy", "test")
	r := RoleManagementPolicyResource{}

	// Ignore the dangling resource post-test as the policy remains while the group is in a pending deletion state
	data.ResourceTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.managementGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("active_assignment_rules.0.expire_after").HasValue("P90D"),
				check.That(data.ResourceName).Key("eligible_assignment_rules.0.expiration_required").HasValue("false"),
				check.That(data.ResourceName).Key("notification_rules.0.active_assignments.0.admin_notifications.0.notification_level").HasValue("Critical"),
			),
		},
		data.ImportStep(),
	})
}

func (RoleManagementPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	client := clients.Authorization.RoleManagementPoliciesClient

	id, err := parse.RoleManagementPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, rolemanagementpolicies.NewScopedRoleManagementPolicyID(id.Scope, id.Name))
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("failed to retrieve role management policy with ID %q: %+v", state.ID, err)
	}

	return pointer.To(true), nil
}

func resourceGroupBase(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]s"
  location = "%[2]s"
}

data "azurerm_role_definition" "contributor" {
  name  = "Contributor"
  scope = azurerm_resource_group.test.id
}
`, data.RandomString, data.Locations.Primary)
}

func (RoleManagementPolicyResource) resourceGroupCreate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_role_management_policy" "test" {
  scope              = azurerm_resource_group.test.id
  role_definition_id = data.azurerm_role_definition.contributor.id

  active_assignment_rules {
    expire_after = "P180D"
  }

  eligible_assignment_rules {
    expiration_required = false
  }

  notification_rules {
    eligible_assignments {
      approver_notifications {
        notification_level    = "Critical"
        default_recipients    = false
        additional_recipients = ["someone@example.com"]
      }
    }
  }
}
`, resourceGroupBase(data), data.RandomString)
}

func (RoleManagementPolicyResource) resourceGroupUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azuread" {}

resource "azuread_group" "approver" {
  display_name     = "PIM Approver Test %[2]s"
  mail_enabled     = false
  security_enabled = true
}

resource "azurerm_role_management_policy" "test" {
  scope              = azurerm_resource_group.test.id
  role_definition_id = data.azurerm_role_definition.contributor.id

  active_assignment_rules {
    expire_after = "P365D"
  }

  eligible_assignment_rules {
    expiration_required = false
  }

  activation_rules {
    maximum_duration = "PT1H"
    require_approval = true
    approval_stage {
      primary_approver {
        object_id = azuread_group.approver.object_id
        type      = "Group"
      }
    }
  }

  notification_rules {
    eligible_assignments {
      approver_notifications {
        notification_level    = "Critical"
        default_recipients    = false
        additional_recipients = ["someone@example.com"]
      }
    }
    eligible_activations {
      assignee_notifications {
        notification_level    = "All"
        default_recipients    = true
        additional_recipients = ["someone.else@example.com"]
      }
    }
  }
}
`, resourceGroupBase(data), data.RandomString)
}

func (RoleManagementPolicyResource) managementGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {}

provider "azuread" {}

data "azuread_domains" "test" {
  only_initial = true
}

resource "azuread_user" "approver" {
  user_principal_name = "pam-approver-%[1]s@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "PAM Approver Test %[1]s"
  password            = "p@$$Wd%[1]s"
}

data "azurerm_client_config" "current" {
}

resource "azurerm_management_group" "test" {
  name = "acctest%[1]s"
}

data "azurerm_role_definition" "contributor" {
  name  = "Contributor"
  scope = azurerm_management_group.test.id
}

resource "azurerm_role_management_policy" "test" {
  scope              = azurerm_management_group.test.id
  role_definition_id = data.azurerm_role_definition.contributor.id

  eligible_assignment_rules {
    expiration_required = false
  }

  active_assignment_rules {
    expire_after = "P90D"
  }

  activation_rules {
    maximum_duration = "PT1H"
    require_approval = true
    approval_stage {
      primary_approver {
        object_id = azuread_user.approver.object_id
        type      = "User"
      }
    }
  }

  notification_rules {
    active_assignments {
      admin_notifications {
        notification_level    = "Critical"
        default_recipients    = false
        additional_recipients = ["someone@example.com"]
      }
    }
  }
}
`, data.RandomString)
}
