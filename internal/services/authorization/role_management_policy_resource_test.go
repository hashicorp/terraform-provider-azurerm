// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
)

type RoleManagementPolicyResource struct{}

func TestAccRoleManagementPolicy_managementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_management_policy", "test")
	r := RoleManagementPolicyResource{}

	// Ignore the dangling resource post-test as the policy remains while the management group exists, or is in a pending deletion state
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

func TestAccRoleManagementPolicy_resourceGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_management_policy", "test")
	r := RoleManagementPolicyResource{}

	// Ignore the dangling resource post-test as the policy remains while the resource group exists, or is in a pending deletion state
	data.ResourceTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.resourceGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("active_assignment_rules.0.expire_after").HasValue("P30D"),
				check.That(data.ResourceName).Key("eligible_assignment_rules.0.expiration_required").HasValue("false"),
				check.That(data.ResourceName).Key("notification_rules.0.eligible_assignments.0.approver_notifications.0.notification_level").HasValue("All"),
			),
		},
		data.ImportStep(),
		{
			Config: r.resourceGroupUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("active_assignment_rules.0.expire_after").HasValue("P15D"),
				check.That(data.ResourceName).Key("eligible_assignment_rules.0.expiration_required").HasValue("true"),
				check.That(data.ResourceName).Key("activation_rules.0.approval_stage.0.primary_approver.0.type").HasValue("Group"),
				check.That(data.ResourceName).Key("notification_rules.0.eligible_assignments.0.approver_notifications.0.notification_level").HasValue("Critical"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRoleManagementPolicy_resourceGroup_activationRulesUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_management_policy", "test")
	r := RoleManagementPolicyResource{}

	data.ResourceTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.resourceGroupActivationRules(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("activation_rules.0.require_approval").HasValue("true"),
				check.That(data.ResourceName).Key("notification_rules.0.eligible_assignments.0.approver_notifications.0.notification_level").HasValue("All"),
			),
		},
		data.ImportStep(),
		{
			Config: r.resourceGroupActivationRules(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("activation_rules.0.require_approval").HasValue("false"),
				check.That(data.ResourceName).Key("notification_rules.0.eligible_assignments.0.approver_notifications.0.notification_level").HasValue("All"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRoleManagementPolicy_subscription(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_management_policy", "test")
	r := RoleManagementPolicyResource{}

	// Ignore the dangling resource post-test as the policy remains while the subscription exists, or is in a pending deletion state
	data.ResourceTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.subscription(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("active_assignment_rules.0.expire_after").HasValue("P180D"),
				check.That(data.ResourceName).Key("eligible_assignment_rules.0.expiration_required").HasValue("false"),
				check.That(data.ResourceName).Key("notification_rules.0.eligible_assignments.0.approver_notifications.0.notification_level").HasValue("Critical"),
			),
		},
		data.ImportStep(),
		{
			Config: r.subscriptionUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("active_assignment_rules.0.expire_after").HasValue("P365D"),
				check.That(data.ResourceName).Key("eligible_assignment_rules.0.expiration_required").HasValue("false"),
				check.That(data.ResourceName).Key("activation_rules.0.approval_stage.0.primary_approver.0.type").HasValue("Group"),
				check.That(data.ResourceName).Key("notification_rules.0.eligible_assignments.0.approver_notifications.0.notification_level").HasValue("Critical"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRoleManagementPolicy_resource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_management_policy", "test")
	r := RoleManagementPolicyResource{}

	// Ignore the dangling resource post-test as the policy remains while the resource exists, or is in a pending deletion state
	data.ResourceTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.resource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("active_assignment_rules.0.expire_after").HasValue("P30D"),
				check.That(data.ResourceName).Key("eligible_assignment_rules.0.expiration_required").HasValue("false"),
				check.That(data.ResourceName).Key("notification_rules.0.eligible_assignments.0.approver_notifications.0.notification_level").HasValue("All"),
			),
		},
		data.ImportStep(),
		{
			Config: r.resourceUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("active_assignment_rules.0.expire_after").HasValue("P15D"),
				check.That(data.ResourceName).Key("eligible_assignment_rules.0.expiration_required").HasValue("true"),
				check.That(data.ResourceName).Key("activation_rules.0.approval_stage.0.primary_approver.0.type").HasValue("Group"),
				check.That(data.ResourceName).Key("notification_rules.0.eligible_assignments.0.approver_notifications.0.notification_level").HasValue("Critical"),
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

	policyId, err := authorization.FindRoleManagementPolicyId(ctx, clients.Authorization.RoleManagementPoliciesClient, id.Scope, id.RoleDefinitionId)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, *policyId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", policyId, err)
	}

	return pointer.To(true), nil
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

func (RoleManagementPolicyResource) resourceGroupTemplate(data acceptance.TestData) string {
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

func (r RoleManagementPolicyResource) resourceGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_role_management_policy" "test" {
  scope              = azurerm_resource_group.test.id
  role_definition_id = data.azurerm_role_definition.contributor.id

  active_assignment_rules {
    expire_after = "P30D"
  }

  eligible_assignment_rules {
    expiration_required = false
  }

  notification_rules {
    eligible_assignments {
      approver_notifications {
        notification_level    = "All"
        default_recipients    = false
        additional_recipients = ["someone@example.com"]
      }
    }
  }
}
`, r.resourceGroupTemplate(data), data.RandomString)
}

func (r RoleManagementPolicyResource) resourceGroupUpdate(data acceptance.TestData) string {
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
    expire_after = "P15D"
  }

  eligible_assignment_rules {
    expiration_required = true
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
`, r.resourceGroupTemplate(data), data.RandomString)
}

func (RoleManagementPolicyResource) subscriptionTemplate(data acceptance.TestData) string {
	return `
provider "azurerm" {}

data "azurerm_subscription" "test" {}

data "azurerm_role_definition" "contributor" {
  name  = "Contributor"
  scope = data.azurerm_subscription.test.id
}
`
}

func (r RoleManagementPolicyResource) subscription(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_role_management_policy" "test" {
  scope              = data.azurerm_subscription.test.id
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
`, r.subscriptionTemplate(data), data.RandomString)
}

func (r RoleManagementPolicyResource) subscriptionUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azuread" {}

resource "azuread_group" "approver" {
  display_name     = "PIM Approver Test %[2]s"
  mail_enabled     = false
  security_enabled = true
}

resource "azurerm_role_management_policy" "test" {
  scope              = data.azurerm_subscription.test.id
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
`, r.subscriptionTemplate(data), data.RandomString)
}

func (r RoleManagementPolicyResource) resourceGroupActivationRules(data acceptance.TestData, requireApproval bool) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_group" "approver" {
  display_name     = "PIM Approver Test %[2]s"
  mail_enabled     = false
  security_enabled = true
}

