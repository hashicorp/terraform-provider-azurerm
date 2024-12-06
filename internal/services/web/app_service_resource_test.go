// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AppServiceResource struct{}

func TestAccAppService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("outbound_ip_addresses").Exists(),
				check.That(data.ResourceName).Key("outbound_ip_address_list.#").Exists(),
				check.That(data.ResourceName).Key("possible_outbound_ip_addresses").Exists(),
				check.That(data.ResourceName).Key("possible_outbound_ip_address_list.#").Exists(),
				check.That(data.ResourceName).Key("custom_domain_verification_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

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

func TestAccAppService_movingAppService(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.moved(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccAppService_freeTier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.freeTier(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_sharedTier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sharedTier(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_32Bit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.service32Bit(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.use_32_bit_worker_process").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_backup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.backupUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Enabled
			Config: r.backup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			// Disabled
			Config: r.backupDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			// remove it
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_http2Enabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.http2Enabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.http2_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_alwaysOn(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.alwaysOn(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.always_on").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_appCommandLine(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appCommandLine(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.app_command_line").HasValue("/sbin/myserver -b 0.0.0.0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_httpsOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.httpsOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_only").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_clientCertEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.clientCertEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_cert_enabled").HasValue("true"),
			),
		},
		{
			Config: r.clientCertEnabledNotSet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_cert_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_clientCertEnabledWithMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.clientCertEnabledWithMode(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_cert_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("client_cert_mode").HasValue("Optional"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_appSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.foo").HasValue("bar"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_appSettingsVnetRouteAll(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettingsVnetRouteAll(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.WEBSITE_VNET_ROUTE_ALL").HasValue("true"),
				check.That(data.ResourceName).Key("site_config.0.vnet_route_all_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_siteConfigVnetRouteAll(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.siteConfigVnetRouteAll(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.vnet_route_all_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_clientAffinityEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.clientAffinityEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_affinity_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_clientAffinityDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.clientAffinityDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_affinity_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_enableManageServiceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.mangedServiceIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
			),
		},
	})
}

func TestAccAppService_updateResourceByEnablingManageServiceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("0"),
			),
		},
		{
			Config: r.mangedServiceIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
			),
		},
	})
}

func TestAccAppService_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsEmpty(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsEmpty(),
			),
		},
	})
}

func TestAccAppService_clientAffinityUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.clientAffinity(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_affinity_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.clientAffinity(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_affinity_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_connectionStrings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.connectionStrings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.connectionStringsUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_storageAccounts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccounts(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_account.#").HasValue("1"),
			),
		},
		{
			Config: r.storageAccountsUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_account.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_oneIpv4Restriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.oneIpv4Restriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.ip_address").HasValue("10.10.10.10/32"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.action").HasValue("Allow"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_oneIpv6Restriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.oneIpv6Restriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.ip_address").HasValue("2400:cb00::/32"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.action").HasValue("Allow"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_completeIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeIpRestriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.#").HasValue("1"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.ip_address").HasValue("10.10.10.10/32"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.name").HasValue("test-restriction"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.priority").HasValue("123"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.action").HasValue("Allow"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_forwarded_for.#").HasValue("2"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_forwarded_for.0").HasValue("2002::1234:abcd:ffff:c0a8:101/64"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_forwarded_for.1").HasValue("9.9.9.9/32"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_azure_fdid.#").HasValue("1"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_azure_fdid.0").HasValue("55ce4ed1-4b06-4bf1-b40e-4638452104da"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_fd_health_probe.#").HasValue("1"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_fd_health_probe.0").HasValue("1"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_forwarded_host.#").HasValue("1"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_forwarded_host.0").HasValue("example.com"),
			),
		},
		data.ImportStep(),
		{
			Config: r.manyCompleteIpRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.#").HasValue("4"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.ip_address").HasValue("10.10.10.10/32"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.name").HasValue("test-restriction"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.priority").HasValue("123"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.action").HasValue("Allow"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.1.ip_address").HasValue("20.20.20.0/24"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.1.name").HasValue("test-restriction-2"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.1.priority").HasValue("1234"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.1.action").HasValue("Deny"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.2.ip_address").HasValue("2400:cb00::/32"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.2.name").HasValue("test-restriction-3"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.2.priority").HasValue("65000"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.2.action").HasValue("Deny"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.3.service_tag").HasValue("AzureEventGrid"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.3.name").HasValue("test-restriction-4"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.3.priority").HasValue("65000"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.3.action").HasValue("Allow"),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeIpRestriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.#").HasValue("1"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.ip_address").HasValue("10.10.10.10/32"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.name").HasValue("test-restriction"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.priority").HasValue("123"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.action").HasValue("Allow"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_oneVNetSubnetIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.oneVNetSubnetIpRestriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_multipleVNetSubnetIpRestrictionUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.oneVNetSubnetIpRestriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.multipleVNetSubnetIpRestriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.#").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.oneVNetSubnetIpRestriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_mixedIpRestrictionUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.mixedIpRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.#").HasValue("4"),
			),
		},
		data.ImportStep(),
		{
			Config: r.mixedIpRestrictionsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.#").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.mixedIpRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.#").HasValue("4"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_zeroedIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// This configuration includes a single explicit ip_restriction
			Config: r.oneIpv4Restriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.#").HasValue("1"),
			),
		},
		{
			// This configuration has no site_config blocks at all.
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.#").HasValue("1"),
			),
		},
		{
			// This configuration explicitly sets ip_restriction to [] using attribute syntax.
			Config: r.zeroedIpRestriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.#").HasValue("0"),
			),
		},
	})
}

