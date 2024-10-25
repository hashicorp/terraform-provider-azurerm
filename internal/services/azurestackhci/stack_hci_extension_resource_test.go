package azurestackhci_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StackHCIExtensionResource struct{}

const arcSettingIdEnv = "ARM_TEST_STACK_HCI_ARC_SETTING_ID"

func TestAccStackHCIExtension_basic(t *testing.T) {
	arcSettingId := os.Getenv(arcSettingIdEnv)
	if arcSettingId == "" {
		t.Skipf("skipping since %q has not been set", arcSettingIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_extension", "test")
	r := StackHCIExtensionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("protected_settings"),
	})
}

func TestAccStackHCIExtension_complete(t *testing.T) {
	arcSettingId := os.Getenv(arcSettingIdEnv)
	if arcSettingId == "" {
		t.Skipf("skipping since %q has not been set", arcSettingIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_extension", "test")
	r := StackHCIExtensionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("protected_settings"),
	})
}

func TestAccStackHCIExtension_update(t *testing.T) {
	arcSettingId := os.Getenv(arcSettingIdEnv)
	if arcSettingId == "" {
		t.Skipf("skipping since %q has not been set", arcSettingIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_extension", "test")
	r := StackHCIExtensionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("protected_settings"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("protected_settings"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("protected_settings"),
	})
}

func TestAccStackHCIExtension_requiresImport(t *testing.T) {
	arcSettingId := os.Getenv(arcSettingIdEnv)
	if arcSettingId == "" {
		t.Skipf("skipping since %q has not been set", arcSettingIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_extension", "test")
	r := StackHCIExtensionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r StackHCIExtensionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	clusterClient := client.AzureStackHCI.Extensions
	id, err := extensions.ParseExtensionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clusterClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r StackHCIExtensionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_extension" "test" {
  name           = "acctest-shce-%[2]s"
  arc_setting_id = %[3]q
  publisher      = "Microsoft.EnterpriseCloud.Monitoring"
  type           = "MicrosoftMonitoringAgent"
}
`, r.template(data), data.RandomString, os.Getenv(arcSettingIdEnv))
}

func (r StackHCIExtensionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_extension" "test" {
  name                               = "acctest-shce-%[2]s"
  arc_setting_id                     = %[3]q
  publisher                          = "Microsoft.EnterpriseCloud.Monitoring"
  type                               = "MicrosoftMonitoringAgent"
  auto_upgrade_minor_version_enabled = true
  automatic_upgrade_enabled          = false
  type_handler_version               = "1.22.0"

  protected_settings = <<PROTECTED_SETTINGS
{
	"workspaceKey": "${azurerm_log_analytics_workspace.test.primary_shared_key}"
}
PROTECTED_SETTINGS

  settings = <<SETTINGS
{
	"workspaceId": "${azurerm_log_analytics_workspace.test.workspace_id}"
}
SETTINGS
}
`, r.template(data), data.RandomString, os.Getenv(arcSettingIdEnv))
}

func (r StackHCIExtensionResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)

	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_extension" "import" {
  name               = azurerm_stack_hci_extension.test.name
  arc_setting_id     = azurerm_stack_hci_extension.test.arc_setting_id
  publisher          = azurerm_stack_hci_extension.test.publisher
  type               = azurerm_stack_hci_extension.test.type
  protected_settings = azurerm_stack_hci_extension.test.protected_settings
  settings           = azurerm_stack_hci_extension.test.settings
}
`, config)
}

func (r StackHCIExtensionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-hciext-%s"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctesthci-law-%[1]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}
`, data.RandomString, data.Locations.Primary)
}
