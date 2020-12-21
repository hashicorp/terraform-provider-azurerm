package eventhub_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type EventHubNamespaceDisasterRecoveryConfigResource struct {
}

func TestAccEventHubNamespaceDisasterRecoveryConfig_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_disaster_recovery_config", "test")
	r := EventHubNamespaceDisasterRecoveryConfigResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespaceDisasterRecoveryConfig_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_disaster_recovery_config", "test")
	r := EventHubNamespaceDisasterRecoveryConfigResource{}

	// skipping due to there being no way to delete a DRC once an alternate name has been set
	// sdk bug: https://github.com/Azure/azure-sdk-for-go/issues/5893
	t.Skip()

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespaceDisasterRecoveryConfig_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_disaster_recovery_config", "test")
	r := EventHubNamespaceDisasterRecoveryConfigResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated_removed(data),
		},
	})
}

func (EventHubNamespaceDisasterRecoveryConfigResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	name := id.Path["disasterRecoveryConfigs"]
	resourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]

	resp, err := clients.Eventhub.DisasterRecoveryConfigsClient.Get(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		return nil, fmt.Errorf("retrieving EventHub Namespace Disaster Recovery Configs %q (namespace %q / resource group: %q): %v", name, namespaceName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.ArmDisasterRecoveryProperties != nil), nil
}

func (EventHubNamespaceDisasterRecoveryConfigResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "testa" {
  name                = "acctest-EHN-%[1]d-a"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace" "testb" {
  name                = "acctest-EHN-%[1]d-b"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace_disaster_recovery_config" "test" {
  name                 = "acctest-EHN-DRC-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  namespace_name       = azurerm_eventhub_namespace.testa.name
  partner_namespace_id = azurerm_eventhub_namespace.testb.id
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

// nolint unused - mistakenly marked as unused
func (EventHubNamespaceDisasterRecoveryConfigResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "testa" {
  name                = "acctest-EHN-%[1]d-a"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace" "testb" {
  name                = "acctest-EHN-%[1]d-b"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace_disaster_recovery_config" "test" {
  name                 = "${azurerm_eventhub_namespace.testa.name}-111"
  resource_group_name  = azurerm_resource_group.test.name
  namespace_name       = azurerm_eventhub_namespace.testa.name
  partner_namespace_id = azurerm_eventhub_namespace.testb.id
  alternate_name       = "acctest-EHN-DRC-%[1]d-alt"
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (EventHubNamespaceDisasterRecoveryConfigResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "testa" {
  name                = "acctest-EHN-%[1]d-a"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace" "testb" {
  name                = "acctest-EHN-%[1]d-b"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace" "testc" {
  name                = "acctest-EHN-%[1]d-c"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace_disaster_recovery_config" "test" {
  name                 = "acctest-EHN-DRC-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  namespace_name       = azurerm_eventhub_namespace.testa.name
  partner_namespace_id = azurerm_eventhub_namespace.testc.id
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (EventHubNamespaceDisasterRecoveryConfigResource) updated_removed(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "testa" {
  name                = "acctest-EHN-%[1]d-a"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace" "testb" {
  name                = "acctest-EHN-%[1]d-b"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace" "testc" {
  name                = "acctest-EHN-%[1]d-c"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}
