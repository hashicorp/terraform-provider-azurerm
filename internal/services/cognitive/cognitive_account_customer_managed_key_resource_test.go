package cognitive_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2021-04-30/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CognitiveAccountCustomerManagedKeyResource struct{}

func TestAccCognitiveAccountCustomerManagedKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_customer_managed_key", "test")
	r := CognitiveAccountCustomerManagedKeyResource{}

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

func TestAccCognitiveAccountCustomerManagedKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_customer_managed_key", "test")
	r := CognitiveAccountCustomerManagedKeyResource{}

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

func TestAccCognitiveAccountCustomerManagedKey_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_customer_managed_key", "test")
	r := CognitiveAccountCustomerManagedKeyResource{}

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

func TestAccCognitiveAccountCustomerManagedKey_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_customer_managed_key", "test")
	r := CognitiveAccountCustomerManagedKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CognitiveAccountCustomerManagedKeyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := cognitiveservicesaccounts.ParseAccountID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cognitive.AccountsClient.AccountsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", *id, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Encryption == nil || resp.Model.Properties.Encryption.KeySource == nil {
		return utils.Bool(false), nil
	}

	if *resp.Model.Properties.Encryption.KeySource == cognitiveservicesaccounts.KeySourceMicrosoftPointCognitiveServices {
		return utils.Bool(false), nil
	}

	return utils.Bool(true), nil
}

func (r CognitiveAccountCustomerManagedKeyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_cognitive_account_customer_managed_key" "test" {
  cognitive_account_id = azurerm_cognitive_account.test.id
  key_vault_key_id     = azurerm_key_vault_key.test.id
}
`, r.template(data))
}

func (r CognitiveAccountCustomerManagedKeyResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s
resource "azurerm_cognitive_account_customer_managed_key" "import" {
  cognitive_account_id = azurerm_cognitive_account_customer_managed_key.test.cognitive_account_id
  key_vault_key_id     = azurerm_cognitive_account_customer_managed_key.test.key_vault_key_id
}
`, template)
}

func (r CognitiveAccountCustomerManagedKeyResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_cognitive_account_customer_managed_key" "test" {
  cognitive_account_id = azurerm_cognitive_account.test.id
  key_vault_key_id     = azurerm_key_vault_key.test.id
  identity_client_id   = azurerm_user_assigned_identity.test.client_id
}
`, r.template(data))
}

func (r CognitiveAccountCustomerManagedKeyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "%s"
}

resource "azurerm_cognitive_account" "test" {
  name                  = "acctest-cogacc-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  kind                  = "SpeechServices"
  sku_name              = "S0"
  custom_subdomain_name = "acctest-cogacc-%d"
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true

  access_policy {
    tenant_id = azurerm_cognitive_account.test.identity.0.tenant_id
    object_id = azurerm_cognitive_account.test.identity.0.principal_id
    key_permissions = [
      "Get", "Create", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify"
    ]
    secret_permissions = [
      "Get",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id
    key_permissions = [
      "Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify"
    ]
    secret_permissions = [
      "Get",
    ]
  }

  access_policy {
    tenant_id = azurerm_user_assigned_identity.test.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id
    key_permissions = [
      "Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify"
    ]
    secret_permissions = [
      "Get",
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctestkvkey%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomString)
}
