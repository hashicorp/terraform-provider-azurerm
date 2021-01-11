package loganalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type LogAnalyticsClusterCustomerManagedKeyResource struct {
}

func TestAccLogAnalyticsClusterCustomerManagedKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_cluster_customer_managed_key", "test")
	r := LogAnalyticsClusterCustomerManagedKeyResource{}

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

func (t LogAnalyticsClusterCustomerManagedKeyResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.LogAnalyticsClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.LogAnalytics.ClusterClient.Get(ctx, id.ResourceGroup, id.ClusterName)
	if err != nil {
		return nil, fmt.Errorf("readingLog Analytics Cluster Customer Managed Key (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (LogAnalyticsClusterCustomerManagedKeyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-la-%[1]d"
  location = "%[2]s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_log_analytics_cluster" "test" {
  name                = "acctest-LA-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}


resource "azurerm_key_vault" "test" {
  name                = "vault%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  soft_delete_enabled        = true
  soft_delete_retention_days = 7
  purge_protection_enabled   = true
}


resource "azurerm_key_vault_access_policy" "terraform" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "create",
    "delete",
    "get",
    "list",
    "purge",
    "update",
  ]

  secret_permissions = [
    "get",
    "delete",
    "set",
  ]

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_key" "test" {
  name         = "key-%[3]s"
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

  depends_on = [azurerm_key_vault_access_policy.terraform]
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "get",
    "unwrapkey",
    "wrapkey"
  ]

  tenant_id = azurerm_log_analytics_cluster.test.identity.0.tenant_id
  object_id = azurerm_log_analytics_cluster.test.identity.0.principal_id

  depends_on = [azurerm_key_vault_access_policy.terraform]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r LogAnalyticsClusterCustomerManagedKeyResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_cluster_customer_managed_key" "test" {
  log_analytics_cluster_id = azurerm_log_analytics_cluster.test.id
  key_vault_key_id         = azurerm_key_vault_key.test.id

  depends_on = [azurerm_key_vault_access_policy.test]
}

`, r.template(data))
}
