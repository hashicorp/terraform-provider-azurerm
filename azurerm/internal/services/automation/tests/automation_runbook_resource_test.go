package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAutomationRunbook_PSWorkflow(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_runbook", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationRunbookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationRunbook_PSWorkflow(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationRunbookExists(data.ResourceName),
				),
			},
			data.ImportStep("publish_content_link"),
		},
	})
}

func TestAccAzureRMAutomationRunbook_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_runbook", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationRunbookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationRunbook_PSWorkflow(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationRunbookExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMAutomationRunbook_requiresImport),
		},
	})
}

func TestAccAzureRMAutomationRunbook_PSWorkflowWithHash(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_runbook", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationRunbookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationRunbook_PSWorkflowWithHash(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationRunbookExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "runbook_type", "PowerShellWorkflow"),
				),
			},
			data.ImportStep("publish_content_link"),
		},
	})
}

func TestAccAzureRMAutomationRunbook_PSWithContent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_runbook", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationRunbookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationRunbook_PSWithContent(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationRunbookExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "runbook_type", "PowerShell"),
					resource.TestCheckResourceAttr(data.ResourceName, "content", "# Some test content\n# for Terraform acceptance test\n"),
				),
			},
			data.ImportStep("publish_content_link"),
		},
	})
}

func TestAccAzureRMAutomationRunbook_PSWorkflowWithoutUri(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_runbook", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationRunbookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationRunbook_PSWorkflowWithoutUri(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationRunbookExists(data.ResourceName),
				),
			},
			data.ImportStep("publish_content_link"),
		},
	})
}

func TestAccAzureRMAutomationRunbook_withJobSchedule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_runbook", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationRunbookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationRunbook_PSWorkflow(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationRunbookExists(data.ResourceName),
				),
			},
			data.ImportStep("publish_content_link"),
			{
				Config: testAccAzureRMAutomationRunbook_withJobSchedule(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationRunbookExists(data.ResourceName),
				),
			},
			data.ImportStep("publish_content_link"),
			{
				Config: testAccAzureRMAutomationRunbook_withJobScheduleUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationRunbookExists(data.ResourceName),
				),
			},
			data.ImportStep("publish_content_link"),
			{
				Config: testAccAzureRMAutomationRunbook_withoutJobSchedule(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationRunbookExists(data.ResourceName),
				),
			},
			data.ImportStep("publish_content_link"),
		},
	})
}

func testCheckAzureRMAutomationRunbookDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Automation.RunbookClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_runbook" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["automation_account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Runbook: '%s'", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, accName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Automation Runbook still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMAutomationRunbookExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Automation.RunbookClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["automation_account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Runbook: '%s'", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, accName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Automation Runbook '%s' (resource group: '%s') does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on automationRunbookClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMAutomationRunbook_PSWorkflow(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_runbook" "test" {
  name                    = "Get-AzureVMTutorial"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name

  log_verbose  = "true"
  log_progress = "true"
  description  = "This is a test runbook for terraform acceptance test"
  runbook_type = "PowerShell"

  content = <<CONTENT
# Some test content
# for Terraform acceptance test
CONTENT
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAutomationRunbook_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMAutomationRunbook_PSWorkflow(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_runbook" "import" {
  name                    = azurerm_automation_runbook.test.name
  location                = azurerm_automation_runbook.test.location
  resource_group_name     = azurerm_automation_runbook.test.resource_group_name
  automation_account_name = azurerm_automation_runbook.test.automation_account_name

  log_verbose  = "true"
  log_progress = "true"
  description  = "This is a test runbook for terraform acceptance test"
  runbook_type = "PowerShell"

  content = <<CONTENT
# Some test content
# for Terraform acceptance test
CONTENT
}
`, template)
}

func testAccAzureRMAutomationRunbook_PSWorkflowWithHash(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_runbook" "test" {
  name                    = "Get-AzureVMTutorial"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name

  log_verbose  = "true"
  log_progress = "true"
  description  = "This is a test runbook for terraform acceptance test"
  runbook_type = "PowerShellWorkflow"

  publish_content_link {
    uri     = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/c4935ffb69246a6058eb24f54640f53f69d3ac9f/101-automation-runbook-getvms/Runbooks/Get-AzureVMTutorial.ps1"
    version = "1.0.0.0"

    hash {
      algorithm = "SHA256"
      value     = "115775B8FF2BE672D8A946BD0B489918C724DDE15A440373CA54461D53010A80"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAutomationRunbook_PSWithContent(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_runbook" "test" {
  name                    = "Get-AzureVMTutorial"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name

  log_verbose  = "true"
  log_progress = "true"
  description  = "This is a test runbook for terraform acceptance test"
  runbook_type = "PowerShell"

  publish_content_link {
    uri = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/c4935ffb69246a6058eb24f54640f53f69d3ac9f/101-automation-runbook-getvms/Runbooks/Get-AzureVMTutorial.ps1"
  }

  content = <<CONTENT
# Some test content
# for Terraform acceptance test
CONTENT

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAutomationRunbook_PSWorkflowWithoutUri(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_runbook" "test" {
  name                    = "Get-AzureVMTutorial"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name

  log_verbose  = "true"
  log_progress = "true"
  description  = "This is a test runbook for terraform acceptance test"
  runbook_type = "PowerShell"

  publish_content_link {
    uri = ""
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAutomationRunbook_withJobSchedule(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%[1]d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  frequency               = "OneTime"
}

resource "azurerm_automation_runbook" "test" {
  name                    = "Get-AzureVMTutorial"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name

  log_verbose  = "true"
  log_progress = "true"
  description  = "This is a test runbook for terraform acceptance test"
  runbook_type = "PowerShell"

  content = <<CONTENT
# Some test content
# for Terraform acceptance test
CONTENT

  job_schedule {
    schedule_name = azurerm_automation_schedule.test.name
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMAutomationRunbook_withJobScheduleUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%[1]d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  frequency               = "OneTime"
}

resource "azurerm_automation_runbook" "test" {
  name                    = "Get-AzureVMTutorial"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name

  log_verbose  = "true"
  log_progress = "true"
  description  = "This is a test runbook for terraform acceptance test"
  runbook_type = "PowerShell"

  content = <<CONTENT
param(
    [string]$Output = "World",
  )
  "Hello, " + $Output + "!"
CONTENT

  job_schedule {
    schedule_name = azurerm_automation_schedule.test.name
    parameters = {
      output     = "Earth"
      case       = "MATTERS"
      keepcount  = 20
      webhookuri = "http://www.example.com/hook"
      url        = "https://www.Example.com"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMAutomationRunbook_withoutJobSchedule(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%[1]d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  frequency               = "OneTime"
}

resource "azurerm_automation_runbook" "test" {
  name                    = "Get-AzureVMTutorial"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name

  log_verbose  = "true"
  log_progress = "true"
  description  = "This is a test runbook for terraform acceptance test"
  runbook_type = "PowerShell"

  content = <<CONTENT
param(
    [string]$Output = "World",
  )
  "Hello, " + $Output + "!"
CONTENT

  job_schedule = []
}
`, data.RandomInteger, data.Locations.Primary)
}
