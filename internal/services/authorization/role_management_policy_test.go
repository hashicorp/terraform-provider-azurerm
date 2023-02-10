package authorization_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type RoleManagementPolicyResource struct{}

func TestAccRoleManagementPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_management_policy", "test")
	r := RoleManagementPolicyResource{}

	data.ResourceTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				check.That(data.ResourceName).Key("role_definition_id").Exists(),
				check.That(data.ResourceName).Key("activation.0.maximum_duration_hours").HasValue("12"),
			),
		},
	})
}

func TestAccRoleManagementPolicy_partial(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_management_policy", "test")
	r := RoleManagementPolicyResource{}
	ri := acceptance.RandTimeInt()
	rPassword := fmt.Sprintf("%s%s", "p@$$Wd", acceptance.RandString(6))

	data.ResourceTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.partial1(ri),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				check.That(data.ResourceName).Key("role_definition_id").Exists(),

				check.That(data.ResourceName).Key("activation.0.require_multi_factor_authentication").HasValue("true"),

				check.That(data.ResourceName).Key("assignment.0.eligible.0.expire_after_hours").HasValue("6"),
				check.That(data.ResourceName).Key("assignment.0.active.0.expire_after_hours").HasValue("5"),

				check.That(data.ResourceName).Key("notifications.0.member_assigned_eligible.0.role_assignment_alert.0.default_recipients").HasValue("false"),

				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.role_assignment_alert.0.default_recipients").HasValue("false"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.role_assignment_alert.0.additional_recipients.#").HasValue("1"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.role_assignment_alert.0.critical_emails_only").HasValue("false"),

				check.That(data.ResourceName).Key("activation.0.approvers.0.group.#").HasValue("1"),
				check.That(data.ResourceName).Key("activation.0.approvers.0.user.#").HasValue("0"),
			),
		},
		{
			Config: r.partial2(ri, rPassword),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				check.That(data.ResourceName).Key("role_definition_id").Exists(),

				check.That(data.ResourceName).Key("activation.0.require_justification").HasValue("true"),
				check.That(data.ResourceName).Key("activation.0.require_multi_factor_authentication").HasValue("true"),

				check.That(data.ResourceName).Key("assignment.0.eligible.0.expire_after_days").HasValue("7"),
				check.That(data.ResourceName).Key("assignment.0.active.0.expire_after_days").HasValue("4"),
				check.That(data.ResourceName).Key("assignment.0.active.0.require_justification").HasValue("true"),
				check.That(data.ResourceName).Key("assignment.0.active.0.require_multi_factor_authentication").HasValue("true"),

				check.That(data.ResourceName).Key("notifications.0.member_assigned_eligible.0.role_assignment_alert.0.default_recipients").HasValue("false"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_eligible.0.role_assignment_alert.0.critical_emails_only").HasValue("false"),

				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.role_assignment_alert.0.default_recipients").HasValue("false"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.role_assignment_alert.0.additional_recipients.#").HasValue("1"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.role_assignment_alert.0.critical_emails_only").HasValue("false"),

				check.That(data.ResourceName).Key("notifications.0.eligible_member_activate.0.role_assignment_alert.0.default_recipients").HasValue("false"),
				check.That(data.ResourceName).Key("notifications.0.eligible_member_activate.0.role_assignment_alert.0.additional_recipients.#").HasValue("1"),
				check.That(data.ResourceName).Key("notifications.0.eligible_member_activate.0.role_assignment_alert.0.critical_emails_only").HasValue("false"),

				check.That(data.ResourceName).Key("activation.0.approvers.0.group.#").HasValue("0"),
				check.That(data.ResourceName).Key("activation.0.approvers.0.user.#").HasValue("1"),
			),
		},
	})
}

func TestAccRoleManagementPolicy_managementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_management_policy", "test")
	r := RoleManagementPolicyResource{}

	data.ResourceTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.managementGroup(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				check.That(data.ResourceName).Key("role_definition_id").Exists(),
				check.That(data.ResourceName).Key("activation.0.maximum_duration_hours").HasValue("12"),
			),
		},
	})
}

func TestAccRoleManagementPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_management_policy", "test")
	r := RoleManagementPolicyResource{}
	ri := acceptance.RandTimeInt()
	rPassword := fmt.Sprintf("%s%s", "p@$$Wd", acceptance.RandString(6))

	data.ResourceTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				check.That(data.ResourceName).Key("role_definition_id").Exists(),
				check.That(data.ResourceName).Key("activation.0.maximum_duration_hours").HasValue("12"),
			),
		},
		{
			Config: r.complete(ri, rPassword),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				check.That(data.ResourceName).Key("role_definition_id").Exists(),

				check.That(data.ResourceName).Key("activation.0.maximum_duration_hours").HasValue("4"),
				check.That(data.ResourceName).Key("activation.0.require_multi_factor_authentication").HasValue("true"),
				check.That(data.ResourceName).Key("activation.0.require_justification").HasValue("true"),
				check.That(data.ResourceName).Key("activation.0.require_ticket_information").HasValue("true"),
				check.That(data.ResourceName).Key("activation.0.require_ticket_information").HasValue("true"),
				check.That(data.ResourceName).Key("activation.0.approvers.0.group.#").HasValue("2"),
				check.That(data.ResourceName).Key("activation.0.approvers.0.user.#").HasValue("2"),

				check.That(data.ResourceName).Key("assignment.0.eligible.0.allow_permanent").HasValue("true"),
				check.That(data.ResourceName).Key("assignment.0.active.0.allow_permanent").HasValue("true"),
				check.That(data.ResourceName).Key("assignment.0.active.0.require_multi_factor_authentication").HasValue("true"),
				check.That(data.ResourceName).Key("assignment.0.active.0.require_justification").HasValue("true"),

				check.That(data.ResourceName).Key("notifications.0.member_assigned_eligible.0.role_assignment_alert.0.default_recipients").HasValue("true"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_eligible.0.role_assignment_alert.0.additional_recipients.#").HasValue("1"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_eligible.0.role_assignment_alert.0.critical_emails_only").HasValue("true"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_eligible.0.assigned_user.0.default_recipients").HasValue("true"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_eligible.0.assigned_user.0.additional_recipients.#").HasValue("1"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_eligible.0.assigned_user.0.critical_emails_only").HasValue("true"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_eligible.0.request_for_extension_or_approval.0.default_recipients").HasValue("true"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_eligible.0.request_for_extension_or_approval.0.additional_recipients.#").HasValue("1"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_eligible.0.request_for_extension_or_approval.0.critical_emails_only").HasValue("true"),

				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.role_assignment_alert.0.default_recipients").HasValue("true"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.role_assignment_alert.0.additional_recipients.#").HasValue("1"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.role_assignment_alert.0.critical_emails_only").HasValue("true"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.assigned_user.0.default_recipients").HasValue("true"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.assigned_user.0.additional_recipients.#").HasValue("1"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.assigned_user.0.critical_emails_only").HasValue("true"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.request_for_extension_or_approval.0.default_recipients").HasValue("true"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.request_for_extension_or_approval.0.additional_recipients.#").HasValue("1"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.request_for_extension_or_approval.0.critical_emails_only").HasValue("true"),

				check.That(data.ResourceName).Key("notifications.0.eligible_member_activate.0.role_assignment_alert.0.default_recipients").HasValue("true"),
				check.That(data.ResourceName).Key("notifications.0.eligible_member_activate.0.role_assignment_alert.0.additional_recipients.#").HasValue("1"),
				check.That(data.ResourceName).Key("notifications.0.eligible_member_activate.0.role_assignment_alert.0.critical_emails_only").HasValue("true"),
				check.That(data.ResourceName).Key("notifications.0.eligible_member_activate.0.assigned_user.0.default_recipients").HasValue("true"),
				check.That(data.ResourceName).Key("notifications.0.eligible_member_activate.0.assigned_user.0.additional_recipients.#").HasValue("1"),
				check.That(data.ResourceName).Key("notifications.0.eligible_member_activate.0.assigned_user.0.critical_emails_only").HasValue("true"),
				check.That(data.ResourceName).Key("notifications.0.eligible_member_activate.0.request_for_extension_or_approval.0.default_recipients").HasValue("true"),
				check.That(data.ResourceName).Key("notifications.0.eligible_member_activate.0.request_for_extension_or_approval.0.critical_emails_only").HasValue("true"),
			),
		},
		{
			Config: r.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				check.That(data.ResourceName).Key("role_definition_id").Exists(),
				check.That(data.ResourceName).Key("activation.0.maximum_duration_hours").HasValue("12"),
				check.That(data.ResourceName).Key("activation.0.approvers.0.group.#").HasValue("0"),
				check.That(data.ResourceName).Key("activation.0.approvers.0.user.#").HasValue("0"),
			),
		},
	})
}

