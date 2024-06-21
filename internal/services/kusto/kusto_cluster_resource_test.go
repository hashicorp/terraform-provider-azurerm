// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KustoClusterResource struct{}

func TestAccKustoCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

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

func TestAccKustoCluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allowed_fqdns.#").HasValue("1"),
				check.That(data.ResourceName).Key("allowed_fqdns.0").HasValue("255.255.255.0/24"),
				check.That(data.ResourceName).Key("allowed_ip_ranges.#").HasValue("1"),
				check.That(data.ResourceName).Key("allowed_ip_ranges.0").HasValue("0.0.0.0/0"),
				check.That(data.ResourceName).Key("outbound_network_access_restricted").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disk_encryption_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("streaming_ingestion_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("purge_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("public_ip_type").HasValue("IPv4"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disk_encryption_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("streaming_ingestion_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("purge_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("public_ip_type").HasValue("DualStack"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disk_encryption_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("streaming_ingestion_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("purge_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("public_ip_type").HasValue("IPv4"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoCluster_doubleEncryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.doubleEncryption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoCluster_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.label").HasValue("test"),
			),
		},
		{
			Config: r.withTagsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.label").HasValue("test1"),
				check.That(data.ResourceName).Key("tags.ENV").HasValue("prod"),
			),
		},
	})
}

func TestAccKustoCluster_sku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Dev(No SLA)_Standard_D11_v2"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("1"),
			),
		},
		{
			Config: r.skuUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_D11_v2"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("2"),
			),
		},
	})
}

func TestAccKustoCluster_zones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withZones(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zones.#").HasValue("1"),
				check.That(data.ResourceName).Key("zones.0").HasValue("1"),
			),
		},
	})
}

func TestAccKustoCluster_identitySystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("0"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoCluster_UserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.principal_id").HasValue(""),
			),
		},
	})
}

func TestAccKustoCluster_multipleAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned, UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
			),
		},
	})
}

