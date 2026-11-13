// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sreagents_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-07-01/agents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SreAgentResource struct{}

func TestAccSreAgent_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sre_agent", "test")
	r := SreAgentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep(),
	})
}

func TestAccSreAgent_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sre_agent", "test")
	r := SreAgentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSreAgent_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sre_agent", "test")
	r := SreAgentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("endpoint").Exists(),
				check.That(data.ResourceName).Key("outbound_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
			),
		},
		data.ImportStep("application_insights.0.connection_string", "incident_management.0.connection_key"),
	})
}

func TestAccSreAgent_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sre_agent", "test")
	r := SreAgentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep(),
		{
			Config: r.minimalBlocks(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("application_insights.0.connection_string"),
		{
			Config: r.updated(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("application_insights.0.connection_string", "incident_management.0.connection_key"),
		{
			Config: r.complete(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("application_insights.0.connection_string", "incident_management.0.connection_key"),
		{
			Config:      r.basic(data),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile("`subnet_id` is required when `sandbox\\.egress\\.mode` is set to `AzureVNet` and must not be set otherwise"),
		},
		{
			Config: r.updated(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("application_insights.0.connection_string", "incident_management.0.connection_key"),
		{
			Config: r.minimalBlocks(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("application_insights.0.connection_string"),
		{
			Config: r.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep(),
	})
}

func TestAccSreAgent_invalidNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sre_agent", "test")
	r := SreAgentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.networkWithoutSubnet(data),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile("`subnet_id` is required when `sandbox\\.egress\\.mode` is set to `AzureVNet` and must not be set otherwise"),
		},
		{
			Config:      r.subnetWithoutNetwork(data),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile("`subnet_id` is required when `sandbox\\.egress\\.mode` is set to `AzureVNet` and must not be set otherwise"),
		},
	})
}

func TestAccSreAgent_invalidIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sre_agent", "test")
	r := SreAgentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.unassignedActionIdentity(data),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile("action.user_assigned_identity_id"),
		},
		{
			Config:      r.unassignedKnowledgeIdentity(data),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile("knowledge_graph.user_assigned_identity_id"),
		},
	})
}

