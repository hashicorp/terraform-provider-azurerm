package logic_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogicAppStandardResource struct{}

func TestAccLogicAppStandard_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data.ResourceName, "kind", "functionapp,workflowapp"),
				check.That(data.ResourceName).Key("version").HasValue("~3"),
				check.That(data.ResourceName).Key("outbound_ip_addresses").Exists(),
				check.That(data.ResourceName).Key("possible_outbound_ip_addresses").Exists(),
				check.That(data.ResourceName).Key("custom_domain_verification_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppStandard_containerized(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.containerized(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data.ResourceName, "kind", "functionapp,linux,container,workflowapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppStandard_extensionBundle(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.extensionBundle(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.hasExtensionBundleAppSetting(true)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppStandard_siteConfigVnetRouteAllEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.siteConfigVnetRouteAllEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.vnet_route_all_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppStandard_appSettingsVnetRouteAllEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettingsVnetRouteAllEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.WEBSITE_VNET_ROUTE_ALL").HasValue("true"),
				check.That(data.ResourceName).Key("site_config.0.vnet_route_all_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppStandard_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_tagsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_appSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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
	})
}

func TestAccLogicAppStandard_customShare(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customShare(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppStandard_siteConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_healthCheck(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_connectionStrings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.connectionStrings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("connection_string.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppStandard_updateVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_3264bit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_httpsOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_createIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_updateIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_corsSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_enableHttp2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_minTlsVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_ftpsState(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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
func TestAccLogicAppStandard_preWarmedInstanceCount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_computedPreWarmedInstanceCount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_oneIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_oneServiceTagIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_changeIpToServiceTagIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_oneVNetSubnetIpRestriction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_ipRestrictionRemoved(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_manyIpRestrictions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_updateStorageAccountKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

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

func TestAccLogicAppStandard_clientCertMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_certificate_mode").HasValue(""),
			),
		},
		{
			Config: r.clientCertMode(data, "Required"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_certificate_mode").HasValue("Required"),
			),
		},
		{
			Config: r.clientCertMode(data, "Optional"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_certificate_mode").HasValue("Optional"),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_certificate_mode").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppStandard_elasticInstanceMinimum(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.elasticInstanceMinimum(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.elastic_instance_minimum").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppStandard_appScaleLimit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appScaleLimit(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.app_scale_limit").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppStandard_runtimeScaleMonitoringEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.runtimeScaleMonitoringEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.runtime_scale_monitoring_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppStandard_dotnetVersion4(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dotnetVersion(data, "~1", "v4.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.dotnet_framework_version").HasValue("v4.0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppStandard_dotnetVersion5(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dotnetVersion(data, "~3", "v5.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.dotnet_framework_version").HasValue("v5.0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppStandard_dotnetVersion6(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_standard", "test")
	r := LogicAppStandardResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dotnetVersion(data, "~4", "v6.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.dotnet_framework_version").HasValue("v6.0"),
			),
		},
		data.ImportStep(),
	})
}

func (r LogicAppStandardResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.LogicAppStandardID(state.ID)
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

func (r LogicAppStandardResource) hasExtensionBundleAppSetting(shouldExist bool) func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		id, err := parse.LogicAppStandardID(state.ID)
		if err != nil {
			return err
		}

		appSettingsResp, err := clients.Web.AppServicesClient.ListApplicationSettings(ctx, id.ResourceGroup, id.SiteName)
		if err != nil {
			return fmt.Errorf("listing AppSettings: %+v", err)
		}

		exists := false
		for k := range appSettingsResp.Properties {
			if strings.EqualFold("AzureFunctionsJobHost__extensionBundle__id", k) {
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

func (r LogicAppStandardResource) basic(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r LogicAppStandardResource) containerized(data acceptance.TestData) string {
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
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_logic_app_standard" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    linux_fx_version = "DOCKER|mcr.microsoft.com/azure-functions/dotnet:3.0-appservice"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r LogicAppStandardResource) extensionBundle(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  use_extension_bundle       = true
  bundle_version             = "[1.*, 2.0.0)"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r LogicAppStandardResource) siteConfigVnetRouteAllEnabled(data acceptance.TestData) string {
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
resource "azurerm_logic_app_standard" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  site_config {
    vnet_route_all_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r LogicAppStandardResource) appSettingsVnetRouteAllEnabled(data acceptance.TestData) string {
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
resource "azurerm_logic_app_standard" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  app_settings = {
    "WEBSITE_VNET_ROUTE_ALL" = "true"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r LogicAppStandardResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_standard" "import" {
  name                       = azurerm_logic_app_standard.test.name
  location                   = azurerm_logic_app_standard.test.location
  resource_group_name        = azurerm_logic_app_standard.test.resource_group_name
  app_service_plan_id        = azurerm_logic_app_standard.test.app_service_plan_id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}
`, template)
}

func (r LogicAppStandardResource) tags(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
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

func (r LogicAppStandardResource) tagsUpdated(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
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

func (r LogicAppStandardResource) version(data acceptance.TestData, version string) string {
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

resource "azurerm_logic_app_standard" "test" {
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

func (r LogicAppStandardResource) appSettings(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  storage_account_share_name = "acctest-%[1]d-func-content"

  app_settings = {
    "hello"                          = "world"
    "APPINSIGHTS_INSTRUMENTATIONKEY" = azurerm_storage_account.test.primary_connection_string
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r LogicAppStandardResource) customShare(data acceptance.TestData) string {
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

resource "azurerm_storage_share" "custom" {
  name                 = "customshare"
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_logic_app_standard" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  storage_account_share_name = azurerm_storage_share.custom.name

  app_settings = {
    "hello"                          = "world"
    "APPINSIGHTS_INSTRUMENTATIONKEY" = azurerm_storage_account.test.primary_connection_string
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r LogicAppStandardResource) alwaysOn(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
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

func (r LogicAppStandardResource) healthCheck(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
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

func (r LogicAppStandardResource) connectionStrings(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
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

func (r LogicAppStandardResource) app64bit(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
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

func (r LogicAppStandardResource) httpsOnly(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
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

func (r LogicAppStandardResource) basicIdentity(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
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

func (r LogicAppStandardResource) corsSettings(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
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

func (r LogicAppStandardResource) enableHttp2(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
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

func (r LogicAppStandardResource) minTlsVersion(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
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

func (r LogicAppStandardResource) ftpsState(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
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

func (r LogicAppStandardResource) preWarmedInstanceCount(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    pre_warmed_instance_count = 1
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r LogicAppStandardResource) computedPreWarmedInstanceCount(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r LogicAppStandardResource) oneIpRestriction(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
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

func (r LogicAppStandardResource) oneServiceTagIpRestriction(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
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

func (r LogicAppStandardResource) oneVNetSubnetIpRestriction(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
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

func (r LogicAppStandardResource) manyIpRestrictions(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
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

func (r LogicAppStandardResource) ipRestrictionRemoved(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
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

func (r LogicAppStandardResource) updateStorageAccountKey(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.secondary_access_key
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r LogicAppStandardResource) clientCertMode(data acceptance.TestData, modeValue string) string {
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

resource "azurerm_logic_app_standard" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  client_certificate_mode    = "%[4]s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, modeValue)
}

func (r LogicAppStandardResource) elasticInstanceMinimum(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    elastic_instance_minimum = 1
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r LogicAppStandardResource) appScaleLimit(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    app_scale_limit = 1
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r LogicAppStandardResource) runtimeScaleMonitoringEnabled(data acceptance.TestData) string {
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

resource "azurerm_logic_app_standard" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  version                    = "~3"

  site_config {
    pre_warmed_instance_count        = 1
    runtime_scale_monitoring_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r LogicAppStandardResource) dotnetVersion(data acceptance.TestData, functionVersion string, version string) string {
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

resource "azurerm_logic_app_standard" "test" {
  name                       = "acctest-%d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  version = "%s"

  site_config {
    dotnet_framework_version = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger, functionVersion, version)
}