func TestAccRoleManagementPolicy_allOff(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_management_policy", "test")
	r := RoleManagementPolicyResource{}

	data.ResourceTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.allOff(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				check.That(data.ResourceName).Key("role_definition_id").Exists(),

				check.That(data.ResourceName).Key("activation.0.maximum_duration_hours").HasValue("3"),
				check.That(data.ResourceName).Key("activation.0.require_multi_factor_authentication").HasValue("false"),
				check.That(data.ResourceName).Key("activation.0.require_justification").HasValue("false"),
				check.That(data.ResourceName).Key("activation.0.require_ticket_information").HasValue("false"),
				check.That(data.ResourceName).Key("activation.0.approvers.0.group.#").HasValue("0"),
				check.That(data.ResourceName).Key("activation.0.approvers.0.user.#").HasValue("0"),

				check.That(data.ResourceName).Key("assignment.0.eligible.0.expire_after_hours").HasValue("3"),
				check.That(data.ResourceName).Key("assignment.0.active.0.expire_after_days").HasValue("15"),
				check.That(data.ResourceName).Key("assignment.0.active.0.require_multi_factor_authentication").HasValue("false"),
				check.That(data.ResourceName).Key("assignment.0.active.0.require_justification").HasValue("false"),

				check.That(data.ResourceName).Key("notifications.0.member_assigned_eligible.0.role_assignment_alert.0.default_recipients").HasValue("false"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_eligible.0.role_assignment_alert.0.additional_recipients.#").HasValue("0"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_eligible.0.role_assignment_alert.0.critical_emails_only").HasValue("false"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_eligible.0.assigned_user.0.default_recipients").HasValue("false"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_eligible.0.assigned_user.0.additional_recipients.#").HasValue("0"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_eligible.0.assigned_user.0.critical_emails_only").HasValue("false"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_eligible.0.request_for_extension_or_approval.0.default_recipients").HasValue("false"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_eligible.0.request_for_extension_or_approval.0.additional_recipients.#").HasValue("0"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_eligible.0.request_for_extension_or_approval.0.critical_emails_only").HasValue("false"),

				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.role_assignment_alert.0.default_recipients").HasValue("false"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.role_assignment_alert.0.additional_recipients.#").HasValue("0"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.role_assignment_alert.0.critical_emails_only").HasValue("false"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.assigned_user.0.default_recipients").HasValue("false"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.assigned_user.0.additional_recipients.#").HasValue("0"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.assigned_user.0.critical_emails_only").HasValue("false"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.request_for_extension_or_approval.0.default_recipients").HasValue("false"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.request_for_extension_or_approval.0.additional_recipients.#").HasValue("0"),
				check.That(data.ResourceName).Key("notifications.0.member_assigned_active.0.request_for_extension_or_approval.0.critical_emails_only").HasValue("false"),

				check.That(data.ResourceName).Key("notifications.0.eligible_member_activate.0.role_assignment_alert.0.default_recipients").HasValue("false"),
				check.That(data.ResourceName).Key("notifications.0.eligible_member_activate.0.role_assignment_alert.0.additional_recipients.#").HasValue("0"),
				check.That(data.ResourceName).Key("notifications.0.eligible_member_activate.0.role_assignment_alert.0.critical_emails_only").HasValue("false"),
				check.That(data.ResourceName).Key("notifications.0.eligible_member_activate.0.assigned_user.0.default_recipients").HasValue("false"),
				check.That(data.ResourceName).Key("notifications.0.eligible_member_activate.0.assigned_user.0.additional_recipients.#").HasValue("0"),
				check.That(data.ResourceName).Key("notifications.0.eligible_member_activate.0.assigned_user.0.critical_emails_only").HasValue("false"),
				check.That(data.ResourceName).Key("notifications.0.eligible_member_activate.0.request_for_extension_or_approval.0.default_recipients").HasValue("false"),
				check.That(data.ResourceName).Key("notifications.0.eligible_member_activate.0.request_for_extension_or_approval.0.critical_emails_only").HasValue("false"),
			),
		},
	})
}

