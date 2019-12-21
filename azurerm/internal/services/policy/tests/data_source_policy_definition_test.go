package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMPolicyDefinition_builtIn(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBuiltInPolicyDefinition("Allowed resource types"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "id", "/providers/Microsoft.Authorization/policyDefinitions/a08ec900-254a-4555-9bf5-e42af04b5c5c"),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "a08ec900-254a-4555-9bf5-e42af04b5c5c"),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", "Allowed resource types"),
					resource.TestCheckResourceAttr(data.ResourceName, "type", "Microsoft.Authorization/policyDefinitions"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "This policy enables you to specify the resource types that your organization can deploy. Only resource types that support 'tags' and 'location' will be affected by this policy. To restrict all resources please duplicate this policy and change the 'mode' to 'All'."),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMPolicyDefinition_builtIn_AtManagementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBuiltInPolicyDefinitionAtManagementGroup("Allowed resource types"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "id", "/providers/Microsoft.Authorization/policyDefinitions/a08ec900-254a-4555-9bf5-e42af04b5c5c"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMPolicyDefinition_custom(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCustomPolicyDefinition(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctestpol-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", fmt.Sprintf("acctestpol-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "type", "Microsoft.Authorization/policyDefinitions"),
					resource.TestCheckResourceAttr(data.ResourceName, "policy_type", "Custom"),
					resource.TestCheckResourceAttr(data.ResourceName, "policy_rule", "{\"if\":{\"not\":{\"field\":\"location\",\"in\":\"[parameters('allowedLocations')]\"}},\"then\":{\"effect\":\"audit\"}}"),
					resource.TestCheckResourceAttr(data.ResourceName, "parameters", "{\"allowedLocations\":{\"metadata\":{\"description\":\"The list of allowed locations for resources.\",\"displayName\":\"Allowed locations\",\"strongType\":\"location\"},\"type\":\"Array\"}}"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "metadata"),
				),
			},
		},
	})
}

func testAccDataSourceBuiltInPolicyDefinition(name string) string {
	return fmt.Sprintf(`
data "azurerm_policy_definition" "test" {
  display_name = "%s"
}
`, name)
}

func testAccDataSourceBuiltInPolicyDefinitionAtManagementGroup(name string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

data "azurerm_policy_definition" "test" {
  display_name        = "%s"
  management_group_id = "${data.azurerm_client_config.current.tenant_id}"
}
`, name)
}

func testAccDataSourceCustomPolicyDefinition(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_policy_definition" "test_policy" {
  name         = "acctestpol-%d"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "acctestpol-%d"

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

  metadata = <<METADATA
  {
	"note":"azurerm acceptance test"
  }
METADATA
}

data "azurerm_policy_definition" "test" {
  display_name = "${azurerm_policy_definition.test_policy.display_name}"
}
`, data.RandomInteger, data.RandomInteger)
}