func TestAccAppService_zeroedIpRestrictionHeaders(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// This configuration includes a single explicit ip_restriction with headers
			Config: r.completeIpRestriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.#").HasValue("1"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.#").HasValue("1"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_forwarded_for.#").HasValue("2"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_forwarded_for.0").HasValue("2002::1234:abcd:ffff:c0a8:101/64"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_forwarded_for.1").HasValue("9.9.9.9/32"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_azure_fdid.#").HasValue("1"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_azure_fdid.0").HasValue("55ce4ed1-4b06-4bf1-b40e-4638452104da"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_fd_health_probe.#").HasValue("1"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_fd_health_probe.0").HasValue("1"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_forwarded_host.#").HasValue("1"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_forwarded_host.0").HasValue("example.com"),
			),
		},
		data.ImportStep(),
		{
			// This configuration includes a single explicit ip_restriction with one header
			Config: r.completeIpRestrictionOneHeader(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.#").HasValue("1"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.#").HasValue("1"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_forwarded_for.#").HasValue("1"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_forwarded_for.0").HasValue("9.9.9.9/32"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_forwarded_for.1").DoesNotExist(),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_azure_fdid.#").DoesNotExist(),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_fd_health_probe.#").DoesNotExist(),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_forwarded_host.#").DoesNotExist(),
			),
		},
		data.ImportStep(),
		{
			// This configuration has no site_config blocks at all.
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.#").HasValue("1"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.#").HasValue("1"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_forwarded_for.#").HasValue("1"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_forwarded_for.0").HasValue("9.9.9.9/32"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_forwarded_for.1").DoesNotExist(),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_azure_fdid.#").DoesNotExist(),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_fd_health_probe.#").DoesNotExist(),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.0.x_forwarded_host.#").DoesNotExist(),
			),
		},
		data.ImportStep(),
		{
			// This configuration explicitly sets ip_restriction.headers to [] using attribute syntax.
			Config: r.zeroedIpRestrictionHeaders(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.#").HasValue("1"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.headers.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_manyIpRestrictions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.manyIpRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.ip_address").HasValue("10.10.10.10/32"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.1.ip_address").HasValue("20.20.20.0/24"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.2.ip_address").HasValue("30.30.0.0/16"),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.3.ip_address").HasValue("192.168.1.2/24"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_scmUseMainIPRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.scmUseMainIPRestriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_scmOneIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.scmOneIpRestriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_completeScmIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeScmIpRestriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.manyCompleteScmIpRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeScmIpRestriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_oneVNetSubnetScmIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.oneVNetSubnetScmIpRestriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_zeroedScmIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// This configuration includes a single explicit scm_ip_restriction
			Config: r.scmOneIpRestriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.scm_ip_restriction.#").HasValue("1"),
			),
		},
		{
			// This configuration has no site_config blocks at all.
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.scm_ip_restriction.#").HasValue("1"),
			),
		},
		{
			// This configuration explicitly sets scm_ip_restriction to [] using attribute syntax.
			Config: r.zeroedScmIpRestriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.scm_ip_restriction.#").HasValue("0"),
			),
		},
	})
}

func TestAccAppService_manyScmIpRestrictions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.manyScmIpRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_defaultDocuments(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.defaultDocuments(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.default_documents.0").HasValue("first.html"),
				check.That(data.ResourceName).Key("site_config.0.default_documents.1").HasValue("second.jsp"),
				check.That(data.ResourceName).Key("site_config.0.default_documents.2").HasValue("third.aspx"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_enabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.enabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_localMySql(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.localMySql(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.local_mysql_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_applicationBlobStorageLogs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.applicationBlobStorageLogs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("logs.0.application_logs.0.file_system_level").HasValue("Warning"),
				check.That(data.ResourceName).Key("logs.0.application_logs.0.azure_blob_storage.0.level").HasValue("Information"),
				check.That(data.ResourceName).Key("logs.0.application_logs.0.azure_blob_storage.0.sas_url").HasValue("http://x.com/"),
				check.That(data.ResourceName).Key("logs.0.application_logs.0.azure_blob_storage.0.retention_in_days").HasValue("3"),
			),
		},
		{
			Config: r.applicationBlobStorageLogsWithAppSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("logs.0.application_logs.0.file_system_level").HasValue("Warning"),
				check.That(data.ResourceName).Key("logs.0.application_logs.0.azure_blob_storage.0.level").HasValue("Information"),
				check.That(data.ResourceName).Key("logs.0.application_logs.0.azure_blob_storage.0.sas_url").HasValue("http://x.com/"),
				check.That(data.ResourceName).Key("logs.0.application_logs.0.azure_blob_storage.0.retention_in_days").HasValue("3"),
				check.That(data.ResourceName).Key("app_settings.foo").HasValue("bar"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_httpFileSystemLogs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.httpFileSystemLogs(data),
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

func TestAccAppService_httpBlobStorageLogs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.httpBlobStorageLogs(data),
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

func TestAccAppService_httpFileSystemAndStorageBlobLogs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.httpFileSystemAndStorageBlobLogs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_managedPipelineMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managedPipelineMode(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.managed_pipeline_mode").HasValue("Classic"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_detailedErrorMessagesLogs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.detailedErrorMessages(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("logs.0.detailed_error_messages_enabled").HasValue("true"),
			),
		},
		{
			Config: r.detailedErrorMessages(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data.ResourceName, "logs.0.detailed_error_messages_enabled", "false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMAppService_failedRequestTracingLogs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.failedRequestTracing(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("logs.0.failed_request_tracing_enabled").HasValue("true"),
			),
		},
		{
			Config: r.failedRequestTracing(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("logs.0.failed_request_tracing_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_tagsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Hello").HasValue("World"),
			),
		},
		{
			Config: r.tagsUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.Hello").HasValue("World"),
				check.That(data.ResourceName).Key("tags.Terraform").HasValue("AcceptanceTests"),
			),
		},
	})
}

