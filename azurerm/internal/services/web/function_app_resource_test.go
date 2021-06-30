package web_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type FunctionAppResource struct{}

func TestAccFunctionApp_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.hasContentShareAppSetting(false)),
				check.That(data.ResourceName).Key("version").HasValue("~1"),
				check.That(data.ResourceName).Key("outbound_ip_addresses").Exists(),
				check.That(data.ResourceName).Key("possible_outbound_ip_addresses").Exists(),
				check.That(data.ResourceName).Key("custom_domain_verification_id").Exists(),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

// TODO remove in 3.0
func TestAccFunctionApp_deprecatedConnectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.deprecatedConnectionString(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// TODO remove in 3.0
func TestAccFunctionApp_deprecatedConnectionStringMissingError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.deprecatedConnectionStringMissingError(data),
			ExpectError: regexp.MustCompile("one of `storage_connection_string` or `storage_account_name` and `storage_account_access_key` must be specified"),
		},
	})
}

// TODO remove in 3.0
func TestAccFunctionApp_deprecatedNeedBothSAAtrributesError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.deprecatedConnectionStringBothSpecifiedError(data),
			ExpectError: regexp.MustCompile("both `storage_account_name` and `storage_account_access_key` must be specified"),
		},
	})
}

func TestAccFunctionApp_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.hasContentShareAppSetting(false)),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccFunctionApp_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("production"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFunctionApp_tagsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("production"),
			),
		},
		{
			Config: r.tagsUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("production"),
				check.That(data.ResourceName).Key("tags.hello").HasValue("Berlin"),
			),
		},
	})
}

func TestAccFunctionApp_appSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.appSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.appSettingsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("app_settings.%", "app_settings.AzureWebJobsDashboard", "app_settings.AzureWebJobsStorage"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFunctionApp_siteConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

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

func TestAccFunctionApp_healthCheck(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.healthCheck(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.health_check_path").HasValue("/"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFunctionApp_linuxFxVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxFxVersion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data.ResourceName, "kind", "functionapp,linux,container"),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("DOCKER|(golang:latest)"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFunctionApp_connectionStrings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.connectionStrings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("connection_string.#").HasValue("1"),
				check.That(data.ResourceName).Key("connection_string.163594034.name").HasValue("Example"),
				check.That(data.ResourceName).Key("connection_string.163594034.type").HasValue("PostgreSQL"),
				check.That(data.ResourceName).Key("connection_string.163594034.value").HasValue("some-postgresql-connection-string"),
			),
		},
		data.ImportStep(),
	})
}

// TODO - Refactor this into more granular tests - currently fails due to race condition in a `ForceNew` step when changed to `kind = linux`
func TestAccFunctionApp_siteConfigMulti(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.%").HasValue("0"),
			),
		},
		{
			Config: r.appSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.%").HasValue("1"),
				check.That(data.ResourceName).Key("app_settings.hello").HasValue("world"),
			),
		},
		{
			Config: r.appSettingsAlwaysOn(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.%").HasValue("1"),
				check.That(data.ResourceName).Key("app_settings.hello").HasValue("world"),
				check.That(data.ResourceName).Key("site_config.0.always_on").HasValue("true"),
			),
		},
		{
			Config: r.appSettingsAlwaysOnLinuxFxVersion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data.ResourceName, "kind", "functionapp,linux,container"),
				check.That(data.ResourceName).Key("app_settings.%").HasValue("1"),
				check.That(data.ResourceName).Key("app_settings.hello").HasValue("world"),
				check.That(data.ResourceName).Key("site_config.0.always_on").HasValue("true"),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("DOCKER|(golang:latest)"),
			),
		},
		{
			Config: r.appSettingsAlwaysOnLinuxFxVersionConnectionStrings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data.ResourceName, "kind", "functionapp,linux,container"),
				check.That(data.ResourceName).Key("app_settings.%").HasValue("1"),
				check.That(data.ResourceName).Key("app_settings.hello").HasValue("world"),
				check.That(data.ResourceName).Key("site_config.0.always_on").HasValue("true"),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("DOCKER|(golang:latest)"),
				check.That(data.ResourceName).Key("connection_string.163594034.name").HasValue("Example"),
				check.That(data.ResourceName).Key("connection_string.163594034.type").HasValue("PostgreSQL"),
				check.That(data.ResourceName).Key("connection_string.163594034.value").HasValue("some-postgresql-connection-string"),
			),
		},
	})
}

func TestAccFunctionApp_updateVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.version(data, "~1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("version").HasValue("~1"),
			),
		},
		{
			Config: r.version(data, "~2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("version").HasValue("~2"),
			),
		},
	})
}

