package azurerm

import (
	"fmt"
	"github.com/hashicorp/terraform/terraform"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMPolicySetDefinition_basic(t *testing.T) {
	resourceName := "azurerm_policy_set_definition.test"

	ri := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPolicySetDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicySetDefinition_basic(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicySetDefinitionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAzureRMPolicySetDefinition_basic(ri int) string {
	return fmt.Sprintf(`
resource "azurerm_policy_set_definition" "test" {
  name         = "acctestpolset-%d"
  policy_type  = "Custom"
  display_name = "acctestpolset-%d"

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
`, ri, ri)
}

func testCheckAzureRMPolicySetDefinitionExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		policySetName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).policySetDefinitionsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, policySetName)
		if err != nil {
			return fmt.Errorf("Bad: Get on policySetDefinitionsClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("policy set does not exist: %s", name)
		}

		return nil
	}
}

func testCheckAzureRMPolicySetDefinitionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).policySetDefinitionsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_policy_set_definition" {
			continue
		}

		name := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("policy set still exists:%s", *resp.Name)
		}
	}

	return nil
}
