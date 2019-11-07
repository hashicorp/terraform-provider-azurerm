package azurerm

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/preview/blueprint/mgmt/2018-11-01-preview/blueprint"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"net/http"
	"testing"
)

func TestAccAzureRMBlueprint_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	subscriptionResourceName := "azurerm_blueprint.test_subscription"
	managementGroupResourceName := "azurerm_blueprint.test_managementGroup"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBlueprintDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBlueprint_basic_subscription(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBlueprintExists(subscriptionResourceName),
					resource.TestCheckResourceAttr(subscriptionResourceName, "name", fmt.Sprintf("acctestbp-sub-%d", ri)),
					// Check the computed values are back correctly
					testCheckAzureRMBlueprintPropertiesStatusSet(subscriptionResourceName),
				),
			},
			// Following test should have target_scope as `managementGroup` but fails an enum check in the API despite being in the spec
			// https://github.com/Azure/azure-rest-api-specs/blob/282efa7dd8301ba615d8741f740f1ed7f500fed1/specification/blueprint/resource-manager/Microsoft.Blueprint/preview/2018-11-01-preview/blueprintDefinition.json#L835
			{
				Config: testAccAzureRMBlueprint_basic_managementGroup(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBlueprintExists(managementGroupResourceName),
					resource.TestCheckResourceAttr(managementGroupResourceName, "name", fmt.Sprintf("acctestbp-mg-%d", ri)),
					// Check the computed values are back correctly
					testCheckAzureRMBlueprintPropertiesStatusSet(managementGroupResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMBlueprint_fullProperties(t *testing.T) {
	resourceName := "azurerm_blueprint.test_properties"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMBlueprint_fullProperties(ri)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBlueprintDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBlueprintExists(resourceName),
					testCheckAzureRMBlueprintPropertiesParameters(resourceName, 1),
					testCheckAzureRMBlueprintPropertiesResourceGroups(resourceName, 1),
				),
			},
		},
	})
}

func testCheckAzureRMBlueprintDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).Blueprint.BlueprintsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_blueprint" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		scope := rs.Primary.Attributes["scope"]

		resp, err := conn.Get(ctx, scope, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}
	}
	return nil
}

func testCheckAzureRMBlueprintExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		scope := rs.Primary.Attributes["scope"]

		client := testAccProvider.Meta().(*ArmClient).Blueprint.BlueprintsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, scope, name)

		if err != nil {
			return fmt.Errorf("Bad: Get on blueprintClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: blueprint %q (scope: %q) does not exist", name, scope)
		}
		return nil
	}

}

func testCheckAzureRMBlueprintPropertiesStatusSet(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		resp, err := testGetAzureRMBlueprint(s, resourceName)
		if err != nil {
			return err
		}

		created := resp.Status.TimeCreated
		lastModified := resp.Status.LastModified
		if created == nil || lastModified == nil {
			return fmt.Errorf("Bad: Computed values for blueprint %q status not found", resourceName)
		}

		return nil
	}
}

func testCheckAzureRMBlueprintPropertiesResourceGroups(resourceName string, groupCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		resp, err := testGetAzureRMBlueprint(s, resourceName)
		if err != nil {
			return err
		}

		if len(resp.ResourceGroups) != groupCount {
			return fmt.Errorf("Bad: Resource group count should have been %v, got %v", groupCount, len(resp.ResourceGroups))
		}

		return nil
	}
}

func testCheckAzureRMBlueprintPropertiesParameters(resourceName string, paramCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		resp, err := testGetAzureRMBlueprint(s, resourceName)
		if err != nil {
			return err
		}

		if len(resp.Parameters) != paramCount {
			return fmt.Errorf("Bad: Resource group count should have been %v, got %v", paramCount, len(resp.Parameters))
		}

		return nil
	}
}

func testGetAzureRMBlueprint(s *terraform.State, resourceName string) (result *blueprint.Model, err error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("Not found: %s", resourceName)
	}

	name := rs.Primary.Attributes["name"]
	scope := rs.Primary.Attributes["scope"]

	client := testAccProvider.Meta().(*ArmClient).Blueprint.BlueprintsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	bp, err := client.Get(ctx, scope, name)

	if err != nil {
		return nil, fmt.Errorf("Bad: Get on blueprintClient: %+v", err)
	}

	if bp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("Bad: Blueprint %q not found in scope %q", name, scope)
	}

	return &bp, err
}

func testAccAzureRMBlueprint_basic_subscription(ri int) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "test" {}

resource "azurerm_blueprint" "test_subscription" {
  name  = "acctestbp-sub-%d"
  scope = join("",["/subscriptions/",data.azurerm_client_config.test.subscription_id])
  type  = "Microsoft.Blueprint/blueprints"
  properties {
    description  = "accTest blueprint %d"
    display_name = "accTest blueprint"
    target_scope = "subscription"
  }
}
`, ri, ri)
}

func testAccAzureRMBlueprint_basic_managementGroup(ri int) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "test" {}

resource "azurerm_blueprint" "test_managementGroup" {
 name  = "acctestbp-mg-%d"
 scope = join("",["/providers/Microsoft.Management/managementGroups/",data.azurerm_client_config.test.tenant_id])
 type  = "Microsoft.Blueprint/blueprints"
 properties {
   description  = "accTest blueprint %d"
   display_name = "accTest blueprint"
   target_scope = "subscription"
 }
}
`, ri, ri)
}

func testAccAzureRMBlueprint_fullProperties(ri int) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "test" {}

resource "azurerm_blueprint" "test_properties" {
  name  = "acctestbp-sub-%d"
  scope = join("",["/subscriptions/",data.azurerm_client_config.test.subscription_id])
  type  = "Microsoft.Blueprint/blueprints"
  properties {
    description  = "accTest blueprint %d"
    display_name = "accTest blueprint"
    target_scope = "subscription"
	resource_groups {
      name         = "accTest-rg"
      location     = "westeurope"
      display_name = "blueprints acceptance test resource group"
      description  = "blueprints acceptance test resource group full description"
	  tags         = {
        accTest = "true"
      }
	}
    parameters {
	  name           = "accTest_Parameter"
	  type           = "string"
	  display_name   = "Acceptance Test parameter"
	  description    = "Acceptance Test parameter full description"
	  default_value  = "accTest"
	  allowed_values = ["accTest", "accTest2"]
	}
  }
}
`, ri, ri)
}
