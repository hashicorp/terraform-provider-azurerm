// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_manager_test

// NOTE: this file is generated - manual changes will be overwritten.
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.
import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KubernetesFleetManagerTestResource struct{}

func TestAccKubernetesFleetManager_regressionTest(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_manager", "test")
	r := KubernetesFleetManagerTestResource{}

	data.ResourceRegressionTest(t, r, []acceptance.TestStep{
		{
			Config: r.updated(data),
		},
	}, "")
}

func TestAccKubernetesFleetManager_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_manager", "test")
	r := KubernetesFleetManagerTestResource{}

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

func TestAccKubernetesFleetManager_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_manager", "test")
	r := KubernetesFleetManagerTestResource{}

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

func TestAccKubernetesFleetManager_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_manager", "test")
	r := KubernetesFleetManagerTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("hub_profile.#").HasValue("1"),
				check.That(data.ResourceName).Key("hub_profile.0.agent_profile.#").HasValue("1"),
				check.That(data.ResourceName).Key("hub_profile.0.agent_profile.0.subnet_id").HasValue(""),
				check.That(data.ResourceName).Key("hub_profile.0.api_server_access_profile.#").HasValue("1"),
				check.That(data.ResourceName).Key("hub_profile.0.fqdn").Exists(),
				check.That(data.ResourceName).Key("hub_profile.0.kubernetes_version").Exists(),
				check.That(data.ResourceName).Key("hub_profile.0.portal_fqdn").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesFleetManager_hubProfileDefaults(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_manager", "test")
	r := KubernetesFleetManagerTestResource{}

	hubCheck := acceptance.ComposeTestCheckFunc(
		check.That(data.ResourceName).ExistsInAzure(r),
		check.That(data.ResourceName).Key("hub_profile.#").HasValue("1"),
		check.That(data.ResourceName).Key("hub_profile.0.fqdn").Exists(),
		check.That(data.ResourceName).Key("hub_profile.0.kubernetes_version").Exists(),
		check.That(data.ResourceName).Key("hub_profile.0.portal_fqdn").Exists(),
	)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.hubProfileDefaults(data, "initial"),
			Check: acceptance.ComposeTestCheckFunc(
				hubCheck,
				check.That(data.ResourceName).Key("tags.phase").HasValue("initial"),
			),
		},
		data.ImportStep(),
		{
			Config:   r.hubProfileDefaults(data, "initial"),
			PlanOnly: true,
		},
		{
			Config: r.hubProfileDefaults(data, "updated"),
			Check: acceptance.ComposeTestCheckFunc(
				hubCheck,
				check.That(data.ResourceName).Key("tags.phase").HasValue("updated"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesFleetManager_privateHub(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_manager", "test")
	r := KubernetesFleetManagerTestResource{}

	checks := []acceptance.TestCheckFunc{
		check.That(data.ResourceName).ExistsInAzure(r),
		check.That(data.ResourceName).Key("hub_profile.#").HasValue("1"),
		check.That(data.ResourceName).Key("hub_profile.0.agent_profile.#").HasValue("1"),
		check.That(data.ResourceName).Key("hub_profile.0.agent_profile.0.subnet_id").MatchesOtherKey(
			check.That("azurerm_subnet.test").Key("id"),
		),
		check.That(data.ResourceName).Key("hub_profile.0.agent_profile.0.virtual_machine_size").HasValue("Standard_D2as_v7"),
		check.That(data.ResourceName).Key("hub_profile.0.api_server_access_profile.#").HasValue("1"),
		check.That(data.ResourceName).Key("hub_profile.0.api_server_access_profile.0.private_cluster_enabled").HasValue("true"),
		check.That(data.ResourceName).Key("hub_profile.0.dns_prefix").HasValue(""),
		check.That(data.ResourceName).Key("hub_profile.0.fqdn").Exists(),
		check.That(data.ResourceName).Key("hub_profile.0.kubernetes_version").Exists(),
		check.That(data.ResourceName).Key("hub_profile.0.portal_fqdn").Exists(),
	}
	for _, key := range []string{
		"hub_profile.#",
		"hub_profile.0.agent_profile.#",
		"hub_profile.0.agent_profile.0.subnet_id",
		"hub_profile.0.agent_profile.0.virtual_machine_size",
		"hub_profile.0.api_server_access_profile.#",
		"hub_profile.0.api_server_access_profile.0.private_cluster_enabled",
		"hub_profile.0.dns_prefix",
		"hub_profile.0.fqdn",
		"hub_profile.0.kubernetes_version",
		"hub_profile.0.portal_fqdn",
		"tags.phase",
	} {
		checks = append(checks, check.That("data.azurerm_kubernetes_fleet_manager.test").Key(key).MatchesOtherKey(
			check.That(data.ResourceName).Key(key),
		))
	}
	hubCheck := acceptance.ComposeTestCheckFunc(checks...)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateHub(data, true, "initial"),
			Check: acceptance.ComposeTestCheckFunc(
				hubCheck,
				check.That(data.ResourceName).Key("tags.phase").HasValue("initial"),
			),
		},
		data.ImportStep(),
		{
			Config:   r.privateHub(data, false, "initial"),
			PlanOnly: true,
		},
		{
			Config: r.privateHub(data, false, "updated"),
			Check: acceptance.ComposeTestCheckFunc(
				hubCheck,
				check.That(data.ResourceName).Key("tags.phase").HasValue("updated"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesFleetManager_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_manager", "test")
	r := KubernetesFleetManagerTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated(data),
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

func (r KubernetesFleetManagerTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := fleets.ParseFleetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ContainerService.V20240401.Fleets.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r KubernetesFleetManagerTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_kubernetes_fleet_manager" "test" {
  location            = azurerm_resource_group.test.location
  name                = "acctestkfm-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
}
`, r.template(data))
}

func (r KubernetesFleetManagerTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_fleet_manager" "import" {
  location            = azurerm_kubernetes_fleet_manager.test.location
  name                = azurerm_kubernetes_fleet_manager.test.name
  resource_group_name = azurerm_kubernetes_fleet_manager.test.resource_group_name
}
	`, r.basic(data))
}

func (r KubernetesFleetManagerTestResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_kubernetes_fleet_manager" "test" {
  location            = azurerm_resource_group.test.location
  name                = "acctestkfm-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  tags = {
    environment = "terraform-acctests"
    some_key    = "some-value"
  }
}
`, r.template(data))
}

func (r KubernetesFleetManagerTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_kubernetes_fleet_manager" "test" {
  location            = azurerm_resource_group.test.location
  name                = "acctestkfm-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  tags = {
    environment = "terraform-acctests"
    some_key    = "some-value"
  }
  hub_profile {
    agent_profile {
      virtual_machine_size = "Standard_DS2_v2"
    }

    api_server_access_profile {
      private_cluster_enabled = false
    }

    dns_prefix = "val-${var.random_string}"
  }
}
`, r.template(data))
}

func (r KubernetesFleetManagerTestResource) hubProfileDefaults(data acceptance.TestData, phase string) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_kubernetes_fleet_manager" "test" {
  location            = azurerm_resource_group.test.location
  name                = "acctestkfm-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  tags = {
    phase = %q
  }
  hub_profile {}
}
`, r.template(data), phase)
}

func (r KubernetesFleetManagerTestResource) privateHub(data acceptance.TestData, configureProfiles bool, phase string) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

data "azuread_service_principal" "fleet" {
  display_name = "Azure Kubernetes Service - Fleet RP"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-${var.random_integer}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["192.168.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsnet-${var.random_integer}"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["192.168.0.0/24"]
}

resource "azurerm_role_assignment" "fleet_network" {
  scope                = azurerm_subnet.test.id
  role_definition_name = "Network Contributor"
  principal_id         = data.azuread_service_principal.fleet.object_id
}

variable "configure_profiles" {
  default = %t
}

data "azurerm_kubernetes_fleet_manager" "test" {
  name                = azurerm_kubernetes_fleet_manager.test.name
  resource_group_name = azurerm_kubernetes_fleet_manager.test.resource_group_name
}

resource "azurerm_kubernetes_fleet_manager" "test" {
  location            = azurerm_resource_group.test.location
  name                = "acctestkfm-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  tags = {
    phase = %q
  }
  hub_profile {
    dynamic "agent_profile" {
      for_each = var.configure_profiles ? [1] : []
      content {
        subnet_id            = azurerm_subnet.test.id
        virtual_machine_size = "Standard_D2as_v7"
      }
    }
    dynamic "api_server_access_profile" {
      for_each = var.configure_profiles ? [1] : []
      content {
        private_cluster_enabled = true
      }
    }
  }

  depends_on = [azurerm_role_assignment.fleet_network]
}
`, r.template(data), configureProfiles, phase)
}

func (r KubernetesFleetManagerTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "primary_location" {
  default = %q
}
variable "random_integer" {
  default = %d
}
variable "random_string" {
  default = %q
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-${var.random_integer}"
  location = var.primary_location
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