func TestAccRoleManagementPolicy_resourceGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_management_policy", "test")
	r := RoleManagementPolicyResource{}
	ri := acceptance.RandTimeInt()

	data.ResourceTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.resourceGroup(data, ri),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				check.That(data.ResourceName).Key("role_definition_id").Exists(),
			),
		},
	})
}

func TestAccRoleManagementPolicy_resource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_management_policy", "test")
	r := RoleManagementPolicyResource{}
	ri := acceptance.RandTimeInt()

	data.ResourceTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.resource(data, ri),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").Exists(),
				check.That(data.ResourceName).Key("role_definition_id").Exists(),
			),
		},
	})
}

func (r RoleManagementPolicyResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.RoleManagementPolicyId(state.ID)

	if err != nil {
		return nil, err
	}

	resp, err := client.Authorization.RoleManagementPoliciesClient.Get(ctx, id.ScopedRoleManagementPolicyId())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, nil
		}

		return nil, fmt.Errorf("loading Role Management Policy %q: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (RoleManagementPolicyResource) basic() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {}

data "azurerm_role_definition" "test" {
  name = "Monitoring Reader"
}

resource "azurerm_role_management_policy" "test" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"

  activation {
    maximum_duration_hours = 12
  }
}`

}

func (RoleManagementPolicyResource) partial1(id int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {}

data "azurerm_role_definition" "test" {
  name = "Backup Reader"
}

resource "azuread_group" "test1" {
  display_name     = "acctest-group-%[1]d1"
  security_enabled = true
}

resource "azurerm_role_management_policy" "test" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"

  activation {
    require_multi_factor_authentication = true

    approvers {
      group {
        name = azuread_group.test1.display_name
        id   = azuread_group.test1.id
      }
    }
  }

  assignment {
    eligible {
      expire_after_hours = 6
    }
    active {
      expire_after_hours = 5
      require_multi_factor_authentication = true
    }
  }

  notifications {
    member_assigned_eligible {
      role_assignment_alert {
        default_recipients = false
      }
    }

    member_assigned_active {
      role_assignment_alert {
        default_recipients = false
        additional_recipients = ["role_assignment_alert@member_assigned_active.com"]
        critical_emails_only = false
      }
    }
  }
}`, id)

}

func (RoleManagementPolicyResource) partial2(id int, password string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {}

data "azurerm_role_definition" "test" {
  name = "Backup Reader"
}

data "azuread_client_config" "test" {}

data "azuread_domains" "test" {
  only_initial = true
}

resource "azuread_user" "test1" {
  user_principal_name = "acctestUser-%[1]d1@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestUser-%[1]d1"
  password            = "%[2]s"
}

resource "azurerm_role_management_policy" "test" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"

  activation {
    require_justification = true

    approvers {
      user {
        name = azuread_user.test1.display_name
        id   = azuread_user.test1.id
      }
    }
  }

  assignment {
    eligible {
      expire_after_days = 7
    }
    active {
      expire_after_days = 4
      require_justification = true
    }
  }

  notifications {
    member_assigned_eligible {
      role_assignment_alert {
        critical_emails_only = false
      }
    }

    eligible_member_activate {
      role_assignment_alert {
        default_recipients    = false
        additional_recipients = ["role_assignment_alert@eligible_member_activate.com"]
        critical_emails_only  = false
      }
    }
  }
}`, id, password)

}

func (RoleManagementPolicyResource) managementGroup() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {}

data "azurerm_role_definition" "test" {
  name = "Monitoring Reader"
}

resource "azurerm_management_group" "test" {
}

resource "azurerm_role_management_policy" "test" {
  scope              = azurerm_management_group.test.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"

  activation {
    maximum_duration_hours = 12
  }
}`

}

