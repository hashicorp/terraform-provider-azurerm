// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	assignments "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-06-01/policyassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ResourceAssignmentTestResource struct{}

func TestAccResourcePolicyAssignment_basicWithBuiltInPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_policy_assignment", "test")
	r := ResourceAssignmentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withBuiltInPolicyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBuiltInPolicyUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBuiltInPolicyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourcePolicyAssignment_basicWithBuiltInPolicyNonComplianceMessage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_policy_assignment", "test")
	r := ResourceAssignmentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withBuiltInPolicyNonComplianceMessage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("non_compliance_message.#").HasValue("1"),
				check.That(data.ResourceName).Key("non_compliance_message.0.content").HasValue("test"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBuiltInPolicyNonComplianceMessageUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("non_compliance_message.0").DoesNotExist(),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBuiltInPolicyNonComplianceMessage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("non_compliance_message.#").HasValue("1"),
				check.That(data.ResourceName).Key("non_compliance_message.0.content").HasValue("test"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourcePolicyAssignment_basicWithBuiltInPolicySet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_policy_assignment", "test")
	r := ResourceAssignmentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withBuiltInPolicySetBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBuiltInPolicySetUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBuiltInPolicySetBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourcePolicyAssignment_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_policy_assignment", "test")
	r := ResourceAssignmentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.userAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourcePolicyAssignment_basicWithBuiltInPolicySetNonComplianceMessage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_policy_assignment", "test")
	r := ResourceAssignmentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withBuiltInPolicySetNonComplianceMessage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("non_compliance_message.#").HasValue("1"),
				check.That(data.ResourceName).Key("non_compliance_message.0.content").HasValue("test"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBuiltInPolicySetNonComplianceMessageUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("non_compliance_message.#").HasValue("2"),
				check.That(data.ResourceName).Key("non_compliance_message.0.content").HasValue("test"),
				check.That(data.ResourceName).Key("non_compliance_message.1.content").HasValue("test2"),
				check.That(data.ResourceName).Key("non_compliance_message.1.policy_definition_reference_id").HasValue("AINE_MinimumPasswordLength"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBuiltInPolicySetNonComplianceMessage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("non_compliance_message.#").HasValue("1"),
				check.That(data.ResourceName).Key("non_compliance_message.0.content").HasValue("test"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourcePolicyAssignment_basicWithCustomPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_policy_assignment", "test")
	r := ResourceAssignmentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withCustomPolicyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withCustomPolicyUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withCustomPolicyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourcePolicyAssignment_basicWithCustomPolicyRequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_policy_assignment", "test")
	r := ResourceAssignmentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withCustomPolicyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		data.RequiresImportErrorStep(r.withCustomPolicyRequiresImport),
	})
}

func TestAccResourcePolicyAssignment_basicWithCustomPolicyComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_policy_assignment", "test")
	r := ResourceAssignmentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withCustomPolicyComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withCustomPolicyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withCustomPolicyComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourcePolicyAssignment_overridesAndResourceSelector(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_policy_assignment", "test")
	r := ResourceAssignmentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withOverrideAndSelectorsBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withOverrideAndSelectorsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withOverrideAndSelectorsBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ResourceAssignmentTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := assignments.ParseScopedPolicyAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}

	assignment, err := client.Policy.AssignmentsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(assignment.HttpResponse) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (r ResourceAssignmentTestResource) withBuiltInPolicyBasic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_definition" "test" {
  display_name = "Allowed locations"
}

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = data.azurerm_policy_definition.test.id
  parameters = jsonencode({
    "listOfAllowedLocations" = {
      "value" = [azurerm_resource_group.test.location]
    }
  })
}
`, template, data.RandomInteger)
}

func (r ResourceAssignmentTestResource) withBuiltInPolicyUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_definition" "test" {
  display_name = "Allowed locations"
}

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = data.azurerm_policy_definition.test.id
  parameters = jsonencode({
    "listOfAllowedLocations" = {
      "value" = [azurerm_resource_group.test.location, "%[3]s"]
    }
  })
}
`, template, data.RandomInteger, data.Locations.Secondary)
}

func (r ResourceAssignmentTestResource) withBuiltInPolicyNonComplianceMessage(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_definition" "test" {
  display_name = "Allowed locations"
}

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = data.azurerm_policy_definition.test.id

  non_compliance_message {
    content = "test"
  }

  parameters = jsonencode({
    "listOfAllowedLocations" = {
      "value" = [azurerm_resource_group.test.location]
    }
  })
}
`, template, data.RandomInteger)
}

func (r ResourceAssignmentTestResource) withBuiltInPolicyNonComplianceMessageUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_definition" "test" {
  display_name = "Allowed locations"
}

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = data.azurerm_policy_definition.test.id
  parameters = jsonencode({
    "listOfAllowedLocations" = {
      "value" = [azurerm_resource_group.test.location, "%[3]s"]
    }
  })
}
`, template, data.RandomInteger, data.Locations.Secondary)
}

