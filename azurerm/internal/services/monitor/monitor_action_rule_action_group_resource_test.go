package monitor_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/monitor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMonitorActionRuleActionGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionRuleActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionRuleActionGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionRuleActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionRuleActionGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule_action_group", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionRuleActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionRuleActionGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionRuleActionGroupExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMMonitorActionRuleActionGroup_requiresImport),
		},
	})
}

func TestAccAzureRMMonitorActionRuleActionGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionRuleActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionRuleActionGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionRuleActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionRuleActionGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule_action_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionRuleActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActionRuleActionGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionRuleActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorActionRuleActionGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionRuleActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorActionRuleActionGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionRuleActionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMMonitorActionRuleActionGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.ActionRulesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can not found Monitor ActionRule: %s", resourceName)
		}
		id, err := parse.ActionRuleID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.GetByName(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: monitor action_rule %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on Monitor ActionRulesClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMMonitorActionRuleActionGroupDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.ActionRulesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_monitor_action_rule_action_group" {
			continue
		}
		id, err := parse.ActionRuleID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.GetByName(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Monitor ActionRulesClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMMonitorActionRuleActionGroup_basic(data acceptance.TestData) string {
	template := testAccAzureRMMonitorActionRuleActionGroup_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_rule_action_group" "test" {
  name                = "acctest-moniter-%d"
  resource_group_name = azurerm_resource_group.test.name
  action_group_id     = azurerm_monitor_action_group.test.id
}
`, template, data.RandomInteger)
}

func testAccAzureRMMonitorActionRuleActionGroup_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMonitorActionRuleActionGroup_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_rule_action_group" "import" {
  name                = azurerm_monitor_action_rule_action_group.test.name
  resource_group_name = azurerm_monitor_action_rule_action_group.test.resource_group_name
  action_group_id     = azurerm_monitor_action_rule_action_group.test.action_group_id
}
`, template)
}

func testAccAzureRMMonitorActionRuleActionGroup_complete(data acceptance.TestData) string {
	template := testAccAzureRMMonitorActionRuleActionGroup_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_rule_action_group" "test" {
  name                = "acctest-moniter-%d"
  resource_group_name = azurerm_resource_group.test.name
  action_group_id     = azurerm_monitor_action_group.test.id
  enabled             = false
  description         = "actionRule-test"

  scope {
    type         = "ResourceGroup"
    resource_ids = [azurerm_resource_group.test.id]
  }

  condition {
    alert_context {
      operator = "Contains"
      values   = ["context1", "context2"]
    }

    alert_rule_id {
      operator = "Contains"
      values   = ["ruleId1", "ruleId2"]
    }

    description {
      operator = "DoesNotContain"
      values   = ["description1", "description2"]
    }

    monitor {
      operator = "NotEquals"
      values   = ["Fired"]
    }

    monitor_service {
      operator = "Equals"
      values   = ["Data Box Edge", "Data Box Gateway", "Resource Health"]
    }

    severity {
      operator = "Equals"
      values   = ["Sev0", "Sev1", "Sev2"]
    }

    target_resource_type {
      operator = "Equals"
      values   = ["Microsoft.Compute/VirtualMachines", "microsoft.batch/batchaccounts"]
    }
  }

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMMonitorActionRuleActionGroup_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-monitor-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