func TestAccKustoCluster_vnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vnet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("virtual_network_configuration.#").HasValue("1"),
				check.That(data.ResourceName).Key("virtual_network_configuration.0.subnet_id").Exists(),
				check.That(data.ResourceName).Key("virtual_network_configuration.0.engine_public_ip_id").Exists(),
				check.That(data.ResourceName).Key("virtual_network_configuration.0.data_management_public_ip_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.vnetRemoved(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("virtual_network_configuration.#").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.vnet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("virtual_network_configuration.#").HasValue("1"),
				check.That(data.ResourceName).Key("virtual_network_configuration.0.subnet_id").Exists(),
				check.That(data.ResourceName).Key("virtual_network_configuration.0.engine_public_ip_id").Exists(),
				check.That(data.ResourceName).Key("virtual_network_configuration.0.data_management_public_ip_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoCluster_languageExtensions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.languageExtensions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("language_extensions.#").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.languageExtensionsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("language_extensions.#").HasValue("2"),
				check.That(data.ResourceName).Key("language_extensions.1").HasValue("R"),
			),
		},
		{
			Config: r.languageExtensionsRemove(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("language_extensions.#").HasValue("1"),
				check.That(data.ResourceName).Key("language_extensions.0").HasValue("R"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoCluster_optimizedAutoScale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.optimizedAutoScale(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("optimized_auto_scale.#").HasValue("1"),
				check.That(data.ResourceName).Key("optimized_auto_scale.0.minimum_instances").HasValue("2"),
				check.That(data.ResourceName).Key("optimized_auto_scale.0.maximum_instances").HasValue("3"),
			),
		},
		data.ImportStep(),
		{
			Config: r.optimizedAutoScaleUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("optimized_auto_scale.#").HasValue("1"),
				check.That(data.ResourceName).Key("optimized_auto_scale.0.minimum_instances").HasValue("3"),
				check.That(data.ResourceName).Key("optimized_auto_scale.0.maximum_instances").HasValue("4"),
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

func TestAccKustoCluster_updateSkuAndOptimizedAutoScale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.noOptimizedAutoScale(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.optimizedAutoScale(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.noOptimizedAutoScale(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoCluster_trustedExternalTenants(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.trustedExternalTenants(data, "[\"*\"]"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.unsetTrustedExternalTenants(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.trustedExternalTenants(data, "[data.azurerm_client_config.current.tenant_id]"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoCluster_newSkus(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.newSkus(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_L8s_v3"),
				check.That(data.ResourceName).Key("allowed_fqdns.#").HasValue("1"),
				check.That(data.ResourceName).Key("allowed_fqdns.0").HasValue("255.255.255.0/24"),
				check.That(data.ResourceName).Key("allowed_ip_ranges.#").HasValue("1"),
				check.That(data.ResourceName).Key("allowed_ip_ranges.0").HasValue("0.0.0.0/0"),
				check.That(data.ResourceName).Key("outbound_network_access_restricted").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func (KustoClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseKustoClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Kusto.ClustersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	if resp.Model == nil {
		return nil, fmt.Errorf("response model is empty")
	}

	exists := resp.Model.Properties != nil

	return &exists, nil
}

func (KustoClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                               = "acctestkc%s"
  location                           = azurerm_resource_group.test.location
  resource_group_name                = azurerm_resource_group.test.name
  allowed_fqdns                      = ["255.255.255.0/24"]
  allowed_ip_ranges                  = ["0.0.0.0/0"]
  public_network_access_enabled      = false
  public_ip_type                     = "DualStack"
  outbound_network_access_restricted = true
  sku {
    name     = "Standard_D13_v2"
    capacity = 2
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) trustedExternalTenants(data acceptance.TestData, tenantConfig string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  trusted_external_tenants = %s
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, tenantConfig)
}

func (KustoClusterResource) unsetTrustedExternalTenants(data acceptance.TestData) string {
	if !features.FourPointOhBeta() {
		return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  trusted_external_tenants = []
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) doubleEncryption(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                      = "acctestkc%s"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  double_encryption_enabled = true

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  tags = {
    label = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  tags = {
    label = "test1"
    ENV   = "prod"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) skuUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Standard_D11_v2"
    capacity = 2
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) withZones(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  zones = ["1"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                        = "acctestkc%s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  auto_stop_enabled           = true
  disk_encryption_enabled     = true
  streaming_ingestion_enabled = true
  purge_enabled               = true
  public_ip_type              = "DualStack"

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) identitySystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) userAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (KustoClusterResource) multipleAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (KustoClusterResource) languageExtensions(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Standard_E4d_v4"
    capacity = 2
  }

  language_extensions = ["PYTHON", "R"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) languageExtensionsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Standard_E4d_v4"
    capacity = 2
  }

  language_extensions = ["PYTHON_3.10.8", "R"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) languageExtensionsRemove(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Standard_E4d_v4"
    capacity = 2
  }

  language_extensions = ["R"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) noOptimizedAutoScale(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_E2a_v4"
    capacity = 1
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) optimizedAutoScale(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Standard_L8s_v3"
    capacity = 2
  }

  optimized_auto_scale {
    minimum_instances = 2
    maximum_instances = 3
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) optimizedAutoScaleUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name = "Standard_D11_v2"
  }

  optimized_auto_scale {
    minimum_instances = 3
    maximum_instances = 4
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) vnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestkc%s-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestkc%s-subnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "delegation"

    service_delegation {
      name    = "Microsoft.Kusto/clusters"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action", "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action", "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action"]
    }
  }
}

resource "azurerm_route_table" "test" {
  name                = "acctestkc%s-rt"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet_route_table_association" "test" {
  subnet_id      = azurerm_subnet.test.id
  route_table_id = azurerm_route_table.test.id
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestkc%s-nsg"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_rule" "test_allow_management_inbound" {
  name                        = "AllowAzureDataExplorerManagement"
  priority                    = 1000
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "443"
  source_address_prefix       = "AzureDataExplorerManagement"
  destination_address_prefix  = "VirtualNetwork"
  resource_group_name         = azurerm_resource_group.test.name
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_public_ip" "engine_pip" {
  name                = "acctestkc%s-engine-pip"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_public_ip" "management_pip" {
  name                = "acctestkc%s-management-pip"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  virtual_network_configuration {
    subnet_id                    = azurerm_subnet.test.id
    engine_public_ip_id          = azurerm_public_ip.engine_pip.id
    data_management_public_ip_id = azurerm_public_ip.management_pip.id
  }

  depends_on = [
    azurerm_subnet_route_table_association.test,
    azurerm_subnet_network_security_group_association.test,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString, data.RandomString, data.RandomString, data.RandomString, data.RandomString)
}

func (KustoClusterResource) vnetRemoved(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestkc%s-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestkc%s-subnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
  delegation {
    name = "delegation"
    service_delegation {
      name    = "Microsoft.Kusto/clusters"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action", "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action", "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action"]
    }
  }
}

resource "azurerm_route_table" "test" {
  name                = "acctestkc%s-rt"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet_route_table_association" "test" {
  subnet_id      = azurerm_subnet.test.id
  route_table_id = azurerm_route_table.test.id
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestkc%s-nsg"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_rule" "test_allow_management_inbound" {
  name                        = "AllowAzureDataExplorerManagement"
  priority                    = 1000
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "443"
  source_address_prefix       = "AzureDataExplorerManagement"
  destination_address_prefix  = "VirtualNetwork"
  resource_group_name         = azurerm_resource_group.test.name
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_public_ip" "engine_pip" {
  name                = "acctestkc%s-engine-pip"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_public_ip" "management_pip" {
  name                = "acctestkc%s-management-pip"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
  public_network_access_enabled = false
  depends_on = [
    azurerm_subnet_route_table_association.test,
    azurerm_subnet_network_security_group_association.test,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString, data.RandomString, data.RandomString, data.RandomString, data.RandomString)
}

func (KustoClusterResource) newSkus(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
resource "azurerm_kusto_cluster" "test" {
  name                               = "acctestkc%s"
  location                           = azurerm_resource_group.test.location
  resource_group_name                = azurerm_resource_group.test.name
  allowed_fqdns                      = ["255.255.255.0/24"]
  allowed_ip_ranges                  = ["0.0.0.0/0"]
  public_network_access_enabled      = false
  public_ip_type                     = "DualStack"
  outbound_network_access_restricted = true
  sku {
    name     = "Standard_L8s_v3"
    capacity = 2
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
