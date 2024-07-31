// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AppServiceEnvironmentV3Resource struct{}

func TestAccAppServiceEnvironmentV3_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment_v3", "test")
	r := AppServiceEnvironmentV3Resource{}

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

func TestAccAppServiceEnvironmentV3_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment_v3", "test")
	r := AppServiceEnvironmentV3Resource{}

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

func TestAccAppServiceEnvironmentV3_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment_v3", "test")
	r := AppServiceEnvironmentV3Resource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("pricing_tier").Exists(),
				check.That(data.ResourceName).Key("cluster_setting.#").HasValue("3"),
				check.That(data.ResourceName).Key("dns_suffix").HasValue(fmt.Sprintf("acctest-ase-%d.appserviceenvironment.net", data.RandomInteger)),
				check.That(data.ResourceName).Key("inbound_network_dependencies.#").HasValue("3"),
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("windows_outbound_ip_addresses.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServiceEnvironmentV3_completeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment_v3", "test")
	r := AppServiceEnvironmentV3Resource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServiceEnvironmentV3_updateVnet(t *testing.T) {
	if true {
		t.Skip("Skipping as `ChangeVnet` func not actually implemented in service")
	}
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment_v3", "test")
	r := AppServiceEnvironmentV3Resource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateVnet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG2-ase-%[2]d/providers/Microsoft.Network/virtualNetworks/acctest-vnet-%[2]d/subnets/acctest-subnet2-%[2]d", data.Client().SubscriptionID, data.RandomInteger)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServiceEnvironmentV3_zoneRedundant(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment_v3", "test")
	r := AppServiceEnvironmentV3Resource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zoneRedundant(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppServiceEnvironmentV3_dedicatedHosts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment_v3", "test")
	r := AppServiceEnvironmentV3Resource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dedicatedHosts(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (AppServiceEnvironmentV3Resource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AppServiceEnvironmentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Web.AppServiceEnvironmentsClient.Get(ctx, id.ResourceGroup, id.HostingEnvironmentName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving App Service Environment V3 %q (Resource Group %q): %+v", id.HostingEnvironmentName, id.ResourceGroup, err)
	}

	return utils.Bool(true), nil
}

func (r AppServiceEnvironmentV3Resource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_app_service_environment_v3" "test" {
  name                = "acctest-ase-%d"
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.test.id
}
`, template, data.RandomInteger)
}

func (r AppServiceEnvironmentV3Resource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_app_service_environment_v3" "test" {
  name                         = "acctest-ase-%d"
  resource_group_name          = azurerm_resource_group.test2.name
  subnet_id                    = azurerm_subnet.test.id
  internal_load_balancing_mode = "Web, Publishing"
  remote_debugging_enabled     = true

  cluster_setting {
    name  = "InternalEncryption"
    value = "true"
  }

  cluster_setting {
    name  = "DisableTls1.0"
    value = "1"
  }

  cluster_setting {
    name  = "FrontEndSSLCipherSuiteOrder"
    value = "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"
  }

  tags = {
    accTest = "1"
    env     = "Test"
  }
}
`, template, data.RandomInteger)
}

func (r AppServiceEnvironmentV3Resource) completeUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_app_service_environment_v3" "test" {
  name                         = "acctest-ase-%d"
  resource_group_name          = azurerm_resource_group.test2.name
  subnet_id                    = azurerm_subnet.test.id
  internal_load_balancing_mode = "Web, Publishing"

  allow_new_private_endpoint_connections = false

  cluster_setting {
    name  = "InternalEncryption"
    value = "true"
  }

  cluster_setting {
    name  = "DisableTls1.0"
    value = "1"
  }

  cluster_setting {
    name  = "FrontEndSSLCipherSuiteOrder"
    value = "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384_P256,TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256_P256,TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384_P256,TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256_P256,TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA_P256,TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA_P256"
  }

  tags = {
    accTest = "1"
    env     = "Test"
    tags    = "Updated"
  }
}
`, template, data.RandomInteger)
}

func (r AppServiceEnvironmentV3Resource) updateVnet(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`

%s

resource "azurerm_virtual_network" "test2" {
  name                = "acctest-vnet2-%[2]d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
  address_space       = ["20.0.0.0/16"]
}

resource "azurerm_subnet" "test2" {
  name                 = "acctest-subnet2-%[2]d"
  resource_group_name  = azurerm_virtual_network.test2.resource_group_name
  virtual_network_name = azurerm_virtual_network.test2.name
  address_prefixes     = ["20.0.2.0/24"]
  delegation {
    name = "asedelegation"
    service_delegation {
      name    = "Microsoft.Web/hostingEnvironments"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_app_service_environment_v3" "test" {
  name                = "acctest-ase-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.test2.id
}
`, template, data.RandomInteger)
}

func (r AppServiceEnvironmentV3Resource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_environment_v3" "import" {
  name                = azurerm_app_service_environment_v3.test.name
  resource_group_name = azurerm_app_service_environment_v3.test.resource_group_name
  subnet_id           = azurerm_app_service_environment_v3.test.subnet_id
}
`, template)
}

func (r AppServiceEnvironmentV3Resource) zoneRedundant(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_app_service_environment_v3" "test" {
  name                = "acctest-ase-%d"
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.test.id
  zone_redundant      = true
}
`, template, data.RandomInteger)
}

func (r AppServiceEnvironmentV3Resource) dedicatedHosts(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_app_service_environment_v3" "test" {
  name                 = "acctest-ase-%d"
  resource_group_name  = azurerm_resource_group.test.name
  subnet_id            = azurerm_subnet.test.id
  dedicated_host_count = 2
}
`, template, data.RandomInteger)
}

func (r AppServiceEnvironmentV3Resource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ase-%[1]d"
  location = "%[2]s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG2-ase-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-subnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  delegation {
    name = "asedelegation"
    service_delegation {
      name    = "Microsoft.Web/hostingEnvironments"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
