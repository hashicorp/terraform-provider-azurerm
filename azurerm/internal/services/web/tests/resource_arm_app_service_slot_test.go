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

func TestAccAzureRMAppServiceSlot_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceSlot_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMAppServiceSlot_requiresImport),
		},
	})
}

func TestAccAzureRMAppServiceSlot_32Bit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_32Bit(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.use_32_bit_worker_process", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceSlot_alwaysOn(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_alwaysOn(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.always_on", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceSlot_appCommandLine(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_appCommandLine(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.app_command_line", "/sbin/myservice -b 0.0.0.0"),
				),
			},

			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceSlot_appSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_appSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "app_settings.foo", "bar"),
				),
			},

			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceSlot_clientAffinityEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_clientAffinityEnabled(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "client_affinity_enabled", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_clientAffinityEnabledUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_clientAffinityEnabled(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "client_affinity_enabled", "true"),
				),
			},
			{
				Config: testAccAzureRMAppServiceSlot_clientAffinityEnabled(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "client_affinity_enabled", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_connectionStrings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_connectionStrings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.3173438943.name", "First"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.3173438943.value", "first-connection-string"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.3173438943.type", "Custom"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.2442860602.name", "Second"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.2442860602.value", "some-postgresql-connection-string"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.2442860602.type", "PostgreSQL"),
				),
			},
			{
				Config: testAccAzureRMAppServiceSlot_connectionStringsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
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

func TestAccAzureRMAppServiceSlot_corsSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_corsSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceSlot_authSettingsAdditionalLoginParams(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")
	tenantID := os.Getenv("ARM_TENANT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_authSettingsAdditionalLoginParams(data, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
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

func TestAccAzureRMAppServiceSlot_authSettingsAdditionalAllowedExternalRedirectUrls(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")
	tenantID := os.Getenv("ARM_TENANT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_authSettingsAdditionalAllowedExternalRedirectUrls(data, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
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

func TestAccAzureRMAppServiceSlot_authSettingsRuntimeVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")
	tenantID := os.Getenv("ARM_TENANT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_authSettingsRuntimeVersion(data, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
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

func TestAccAzureRMAppServiceSlot_authSettingsTokenRefreshExtensionHours(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")
	tenantID := os.Getenv("ARM_TENANT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_authSettingsTokenRefreshExtensionHours(data, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
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

func TestAccAzureRMAppServiceSlot_authSettingsUnauthenticatedClientAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")
	tenantID := os.Getenv("ARM_TENANT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_authSettingsUnauthenticatedClientAction(data, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
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

func TestAccAzureRMAppServiceSlot_authSettingsTokenStoreEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")
	tenantID := os.Getenv("ARM_TENANT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_authSettingsTokenStoreEnabled(data, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
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

func TestAccAzureRMAppServiceSlot_aadAuthSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")
	tenantID := os.Getenv("ARM_TENANT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_aadAuthSettings(data, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
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

func TestAccAzureRMAppServiceSlot_facebookAuthSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_facebookAuthSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
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

func TestAccAzureRMAppServiceSlot_googleAuthSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_googleAuthSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.google.0.client_id", "googleclientid"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.google.0.client_secret", "googleclientsecret"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.google.0.oauth_scopes.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_microsoftAuthSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_microsoftAuthSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
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

func TestAccAzureRMAppServiceSlot_twitterAuthSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_twitterAuthSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.twitter.0.consumer_key", "twitterconsumerkey"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.twitter.0.consumer_secret", "twitterconsumersecret"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceSlot_multiAuthSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")
	tenantID := os.Getenv("ARM_TENANT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_aadAuthSettings(data, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAppServiceSlot_aadMicrosoftAuthSettings(data, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
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

func TestAccAzureRMAppServiceSlot_defaultDocuments(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_defaultDocuments(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.default_documents.0", "first.html"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.default_documents.1", "second.jsp"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.default_documents.2", "third.aspx"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_enabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_enabled(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_enabledUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_enabled(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "false"),
				),
			},
			{
				Config: testAccAzureRMAppServiceSlot_enabled(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_httpsOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_httpsOnly(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "https_only", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_httpsOnlyUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_httpsOnly(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "https_only", "true"),
				),
			},
			{
				Config: testAccAzureRMAppServiceSlot_httpsOnly(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "https_only", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_http2Enabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_http2Enabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.http2_enabled", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_oneIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_oneIpRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.ip_address", "10.10.10.10/32"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_oneVNetSubnetIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_oneVNetSubnetIpRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceSlot_zeroedIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				// This configuration includes a single explicit ip_restriction
				Config: testAccAzureRMAppServiceSlot_oneIpRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.#", "1"),
				),
			},
			{
				// This configuration has no site_config blocks at all.
				Config: testAccAzureRMAppServiceSlot_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.#", "1"),
				),
			},
			{
				// This configuration explicitly sets ip_restriction to [] using attribute syntax.
				Config: testAccAzureRMAppServiceSlot_zeroedIpRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_manyIpRestrictions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_manyIpRestrictions(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceSlot_scmUseMainIPRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_scmUseMainIPRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceSlot_scmOneIPRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_scmOneIPRestriction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceSlot_localMySql(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_localMySql(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.local_mysql_enabled", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_managedPipelineMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_managedPipelineMode(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.managed_pipeline_mode", "Classic"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_tagsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_tags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Hello", "World"),
				),
			},
			{
				Config: testAccAzureRMAppServiceSlot_tagsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Hello", "World"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Terraform", "AcceptanceTests"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_remoteDebugging(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_remoteDebugging(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.remote_debugging_enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.remote_debugging_version", "VS2019"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_windowsDotNet2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_windowsDotNet(data, "v2.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.dotnet_framework_version", "v2.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_updateManageServiceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMAppServiceSlot_enableManageServiceIdentity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.tenant_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_windowsDotNet4(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_windowsDotNet(data, "v4.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.dotnet_framework_version", "v4.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_userAssignedIdentity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "UserAssigned"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.identity_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.principal_id", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.tenant_id", ""),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_windowsDotNetUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_windowsDotNet(data, "v2.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.dotnet_framework_version", "v2.0"),
				),
			},
			{
				Config: testAccAzureRMAppServiceSlot_windowsDotNet(data, "v4.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.dotnet_framework_version", "v4.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_windowsJava7Jetty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_windowsJava(data, "1.7", "JETTY", "9.3"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_version", "1.7"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container", "JETTY"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container_version", "9.3"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_windowsJava8Jetty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_windowsJava(data, "1.8", "JETTY", "9.3"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_version", "1.8"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container", "JETTY"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container_version", "9.3"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_windowsJava11Jetty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_windowsJava(data, "11", "JETTY", "9.3"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_version", "11"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container", "JETTY"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container_version", "9.3"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_windowsJava7Tomcat(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_windowsJava(data, "1.7", "TOMCAT", "9.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_version", "1.7"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container", "TOMCAT"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container_version", "9.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_windowsJava8Tomcat(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_windowsJava(data, "1.8", "TOMCAT", "9.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_version", "1.8"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container", "TOMCAT"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container_version", "9.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_windowsJava11Tomcat(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_windowsJava(data, "11", "TOMCAT", "9.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_version", "11"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container", "TOMCAT"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container_version", "9.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_windowsJava7Minor(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_windowsJava(data, "1.7.0_80", "TOMCAT", "9.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_version", "1.7.0_80"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container", "TOMCAT"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container_version", "9.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_windowsJava8Minor(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_windowsJava(data, "1.8.0_181", "TOMCAT", "9.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_version", "1.8.0_181"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container", "TOMCAT"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.java_container_version", "9.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_windowsPHP7(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_windowsPHP(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.php_version", "7.3"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_windowsPython(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_windowsPython(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.python_version", "3.4"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_webSockets(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_webSockets(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.websockets_enabled", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_enableManageServiceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_enableManageServiceIdentity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.tenant_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceSlot_minTls(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_minTls(data, "1.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.min_tls_version", "1.0"),
				),
			},
			{
				Config: testAccAzureRMAppServiceSlot_minTls(data, "1.1"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.min_tls_version", "1.1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceSlot_applicationBlobStorageLogs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_applicationBlobStorageLogs(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "logs.0.application_logs.0.file_system_level", "Warning"),
					resource.TestCheckResourceAttr(data.ResourceName, "logs.0.application_logs.0.azure_blob_storage.0.level", "Information"),
					resource.TestCheckResourceAttr(data.ResourceName, "logs.0.application_logs.0.azure_blob_storage.0.sas_url", "https://example.com/"),
					resource.TestCheckResourceAttr(data.ResourceName, "logs.0.application_logs.0.azure_blob_storage.0.retention_in_days", "3"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceSlot_httpFileSystemLogs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_httpFileSystemLogs(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "logs.0.http_logs.0.file_system.0.retention_in_days", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "logs.0.http_logs.0.file_system.0.retention_in_mb", "25"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceSlot_httpBlobStorageLogs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_httpBlobStorageLogs(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "logs.0.http_logs.0.azure_blob_storage.0.sas_url", "https://example.com/"),
					resource.TestCheckResourceAttr(data.ResourceName, "logs.0.http_logs.0.azure_blob_storage.0.retention_in_days", "3"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceSlot_autoSwap(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_slot", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSlot_autoSwap(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceSlotExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.auto_swap_slot_name", "production"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMAppServiceSlotDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_slot" {
			continue
		}

		slot := rs.Primary.Attributes["name"]
		appServiceName := rs.Primary.Attributes["app_service_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetSlot(ctx, resourceGroup, appServiceName, slot)

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

func testCheckAzureRMAppServiceSlotExists(slot string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[slot]
		if !ok {
			return fmt.Errorf("Slot Not found: %q", slot)
		}

		appServiceName := rs.Primary.Attributes["app_service_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Service Slot: %q/%q", appServiceName, slot)
		}

		resp, err := client.GetSlot(ctx, resourceGroup, appServiceName, slot)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: App Service slot %q/%q (resource group: %q) does not exist", appServiceName, slot, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on appServicesClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMAppServiceSlot_basic(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMAppServiceSlot_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_slot" "import" {
  name                = azurerm_app_service_slot.test.name
  location            = azurerm_app_service_slot.test.location
  resource_group_name = azurerm_app_service_slot.test.resource_group_name
  app_service_plan_id = azurerm_app_service_slot.test.app_service_plan_id
  app_service_name    = azurerm_app_service_slot.test.app_service_name
}
`, template)
}

func testAccAzureRMAppServiceSlot_32Bit(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  site_config {
    use_32_bit_worker_process = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_alwaysOn(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  site_config {
    always_on = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_appCommandLine(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  site_config {
    app_command_line = "/sbin/myservice -b 0.0.0.0"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_appSettings(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  app_settings = {
    "foo" = "bar"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_clientAffinityEnabled(data acceptance.TestData, clientAffinityEnabled bool) string {
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

resource "azurerm_app_service_slot" "test" {
  name                    = "acctestASSlot-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  app_service_plan_id     = azurerm_app_service_plan.test.id
  app_service_name        = azurerm_app_service.test.name
  client_affinity_enabled = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, clientAffinityEnabled)
}

func testAccAzureRMAppServiceSlot_connectionStrings(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_connectionStringsUpdated(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_corsSettings(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_authSettingsAdditionalLoginParams(data acceptance.TestData, tenantID string) string {
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
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, tenantID)
}

func testAccAzureRMAppServiceSlot_authSettingsAdditionalAllowedExternalRedirectUrls(data acceptance.TestData, tenantID string) string {
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
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, tenantID)
}

func testAccAzureRMAppServiceSlot_authSettingsRuntimeVersion(data acceptance.TestData, tenantID string) string {
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
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, tenantID)
}

func testAccAzureRMAppServiceSlot_authSettingsTokenRefreshExtensionHours(data acceptance.TestData, tenantID string) string {
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
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, tenantID)
}

func testAccAzureRMAppServiceSlot_authSettingsTokenStoreEnabled(data acceptance.TestData, tenantID string) string {
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
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, tenantID)
}

func testAccAzureRMAppServiceSlot_authSettingsUnauthenticatedClientAction(data acceptance.TestData, tenantID string) string {
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
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, tenantID)
}

func testAccAzureRMAppServiceSlot_aadAuthSettings(data acceptance.TestData, tenantID string) string {
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
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, tenantID)
}

func testAccAzureRMAppServiceSlot_facebookAuthSettings(data acceptance.TestData) string {
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
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_googleAuthSettings(data acceptance.TestData) string {
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
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_microsoftAuthSettings(data acceptance.TestData) string {
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
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_twitterAuthSettings(data acceptance.TestData) string {
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
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_aadMicrosoftAuthSettings(data acceptance.TestData, tenantID string) string {
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
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, tenantID, web.AzureActiveDirectory)
}

func testAccAzureRMAppServiceSlot_defaultDocuments(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  site_config {
    default_documents = [
      "first.html",
      "second.jsp",
      "third.aspx",
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_enabled(data acceptance.TestData, enabled bool) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name
  enabled             = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, enabled)
}

func testAccAzureRMAppServiceSlot_httpsOnly(data acceptance.TestData, httpsOnly bool) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name
  https_only          = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, httpsOnly)
}

func testAccAzureRMAppServiceSlot_http2Enabled(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  site_config {
    http2_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_oneIpRestriction(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10/32"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_oneVNetSubnetIpRestriction(data acceptance.TestData) string {
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
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  site_config {
    ip_restriction {
      virtual_network_subnet_id = azurerm_subnet.test.id
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_scmUseMainIPRestriction(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-web-%d"
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
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  site_config {
    scm_use_main_ip_restriction = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_scmOneIPRestriction(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-web-%d"
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
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  site_config {
    scm_ip_restriction {
      ip_address = "10.10.10.10/32"
      action     = "Allow"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_zeroedIpRestriction(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  site_config {
    ip_restriction = []
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_manyIpRestrictions(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_localMySql(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  site_config {
    local_mysql_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_managedPipelineMode(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  site_config {
    managed_pipeline_mode = "Classic"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_tags(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_tagsUpdated(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  tags = {
    "Hello"     = "World"
    "Terraform" = "AcceptanceTests"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_remoteDebugging(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  site_config {
    remote_debugging_enabled = true
    remote_debugging_version = "VS2019"
  }

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_windowsDotNet(data acceptance.TestData, version string) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  site_config {
    dotnet_framework_version = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, version)
}

func testAccAzureRMAppServiceSlot_windowsJava(data acceptance.TestData, javaVersion, container, containerVersion string) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  site_config {
    java_version           = "%s"
    java_container         = "%s"
    java_container_version = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, javaVersion, container, containerVersion)
}

func testAccAzureRMAppServiceSlot_windowsPHP(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  site_config {
    php_version = "7.3"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_windowsPython(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  site_config {
    python_version = "3.4"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_webSockets(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  site_config {
    websockets_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_enableManageServiceIdentity(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_userAssignedIdentity(data acceptance.TestData) string {
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
}

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_minTls(data acceptance.TestData, tlsVersion string) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  site_config {
    min_tls_version = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, tlsVersion)
}

func testAccAzureRMAppServiceSlot_applicationBlobStorageLogs(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  logs {
    application_logs {
      file_system_level = "Warning"
      azure_blob_storage {
        level             = "Information"
        sas_url           = "https://example.com/"
        retention_in_days = 3
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_httpFileSystemLogs(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  logs {
    http_logs {
      file_system {
        retention_in_days = 4
        retention_in_mb   = 25
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_httpBlobStorageLogs(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  logs {
    http_logs {
      azure_blob_storage {
        sas_url           = "https://example.com/"
        retention_in_days = 3
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceSlot_autoSwap(data acceptance.TestData) string {
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

resource "azurerm_app_service_slot" "test" {
  name                = "acctestASSlot-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
  app_service_name    = azurerm_app_service.test.name

  site_config {
    auto_swap_slot_name = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
