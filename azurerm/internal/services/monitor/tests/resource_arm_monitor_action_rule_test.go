package tests

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

func TestAccAzureRMMonitorActionRule_diagnostics_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAlertsManagementActionRule_diagnostics_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionRule_diagnostics_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAlertsManagementActionRule_diagnostics_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionRuleExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMAlertsManagementActionRule_diagnostics_requiresImport),
		},
	})
}

func TestAccAzureRMMonitorActionRule_diagnostics_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAlertsManagementActionRule_diagnostics_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionRule_diagnostics_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAlertsManagementActionRule_diagnostics_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAlertsManagementActionRule_diagnostics_complete(data),
				Check:  resource.ComposeTestCheckFunc(),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAlertsManagementActionRule_diagnostics_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionRule_actionGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAlertsManagementActionRule_actionGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionRule_actionGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAlertsManagementActionRule_actionGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionRule_actionGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAlertsManagementActionRule_actionGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAlertsManagementActionRule_actionGroup_complete(data),
				Check:  resource.ComposeTestCheckFunc(),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAlertsManagementActionRule_actionGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionRule_suppression_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAlertsManagementActionRule_suppression_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionRule_suppression_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAlertsManagementActionRule_suppression_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActionRule_suppression_updateSuppressionConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAlertsManagementActionRule_suppression_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAlertsManagementActionRule_suppression_dailyRecurrence(data),
				Check:  resource.ComposeTestCheckFunc(),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAlertsManagementActionRule_suppression_monthlyRecurrence(data),
				Check:  resource.ComposeTestCheckFunc(),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAlertsManagementActionRule_suppression_complete(data),
				Check:  resource.ComposeTestCheckFunc(),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAlertsManagementActionRule_suppression_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActionRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMMonitorActionRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.ActionRulesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can not found AlertsManagement ActionRule: %s", resourceName)
		}
		id, err := parse.ActionRuleID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.GetByName(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: alerts_management action_rule %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on AlertsManagement ActionRulesClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMMonitorActionRuleDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.ActionRulesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_monitor_action_rule" {
			continue
		}
		id, err := parse.ActionRuleID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.GetByName(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on AlertsManagement ActionRulesClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMAlertsManagementActionRule_diagnostics_basic(data acceptance.TestData) string {
	template := testAccAzureRMAlertsManagementActionRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_rule" "test" {
  name                = "acctest-ar-%d"
  resource_group_name = azurerm_resource_group.test.name
  type                = "Diagnostics"
}
`, template, data.RandomInteger)
}

func testAccAzureRMAlertsManagementActionRule_diagnostics_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMAlertsManagementActionRule_diagnostics_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_rule" "import" {
  name                = azurerm_monitor_action_rule.test.name
  resource_group_name = azurerm_monitor_action_rule.test.resource_group_name
  type                = "Diagnostics"

  scope {
    type         = "ResourceGroup"
    resource_ids = [azurerm_resource_group.test.id]
  }
}
`, template)
}

func testAccAzureRMAlertsManagementActionRule_diagnostics_complete(data acceptance.TestData) string {
	template := testAccAzureRMAlertsManagementActionRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_rule" "test" {
  name                = "acctest-ar-%d"
  resource_group_name = azurerm_resource_group.test.name
  type                = "Diagnostics"
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

func testAccAzureRMAlertsManagementActionRule_actionGroup_basic(data acceptance.TestData) string {
	template := testAccAzureRMAlertsManagementActionRule_actionGroup_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_rule" "test" {
  name                = "acctest-ar-%d"
  resource_group_name = azurerm_resource_group.test.name
  type                = "ActionGroup"
  action_group_id     = azurerm_monitor_action_group.test.id
}
`, template, data.RandomInteger)
}

func testAccAzureRMAlertsManagementActionRule_actionGroup_complete(data acceptance.TestData) string {
	template := testAccAzureRMAlertsManagementActionRule_actionGroup_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_rule" "test" {
  name                = "acctest-ar-%d"
  resource_group_name = azurerm_resource_group.test.name
  type                = "ActionGroup"
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

func testAccAzureRMAlertsManagementActionRule_suppression_basic(data acceptance.TestData) string {
	template := testAccAzureRMAlertsManagementActionRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_rule" "test" {
  name                = "acctest-ar-%d"
  resource_group_name = azurerm_resource_group.test.name
  type                = "Suppression"

  suppression {
    recurrence_type = "Always"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMAlertsManagementActionRule_suppression_complete(data acceptance.TestData) string {
	template := testAccAzureRMAlertsManagementActionRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctest%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_action_rule" "test" {
  name                = "acctest-ar-%d"
  resource_group_name = azurerm_resource_group.test.name
  type                = "Suppression"
  enabled             = false
  description         = "actionRule-test"

  scope {
    type         = "ResourceGroup"
    resource_ids = [azurerm_resource_group.test.id]
  }

  suppression {
    recurrence_type = "Weekly"

    schedule {
      start_date = "12/09/2018"
      start_time = "06:00:00"
      end_date   = "12/18/2018"
      end_time   = "14:00:00"

      recurrence_weekly = ["Sunday", "Monday", "Friday", "Saturday"]
    }
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
`, template, data.RandomString, data.RandomInteger)
}

func testAccAzureRMAlertsManagementActionRule_suppression_dailyRecurrence(data acceptance.TestData) string {
	template := testAccAzureRMAlertsManagementActionRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_rule" "test" {
  name                = "acctest-ar-%d"
  resource_group_name = azurerm_resource_group.test.name
  type                = "Suppression"

  scope {
    type         = "ResourceGroup"
    resource_ids = [azurerm_resource_group.test.id]
  }

  suppression {
    recurrence_type = "Daily"

    schedule {
      start_date = "12/09/2018"
      start_time = "06:00:00"
      end_date   = "12/18/2018"
      end_time   = "14:00:00"
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMAlertsManagementActionRule_suppression_monthlyRecurrence(data acceptance.TestData) string {
	template := testAccAzureRMAlertsManagementActionRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_rule" "test" {
  name                = "acctest-ar-%d"
  resource_group_name = azurerm_resource_group.test.name
  type                = "Suppression"

  scope {
    type         = "ResourceGroup"
    resource_ids = [azurerm_resource_group.test.id]
  }

  suppression {
    recurrence_type = "Monthly"

    schedule {
      start_date         = "12/09/2018"
      start_time         = "06:00:00"
      end_date           = "12/18/2018"
      end_time           = "14:00:00"
      recurrence_monthly = [1, 2, 15, 30, 31]
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMAlertsManagementActionRule_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMAlertsManagementActionRule_actionGroup_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
