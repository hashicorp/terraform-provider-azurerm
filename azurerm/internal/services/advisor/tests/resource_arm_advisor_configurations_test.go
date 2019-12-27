package tests

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAdvisorConfigurations_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_advisor_configurations", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAdvisorConfigurationsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAdvisorConfigurations_basic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdvisorConfigurationsExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "exclude", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMAdvisorConfigurations_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_advisor_configurations", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAdvisorConfigurationsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAdvisorConfigurations_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdvisorConfigurationsExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "exclude", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "low_cpu_threshold", "5"),
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

func testAccAzureRMAdvisorConfigurations_complete(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary)
}