func TestAccAppService_remoteDebugging(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.remoteDebugging(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.remote_debugging_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("site_config.0.remote_debugging_version").HasValue("VS2022"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_windowsDotNet2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsDotNet(data, "v2.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.dotnet_framework_version").HasValue("v2.0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_windowsDotNet4(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsDotNet(data, "v4.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.dotnet_framework_version").HasValue("v4.0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_windowsDotNet5(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsDotNet(data, "v5.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.dotnet_framework_version").HasValue("v5.0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_windowsDotNet6(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsDotNet(data, "v6.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.dotnet_framework_version").HasValue("v6.0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_windowsDotNetUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsDotNet(data, "v2.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.dotnet_framework_version").HasValue("v2.0"),
			),
		},
		{
			Config: r.windowsDotNet(data, "v4.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.dotnet_framework_version").HasValue("v4.0"),
			),
		},
		{
			Config: r.windowsDotNet(data, "v5.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.dotnet_framework_version").HasValue("v5.0"),
			),
		},
		{
			Config: r.windowsDotNet(data, "v6.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.dotnet_framework_version").HasValue("v6.0"),
			),
		},
	})
}

func TestAccAppService_windowsJava7Java(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsJava(data, "1.7", "JAVA", "9.3"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.java_version").HasValue("1.7"),
				check.That(data.ResourceName).Key("site_config.0.java_container").HasValue("JAVA"),
				check.That(data.ResourceName).Key("site_config.0.java_container_version").HasValue("9.3"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_windowsJava8Java(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsJava(data, "1.8", "JAVA", "9.3"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.java_version").HasValue("1.8"),
				check.That(data.ResourceName).Key("site_config.0.java_container").HasValue("JAVA"),
				check.That(data.ResourceName).Key("site_config.0.java_container_version").HasValue("9.3"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_windowsJava11Java(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsJava(data, "11", "JAVA", "9.3"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.java_version").HasValue("11"),
				check.That(data.ResourceName).Key("site_config.0.java_container").HasValue("JAVA"),
				check.That(data.ResourceName).Key("site_config.0.java_container_version").HasValue("9.3"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_windowsJava7Jetty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsJava(data, "1.7", "JETTY", "9.3"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.java_version").HasValue("1.7"),
				check.That(data.ResourceName).Key("site_config.0.java_container").HasValue("JETTY"),
				check.That(data.ResourceName).Key("site_config.0.java_container_version").HasValue("9.3"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_windowsJava8Jetty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsJava(data, "1.8", "JETTY", "9.3"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.java_version").HasValue("1.8"),
				check.That(data.ResourceName).Key("site_config.0.java_container").HasValue("JETTY"),
				check.That(data.ResourceName).Key("site_config.0.java_container_version").HasValue("9.3"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_windowsJava11Jetty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsJava(data, "11", "JETTY", "9.3"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.java_version").HasValue("11"),
				check.That(data.ResourceName).Key("site_config.0.java_container").HasValue("JETTY"),
				check.That(data.ResourceName).Key("site_config.0.java_container_version").HasValue("9.3"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_windowsJava7Tomcat(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsJava(data, "1.7", "TOMCAT", "9.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.java_version").HasValue("1.7"),
				check.That(data.ResourceName).Key("site_config.0.java_container").HasValue("TOMCAT"),
				check.That(data.ResourceName).Key("site_config.0.java_container_version").HasValue("9.0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_windowsJava8Tomcat(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsJava(data, "1.8", "TOMCAT", "9.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.java_version").HasValue("1.8"),
				check.That(data.ResourceName).Key("site_config.0.java_container").HasValue("TOMCAT"),
				check.That(data.ResourceName).Key("site_config.0.java_container_version").HasValue("9.0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_windowsJava11Tomcat(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsJava(data, "11", "TOMCAT", "9.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.java_version").HasValue("11"),
				check.That(data.ResourceName).Key("site_config.0.java_container").HasValue("TOMCAT"),
				check.That(data.ResourceName).Key("site_config.0.java_container_version").HasValue("9.0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_windowsPHP7(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsPHP(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.php_version").HasValue("7.4"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_windowsPython(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsPython(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.python_version").HasValue("3.4"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_webSockets(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.webSockets(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.websockets_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_scmType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.scmType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.scm_type").HasValue("LocalGit"),
				check.That(data.ResourceName).Key("source_control.#").HasValue("1"),
				check.That(data.ResourceName).Key("site_credential.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_withSourceControl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withSourceControl(data, "main"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_withSourceControlUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withSourceControl(data, "main"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withSourceControl(data, "development"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_ftpsState(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ftpsState(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ftps_state").HasValue("AllAllowed"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_healthCheckPath(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.healthCheckPath(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.health_check_path").HasValue("/health"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_numberOfWorkers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.numberOfWorkers(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.number_of_workers").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

// Note: to specify `linux_fx_version` the App Service Plan must be of `kind = "Linux"`, and `reserved = true`
func TestAccAppService_linuxFxVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxFxVersion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_minTls(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.minTls(data, "1.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.min_tls_version").HasValue("1.0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.minTls(data, "1.1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.min_tls_version").HasValue("1.1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_corsSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.corsSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.cors.#").HasValue("1"),
				check.That(data.ResourceName).Key("site_config.0.cors.0.support_credentials").HasValue("true"),
				check.That(data.ResourceName).Key("site_config.0.cors.0.allowed_origins.#").HasValue("3"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_authSettingsAdditionalLoginParams(t *testing.T) {
	tenantID := os.Getenv("ARM_TENANT_ID")

	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authSettingsAdditionalLoginParams(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_settings.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("auth_settings.0.additional_login_params.test_key").HasValue("test_value"),
				check.That(data.ResourceName).Key("auth_settings.0.issuer").HasValue(fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.client_id").HasValue("aadclientid"),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.client_secret").HasValue("aadsecret"),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.allowed_audiences.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_authSettingsAdditionalAllowedExternalRedirectUrls(t *testing.T) {
	tenantID := os.Getenv("ARM_TENANT_ID")

	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authSettingsAdditionalAllowedExternalRedirectUrls(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_settings.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("auth_settings.0.allowed_external_redirect_urls.#").HasValue("1"),
				check.That(data.ResourceName).Key("auth_settings.0.allowed_external_redirect_urls.0").HasValue("https://terra.form"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.client_id").HasValue("aadclientid"),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.client_secret").HasValue("aadsecret"),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.allowed_audiences.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_authSettingsRuntimeVersion(t *testing.T) {
	tenantID := os.Getenv("ARM_TENANT_ID")

	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authSettingsRuntimeVersion(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_settings.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("auth_settings.0.runtime_version").HasValue("1.0"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.client_id").HasValue("aadclientid"),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.client_secret").HasValue("aadsecret"),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.allowed_audiences.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_authSettingsTokenRefreshExtensionHours(t *testing.T) {
	tenantID := os.Getenv("ARM_TENANT_ID")

	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authSettingsTokenRefreshExtensionHours(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_settings.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("auth_settings.0.token_refresh_extension_hours").HasValue("75"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.client_id").HasValue("aadclientid"),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.client_secret").HasValue("aadsecret"),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.allowed_audiences.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_authSettingsUnauthenticatedClientAction(t *testing.T) {
	tenantID := os.Getenv("ARM_TENANT_ID")

	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authSettingsUnauthenticatedClientAction(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_settings.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("auth_settings.0.unauthenticated_client_action").HasValue("RedirectToLoginPage"),
				check.That(data.ResourceName).Key("auth_settings.0.issuer").HasValue(fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.client_id").HasValue("aadclientid"),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.client_secret").HasValue("aadsecret"),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.allowed_audiences.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_authSettingsTokenStoreEnabled(t *testing.T) {
	tenantID := os.Getenv("ARM_TENANT_ID")

	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authSettingsTokenStoreEnabled(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_settings.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("auth_settings.0.token_store_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("auth_settings.0.issuer").HasValue(fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.client_id").HasValue("aadclientid"),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.client_secret").HasValue("aadsecret"),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.allowed_audiences.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_aadAuthSettings(t *testing.T) {
	tenantID := os.Getenv("ARM_TENANT_ID")

	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.aadAuthSettings(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_settings.0.enabled").HasValue("true"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.client_id").HasValue("aadclientid"),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.client_secret").HasValue("aadsecret"),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.allowed_audiences.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_facebookAuthSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.facebookAuthSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_settings.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("auth_settings.0.facebook.0.app_id").HasValue("facebookappid"),
				check.That(data.ResourceName).Key("auth_settings.0.facebook.0.app_secret").HasValue("facebookappsecret"),
				check.That(data.ResourceName).Key("auth_settings.0.facebook.0.oauth_scopes.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_googleAuthSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.googleAuthSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_settings.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("auth_settings.0.google.0.client_id").HasValue("googleclientid"),
				check.That(data.ResourceName).Key("auth_settings.0.google.0.client_secret").HasValue("googleclientsecret"),
				check.That(data.ResourceName).Key("auth_settings.0.google.0.oauth_scopes.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_microsoftAuthSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.microsoftAuthSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_settings.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("auth_settings.0.microsoft.0.client_id").HasValue("microsoftclientid"),
				check.That(data.ResourceName).Key("auth_settings.0.microsoft.0.client_secret").HasValue("microsoftclientsecret"),
				check.That(data.ResourceName).Key("auth_settings.0.microsoft.0.oauth_scopes.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_twitterAuthSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.twitterAuthSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_settings.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("auth_settings.0.twitter.0.consumer_key").HasValue("twitterconsumerkey"),
				check.That(data.ResourceName).Key("auth_settings.0.twitter.0.consumer_secret").HasValue("twitterconsumersecret"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_multiAuthSettings(t *testing.T) {
	tenantID := os.Getenv("ARM_TENANT_ID")

	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.aadAuthSettings(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_settings.0.enabled").HasValue("true"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.client_id").HasValue("aadclientid"),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.client_secret").HasValue("aadsecret"),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.allowed_audiences.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.aadMicrosoftAuthSettings(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_settings.0.enabled").HasValue("true"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.client_id").HasValue("aadclientid"),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.client_secret").HasValue("aadsecret"),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.allowed_audiences.#").HasValue("1"),
				check.That(data.ResourceName).Key("auth_settings.0.microsoft.0.client_id").HasValue("microsoftclientid"),
				check.That(data.ResourceName).Key("auth_settings.0.microsoft.0.client_secret").HasValue("microsoftclientsecret"),
				check.That(data.ResourceName).Key("auth_settings.0.microsoft.0.oauth_scopes.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_basicWindowsContainer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWindowsContainer(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.windows_fx_version").HasValue("DOCKER|mcr.microsoft.com/azure-app-service/samples/aspnethelloworld:latest"),
				check.That(data.ResourceName).Key("app_settings.DOCKER_REGISTRY_SERVER_URL").HasValue("https://mcr.microsoft.com"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_AcrManageIdentityCredentials(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.AcrUseManagedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_AcrUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.AcrUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppService_keyVaultUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	r := AppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.KeyVaultUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AppServiceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AppServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Web.AppServicesClient.Get(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving App Service %q (Resource Group %q): %+v", id.SiteName, id.ResourceGroup, err)
	}

	// The SDK defines 404 as an "ok" status code..
	if utils.ResponseWasNotFound(resp.Response) {
		return utils.Bool(false), nil
	}

	return utils.Bool(true), nil
}

func (r AppServiceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service" "import" {
  name                = azurerm_app_service.test.name
  location            = azurerm_app_service.test.location
  resource_group_name = azurerm_app_service.test.resource_group_name
  app_service_plan_id = azurerm_app_service.test.app_service_plan_id
}
`, template)
}

func (r AppServiceResource) freeTier(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Free"
    size = "F1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    use_32_bit_worker_process = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) moved(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service_plan" "other" {
  name                = "acctestASP2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.other.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) sharedTier(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Free"
    size = "F1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    use_32_bit_worker_process = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) alwaysOn(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    always_on = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) appCommandLine(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    app_command_line = "/sbin/myserver -b 0.0.0.0"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) httpsOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  https_only          = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) clientCertEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  client_cert_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) clientCertEnabledWithMode(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  client_cert_enabled = true
  client_cert_mode    = "Optional"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) clientCertEnabledNotSet(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) service32Bit(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    use_32_bit_worker_process = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) backup(data acceptance.TestData) string {
	template := r.backupTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  backup {
    name                = "acctest"
    storage_account_url = "https://${azurerm_storage_account.test.name}.blob.core.windows.net/${azurerm_storage_container.test.name}${data.azurerm_storage_account_sas.test.sas}&sr=b"
    schedule {
      frequency_interval = 1
      frequency_unit     = "Day"
    }
  }
}
`, template, data.RandomInteger)
}

func (r AppServiceResource) backupDisabled(data acceptance.TestData) string {
	template := r.backupTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  backup {
    name                = "acctest"
    storage_account_url = "https://${azurerm_storage_account.test.name}.blob.core.windows.net/${azurerm_storage_container.test.name}${data.azurerm_storage_account_sas.test.sas}&sr=b"
    enabled             = false
    schedule {
      frequency_interval = 1
      frequency_unit     = "Day"
    }
  }
}
`, template, data.RandomInteger)
}

func (r AppServiceResource) backupUpdated(data acceptance.TestData) string {
	template := r.backupTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  backup {
    name                = "acctest"
    storage_account_url = "https://${azurerm_storage_account.test.name}.blob.core.windows.net/${azurerm_storage_container.test.name}${data.azurerm_storage_account_sas.test.sas}&sr=b"
    schedule {
      frequency_interval = 2
      frequency_unit     = "Hour"
    }
  }
}
`, template, data.RandomInteger)
}

func (r AppServiceResource) backupTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "example"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

data "azurerm_storage_account_sas" "test" {
  connection_string = azurerm_storage_account.test.primary_connection_string
  https_only        = true

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

  start  = "2019-03-21"
  expiry = "2022-03-21"

  permissions {
    read    = false
    write   = true
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

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r AppServiceResource) http2Enabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    http2_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) appSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  app_settings = {
    "foo" = "bar"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) appSettingsVnetRouteAll(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  app_settings = {
    "WEBSITE_VNET_ROUTE_ALL" = "true"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) siteConfigVnetRouteAll(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    vnet_route_all_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) clientAffinityEnabled(data acceptance.TestData) string {
	return r.clientAffinity(data, true)
}

func (r AppServiceResource) clientAffinityDisabled(data acceptance.TestData) string {
	return r.clientAffinity(data, false)
}

func (r AppServiceResource) clientAffinity(data acceptance.TestData, clientAffinity bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                    = "acctestAS-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  app_service_plan_id     = azurerm_app_service_plan.test.id
  client_affinity_enabled = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, clientAffinity)
}

func (r AppServiceResource) mangedServiceIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) userAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) connectionStrings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  connection_string {
    name  = "First"
    value = "first-connection-string"
    type  = "Custom"
  }

  connection_string {
    name  = "Second"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) connectionStringsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  connection_string {
    name  = "Second"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  connection_string {
    name  = "First"
    value = "first-connection-string"
    type  = "Custom"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) storageAccounts(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acct%d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                 = "acctestcontainer"
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  storage_account {
    name         = "blobs"
    type         = "AzureBlob"
    account_name = azurerm_storage_account.test.name
    share_name   = azurerm_storage_container.test.name
    access_key   = azurerm_storage_account.test.primary_access_key
    mount_path   = "/blobs"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) storageAccountsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acct%d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                 = "acctestcontainer"
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_storage_share" "test" {
  name                 = "acctestshare"
  storage_account_name = azurerm_storage_account.test.name
  quota                = 1
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  storage_account {
    name         = "blobs"
    type         = "AzureBlob"
    account_name = azurerm_storage_account.test.name
    share_name   = azurerm_storage_container.test.name
    access_key   = azurerm_storage_account.test.primary_access_key
    mount_path   = "/blobs"
  }

  storage_account {
    name         = "files"
    type         = "AzureFiles"
    account_name = azurerm_storage_account.test.name
    share_name   = azurerm_storage_share.test.name
    access_key   = azurerm_storage_account.test.primary_access_key
    mount_path   = "/files"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) oneIpv4Restriction(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10/32"
      action     = "Allow"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) oneIpv6Restriction(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    ip_restriction {
      ip_address = "2400:cb00::/32"
      action     = "Allow"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) completeIpRestriction(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10/32"
      name       = "test-restriction"
      priority   = 123
      action     = "Allow"
      headers {
        x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da"]
        x_fd_health_probe = ["1"]
        x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64"]
        x_forwarded_host  = ["example.com"]
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) completeIpRestrictionOneHeader(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10/32"
      name       = "test-restriction"
      priority   = 123
      action     = "Allow"
      headers {
        x_forwarded_for = ["9.9.9.9/32"]
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) manyCompleteIpRestrictions(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10/32"
      name       = "test-restriction"
      priority   = 123
      action     = "Allow"
    }

    ip_restriction {
      ip_address = "20.20.20.0/24"
      name       = "test-restriction-2"
      priority   = 1234
      action     = "Deny"
    }

    ip_restriction {
      ip_address = "2400:cb00::/32"
      name       = "test-restriction-3"
      action     = "Deny"
    }

    ip_restriction {
      service_tag = "AzureEventGrid"
      name        = "test-restriction-4"
      action      = "Allow"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) oneVNetSubnetIpRestriction(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    ip_restriction {
      virtual_network_subnet_id = azurerm_subnet.test.id
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) multipleVNetSubnetIpRestriction(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_subnet" "test2" {
  name                 = "acctestsubnet%d-2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.3.0/24"]
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    ip_restriction {
      virtual_network_subnet_id = azurerm_subnet.test.id
    }
    ip_restriction {
      virtual_network_subnet_id = azurerm_subnet.test2.id
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) mixedIpRestrictions(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_subnet" "test2" {
  name                 = "acctestsubnet%d-2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.3.0/24"]
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    ip_restriction {
      virtual_network_subnet_id = azurerm_subnet.test.id
    }

    ip_restriction {
      virtual_network_subnet_id = azurerm_subnet.test2.id
    }

    ip_restriction {
      ip_address = "30.30.0.0/16"
    }

    ip_restriction {
      ip_address = "192.168.1.2/24"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) mixedIpRestrictionsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_subnet" "test2" {
  name                 = "acctestsubnet%d-2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.3.0/24"]
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    ip_restriction {
      virtual_network_subnet_id = azurerm_subnet.test.id
    }

    ip_restriction {
      ip_address = "30.30.0.0/16"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) manyIpRestrictions(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10/32"
    }

    ip_restriction {
      ip_address = "20.20.20.0/24"
    }

    ip_restriction {
      ip_address = "30.30.0.0/16"
    }

    ip_restriction {
      ip_address = "192.168.1.2/24"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) zeroedIpRestriction(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    ip_restriction = []
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) zeroedIpRestrictionHeaders(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10/32"
      name       = "test-restriction"
      priority   = 123
      action     = "Allow"
      headers    = []
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) scmUseMainIPRestriction(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    scm_use_main_ip_restriction = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) scmOneIpRestriction(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    scm_ip_restriction {
      ip_address = "10.10.10.10/32"
      action     = "Allow"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) completeScmIpRestriction(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    scm_ip_restriction {
      ip_address = "10.10.10.10/32"
      name       = "test-restriction"
      priority   = 123
      action     = "Allow"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) manyCompleteScmIpRestrictions(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    scm_ip_restriction {
      ip_address = "10.10.10.10/32"
      name       = "test-restriction"
      priority   = 123
      action     = "Allow"
    }

    scm_ip_restriction {
      ip_address = "20.20.20.0/24"
      name       = "test-restriction-2"
      priority   = 1234
      action     = "Deny"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) oneVNetSubnetScmIpRestriction(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    scm_ip_restriction {
      virtual_network_subnet_id = azurerm_subnet.test.id
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) zeroedScmIpRestriction(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    scm_ip_restriction = []
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) manyScmIpRestrictions(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    scm_ip_restriction {
      ip_address = "10.10.10.10/32"
    }

    scm_ip_restriction {
      ip_address = "20.20.20.0/24"
    }

    scm_ip_restriction {
      ip_address = "30.30.0.0/16"
    }

    scm_ip_restriction {
      ip_address = "192.168.1.2/24"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) defaultDocuments(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    default_documents = [
      "first.html",
      "second.jsp",
      "third.aspx",
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) enabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  enabled             = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) localMySql(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    local_mysql_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) applicationBlobStorageLogs(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  logs {
    application_logs {
      file_system_level = "Warning"
      azure_blob_storage {
        level             = "Information"
        sas_url           = "http://x.com/"
        retention_in_days = 3
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) applicationBlobStorageLogsWithAppSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_settings = {
    foo = "bar"
  }
  logs {
    application_logs {
      file_system_level = "Warning"
      azure_blob_storage {
        level             = "Information"
        sas_url           = "http://x.com/"
        retention_in_days = 3
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) httpFileSystemLogs(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  logs {
    http_logs {
      file_system {
        retention_in_days = 4
        retention_in_mb   = 25
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) httpBlobStorageLogs(data acceptance.TestData) string {
	template := r.backupTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  logs {
    http_logs {
      azure_blob_storage {
        sas_url           = "https://${azurerm_storage_account.test.name}.blob.core.windows.net/${azurerm_storage_container.test.name}${data.azurerm_storage_account_sas.test.sas}&sr=b"
        retention_in_days = 3
      }
    }
  }
}
`, template, data.RandomInteger)
}

func (r AppServiceResource) httpFileSystemAndStorageBlobLogs(data acceptance.TestData) string {
	template := r.backupTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  logs {
    application_logs {
      azure_blob_storage {
        level             = "Information"
        sas_url           = "https://${azurerm_storage_account.test.name}.blob.core.windows.net/${azurerm_storage_container.test.name}${data.azurerm_storage_account_sas.test.sas}&sr=b"
        retention_in_days = 3
      }
    }
    http_logs {
      file_system {
        retention_in_days = 4
        retention_in_mb   = 25
      }
    }
  }
}
`, template, data.RandomInteger)
}

func (r AppServiceResource) detailedErrorMessages(data acceptance.TestData, detailedErrorEnabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  logs {
    detailed_error_messages_enabled = %t
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, detailedErrorEnabled)
}

func (r AppServiceResource) failedRequestTracing(data acceptance.TestData, failedRequestEnabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  logs {
    failed_request_tracing_enabled = %t
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, failedRequestEnabled)
}

func (r AppServiceResource) managedPipelineMode(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    managed_pipeline_mode = "Classic"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) remoteDebugging(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    remote_debugging_enabled = true
    remote_debugging_version = "VS2022"
  }

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) tagsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  tags = {
    "Hello"     = "World"
    "Terraform" = "AcceptanceTests"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) windowsDotNet(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    dotnet_framework_version = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, version)
}

func (r AppServiceResource) windowsJava(data acceptance.TestData, javaVersion, container, containerVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    java_version           = "%s"
    java_container         = "%s"
    java_container_version = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, javaVersion, container, containerVersion)
}

func (r AppServiceResource) windowsPHP(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    php_version = "7.4"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) windowsPython(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    python_version = "3.4"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) webSockets(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    websockets_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) scmType(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    scm_type = "LocalGit"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) withSourceControl(data acceptance.TestData, branch string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-web-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  source_control {
    repo_url           = "https://github.com/jackofallops/azure-app-service-static-site-tests.git"
    branch             = "%[5]s"
    manual_integration = true
    rollback_enabled   = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, branch)
}

func (r AppServiceResource) ftpsState(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    ftps_state = "AllAllowed"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) healthCheckPath(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    health_check_path = "/health"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) numberOfWorkers(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  per_site_scaling    = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    number_of_workers = 1
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) linuxFxVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    always_on        = true
    linux_fx_version = "DOCKER|golang:latest"
  }

  app_settings = {
    "WEBSITES_ENABLE_APP_SERVICE_STORAGE" = "false"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) minTls(data acceptance.TestData, tlsVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    min_tls_version = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, tlsVersion)
}

func (r AppServiceResource) corsSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    cors {
      allowed_origins = [
        "http://www.contoso.com",
        "www.contoso.com",
        "contoso.com",
      ]

      support_credentials = true
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) authSettingsAdditionalLoginParams(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    always_on = true
  }

  auth_settings {
    enabled = true
    issuer  = "https://sts.windows.net/%s"

    additional_login_params = {
      test_key = "test_value"
    }

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, tenantID)
}

func (r AppServiceResource) authSettingsAdditionalAllowedExternalRedirectUrls(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    always_on = true
  }

  auth_settings {
    enabled = true
    issuer  = "https://sts.windows.net/%s"

    allowed_external_redirect_urls = [
      "https://terra.form",
    ]

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, tenantID)
}

func (r AppServiceResource) authSettingsRuntimeVersion(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    always_on = true
  }

  auth_settings {
    enabled         = true
    issuer          = "https://sts.windows.net/%s"
    runtime_version = "1.0"

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, tenantID)
}

func (r AppServiceResource) authSettingsTokenRefreshExtensionHours(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    always_on = true
  }

  auth_settings {
    enabled                       = true
    issuer                        = "https://sts.windows.net/%s"
    token_refresh_extension_hours = 75

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, tenantID)
}

func (r AppServiceResource) authSettingsTokenStoreEnabled(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    always_on = true
  }

  auth_settings {
    enabled             = true
    issuer              = "https://sts.windows.net/%s"
    token_store_enabled = true

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, tenantID)
}

func (r AppServiceResource) authSettingsUnauthenticatedClientAction(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    always_on = true
  }

  auth_settings {
    enabled                       = true
    issuer                        = "https://sts.windows.net/%s"
    unauthenticated_client_action = "RedirectToLoginPage"

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, tenantID)
}

func (r AppServiceResource) aadAuthSettings(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    always_on = true
  }

  auth_settings {
    enabled = true
    issuer  = "https://sts.windows.net/%s"

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, tenantID)
}

func (r AppServiceResource) facebookAuthSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    always_on = true
  }

  auth_settings {
    enabled = true

    facebook {
      app_id     = "facebookappid"
      app_secret = "facebookappsecret"

      oauth_scopes = [
        "facebookscope",
      ]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) googleAuthSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    always_on = true
  }

  auth_settings {
    enabled = true

    google {
      client_id     = "googleclientid"
      client_secret = "googleclientsecret"

      oauth_scopes = [
        "googlescope",
      ]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) microsoftAuthSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    always_on = true
  }

  auth_settings {
    enabled = true

    microsoft {
      client_id     = "microsoftclientid"
      client_secret = "microsoftclientsecret"

      oauth_scopes = [
        "microsoftscope",
      ]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) twitterAuthSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    always_on = true
  }

  auth_settings {
    enabled = true

    twitter {
      consumer_key    = "twitterconsumerkey"
      consumer_secret = "twitterconsumersecret"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) aadMicrosoftAuthSettings(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctestRG-%d"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestRG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    always_on = true
  }

  auth_settings {
    enabled          = true
    issuer           = "https://sts.windows.net/%s"
    default_provider = "%s"

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }

    microsoft {
      client_id     = "microsoftclientid"
      client_secret = "microsoftclientsecret"

      oauth_scopes = [
        "microsoftscope",
      ]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, tenantID, web.BuiltInAuthenticationProviderAzureActiveDirectory)
}

func (r AppServiceResource) basicWindowsContainer(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  is_xenon            = true
  kind                = "xenon"

  sku {
    tier = "PremiumV3"
    size = "P1v3"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    windows_fx_version = "DOCKER|mcr.microsoft.com/azure-app-service/samples/aspnethelloworld:latest"
  }

  app_settings = {
    "DOCKER_REGISTRY_SERVER_URL"      = "https://mcr.microsoft.com"
    "DOCKER_REGISTRY_SERVER_USERNAME" = ""
    "DOCKER_REGISTRY_SERVER_PASSWORD" = ""
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) AcrUseManagedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  identity {
    type = "SystemAssigned"
  }

  site_config {
    acr_use_managed_identity_credentials = true
  }

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) AcrUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    acr_use_managed_identity_credentials = true
    acr_user_managed_identity_client_id  = azurerm_user_assigned_identity.test.client_id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r AppServiceResource) KeyVaultUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  key_vault_reference_identity_id = azurerm_user_assigned_identity.test.id

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
