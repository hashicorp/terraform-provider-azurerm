package devtestlabs_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccDevTestPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevTestPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevTestPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccDevTestPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevTestPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevTestPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			{
				Config:      testAccDevTestPolicy_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_dev_test_policy"),
			},
		},
	})
}

func TestAccDevTestPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevTestPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevTestPolicy_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Acceptance", "Test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckDevTestPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).DevTestLabs.PoliciesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		policyName := rs.Primary.Attributes["name"]
		policySetName := rs.Primary.Attributes["policy_set_name"]
		labName := rs.Primary.Attributes["lab_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

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

func testCheckDevTestPolicyDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).DevTestLabs.PoliciesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		policyName := rs.Primary.Attributes["name"]
		policySetName := rs.Primary.Attributes["policy_set_name"]
		labName := rs.Primary.Attributes["lab_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Delete(ctx, resourceGroup, labName, policySetName, policyName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Failed deleting policy - it wasn't found %q: %+v", resourceGroup, err)
			}
			return fmt.Errorf("Failed deleting policy %q: %+v", resourceGroup, err)
		}

		// must be a better way
		time.Sleep(15 * time.Second)

		// check the policy is NOT there
		resp1, err := conn.Get(ctx, resourceGroup, labName, policySetName, policyName, "")
		if err != nil {
			if resp1.StatusCode == http.StatusNotFound {
				return nil
			}
			return fmt.Errorf("Bad: Get devTestPoliciesClient: %+v", err)
		}

		return nil
	}
}

func testCheckDevTestPolicyDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).DevTestLabs.PoliciesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccDevTestPolicy_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dev_test_policy" "test" {
  name                = "LabVmCount"
  policy_set_name     = "default"
  lab_name            = azurerm_dev_test_lab.test.name
  resource_group_name = azurerm_resource_group.test.name
  threshold           = "999"
  evaluator_type      = "MaxValuePolicy"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccDevTestPolicy_requiresImport(data acceptance.TestData) string {
	template := testAccDevTestPolicy_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_policy" "import" {
  name                = azurerm_dev_test_policy.test.name
  policy_set_name     = "$[azurerm_dev_test_policy.test.policy_set_name}"
  lab_name            = azurerm_dev_test_policy.test.lab_name
  resource_group_name = azurerm_dev_test_policy.test.resource_group_name
  threshold           = "999"
  evaluator_type      = "MaxValuePolicy"
}
`, template)
}

func testAccDevTestPolicy_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dev_test_policy" "test" {
  name                = "LabVmCount"
  policy_set_name     = "default"
  lab_name            = azurerm_dev_test_lab.test.name
  resource_group_name = azurerm_resource_group.test.name
  threshold           = "999"
  evaluator_type      = "MaxValuePolicy"
  description         = "Aloha this is the max number of VM's'"

  tags = {
    "Acceptance" = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
