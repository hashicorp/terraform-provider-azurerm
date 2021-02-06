package consumption

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMConsumptionBudgetResourceGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_consumption_budget_resource_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMConsumptionBudgetResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMConsumptionBudgetResourceGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMConsumptionBudgetResourceGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMConsumptionBudgetResourceGroup_basicUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_consumption_budget_resource_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMConsumptionBudgetResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMConsumptionBudgetResourceGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMConsumptionBudgetResourceGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMConsumptionBudgetResourceGroup_basicUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMConsumptionBudgetResourceGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMConsumptionBudgetResourceGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_consumption_budget_resource_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMConsumptionBudgetResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMConsumptionBudgetResourceGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMConsumptionBudgetResourceGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}
func TestAccAzureRMConsumptionBudgetResourceGroup_completeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_consumption_budget_resource_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMConsumptionBudgetResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMConsumptionBudgetResourceGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMConsumptionBudgetResourceGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMConsumptionBudgetResourceGroup_completeUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMConsumptionBudgetResourceGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMConsumptionBudgetResourceGroup_usageCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_consumption_budget_resource_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMConsumptionBudgetResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMConsumptionBudgetResourceGroup_usageCategory(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMConsumptionBudgetResourceGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMConsumptionBudgetResourceGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Consumption.BudgetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		consumptionBudgetName := rs.Primary.Attributes["name"]
		subscriptionId, hasSubscriptionId := rs.Primary.Attributes["subscription_id"]
		resourceGroupName, hasResourceGroupName := rs.Primary.Attributes["resource_group_name"]

		if !hasSubscriptionId {
			return fmt.Errorf("bad: no subscription id found in state for Consumption Budget: %s", consumptionBudgetName)
		}

		if !hasResourceGroupName {
			return fmt.Errorf("bad: no resource group name found in state for Consumption Budget: %s", consumptionBudgetName)
		}

		scope := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", subscriptionId, resourceGroupName)

		resp, err := conn.Get(ctx, scope, consumptionBudgetName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Consumption Budget %q for scope %q does not exist", consumptionBudgetName, scope)
			}

			return fmt.Errorf("bad: Get on consumptionBudgetClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMConsumptionBudgetResourceGroupDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Consumption.BudgetsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_consumption_budget_resource_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		subscriptionId := rs.Primary.Attributes["subscription_id"]
		resourceGroupName := rs.Primary.Attributes["resource_group_name"]
		scope := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", subscriptionId, resourceGroupName)

		resp, err := conn.Get(ctx, scope, name)

		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}
	}

	return nil
}

func testAccAzureRMConsumptionBudgetResourceGroup_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_consumption_budget_resource_group" "test" {
  name                = "acctestconsumptionbudgetresourcegroup-%d"
  subscription_id     = data.azurerm_subscription.current.subscription_id
  resource_group_name = azurerm_resource_group.test.name

  amount     = 1000
  category   = "Cost"
  time_grain = "Monthly"

  time_period {
    start_date = "%s"
  }

  notification {
    enabled   = true
    threshold = 90.0
    operator  = "EqualTo"

    contact_emails = [
      "foo@example.com",
      "bar@example.com",
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, consumptionBudgetTestStartDate().Format(time.RFC3339))
}

func testAccAzureRMConsumptionBudgetResourceGroup_basicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_consumption_budget_resource_group" "test" {
  name                = "acctestconsumptionbudgetresourcegroup-%d"
  subscription_id     = data.azurerm_subscription.current.subscription_id
  resource_group_name = azurerm_resource_group.test.name

  // Changed the amount from 1000 to 2000
  amount     = 3000
  category   = "Cost"
  time_grain = "Monthly"

  // Add end_date
  time_period {
    start_date = "%s"
    end_date   = "%s"
  }

  // Changed threshold and operator
  notification {
    enabled   = true
    threshold = 95.0
    operator  = "GreaterThan"

    contact_emails = [
      "foo@example.com",
      "bar@example.com",
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, consumptionBudgetTestStartDate().Format(time.RFC3339), consumptionBudgetTestStartDate().AddDate(1, 1, 0).Format(time.RFC3339))
}

