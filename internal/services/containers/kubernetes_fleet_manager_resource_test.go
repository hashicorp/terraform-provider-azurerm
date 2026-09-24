package containers_test

// NOTE: this file is generated - manual changes will be overwritten.
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.
import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KubernetesFleetManagerTestResource struct{}

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
				check.That(data.ResourceName).Key("hub_profile.0.agent_profile.0.virtual_machine_size").HasValue("Standard_DS2_v2"),
				check.That(data.ResourceName).Key("hub_profile.0.api_server_access_profile.#").HasValue("1"),
				check.That(data.ResourceName).Key("hub_profile.0.api_server_access_profile.0.enable_private_cluster").HasValue("false"),
				check.That(data.ResourceName).Key("hub_profile.0.dns_prefix").HasValue(fmt.Sprintf("val-%s", data.RandomString)),
				check.That(data.ResourceName).Key("hub_profile.0.fqdn").Exists(),
				check.That(data.ResourceName).Key("hub_profile.0.kubernetes_version").Exists(),
				check.That(data.ResourceName).Key("hub_profile.0.portal_fqdn").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesFleetManager_privateHub(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_fleet_manager", "test")
	r := KubernetesFleetManagerTestResource{}
	subnetId := r.preCheckPrivateHub(t)

	checks := []acceptance.TestCheckFunc{
		check.That(data.ResourceName).ExistsInAzure(r),
		check.That(data.ResourceName).Key("hub_profile.#").HasValue("1"),
		check.That(data.ResourceName).Key("hub_profile.0.agent_profile.#").HasValue("1"),
		check.That(data.ResourceName).Key("hub_profile.0.agent_profile.0.subnet_id").HasValue(subnetId),
		check.That(data.ResourceName).Key("hub_profile.0.agent_profile.0.virtual_machine_size").HasValue("Standard_D2as_v7"),
		check.That(data.ResourceName).Key("hub_profile.0.api_server_access_profile.#").HasValue("1"),
		check.That(data.ResourceName).Key("hub_profile.0.api_server_access_profile.0.enable_private_cluster").HasValue("true"),
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
		"hub_profile.0.api_server_access_profile.0.enable_private_cluster",
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
			Config: r.privateHub(data, subnetId, true, "initial"),
			Check: acceptance.ComposeTestCheckFunc(
				hubCheck,
				check.That(data.ResourceName).Key("tags.phase").HasValue("initial"),
			),
		},
		data.ImportStep(),
		{
			Config: r.privateHub(data, subnetId, false, "updated"),
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

func TestKubernetesFleetManagerPrivateHubPreCheck(t *testing.T) {
	expected := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctest/providers/Microsoft.Network/virtualNetworks/test/subnets/agents"
	for name, subnetId := range map[string]string{
		"canonical":           expected,
		"mixed_case_segments": "/SUBSCRIPTIONS/00000000-0000-0000-0000-000000000000/RESOURCEGROUPS/acctest/PROVIDERS/microsoft.network/VIRTUALNETWORKS/test/SUBNETS/agents",
	} {
		t.Run(name, func(t *testing.T) {
			t.Setenv("ARM_TEST_FLEET_SUBNET_ID", subnetId)
			actual := (KubernetesFleetManagerTestResource{}).preCheckPrivateHub(t)
			if actual != expected {
				t.Fatalf("expected normalized subnet ID %q, got %q", expected, actual)
			}
		})
	}
}

func TestKubernetesFleetManagerPrivateHubConfig(t *testing.T) {
	data := acceptance.TestData{RandomInteger: 123, RandomString: "abcde"}
	data.Locations.Primary = "eastus"
	subnetId := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctest/providers/Microsoft.Network/virtualNetworks/test/subnets/agents"
	r := KubernetesFleetManagerTestResource{}

	for _, test := range []struct {
		name              string
		configureProfiles bool
		phase             string
	}{
		{name: "create", configureProfiles: true, phase: "initial"},
		{name: "profile_omission_update", configureProfiles: false, phase: "updated"},
	} {
		t.Run(test.name, func(t *testing.T) {
			config := r.privateHub(data, subnetId, test.configureProfiles, test.phase)
			for _, block := range []string{
				`provider "azurerm"`,
				`resource "azurerm_resource_group" "test"`,
				`data "azurerm_kubernetes_fleet_manager" "test"`,
				`resource "azurerm_kubernetes_fleet_manager" "test"`,
				`dynamic "agent_profile"`,
				`dynamic "api_server_access_profile"`,
			} {
				if count := strings.Count(config, block); count != 1 {
					t.Errorf("expected one %s block, got %d", block, count)
				}
			}
			for _, value := range []string{
				`default = "eastus"`,
				fmt.Sprintf("default = %t", test.configureProfiles),
				fmt.Sprintf("phase = %q", test.phase),
				fmt.Sprintf("subnet_id            = %q", subnetId),
				`virtual_machine_size = "Standard_D2as_v7"`,
				"enable_private_cluster = true",
			} {
				if !strings.Contains(config, value) {
					t.Errorf("expected rendered configuration to contain %q", value)
				}
			}
			if count := strings.Count(config, "for_each = var.configure_profiles ? [1] : []"); count != 2 {
				t.Errorf("expected both optional profiles to use the omission switch, got %d", count)
			}
			if strings.Contains(config, "dns_prefix") {
				t.Error("private hubs must omit dns_prefix")
			}
			if strings.Contains(config, "%!") {
				t.Error("configuration contains an unresolved formatting directive")
			}
		})
	}
}

func (KubernetesFleetManagerTestResource) preCheckPrivateHub(t *testing.T) string {
	t.Helper()
	subnetId := os.Getenv("ARM_TEST_FLEET_SUBNET_ID")
	if subnetId == "" {
		t.Skip("ARM_TEST_FLEET_SUBNET_ID must reference a dedicated subnet with Network Contributor assigned to the Fleet resource provider")
	}

	id, err := commonids.ParseSubnetIDInsensitively(subnetId)
	if err != nil {
		t.Fatalf("parsing ARM_TEST_FLEET_SUBNET_ID: %+v", err)
	}
	return id.ID()
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
      enable_private_cluster = false
    }

    dns_prefix = "val-${var.random_string}"
  }
}
`, r.template(data))
}

func (r KubernetesFleetManagerTestResource) privateHub(data acceptance.TestData, subnetId string, configureProfiles bool, phase string) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
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
        subnet_id            = %q
        virtual_machine_size = "Standard_D2as_v7"
      }
    }
    dynamic "api_server_access_profile" {
      for_each = var.configure_profiles ? [1] : []
      content {
        enable_private_cluster = true
      }
    }
  }
}
`, r.template(data), configureProfiles, phase, subnetId)
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