func (RoleManagementPolicyResource) complete(id int, password string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {}

data "azurerm_role_definition" "test" {
  name = "Monitoring Reader"
}

data "azuread_client_config" "test" {}

data "azuread_domains" "test" {
  only_initial = true
}

resource "azuread_user" "test1" {
  user_principal_name = "acctestUser-%[1]d1@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestUser-%[1]d1"
  password            = "%[2]s"
}

resource "azuread_user" "test2" {
  user_principal_name = "acctestUser-%[1]d2@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestUser-%[1]d2"
  password            = "%[2]s"
}

resource "azuread_group" "test1" {
  display_name     = "acctest-group-%[1]d1"
  security_enabled = true
}

resource "azuread_group" "test2" {
  display_name     = "acctest-group-%[1]d2"
  security_enabled = true
}

resource "azurerm_role_management_policy" "test" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"

  activation {
    maximum_duration_hours = 4

    require_multi_factor_authentication = true
    require_justification = true
    require_ticket_information = true

    approvers {
      group {
        name = azuread_group.test1.display_name
        id   = azuread_group.test1.id
      }
      group {
        name = azuread_group.test2.display_name
        id   = azuread_group.test2.id
      }
      user {
        name = azuread_user.test1.display_name
        id   = azuread_user.test1.id
      }
      user {
        name = azuread_user.test2.display_name
        id   = azuread_user.test2.id
      }
    }
  }

  assignment {
    eligible {
      allow_permanent = true
    }

    active {
      allow_permanent = true
      require_multi_factor_authentication = true
      require_justification = true
    }
  }

  notifications {
    member_assigned_eligible {
      role_assignment_alert {
        default_recipients = true
        additional_recipients = ["role_assignment_alert@member_assigned_eligible.com"]
        critical_emails_only = true
      }
      assigned_user {
        default_recipients = true
        additional_recipients = ["assigned_user@member_assigned_eligible.com"]
        critical_emails_only = true
      }
      request_for_extension_or_approval {
        default_recipients = true
        additional_recipients = ["request_for_extension_or_approval@member_assigned_eligible.com"]
        critical_emails_only = true
      }
    }
    member_assigned_active {
      role_assignment_alert {
        default_recipients = true
        additional_recipients = ["role_assignment_alert@member_assigned_active.com"]
        critical_emails_only = true
      }
      assigned_user {
        default_recipients = true
        additional_recipients = ["assigned_user@member_assigned_active.com"]
        critical_emails_only = true
      }
      request_for_extension_or_approval {
        default_recipients = true
        additional_recipients = ["request_for_extension_or_approval@member_assigned_active.com"]
        critical_emails_only = true
      }
    }
    eligible_member_activate {
      role_assignment_alert {
        default_recipients    = true
        additional_recipients = ["role_assignment_alert@eligible_member_activate.com"]
        critical_emails_only  = true
      }
      assigned_user {
        default_recipients    = true
        additional_recipients = ["assigned_user@eligible_member_activate.com"]
        critical_emails_only  = true
      }
      request_for_extension_or_approval {
        default_recipients   = true
        critical_emails_only = true
      }
    }
  }
}