resource "azurerm_role_management_policy" "test" {
  scope              = azurerm_resource_group.test.id
  role_definition_id = data.azurerm_role_definition.contributor.id

  activation_rules {
    maximum_duration = "PT1H"
    require_approval = %t
    approval_stage {
      primary_approver {
        object_id = azuread_group.approver.object_id
        type      = "Group"
      }
    }
  }

  eligible_assignment_rules {
    expiration_required = false
  }

  notification_rules {
    eligible_assignments {
      approver_notifications {
        notification_level    = "All"
        default_recipients    = false
        additional_recipients = ["someone@example.com"]
      }
    }
  }
}
`, r.resourceGroupTemplate(data), data.RandomString, requireApproval)
}

func (RoleManagementPolicyResource) resourceTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]s"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "accteststg%[1]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

data "azurerm_role_definition" "contributor" {
  name  = "Contributor"
  scope = azurerm_storage_account.test.id
}
`, data.RandomString, data.Locations.Primary)
}

func (r RoleManagementPolicyResource) resource(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_role_management_policy" "test" {
  scope              = azurerm_storage_account.test.id
  role_definition_id = data.azurerm_role_definition.contributor.id

  active_assignment_rules {
    expire_after = "P30D"
  }

  eligible_assignment_rules {
    expiration_required = false
  }

  notification_rules {
    eligible_assignments {
      approver_notifications {
        notification_level    = "All"
        default_recipients    = false
        additional_recipients = ["someone@example.com"]
      }
    }
  }
}
`, r.resourceTemplate(data), data.RandomString)
}

func (r RoleManagementPolicyResource) resourceUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azuread" {}

resource "azuread_group" "approver" {
  display_name     = "PIM Approver Test %[2]s"
  mail_enabled     = false
  security_enabled = true
}

resource "azurerm_role_management_policy" "test" {
  scope              = azurerm_storage_account.test.id
  role_definition_id = data.azurerm_role_definition.contributor.id

  active_assignment_rules {
    expire_after = "P15D"
  }

  eligible_assignment_rules {
    expiration_required = true
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
`, r.resourceTemplate(data), data.RandomString)
}
