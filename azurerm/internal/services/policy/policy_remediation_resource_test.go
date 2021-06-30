package policy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type PolicyRemediationResource struct{}

func TestAccAzureRMPolicyRemediation_atSubscription(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_remediation", "test")
	r := PolicyRemediationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.atSubscription(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicyRemediation_atSubscriptionWithDefinitionSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_remediation", "test")
	r := PolicyRemediationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.atSubscriptionWithDefinitionSet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicyRemediation_atResourceGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_remediation", "test")
	r := PolicyRemediationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.atResourceGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicyRemediation_atResourceGroupWithDiscoveryMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_remediation", "test")
	r := PolicyRemediationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.atResourceGroupWithDiscoveryMode(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicyRemediation_atManagementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_remediation", "test")
	r := PolicyRemediationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.atManagementGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicyRemediation_atResource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_remediation", "test")
	r := PolicyRemediationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.atResource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicyRemediation_updateLocation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_remediation", "test")
	r := PolicyRemediationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.atResourceGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateLocation(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicyRemediation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_remediation", "test")
	r := PolicyRemediationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.atResourceGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r PolicyRemediationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.PolicyRemediationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := policy.RemediationGetAtScope(ctx, client.Policy.RemediationsClient, id.Name, id.PolicyScopeId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Policy Remediation %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.RemediationProperties != nil), nil
}

func (r PolicyRemediationResource) atSubscription(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_policy_definition" "test" {
  name         = "acctestDef-%[1]s"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "my-policy-definition"

  policy_rule = <<POLICY_RULE
    {
    "if": {
      "not": {
        "field": "location",
        "in": "[parameters('allowedLocations')]"
      }
    },
    "then": {
      "effect": "audit"
    }
  }
POLICY_RULE

  parameters = <<PARAMETERS
    {
    "allowedLocations": {
      "type": "Array",
      "metadata": {
        "description": "The list of allowed locations for resources.",
        "displayName": "Allowed locations",
        "strongType": "location"
      }
    }
  }
PARAMETERS
}

resource "azurerm_policy_assignment" "test" {
  name                 = "acctestAssign-%[1]s"
  scope                = data.azurerm_subscription.current.id
  policy_definition_id = azurerm_policy_definition.test.id
  description          = "Policy Assignment created via an Acceptance Test"
  display_name         = "My Example Policy Assignment"

  parameters = <<PARAMETERS
{
  "allowedLocations": {
    "value": [ "West Europe" ]
  }
}
PARAMETERS
}

resource "azurerm_policy_remediation" "test" {
  name                 = "acctestremediation-%[1]s"
  scope                = azurerm_policy_assignment.test.scope
  policy_assignment_id = azurerm_policy_assignment.test.id
}
`, data.RandomString)
}

func (r PolicyRemediationResource) atSubscriptionWithDefinitionSet(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_policy_set_definition" "test" {
  name         = "testPolicySet-%[1]s"
  policy_type  = "Custom"
  display_name = "testPolicySet-%[1]s"

  parameters = <<PARAMETERS
    {
        "allowedLocations": {
            "type": "Array",
            "metadata": {
                "description": "The list of allowed locations for resources.",
                "displayName": "Allowed locations",
                "strongType": "location"
            }
        }
    }
PARAMETERS

  policy_definitions = <<POLICY_DEFINITIONS
    [
        {
            "parameters": {
                "listOfAllowedLocations": {
                    "value": "[parameters('allowedLocations')]"
                }
            },
            "policyDefinitionId": "/providers/Microsoft.Authorization/policyDefinitions/e765b5de-1225-4ba3-bd56-1ac6695af988"
        }
    ]
POLICY_DEFINITIONS
}

resource "azurerm_policy_definition" "test" {
  name         = "acctestDef-%[1]s"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "acctestDef-%[1]s"

  policy_rule = <<POLICY_RULE
    {
    "if": {
      "not": {
        "field": "location",
        "in": "[parameters('allowedLocations')]"
      }
    },
    "then": {
      "effect": "audit"
    }
  }
POLICY_RULE

  parameters = <<PARAMETERS
    {
    "allowedLocations": {
      "type": "Array",
      "metadata": {
        "description": "The list of allowed locations for resources.",
        "displayName": "Allowed locations",
        "strongType": "location"
      }
    }
  }
PARAMETERS
}

resource "azurerm_policy_assignment" "test" {
  name                 = "acctestAssign-%[1]s"
  scope                = data.azurerm_subscription.current.id
  policy_definition_id = azurerm_policy_set_definition.test.id
  description          = "Policy Assignment created via an Acceptance Test"
  display_name         = "acctestAssign-%[1]s"

  parameters = <<PARAMETERS
{
  "allowedLocations": {
    "value": [ "West Europe" ]
  }
}
PARAMETERS
}

resource "azurerm_policy_remediation" "test" {
  name                           = "acctestremediation-%[1]s"
  scope                          = azurerm_policy_assignment.test.scope
  policy_assignment_id           = azurerm_policy_assignment.test.id
  policy_definition_reference_id = azurerm_policy_definition.test.id
}
`, data.RandomString)
}

func (r PolicyRemediationResource) atResourceGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-policy-%[1]s"
  location = "%[2]s"
}