`, id, password)
}

func (RoleManagementPolicyResource) allOff() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {}

data "azurerm_role_definition" "test" {
  name = "Reader"
}

resource "azurerm_role_management_policy" "test" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"

  activation {
    maximum_duration_hours = 3

    require_multi_factor_authentication = false
    require_justification = false
    require_ticket_information = false
  }

  assignment {
    eligible {
      expire_after_hours = 3
    }

    active {
      expire_after_days = 15
      require_multi_factor_authentication = false
      require_justification = false
    }
  }

  notifications {
    member_assigned_eligible {
      role_assignment_alert {
        default_recipients = false
        additional_recipients = []
        critical_emails_only = false
      }
      assigned_user {
        default_recipients = false
        additional_recipients = []
        critical_emails_only = false
      }
      request_for_extension_or_approval {
        default_recipients = false
        additional_recipients = []
        critical_emails_only = false
      }
    }
    member_assigned_active {
      role_assignment_alert {
        default_recipients = false
        additional_recipients = []
        critical_emails_only = false
      }
      assigned_user {
        default_recipients = false
        additional_recipients = []
        critical_emails_only = false
      }
      request_for_extension_or_approval {
        default_recipients = false
        additional_recipients = []
        critical_emails_only = false
      }
    }
    eligible_member_activate {
      role_assignment_alert {
        default_recipients    = false
        additional_recipients = []
        critical_emails_only  = false
      }
      assigned_user {
        default_recipients    = false
        additional_recipients = []
        critical_emails_only  = false
      }
      request_for_extension_or_approval {
        default_recipients   = false
        critical_emails_only = false
      }
    }
  }
}`
}

func (RoleManagementPolicyResource) resourceGroup(data acceptance.TestData, id int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {}

data "azurerm_role_definition" "test" {
  name = "Monitoring Reader"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_role_management_policy" "test" {
  scope              = azurerm_resource_group.test.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"

  activation {
    maximum_duration_hours = 8

    require_multi_factor_authentication = true
    require_justification = true
    require_ticket_information = true
  }

  assignment {
    eligible {
      allow_permanent = true
    }

    active {
      allow_permanent = true
      require_multi_factor_authentication = true
      require_justification = true
    }
  }

  notifications {
    member_assigned_eligible {
      role_assignment_alert {
        default_recipients = true
        additional_recipients = ["role_assignment_alert@member_assigned_eligible.com"]
        critical_emails_only = true
      }
      assigned_user {
        default_recipients = true
        additional_recipients = ["assigned_user@member_assigned_eligible.com"]
        critical_emails_only = true
      }
      request_for_extension_or_approval {
        default_recipients = true
        additional_recipients = ["request_for_extension_or_approval@member_assigned_eligible.com"]
        critical_emails_only = true
      }
    }
    member_assigned_active {
      role_assignment_alert {
        default_recipients = true
        additional_recipients = ["role_assignment_alert@member_assigned_active.com"]
        critical_emails_only = true
      }
      assigned_user {
        default_recipients = true
        additional_recipients = ["assigned_user@member_assigned_active.com"]
        critical_emails_only = true
      }
      request_for_extension_or_approval {
        default_recipients = true
        additional_recipients = ["request_for_extension_or_approval@member_assigned_active.com"]
        critical_emails_only = true
      }
    }
    eligible_member_activate {
      role_assignment_alert {
        default_recipients    = true
        additional_recipients = ["role_assignment_alert@eligible_member_activate.com"]
        critical_emails_only  = true
      }
      assigned_user {
        default_recipients    = true
        additional_recipients = ["assigned_user@eligible_member_activate.com"]
        critical_emails_only  = true
      }
      request_for_extension_or_approval {
        default_recipients   = true
        critical_emails_only = true
      }
    }
  }
}

`, id, data.Locations.Primary)
}

func (RoleManagementPolicyResource) resource(data acceptance.TestData, id int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {}

data "azurerm_role_definition" "test" {
  name = "Monitoring Reader"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "amtestVNET1-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_role_management_policy" "test" {
  scope              = azurerm_virtual_network.test.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"

  activation {
    maximum_duration_hours = 16

    require_multi_factor_authentication = false
    require_justification = false
    require_ticket_information = false
  }

  assignment {
    eligible {
      allow_permanent = true
    }

    active {
      allow_permanent = true
      require_multi_factor_authentication = false
      require_justification = false
    }
  }
}
`, id, data.Locations.Primary)
}
