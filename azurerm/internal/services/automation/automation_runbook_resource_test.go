package automation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AutomationRunbookResource struct {
}

func TestAccAutomationRunbook_PSWorkflow(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_runbook", "test")
	r := AutomationRunbookResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.PSWorkflow(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("publish_content_link"),
	})
}

func TestAccAutomationRunbook_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_runbook", "test")
	r := AutomationRunbookResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.PSWorkflow(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAutomationRunbook_PSWorkflowWithHash(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_runbook", "test")
	r := AutomationRunbookResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.PSWorkflowWithHash(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("runbook_type").HasValue("PowerShellWorkflow"),
			),
		},
		data.ImportStep("publish_content_link"),
	})
}

func TestAccAutomationRunbook_PSWithContent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_runbook", "test")
	r := AutomationRunbookResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.PSWithContent(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("runbook_type").HasValue("PowerShell"),
				check.That(data.ResourceName).Key("content").HasValue("# Some test content\n# for Terraform acceptance test\n"),
			),
		},
		data.ImportStep("publish_content_link"),
	})
}

func TestAccAutomationRunbook_PSWorkflowWithoutUri(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_runbook", "test")
	r := AutomationRunbookResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.PSWorkflowWithoutUri(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("publish_content_link"),
	})
}

func TestAccAutomationRunbook_withJobSchedule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_runbook", "test")
	r := AutomationRunbookResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.PSWorkflow(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("publish_content_link"),
		{
			Config: r.withJobSchedule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("publish_content_link"),
		{
			Config: r.withJobScheduleUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("publish_content_link"),
		{
			Config: r.withoutJobSchedule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("publish_content_link"),
	})
}

func (t AutomationRunbookResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resGroup := id.ResourceGroup
	accName := id.Path["automationAccounts"]
	name := id.Path["runbooks"]

	resp, err := clients.Automation.RunbookClient.Get(ctx, resGroup, accName, name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Automation Runbook '%s' (resource group: '%s') does not exist", name, id.ResourceGroup)
	}

	return utils.Bool(resp.RunbookProperties != nil), nil
}

func (AutomationRunbookResource) PSWorkflow(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
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

func (AutomationRunbookResource) requiresImport(data acceptance.TestData) string {
	template := AutomationRunbookResource{}.PSWorkflow(data)
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

func (AutomationRunbookResource) PSWorkflowWithHash(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
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

func (AutomationRunbookResource) PSWithContent(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
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

func (AutomationRunbookResource) PSWorkflowWithoutUri(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
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

func (AutomationRunbookResource) withJobSchedule(data acceptance.TestData) string {
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

func (AutomationRunbookResource) withJobScheduleUpdated(data acceptance.TestData) string {
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

func (AutomationRunbookResource) withoutJobSchedule(data acceptance.TestData) string {
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