resource "azurerm_policy_definition" "test" {
  name         = "acctestDef-%[1]s"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "my-policy-definition"

  policy_rule = <<POLICY_RULE
    {
    "if": {
      "not": {
        "field": "location",
        "in": "[parameters('allowedLocations')]"
      }
    },
    "then": {
      "effect": "audit"
    }
  }
POLICY_RULE

  parameters = <<PARAMETERS
    {
    "allowedLocations": {
      "type": "Array",
      "metadata": {
        "description": "The list of allowed locations for resources.",
        "displayName": "Allowed locations",
        "strongType": "location"
      }
    }
  }
PARAMETERS
}

resource "azurerm_policy_assignment" "test" {
  name                 = "acctestAssign-%[1]s"
  scope                = azurerm_resource_group.test.id
  policy_definition_id = azurerm_policy_definition.test.id
  description          = "Policy Assignment created via an Acceptance Test"
  display_name         = "acctestAssign-%[1]s"

  parameters = <<PARAMETERS
{
  "allowedLocations": {
    "value": [ "West Europe" ]
  }
}
PARAMETERS
}

resource "azurerm_policy_remediation" "test" {
  name                 = "acctestremediation-%[1]s"
  scope                = azurerm_policy_assignment.test.scope
  policy_assignment_id = azurerm_policy_assignment.test.id
}
`, data.RandomString, data.Locations.Primary)
}

func (r PolicyRemediationResource) atResourceGroupWithDiscoveryMode(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-policy-%[1]s"
  location = "%[2]s"
}

resource "azurerm_policy_definition" "test" {
  name         = "acctestDef-%[1]s"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "my-policy-definition"

  policy_rule = <<POLICY_RULE
    {
    "if": {
      "not": {
        "field": "location",
        "in": "[parameters('allowedLocations')]"
      }
    },
    "then": {
      "effect": "audit"
    }
  }
POLICY_RULE

  parameters = <<PARAMETERS
    {
    "allowedLocations": {
      "type": "Array",
      "metadata": {
        "description": "The list of allowed locations for resources.",
        "displayName": "Allowed locations",
        "strongType": "location"
      }
    }
  }
PARAMETERS
}

resource "azurerm_policy_assignment" "test" {
  name                 = "acctestAssign-%[1]s"
  scope                = azurerm_resource_group.test.id
  policy_definition_id = azurerm_policy_definition.test.id
  description          = "Policy Assignment created via an Acceptance Test"
  display_name         = "My Example Policy Assignment"

  parameters = <<PARAMETERS
{
  "allowedLocations": {
    "value": [ "West Europe" ]
  }
}
PARAMETERS
}

resource "azurerm_policy_remediation" "test" {
  name                 = "acctestremediation-%[1]s"
  scope                = azurerm_policy_assignment.test.scope
  policy_assignment_id = azurerm_policy_assignment.test.id

  resource_discovery_mode = "ReEvaluateCompliance"
}
`, data.RandomString, data.Locations.Primary)
}

func (r PolicyRemediationResource) updateLocation(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-policy-%[1]s"
  location = "%[2]s"
}

resource "azurerm_policy_definition" "test" {
  name         = "acctestDef-%[1]s"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "my-policy-definition"

  policy_rule = <<POLICY_RULE
    {
    "if": {
      "not": {
        "field": "location",
        "in": "[parameters('allowedLocations')]"
      }
    },
    "then": {
      "effect": "audit"
    }
  }
POLICY_RULE

  parameters = <<PARAMETERS
    {
    "allowedLocations": {
      "type": "Array",
      "metadata": {
        "description": "The list of allowed locations for resources.",
        "displayName": "Allowed locations",
        "strongType": "location"
      }
    }
  }
PARAMETERS
}

