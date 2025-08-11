package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataFactoryCustomerManagedKeyTestResource struct{}

func TestAccDataFactoryCustomerManagedKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_customer_managed_key", "test")
	r := DataFactoryCustomerManagedKeyTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemAssignedKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryCustomerManagedKey_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_customer_managed_key", "test")
	r := DataFactoryCustomerManagedKeyTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userAssignedIdentityKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryCustomerManagedKey_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_customer_managed_key", "test")
	r := DataFactoryCustomerManagedKeyTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemAssignedKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.systemAssignedKeyUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (DataFactoryCustomerManagedKeyTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := factories.ParseFactoryID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.DataFactory.Factories.Get(ctx, *id, factories.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.Encryption != nil), nil
}

func (r DataFactoryCustomerManagedKeyTestResource) systemAssignedKey(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_data_factory_customer_managed_key" "test" {
  data_factory_id         = azurerm_data_factory.test.id
  customer_managed_key_id = azurerm_key_vault_key.test.id
}
`, r.systemAssignedTemplate(data))
}

func (r DataFactoryCustomerManagedKeyTestResource) systemAssignedKeyUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_customer_managed_key" "test" {
  data_factory_id         = azurerm_data_factory.test.id
  customer_managed_key_id = azurerm_key_vault_key.test2.id
}
`, r.systemAssignedTemplate(data))
}

func (r DataFactoryCustomerManagedKeyTestResource) userAssignedIdentityKey(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_customer_managed_key" "test" {
  data_factory_id           = azurerm_data_factory.test.id
  customer_managed_key_id   = azurerm_key_vault_key.test.id
  user_assigned_identity_id = azurerm_user_assigned_identity.test.id
}
`, r.userAssignedTemplate(data))
}

func (r DataFactoryCustomerManagedKeyTestResource) systemAssignedTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_access_policy" "datafactory" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_data_factory.test.identity[0].tenant_id
  object_id    = azurerm_data_factory.test.identity[0].principal_id

  key_permissions = [
    "Get",
    "WrapKey",
    "UnwrapKey",
  ]
}

resource "azurerm_data_factory" "test" {
  name                = "acctest%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "SystemAssigned"
  }
}


`, r.template(data), data.RandomInteger)
}

func (r DataFactoryCustomerManagedKeyTestResource) userAssignedTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory" "test" {
  name                = "acctest%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }
}
`, r.template(data), data.RandomInteger)
}

func (r DataFactoryCustomerManagedKeyTestResource) template(data acceptance.TestData) string {
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
  name     = "acctestRG-df-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_key_vault" "test" {
  name                     = "acckv%[1]d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Recover",
    "Update",
    "GetRotationPolicy",
  ]

  secret_permissions = [
    "Delete",
    "Get",
    "Set",
  ]
}

resource "azurerm_key_vault_access_policy" "userassigned" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_user_assigned_identity.test.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id

  key_permissions = [
    "Get",
    "UnwrapKey",
    "WrapKey"
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "key"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "unwrapKey",
    "wrapKey"
  ]

  depends_on = [azurerm_key_vault_access_policy.test, azurerm_key_vault_access_policy.userassigned]
}

resource "azurerm_key_vault_key" "test2" {
  name         = "key2"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "unwrapKey",
    "wrapKey"
  ]

  depends_on = [azurerm_key_vault_access_policy.test, azurerm_key_vault_access_policy.userassigned]
}
`, data.RandomInteger, data.Locations.Primary)
}
