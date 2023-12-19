// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package firewall_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-06-01/firewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type FirewallPolicyResource struct{}

func TestAccFirewallPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")
	r := FirewallPolicyResource{}

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

func TestAccFirewallPolicy_basicPremium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")
	r := FirewallPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicPremium(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")
	r := FirewallPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("dns.0.servers.#").HasValue("3"),
				check.That(data.ResourceName).Key("dns.0.servers.0").HasValue("1.1.1.1"),
				check.That(data.ResourceName).Key("dns.0.servers.1").HasValue("3.3.3.3"),
				check.That(data.ResourceName).Key("dns.0.servers.2").HasValue("2.2.2.2"),
				check.That(data.ResourceName).Key("dns.0.proxy_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallPolicy_completePremium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")
	r := FirewallPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completePremium(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")
	r := FirewallPolicyResource{}

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

func TestAccFirewallPolicy_updatePremium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")
	r := FirewallPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicPremium(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completePremium(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")
	r := FirewallPolicyResource{}

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

func TestAccFirewallPolicy_inherit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")
	r := FirewallPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.inherit(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallPolicy_insights(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy", "test")
	r := FirewallPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.defaultWorkspaceOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.regionalWorkspace(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.defaultWorkspaceOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (FirewallPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := firewallpolicies.ParseFirewallPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.FirewallPolicies.Get(ctx, *id, firewallpolicies.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (FirewallPolicyResource) basic(data acceptance.TestData) string {
	r := FirewallPolicyResource{}
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_firewall_policy" "test" {
  name                = "acctest-networkfw-Policy-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  tags = {
    Env = "Test"
  }
}
`, template, data.RandomInteger)
}

func (FirewallPolicyResource) basicPremium(data acceptance.TestData) string {
	r := FirewallPolicyResource{}
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_firewall_policy" "test" {
  name                = "acctest-networkfw-Policy-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
}
`, template, data.RandomInteger)
}

func (FirewallPolicyResource) pacFile(data acceptance.TestData) string {
	utcNow := time.Now().UTC()
	startDate := utcNow.Format(time.RFC3339)
	endDate := utcNow.Add(time.Hour * 24).Format(time.RFC3339)

	return fmt.Sprintf(`
resource "azurerm_storage_account" "test" {
  name                            = "acctestacc%[1]s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}

resource "azurerm_storage_blob" "test" {
  name                   = "example.pac"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  source_content         = "function FindProxyForURL(url, host) { return \"DIRECT\"; }"
}

data "azurerm_storage_account_sas" "test" {
  connection_string = azurerm_storage_account.test.primary_connection_string
  https_only        = true
  signed_version    = "2019-10-10"

  resource_types {
    service   = false
    container = false
    object    = true
  }

  services {
    blob  = true
    queue = false
    table = false
    file  = false
  }

  start  = "%[2]s"
  expiry = "%[3]s"

  permissions {
    read    = true
    write   = false
    delete  = false
    list    = false
    add     = false
    create  = false
    update  = false
    process = false
    tag     = false
    filter  = false
  }
}
`, data.RandomString, startDate, endDate)
}

func (FirewallPolicyResource) complete(data acceptance.TestData) string {
	r := FirewallPolicyResource{}
	template := r.template(data)
	return fmt.Sprintf(`

%s

%s

resource "azurerm_firewall_policy" "test" {
  name                     = "acctest-networkfw-Policy-%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  threat_intelligence_mode = "Off"
  threat_intelligence_allowlist {
    ip_addresses = ["1.1.1.1", "2.2.2.2", "10.0.0.0/16"]
    fqdns        = ["foo.com", "bar.com"]
  }
  explicit_proxy {
    enabled         = true
    http_port       = 8087
    https_port      = 8088
    enable_pac_file = true
    pac_file_port   = 8089
    pac_file        = "${azurerm_storage_blob.test.id}${data.azurerm_storage_account_sas.test.sas}&sr=b"
  }
  auto_learn_private_ranges_enabled = true
  dns {
    servers       = ["1.1.1.1", "3.3.3.3", "2.2.2.2"]
    proxy_enabled = true
  }
  tags = {
    env = "Test"
  }
}
`, template, FirewallPolicyResource{}.pacFile(data), data.RandomInteger)
}

func (FirewallPolicyResource) completePremium(data acceptance.TestData) string {
	r := FirewallPolicyResource{}
	template := r.templatePremium(data)
	return fmt.Sprintf(`


%s

%s

resource "azurerm_firewall_policy" "test" {
  name                     = "acctest-networkfw-Policy-%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  sku                      = "Premium"
  threat_intelligence_mode = "Off"
  threat_intelligence_allowlist {
    ip_addresses = ["1.1.1.1", "2.2.2.2", "10.0.0.0/16"]
    fqdns        = ["foo.com", "bar.com"]
  }
  explicit_proxy {
    enabled         = true
    http_port       = 8087
    https_port      = 8088
    enable_pac_file = true
    pac_file_port   = 8089
    pac_file        = "${azurerm_storage_blob.test.id}${data.azurerm_storage_account_sas.test.sas}&sr=b"
  }
  auto_learn_private_ranges_enabled = true
  dns {
    servers       = ["1.1.1.1", "2.2.2.2"]
    proxy_enabled = true
  }
  intrusion_detection {
    mode = "Alert"
    signature_overrides {
      state = "Alert"
      id    = "1"
    }
    private_ranges = ["172.111.111.111"]
    traffic_bypass {
      name                  = "Name bypass traffic settings"
      description           = "Description bypass traffic settings"
      destination_addresses = []
      source_addresses      = []
      protocol              = "Any"
      destination_ports     = ["*"]
      source_ip_groups = [
        azurerm_ip_group.test_source.id,
      ]
      destination_ip_groups = [
        azurerm_ip_group.test_destination.id,
      ]
    }
  }
  sql_redirect_allowed = true
  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
  tls_certificate {
    key_vault_secret_id = azurerm_key_vault_certificate.test.secret_id
    name                = azurerm_key_vault_certificate.test.name
  }
  private_ip_ranges = ["172.16.0.0/12", "192.168.0.0/16"]
  tags = {
    env = "Test"
  }
}
`, template, FirewallPolicyResource{}.pacFile(data), data.RandomInteger)
}

func (FirewallPolicyResource) requiresImport(data acceptance.TestData) string {
	r := FirewallPolicyResource{}
	template := r.basic(data)
	return fmt.Sprintf(`
%s
resource "azurerm_firewall_policy" "import" {
  name                = azurerm_firewall_policy.test.name
  resource_group_name = azurerm_firewall_policy.test.resource_group_name
  location            = azurerm_firewall_policy.test.location
}
`, template)
}

func (FirewallPolicyResource) inherit(data acceptance.TestData) string {
	r := FirewallPolicyResource{}
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_firewall_policy" "test-parent" {
  name                = "acctest-networkfw-Policy-%d-parent"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
resource "azurerm_firewall_policy" "test" {
  name                = "acctest-networkfw-Policy-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  base_policy_id      = azurerm_firewall_policy.test-parent.id
  threat_intelligence_allowlist {
    ip_addresses = ["1.1.1.1", "2.2.2.2"]
    fqdns        = ["foo.com", "bar.com"]
  }
  dns {
    servers       = ["1.1.1.1", "2.2.2.2"]
    proxy_enabled = true
  }
  tags = {
    env = "Test"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (FirewallPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-networkfw-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (FirewallPolicyResource) templatePremium(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-networkfw-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                            = "tlskv%d"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
  enabled_for_disk_encryption     = true
  enabled_for_deployment          = true
  enabled_for_template_deployment = true
  tenant_id                       = data.azurerm_client_config.current.tenant_id
  sku_name                        = "standard"
}

resource "azurerm_ip_group" "test_source" {
  name                = "acctestIpGroupForFirewallNetworkRulesSource"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["1.2.3.4/32", "12.34.56.0/24"]
}

resource "azurerm_ip_group" "test_destination" {
  name                = "acctestIpGroupForFirewallNetworkRulesDestination"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["192.168.0.0/25", "192.168.0.192/26"]
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id

  key_permissions = [
    "Create",
    "Get",
  ]

  certificate_permissions = [
    "Create",
    "Get",
    "List",
    "ManageContacts",
  ]

  secret_permissions = [
    "Get",
    "List",
  ]
}

resource "azurerm_key_vault_access_policy" "test2" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Backup",
    "Create",
    "Delete",
    "Get",
    "Import",
    "List",
    "Purge",
    "Recover",
    "Restore",
    "Update"
  ]

  certificate_permissions = [
    "Backup",
    "Create",
    "Get",
    "List",
    "Import",
    "Purge",
    "Delete",
    "Recover",
    "ManageContacts",
  ]

  secret_permissions = [
    "Get",
    "List",
    "Set",
    "Purge",
    "Delete",
    "Recover"
  ]
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "AzureFirewallPolicyCertificate"
  key_vault_id = azurerm_key_vault.test.id

  certificate {
    contents = filebase64("testdata/certificate.pfx")
    password = "somepassword"
  }

  depends_on = [azurerm_key_vault_access_policy.test, azurerm_key_vault_access_policy.test2]
}
`, data.RandomInteger, "westeurope", data.RandomInteger, data.RandomInteger)
}

func (FirewallPolicyResource) defaultWorkspaceOnly(data acceptance.TestData) string {
	r := FirewallPolicyResource{}
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "default" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_firewall_policy" "test" {
  name                = "acctest-networkfw-Policy-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  insights {
    enabled                            = true
    retention_in_days                  = 7
    default_log_analytics_workspace_id = azurerm_log_analytics_workspace.default.id
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (FirewallPolicyResource) regionalWorkspace(data acceptance.TestData) string {
	r := FirewallPolicyResource{}
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "default" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_log_analytics_workspace" "regional" {
  name                = "acctestLAW-region-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_firewall_policy" "test" {
  name                = "acctest-networkfw-Policy-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  insights {
    enabled                            = true
    retention_in_days                  = 7
    default_log_analytics_workspace_id = azurerm_log_analytics_workspace.default.id
    log_analytics_workspace {
      id                = azurerm_log_analytics_workspace.regional.id
      firewall_location = "%s"
    }
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.Locations.Primary)
}