resource "azurerm_policy_assignment" "test" {
  name                 = "acctestAssign-%[1]s"
  scope                = azurerm_resource_group.test.id
  policy_definition_id = azurerm_policy_definition.test.id
  description          = "Policy Assignment created via an Acceptance Test"
  display_name         = "My Example Policy Assignment"

  parameters = <<PARAMETERS
{
  "allowedLocations": {
    "value": [ "West Europe" ]
  }
}
PARAMETERS
}

resource "azurerm_policy_remediation" "test" {
  name                 = "acctestremediation-%[1]s"
  scope                = azurerm_policy_assignment.test.scope
  policy_assignment_id = azurerm_policy_assignment.test.id
  location_filters     = ["westus"]
}
`, data.RandomString, data.Locations.Primary)
}

func (r PolicyRemediationResource) requiresImport(data acceptance.TestData) string {
	template := r.atResourceGroup(data)
	return fmt.Sprintf(`
%s

resource "azurerm_policy_remediation" "import" {
  name                 = azurerm_policy_remediation.test.name
  scope                = azurerm_policy_remediation.test.scope
  policy_assignment_id = azurerm_policy_remediation.test.policy_assignment_id
}
`, template)
}

func (r PolicyRemediationResource) atManagementGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
  display_name = "acctest-policy-%[1]s"
}

resource "azurerm_policy_definition" "test" {
  name                = "acctestDef-%[1]s"
  policy_type         = "Custom"
  mode                = "All"
  display_name        = "my-policy-definition"
  management_group_id = azurerm_management_group.test.group_id

  policy_rule = <<POLICY_RULE
    {
    "if": {
      "not": {
        "field": "location",
        "in": "[parameters('allowedLocations')]"
      }
    },
    "then": {
      "effect": "audit"
    }
  }
POLICY_RULE

  parameters = <<PARAMETERS
    {
    "allowedLocations": {
      "type": "Array",
      "metadata": {
        "description": "The list of allowed locations for resources.",
        "displayName": "Allowed locations",
        "strongType": "location"
      }
    }
  }
PARAMETERS
}

resource "azurerm_policy_assignment" "test" {
  name = "acctestAssign-%[1]s"
  #   scope                = azurerm_resource_group.test.id
  #   scope                = data.azurerm_subscription.current.id
  scope = azurerm_management_group.test.id
  # scope                = azurerm_virtual_machine.main.id
  policy_definition_id = azurerm_policy_definition.test.id
  description          = "Policy Assignment created via an Acceptance Test"
  display_name         = "My Example Policy Assignment"

  parameters = <<PARAMETERS
{
  "allowedLocations": {
    "value": [ "West Europe" ]
  }
}
PARAMETERS
}

resource "azurerm_policy_remediation" "test" {
  name                 = "acctestremediation-%[1]s"
  scope                = azurerm_policy_assignment.test.scope
  policy_assignment_id = azurerm_policy_assignment.test.id
}
`, data.RandomString)
}

func (r PolicyRemediationResource) atResource(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-policy-%[1]s"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-network-%[1]s"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[1]s"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%[1]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "acctest-vm-%[1]s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  admin_password                  = "P@ssw0rd1234!"
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }
}

resource "azurerm_policy_definition" "test" {
  name         = "acctestDef-%[1]s"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "my-policy-definition"

  policy_rule = <<POLICY_RULE
    {
    "if": {
      "not": {
        "field": "location",
        "in": "[parameters('allowedLocations')]"
      }
    },
    "then": {
      "effect": "audit"
    }
  }
POLICY_RULE

  parameters = <<PARAMETERS
    {
    "allowedLocations": {
      "type": "Array",
      "metadata": {
        "description": "The list of allowed locations for resources.",
        "displayName": "Allowed locations",
        "strongType": "location"
      }
    }
  }
PARAMETERS
}

resource "azurerm_policy_assignment" "test" {
  name                 = "acctestAssign-%[1]s"
  scope                = azurerm_linux_virtual_machine.test.id
  policy_definition_id = azurerm_policy_definition.test.id
  description          = "Policy Assignment created via an Acceptance Test"
  display_name         = "My Example Policy Assignment"

  parameters = <<PARAMETERS
{
  "allowedLocations": {
    "value": [ "West Europe" ]
  }
}
PARAMETERS
}

resource "azurerm_policy_remediation" "test" {
  name                 = "acctestremediation-%[1]s"
  scope                = azurerm_policy_assignment.test.scope
  policy_assignment_id = azurerm_policy_assignment.test.id
}
`, data.RandomString, data.Locations.Primary)
}