func TestAccFunctionApp_3264bit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.use_32_bit_worker_process").HasValue("true"),
			),
		},
		{
			Config: r.app64bit(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.use_32_bit_worker_process").HasValue("false"),
			),
		},
	})
}

func TestAccFunctionApp_httpsOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.httpsOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_only").HasValue("true"),
			),
		},
	})
}

func TestAccFunctionApp_dailyMemoryTimeQuota(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.consumptionPlan(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("daily_memory_time_quota").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.dailyMemoryTimeQuota(data, 1000),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("daily_memory_time_quota").HasValue("1000"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFunctionApp_consumptionPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.consumptionPlan(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.hasContentShareAppSetting(true)),
				check.That(data.ResourceName).Key("site_config.0.use_32_bit_worker_process").HasValue("true"),
			),
		},
	})
}

func TestAccFunctionApp_consumptionPlanUppercaseName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.consumptionPlanUppercaseName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.hasContentShareAppSetting(true)),
				check.That(data.ResourceName).Key("site_config.0.use_32_bit_worker_process").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFunctionApp_createIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				acceptance.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
				acceptance.TestMatchResourceAttr(data.ResourceName, "identity.0.tenant_id", validate.UUIDRegExp),
			),
		},
	})
}

func TestAccFunctionApp_updateIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("0"),
			),
		},
		{
			Config: r.basicIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				acceptance.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
				acceptance.TestMatchResourceAttr(data.ResourceName, "identity.0.tenant_id", validate.UUIDRegExp),
			),
		},
	})
}

func TestAccFunctionApp_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

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
		data.ImportStep(),
		{
			Config: r.userAssignedIdentityUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("2"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsEmpty(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFunctionApp_loggingDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.loggingDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.hasContentShareAppSetting(false)),
				check.That(data.ResourceName).Key("enable_builtin_logging").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFunctionApp_updateLogging(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enable_builtin_logging").HasValue("true"),
			),
		},
		{
			Config: r.loggingDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enable_builtin_logging").HasValue("false"),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enable_builtin_logging").HasValue("true"),
			),
		},
	})
}

func TestAccFunctionApp_authSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	tenantID := os.Getenv("ARM_TENANT_ID")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authSettings(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_settings.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("auth_settings.0.issuer").HasValue(fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
				check.That(data.ResourceName).Key("auth_settings.0.runtime_version").HasValue("1.0"),
				check.That(data.ResourceName).Key("auth_settings.0.unauthenticated_client_action").HasValue("RedirectToLoginPage"),
				check.That(data.ResourceName).Key("auth_settings.0.token_refresh_extension_hours").HasValue("75"),
				check.That(data.ResourceName).Key("auth_settings.0.token_store_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("auth_settings.0.additional_login_params.test_key").HasValue("test_value"),
				check.That(data.ResourceName).Key("auth_settings.0.allowed_external_redirect_urls.#").HasValue("1"),
				check.That(data.ResourceName).Key("auth_settings.0.allowed_external_redirect_urls.0").HasValue("https://terra.form"),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.client_id").HasValue("aadclientid"),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.client_secret").HasValue("aadsecret"),
				check.That(data.ResourceName).Key("auth_settings.0.active_directory.0.allowed_audiences.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFunctionApp_corsSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.corsSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.cors.#").HasValue("1"),
				check.That(data.ResourceName).Key("site_config.0.cors.0.support_credentials").HasValue("true"),
				check.That(data.ResourceName).Key("site_config.0.cors.0.allowed_origins.#").HasValue("4"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFunctionApp_enableHttp2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.enableHttp2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.http2_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFunctionApp_minTlsVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.minTlsVersion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.min_tls_version").HasValue("1.2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFunctionApp_ftpsState(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

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

func TestAccFunctionApp_javaVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.javaVersion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFunctionApp_preWarmedInstanceCount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.preWarmedInstanceCount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.pre_warmed_instance_count").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMFunctionApp_computedPreWarmedInstanceCount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.computedPreWarmedInstanceCount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.pre_warmed_instance_count").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFunctionApp_oneIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.oneIpRestriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.ip_address").HasValue("10.10.10.10/32"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFunctionApp_oneServiceTagIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.oneServiceTagIpRestriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.service_tag").HasValue("AzureEventGrid"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFunctionApp_changeIpToServiceTagIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.oneIpRestriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.ip_address").HasValue("10.10.10.10/32"),
			),
		},
		{
			Config: r.oneServiceTagIpRestriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.service_tag").HasValue("AzureEventGrid"),
			),
		},
		{
			Config: r.oneIpRestriction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.0.ip_address").HasValue("10.10.10.10/32"),
			),
		},
	})
}

