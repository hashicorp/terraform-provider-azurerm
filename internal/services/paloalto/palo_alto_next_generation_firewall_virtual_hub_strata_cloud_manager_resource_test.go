// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package paloalto_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-05-23/firewalls"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NextGenerationFirewallVHubStrataCloudManagerResource struct{}

func TestAccNextGenerationFirewallVHubStrataCloudManager_basic(t *testing.T) {
	if scmTenant := os.Getenv("ARM_PALO_ALTO_SCM_TENANT_NAME"); scmTenant == "" {
		t.Skipf("skipping as Palo Alto Strata Cloud Manager tenant name not set in `ARM_PALO_ALTO_SCM_TENANT_NAME`")
	}

	data := acceptance.BuildTestData(t, "azurerm_palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager", "test")
	r := NextGenerationFirewallVHubStrataCloudManagerResource{}

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

func TestAccNextGenerationFirewallVHubStrataCloudManager_requiresImport(t *testing.T) {
	if scmTenant := os.Getenv("ARM_PALO_ALTO_SCM_TENANT_NAME"); scmTenant == "" {
		t.Skipf("skipping as Palo Alto Strata Cloud Manager tenant name not set in `ARM_PALO_ALTO_SCM_TENANT_NAME`")
	}

	data := acceptance.BuildTestData(t, "azurerm_palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager", "test")
	r := NextGenerationFirewallVHubStrataCloudManagerResource{}

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

func TestAccNextGenerationFirewallVHubStrataCloudManager_complete(t *testing.T) {
	if scmTenant := os.Getenv("ARM_PALO_ALTO_SCM_TENANT_NAME"); scmTenant == "" {
		t.Skipf("skipping as Palo Alto Strata Cloud Manager tenant name not set in `ARM_PALO_ALTO_SCM_TENANT_NAME`")
	}

	data := acceptance.BuildTestData(t, "azurerm_palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager", "test")
	r := NextGenerationFirewallVHubStrataCloudManagerResource{}

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

func TestAccNextGenerationFirewallVHubStrataCloudManager_update(t *testing.T) {
	if scmTenant := os.Getenv("ARM_PALO_ALTO_SCM_TENANT_NAME"); scmTenant == "" {
		t.Skipf("skipping as Palo Alto Strata Cloud Manager tenant name not set in `ARM_PALO_ALTO_SCM_TENANT_NAME`")
	}

	data := acceptance.BuildTestData(t, "azurerm_palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager", "test")
	r := NextGenerationFirewallVHubStrataCloudManagerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data), // Skip plan_id update as it's only applicable when a new plan is introduced, thus it's not able to be tested
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

func (r NextGenerationFirewallVHubStrataCloudManagerResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := firewalls.ParseFirewallID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.PaloAlto.PaloAltoClient_v2025_05_23.Firewalls.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r NextGenerationFirewallVHubStrataCloudManagerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager" "test" {
  name                             = "acctest-ngfwvhscm-%[2]d"
  resource_group_name              = azurerm_resource_group.test.name
  location                         = azurerm_resource_group.test.location
  strata_cloud_manager_tenant_name = "%[3]s"

  network_profile {
    virtual_hub_id               = azurerm_virtual_hub.test.id
    network_virtual_appliance_id = azurerm_palo_alto_virtual_network_appliance.test.id
    public_ip_address_ids        = [azurerm_public_ip.test.id]
  }

  // tags is required in the test subscription account, otherwise it fails
  tags = {
    userid = "terraform-test"
  }
}
`, r.template(data), data.RandomInteger, os.Getenv("ARM_PALO_ALTO_SCM_TENANT_NAME"))
}

func (r NextGenerationFirewallVHubStrataCloudManagerResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager" "import" {
  name                             = azurerm_palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager.test.name
  resource_group_name              = azurerm_palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager.test.resource_group_name
  location                         = azurerm_palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager.test.location
  strata_cloud_manager_tenant_name = azurerm_palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager.test.strata_cloud_manager_tenant_name

  network_profile {
    virtual_hub_id               = azurerm_palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager.test.network_profile[0].virtual_hub_id
    network_virtual_appliance_id = azurerm_palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager.test.network_profile[0].network_virtual_appliance_id
    public_ip_address_ids        = azurerm_palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager.test.network_profile[0].public_ip_address_ids
  }

  tags = {
    userid = "terraform-test"
  }
}
`, template)
}

func (r NextGenerationFirewallVHubStrataCloudManagerResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager" "test" {
  name                             = "acctest-ngfwvhscm-%[2]d"
  resource_group_name              = azurerm_resource_group.test.name
  location                         = azurerm_resource_group.test.location
  strata_cloud_manager_tenant_name = "%[3]s"

  network_profile {
    virtual_hub_id               = azurerm_virtual_hub.test.id
    network_virtual_appliance_id = azurerm_palo_alto_virtual_network_appliance.test.id
    public_ip_address_ids        = [azurerm_public_ip.test.id]
    egress_nat_ip_address_ids    = [azurerm_public_ip.egress.id]
    trusted_address_ranges       = ["20.22.92.11"]
  }

  dns_settings {
    dns_servers = ["8.8.8.8", "8.8.4.4"]
  }

  destination_nat {
    name     = "testDNAT-1"
    protocol = "TCP"
    frontend_config {
      public_ip_address_id = azurerm_public_ip.test.id
      port                 = 8081
    }
    backend_config {
      public_ip_address = "10.0.1.101"
      port              = 18081
    }
  }

  destination_nat {
    name     = "testDNAT-2"
    protocol = "UDP"
    frontend_config {
      public_ip_address_id = azurerm_public_ip.test.id
      port                 = 8082
    }
    backend_config {
      public_ip_address = "10.0.1.102"
      port              = 18082
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    userid = "terraform"
  }
}
`, r.template(data), data.RandomInteger, os.Getenv("ARM_PALO_ALTO_SCM_TENANT_NAME"))
}

func (r NextGenerationFirewallVHubStrataCloudManagerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-PANGFWVHSCM-%[1]d"
  location = "%[2]s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"

  depends_on = [azurerm_public_ip.egress]
}

resource "azurerm_public_ip" "egress" {
  name                = "acctestpublicip-%[1]d-e"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctestVHUB-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"

  tags = {
    hubSaaSPreview = "true"
  }
}

resource "azurerm_palo_alto_virtual_network_appliance" "test" {
  name           = "testAcc-panva-%[1]d"
  virtual_hub_id = azurerm_virtual_hub.test.id
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uaid-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary)
}