func (r ResourceAssignmentTestResource) withBuiltInPolicySetBasic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_set_definition" "test" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  location             = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func (r ResourceAssignmentTestResource) withBuiltInPolicySetUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_set_definition" "test" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  location             = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }

  metadata = jsonencode({
    "category" : "Testing"
  })
}
`, template, data.RandomInteger)
}

func (r ResourceAssignmentTestResource) withBuiltInPolicySetNonComplianceMessage(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_set_definition" "test" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  location             = azurerm_resource_group.test.location

  non_compliance_message {
    content = "test"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func (r ResourceAssignmentTestResource) withBuiltInPolicySetNonComplianceMessageUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_set_definition" "test" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  location             = azurerm_resource_group.test.location

  non_compliance_message {
    content = "test"
  }

  non_compliance_message {
    content                        = "test2"
    policy_definition_reference_id = "AINE_MinimumPasswordLength"
  }

  identity {
    type = "SystemAssigned"
  }

  metadata = jsonencode({
    "category" : "Testing"
  })
}
`, template, data.RandomInteger)
}

func (r ResourceAssignmentTestResource) withCustomPolicyBasic(data acceptance.TestData) string {
	template := r.templateWithCustomPolicy(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = azurerm_policy_definition.test.id
}
`, template, data.RandomInteger)
}

func (r ResourceAssignmentTestResource) withOverrideAndSelectorsBasic(data acceptance.TestData) string {
	template := r.templateWithCustomPolicy(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_set_definition" "test" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  metadata = jsonencode({
    "category" : "Testing"
  })

  overrides {
    value = "Disabled"
    selectors {
      in = [data.azurerm_policy_set_definition.test.policy_definition_reference.0.reference_id]
    }
  }

  resource_selectors {
    selectors {
      not_in = ["eastus", "westus"]
      kind   = "resourceLocation"
    }
  }
}
`, template, data.RandomInteger)
}

func (r ResourceAssignmentTestResource) withOverrideAndSelectorsUpdate(data acceptance.TestData) string {
	template := r.templateWithCustomPolicy(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_set_definition" "test" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  metadata = jsonencode({
    "category" : "Testing"
  })

  overrides {
    value = "AuditIfNotExists"
    selectors {
      in = [data.azurerm_policy_set_definition.test.policy_definition_reference.0.reference_id]
    }
  }

  resource_selectors {
    selectors {
      not_in = ["eastus"]
      kind   = "resourceLocation"
    }
  }
}
`, template, data.RandomInteger)
}

func (r ResourceAssignmentTestResource) withCustomPolicyComplete(data acceptance.TestData) string {
	template := r.templateWithCustomPolicy(data)
	// NOTE: we could include parameters here but it's tested extensively elsewhere
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = azurerm_policy_definition.test.id
  description          = "This is a policy assignment from an acceptance test"
  display_name         = "AccTest Policy %[2]d"
  enforce              = false
  not_scopes = [
    format("%%s/subnets/subnet1", azurerm_virtual_network.test.id)
  ]
  metadata = jsonencode({
    "category" : "Testing"
  })
}
`, template, data.RandomInteger)
}

func (r ResourceAssignmentTestResource) withCustomPolicyRequiresImport(data acceptance.TestData) string {
	template := r.withCustomPolicyBasic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_policy_assignment" "import" {
  name                 = azurerm_resource_policy_assignment.test.name
  resource_id          = azurerm_resource_policy_assignment.test.resource_id
  policy_definition_id = azurerm_resource_policy_assignment.test.policy_definition_id
}
`, template)
}

func (r ResourceAssignmentTestResource) withCustomPolicyUpdated(data acceptance.TestData) string {
	template := r.templateWithCustomPolicy(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = azurerm_policy_definition.test.id
  metadata = jsonencode({
    "category" : "Testing"
  })
}
`, template, data.RandomInteger)
}

func (r ResourceAssignmentTestResource) templateWithCustomPolicy(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_policy_definition" "test" {
  name         = "acctestpol-%[2]d"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "acctestpol-%[2]d"

  policy_rule = <<POLICY_RULE
	{
    "if": {
      "not": {
        "field": "name",
        "equals": "bob"
      }
    },
    "then": {
      "effect": "audit"
    }
  }
POLICY_RULE
}
`, template, data.RandomInteger)
}

func (r ResourceAssignmentTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest%[1]d"
  location = %[2]q
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ResourceAssignmentTestResource) systemAssignedIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_set_definition" "test" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  location             = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func (r ResourceAssignmentTestResource) userAssignedIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_set_definition" "test" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestua%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  location             = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, template, data.RandomInteger)
}