func TestAccFunctionApp_oneVNetSubnetIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

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

func TestAccFunctionApp_ipRestrictionRemoved(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// This configuration includes a single explicit ip_restriction
			Config: r.oneIpRestriction(data),
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
			Config: r.ipRestrictionRemoved(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.ip_restriction.#").HasValue("0"),
			),
		},
	})
}

func TestAccFunctionApp_manyIpRestrictions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.manyIpRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFunctionApp_scmUseMainIPRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

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

func TestAccFunctionApp_scmOneIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

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

func TestAccFunctionApp_updateStorageAccountKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateStorageAccountKey(data),
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

func TestAccFunctionApp_withSourceControl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

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

func TestAccFunctionApp_sourceControlUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

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

func TestAccFunctionApp_scm(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.scm(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFunctionApp_clientCertMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	r := FunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_cert_mode").HasValue(""),
			),
		},
		{
			Config: r.clientCertMode(data, "Required"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_cert_mode").HasValue("Required"),
			),
		},
		{
			Config: r.clientCertMode(data, "Optional"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_cert_mode").HasValue("Optional"),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_cert_mode").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func (r FunctionAppResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FunctionAppID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Web.AppServicesClient.Get(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Function App %q (Resource Group %q): %+v", id.SiteName, id.ResourceGroup, err)
	}

	// The SDK defines 404 as an "ok" status code..
	if utils.ResponseWasNotFound(resp.Response) {
		return utils.Bool(false), nil
	}

	return utils.Bool(true), nil
}

func (r FunctionAppResource) hasContentShareAppSetting(shouldExist bool) func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		id, err := parse.FunctionAppID(state.ID)
		if err != nil {
			return err
		}

		appSettingsResp, err := clients.Web.AppServicesClient.ListApplicationSettings(ctx, id.ResourceGroup, id.SiteName)
		if err != nil {
			return fmt.Errorf("listing AppSettings: %+v", err)
		}

		exists := false
		for k := range appSettingsResp.Properties {
			if strings.EqualFold("WEBSITE_CONTENTSHARE", k) {
				exists = true
				break
			}
		}
		if exists != shouldExist {
			return fmt.Errorf("expected %t but got %t", shouldExist, exists)
		}

		return nil
	}
}

func (r FunctionAppResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) withSourceControl(data acceptance.TestData, branch string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  source_control {
    repo_url           = "https://github.com/jackofallops/azure-app-service-static-site-tests.git"
    branch             = "%[4]s"
    manual_integration = true
    rollback_enabled   = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, branch)
}

func (r FunctionAppResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_function_app" "import" {
  name                       = azurerm_function_app.test.name
  location                   = azurerm_function_app.test.location
  resource_group_name        = azurerm_function_app.test.resource_group_name
  app_service_plan_id        = azurerm_function_app.test.app_service_plan_id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}
