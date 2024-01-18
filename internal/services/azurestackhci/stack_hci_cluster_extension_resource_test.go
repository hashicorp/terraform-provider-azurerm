package azurestackhci_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-08-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StackHCIExtensionResource struct{}

func TestAccStackHCIClusterExtension_basic(t *testing.T) {
	stackHCIClusterName := os.Getenv("ARM_TEST_HCI_CLUSTER_NAME")
	if stackHCIClusterName == "" {
		t.Skipf("skip the test due to missing environment variable ARM_TEST_HCI_CLUSTER_NAME")
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
		data.ImportStep("protected_setting"),
	})
}

func TestAccStackHCIClusterExtension_complete(t *testing.T) {
	stackHCIClusterName := os.Getenv("ARM_TEST_HCI_CLUSTER_NAME")
	if stackHCIClusterName == "" {
		t.Skipf("skip the test due to missing environment variable ARM_TEST_HCI_CLUSTER_NAME")
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
		data.ImportStep("protected_setting"),
	})
}

func TestAccStackHCIClusterExtension_update(t *testing.T) {
	stackHCIClusterName := os.Getenv("ARM_TEST_HCI_CLUSTER_NAME")
	if stackHCIClusterName == "" {
		t.Skipf("skip the test due to missing environment variable ARM_TEST_HCI_CLUSTER_NAME")
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
		data.ImportStep("protected_setting"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("protected_setting"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("protected_setting"),
	})
}

func TestAccStackHCIClusterExtension_requiresImport(t *testing.T) {
	stackHCIClusterName := os.Getenv("ARM_TEST_HCI_CLUSTER_NAME")
	if stackHCIClusterName == "" {
		t.Skipf("skip the test due to missing environment variable ARM_TEST_HCI_CLUSTER_NAME")
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
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r StackHCIExtensionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-%s"
  location            = data.azurerm_resource_group.test.location
  resource_group_name = data.azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_stack_hci_extension" "test" {
  name           = "acctest-shce-%[2]s"
  arc_setting_id = data.azurerm_stack_hci_cluster_arc_setting.test.id
  publisher      = "Microsoft.EnterpriseCloud.Monitoring"
  type           = "MicrosoftMonitoringAgent"

  protected_setting = <<PROTECTED_SETTING
{
	"workspaceKey": "${azurerm_log_analytics_workspace.test.primary_shared_key}"
}
PROTECTED_SETTING

  setting = <<SETTING
{
	"workspaceId": "${azurerm_log_analytics_workspace.test.workspace_id}"
}
SETTING
}
`, r.template(data), data.RandomString)
}

func (r StackHCIExtensionResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-%s"
  location            = data.azurerm_resource_group.test.location
  resource_group_name = data.azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_stack_hci_extension" "test" {
  name                       = "acctest-shce-%[2]s"
  arc_setting_id             = data.azurerm_stack_hci_cluster_arc_setting.test.id
  publisher                  = "Microsoft.EnterpriseCloud.Monitoring"
  type                       = "MicrosoftMonitoringAgent"
  auto_upgrade_minor_version = true
  automatic_upgrade_enabled  = false
  type_handler_version       = "1.22.0"

  protected_setting = <<PROTECTED_SETTING
{
	"workspaceKey": "${azurerm_log_analytics_workspace.test.primary_shared_key}"
}
PROTECTED_SETTING

  setting = <<SETTING
{
	"workspaceId": "${azurerm_log_analytics_workspace.test.workspace_id}"
}
SETTING
}
`, r.template(data), data.RandomString)
}

func (r StackHCIExtensionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-%s"
  location            = data.azurerm_resource_group.test.location
  resource_group_name = data.azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_stack_hci_extension" "test" {
  name                       = "acctest-shce-%[2]s"
  arc_setting_id             = data.azurerm_stack_hci_cluster_arc_setting.test.id
  publisher                  = "Microsoft.EnterpriseCloud.Monitoring"
  type                       = "MicrosoftMonitoringAgent"
  auto_upgrade_minor_version = true
  automatic_upgrade_enabled  = false
  force_update_tag           = "1"
  type_handler_version       = "1.22.0"

  protected_setting = <<PROTECTED_SETTING
{
	"workspaceKey": "${azurerm_log_analytics_workspace.test.primary_shared_key}"
}
PROTECTED_SETTING

  setting = <<SETTING
{
	"workspaceId": "${azurerm_log_analytics_workspace.test.workspace_id}"
}
SETTING
}
`, r.template(data), data.RandomString)
}

func (r StackHCIExtensionResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)

	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_extension" "import" {
  name              = azurerm_stack_hci_extension.test.name
  arc_setting_id    = azurerm_stack_hci_extension.test.arc_setting_id
  publisher         = azurerm_stack_hci_extension.test.publisher
  type              = azurerm_stack_hci_extension.test.type
  protected_setting = azurerm_stack_hci_extension.test.protected_setting
  setting           = azurerm_stack_hci_extension.test.setting
}
`, config)
}

// nolint: unparam
func (r StackHCIExtensionResource) template(data acceptance.TestData) string {
	stackHCIClusterName := os.Getenv("ARM_TEST_HCI_CLUSTER_NAME")
	resourceGroupName := os.Getenv("ARM_TEST_HCI_RESOURCE_GROUP_NAME")
	arcSettingName := os.Getenv("ARM_TEST_HCI_ARC_SETTING_NAME")

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_stack_hci_cluster" "test" {
  name                = "%s"
  resource_group_name = "%s"
}

data "azurerm_stack_hci_cluster_arc_setting" "test" {
  name                 = "%s"
  stack_hci_cluster_id = data.azurerm_stack_hci_cluster.test.id
}

data "azurerm_resource_group" "test" {
  name = "%[2]s"
}
`, stackHCIClusterName, resourceGroupName, arcSettingName)
}
