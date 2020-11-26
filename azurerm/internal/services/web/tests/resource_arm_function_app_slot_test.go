package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMFunctionAppSlot_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionAppSlot_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMFunctionAppSlot_requiresImport),
		},
	})
}

func TestAccAzureRMFunctionAppSlot_32Bit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_32Bit(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.use_32_bit_worker_process", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionAppSlot_alwaysOn(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_alwaysOn(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.always_on", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionAppSlot_appSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_appSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "app_settings.foo", "bar"),
				),
			},

			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionAppSlot_clientAffinityEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_clientAffinityEnabled(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "client_affinity_enabled", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionAppSlot_clientAffinityEnabledUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_clientAffinityEnabled(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "client_affinity_enabled", "true"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMFunctionAppSlot_clientAffinityEnabled(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "client_affinity_enabled", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionAppSlot_connectionStrings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_connectionStrings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.3173438943.name", "First"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.3173438943.value", "first-connection-string"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.3173438943.type", "Custom"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.2442860602.name", "Second"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.2442860602.value", "some-postgresql-connection-string"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.2442860602.type", "PostgreSQL"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMFunctionAppSlot_connectionStringsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.3173438943.name", "First"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.3173438943.value", "first-connection-string"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.3173438943.type", "Custom"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.2442860602.name", "Second"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.2442860602.value", "some-postgresql-connection-string"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.2442860602.type", "PostgreSQL"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionAppSlot_corsSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_corsSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionAppSlot_autoSwap(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_autoSwap(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.auto_swap_slot_name", "production"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionAppSlot_authSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	tenantID := os.Getenv("ARM_TENANT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_authSettings(data, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.runtime_version", "1.0"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.unauthenticated_client_action", "RedirectToLoginPage"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.token_refresh_extension_hours", "75"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.token_store_enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.additional_login_params.test_key", "test_value"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.allowed_external_redirect_urls.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.allowed_external_redirect_urls.0", "https://terra.form"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionAppSlot_enabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_enabled(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionAppSlot_enabledUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_enabled(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "false"),
				),
			},
			{
				Config: testAccAzureRMFunctionAppSlot_enabled(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionAppSlot_httpsOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_httpsOnly(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "https_only", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionAppSlot_httpsOnlyUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_httpsOnly(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "https_only", "true"),
				),
			},
			{
				Config: testAccAzureRMFunctionAppSlot_httpsOnly(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "https_only", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionAppSlot_http2Enabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_http2Enabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.http2_enabled", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionAppSlot_oneIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_oneIpRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.ip_address", "10.10.10.10/32"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionAppSlot_oneVNetSubnetIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_oneVNetSubnetIpRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionAppSlot_zeroedIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				// This configuration includes a single explicit ip_restriction
				Config: testAccAzureRMFunctionAppSlot_oneIpRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.#", "1"),
				),
			},
			{
				// This configuration has no site_config blocks at all.
				Config: testAccAzureRMFunctionAppSlot_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.#", "1"),
				),
			},
			{
				// This configuration explicitly sets ip_restriction to [] using attribute syntax.
				Config: testAccAzureRMFunctionAppSlot_zeroedIpRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionAppSlot_manyIpRestrictions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_manyIpRestrictions(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionAppSlot_scmUseMainIPRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_scmUseMainIPRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionAppSlot_scmIPRestrictionComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_scmIPRestrictionComplete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionAppSlot_tagsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_tags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Hello", "World"),
				),
			},
			{
				Config: testAccAzureRMFunctionAppSlot_tagsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Hello", "World"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Terraform", "AcceptanceTests"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionAppSlot_updateManageServiceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMFunctionAppSlot_enableManageServiceIdentity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.tenant_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionAppSlot_updateVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_version(data, "~1"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "~1"),
				),
			},
			{
				Config: testAccAzureRMFunctionAppSlot_version(data, "~2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "~2"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionAppSlot_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_userAssignedIdentity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "UserAssigned"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.identity_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.principal_id", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.tenant_id", ""),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionAppSlot_webSockets(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_webSockets(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.websockets_enabled", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionAppSlot_enableManageServiceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_enableManageServiceIdentity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.tenant_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionAppSlot_minTls(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_minTls(data, "1.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.min_tls_version", "1.0"),
				),
			},
			{
				Config: testAccAzureRMFunctionAppSlot_minTls(data, "1.1"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.min_tls_version", "1.1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionAppSlot_preWarmedInstanceCount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionAppSlot_preWarmedInstanceCount(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.pre_warmed_instance_count", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMFunctionAppSlotDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_function_app_slot" {
			continue
		}

		slot := rs.Primary.Attributes["name"]
		FunctionAppName := rs.Primary.Attributes["function_app_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetSlot(ctx, resourceGroup, FunctionAppName, slot)
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

func testCheckAzureRMFunctionAppSlotExists(slot string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[slot]
		if !ok {
			return fmt.Errorf("Slot Not found: %q", slot)
		}

		FunctionAppName := rs.Primary.Attributes["function_app_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Function App Slot: %q/%q", FunctionAppName, slot)
		}

		resp, err := client.GetSlot(ctx, resourceGroup, FunctionAppName, slot)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Function App slot %q/%q (resource group: %q) does not exist", FunctionAppName, slot, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on AppServicesClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMFunctionAppSlot_basic(data acceptance.TestData) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionAppSlot_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMFunctionAppSlot_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_function_app_slot" "import" {
  name                       = azurerm_function_app_slot.test.name
  location                   = azurerm_function_app_slot.test.location
  resource_group_name        = azurerm_function_app_slot.test.resource_group_name
  app_service_plan_id        = azurerm_function_app_slot.test.app_service_plan_id
  function_app_name          = azurerm_function_app_slot.test.function_app_name
  storage_account_name       = azurerm_function_app_slot.test.storage_account_name
  storage_account_access_key = azurerm_function_app_slot.test.storage_account_access_key
}
`, template)
}

func testAccAzureRMFunctionAppSlot_32Bit(data acceptance.TestData) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    use_32_bit_worker_process = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionAppSlot_alwaysOn(data acceptance.TestData) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    always_on = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionAppSlot_appSettings(data acceptance.TestData) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  app_settings = {
    "foo" = "bar"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionAppSlot_clientAffinityEnabled(data acceptance.TestData, clientAffinityEnabled bool) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  client_affinity_enabled    = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger, clientAffinityEnabled)
}

func testAccAzureRMFunctionAppSlot_connectionStrings(data acceptance.TestData) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionAppSlot_connectionStringsUpdated(data acceptance.TestData) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionAppSlot_corsSettings(data acceptance.TestData) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionAppSlot_autoSwap(data acceptance.TestData) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    auto_swap_slot_name = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionAppSlot_authSettings(data acceptance.TestData, tenantID string) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  auth_settings {
    enabled                       = true
    issuer                        = "https://sts.windows.net/%s"
    runtime_version               = "1.0"
    unauthenticated_client_action = "RedirectToLoginPage"
    token_refresh_extension_hours = 75
    token_store_enabled           = true

    additional_login_params = {
      test_key = "test_value"
    }

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger, tenantID)
}

func testAccAzureRMFunctionAppSlot_enabled(data acceptance.TestData, enabled bool) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  enabled                    = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger, enabled)
}

func testAccAzureRMFunctionAppSlot_httpsOnly(data acceptance.TestData, httpsOnly bool) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  https_only                 = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger, httpsOnly)
}

func testAccAzureRMFunctionAppSlot_http2Enabled(data acceptance.TestData) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    http2_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionAppSlot_oneIpRestriction(data acceptance.TestData) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10/32"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionAppSlot_oneVNetSubnetIpRestriction(data acceptance.TestData) string {
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

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    ip_restriction {
      virtual_network_subnet_id = azurerm_subnet.test.id
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionAppSlot_zeroedIpRestriction(data acceptance.TestData) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    ip_restriction = []
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionAppSlot_manyIpRestrictions(data acceptance.TestData) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionAppSlot_scmUseMainIPRestriction(data acceptance.TestData) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10/32"
    }

    scm_use_main_ip_restriction = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionAppSlot_scmIPRestrictionComplete(data acceptance.TestData) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    scm_ip_restriction {
      ip_address = "10.10.10.10/32"
      name       = "test-restriction"
      priority   = 123
      action     = "Allow"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionAppSlot_tags(data acceptance.TestData) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionAppSlot_tagsUpdated(data acceptance.TestData) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  tags = {
    "Hello"     = "World"
    "Terraform" = "AcceptanceTests"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionAppSlot_webSockets(data acceptance.TestData) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    websockets_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionAppSlot_enableManageServiceIdentity(data acceptance.TestData) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionAppSlot_version(data acceptance.TestData, version string) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  version                    = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger, version)
}

func testAccAzureRMFunctionAppSlot_userAssignedIdentity(data acceptance.TestData) string {
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

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionAppSlot_minTls(data acceptance.TestData, tlsVersion string) string {
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
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    min_tls_version = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger, tlsVersion)
}

func testAccAzureRMFunctionAppSlot_preWarmedInstanceCount(data acceptance.TestData) string {
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

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  kind                = "elastic"
  sku {
    tier = "ElasticPremium"
    size = "EP1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctestFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function_app_slot" "test" {
  name                       = "acctestFASlot-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  function_app_name          = azurerm_function_app.test.name
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    pre_warmed_instance_count = "1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