func (SreAgentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := agents.ParseAgentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SREAgents.AgentsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (SreAgentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-sre-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-identity-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-law-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "PerGB2018"
}

resource "azurerm_application_insights" "test" {
  name                = "acctest-appinsights-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  application_type    = "web"
  workspace_id        = azurerm_log_analytics_workspace.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r SreAgentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    application_insights {
      disable_generated_rule = true
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-sre-%[1]d"
  location = "%[2]s"
}

resource "azurerm_sre_agent" "test" {
  name                = "acctest-sre-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r SreAgentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_sre_agent" "import" {
  name                = azurerm_sre_agent.test.name
  resource_group_name = azurerm_sre_agent.test.resource_group_name
  location            = azurerm_sre_agent.test.location

  identity {
    type = azurerm_sre_agent.test.identity[0].type
  }
}
`, r.basic(data))
}

func (r SreAgentResource) minimalBlocks(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    application_insights {
      disable_generated_rule = true
    }
  }
}

%[1]s

resource "azurerm_sre_agent" "test" {
  name                = "acctest-sre-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  action {
    user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }

  application_insights {
    application_id    = azurerm_application_insights.test.app_id
    connection_string = azurerm_application_insights.test.connection_string
  }

  default_model {
    name     = "Automatic"
    provider = "MicrosoftFoundry"
  }

  incident_management {
    type = "AzMonitor"
  }

  knowledge_graph {
    user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }

  sandbox {
    egress {
      mode = "Limited"

      http_proxy {
        https_url = "http://proxy-updated.example.invalid:8080"
      }
    }

    package {
      name            = "requests"
      package_manager = "pip"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r SreAgentResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    application_insights {
      disable_generated_rule = true
    }
  }
}

%[1]s

resource "azurerm_sre_agent" "test" {
  name                = "acctest-sre-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  upgrade_channel     = "Preview"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  action {
    access_level              = "Low"
    mode                      = "Review"
    user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }

  application_insights {
    application_id    = azurerm_application_insights.test.app_id
    connection_string = azurerm_application_insights.test.connection_string
  }

  default_model {
    name     = "Automatic"
    provider = "MicrosoftFoundry"
  }

  incident_management {
    connection_key    = "acceptance-not-a-real-credential-updated"
    connection_name   = "azmonitor"
    connection_url    = "https://incident-updated.example.invalid"
    on_behalf_of_user = "sre-updated@example.invalid"
    type              = "AzMonitor"
  }

  knowledge_graph {
    managed_resource_ids      = [azurerm_resource_group.test.id]
    user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }

  sandbox {
    egress {
      mode                                   = "Limited"
      http_mcp_server_network_access_enabled = true
      allowed_code_repositories              = ["Github"]
      allowed_hosts                          = ["*.example.com"]
      allowed_registries                     = ["pypi"]

      http_proxy {
        https_url = "http://proxy-updated.example.invalid:8080"
        http_url  = "http://proxy-updated.example.invalid:8080"
        no_proxy  = ["localhost"]
      }
    }

    package {
      name            = "requests"
      package_manager = "pip"
      version         = "2.32.5"
    }
  }

  tags = {
    environment = "acceptance"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r SreAgentResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    application_insights {
      disable_generated_rule = true
    }
  }
}

%[1]s

resource "azurerm_user_assigned_identity" "second" {
  name                = "acctest-identity2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_application_insights" "second" {
  name                = "acctest-appinsights2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  application_type    = "web"
  workspace_id        = azurerm_log_analytics_workspace.test.id
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-subnet-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.0.0/24"]

  delegation {
    name = "sre-agent"

    service_delegation {
      name    = "Microsoft.App/environments"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_sre_agent" "test" {
  name                = "acctest-sre-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.test.id
  upgrade_channel     = "Stable"

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id, azurerm_user_assigned_identity.second.id]
  }

  action {
    access_level              = "High"
    mode                      = "Autonomous"
    user_assigned_identity_id = azurerm_user_assigned_identity.second.id
  }

  application_insights {
    application_id    = azurerm_application_insights.second.app_id
    connection_string = azurerm_application_insights.second.connection_string
  }

  default_model {
    name     = "Automatic"
    provider = "MicrosoftFoundry"
  }

  incident_management {
    connection_key    = "acceptance-not-a-real-credential-complete"
    connection_name   = "azmonitor"
    connection_url    = "https://incident-complete.example.invalid"
    on_behalf_of_user = "sre-complete@example.invalid"
    type              = "AzMonitor"
  }

  knowledge_graph {
    managed_resource_ids      = [azurerm_resource_group.test.id, azurerm_virtual_network.test.id]
    user_assigned_identity_id = azurerm_user_assigned_identity.second.id
  }

  sandbox {
    egress {
      mode                                   = "AzureVNet"
      http_mcp_server_network_access_enabled = false
      allowed_code_repositories              = ["Github", "AzureDevOps"]
      allowed_hosts                          = ["*.example.com", "example.org"]
      allowed_registries                     = ["pypi", "npmjs"]

      http_proxy {
        https_url = "http://proxy-complete.example.invalid:8080"
        http_url  = "http://proxy-complete.example.invalid:8080"
        no_proxy  = ["localhost", "127.0.0.1"]
      }
    }

    package {
      name            = "requests"
      package_manager = "pip"
      version         = "2.32.4"
    }

    package {
      name            = "urllib3"
      package_manager = "pip"
      version         = "2.5.0"
    }
  }

  tags = {
    environment = "acceptance"
    purpose     = "complete"
  }
}


`, r.template(data), data.RandomInteger)
}

func (r SreAgentResource) networkWithoutSubnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    application_insights {
      disable_generated_rule = true
    }
  }
}

%[1]s

resource "azurerm_sre_agent" "test" {
  name                = "acctest-sre-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }

  sandbox {
    egress {
      mode = "AzureVNet"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r SreAgentResource) subnetWithoutNetwork(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    application_insights {
      disable_generated_rule = true
    }
  }
}

%[1]s

resource "azurerm_sre_agent" "test" {
  name                = "acctest-sre-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Network/virtualNetworks/test/subnets/test"

  identity {
    type = "SystemAssigned"
  }

  sandbox {
    egress {
      mode = "Limited"
    }

    package {
      name            = "requests"
      package_manager = "pip"
      version         = azurerm_user_assigned_identity.test.principal_id
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r SreAgentResource) unassignedActionIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    application_insights {
      disable_generated_rule = true
    }
  }
}

%[1]s

resource "azurerm_sre_agent" "test" {
  name                = "acctest-sre-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = ["/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.ManagedIdentity/userAssignedIdentities/test"]
  }

  action {
    user_assigned_identity_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.ManagedIdentity/userAssignedIdentities/other"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r SreAgentResource) unassignedKnowledgeIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    application_insights {
      disable_generated_rule = true
    }
  }
}

%[1]s

resource "azurerm_sre_agent" "test" {
  name                = "acctest-sre-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = ["/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.ManagedIdentity/userAssignedIdentities/test"]
  }

  knowledge_graph {
    user_assigned_identity_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.ManagedIdentity/userAssignedIdentities/other"
  }
}
`, r.template(data), data.RandomInteger)
}
