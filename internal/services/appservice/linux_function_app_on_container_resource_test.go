package appservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LinuxFunctionAppOnContainerResource struct{}

func TestAccLinuxFunctionAppOnContainer_basicPublicRegistry(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_function_app_on_container", "test")
	r := LinuxFunctionAppOnContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicPublicRegistry(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxFunctionAppOnContainer_completePublicRegistry(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_function_app_on_container", "test")
	r := LinuxFunctionAppOnContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completePublicRegistry(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxFunctionAppOnContainer_storageUsesMSI(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_function_app_on_container", "test")
	r := LinuxFunctionAppOnContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccountUsesMsi(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxFunctionAppOnContainer_storageAccountKeyVaultAccess(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_function_app_on_container", "test")
	r := LinuxFunctionAppOnContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccountKeyVaultAccess(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxFunctionAppOnContainer_basicMCRUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_function_app_on_container", "test")
	r := LinuxFunctionAppOnContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicPublicRegistry(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basicMcrUpdateAppSetting(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basicPublicRegistry(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (r LinuxFunctionAppOnContainerResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseFunctionAppID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.AppService.WebAppsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Linux Functions App %s: %+v", id, err)
	}
	if response.WasNotFound(resp.HttpResponse) {
		return utils.Bool(false), nil
	}
	return utils.Bool(true), nil
}

func (r LinuxFunctionAppOnContainerResource) basicPublicRegistry(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_function_app_on_container" "test" {
  name                         = "acctest-lfa-%s"
  resource_group_name          = azurerm_resource_group.test.name
  container_app_environment_id = azurerm_container_app_environment.test.id


  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    application_insights_connection_string = azurerm_application_insights.test.connection_string
    app_scale_limit                        = 10
  }

  app_settings = {
    foo    = "bar"
    secret = "sauce"
  }

  container_image = "azure-functions/dotnet:3.0-appservice-quickstart"

  registry {
    registry_server_url = "mcr.microsoft.com"
    registry_username   = ""
    registry_password   = ""
  }
}
`, r.template(data), data.RandomStringOfLength(5))
}

func (r LinuxFunctionAppOnContainerResource) completePublicRegistry(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_function_app_on_container" "test" {
  name                         = "acctest-lfa-%s"
  resource_group_name          = azurerm_resource_group.test.name
  container_app_environment_id = azurerm_container_app_environment.test.id


  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    application_insights_connection_string = azurerm_application_insights.test.connection_string
    app_scale_limit                        = 10
    elastic_instance_minimum               = 10
  }

  app_settings = {
    foo    = "bar"
    secret = "sauce"
  }

  container_image = "azure-functions/dotnet:3.0-appservice-quickstart"

  registry {
    registry_server_url = "mcr.microsoft.com"
    registry_username   = ""
    registry_password   = ""
  }

  identity {
    type = "SystemAssigned"
  }
}
`, r.template(data), data.RandomStringOfLength(5))
}

func (r LinuxFunctionAppOnContainerResource) storageAccountUsesMsi(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_function_app_on_container" "test" {
  name                         = "acctest-lfa-%s"
  resource_group_name          = azurerm_resource_group.test.name
  container_app_environment_id = azurerm_container_app_environment.test.id


  storage_account_name          = azurerm_storage_account.test.name
  storage_uses_managed_identity = true

  site_config {
  }

  app_settings = {
    "test" = "value"
  }

  container_image = "azure-functions/dotnet8-quickstart-demo:1.0"

  registry {
    registry_server_url = "mcr.microsoft.com"
    registry_username   = ""
    registry_password   = ""
  }
}
`, r.template(data), data.RandomStringOfLength(5))
}

func (r LinuxFunctionAppOnContainerResource) storageAccountKeyVaultAccess(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy    = true
      recover_soft_deleted_key_vaults = true
    }
  }
}

%s

data "azurerm_client_config" "current" {}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv-%[2]s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Get",
    ]

    secret_permissions = [
      "Get",
      "Delete",
      "List",
      "Purge",
      "Recover",
      "Set",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

    secret_permissions = [
      "Get",
      "List",
    ]
  }

  tags = {
    environment = "AccTest"
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%[2]s"
  value        = "DefaultEndpointsProtocol=https;AccountName=${azurerm_storage_account.test.name};AccountKey=${azurerm_storage_account.test.primary_access_key};EndpointSuffix=core.windows.net"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_linux_function_app_on_container" "test" {
  name                         = "acctest-lfa-%[3]s"
  resource_group_name          = azurerm_resource_group.test.name
  container_app_environment_id = azurerm_container_app_environment.test.id

  key_vault_reference_identity_id = azurerm_user_assigned_identity.test.id
  storage_key_vault_secret_id     = azurerm_key_vault_secret.test.versionless_id

  site_config {
  }

  app_settings = {
    "test" = "value"
  }

  container_image = "azure-functions/dotnet8-quickstart-demo:1.0"

  registry {
    registry_server_url = "mcr.microsoft.com"
    registry_username   = ""
    registry_password   = ""
  }

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, r.template(data), data.RandomString, data.RandomStringOfLength(5))
}

func (r LinuxFunctionAppOnContainerResource) basicMcrUpdateAppSetting(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_function_app_on_container" "test" {
  name                         = "acctest-lfa-%s"
  resource_group_name          = azurerm_resource_group.test.name
  container_app_environment_id = azurerm_container_app_environment.test.id


  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
  }

  app_settings = {
    "test"      = "value"
    "updatekey" = "updatevalue"
  }

  container_image = "azure-functions/dotnet8-quickstart-demo:1.0"

  registry {
    registry_server_url = "mcr.microsoft.com"
    registry_username   = ""
    registry_password   = ""
  }
}
`, r.template(data), data.RandomStringOfLength(5))
}

func (r LinuxFunctionAppOnContainerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "rg-sfi-acctestRG-LFA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_container_app_environment" "test" {
  name                = "acctest-CAEnv%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}
