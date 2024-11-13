// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/managedinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MsSqlManagedInstanceResource struct{}

func TestAccMsSqlManagedInstance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")
	r := MsSqlManagedInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlManagedInstance_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")
	r := MsSqlManagedInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicZRS(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.basicZRS(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlManagedInstance_premium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")
	r := MsSqlManagedInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.premium(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlManagedInstance_backupRedundancyLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")
	r := MsSqlManagedInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageType(data, "LRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlManagedInstance_backupRedundancyGZRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")
	r := MsSqlManagedInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageType(data, "GZRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlManagedInstance_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")
	r := MsSqlManagedInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("0"),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.identity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
			),
		},
		data.ImportStep("administrator_login_password"),

		// TODO: uncomment this when https://github.com/Azure/azure-rest-api-specs/issues/16838 is resolved
		/*
			{
				Config: r.basic(data),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
					check.That(data.ResourceName).Key("identity.#").HasValue("0"),
				),
			},
			data.ImportStep("administrator_login_password"),
		*/
	})
}

func TestAccMsSqlManagedInstance_multiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")
	r := MsSqlManagedInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlManagedInstance_dnsZonePartner(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")
	r := MsSqlManagedInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dnsZonePartnerPrep(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.dnsZonePartner(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password", "dns_zone_partner_id"),
		{
			// DNS Zone Partner empty makes delete faster as MI can be destroyed simultaneously
			Config: r.emptyDnsZonePartner(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlManagedInstance_multipleDnsZonePartners(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")
	r := MsSqlManagedInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dnsZonePartnersPrep(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.dnsZonePartners(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password", "dns_zone_partner_id"),
		{
			// DNS Zone Partner empty makes delete faster as MI can be destroyed simultaneously
			Config: r.emptyDnsZonePartners(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlManagedInstance_withMaintenanceConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")
	r := MsSqlManagedInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withMaintenanceConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlManagedInstance_withServicePrincipal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")
	r := MsSqlManagedInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withServicePrincipal(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.withServicePrincipalUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlManagedInstance_backupRedundancyUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance", "test")
	r := MsSqlManagedInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageType(data, "GRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.storageType(data, "GZRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func (r MsSqlManagedInstanceResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseSqlManagedInstanceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQLManagedInstance.ManagedInstancesClient.Get(ctx, *id, managedinstances.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r MsSqlManagedInstanceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {
    resource_group {
      /* Due to the creation of unmanaged Microsoft.Network/networkIntentPolicies in this service,
      prevent_deletion_if_contains_resources has been added here to allow the test resources to be
       deleted until this can be properly investigated
      */
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_mssql_managed_instance" "test" {
  name                = "acctestsqlserver%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  license_type       = "BasePrice"
  sku_name           = "GP_Gen5"
  storage_size_in_gb = 32
  subnet_id          = azurerm_subnet.test.id
  vcores             = 4

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  depends_on = [
    azurerm_subnet_network_security_group_association.test,
    azurerm_subnet_route_table_association.test,
  ]

  tags = {
    environment = "staging"
    database    = "test"
  }
}
`, r.template(data, data.Locations.Primary), data.RandomInteger)
}

func (r MsSqlManagedInstanceResource) basicZRS(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {
    resource_group {
      /* Due to the creation of unmanaged Microsoft.Network/networkIntentPolicies in this service,
      prevent_deletion_if_contains_resources has been added here to allow the test resources to be
       deleted until this can be properly investigated
      */
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_mssql_managed_instance" "test" {
  name                = "acctestsqlserver%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  license_type         = "BasePrice"
  sku_name             = "GP_Gen5"
  storage_account_type = "ZRS"
  storage_size_in_gb   = 32
  subnet_id            = azurerm_subnet.test.id
  vcores               = 4

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  depends_on = [
    azurerm_subnet_network_security_group_association.test,
    azurerm_subnet_route_table_association.test,
  ]

  tags = {
    environment = "staging"
    database    = "test"
  }
}
`, r.template(data, data.Locations.Primary), data.RandomInteger)
}

func (r MsSqlManagedInstanceResource) premium(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {
    resource_group {
      /* Due to the creation of unmanaged Microsoft.Network/networkIntentPolicies in this service,
      prevent_deletion_if_contains_resources has been added here to allow the test resources to be
       deleted until this can be properly investigated
      */
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_mssql_managed_instance" "test" {
  name                = "acctestsqlserver%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  license_type       = "BasePrice"
  sku_name           = "GP_Gen8IM"
  storage_size_in_gb = 32
  subnet_id          = azurerm_subnet.test.id
  vcores             = 4

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  depends_on = [
    azurerm_subnet_network_security_group_association.test,
    azurerm_subnet_route_table_association.test,
  ]

  tags = {
    environment = "staging"
    database    = "test"
  }
}
`, r.template(data, data.Locations.Secondary), data.RandomInteger)
}

func (r MsSqlManagedInstanceResource) storageType(data acceptance.TestData, storageAccountType string) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {
    resource_group {
      /* Due to the creation of unmanaged Microsoft.Network/networkIntentPolicies in this service,
      prevent_deletion_if_contains_resources has been added here to allow the test resources to be
       deleted until this can be properly investigated
      */
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_mssql_managed_instance" "test" {
  name                = "acctestsqlserver%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  license_type         = "BasePrice"
  sku_name             = "GP_Gen5"
  storage_account_type = "%[3]s"
  storage_size_in_gb   = 32
  subnet_id            = azurerm_subnet.test.id
  vcores               = 4

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  depends_on = [
    azurerm_subnet_network_security_group_association.test,
    azurerm_subnet_route_table_association.test,
  ]

  tags = {
    environment = "staging"
    database    = "test"
  }
}
`, r.template(data, data.Locations.Primary), data.RandomInteger, storageAccountType)
}

func (r MsSqlManagedInstanceResource) identity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {
    resource_group {
      /* Due to the creation of unmanaged Microsoft.Network/networkIntentPolicies in this service,
      prevent_deletion_if_contains_resources has been added here to allow the test resources to be
       deleted until this can be properly investigated
      */
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_mssql_managed_instance" "test" {
  name                = "acctestsqlserver%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  license_type       = "BasePrice"
  sku_name           = "GP_Gen5"
  storage_size_in_gb = 32
  subnet_id          = azurerm_subnet.test.id
  vcores             = 4

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  depends_on = [
    azurerm_subnet_network_security_group_association.test,
    azurerm_subnet_route_table_association.test,
  ]

  identity {
    type = "SystemAssigned"
  }

  tags = {
    environment = "staging"
    database    = "test"
  }
}
`, r.template(data, data.Locations.Primary), data.RandomInteger)
}

func (r MsSqlManagedInstanceResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {
    resource_group {
      /* Due to the creation of unmanaged Microsoft.Network/networkIntentPolicies in this service,
      prevent_deletion_if_contains_resources has been added here to allow the test resources to be
       deleted until this can be properly investigated
      */
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_mssql_managed_instance" "test" {
  name                = "acctestsqlserver%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  license_type                 = "BasePrice"
  minimum_tls_version          = "1.0"
  proxy_override               = "Proxy"
  public_data_endpoint_enabled = true
  sku_name                     = "GP_Gen5"
  storage_account_type         = "ZRS"
  storage_size_in_gb           = 64
  subnet_id                    = azurerm_subnet.test.id
  timezone_id                  = "Pacific Standard Time"
  vcores                       = 8
  zone_redundant_enabled       = true

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  depends_on = [
    azurerm_subnet_network_security_group_association.test,
    azurerm_subnet_route_table_association.test,
  ]

  tags = {
    environment = "production"
  }
}
`, r.template(data, data.Locations.Primary), data.RandomInteger)
}

func (r MsSqlManagedInstanceResource) multiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {
    resource_group {
      /* Due to the creation of unmanaged Microsoft.Network/networkIntentPolicies in this service,
      prevent_deletion_if_contains_resources has been added here to allow the test resources to be
       deleted until this can be properly investigated
      */
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_mssql_managed_instance" "test" {
  name                = "acctestsqlserver%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  license_type       = "BasePrice"
  sku_name           = "GP_Gen5"
  storage_size_in_gb = 32
  subnet_id          = azurerm_subnet.test.id
  vcores             = 4

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  depends_on = [
    azurerm_subnet_network_security_group_association.test,
    azurerm_subnet_route_table_association.test,
  ]

  tags = {
    environment = "staging"
    database    = "test"
  }
}

resource "azurerm_mssql_managed_instance" "secondary" {
  name                = "acctestsqlserver2%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  license_type       = "BasePrice"
  sku_name           = "GP_Gen5"
  storage_size_in_gb = 32
  subnet_id          = azurerm_subnet.test.id
  vcores             = 4

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  depends_on = [
    azurerm_subnet_network_security_group_association.test,
    azurerm_subnet_route_table_association.test,
  ]

  tags = {
    environment = "prod"
    database    = "test"
  }
}
`, r.template(data, data.Locations.Primary), data.RandomInteger)
}

func (r MsSqlManagedInstanceResource) dnsZonePartnerPrep(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
%[2]s
`, r.basic(data), r.templateSecondary(data))
}

func (r MsSqlManagedInstanceResource) dnsZonePartner(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "azurerm_mssql_managed_instance" "secondary" {
  name                = "acctestsqlserver2%[3]d"
  resource_group_name = azurerm_resource_group.secondary.name
  location            = azurerm_resource_group.secondary.location

  dns_zone_partner_id = azurerm_mssql_managed_instance.test.id
  license_type        = "BasePrice"
  sku_name            = "GP_Gen5"
  storage_size_in_gb  = 32
  subnet_id           = azurerm_subnet.secondary.id
  vcores              = 4

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  depends_on = [
    azurerm_subnet_network_security_group_association.secondary,
    azurerm_subnet_route_table_association.secondary,
  ]

  tags = {
    environment = "prod"
    database    = "test"
  }
}
`, r.basic(data), r.templateSecondary(data), data.RandomInteger)
}

func (r MsSqlManagedInstanceResource) emptyDnsZonePartner(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "azurerm_mssql_managed_instance" "secondary" {
  name                = "acctestsqlserver2%[3]d"
  resource_group_name = azurerm_resource_group.secondary.name
  location            = azurerm_resource_group.secondary.location

  license_type       = "BasePrice"
  sku_name           = "GP_Gen5"
  storage_size_in_gb = 32
  subnet_id          = azurerm_subnet.secondary.id
  vcores             = 4

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  depends_on = [
    azurerm_subnet_network_security_group_association.secondary,
    azurerm_subnet_route_table_association.secondary,
  ]

  tags = {
    environment = "prod"
    database    = "test"
  }
}
`, r.basic(data), r.templateSecondary(data), data.RandomInteger)
}

func (r MsSqlManagedInstanceResource) dnsZonePartnersPrep(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
%[2]s
%[3]s
`, r.basic(data), r.templateSecondary(data), r.templateExtraSecondary(data))
}

func (r MsSqlManagedInstanceResource) dnsZonePartners(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
%[2]s
%[3]s

resource "azurerm_mssql_managed_instance" "secondary" {
  name                = "acctestsqlserver2%[4]d"
  resource_group_name = azurerm_resource_group.secondary.name
  location            = azurerm_resource_group.secondary.location

  dns_zone_partner_id = azurerm_mssql_managed_instance.test.id
  license_type        = "BasePrice"
  sku_name            = "GP_Gen5"
  storage_size_in_gb  = 32
  subnet_id           = azurerm_subnet.secondary.id
  vcores              = 4

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  depends_on = [
    azurerm_subnet_network_security_group_association.secondary,
    azurerm_subnet_route_table_association.secondary,
  ]

  tags = {
    environment = "prod"
    database    = "test"
  }
}

resource "azurerm_mssql_managed_instance" "secondary_2" {
  name                = "acctestsqlserver3%[4]d"
  resource_group_name = azurerm_resource_group.secondary_2.name
  location            = azurerm_resource_group.secondary_2.location

  dns_zone_partner_id = azurerm_mssql_managed_instance.test.id
  license_type        = "BasePrice"
  sku_name            = "GP_Gen5"
  storage_size_in_gb  = 32
  subnet_id           = azurerm_subnet.secondary_2.id
  vcores              = 4

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  depends_on = [
    azurerm_subnet_network_security_group_association.secondary_2,
    azurerm_subnet_route_table_association.secondary_2,
  ]

  tags = {
    environment = "prod"
    database    = "test"
  }
}
`, r.basic(data), r.templateSecondary(data), r.templateExtraSecondary(data), data.RandomInteger)
}

func (r MsSqlManagedInstanceResource) emptyDnsZonePartners(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
%[2]s
%[3]s

resource "azurerm_mssql_managed_instance" "secondary" {
  name                = "acctestsqlserver2%[4]d"
  resource_group_name = azurerm_resource_group.secondary.name
  location            = azurerm_resource_group.secondary.location

  license_type       = "BasePrice"
  sku_name           = "GP_Gen5"
  storage_size_in_gb = 32
  subnet_id          = azurerm_subnet.secondary.id
  vcores             = 4

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  depends_on = [
    azurerm_subnet_network_security_group_association.secondary,
    azurerm_subnet_route_table_association.secondary,
  ]

  tags = {
    environment = "prod"
    database    = "test"
  }
}

resource "azurerm_mssql_managed_instance" "secondary_2" {
  name                = "acctestsqlserver3%[4]d"
  resource_group_name = azurerm_resource_group.secondary_2.name
  location            = azurerm_resource_group.secondary_2.location

  license_type       = "BasePrice"
  sku_name           = "GP_Gen5"
  storage_size_in_gb = 32
  subnet_id          = azurerm_subnet.secondary_2.id
  vcores             = 4

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  depends_on = [
    azurerm_subnet_network_security_group_association.secondary_2,
    azurerm_subnet_route_table_association.secondary_2,
  ]

  tags = {
    environment = "prod"
    database    = "test"
  }
}
`, r.basic(data), r.templateSecondary(data), r.templateExtraSecondary(data), data.RandomInteger)
}

func (r MsSqlManagedInstanceResource) template(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG1-sql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet1-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
}

resource "azurerm_subnet" "test" {
  name                 = "subnet1-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.0.0/24"]

  delegation {
    name = "managedinstancedelegation"

    service_delegation {
      name    = "Microsoft.Sql/managedInstances"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action", "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action", "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action"]
    }
  }
}

resource "azurerm_network_security_group" "test" {
  name                = "mi-security-group1-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_rule" "allow_management_inbound_1" {
  name                        = "allow_management_inbound"
  priority                    = 106
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_ranges     = ["9000", "9003", "1438", "1440", "1452"]
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.test.name
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_network_security_rule" "allow_misubnet_inbound_1" {
  name                        = "allow_misubnet_inbound"
  priority                    = 200
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "10.0.0.0/24"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.test.name
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_network_security_rule" "allow_health_probe_inbound_1" {
  name                        = "allow_health_probe_inbound"
  priority                    = 300
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "AzureLoadBalancer"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.test.name
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_network_security_rule" "allow_tds_inbound_1" {
  name                        = "allow_tds_inbound"
  priority                    = 1000
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "1433"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.test.name
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_network_security_rule" "allow_redirect_inbound_1" {
  name                        = "allow_redirect_inbound"
  priority                    = 1100
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "11000-11999"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.test.name
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_network_security_rule" "allow_geodr_inbound_1" {
  name                        = "allow_geodr_inbound"
  priority                    = 1200
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "5022"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.test.name
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_network_security_rule" "deny_all_inbound_1" {
  name                        = "deny_all_inbound"
  priority                    = 4096
  direction                   = "Inbound"
  access                      = "Deny"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.test.name
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_network_security_rule" "allow_management_outbound_1" {
  name                        = "allow_management_outbound"
  priority                    = 110
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_ranges     = ["80", "443", "12000"]
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.test.name
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_network_security_rule" "allow_misubnet_outbound_1" {
  name                        = "allow_misubnet_outbound"
  priority                    = 200
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "10.0.0.0/24"
  resource_group_name         = azurerm_resource_group.test.name
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_network_security_rule" "allow_redirect_outbound_1" {
  name                        = "allow_redirect_outbound"
  priority                    = 1100
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "11000-11999"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.test.name
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_network_security_rule" "allow_geodr_outbound_1" {
  name                        = "allow_geodr_outbound"
  priority                    = 1200
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "5022"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.test.name
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_network_security_rule" "deny_all_outbound_1" {
  name                        = "deny_all_outbound"
  priority                    = 4096
  direction                   = "Outbound"
  access                      = "Deny"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.test.name
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_route_table" "test" {
  name                = "routetable1-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  route {
    name           = "subnet-to-vnetlocal"
    address_prefix = "10.0.0.0/24"
    next_hop_type  = "VnetLocal"
  }

  depends_on = [
    azurerm_subnet.test,
  ]

  lifecycle {
    ignore_changes = ["route"]
  }
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_subnet_route_table_association" "test" {
  subnet_id      = azurerm_subnet.test.id
  route_table_id = azurerm_route_table.test.id
}
  `, data.RandomInteger, location)
}

func (r MsSqlManagedInstanceResource) templateSecondary(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "secondary" {
  name     = "acctestRG2-sql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "secondary" {
  name                = "acctest-vnet2-%[1]d"
  resource_group_name = azurerm_resource_group.secondary.name
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.secondary.location
}

resource "azurerm_subnet" "secondary" {
  name                 = "subnet2-%[1]d"
  resource_group_name  = azurerm_resource_group.secondary.name
  virtual_network_name = azurerm_virtual_network.secondary.name
  address_prefixes     = ["10.1.0.0/24"]

  delegation {
    name = "managedinstancedelegation"

    service_delegation {
      name    = "Microsoft.Sql/managedInstances"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action", "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action", "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action"]
    }
  }
}

resource "azurerm_network_security_group" "secondary" {
  name                = "mi-security-group2-%[1]d"
  location            = azurerm_resource_group.secondary.location
  resource_group_name = azurerm_resource_group.secondary.name
}

resource "azurerm_network_security_rule" "allow_management_inbound_2" {
  name                        = "allow_management_inbound"
  priority                    = 106
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_ranges     = ["9000", "9003", "1438", "1440", "1452"]
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary.name
  network_security_group_name = azurerm_network_security_group.secondary.name
}

resource "azurerm_network_security_rule" "allow_misubnet_inbound_2" {
  name                        = "allow_misubnet_inbound"
  priority                    = 200
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "10.1.0.0/24"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary.name
  network_security_group_name = azurerm_network_security_group.secondary.name
}

resource "azurerm_network_security_rule" "allow_health_probe_inbound_2" {
  name                        = "allow_health_probe_inbound"
  priority                    = 300
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "AzureLoadBalancer"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary.name
  network_security_group_name = azurerm_network_security_group.secondary.name
}

resource "azurerm_network_security_rule" "allow_tds_inbound_2" {
  name                        = "allow_tds_inbound"
  priority                    = 1000
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "1433"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary.name
  network_security_group_name = azurerm_network_security_group.secondary.name
}

resource "azurerm_network_security_rule" "allow_redirect_inbound_2" {
  name                        = "allow_redirect_inbound"
  priority                    = 1100
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "11000-11999"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary.name
  network_security_group_name = azurerm_network_security_group.secondary.name
}

resource "azurerm_network_security_rule" "allow_geodr_inbound_2" {
  name                        = "allow_geodr_inbound"
  priority                    = 1200
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "5022"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary.name
  network_security_group_name = azurerm_network_security_group.secondary.name
}

resource "azurerm_network_security_rule" "deny_all_inbound_2" {
  name                        = "deny_all_inbound"
  priority                    = 4096
  direction                   = "Inbound"
  access                      = "Deny"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary.name
  network_security_group_name = azurerm_network_security_group.secondary.name
}

resource "azurerm_network_security_rule" "allow_management_outbound_2" {
  name                        = "allow_management_outbound"
  priority                    = 110
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_ranges     = ["80", "443", "12000"]
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary.name
  network_security_group_name = azurerm_network_security_group.secondary.name
}

resource "azurerm_network_security_rule" "allow_misubnet_outbound_2" {
  name                        = "allow_misubnet_outbound"
  priority                    = 200
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "10.1.0.0/24"
  resource_group_name         = azurerm_resource_group.secondary.name
  network_security_group_name = azurerm_network_security_group.secondary.name
}

resource "azurerm_network_security_rule" "allow_redirect_outbound_2" {
  name                        = "allow_redirect_outbound"
  priority                    = 1100
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "11000-11999"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary.name
  network_security_group_name = azurerm_network_security_group.secondary.name
}

resource "azurerm_network_security_rule" "allow_geodr_outbound_2" {
  name                        = "allow_geodr_outbound"
  priority                    = 1200
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "5022"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary.name
  network_security_group_name = azurerm_network_security_group.secondary.name
}

resource "azurerm_network_security_rule" "deny_all_outbound_2" {
  name                        = "deny_all_outbound"
  priority                    = 4096
  direction                   = "Outbound"
  access                      = "Deny"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary.name
  network_security_group_name = azurerm_network_security_group.secondary.name
}

resource "azurerm_route_table" "secondary" {
  name                = "routetable2-%[1]d"
  location            = azurerm_resource_group.secondary.location
  resource_group_name = azurerm_resource_group.secondary.name

  route {
    name           = "subnet-to-vnetlocal"
    address_prefix = "10.1.0.0/24"
    next_hop_type  = "VnetLocal"
  }

  depends_on = [
    azurerm_subnet.secondary,
  ]

  lifecycle {
    ignore_changes = ["route"]
  }
}

resource "azurerm_subnet_network_security_group_association" "secondary" {
  subnet_id                 = azurerm_subnet.secondary.id
  network_security_group_id = azurerm_network_security_group.secondary.id
}

resource "azurerm_subnet_route_table_association" "secondary" {
  subnet_id      = azurerm_subnet.secondary.id
  route_table_id = azurerm_route_table.secondary.id
}
  `, data.RandomInteger, data.Locations.Secondary)
}

func (r MsSqlManagedInstanceResource) templateExtraSecondary(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "secondary_2" {
  name     = "acctestRG3-sql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "secondary_2" {
  name                = "acctest-vnet3-%[1]d"
  resource_group_name = azurerm_resource_group.secondary_2.name
  address_space       = ["10.2.0.0/16"]
  location            = azurerm_resource_group.secondary_2.location
}

resource "azurerm_subnet" "secondary_2" {
  name                 = "subnet3-%[1]d"
  resource_group_name  = azurerm_resource_group.secondary_2.name
  virtual_network_name = azurerm_virtual_network.secondary_2.name
  address_prefixes     = ["10.2.0.0/24"]

  delegation {
    name = "managedinstancedelegation"

    service_delegation {
      name    = "Microsoft.Sql/managedInstances"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action", "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action", "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action"]
    }
  }
}

resource "azurerm_network_security_group" "secondary_2" {
  name                = "mi-security-group3-%[1]d"
  location            = azurerm_resource_group.secondary_2.location
  resource_group_name = azurerm_resource_group.secondary_2.name
}

resource "azurerm_network_security_rule" "allow_management_inbound_3" {
  name                        = "allow_management_inbound"
  priority                    = 106
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_ranges     = ["9000", "9003", "1438", "1440", "1452"]
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary_2.name
  network_security_group_name = azurerm_network_security_group.secondary_2.name
}

resource "azurerm_network_security_rule" "allow_misubnet_inbound_3" {
  name                        = "allow_misubnet_inbound"
  priority                    = 200
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "10.2.0.0/24"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary_2.name
  network_security_group_name = azurerm_network_security_group.secondary_2.name
}

resource "azurerm_network_security_rule" "allow_health_probe_inbound_3" {
  name                        = "allow_health_probe_inbound"
  priority                    = 300
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "AzureLoadBalancer"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary_2.name
  network_security_group_name = azurerm_network_security_group.secondary_2.name
}

resource "azurerm_network_security_rule" "allow_tds_inbound_3" {
  name                        = "allow_tds_inbound"
  priority                    = 1000
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "1433"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary_2.name
  network_security_group_name = azurerm_network_security_group.secondary_2.name
}

resource "azurerm_network_security_rule" "allow_redirect_inbound_3" {
  name                        = "allow_redirect_inbound"
  priority                    = 1100
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "11000-11999"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary_2.name
  network_security_group_name = azurerm_network_security_group.secondary_2.name
}

resource "azurerm_network_security_rule" "allow_geodr_inbound_3" {
  name                        = "allow_geodr_inbound"
  priority                    = 1200
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "5022"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary_2.name
  network_security_group_name = azurerm_network_security_group.secondary_2.name
}

resource "azurerm_network_security_rule" "deny_all_inbound_3" {
  name                        = "deny_all_inbound"
  priority                    = 4096
  direction                   = "Inbound"
  access                      = "Deny"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary_2.name
  network_security_group_name = azurerm_network_security_group.secondary_2.name
}

resource "azurerm_network_security_rule" "allow_management_outbound_3" {
  name                        = "allow_management_outbound"
  priority                    = 110
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_ranges     = ["80", "443", "12000"]
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary_2.name
  network_security_group_name = azurerm_network_security_group.secondary_2.name
}

resource "azurerm_network_security_rule" "allow_misubnet_outbound_3" {
  name                        = "allow_misubnet_outbound"
  priority                    = 200
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "10.2.0.0/24"
  resource_group_name         = azurerm_resource_group.secondary_2.name
  network_security_group_name = azurerm_network_security_group.secondary_2.name
}

resource "azurerm_network_security_rule" "allow_redirect_outbound_3" {
  name                        = "allow_redirect_outbound"
  priority                    = 1100
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "11000-11999"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary_2.name
  network_security_group_name = azurerm_network_security_group.secondary_2.name
}

resource "azurerm_network_security_rule" "allow_geodr_outbound_3" {
  name                        = "allow_geodr_outbound"
  priority                    = 1200
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "5022"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary_2.name
  network_security_group_name = azurerm_network_security_group.secondary_2.name
}

resource "azurerm_network_security_rule" "deny_all_outbound_3" {
  name                        = "deny_all_outbound"
  priority                    = 4096
  direction                   = "Outbound"
  access                      = "Deny"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.secondary_2.name
  network_security_group_name = azurerm_network_security_group.secondary_2.name
}

resource "azurerm_route_table" "secondary_2" {
  name                = "routetable3-%[1]d"
  location            = azurerm_resource_group.secondary_2.location
  resource_group_name = azurerm_resource_group.secondary_2.name

  route {
    name           = "subnet-to-vnetlocal"
    address_prefix = "10.2.0.0/24"
    next_hop_type  = "VnetLocal"
  }

  depends_on = [
    azurerm_subnet.secondary_2,
  ]

  lifecycle {
    ignore_changes = ["route"]
  }
}

resource "azurerm_subnet_network_security_group_association" "secondary_2" {
  subnet_id                 = azurerm_subnet.secondary_2.id
  network_security_group_id = azurerm_network_security_group.secondary_2.id
}

resource "azurerm_subnet_route_table_association" "secondary_2" {
  subnet_id      = azurerm_subnet.secondary_2.id
  route_table_id = azurerm_route_table.secondary_2.id
}
`, data.RandomInteger, data.Locations.Secondary)
}

func (r MsSqlManagedInstanceResource) withMaintenanceConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {
    resource_group {
      /* Due to the creation of unmanaged Microsoft.Network/networkIntentPolicies in this service, 
      prevent_deletion_if_contains_resources has been added here to allow the test resources to be 
       deleted until this can be properly investigated
      */
      prevent_deletion_if_contains_resources = false
    }
  }
}


resource "azurerm_mssql_managed_instance" "test" {
  name                = "acctestsqlserver%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  license_type       = "BasePrice"
  sku_name           = "GP_Gen5"
  storage_size_in_gb = 32
  subnet_id          = azurerm_subnet.test.id
  vcores             = 4

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  maintenance_configuration_name = "SQL_Default"

  depends_on = [
    azurerm_subnet_network_security_group_association.test,
    azurerm_subnet_route_table_association.test,
  ]

  tags = {
    environment = "staging"
    database    = "test"
  }
}
`, r.template(data, data.Locations.Primary), data.RandomInteger)
}

func (r MsSqlManagedInstanceResource) withServicePrincipal(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {
    resource_group {
      /* Due to the creation of unmanaged Microsoft.Network/networkIntentPolicies in this service,
      prevent_deletion_if_contains_resources has been added here to allow the test resources to be
       deleted until this can be properly investigated
      */
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_mssql_managed_instance" "test" {
  name                = "acctestsqlserver%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  license_type           = "BasePrice"
  service_principal_type = "SystemAssigned"
  sku_name               = "GP_Gen5"
  storage_size_in_gb     = 32
  subnet_id              = azurerm_subnet.test.id
  vcores                 = 4

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  depends_on = [
    azurerm_subnet_network_security_group_association.test,
    azurerm_subnet_route_table_association.test,
  ]

  tags = {
    environment = "staging"
    database    = "test"
  }
}
`, r.template(data, data.Locations.Primary), data.RandomInteger)
}

func (r MsSqlManagedInstanceResource) withServicePrincipalUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {
    resource_group {
      /* Due to the creation of unmanaged Microsoft.Network/networkIntentPolicies in this service,
      prevent_deletion_if_contains_resources has been added here to allow the test resources to be
       deleted until this can be properly investigated
      */
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_mssql_managed_instance" "test" {
  name                = "acctestsqlserver%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  license_type       = "BasePrice"
  sku_name           = "GP_Gen5"
  storage_size_in_gb = 32
  subnet_id          = azurerm_subnet.test.id
  vcores             = 4

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  depends_on = [
    azurerm_subnet_network_security_group_association.test,
    azurerm_subnet_route_table_association.test,
  ]

  tags = {
    environment = "staging"
    database    = "test"
  }
}
`, r.template(data, data.Locations.Primary), data.RandomInteger)
}
