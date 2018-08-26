package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMPolicyAssignment_basic(t *testing.T) {
	resourceName := "azurerm_policy_assignment.test"

	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPolicyAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyAssignment_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyAssignmentExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMPolicyAssignment_complete(t *testing.T) {
	resourceName := "azurerm_policy_assignment.test"

	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPolicyAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyAssignment_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyAssignmentExists(resourceName),
				),
			},
		},
	})
}

func testCheckAzureRMPolicyAssignmentExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).policyAssignmentsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		id := rs.Primary.ID
		resp, err := client.GetByID(ctx, id)
		if err != nil {
			return fmt.Errorf("Bad: Get on policyAssignmentsClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Policy Assignment does not exist: %s", name)
		}

		return nil
	}
}

func testCheckAzureRMPolicyAssignmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).policyAssignmentsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_policy_definition" {
			continue
		}

		id := rs.Primary.ID
		resp, err := client.GetByID(ctx, id)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Policy Assignment still exists:%s", *resp.Name)
		}
	}

	return nil
}

func testAzureRMPolicyAssignment_basic(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_policy_definition" "test" {
  name         = "acctestpol-%d"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "acctestpol-%d"
  policy_rule  = <<POLICY_RULE
	{
    "if": {
      "not": {
        "field": "location",
        "equals": "%s"
      }
    },
    "then": {
      "effect": "audit"
    }
  }
POLICY_RULE
}

resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_policy_assignment" "test" {
  name                 = "acctestpa-%d"
  scope                = "${azurerm_resource_group.test.id}"
  policy_definition_id = "${azurerm_policy_definition.test.id}"
}
`, ri, ri, location, ri, location, ri)
}

func testAzureRMPolicyAssignment_complete(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_policy_definition" "test" {
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
}

resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_policy_assignment" "test" {
  name                 = "acctestpa-%d"
  scope                = "${azurerm_resource_group.test.id}"
  policy_definition_id = "${azurerm_policy_definition.test.id}"
  description          = "Policy Assignment created via an Acceptance Test"
  display_name         = "Acceptance Test Run %d"
  parameters = <<PARAMETERS
{
  "allowedLocations": {
    "value": [ "%s" ]
  }
}
PARAMETERS
}
`, ri, ri, ri, location, ri, ri, location)
}
