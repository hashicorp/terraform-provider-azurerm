package aiservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AIServicesHub struct{}

func TestAccAIServicesHub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ai_services_hub", "test")
	r := AIServicesHub{}

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

func TestAccAIServicesHub_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ai_services_hub", "test")
	r := AIServicesHub{}

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

func TestAccAIServicesHub_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ai_services_hub", "test")
	r := AIServicesHub{}

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

func TestAccAIServicesHub_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ai_services_hub", "test")
	r := AIServicesHub{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
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

func TestAccAIServicesHub_encryptionWithSystemAssignedId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ai_services_hub", "test")
	r := AIServicesHub{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.encryptionWithSystemAssignedId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAIServicesHub_encryptionWithUserAssignedId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ai_services_hub", "test")
	r := AIServicesHub{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.encryptionWithUserAssignedId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (AIServicesHub) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := workspaces.ParseWorkspaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.MachineLearning.Workspaces.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r AIServicesHub) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

%s

resource "azurerm_ai_services_hub" "test" {
  name                = "acctestaihub-%[2]d"
  location            = azurerm_ai_services.test.location
  resource_group_name = azurerm_resource_group.test.name
  storage_account_id  = azurerm_storage_account.test.id
  key_vault_id        = azurerm_key_vault.test.id

  identity {
    type = "SystemAssigned"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r AIServicesHub) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

%s

resource "azurerm_application_insights" "test" {
  name                = "acctestai-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  // Premium sku is required when creating a hub with AllowInternetOutbound or AllowOnlyApprovedOutbound isolation mode
  sku = "Premium"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_ai_services_hub" "test" {
  name                = "acctestaihub-%[2]d"
  location            = azurerm_ai_services.test.location
  resource_group_name = azurerm_resource_group.test.name
  storage_account_id  = azurerm_storage_account.test.id
  key_vault_id        = azurerm_key_vault.test.id

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  application_insights_id        = azurerm_application_insights.test.id
  container_registry_id          = azurerm_container_registry.test.id
  primary_user_assigned_identity = azurerm_user_assigned_identity.test.id
  public_network_access          = "Disabled"
  image_build_compute_name       = "buildtest"
  description                    = "AI Hub created by Terraform"
  friendly_name                  = "AI Hub"
  high_business_impact_enabled   = true

  managed_network {
    isolation_mode = "AllowInternetOutbound"
  }

  tags = {
    env = "test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r AIServicesHub) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

%s

resource "azurerm_application_insights" "test" {
  name                = "acctestai-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctestuai2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_ai_services_hub" "test" {
  name                = "acctestaihub-%[2]d"
  location            = azurerm_ai_services.test.location
  resource_group_name = azurerm_resource_group.test.name
  storage_account_id  = azurerm_storage_account.test.id
  key_vault_id        = azurerm_key_vault.test.id

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
      azurerm_user_assigned_identity.test2.id
    ]
  }

  application_insights_id        = azurerm_application_insights.test.id
  container_registry_id          = azurerm_container_registry.test.id
  primary_user_assigned_identity = azurerm_user_assigned_identity.test2.id
  public_network_access          = "Enabled"
  image_build_compute_name       = "buildtestupdated"
  description                    = "AI Hub for Projects"
  friendly_name                  = "AI Hub for OS models"
  high_business_impact_enabled   = true

  managed_network {
    isolation_mode = "AllowOnlyApprovedOutbound"
  }

  tags = {
    env   = "prod"
    model = "regression"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r AIServicesHub) encryptionWithUserAssignedId(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_key_vault_access_policy" "test-uai-policy" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_user_assigned_identity.test.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id
  key_permissions = [
    "Get",
    "Recover",
    "UnwrapKey",
    "WrapKey",
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "acckvKey-%[2]d"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
  depends_on = [azurerm_key_vault.test, azurerm_key_vault_access_policy.test]
}

resource "azurerm_role_assignment" "test_kv" {
  scope                = azurerm_key_vault.test.id
  role_definition_name = "Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_ai_services_hub" "test" {
  name                = "acctestaihub-%[2]d"
  location            = azurerm_ai_services.test.location
  resource_group_name = azurerm_resource_group.test.name
  storage_account_id  = azurerm_storage_account.test.id
  key_vault_id        = azurerm_key_vault.test.id

  primary_user_assigned_identity = azurerm_user_assigned_identity.test.id

  encryption {
    user_assigned_identity_id = azurerm_user_assigned_identity.test.id
    key_vault_id              = azurerm_key_vault.test.id
    key_id                    = azurerm_key_vault_key.test.id
  }

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }

  depends_on = [
    azurerm_role_assignment.test_kv,
    azurerm_key_vault_access_policy.test-uai-policy
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r AIServicesHub) encryptionWithSystemAssignedId(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

%s

resource "azurerm_key_vault_key" "test" {
  name         = "acckvkey-%[2]d"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
  depends_on = [azurerm_key_vault.test, azurerm_key_vault_access_policy.test]
}

resource "azurerm_ai_services_hub" "test" {
  name                = "acctestaihub-%[2]d"
  location            = azurerm_ai_services.test.location
  resource_group_name = azurerm_resource_group.test.name
  storage_account_id  = azurerm_storage_account.test.id
  key_vault_id        = azurerm_key_vault.test.id

  encryption {
    key_vault_id = azurerm_key_vault.test.id
    key_id       = azurerm_key_vault_key.test.id
  }

  identity {
    type = "SystemAssigned"
  }
}
`, r.template(data), data.RandomInteger)
}

func (AIServicesHub) requiresImport(data acceptance.TestData) string {
	template := AIServicesHub{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_ai_services_hub" "import" {
  name                = azurerm_ai_services_hub.test.name
  location            = azurerm_ai_services_hub.test.location
  resource_group_name = azurerm_ai_services_hub.test.resource_group_name
  storage_account_id  = azurerm_ai_services_hub.test.storage_account_id
  key_vault_id        = azurerm_ai_services_hub.test.key_vault_id

  identity {
    type = "SystemAssigned"
  }
}
`, template)
}

func (AIServicesHub) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aiservices-%[1]d"
  location = "%[2]s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestvault%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  purge_protection_enabled   = true
  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Get",
    "Delete",
    "Purge",
    "GetRotationPolicy",
  ]
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_ai_services" "test" {
  name                = "acctestaiservices-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "S0"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
