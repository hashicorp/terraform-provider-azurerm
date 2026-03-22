// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machinelearning_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/onlineendpoint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type OnlineEndpointResource struct{}

func TestAccMachineLearningOnlineEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_online_endpoint", "test")
	r := OnlineEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMachineLearningOnlineEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_online_endpoint", "test")
	r := OnlineEndpointResource{}

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

func TestAccMachineLearningOnlineEndpoint_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_online_endpoint", "test")
	r := OnlineEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMachineLearningOnlineEndpoint_authModeAADToken(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_online_endpoint", "test")
	r := OnlineEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authMode(data, "AADToken"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMachineLearningOnlineEndpoint_authModeAMLToken(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_online_endpoint", "test")
	r := OnlineEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authMode(data, "AMLToken"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMachineLearningOnlineEndpoint_authModeKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_online_endpoint", "test")
	r := OnlineEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authMode(data, "Key"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMachineLearningOnlineEndpoint_publicNetworkAccessDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_online_endpoint", "test")
	r := OnlineEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.publicNetworkAccessDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMachineLearningOnlineEndpoint_identityUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_online_endpoint", "test")
	r := OnlineEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identityUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMachineLearningOnlineEndpoint_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_online_endpoint", "test")
	r := OnlineEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r OnlineEndpointResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	onlineEndpointClient := client.MachineLearning.OnlineEndpoints

	id, err := onlineendpoint.ParseOnlineEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := onlineEndpointClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Properties.AuthMode != ""), nil
}

func (r OnlineEndpointResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_online_endpoint" "test" {
  name                          = "acctest-mle-%[2]d"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  location                      = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomIntOfLength(8))
}

func (r OnlineEndpointResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_online_endpoint" "import" {
  name                          = azurerm_machine_learning_online_endpoint.test.name
  machine_learning_workspace_id = azurerm_machine_learning_online_endpoint.test.machine_learning_workspace_id
  location                      = azurerm_machine_learning_online_endpoint.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, config)
}

func (r OnlineEndpointResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_online_endpoint" "test" {
  name                          = "acctest-mle-%[2]d"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  location                      = azurerm_resource_group.test.location
  authentication_mode           = "AMLToken"
  description                   = "Test Online Endpoint"
  public_network_access_enabled = false

  identity {
    type = "SystemAssigned"
  }

  properties = {
    "test-key" = "test-value"
  }

  tags = {
    environment = "test"
  }
}
`, template, data.RandomIntOfLength(8))
}

func (r OnlineEndpointResource) authMode(data acceptance.TestData, mode string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_online_endpoint" "test" {
  name                          = "acctest-mle-%[2]d"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  location                      = azurerm_resource_group.test.location
  authentication_mode           = "%[3]s"

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomIntOfLength(8), mode)
}

func (r OnlineEndpointResource) publicNetworkAccessDisabled(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_online_endpoint" "test" {
  name                          = "acctest-mle-%[2]d"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  location                      = azurerm_resource_group.test.location
  public_network_access_enabled = false

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomIntOfLength(8))
}

func (r OnlineEndpointResource) identityUserAssigned(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_machine_learning_online_endpoint" "test" {
  name                          = "acctest-mle-%[2]d"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  location                      = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, template, data.RandomIntOfLength(8))
}

func (r OnlineEndpointResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_online_endpoint" "test" {
  name                          = "acctest-mle-%[2]d"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  location                      = azurerm_resource_group.test.location
  public_network_access_enabled = false

  identity {
    type = "SystemAssigned"
  }

  tags = {
    environment = "staging"
    purpose     = "testing"
  }
}
`, template, data.RandomIntOfLength(8))
}

func (r OnlineEndpointResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ml-%[1]d"
  location = "%[2]s"
  tags = {
    "stage" = "test"
  }
}

resource "azurerm_application_insights" "test" {
  name                = "acctestai-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_key_vault" "test" {
  name                = "acctest%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  purge_protection_enabled   = true
  soft_delete_retention_days = 7
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_machine_learning_workspace" "test" {
  name                    = "acctest-MLW%[4]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomIntOfLength(16))
}