`, template)
}

func (r FunctionAppResource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) tagsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  tags = {
    environment = "production"
    hello       = "Berlin"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) version(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  version                    = "%[4]s"
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, version)
}

func (r FunctionAppResource) appSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  app_settings = {
    "hello" = "world"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) appSettingsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  app_settings = {
    "APPINSIGHTS_INSTRUMENTATIONKEY" = azurerm_storage_account.test.primary_connection_string
    "AzureWebJobsDashboard"          = azurerm_storage_account.test.primary_connection_string
    "AzureWebJobsStorage"            = azurerm_storage_account.test.primary_connection_string
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) alwaysOn(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    always_on = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) healthCheck(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    health_check_path = "/"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) linuxFxVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  kind                = "Linux"
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }

  reserved = true
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  version                    = "~2"
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  os_type                    = "linux"

  site_config {
    linux_fx_version = "DOCKER|(golang:latest)"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) connectionStrings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  connection_string {
    name  = "Example"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) appSettingsAlwaysOn(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  app_settings = {
    "hello" = "world"
  }

  site_config {
    always_on = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) appSettingsAlwaysOnLinuxFxVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  kind                = "Linux"
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }

  reserved = true
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  version                    = "~2"
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  app_settings = {
    "hello" = "world"
  }

  site_config {
    always_on        = true
    linux_fx_version = "DOCKER|(golang:latest)"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) appSettingsAlwaysOnLinuxFxVersionConnectionStrings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  kind                = "Linux"
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }

  reserved = true
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  version                    = "~2"
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  app_settings = {
    "hello" = "world"
  }

  site_config {
    always_on        = true
    linux_fx_version = "DOCKER|(golang:latest)"
  }

  connection_string {
    name  = "Example"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) app64bit(data acceptance.TestData) string {
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    use_32_bit_worker_process = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r FunctionAppResource) httpsOnly(data acceptance.TestData) string {
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  https_only                 = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r FunctionAppResource) dailyMemoryTimeQuota(data acceptance.TestData, dailyMemoryTimeQuota int) string {
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
  kind                = "FunctionApp"

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  daily_memory_time_quota    = %d
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger, dailyMemoryTimeQuota)
}

func (r FunctionAppResource) consumptionPlan(data acceptance.TestData) string {
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
  kind                = "FunctionApp"

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r FunctionAppResource) consumptionPlanUppercaseName(data acceptance.TestData) string {
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
  kind                = "FunctionApp"

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%d-FuncWithUppercase"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r FunctionAppResource) basicIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) userAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_user_assigned_identity" "first" {
  name                = "acctest1%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.first.id]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) userAssignedIdentityUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_user_assigned_identity" "first" {
  name                = "acctest1%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_user_assigned_identity" "second" {
  name                = "acctest2%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.first.id, azurerm_user_assigned_identity.second.id]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) loggingDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  enable_builtin_logging     = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) authSettings(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  auth_settings {
    enabled                       = true
    issuer                        = "https://sts.windows.net/%[4]s"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, tenantID)
}

func (r FunctionAppResource) corsSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    cors {
      allowed_origins = [
        "http://www.contoso.com",
        "www.contoso.com",
        "contoso.com",
        "http://localhost:4201",
      ]

      support_credentials = true
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) enableHttp2(data acceptance.TestData) string {
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    http2_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r FunctionAppResource) minTlsVersion(data acceptance.TestData) string {
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    min_tls_version = "1.2"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r FunctionAppResource) ftpsState(data acceptance.TestData) string {
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    ftps_state = "AllAllowed"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r FunctionAppResource) javaVersion(data acceptance.TestData) string {
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  version = "~3"

  site_config {
    java_version = "1.8"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r FunctionAppResource) preWarmedInstanceCount(data acceptance.TestData) string {
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
  name                      = "acctest-%d-func"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  app_service_plan_id       = azurerm_app_service_plan.test.id
  storage_connection_string = azurerm_storage_account.test.primary_connection_string

  site_config {
    pre_warmed_instance_count = 1
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r FunctionAppResource) computedPreWarmedInstanceCount(data acceptance.TestData) string {
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
  name                      = "acctest-%d-func"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  app_service_plan_id       = azurerm_app_service_plan.test.id
  storage_connection_string = azurerm_storage_account.test.primary_connection_string
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r FunctionAppResource) oneIpRestriction(data acceptance.TestData) string {
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10/32"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r FunctionAppResource) oneServiceTagIpRestriction(data acceptance.TestData) string {
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    ip_restriction {
      service_tag = "AzureEventGrid"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r FunctionAppResource) oneVNetSubnetIpRestriction(data acceptance.TestData) string {
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

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
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

resource "azurerm_function_app" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    ip_restriction {
      virtual_network_subnet_id = azurerm_subnet.test.id
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r FunctionAppResource) manyIpRestrictions(data acceptance.TestData) string {
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r FunctionAppResource) ipRestrictionRemoved(data acceptance.TestData) string {
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    ip_restriction = []
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r FunctionAppResource) scmUseMainIPRestriction(data acceptance.TestData) string {
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10/32"
    }
    scm_use_main_ip_restriction = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r FunctionAppResource) scmOneIpRestriction(data acceptance.TestData) string {
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    scm_ip_restriction {
      ip_address = "10.10.10.10/32"
      action     = "Allow"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r FunctionAppResource) deprecatedConnectionString(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  app_service_plan_id       = azurerm_app_service_plan.test.id
  storage_connection_string = azurerm_storage_account.test.primary_connection_string
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) deprecatedConnectionStringMissingError(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                = "acctest-%[1]d-func"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) deprecatedConnectionStringBothSpecifiedError(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                 = "acctest-%[1]d-func"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  app_service_plan_id  = azurerm_app_service_plan.test.id
  storage_account_name = azurerm_storage_account.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) updateStorageAccountKey(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.secondary_access_key
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) scm(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    always_on                 = true
    scm_type                  = "LocalGit"
    use_32_bit_worker_process = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppResource) clientCertMode(data acceptance.TestData, modeValue string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  client_cert_mode           = "%[4]s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, modeValue)
}
