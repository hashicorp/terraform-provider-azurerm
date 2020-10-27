package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAppService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "outbound_ip_addresses"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "possible_outbound_ip_addresses"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMAppService_requiresImport),
		},
	})
}

func TestAccAzureRMAppService_movingAppService(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMAppService_moved(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_freeTier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_freeTier(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_sharedTier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_sharedTier(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_32Bit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_32Bit(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.use_32_bit_worker_process", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_backup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_backup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAppService_backupUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				// Enabled
				Config: testAccAzureRMAppService_backup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			{
				// Disabled
				Config: testAccAzureRMAppService_backupDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			{
				// remove it
				Config: testAccAzureRMAppService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_http2Enabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_http2Enabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.http2_enabled", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_alwaysOn(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_alwaysOn(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.always_on", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_appCommandLine(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_appCommandLine(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.app_command_line", "/sbin/myserver -b 0.0.0.0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_httpsOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_httpsOnly(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "https_only", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_clientCertEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_clientCertEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "client_cert_enabled", "true"),
				),
			},
			{
				Config: testAccAzureRMAppService_clientCertEnabledNotSet(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "client_cert_enabled", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_appSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_appSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "app_settings.foo", "bar"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_clientAffinityEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_clientAffinityEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "client_affinity_enabled", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_clientAffinityDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_clientAffinityDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "client_affinity_enabled", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_enableManageServiceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_mangedServiceIdentity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.tenant_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_updateResourceByEnablingManageServiceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.#", "0"),
				),
			},
			{
				Config: testAccAzureRMAppService_mangedServiceIdentity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.tenant_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_userAssignedIdentity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "UserAssigned"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.identity_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.principal_id", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.tenant_id", ""),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_clientAffinityUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_clientAffinity(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "client_affinity_enabled", "true"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAppService_clientAffinity(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "client_affinity_enabled", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_connectionStrings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_connectionStrings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.3173438943.name", "First"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.3173438943.value", "first-connection-string"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.3173438943.type", "Custom"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.2442860602.name", "Second"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.2442860602.value", "some-postgresql-connection-string"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.2442860602.type", "PostgreSQL"),
				),
			},
			{
				Config: testAccAzureRMAppService_connectionStringsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.3173438943.name", "First"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.3173438943.value", "first-connection-string"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.3173438943.type", "Custom"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.2442860602.name", "Second"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.2442860602.value", "some-postgresql-connection-string"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.2442860602.type", "PostgreSQL"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_storageAccounts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_storageAccounts(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_account.#", "1"),
				),
			},
			{
				Config: testAccAzureRMAppService_storageAccountsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_account.#", "2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_oneIpv4Restriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_oneIpv4Restriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.ip_address", "10.10.10.10/32"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.action", "Allow"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_oneIpv6Restriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_oneIpv6Restriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.ip_address", "2400:cb00::/32"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.action", "Allow"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_completeIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_completeIpRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.ip_address", "10.10.10.10/32"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.name", "test-restriction"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.priority", "123"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.action", "Allow"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAppService_manyCompleteIpRestrictions(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.#", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.ip_address", "10.10.10.10/32"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.name", "test-restriction"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.priority", "123"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.1.ip_address", "20.20.20.0/24"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.1.name", "test-restriction-2"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.1.priority", "1234"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.1.action", "Deny"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.2.ip_address", "2400:cb00::/32"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.2.name", "test-restriction-3"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.2.priority", "65000"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.2.action", "Deny"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAppService_completeIpRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.ip_address", "10.10.10.10/32"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.name", "test-restriction"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.priority", "123"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.action", "Allow"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_oneVNetSubnetIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_oneVNetSubnetIpRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_zeroedIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				// This configuration includes a single explicit ip_restriction
				Config: testAccAzureRMAppService_oneIpv4Restriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.#", "1"),
				),
			},
			{
				// This configuration has no site_config blocks at all.
				Config: testAccAzureRMAppService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.#", "1"),
				),
			},
			{
				// This configuration explicitly sets ip_restriction to [] using attribute syntax.
				Config: testAccAzureRMAppService_zeroedIpRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_manyIpRestrictions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_manyIpRestrictions(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.ip_address", "10.10.10.10/32"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.1.ip_address", "20.20.20.0/24"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.2.ip_address", "30.30.0.0/16"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.3.ip_address", "192.168.1.2/24"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_scmUseMainIPRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_scmUseMainIPRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_scmOneIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_scmOneIpRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_completeScmIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_completeScmIpRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAppService_manyCompleteScmIpRestrictions(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAppService_completeScmIpRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_oneVNetSubnetScmIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_oneVNetSubnetScmIpRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_zeroedScmIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				// This configuration includes a single explicit scm_ip_restriction
				Config: testAccAzureRMAppService_scmOneIpRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.scm_ip_restriction.#", "1"),
				),
			},
			{
				// This configuration has no site_config blocks at all.
				Config: testAccAzureRMAppService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.scm_ip_restriction.#", "1"),
				),
			},
			{
				// This configuration explicitly sets scm_ip_restriction to [] using attribute syntax.
				Config: testAccAzureRMAppService_zeroedScmIpRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.scm_ip_restriction.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_manyScmIpRestrictions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_manyScmIpRestrictions(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_defaultDocuments(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_defaultDocuments(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.default_documents.0", "first.html"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.default_documents.1", "second.jsp"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.default_documents.2", "third.aspx"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_enabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_enabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_localMySql(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_localMySql(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.local_mysql_enabled", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_applicationBlobStorageLogs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_applicationBlobStorageLogs(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "logs.0.application_logs.0.file_system_level", "Warning"),
					resource.TestCheckResourceAttr(data.ResourceName, "logs.0.application_logs.0.azure_blob_storage.0.level", "Information"),
					resource.TestCheckResourceAttr(data.ResourceName, "logs.0.application_logs.0.azure_blob_storage.0.sas_url", "http://x.com/"),
					resource.TestCheckResourceAttr(data.ResourceName, "logs.0.application_logs.0.azure_blob_storage.0.retention_in_days", "3"),
				),
			},
			{
				Config: testAccAzureRMAppService_applicationBlobStorageLogsWithAppSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "logs.0.application_logs.0.file_system_level", "Warning"),
					resource.TestCheckResourceAttr(data.ResourceName, "logs.0.application_logs.0.azure_blob_storage.0.level", "Information"),
					resource.TestCheckResourceAttr(data.ResourceName, "logs.0.application_logs.0.azure_blob_storage.0.sas_url", "http://x.com/"),
					resource.TestCheckResourceAttr(data.ResourceName, "logs.0.application_logs.0.azure_blob_storage.0.retention_in_days", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "app_settings.foo", "bar"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_httpFileSystemLogs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_httpFileSystemLogs(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAppService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_httpBlobStorageLogs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_httpBlobStorageLogs(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAppService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_httpFileSystemAndStorageBlobLogs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_httpFileSystemAndStorageBlobLogs(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_managedPipelineMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_managedPipelineMode(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.managed_pipeline_mode", "Classic"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_tagsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_tags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Hello", "World"),
				),
			},
			{
				Config: testAccAzureRMAppService_tagsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Hello", "World"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Terraform", "AcceptanceTests"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_remoteDebugging(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_remoteDebugging(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.remote_debugging_enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.remote_debugging_version", "VS2019"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_windowsDotNet2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_windowsDotNet(data, "v2.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.dotnet_framework_version", "v2.0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_windowsDotNet4(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_windowsDotNet(data, "v4.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.dotnet_framework_version", "v4.0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_windowsDotNetUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_windowsDotNet(data, "v2.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.dotnet_framework_version", "v2.0"),
				),
			},
			{
				Config: testAccAzureRMAppService_windowsDotNet(data, "v4.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.dotnet_framework_version", "v4.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_windowsJava7Java(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_windowsJava(data, "1.7", "JAVA", "9.3"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_version", "1.7"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container", "JAVA"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container_version", "9.3"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_windowsJava8Java(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_windowsJava(data, "1.8", "JAVA", "9.3"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_version", "1.8"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container", "JAVA"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container_version", "9.3"),
				),
			},
			data.ImportStep(),
		},
	})
}
func TestAccAzureRMAppService_windowsJava11Java(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_windowsJava(data, "11", "JAVA", "9.3"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_version", "11"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container", "JAVA"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container_version", "9.3"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_windowsJava7Jetty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_windowsJava(data, "1.7", "JETTY", "9.3"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_version", "1.7"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container", "JETTY"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container_version", "9.3"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_windowsJava8Jetty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_windowsJava(data, "1.8", "JETTY", "9.3"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_version", "1.8"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container", "JETTY"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container_version", "9.3"),
				),
			},
			data.ImportStep(),
		},
	})
}
func TestAccAzureRMAppService_windowsJava11Jetty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_windowsJava(data, "11", "JETTY", "9.3"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_version", "11"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container", "JETTY"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container_version", "9.3"),
				),
			},
			data.ImportStep(),
		},
	})
}
func TestAccAzureRMAppService_windowsJava7Tomcat(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_windowsJava(data, "1.7", "TOMCAT", "9.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_version", "1.7"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container", "TOMCAT"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container_version", "9.0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_windowsJava8Tomcat(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_windowsJava(data, "1.8", "TOMCAT", "9.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_version", "1.8"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container", "TOMCAT"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container_version", "9.0"),
				),
			},
			data.ImportStep(),
		},
	})
}
func TestAccAzureRMAppService_windowsJava11Tomcat(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_windowsJava(data, "11", "TOMCAT", "9.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_version", "11"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container", "TOMCAT"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container_version", "9.0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_windowsJava7Minor(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_windowsJava(data, "1.7.0_80", "TOMCAT", "9.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_version", "1.7.0_80"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container", "TOMCAT"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container_version", "9.0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_windowsJava8Minor(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_windowsJava(data, "1.8.0_181", "TOMCAT", "9.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_version", "1.8.0_181"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container", "TOMCAT"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container_version", "9.0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_windowsPHP7(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_windowsPHP(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.php_version", "7.3"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_windowsPython(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_windowsPython(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.python_version", "3.4"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_webSockets(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_webSockets(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.websockets_enabled", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_scmType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_scmType(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.scm_type", "LocalGit"),
					resource.TestCheckResourceAttr(data.ResourceName, "source_control.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_credential.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_withSourceControl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_withSourceControl(data, "main"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_withSourceControlUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_withSourceControl(data, "main"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAppService_withSourceControl(data, "development"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_ftpsState(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_ftpsState(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ftps_state", "AllAllowed"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_healthCheckPath(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_healthCheckPath(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.health_check_path", "/health"),
				),
			},
			data.ImportStep(),
		},
	})
}

// Note: to specify `linux_fx_version` the App Service Plan must be of `kind = "Linux"`, and `reserved = true`
func TestAccAzureRMAppService_linuxFxVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_linuxFxVersion(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_minTls(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_minTls(data, "1.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.min_tls_version", "1.0"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAppService_minTls(data, "1.1"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.min_tls_version", "1.1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_corsSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_corsSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.cors.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.cors.0.support_credentials", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.cors.0.allowed_origins.#", "3"),
				)},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_authSettingsAdditionalLoginParams(t *testing.T) {
	tenantID := os.Getenv("ARM_TENANT_ID")

	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_authSettingsAdditionalLoginParams(data, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.additional_login_params.test_key", "test_value"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_authSettingsAdditionalAllowedExternalRedirectUrls(t *testing.T) {
	tenantID := os.Getenv("ARM_TENANT_ID")

	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_authSettingsAdditionalAllowedExternalRedirectUrls(data, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.allowed_external_redirect_urls.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.allowed_external_redirect_urls.0", "https://terra.form"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_authSettingsRuntimeVersion(t *testing.T) {
	tenantID := os.Getenv("ARM_TENANT_ID")

	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_authSettingsRuntimeVersion(data, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.runtime_version", "1.0"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_authSettingsTokenRefreshExtensionHours(t *testing.T) {
	tenantID := os.Getenv("ARM_TENANT_ID")

	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_authSettingsTokenRefreshExtensionHours(data, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.token_refresh_extension_hours", "75"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_authSettingsUnauthenticatedClientAction(t *testing.T) {
	tenantID := os.Getenv("ARM_TENANT_ID")

	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_authSettingsUnauthenticatedClientAction(data, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.unauthenticated_client_action", "RedirectToLoginPage"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_authSettingsTokenStoreEnabled(t *testing.T) {
	tenantID := os.Getenv("ARM_TENANT_ID")

	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_authSettingsTokenStoreEnabled(data, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.token_store_enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_aadAuthSettings(t *testing.T) {
	tenantID := os.Getenv("ARM_TENANT_ID")

	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_aadAuthSettings(data, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_facebookAuthSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_facebookAuthSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.facebook.0.app_id", "facebookappid"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.facebook.0.app_secret", "facebookappsecret"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.facebook.0.oauth_scopes.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_googleAuthSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_googleAuthSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.google.0.client_id", "googleclientid"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.google.0.client_secret", "googleclientsecret"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.google.0.oauth_scopes.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_microsoftAuthSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_microsoftAuthSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.microsoft.0.client_id", "microsoftclientid"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.microsoft.0.client_secret", "microsoftclientsecret"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.microsoft.0.oauth_scopes.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_twitterAuthSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_twitterAuthSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.twitter.0.consumer_key", "twitterconsumerkey"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.twitter.0.consumer_secret", "twitterconsumersecret"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_multiAuthSettings(t *testing.T) {
	tenantID := os.Getenv("ARM_TENANT_ID")

	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_aadAuthSettings(data, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAppService_aadMicrosoftAuthSettings(data, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.microsoft.0.client_id", "microsoftclientid"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.microsoft.0.client_secret", "microsoftclientsecret"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.microsoft.0.oauth_scopes.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_basicWindowsContainer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_basicWindowsContainer(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.windows_fx_version", "DOCKER|mcr.microsoft.com/azure-app-service/samples/aspnethelloworld:latest"),
					resource.TestCheckResourceAttr(data.ResourceName, "app_settings.DOCKER_REGISTRY_SERVER_URL", "https://mcr.microsoft.com"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppService_aseScopeNameCheck(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppService_inAppServiceEnvironment(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMAppServiceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return nil
	}

	return nil
}

func testCheckAzureRMAppServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		appServiceName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Service: %s", appServiceName)
		}

		resp, err := client.Get(ctx, resourceGroup, appServiceName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: App Service %q (resource group: %q) does not exist", appServiceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on appServicesClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMAppService_basic(data acceptance.TestData) string {
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

func testAccAzureRMAppService_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMAppService_basic(data)
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

func testAccAzureRMAppService_freeTier(data acceptance.TestData) string {
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

func testAccAzureRMAppService_moved(data acceptance.TestData) string {
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

func testAccAzureRMAppService_sharedTier(data acceptance.TestData) string {
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

func testAccAzureRMAppService_alwaysOn(data acceptance.TestData) string {
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

func testAccAzureRMAppService_appCommandLine(data acceptance.TestData) string {
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

func testAccAzureRMAppService_httpsOnly(data acceptance.TestData) string {
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

func testAccAzureRMAppService_clientCertEnabled(data acceptance.TestData) string {
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

func testAccAzureRMAppService_clientCertEnabledNotSet(data acceptance.TestData) string {
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

func testAccAzureRMAppService_32Bit(data acceptance.TestData) string {
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

func testAccAzureRMAppService_backup(data acceptance.TestData) string {
	template := testAccAzureRMAppService_backupTemplate(data)
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

func testAccAzureRMAppService_backupDisabled(data acceptance.TestData) string {
	template := testAccAzureRMAppService_backupTemplate(data)
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

func testAccAzureRMAppService_backupUpdated(data acceptance.TestData) string {
	template := testAccAzureRMAppService_backupTemplate(data)
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

func testAccAzureRMAppService_backupTemplate(data acceptance.TestData) string {
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

func testAccAzureRMAppService_http2Enabled(data acceptance.TestData) string {
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

func testAccAzureRMAppService_appSettings(data acceptance.TestData) string {
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

func testAccAzureRMAppService_clientAffinityEnabled(data acceptance.TestData) string {
	return testAccAzureRMAppService_clientAffinity(data, true)
}

func testAccAzureRMAppService_clientAffinityDisabled(data acceptance.TestData) string {
	return testAccAzureRMAppService_clientAffinity(data, false)
}

func testAccAzureRMAppService_clientAffinity(data acceptance.TestData, clientAffinity bool) string {
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

func testAccAzureRMAppService_mangedServiceIdentity(data acceptance.TestData) string {
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

func testAccAzureRMAppService_userAssignedIdentity(data acceptance.TestData) string {
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

func testAccAzureRMAppService_connectionStrings(data acceptance.TestData) string {
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

func testAccAzureRMAppService_connectionStringsUpdated(data acceptance.TestData) string {
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

func testAccAzureRMAppService_storageAccounts(data acceptance.TestData) string {
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

func testAccAzureRMAppService_storageAccountsUpdated(data acceptance.TestData) string {
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

func testAccAzureRMAppService_oneIpv4Restriction(data acceptance.TestData) string {
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

func testAccAzureRMAppService_oneIpv6Restriction(data acceptance.TestData) string {
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

func testAccAzureRMAppService_completeIpRestriction(data acceptance.TestData) string {
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
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppService_manyCompleteIpRestrictions(data acceptance.TestData) string {
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
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppService_oneVNetSubnetIpRestriction(data acceptance.TestData) string {
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
  address_prefix       = "10.0.2.0/24"
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

func testAccAzureRMAppService_manyIpRestrictions(data acceptance.TestData) string {
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

func testAccAzureRMAppService_zeroedIpRestriction(data acceptance.TestData) string {
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

func testAccAzureRMAppService_scmUseMainIPRestriction(data acceptance.TestData) string {
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

func testAccAzureRMAppService_scmOneIpRestriction(data acceptance.TestData) string {
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

func testAccAzureRMAppService_completeScmIpRestriction(data acceptance.TestData) string {
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

func testAccAzureRMAppService_manyCompleteScmIpRestrictions(data acceptance.TestData) string {
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

func testAccAzureRMAppService_oneVNetSubnetScmIpRestriction(data acceptance.TestData) string {
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
  address_prefix       = "10.0.2.0/24"
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

func testAccAzureRMAppService_zeroedScmIpRestriction(data acceptance.TestData) string {
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

func testAccAzureRMAppService_manyScmIpRestrictions(data acceptance.TestData) string {
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

func testAccAzureRMAppService_defaultDocuments(data acceptance.TestData) string {
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

func testAccAzureRMAppService_enabled(data acceptance.TestData) string {
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

func testAccAzureRMAppService_localMySql(data acceptance.TestData) string {
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

func testAccAzureRMAppService_applicationBlobStorageLogs(data acceptance.TestData) string {
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

func testAccAzureRMAppService_applicationBlobStorageLogsWithAppSettings(data acceptance.TestData) string {
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

func testAccAzureRMAppService_httpFileSystemLogs(data acceptance.TestData) string {
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

func testAccAzureRMAppService_httpBlobStorageLogs(data acceptance.TestData) string {
	template := testAccAzureRMAppService_backupTemplate(data)
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

func testAccAzureRMAppService_httpFileSystemAndStorageBlobLogs(data acceptance.TestData) string {
	template := testAccAzureRMAppService_backupTemplate(data)
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

func testAccAzureRMAppService_managedPipelineMode(data acceptance.TestData) string {
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

func testAccAzureRMAppService_remoteDebugging(data acceptance.TestData) string {
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
    remote_debugging_version = "VS2019"
  }

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppService_tags(data acceptance.TestData) string {
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

func testAccAzureRMAppService_tagsUpdated(data acceptance.TestData) string {
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

func testAccAzureRMAppService_windowsDotNet(data acceptance.TestData, version string) string {
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

func testAccAzureRMAppService_windowsJava(data acceptance.TestData, javaVersion, container, containerVersion string) string {
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

func testAccAzureRMAppService_windowsPHP(data acceptance.TestData) string {
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
    php_version = "7.3"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppService_windowsPython(data acceptance.TestData) string {
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

func testAccAzureRMAppService_webSockets(data acceptance.TestData) string {
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

func testAccAzureRMAppService_scmType(data acceptance.TestData) string {
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

func testAccAzureRMAppService_withSourceControl(data acceptance.TestData, branch string) string {
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

func testAccAzureRMAppService_ftpsState(data acceptance.TestData) string {
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

func testAccAzureRMAppService_healthCheckPath(data acceptance.TestData) string {
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

func testAccAzureRMAppService_linuxFxVersion(data acceptance.TestData) string {
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

func testAccAzureRMAppService_minTls(data acceptance.TestData, tlsVersion string) string {
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

func testAccAzureRMAppService_corsSettings(data acceptance.TestData) string {
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

func testAccAzureRMAppService_authSettingsAdditionalLoginParams(data acceptance.TestData, tenantID string) string {
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

func testAccAzureRMAppService_authSettingsAdditionalAllowedExternalRedirectUrls(data acceptance.TestData, tenantID string) string {
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

func testAccAzureRMAppService_authSettingsRuntimeVersion(data acceptance.TestData, tenantID string) string {
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

func testAccAzureRMAppService_authSettingsTokenRefreshExtensionHours(data acceptance.TestData, tenantID string) string {
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

func testAccAzureRMAppService_authSettingsTokenStoreEnabled(data acceptance.TestData, tenantID string) string {
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

func testAccAzureRMAppService_authSettingsUnauthenticatedClientAction(data acceptance.TestData, tenantID string) string {
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

func testAccAzureRMAppService_aadAuthSettings(data acceptance.TestData, tenantID string) string {
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

func testAccAzureRMAppService_facebookAuthSettings(data acceptance.TestData) string {
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

func testAccAzureRMAppService_googleAuthSettings(data acceptance.TestData) string {
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

func testAccAzureRMAppService_microsoftAuthSettings(data acceptance.TestData) string {
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

func testAccAzureRMAppService_twitterAuthSettings(data acceptance.TestData) string {
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

func testAccAzureRMAppService_aadMicrosoftAuthSettings(data acceptance.TestData, tenantID string) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, tenantID, web.AzureActiveDirectory)
}

func testAccAzureRMAppService_basicWindowsContainer(data acceptance.TestData) string {
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
    tier = "PremiumContainer"
    size = "PC2"
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

func testAccAzureRMAppService_inAppServiceEnvironment(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "ase" {
  name                 = "asesubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_subnet" "gateway" {
  name                 = "gatewaysubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_app_service_environment" "test" {
  name      = "acctest-ase-%d"
  subnet_id = azurerm_subnet.ase.id
}

resource "azurerm_app_service_plan" "test" {
  name                       = "acctest-ASP-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_environment_id = azurerm_app_service_environment.test.id

  sku {
    tier     = "Isolated"
    size     = "I1"
    capacity = 1
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}

`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
