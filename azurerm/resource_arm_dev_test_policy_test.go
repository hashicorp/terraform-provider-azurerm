package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMDevTestPolicy_basic(t *testing.T) {
	resourceName := "azurerm_dev_test_policy.test"
	rInt := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDevTestPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDevTestPolicy_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func TestAccAzureRMDevTestPolicy_complete(t *testing.T) {
	resourceName := "azurerm_dev_test_policy.test"
	rInt := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDevTestPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDevTestPolicy_complete(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Acceptance", "Test"),
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

func testCheckAzureRMDevTestPolicyExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		policyName := rs.Primary.Attributes["name"]
		policySetName := rs.Primary.Attributes["policy_set_name"]
		labName := rs.Primary.Attributes["lab_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		conn := testAccProvider.Meta().(*ArmClient).devTestPoliciesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, labName, policySetName, policyName, "")
		if err != nil {
			return fmt.Errorf("Bad: Get devTestPoliciesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DevTest Policy %q (Policy Set %q / Lab %q / Resource Group: %q) does not exist", policyName, policySetName, labName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDevTestPolicyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).devTestPoliciesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dev_test_policy" {
			continue
		}

		policyName := rs.Primary.Attributes["name"]
		policySetName := rs.Primary.Attributes["policy_set_name"]
		labName := rs.Primary.Attributes["lab_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, labName, policySetName, policyName, "")

		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("DevTest Policy still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMDevTestPolicy_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dev_test_policy" "test" {
  name                = "LabVmCount"
  policy_set_name     = "default"
  lab_name            = "${azurerm_dev_test_lab.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  threshold           = "999"
  evaluator_type      = "MaxValuePolicy"
}
`, rInt, location, rInt)
}

func testAccAzureRMDevTestPolicy_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dev_test_policy" "test" {
  name                = "LabVmCount"
  policy_set_name     = "default"
  lab_name            = "${azurerm_dev_test_lab.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  threshold           = "999"
  evaluator_type      = "MaxValuePolicy"

  tags {
    "Acceptance" = "Test"
  }
}
`, rInt, location, rInt)
}
