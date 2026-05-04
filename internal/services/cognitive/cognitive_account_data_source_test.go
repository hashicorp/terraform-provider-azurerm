// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type CognitiveAccountDataSource struct{}

func TestAccCognitiveAccountDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cognitive_account", "test")
	r := CognitiveAccountDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kind").HasValue("Face"),
				check.That(data.ResourceName).Key("local_auth_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("outbound_network_access_restricted").HasValue("false"),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("project_management_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
	})
}

func TestAccCognitiveAccountDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cognitive_account", "test")
	r := CognitiveAccountDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kind").HasValue("AIServices"),
				check.That(data.ResourceName).Key("outbound_network_access_restricted").HasValue("false"),
				check.That(data.ResourceName).Key("dynamic_throttling_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("project_management_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("local_auth_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned, UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("sku_name").HasValue("S0"),
				check.That(data.ResourceName).Key("tags.Acceptance").HasValue("Test"),
				check.That(data.ResourceName).Key("custom_subdomain_name").IsNotEmpty(),
				check.That(data.ResourceName).Key("fqdns.#").HasValue("2"),
				check.That(data.ResourceName).Key("network_acls.0.default_action").HasValue("Deny"),
				check.That(data.ResourceName).Key("network_acls.0.bypass").HasValue("AzureServices"),
				check.That(data.ResourceName).Key("network_acls.0.virtual_network_rules.0.ignore_missing_vnet_service_endpoint").HasValue("false"),
				check.That(data.ResourceName).Key("network_acls.0.virtual_network_rules.0.subnet_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("network_acls.0.ip_rules.#").HasValue("1"),
				check.That(data.ResourceName).Key("network_injection.0.scenario").HasValue(string(cognitiveservicesaccounts.ScenarioTypeAgent)),
				check.That(data.ResourceName).Key("network_injection.0.subnet_id").Exists(),
			),
		},
	})
}

func TestAccCognitiveAccountDataSource_customerManagedKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cognitive_account", "test")
	r := CognitiveAccountDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.customerManagedKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("customer_managed_key.0.key_vault_key_id").Exists(),
				check.That(data.ResourceName).Key("customer_managed_key.0.identity_client_id").IsUUID(),
			),
		},
	})
}

func TestAccCognitiveAccountDataSource_speechServicesWithStorage(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cognitive_account", "test")
	r := CognitiveAccountDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.speechServicesWithStorage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kind").HasValue("SpeechServices"),
				check.That(data.ResourceName).Key("storage.0.storage_account_id").Exists(),
				check.That(data.ResourceName).Key("storage.0.identity_client_id").Exists(),
			),
		},
	})
}

func TestAccCognitiveAccountDataSource_qnaRuntimeEndpoint(t *testing.T) {
	t.Skipf("skipping as there is no available quota for kind `QnAMaker`")
	data := acceptance.BuildTestData(t, "data.azurerm_cognitive_account", "test")
	r := CognitiveAccountDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.qnaRuntimeEndpoint(data, "https://localhost:8080/"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kind").HasValue("QnAMaker"),
				check.That(data.ResourceName).Key("qna_runtime_endpoint").HasValue("https://localhost:8080/"),
				check.That(data.ResourceName).Key("endpoint").Exists(),
			),
		},
	})
}

func TestAccCognitiveAccountDataSource_metricsAdvisor(t *testing.T) {
	t.Skipf("skipping as there is no available quota for kind `MetricsAdvisor`")
	data := acceptance.BuildTestData(t, "data.azurerm_cognitive_account", "test")
	r := CognitiveAccountDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.metricsAdvisor(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kind").HasValue("MetricsAdvisor"),
				check.That(data.ResourceName).Key("metrics_advisor_aad_client_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("metrics_advisor_aad_tenant_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("metrics_advisor_super_user_name").IsNotEmpty(),
				check.That(data.ResourceName).Key("metrics_advisor_website_name").IsNotEmpty(),
			),
		},
	})
}

func (CognitiveAccountDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Face"
  sku_name            = "S0"

  tags = {
    Acceptance = "Test"
  }
}

data "azurerm_cognitive_account" "test" {
  name                = azurerm_cognitive_account.test.name
  resource_group_name = azurerm_cognitive_account.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (CognitiveAccountDataSource) complete(data acceptance.TestData) string {
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
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network" "test_agent" {
  name                = "acctestvirtnetagent%[1]d"
  address_space       = ["192.168.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test_a" {
  name                 = "acctestsubneta%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  service_endpoints    = ["Microsoft.CognitiveServices"]
}

resource "azurerm_subnet" "test_b" {
  name                 = "acctestsubnetb%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.4.0/24"]
  service_endpoints    = ["Microsoft.CognitiveServices"]
}

resource "azurerm_subnet" "test_agent" {
  name                 = "acctestsubnetaagent%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test_agent.name
  address_prefixes     = ["192.168.0.0/24"]

  delegation {
    name = "Microsoft.App/environments"

    service_delegation {
      name = "Microsoft.App/environments"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action"
      ]
    }
  }
}

resource "azurerm_cognitive_account" "test" {
  name                               = "acctestaiservices-%[1]d"
  location                           = azurerm_resource_group.test.location
  resource_group_name                = azurerm_resource_group.test.name
  kind                               = "AIServices"
  sku_name                           = "S0"
  fqdns                              = ["foo.com", "bar.com"]
  local_auth_enabled                 = true
  outbound_network_access_restricted = false
  project_management_enabled         = true
  public_network_access_enabled      = true
  custom_subdomain_name              = "acctestaiservices-%[1]d"

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }

  network_acls {
    default_action = "Deny"
    bypass         = "AzureServices"
    ip_rules       = ["123.0.0.101"]
    virtual_network_rules {
      subnet_id = azurerm_subnet.test_a.id
    }
    virtual_network_rules {
      subnet_id = azurerm_subnet.test_b.id
    }
  }

  network_injection {
    scenario  = "agent"
    subnet_id = azurerm_subnet.test_agent.id
  }

  tags = {
    Acceptance = "Test"
  }
}

data "azurerm_cognitive_account" "test" {
  name                = azurerm_cognitive_account.test.name
  resource_group_name = azurerm_cognitive_account.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (CognitiveAccountDataSource) qnaRuntimeEndpoint(data acceptance.TestData, url string) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_cognitive_account" "test" {
  name                = azurerm_cognitive_account.test.name
  resource_group_name = azurerm_cognitive_account.test.resource_group_name
}
`, CognitiveAccountResource{}.qnaRuntimeEndpoint(data, url))
}

func (CognitiveAccountDataSource) metricsAdvisor(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_cognitive_account" "test" {
  name                = azurerm_cognitive_account.test.name
  resource_group_name = azurerm_cognitive_account.test.resource_group_name
}
`, CognitiveAccountResource{}.metricsAdvisor(data))
}

func (CognitiveAccountDataSource) customerManagedKey(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_cognitive_account" "test" {
  name                = azurerm_cognitive_account.test.name
  resource_group_name = azurerm_cognitive_account.test.resource_group_name
}
`, CognitiveAccountResource{}.customerManagedKey(data))
}

func (CognitiveAccountDataSource) speechServicesWithStorage(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_cognitive_account" "test" {
  name                = azurerm_cognitive_account.test.name
  resource_group_name = azurerm_cognitive_account.test.resource_group_name
}
`, CognitiveAccountResource{}.speechServicesWithStorage(data))
}
