package tests

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/policyinsights/mgmt/2019-10-01/policyinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	azpolicyinsight "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policyinsights"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPolicyInsightsRemediation_atSubscription(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_remediation", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyInsightsRemediationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPolicyInsightsRemediation_atSubscription(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyInsightsRemediationExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scope"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "policy_assignment_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicyInsightsRemediation_atResourceGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_remediation", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyInsightsRemediationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPolicyInsightsRemediation_atResourceGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyInsightsRemediationExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scope"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "policy_assignment_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicyInsightsRemediation_atManagementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_remediation", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyInsightsRemediationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPolicyInsightsRemediation_atManagementGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyInsightsRemediationExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scope"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "policy_assignment_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicyInsightsRemediation_atResource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_remediation", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyInsightsRemediationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPolicyInsightsRemediation_atResource(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyInsightsRemediationExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scope"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "policy_assignment_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMPolicyInsightsRemediationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy Insights Remediation not found: %s", resourceName)
		}

		rId, err := azpolicyinsight.ParseRemediationId(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Bad: Cannot parse remediation ID %s", rs.Primary.ID)
		}

		name := rs.Primary.Attributes["name"]
		scope := rs.Primary.Attributes["scope"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).PolicyInsights.RemediationsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		var resp policyinsights.Remediation
		switch rId.Type {
		case azpolicyinsight.AtSubscription:
			resp, err = client.GetAtSubscription(ctx, *rId.SubscriptionId, name)
		case azpolicyinsight.AtResourceGroup:
			resp, err = client.GetAtResourceGroup(ctx, *rId.SubscriptionId, *rId.ResourceGroup, name)
		case azpolicyinsight.AtManagementGroup:
			resp, err = client.GetAtManagementGroup(ctx, *rId.ManagementGroupId, name)
		case azpolicyinsight.AtResource:
			resp, err = client.GetAtResource(ctx, rId.Scope, name)
		default:
			return fmt.Errorf("Bad: Cannot recognize scope %s as Subscription ID, Resource Group ID, Management Group ID or Resource ID", scope)
		}
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Policy Insights Remediation %q (Scope %q) does not exist", name, scope)
			}
			return fmt.Errorf("Bad: Get on RemediationsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMPolicyInsightsRemediationDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).PolicyInsights.RemediationsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_policy_remediation" {
			continue
		}

		rId, err := azpolicyinsight.ParseRemediationId(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Bad: Cannot parse remediation ID %s", rs.Primary.ID)
		}

		name := rs.Primary.Attributes["name"]
		scope := rs.Primary.Attributes["scope"]

		var resp policyinsights.Remediation
		switch rId.Type {
		case azpolicyinsight.AtSubscription:
			resp, err = client.GetAtSubscription(ctx, *rId.SubscriptionId, name)
		case azpolicyinsight.AtResourceGroup:
			resp, err = client.GetAtResourceGroup(ctx, *rId.SubscriptionId, *rId.ResourceGroup, name)
		case azpolicyinsight.AtManagementGroup:
			resp, err = client.GetAtManagementGroup(ctx, *rId.ManagementGroupId, name)
		case azpolicyinsight.AtResource:
			resp, err = client.GetAtResource(ctx, rId.Scope, name)
		default:
			return fmt.Errorf("Bad: Cannot recognize scope %s as Subscription ID, Resource Group ID, Management Group ID or Resource ID", scope)
		}

		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on remediationsClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMPolicyInsightsRemediation_atSubscription(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "current" {}

resource "azurerm_policy_definition" "test" {
  name         = "acctestDef-%d"
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
  name                 = "acctestAssign-%d"
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
  name = "acctestRemediation-%d"
  scope = azurerm_policy_assignment.test.scope
  policy_assignment_id = azurerm_policy_assignment.test.id
}
`, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPolicyInsightsRemediation_atResourceGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_policy_definition" "test" {
  name         = "acctestDef-%d"
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
  name                 = "acctestAssign-%d"
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
  name = "acctestRemediation-%d"
  scope = azurerm_policy_assignment.test.scope
  policy_assignment_id = azurerm_policy_assignment.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPolicyInsightsRemediation_atManagementGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_management_group" "test" {
  display_name = "acctest-mgmt-%d"
}

resource "azurerm_policy_definition" "test" {
  name         = "acctestDef-%d"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "my-policy-definition"
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
  name                 = "acctestAssign-%d"
#   scope                = azurerm_resource_group.test.id
#   scope                = data.azurerm_subscription.current.id
  scope                = azurerm_management_group.test.id
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
  name = "acctestRemediation-%d"
  scope = azurerm_policy_assignment.test.scope
  policy_assignment_id = azurerm_policy_assignment.test.id
}
`, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPolicyInsightsRemediation_atResource(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-network-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctest-vm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_DS1_v2"

  delete_os_disk_on_termination = true
  delete_data_disks_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }
  os_profile {
    computer_name  = "hostname"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }
  os_profile_linux_config {
    disable_password_authentication = false
  }
  tags = {
    environment = "staging"
  }
}

resource "azurerm_policy_definition" "test" {
  name         = "acctestDef-%d"
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
  name                 = "acctestAssign-%d"
  scope                = azurerm_virtual_machine.test.id
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
  name = "acctestRemediation-%d"
  scope = azurerm_policy_assignment.test.scope
  policy_assignment_id = azurerm_policy_assignment.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger,
		data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
