package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPolicyRemediation_atSubscription(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_remediation", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyRemediationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPolicyRemediation_atSubscription(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyRemediationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicyRemediation_atSubscriptionWithDefinitionSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_remediation", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyRemediationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPolicyRemediation_atSubscriptionWithDefinitionSet(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyRemediationExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scope"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicyRemediation_atResourceGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_remediation", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyRemediationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPolicyRemediation_atResourceGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyRemediationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicyRemediation_atResourceGroupWithDiscoveryMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_remediation", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyRemediationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPolicyRemediation_atResourceGroupWithDiscoveryMode(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyRemediationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicyRemediation_atManagementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_remediation", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyRemediationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPolicyRemediation_atManagementGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyRemediationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicyRemediation_atResource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_remediation", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyRemediationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPolicyRemediation_atResource(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyRemediationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicyRemediation_updateLocation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_remediation", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyRemediationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPolicyRemediation_atResourceGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyRemediationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPolicyRemediation_updateLocation(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyRemediationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicyRemediation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_remediation", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyRemediationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPolicyRemediation_atResourceGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyRemediationExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMPolicyRemediation_requiresImport),
		},
	})
}

func testCheckAzureRMPolicyRemediationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Policy.RemediationsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy Insights Remediation not found: %s", resourceName)
		}

		id, err := parse.PolicyRemediationID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := policy.RemediationGetAtScope(ctx, client, id.Name, id.PolicyScopeId); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Policy Insights Remediation %q (Scope %q) does not exist", id.Name, id.ScopeId())
			}
			return fmt.Errorf("Bad: Get on PolicyInsights.RemediationsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMPolicyRemediationDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Policy.RemediationsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_policy_remediation" {
			continue
		}

		id, err := parse.PolicyRemediationID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := policy.RemediationGetAtScope(ctx, client, id.Name, id.PolicyScopeId); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on Policy.RemediationsClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMPolicyRemediation_atSubscription(data acceptance.TestData) string {
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

func testAccAzureRMPolicyRemediation_atSubscriptionWithDefinitionSet(data acceptance.TestData) string {
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

func testAccAzureRMPolicyRemediation_atResourceGroup(data acceptance.TestData) string {
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

func testAccAzureRMPolicyRemediation_atResourceGroupWithDiscoveryMode(data acceptance.TestData) string {
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

func testAccAzureRMPolicyRemediation_updateLocation(data acceptance.TestData) string {
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

func testAccAzureRMPolicyRemediation_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMPolicyRemediation_atResourceGroup(data)
	return fmt.Sprintf(`
%s

resource "azurerm_policy_remediation" "import" {
  name                 = azurerm_policy_remediation.test.name
  scope                = azurerm_policy_remediation.test.scope
  policy_assignment_id = azurerm_policy_remediation.test.policy_assignment_id
}
`, template)
}

func testAccAzureRMPolicyRemediation_atManagementGroup(data acceptance.TestData) string {
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

func testAccAzureRMPolicyRemediation_atResource(data acceptance.TestData) string {
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
