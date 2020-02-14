package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMPolicyRemediation_atSubscription(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_remediation", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePolicyRemediation_atSubscription(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "scope"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "policy_assignment_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMPolicyRemediation_atResourceGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_remediation", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePolicyRemediation_atResourceGroup(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "scope"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "policy_assignment_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMPolicyRemediation_atManagementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_remediation", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePolicyRemediation_atManagementGroup(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "scope"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "policy_assignment_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMPolicyRemediation_atResource(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_remediation", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePolicyRemediation_atResource(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "scope"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "policy_assignment_id"),
				),
			},
		},
	})
}

func testAccDataSourcePolicyRemediation_atSubscription(data acceptance.TestData) string {
	config := testAccAzureRMPolicyRemediation_atSubscription(data)
	return fmt.Sprintf(`
%s

data "azurerm_policy_remediation" "test" {
  name  = azurerm_policy_remediation.test.name
  scope = data.azurerm_subscription.current.id
}
`, config)
}

func testAccDataSourcePolicyRemediation_atResourceGroup(data acceptance.TestData) string {
	config := testAccAzureRMPolicyRemediation_atResourceGroup(data)
	return fmt.Sprintf(`
%s

data "azurerm_policy_remediation" "test" {
  name  = azurerm_policy_remediation.test.name
  scope = azurerm_resource_group.test.id
}
`, config)
}

func testAccDataSourcePolicyRemediation_atManagementGroup(data acceptance.TestData) string {
	config := testAccAzureRMPolicyRemediation_atManagementGroup(data)
	return fmt.Sprintf(`
%s

data "azurerm_policy_remediation" "test" {
  name  = azurerm_policy_remediation.test.name
  scope = azurerm_management_group.test.id
}
`, config)
}

func testAccDataSourcePolicyRemediation_atResource(data acceptance.TestData) string {
	config := testAccAzureRMPolicyRemediation_atResource(data)
	return fmt.Sprintf(`
%s

data "azurerm_policy_remediation" "test" {
  name  = azurerm_policy_remediation.test.name
  scope = azurerm_virtual_machine.test.id
}
`, config)
}
