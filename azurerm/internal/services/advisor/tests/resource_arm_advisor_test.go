package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMAdvisor_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_advisor", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAdvisorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAdvisor_basic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdvisorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAdvisor_updateBasic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdvisorExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "low_cpu_threshold", "5"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAdvisor_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_advisor", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAdvisorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAdvisor_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdvisorExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "exclude_resource_groups.#", "1"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAdvisor_updateComplete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "low_cpu_threshold", "5"),
					resource.TestCheckResourceAttr(data.ResourceName, "exclude_resource_groups.#", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMAdvisorExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		client := acceptance.AzureProvider.Meta().(*clients.Client).Advisor.ConfigurationsClient

		resp, err := client.ListBySubscriptionComplete(ctx)
		if err != nil {
			return fmt.Errorf("Bad: Get on Advisor: %+v", err)
		}
		if !resp.NotDone() {
			return fmt.Errorf("Bad: Advisor does not exist")
		}

		return nil
	}
}

func testCheckAzureRMAdvisorDestroy(s *terraform.State) error {
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
	client := acceptance.AzureProvider.Meta().(*clients.Client).Advisor.ConfigurationsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_advisor" {
			continue
		}

		resp, err := client.ListBySubscriptionComplete(ctx)
		if err != nil {
			return err
		}

		if !resp.NotDone() || !*(resp.Value().Properties.Exclude) {
			return fmt.Errorf("Error deleting Advisor Configuration")
		}
		return nil
	}
	return nil
}

func testAccAzureRMAdvisor_basic() string {
	return fmt.Sprintf(`
resource "azurerm_advisor" "test" {
}
`)
}

func testAccAzureRMAdvisor_updateBasic() string {
	return fmt.Sprintf(`
resource "azurerm_advisor" "test" {
  low_cpu_threshold = "5"
}
`)
}

func testAccAzureRMAdvisor_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d"
  location = "%s"
}

resource "azurerm_advisor" "test" {
  exclude_resource_groups = [azurerm_resource_group.test.name]
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMAdvisor_updateComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d"
  location = "%s"
}

resource "azurerm_advisor" "test" {
  low_cpu_threshold       = "5"
  exclude_resource_groups = []
}
`, data.RandomInteger, data.Locations.Primary)
}
