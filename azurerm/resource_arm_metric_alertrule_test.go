package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestValidateMetricAlertRuleTags(t *testing.T) {
	cases := []struct {
		Name     string
		Value    map[string]interface{}
		ErrCount int
	}{
		{
			Name: "Single Valid",
			Value: map[string]interface{}{
				"hello": "world",
			},
			ErrCount: 0,
		},
		{
			Name: "Single Invalid",
			Value: map[string]interface{}{
				"$Type": "hello/world",
			},
			ErrCount: 1,
		},
		{
			Name: "Single Invalid lowercase",
			Value: map[string]interface{}{
				"$type": "hello/world",
			},
			ErrCount: 1,
		},
		{
			Name: "Multiple Valid",
			Value: map[string]interface{}{
				"hello": "world",
				"foo":   "bar",
			},
			ErrCount: 0,
		},
		{
			Name: "Multiple Invalid",
			Value: map[string]interface{}{
				"hello": "world",
				"$type": "Microsoft.Foo/Bar",
			},
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateMetricAlertRuleTags(tc.Value, "azurerm_metric_alert_rule")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected %q to return %d errors but returned %d", tc.Name, tc.ErrCount, len(errors))
		}
	}
}

func TestAccAzureRMMetricAlertRule_virtualMachineCpu(t *testing.T) {
	resourceName := "azurerm_metric_alertrule.test"
	ri := tf.AccRandTimeInt()
	preConfig := testAccAzureRMMetricAlertRule_virtualMachineCpu(ri, acceptance.Location(), true)
	postConfig := testAccAzureRMMetricAlertRule_virtualMachineCpu(ri, acceptance.Location(), false)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMetricAlertRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMetricAlertRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckNoResourceAttr(resourceName, "tags.$type"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMetricAlertRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceName, "tags.$type"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMetricAlertRuleExists(resourceName),
					resource.TestCheckNoResourceAttr(resourceName, "tags.$type"),
				),
			},
		},
	})
}

func TestAccAzureRMMetricAlertRule_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_metric_alertrule.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMetricAlertRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMetricAlertRule_virtualMachineCpu(ri, location, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMetricAlertRuleExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMMetricAlertRule_requiresImport(ri, location, true),
				ExpectError: acceptance.RequiresImportError("azurerm_metric_alertrule"),
			},
		},
	})
}

func TestAccAzureRMMetricAlertRule_sqlDatabaseStorage(t *testing.T) {
	resourceName := "azurerm_metric_alertrule.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMetricAlertRule_sqlDatabaseStorage(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMetricAlertRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMetricAlertRuleExists(resourceName),
					resource.TestCheckNoResourceAttr(resourceName, "tags.$type"),
				),
			},
		},
	})
}

func testCheckAzureRMMetricAlertRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Alert Rule: %s", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.AlertRulesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Alert Rule %q (resource group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on monitorAlertRulesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMMetricAlertRuleDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.AlertRulesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_metric_alertrule" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Alert Rule still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMMetricAlertRule_virtualMachineCpu(rInt int, location string, enabled bool) string {
	template := testAccAzureRMVirtualMachine_basicLinuxMachine_managedDisk_explicit(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_metric_alertrule" "test" {
  name                = "${azurerm_virtual_machine.test.name}-cpu"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  description = "An alert rule to watch the metric Percentage CPU"

  enabled = %t

  resource_id = "${azurerm_virtual_machine.test.id}"
  metric_name = "Percentage CPU"
  operator    = "GreaterThan"
  threshold   = 75
  aggregation = "Average"
  period      = "PT5M"

  email_action {
    send_to_service_owners = false

    custom_emails = [
      "support@azure.microsoft.com",
    ]
  }

  webhook_action {
    service_uri = "https://requestb.in/18jamc41"

    properties = {
      severity        = "incredible"
      acceptance_test = "true"
    }
  }
}
`, template, enabled)
}

func testAccAzureRMMetricAlertRule_requiresImport(rInt int, location string, enabled bool) string {
	template := testAccAzureRMMetricAlertRule_virtualMachineCpu(rInt, location, enabled)
	return fmt.Sprintf(`
%s

resource "azurerm_metric_alertrule" "import" {
  name                = "${azurerm_metric_alertrule.test.name}"
  resource_group_name = "${azurerm_metric_alertrule.test.resource_group_name}"
  location            = "${azurerm_metric_alertrule.test.location}"
  description         = "${azurerm_metric_alertrule.test.description}"
  enabled             = "${azurerm_metric_alertrule.test.enabled}"

  resource_id = "${azurerm_virtual_machine.test.id}"
  metric_name = "Percentage CPU"
  operator    = "GreaterThan"
  threshold   = 75
  aggregation = "Average"
  period      = "PT5M"

  email_action {
    send_to_service_owners = false

    custom_emails = [
      "support@azure.microsoft.com",
    ]
  }

  webhook_action {
    service_uri = "https://requestb.in/18jamc41"

    properties = {
      severity        = "incredible"
      acceptance_test = "true"
    }
  }
}
`, template)
}

func testAccAzureRMMetricAlertRule_sqlDatabaseStorage(rInt int, location string) string {
	basicSqlServerDatabase := testAccAzureRMSqlDatabase_basic(rInt, location)

	return fmt.Sprintf(`
%s

resource "azurerm_metric_alertrule" "test" {
  name                = "${azurerm_sql_database.test.name}-storage"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  description = "An alert rule to watch the metric Storage"

  enabled = true

  resource_id = "${azurerm_sql_database.test.id}"
  metric_name = "storage"
  operator    = "GreaterThan"
  threshold   = 1073741824
  aggregation = "Maximum"
  period      = "PT10M"

  email_action {
    send_to_service_owners = false

    custom_emails = [
      "support@azure.microsoft.com",
    ]
  }

  webhook_action {
    service_uri = "https://requestb.in/18jamc41"

    properties = {
      severity        = "incredible"
      acceptance_test = "true"
    }
  }
}
`, basicSqlServerDatabase)
}
