package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMDevTestPolicy_basic(t *testing.T) {
	resourceName := "azurerm_dev_test_policy.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
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

func TestAccAzureRMDevTestPolicy_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_dev_test_policy.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
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
				Config:      testAccAzureRMDevTestPolicy_requiresImport(rInt, location),
				ExpectError: testRequiresImportError("azurerm_dev_test_policy"),
			},
		},
	})
}

func TestAccAzureRMDevTestPolicy_complete(t *testing.T) {
	resourceName := "azurerm_dev_test_policy.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
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

func testCheckAzureRMDevTestPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		policyName := rs.Primary.Attributes["name"]
		policySetName := rs.Primary.Attributes["policy_set_name"]
		labName := rs.Primary.Attributes["lab_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		conn := testAccProvider.Meta().(*ArmClient).devTestLabs.PoliciesClient
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
	conn := testAccProvider.Meta().(*ArmClient).devTestLabs.PoliciesClient
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

func testAccAzureRMDevTestPolicy_requiresImport(rInt int, location string) string {
	template := testAccAzureRMDevTestPolicy_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_policy" "import" {
  name                = "${azurerm_dev_test_policy.test.name}"
  policy_set_name     = "$[azurerm_dev_test_policy.test.policy_set_name}"
  lab_name            = "${azurerm_dev_test_policy.test.lab_name}"
  resource_group_name = "${azurerm_dev_test_policy.test.resource_group_name}"
  threshold           = "999"
  evaluator_type      = "MaxValuePolicy"
}
`, template)
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
  description         = "Aloha this is the max number of VM's'"

  tags = {
    "Acceptance" = "Test"
  }
}
`, rInt, location, rInt)
}
