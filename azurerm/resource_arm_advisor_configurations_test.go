package azurerm

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAdvisorConfigurations_basic(t *testing.T) {
	resourceName := "azurerm_advisor_configurations.test"

	config := testAccAzureRMAdvisorConfigurations_basic()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAdvisorConfigurationsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdvisorConfigurationsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "exclude", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMAdvisorConfigurations_complete(t *testing.T) {
	resourceName := "azurerm_advisor_configurations.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	config := testAccAzureRMAdvisorConfigurations_complete(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAdvisorConfigurationsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdvisorConfigurationsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "exclude", "false"),
					resource.TestCheckResourceAttr(resourceName, "low_cpu_threshold", "5"),
				),
			},
		},
	})
}

func testCheckAzureRMAdvisorConfigurationsExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		client := acceptance.AzureProvider.Meta().(*clients.Client).Advisor.ConfigurationsClient

		if resourceGroup, ok := rs.Primary.Attributes["resource_group_name"]; ok {
			resp, err := client.ListByResourceGroup(ctx, resourceGroup)
			if err != nil {
				return fmt.Errorf("Bad: Get on Advisor Configurations: %+v", err)
			}
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Advisor Configurations of resource group: %q does not exist", resourceGroup)
			}
		} else {
			resp, err := client.ListBySubscription(ctx)
			if err != nil {
				return fmt.Errorf("Bad: Get on Advisor Configurations: %+v", err)
			}
			if !resp.NotDone() {
				return fmt.Errorf("Bad: Advisor Configurations does not exist")
			}
		}

		return nil
	}
}

func testCheckAzureRMAdvisorConfigurationsDestroy(s *terraform.State) error {
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
	client := acceptance.AzureProvider.Meta().(*clients.Client).Advisor.ConfigurationsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_advisor_configurations" {
			continue
		}

		if resourceGroup, ok := rs.Primary.Attributes["resource_group_name"]; ok {
			resp, err := client.ListByResourceGroup(ctx, resourceGroup)
			if err != nil {
				return err
			}

			if resp.IsEmpty() || (*(*resp.Value)[0].Properties.Exclude) != true {
				return fmt.Errorf("Error deleting Advisor Configuration")
			}

			return nil
		} else {
			resp, err := client.ListBySubscription(ctx)
			if err != nil {
				return err
			}

			if !resp.NotDone() || (*resp.Values()[0].Properties.Exclude) != true {
				return fmt.Errorf("Error deleting Advisor Configuration")
			}

			return nil
		}
	}

	return nil
}

func testAccAzureRMAdvisorConfigurations_basic() string {
	return fmt.Sprintf(`
resource "azurerm_advisor_configurations" "test" {
  exclude           = false
}
`)
}

func testAccAzureRMAdvisorConfigurations_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d"
  location = "%s"
}

resource "azurerm_advisor_configurations" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  exclude             = false
  low_cpu_threshold   = "5"
}
`, rInt, location)
}
