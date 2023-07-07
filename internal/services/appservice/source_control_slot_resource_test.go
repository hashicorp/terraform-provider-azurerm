// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// (@jackofallops) Note: Some tests require a valid GitHub token for ARM_GITHUB_ACCESS_TOKEN. This token needs the `repo` and `workflow` permissions to be applicable on the referenced repositories.

type SourceControlSlotResource struct{}

func TestAccSourceControlSlotResource_windowsExternalGit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_source_control_slot", "test")
	r := SourceControlSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsExternalGit(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSourceControlSlotResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_source_control_slot", "test")
	r := SourceControlSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsExternalGit(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSourceControlSlotResource_windowsLocalGit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_source_control_slot", "test")
	r := SourceControlSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsLocalGit(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scm_type").HasValue("LocalGit"),
				check.That(data.ResourceName).Key("repo_url").HasValue(fmt.Sprintf("https://acctestwa-%[1]d-acctestwas-%[1]d.scm.azurewebsites.net", data.RandomInteger)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSourceControlSlotResource_windowsGitHubAction(t *testing.T) {
	if ok := os.Getenv("ARM_GITHUB_ACCESS_TOKEN"); ok == "" {
		t.Skip("Skipping as `ARM_GITHUB_ACCESS_TOKEN` is not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_app_service_source_control_slot", "test")
	r := SourceControlSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsGitHubAction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSourceControlSlotResource_windowsGitHub(t *testing.T) {
	if ok := os.Getenv("ARM_GITHUB_ACCESS_TOKEN"); ok == "" {
		t.Skip("Skipping as `ARM_GITHUB_ACCESS_TOKEN` is not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_app_service_source_control_slot", "test")
	r := SourceControlSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsGitHub(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSourceControlSlotResource_linuxExternalGit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_source_control_slot", "test")
	r := SourceControlSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxExternalGit(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSourceControlSlotResource_linuxLocalGit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_source_control_slot", "test")
	r := SourceControlSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxLocalGit(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scm_type").HasValue("LocalGit"),
				check.That(data.ResourceName).Key("repo_url").HasValue(fmt.Sprintf("https://acctestwa-%[1]d-acctestwas-%[1]d.scm.azurewebsites.net", data.RandomInteger)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSourceControlSlotResource_linuxGitHubAction(t *testing.T) {
	if ok := os.Getenv("ARM_GITHUB_ACCESS_TOKEN"); ok == "" {
		t.Skip("Skipping as `ARM_GITHUB_ACCESS_TOKEN` is not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_app_service_source_control_slot", "test")
	r := SourceControlSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxGitHubAction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("uses_github_action").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSourceControlSlotResource_linuxGitHub(t *testing.T) {
	if ok := os.Getenv("ARM_GITHUB_ACCESS_TOKEN"); ok == "" {
		t.Skip("Skipping as `ARM_GITHUB_ACCESS_TOKEN` is not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_app_service_source_control_slot", "test")
	r := SourceControlSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxGitHub(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SourceControlSlotResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.WebAppSlotID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.AppService.WebAppsClient.GetSourceControlSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Source Control for %s: %v", id, err)
	}
	if resp.SiteSourceControlProperties == nil || resp.SiteSourceControlProperties.RepoURL == nil {
		return utils.Bool(false), nil
	}

	return utils.Bool(true), nil
}

func (r SourceControlSlotResource) windowsExternalGit(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_app_service_source_control_slot" "test" {
  slot_id                = azurerm_windows_web_app_slot.test.id
  repo_url               = "https://github.com/Azure-Samples/app-service-web-dotnet-get-started.git"
  branch                 = "master"
  use_manual_integration = true
}
`, r.baseWindowsAppTemplate(data))
}

func (r SourceControlSlotResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_source_control_slot" "import" {
  slot_id                = azurerm_app_service_source_control_slot.test.slot_id
  repo_url               = azurerm_app_service_source_control_slot.test.repo_url
  branch                 = azurerm_app_service_source_control_slot.test.branch
  use_manual_integration = azurerm_app_service_source_control_slot.test.use_manual_integration
}
`, r.windowsExternalGit(data))
}

func (r SourceControlSlotResource) linuxExternalGit(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_app_service_source_control_slot" "test" {
  slot_id                = azurerm_linux_web_app_slot.test.id
  repo_url               = "https://github.com/Azure-Samples/python-docs-hello-world.git"
  branch                 = "master"
  use_manual_integration = true
}
`, r.baseLinuxAppTemplate(data))
}

func (r SourceControlSlotResource) windowsLocalGit(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_app_service_source_control_slot" "test" {
  slot_id       = azurerm_windows_web_app_slot.test.id
  use_local_git = true
}
`, r.baseWindowsAppTemplate(data))
}

func (r SourceControlSlotResource) linuxLocalGit(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_app_service_source_control_slot" "test" {
  slot_id       = azurerm_linux_web_app_slot.test.id
  use_local_git = true
}
`, r.baseLinuxAppTemplate(data))
}

func (r SourceControlSlotResource) windowsGitHubAction(data acceptance.TestData) string {
	token := os.Getenv("ARM_GITHUB_ACCESS_TOKEN")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource azurerm_source_control_token test {
  type  = "GitHub"
  token = "%s"
}

resource "azurerm_app_service_source_control_slot" "test" {
  slot_id  = azurerm_windows_web_app.test.id
  repo_url = "https://github.com/Azure-Samples/app-service-web-dotnet-get-started.git"
  branch   = "master"

  github_action_configuration {
    generate_workflow_file = true

    code_configuration {
      runtime_stack   = "dotnetcore"
      runtime_version = "5.0.x"
    }
  }
}
`, r.baseWindowsAppTemplate(data), token)
}

func (r SourceControlSlotResource) windowsGitHub(data acceptance.TestData) string {
	token := os.Getenv("ARM_GITHUB_ACCESS_TOKEN")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource azurerm_source_control_token test {
  type  = "GitHub"
  token = "%s"
}

resource "azurerm_app_service_source_control_slot" "test" {
  slot_id  = azurerm_windows_web_app_slot.test.id
  repo_url = "https://github.com/jackofallops/azure-app-service-static-site-tests.git"
  branch   = "main"

  depends_on = [
    azurerm_source_control_token.test,
  ]
}
`, r.baseWindowsAppTemplate(data), token)
}

func (r SourceControlSlotResource) linuxGitHubAction(data acceptance.TestData) string {
	token := os.Getenv("ARM_GITHUB_ACCESS_TOKEN")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource azurerm_source_control_token test {
  type  = "GitHub"
  token = "%s"
}

resource "azurerm_app_service_source_control_slot" "test" {
  slot_id  = azurerm_linux_web_app_slot.test.id
  repo_url = "https://github.com/jackofallops/app-service-web-dotnet-get-started.git"
  branch   = "master"

  github_action_configuration {
    generate_workflow_file = true

    code_configuration {
      runtime_stack   = "dotnetcore"
      runtime_version = "5.0.x"
    }
  }
}
`, r.baseLinuxAppTemplate(data), token)
}

func (r SourceControlSlotResource) linuxGitHub(data acceptance.TestData) string {
	token := os.Getenv("ARM_GITHUB_ACCESS_TOKEN")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_source_control_token" "test" {
  type  = "GitHub"
  token = "%s"
}

resource "azurerm_app_service_source_control_slot" "test" {
  slot_id                = azurerm_linux_web_app_slot.test.id
  repo_url               = "https://github.com/Azure-Samples/python-docs-hello-world.git"
  branch                 = "master"
  use_manual_integration = true

  depends_on = [
    azurerm_source_control_token.test,
  ]
}
`, r.baseLinuxAppTemplate(data), token)
}

func (r SourceControlSlotResource) baseWindowsAppTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ASSC-%[1]d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASSC-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Windows"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_app_service_plan.test.id

  site_config {}
}

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%[1]d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r SourceControlSlotResource) baseLinuxAppTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ASSC-%[1]d"
  location = "%[2]s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_app_service_plan.test.id

  site_config {
    application_stack {
      python_version = "3.8"
    }
  }
}

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%[1]d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {
    application_stack {
      python_version = "3.8"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
