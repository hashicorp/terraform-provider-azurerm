package cognitive_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/accountconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CognitiveAccountConnectionAADResource struct{}

func TestAccCognitiveAccountConnectionAAD_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_aad", "test")
	r := CognitiveAccountConnectionAADResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("metadata"),
	})
}

func TestAccCognitiveAccountConnectionAAD_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_aad", "test")
	r := CognitiveAccountConnectionAADResource{}

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

func TestAccCognitiveAccountConnectionAAD_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_aad", "test")
	r := CognitiveAccountConnectionAADResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("metadata"),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("metadata"),
	})
}

func (r CognitiveAccountConnectionAADResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := accountconnectionresource.ParseConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cognitive.AccountConnectionResourceClient.AccountConnectionsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r CognitiveAccountConnectionAADResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-ac-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cognitive_account" "test" {
  name                       = "acctest-aiservices-%[1]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  kind                       = "AIServices"
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "acctestaiservices-%[1]d"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestsc%[3]s"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r CognitiveAccountConnectionAADResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cognitive_account_connection_aad" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "AzureBlob"
  target               = azurerm_storage_account.test.primary_blob_endpoint

  metadata = {
    accountName   = azurerm_storage_account.test.name
    containerName = azurerm_storage_container.test.name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CognitiveAccountConnectionAADResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account_connection_aad" "import" {
  name                 = azurerm_cognitive_account_connection_aad.test.name
  cognitive_account_id = azurerm_cognitive_account_connection_aad.test.cognitive_account_id
  category             = azurerm_cognitive_account_connection_aad.test.category
  target               = azurerm_cognitive_account_connection_aad.test.target

  metadata = {
    accountName   = azurerm_cognitive_account_connection_aad.test.metadata.accountName
    containerName = azurerm_cognitive_account_connection_aad.test.metadata.containerName
  }
}
`, r.basic(data))
}

func (r CognitiveAccountConnectionAADResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_storage_container" "test2" {
  name                  = "acctestsc2%[3]s"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_cognitive_account_connection_aad" "test" {
  name                 = "acctest-conn-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "AzureBlob"
  target               = azurerm_storage_account.test.primary_blob_endpoint

  metadata = {
    accountName   = azurerm_storage_account.test.name
    containerName = azurerm_storage_container.test2.name
  }
}
`, r.template(data), data.RandomInteger, data.RandomString)
}