func testAccAzureRMConsumptionBudgetResourceGroup_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestAG-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestAG"
}

resource "azurerm_consumption_budget_resource_group" "test" {
  name                = "acctestconsumptionbudgetresourcegroup-%d"
  subscription_id     = data.azurerm_subscription.current.subscription_id
  resource_group_name = azurerm_resource_group.test.name

  amount     = 1000
  category   = "Cost"
  time_grain = "Monthly"

  time_period {
    start_date = "%s"
    end_date   = "%s"
  }

  filter {
    resource_groups = [
      azurerm_resource_group.test.name,
    ]
    resources = [
      azurerm_monitor_action_group.test.id,
    ]
    meters = [
      "00000000-0000-0000-0000-000000000000",
    ]
    tag {
      name = "foo"
      values = [
        "bar",
        "baz",
      ]
    }
  }

  notification {
    enabled   = true
    threshold = 90.0
    operator  = "EqualTo"

    contact_emails = [
      "foo@example.com",
      "bar@example.com",
    ]

    contact_groups = [
      azurerm_monitor_action_group.test.id,
    ]

    contact_roles = [
      "Owner",
    ]
  }

  notification {
    enabled   = false
    threshold = 100.0
    operator  = "GreaterThan"

    contact_emails = [
      "foo@example.com",
      "bar@example.com",
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, consumptionBudgetTestStartDate().Format(time.RFC3339), consumptionBudgetTestStartDate().AddDate(1, 1, 0).Format(time.RFC3339))
}

func testAccAzureRMConsumptionBudgetResourceGroup_completeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestAG-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestAG"
}

resource "azurerm_consumption_budget_resource_group" "test" {
  name                = "acctestconsumptionbudgetresourcegroup-%d"
  subscription_id     = data.azurerm_subscription.current.subscription_id
  resource_group_name = azurerm_resource_group.test.name

  // Changed the amount from 1000 to 2000
  amount     = 2000
  category   = "Cost"
  time_grain = "Monthly"

  // Removed end_date
  time_period {
    start_date = "%s"
  }

  filter {
    resource_groups = [
      azurerm_resource_group.test.name,
    ]
    // Removed resources
    meters = [
      "00000000-0000-0000-0000-000000000000",
    ]
    tag {
      name = "foo"
      values = [
        "bar",
        "baz",
      ]
    }
    // Added tag: zip
    tag {
      name = "zip"
      values = [
        "zap",
        "zop",
      ]
    }
  }

  notification {
    enabled   = true
    threshold = 90.0
    operator  = "EqualTo"

    contact_emails = [
      // Added baz@example.com
      "baz@example.com",
      "foo@example.com",
      "bar@example.com",
    ]

    contact_groups = [
      azurerm_monitor_action_group.test.id,
    ]
    // Removed contact_roles
  }

  notification {
    // Set enabled to true
    enabled   = true
    threshold = 100.0
    // Changed from EqualTo to GreaterThanOrEqualTo 
    operator = "GreaterThanOrEqualTo"

    contact_emails = [
      "foo@example.com",
      "bar@example.com",
    ]

    // Added contact_groups
    contact_groups = [
      azurerm_monitor_action_group.test.id,
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, consumptionBudgetTestStartDate().Format(time.RFC3339))
}

func testAccAzureRMConsumptionBudgetResourceGroup_usageCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_consumption_budget_resource_group" "test" {
  name                = "acctestconsumptionbudgetresourcegroup-%d"
  subscription_id     = data.azurerm_subscription.current.subscription_id
  resource_group_name = azurerm_resource_group.test.name

  amount     = 1000
  category   = "Usage"
  time_grain = "Monthly"

  time_period {
    start_date = "%s"
  }

  filter {
    meters = [
      "00000000-0000-0000-0000-000000000000",
    ]
  }

  notification {
    enabled   = true
    threshold = 90.0
    operator  = "EqualTo"

    contact_emails = [
      "foo@example.com",
      "bar@example.com",
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, consumptionBudgetTestStartDate().Format(time.RFC3339))
}
