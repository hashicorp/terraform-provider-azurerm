package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMPolicyDefinition_builtIn(t *testing.T) {
	dataSourceName := "data.azurerm_policy_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBuiltInPolicyDefinition("Allowed resource types"),
				Check: resource.ComposeTestCheckFunc(
					testAzureRMClientConfigAttr(dataSourceName, "id", "/providers/Microsoft.Authorization/policyDefinitions/a08ec900-254a-4555-9bf5-e42af04b5c5c"),
					testAzureRMClientConfigAttr(dataSourceName, "name", "a08ec900-254a-4555-9bf5-e42af04b5c5c"),
					testAzureRMClientConfigAttr(dataSourceName, "display_name", "Allowed resource types"),
					testAzureRMClientConfigAttr(dataSourceName, "type", "Microsoft.Authorization/policyDefinitions"),
					testAzureRMClientConfigAttr(dataSourceName, "description", "This policy enables you to specify the resource types that your organization can deploy."),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMPolicyDefinition_builtIn_AtManagementGroup(t *testing.T) {
	dataSourceName := "data.azurerm_policy_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBuiltInPolicyDefinitionAtManagementGroup("Allowed resource types"),
				Check: resource.ComposeTestCheckFunc(
					testAzureRMClientConfigAttr(dataSourceName, "id", "/providers/Microsoft.Authorization/policyDefinitions/a08ec900-254a-4555-9bf5-e42af04b5c5c"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMPolicyDefinition_custom(t *testing.T) {
	ri := tf.AccRandTimeInt()
	dataSourceName := "data.azurerm_policy_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCustomPolicyDefinition(ri),
				Check: resource.ComposeTestCheckFunc(
					testAzureRMAttrExists(dataSourceName, "id"),
					testAzureRMClientConfigAttr(dataSourceName, "name", fmt.Sprintf("acctestpol-%d", ri)),
					testAzureRMClientConfigAttr(dataSourceName, "display_name", fmt.Sprintf("acctestpol-%d", ri)),
					testAzureRMClientConfigAttr(dataSourceName, "type", "Microsoft.Authorization/policyDefinitions"),
					testAzureRMClientConfigAttr(dataSourceName, "policy_type", "Custom"),
					testAzureRMClientConfigAttr(dataSourceName, "policy_rule", "{\"if\":{\"not\":{\"field\":\"location\",\"in\":\"[parameters('allowedLocations')]\"}},\"then\":{\"effect\":\"audit\"}}"),
					testAzureRMClientConfigAttr(dataSourceName, "parameters", "{\"allowedLocations\":{\"metadata\":{\"description\":\"The list of allowed locations for resources.\",\"displayName\":\"Allowed locations\",\"strongType\":\"location\"},\"type\":\"Array\"}}"),
					testAzureRMClientConfigAttr(dataSourceName, "metadata", "{\"note\":\"azurerm acceptance test\"}"),
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

func testAccDataSourceCustomPolicyDefinition(ri int) string {
	return fmt.Sprintf(`
resource "azurerm_policy_definition" "test_policy" {
  name         = "acctestpol-%d"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "acctestpol-%d"

  policy_rule  = <<POLICY_RULE
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

`, ri, ri)
}

func testAzureRMAttrExists(name, key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return resource.TestCheckResourceAttrSet(name, key)(s)
	}
}
