// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

type ApplicationGatewayResource struct{}

func TestAccApplicationGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_Small"),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("Standard"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("2"),
				check.That(data.ResourceName).Key("waf_configuration.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_autoscaleConfiguration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoscaleConfiguration(data, 0, 10),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_v2"),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("Standard_v2"),
				check.That(data.ResourceName).Key("autoscale_configuration.0.min_capacity").HasValue("0"),
				check.That(data.ResourceName).Key("autoscale_configuration.0.max_capacity").HasValue("10"),
				check.That(data.ResourceName).Key("waf_configuration.#").HasValue("0"),
			),
		},
		{
			Config: r.autoscaleConfiguration(data, 4, 12),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_v2"),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("Standard_v2"),
				check.That(data.ResourceName).Key("autoscale_configuration.0.min_capacity").HasValue("4"),
				check.That(data.ResourceName).Key("autoscale_configuration.0.max_capacity").HasValue("12"),
				check.That(data.ResourceName).Key("waf_configuration.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_autoscaleConfigurationNoMaxCapacity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoscaleConfigurationNoMaxCapacity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_v2"),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("Standard_v2"),
				check.That(data.ResourceName).Key("autoscale_configuration.0.min_capacity").HasValue("2"),
				check.That(data.ResourceName).Key("waf_configuration.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_globalConfiguration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.createGlobalConfiguration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_v2"),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("Standard_v2"),
				check.That(data.ResourceName).Key("global.0.request_buffering_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("global.0.response_buffering_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateGlobalConfiguration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_v2"),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("Standard_v2"),
				check.That(data.ResourceName).Key("global.0.request_buffering_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("global.0.response_buffering_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_zones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zones(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_v2"),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("Standard_v2"),
				check.That(data.ResourceName).Key("zones.#").HasValue("2"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("2"),
				check.That(data.ResourceName).Key("waf_configuration.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_overridePath(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.overridePath(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backend_http_settings.0.path").HasValue("/path1/"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_http2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.http2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enable_http2").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_application_gateway"),
		},
	})
}

func TestAccApplicationGateway_authCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authCertificate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("authentication_certificate.0.name").Exists(),
			),
		},
		// since these are read from the existing state
		data.ImportStep(

			"authentication_certificate.0.data",
		),
		{
			Config: r.authCertificateUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("authentication_certificate.0.name").Exists(),
			),
		},
		// since these are read from the existing state
		data.ImportStep(

			"authentication_certificate.0.data",
		),
	})
}

func TestAccApplicationGateway_customFirewallPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customFirewallPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("firewall_policy_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_customHttpListenerFirewallPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customHttpListenerFirewallPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("http_listener.0.firewall_policy_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_customPathRuleFirewallPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customPathRuleFirewallPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_trustedRootCertificate_keyvault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.trustedRootCertificate_keyvault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("trusted_root_certificate.0.name").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.update_trustedRootCertificate_keyvault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("trusted_root_certificate.0.name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_trustedRootCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.trustedRootCertificate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("trusted_root_certificate.0.name").Exists(),
			),
		},
		// since these are read from the existing state
		data.ImportStep(
			"trusted_root_certificate.0.data",
		),
		{
			Config: r.trustedRootCertificateUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("trusted_root_certificate.0.name").Exists(),
			),
		},
		// since these are read from the existing state
		data.ImportStep(
			"trusted_root_certificate.0.data",
		),
	})
}

func TestAccApplicationGateway_pathBasedRouting(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.pathBasedRouting(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_routingRedirect_httpListener(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.routingRedirect_httpListener(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_routingRedirect_httpListenerError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.routingRedirect_httpListenerError(data),
			ExpectError: regexp.MustCompile("Conflict between `backend_address_pool_name` and `redirect_configuration_name`"),
		},
	})
}

func TestAccApplicationGateway_routingRedirect_pathBased(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.routingRedirect_pathBased(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_routingRedirect_pathBasedWithRedirection(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.pathBasedRoutingwithRedirection(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_customErrorConfigurations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customErrorConfigurations(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_rewriteRuleSets_backend(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.rewriteRuleSets_backend(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rewrite_rule_set.0.name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_rewriteRuleSets_redirect(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.rewriteRuleSets_redirect(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rewrite_rule_set.0.name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_rewriteRuleSets_rewriteUrl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.rewriteRuleSets_rewriteUrl(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rewrite_rule_set.0.name").Exists(),
				check.That(data.ResourceName).Key("rewrite_rule_set.#").HasValue("2"),
				check.That(data.ResourceName).Key("rewrite_rule_set.0.rewrite_rule.0.url.0.components").HasValue(""),
				check.That(data.ResourceName).Key("rewrite_rule_set.1.rewrite_rule.#").HasValue("2"),
				check.That(data.ResourceName).Key("rewrite_rule_set.1.rewrite_rule.0.url.0.components").HasValue("path_only"),
				check.That(data.ResourceName).Key("rewrite_rule_set.1.rewrite_rule.1.url.0.components").HasValue("query_string_only"),
			),
		},
		data.ImportStep(),
		{
			Config: r.rewriteRuleSets_rewriteUrlUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rewrite_rule_set.0.name").Exists(),
				check.That(data.ResourceName).Key("rewrite_rule_set.#").HasValue("2"),
				check.That(data.ResourceName).Key("rewrite_rule_set.0.rewrite_rule.0.url.0.components").HasValue(""),
				check.That(data.ResourceName).Key("rewrite_rule_set.1.rewrite_rule.0.url.0.components").HasValue(""),
				check.That(data.ResourceName).Key("rewrite_rule_set.1.rewrite_rule.1.url.0.components").HasValue("query_string_only"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_probes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.probes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_probesEmptyMatch(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.probesEmptyMatch(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_probesPickHostNameFromBackendHTTPSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.probesPickHostNameFromBackendHTTPSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("probe.0.pick_host_name_from_backend_http_settings").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_probesWithPort(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.probesWithPort(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("probe.0.port").HasValue("8082"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_backendHttpSettingsHostName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}
	hostName := "example.com"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backendHttpSettingsHostName(data, hostName, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backend_http_settings.0.host_name").HasValue(hostName),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_withHttpListenerHostNames(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withHttpListenerHostNames(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_backendHttpSettingsHostNameAndPick(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}
	hostName := "example.com"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.backendHttpSettingsHostName(data, hostName, true),
			ExpectError: regexp.MustCompile("Only one of `host_name` or `pick_host_name_from_backend_address` can be set"),
		},
	})
}

func TestAccApplicationGateway_settingsPickHostNameFromBackendAddress(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.settingsPickHostNameFromBackendAddress(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backend_http_settings.0.pick_host_name_from_backend_address").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_sslCertificate_keyvault_versionless(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sslCertificate_keyvault_versionless(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ssl_certificate.0.key_vault_secret_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_sslCertificate_keyvault_versioned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sslCertificate_keyvault_versioned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ssl_certificate.0.key_vault_secret_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_sslCertificate_EmptyPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sslCertificateEmptyPassword(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// since these are read from the existing state
		data.ImportStep(
			"ssl_certificate.0.data",
			"ssl_certificate.0.password",
		),
	})
}

func TestAccApplicationGateway_sslCertificateManualChangeIgnoreChanges(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.manualSslCertificateChangeIgnoreChangesConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ssl_certificate.0.name").HasValue("acctestcertificate1"),
				data.CheckWithClient(r.changeCert("acctestcertificate2")),
			),
		},
		{
			Config: r.manualSslCertificateChangeIgnoreChangesUpdatedConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ssl_certificate.0.name").HasValue("acctestcertificate2"),
			),
		},
	})
}

func TestAccApplicationGateway_sslCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sslCertificate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// since these are read from the existing state
		data.ImportStep(
			"ssl_certificate.0.data",
			"ssl_certificate.0.password",
		),
		{
			Config: r.sslCertificateUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// since these are read from the existing state
		data.ImportStep(
			"ssl_certificate.0.data",
			"ssl_certificate.0.password",
		),
	})
}

func TestAccApplicationGateway_webApplicationFirewall(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.webApplicationFirewall(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("WAF_Medium"),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("WAF"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("1"),
				check.That(data.ResourceName).Key("waf_configuration.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("waf_configuration.0.firewall_mode").HasValue("Detection"),
				check.That(data.ResourceName).Key("waf_configuration.0.rule_set_type").HasValue("OWASP"),
				check.That(data.ResourceName).Key("waf_configuration.0.rule_set_version").HasValue("3.0"),
				check.That(data.ResourceName).Key("waf_configuration.0.file_upload_limit_mb").HasValue("100"),
				check.That(data.ResourceName).Key("waf_configuration.0.request_body_check").HasValue("true"),
				check.That(data.ResourceName).Key("waf_configuration.0.max_request_body_size_kb").HasValue("100"),
			),
		},
	})
}

func TestAccApplicationGateway_connectionDraining(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.connectionDraining(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backend_http_settings.0.connection_draining.0.enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_Small"),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("Standard"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("2"),
				check.That(data.ResourceName).Key("waf_configuration.#").HasValue("0"),
				acceptance.TestCheckNoResourceAttr(data.ResourceName, "backend_http_settings.0.connection_draining.0.enabled"),
				acceptance.TestCheckNoResourceAttr(data.ResourceName, "backend_http_settings.0.connection_draining.0.drain_timeout_sec"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_webApplicationFirewall_disabledRuleGroups(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.webApplicationFirewall_disabledRuleGroups(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("WAF_v2"),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("WAF_v2"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("1"),
				check.That(data.ResourceName).Key("waf_configuration.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("waf_configuration.0.firewall_mode").HasValue("Detection"),
				check.That(data.ResourceName).Key("waf_configuration.0.rule_set_type").HasValue("OWASP"),
				check.That(data.ResourceName).Key("waf_configuration.0.rule_set_version").HasValue("3.0"),
				check.That(data.ResourceName).Key("waf_configuration.0.request_body_check").HasValue("true"),
				check.That(data.ResourceName).Key("waf_configuration.0.max_request_body_size_kb").HasValue("128"),
				check.That(data.ResourceName).Key("waf_configuration.0.file_upload_limit_mb").HasValue("100"),
				check.That(data.ResourceName).Key("waf_configuration.0.disabled_rule_group.0.rule_group_name").HasValue("REQUEST-921-PROTOCOL-ATTACK"),
				check.That(data.ResourceName).Key("waf_configuration.0.disabled_rule_group.0.rules.0").HasValue("921110"),
				check.That(data.ResourceName).Key("waf_configuration.0.disabled_rule_group.0.rules.1").HasValue("921151"),
				check.That(data.ResourceName).Key("waf_configuration.0.disabled_rule_group.0.rules.2").HasValue("921180"),
				check.That(data.ResourceName).Key("waf_configuration.0.disabled_rule_group.1.rule_group_name").HasValue("REQUEST-930-APPLICATION-ATTACK-LFI"),
				check.That(data.ResourceName).Key("waf_configuration.0.disabled_rule_group.1.rules.0").HasValue("930120"),
				check.That(data.ResourceName).Key("waf_configuration.0.disabled_rule_group.1.rules.1").HasValue("930130"),
				check.That(data.ResourceName).Key("waf_configuration.0.disabled_rule_group.2.rule_group_name").HasValue("REQUEST-942-APPLICATION-ATTACK-SQLI"),
			),
		},
		{
			Config: r.webApplicationFirewall_disabledRuleGroups_enabled_some_rules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("WAF_v2"),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("WAF_v2"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("1"),
				check.That(data.ResourceName).Key("waf_configuration.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("waf_configuration.0.firewall_mode").HasValue("Detection"),
				check.That(data.ResourceName).Key("waf_configuration.0.rule_set_type").HasValue("OWASP"),
				check.That(data.ResourceName).Key("waf_configuration.0.rule_set_version").HasValue("3.0"),
				check.That(data.ResourceName).Key("waf_configuration.0.request_body_check").HasValue("true"),
				check.That(data.ResourceName).Key("waf_configuration.0.max_request_body_size_kb").HasValue("128"),
				check.That(data.ResourceName).Key("waf_configuration.0.file_upload_limit_mb").HasValue("100"),
				check.That(data.ResourceName).Key("waf_configuration.0.disabled_rule_group.0.rule_group_name").HasValue("REQUEST-921-PROTOCOL-ATTACK"),
				check.That(data.ResourceName).Key("waf_configuration.0.disabled_rule_group.0.rules.0").HasValue("921110"),
				check.That(data.ResourceName).Key("waf_configuration.0.disabled_rule_group.0.rules.1").HasValue("921151"),
				check.That(data.ResourceName).Key("waf_configuration.0.disabled_rule_group.0.rules.2").HasValue("921180"),
				check.That(data.ResourceName).Key("waf_configuration.0.disabled_rule_group.1.rule_group_name").HasValue("REQUEST-942-APPLICATION-ATTACK-SQLI"),
			),
		},
	})
}

func TestAccApplicationGateway_webApplicationFirewall_exclusions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.webApplicationFirewall_exclusions_many(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("WAF_v2"),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("WAF_v2"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("1"),
				check.That(data.ResourceName).Key("waf_configuration.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("waf_configuration.0.firewall_mode").HasValue("Detection"),
				check.That(data.ResourceName).Key("waf_configuration.0.rule_set_type").HasValue("OWASP"),
				check.That(data.ResourceName).Key("waf_configuration.0.rule_set_version").HasValue("3.0"),
				check.That(data.ResourceName).Key("waf_configuration.0.request_body_check").HasValue("true"),
				check.That(data.ResourceName).Key("waf_configuration.0.max_request_body_size_kb").HasValue("128"),
				check.That(data.ResourceName).Key("waf_configuration.0.file_upload_limit_mb").HasValue("750"),
				check.That(data.ResourceName).Key("waf_configuration.0.exclusion.0.match_variable").HasValue("RequestArgNames"),
				check.That(data.ResourceName).Key("waf_configuration.0.exclusion.0.selector_match_operator").HasValue("Equals"),
				check.That(data.ResourceName).Key("waf_configuration.0.exclusion.0.selector").HasValue("displayNameHtml"),
				check.That(data.ResourceName).Key("waf_configuration.0.exclusion.1.match_variable").HasValue("RequestCookieNames"),
				check.That(data.ResourceName).Key("waf_configuration.0.exclusion.1.selector_match_operator").HasValue("EndsWith"),
				check.That(data.ResourceName).Key("waf_configuration.0.exclusion.1.selector").HasValue("username"),
				check.That(data.ResourceName).Key("waf_configuration.0.exclusion.2.match_variable").HasValue("RequestHeaderNames"),
				check.That(data.ResourceName).Key("waf_configuration.0.exclusion.2.selector_match_operator").HasValue("StartsWith"),
				check.That(data.ResourceName).Key("waf_configuration.0.exclusion.2.selector").HasValue("ORIGIN"),
				check.That(data.ResourceName).Key("waf_configuration.0.exclusion.3.match_variable").HasValue("RequestHeaderNames"),
				check.That(data.ResourceName).Key("waf_configuration.0.exclusion.3.selector_match_operator").HasValue("Contains"),
				check.That(data.ResourceName).Key("waf_configuration.0.exclusion.3.selector").HasValue("ORIGIN"),
				check.That(data.ResourceName).Key("waf_configuration.0.exclusion.4.match_variable").HasValue("RequestHeaderNames"),
				check.That(data.ResourceName).Key("waf_configuration.0.exclusion.4.selector_match_operator").HasValue(""),
				check.That(data.ResourceName).Key("waf_configuration.0.exclusion.4.selector").HasValue(""),
			),
		},
		{
			Config: r.webApplicationFirewall_exclusions_one(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("WAF_v2"),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("WAF_v2"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("1"),
				check.That(data.ResourceName).Key("waf_configuration.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("waf_configuration.0.firewall_mode").HasValue("Detection"),
				check.That(data.ResourceName).Key("waf_configuration.0.rule_set_type").HasValue("OWASP"),
				check.That(data.ResourceName).Key("waf_configuration.0.rule_set_version").HasValue("3.0"),
				check.That(data.ResourceName).Key("waf_configuration.0.request_body_check").HasValue("true"),
				check.That(data.ResourceName).Key("waf_configuration.0.max_request_body_size_kb").HasValue("128"),
				check.That(data.ResourceName).Key("waf_configuration.0.file_upload_limit_mb").HasValue("750"),
				check.That(data.ResourceName).Key("waf_configuration.0.exclusion.0.match_variable").HasValue("RequestArgNames"),
				check.That(data.ResourceName).Key("waf_configuration.0.exclusion.0.selector_match_operator").HasValue("Equals"),
				check.That(data.ResourceName).Key("waf_configuration.0.exclusion.0.selector").HasValue("displayNameHtml"),
			),
		},
	})
}

func TestAccApplicationGateway_sslPolicy_policyType_predefined(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sslPolicy_policyType_predefined(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ssl_policy.0.policy_type").HasValue("Predefined"),
				check.That(data.ResourceName).Key("ssl_policy.0.policy_name").HasValue("AppGwSslPolicy20170401S"),
			),
		},
	})
}

func TestAccApplicationGateway_sslPolicy_policyType_custom(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sslPolicy_policyType_custom(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ssl_policy.0.policy_type").HasValue("Custom"),
				check.That(data.ResourceName).Key("ssl_policy.0.min_protocol_version").HasValue("TLSv1_1"),
				check.That(data.ResourceName).Key("ssl_policy.0.cipher_suites.0").HasValue("TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256"),
				check.That(data.ResourceName).Key("ssl_policy.0.cipher_suites.1").HasValue("TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384"),
				check.That(data.ResourceName).Key("ssl_policy.0.cipher_suites.2").HasValue("TLS_RSA_WITH_AES_128_GCM_SHA256"),
			),
		},
	})
}

func TestAccApplicationGateway_sslPolicy_policyType_customV2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sslPolicy_policyType_customV2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ssl_policy.0.policy_type").HasValue("CustomV2"),
				check.That(data.ResourceName).Key("ssl_policy.0.min_protocol_version").HasValue("TLSv1_3"),
			),
		},
	})
}

func TestAccApplicationGateway_sslPolicy_disabledProtocols(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sslPolicy_disabledProtocols(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ssl_policy.0.disabled_protocols.0").HasValue("TLSv1_0"),
				check.That(data.ResourceName).Key("ssl_policy.0.disabled_protocols.1").HasValue("TLSv1_1"),
			),
		},
	})
}

func TestAccApplicationGateway_cookieAffinity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cookieAffinity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backend_http_settings.0.affinity_cookie_name").HasValue("testCookieName"),
			),
		},
		data.ImportStep(),
		{
			Config: r.cookieAffinityUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backend_http_settings.0.affinity_cookie_name").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_gatewayIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.gatewayIPUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_UserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.UserDefinedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
			),
		},
	})
}

func TestAccApplicationGateway_V2SKUCapacity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.v2SKUCapacity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_v2"),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("Standard_v2"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("124"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_IncludePathWithTargetURL(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.includePathWithTargetURL(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_backendAddressPoolEmptyIpList(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backendAddressPoolEmptyIpList(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_sslProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sslProfile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ssl_profile.0.verify_client_cert_issuer_dn").HasValue("false"),
				check.That(data.ResourceName).Key("ssl_profile.0.trusted_client_certificate_names.0").DoesNotExist(),
				check.That(data.ResourceName).Key("http_listener.0.ssl_profile_name").Exists(),
				check.That(data.ResourceName).Key("ssl_profile.0.ssl_policy.0.policy_type").HasValue("Predefined"),
				check.That(data.ResourceName).Key("ssl_profile.0.ssl_policy.0.policy_name").HasValue("AppGwSslPolicy20150501"),
			),
		},
		// since these are read from the existing state
		data.ImportStep(
			"ssl_certificate.0.data",
			"ssl_certificate.0.password",
		),
		{
			Config: r.sslProfileUpdateOne(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ssl_profile.0.verify_client_cert_issuer_dn").HasValue("true"),
				check.That(data.ResourceName).Key("ssl_profile.0.trusted_client_certificate_names.0").DoesNotExist(),
				check.That(data.ResourceName).Key("http_listener.0.ssl_profile_name").Exists(),
				check.That(data.ResourceName).Key("ssl_profile.0.ssl_policy.0.policy_type").HasValue("Predefined"),
				check.That(data.ResourceName).Key("ssl_profile.0.ssl_policy.0.policy_name").HasValue("AppGwSslPolicy20170401"),
			),
		},
		// since these are read from the existing state
		data.ImportStep(
			"ssl_certificate.0.data",
			"ssl_certificate.0.password",
		),
	})
}

func TestAccApplicationGateway_sslProfileWithClientCertificateVerification(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sslProfileWithClientCertificateVerification(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ssl_profile.0.verify_client_cert_issuer_dn").HasValue("false"),
				check.That(data.ResourceName).Key("ssl_profile.0.trusted_client_certificate_names.0").Exists(),
				check.That(data.ResourceName).Key("http_listener.0.ssl_profile_name").Exists(),
				check.That(data.ResourceName).Key("ssl_profile.0.ssl_policy.0.policy_type").HasValue("Custom"),
				check.That(data.ResourceName).Key("ssl_profile.0.ssl_policy.0.min_protocol_version").HasValue("TLSv1_1"),
				check.That(data.ResourceName).Key("ssl_profile.0.ssl_policy.0.cipher_suites.0").HasValue("TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256"),
				check.That(data.ResourceName).Key("ssl_profile.0.ssl_policy.0.cipher_suites.1").HasValue("TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384"),
				check.That(data.ResourceName).Key("ssl_profile.0.ssl_policy.0.cipher_suites.2").HasValue("TLS_RSA_WITH_AES_128_GCM_SHA256"),
			),
		},
		// since these are read from the existing state
		data.ImportStep(
			"ssl_certificate.0.data",
			"ssl_certificate.0.password",
		),
		{
			Config: r.sslProfileWithClientCertificateVerificationUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ssl_profile.0.verify_client_cert_issuer_dn").HasValue("true"),
				check.That(data.ResourceName).Key("ssl_profile.0.trusted_client_certificate_names.0").Exists(),
				check.That(data.ResourceName).Key("http_listener.0.ssl_profile_name").Exists(),
				check.That(data.ResourceName).Key("ssl_profile.0.ssl_policy.0.disabled_protocols.0").HasValue("TLSv1_0"),
				check.That(data.ResourceName).Key("ssl_profile.0.ssl_policy.0.disabled_protocols.1").HasValue("TLSv1_1"),
			),
		},
		// since these are read from the existing state
		data.ImportStep(
			"ssl_certificate.0.data",
			"ssl_certificate.0.password",
		),
		{
			Config: r.sslProfileWithClientCertificateVerificationUpdateTwo(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ssl_profile.0.verify_client_cert_issuer_dn").HasValue("false"),
				check.That(data.ResourceName).Key("ssl_profile.0.trusted_client_certificate_names.0").Exists(),
				check.That(data.ResourceName).Key("http_listener.0.ssl_profile_name").Exists(),
				check.That(data.ResourceName).Key("ssl_profile.0.ssl_policy.0").DoesNotExist(),
			),
		},
		// since these are read from the existing state
		data.ImportStep(
			"ssl_certificate.0.data",
			"ssl_certificate.0.password",
		),
	})
}

func TestAccApplicationGateway_requestRoutingRulePriority(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.requestRoutingRulePrioritySet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("request_routing_rule.0.priority").HasValue("1"),
				check.That(data.ResourceName).Key("request_routing_rule.1.priority").HasValue("20000"),
			),
		},
	})
}

func TestAccApplicationGateway_privateLink(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic_v2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_v2"),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("Standard_v2"),
				check.That(data.ResourceName).Key("private_link_configuration.#").HasValue("0"),
				check.That(data.ResourceName).Key("frontend_ip_configuration.1.private_link_configuration_name").DoesNotExist(),
				check.That(data.ResourceName).Key("frontend_ip_configuration.1.private_link_configuration_id").DoesNotExist(),
			),
		},
		{
			Config: r.privateLink(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("frontend_ip_configuration.1.private_link_configuration_name").HasValue("private_link"),
				check.That(data.ResourceName).Key("frontend_ip_configuration.1.private_link_configuration_id").IsSet(),
				check.That(data.ResourceName).Key("private_link_configuration.0.name").HasValue("private_link"),
				check.That(data.ResourceName).Key("private_link_configuration.0.id").IsSet(),
				check.That(data.ResourceName).Key("private_link_configuration.0.ip_configuration.0.name").HasValue("primary"),
			),
		},
		{
			Config: r.basic_v2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_v2"),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("Standard_v2"),
				check.That(data.ResourceName).Key("private_link_configuration.#").HasValue("0"),
				check.That(data.ResourceName).Key("frontend_ip_configuration.1.private_link_configuration_name").DoesNotExist(),
				check.That(data.ResourceName).Key("frontend_ip_configuration.1.private_link_configuration_id").DoesNotExist(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_updateForceFirewallPolicyAssociation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.updateForceFirewallPolicyAssociation(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateForceFirewallPolicyAssociation(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_updateEnableFips(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.updateEnableFips(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateEnableFips(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_updateFeipConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("frontend_ip_configuration.0.public_ip_address_id").IsSet(),
			),
		},
		{
			Config: r.updateFeipConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("frontend_ip_configuration.0.public_ip_address_id").IsEmpty(),
				check.That(data.ResourceName).Key("frontend_ip_configuration.0.subnet_id").IsSet(),
				check.That(data.ResourceName).Key("frontend_ip_configuration.0.private_ip_address_allocation").HasValue("Static"),
				check.That(data.ResourceName).Key("frontend_ip_configuration.0.private_ip_address").HasValue("10.0.0.10"),
				check.That(data.ResourceName).Key("frontend_ip_configuration.1.public_ip_address_id").IsSet(),
			),
		},
		data.ImportStep(),
		{
			Config: r.deletePublicFeip(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("frontend_ip_configuration.0.public_ip_address_id").IsEmpty(),
				check.That(data.ResourceName).Key("frontend_ip_configuration.0.subnet_id").IsSet(),
				check.That(data.ResourceName).Key("frontend_ip_configuration.0.private_ip_address_allocation").HasValue("Static"),
				check.That(data.ResourceName).Key("frontend_ip_configuration.0.private_ip_address").HasValue("10.0.0.10"),
				check.That(data.ResourceName).Key("frontend_ip_configuration.1").DoesNotExist(),
			),
		},
	})
}

func TestAccApplicationGateway_withoutRequestTimeout(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withoutRequestTimeout(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationGateway_removeFirewallPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	r := ApplicationGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic_wafv2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic_v2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t ApplicationGatewayResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ApplicationGatewayID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.ApplicationGatewaysClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Application Gateway (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r ApplicationGatewayResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) basic_wafv2(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-standard-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_web_application_firewall_policy" "test" {
  name                = "acctest-fwp-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  managed_rules {
    managed_rule_set {
      type    = "OWASP"
      version = "3.1"
    }
  }
  policy_settings {
    enabled                     = true
    mode                        = "Prevention"
    request_body_check          = true
    file_upload_limit_in_mb     = 250
    max_request_body_size_in_kb = 128
  }
  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  firewall_policy_id  = azurerm_web_application_firewall_policy.test.id
  sku {
    name     = "WAF_v2"
    tier     = "WAF_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) basic_v2(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-standard-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) createGlobalConfiguration(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-standard-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  global {
    request_buffering_enabled  = true
    response_buffering_enabled = true
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) updateGlobalConfiguration(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-standard-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  global {
    request_buffering_enabled  = false
    response_buffering_enabled = false
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) UserDefinedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "acctest%s"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 1
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) zones(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
  zones               = ["1", "2", "3"]
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  zones               = ["1", "2"]

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) autoscaleConfiguration(data acceptance.TestData, minCapacity int, maxCapacity int) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name = "Standard_v2"
    tier = "Standard_v2"
  }

  autoscale_configuration {
    min_capacity = %d
    max_capacity = %d
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, minCapacity, maxCapacity)
}

func (r ApplicationGatewayResource) autoscaleConfigurationNoMaxCapacity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name = "Standard_v2"
    tier = "Standard_v2"
  }

  autoscale_configuration {
    min_capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) overridePath(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    path                  = "/path1/"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) http2(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  enable_http2        = true

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_application_gateway" "import" {
  name                = azurerm_application_gateway.test.name
  resource_group_name = azurerm_application_gateway.test.resource_group_name
  location            = azurerm_application_gateway.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.basic(data))
}

func (r ApplicationGatewayResource) authCertificate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  auth_cert_name                 = "${azurerm_virtual_network.test.name}-auth"
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 443
    protocol              = "Https"
    request_timeout       = 1

    authentication_certificate {
      name = local.auth_cert_name
    }
  }

  authentication_certificate {
    name = local.auth_cert_name
    data = file("testdata/application_gateway_test.cer")
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger)
}

// nolint unused - mistakenly marked as unused
func (r ApplicationGatewayResource) trustedRootCertificate_keyvault(data acceptance.TestData) string {
	out := fmt.Sprintf(`
%[1]s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  auth_cert_name                 = "${azurerm_virtual_network.test.name}-auth"
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

data "azurerm_client_config" "test" {}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  name = "acctest%[2]d"
}

resource "azurerm_public_ip" "testStd" {
  name                = "acctest-PubIpStd-%[2]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_key_vault" "test" {
  name                = "acct%[2]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.test.tenant_id}"
  sku_name            = "standard"

  access_policy {
    tenant_id               = "${data.azurerm_client_config.test.tenant_id}"
    object_id               = "${data.azurerm_client_config.test.object_id}"
    secret_permissions      = ["Delete", "Get", "Set"]
    certificate_permissions = ["Create", "Delete", "Get", "Import", "Purge"]
  }

  access_policy {
    tenant_id               = "${data.azurerm_client_config.test.tenant_id}"
    object_id               = "${azurerm_user_assigned_identity.test.principal_id}"
    secret_permissions      = ["Get"]
    certificate_permissions = ["Get"]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctest%[2]d"
  key_vault_id = "${azurerm_key_vault.test.id}"

  certificate {
    contents = filebase64("testdata/app_service_certificate.pfx")
    password = "terraform"
  }

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = false
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }
  }
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%[2]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "WAF_v2"
    tier     = "WAF_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  waf_configuration {
    enabled                  = true
    firewall_mode            = "Detection"
    rule_set_type            = "OWASP"
    rule_set_version         = "3.0"
    file_upload_limit_mb     = 100
    request_body_check       = true
    max_request_body_size_kb = 100
  }

  identity {
    type         = "UserAssigned"
    identity_ids = ["${azurerm_user_assigned_identity.test.id}"]
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.testStd.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 443
    protocol              = "Https"
    request_timeout       = 1
  }

  trusted_root_certificate {
    name                = "${local.auth_cert_name}"
    key_vault_secret_id = "${azurerm_key_vault_certificate.test.secret_id}"
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger)
	return out
}

func (r ApplicationGatewayResource) update_trustedRootCertificate_keyvault(data acceptance.TestData) string {
	out := fmt.Sprintf(`
%[1]s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  auth_cert_name                 = "${azurerm_virtual_network.test.name}-auth"
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

data "azurerm_client_config" "test" {}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctest%[2]d"
}

resource "azurerm_public_ip" "testStd" {
  name                = "acctest-PubIpStd-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_key_vault" "test" {
  name                = "acct%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.test.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id               = data.azurerm_client_config.test.tenant_id
    object_id               = data.azurerm_client_config.test.object_id
    secret_permissions      = ["Delete", "Get", "Set"]
    certificate_permissions = ["Create", "Delete", "Get", "Import", "Purge"]
  }

  access_policy {
    tenant_id               = data.azurerm_client_config.test.tenant_id
    object_id               = azurerm_user_assigned_identity.test.principal_id
    secret_permissions      = ["Get"]
    certificate_permissions = ["Get"]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctest%[2]d"
  key_vault_id = azurerm_key_vault.test.id

  certificate {
    contents = filebase64("testdata/app_service_certificate.pfx")
    password = "terraform"
  }

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = false
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }
  }
}

resource "azurerm_key_vault_certificate" "test2" {
  name         = "acctestkvc%[2]d"
  key_vault_id = azurerm_key_vault.test.id

  certificate {
    contents = filebase64("testdata/app_service_certificate.pfx")
    password = "terraform"
  }

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = false
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }
  }
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "WAF_v2"
    tier     = "WAF_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  waf_configuration {
    enabled                  = true
    firewall_mode            = "Detection"
    rule_set_type            = "OWASP"
    rule_set_version         = "3.0"
    file_upload_limit_mb     = 100
    request_body_check       = true
    max_request_body_size_kb = 100
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.testStd.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 443
    protocol              = "Https"
    request_timeout       = 1
  }

  trusted_root_certificate {
    name                = local.auth_cert_name
    key_vault_secret_id = azurerm_key_vault_certificate.test2.secret_id
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger)
	return out
}

func (r ApplicationGatewayResource) updateForceFirewallPolicyAssociation(data acceptance.TestData, forceFirewallPolicyAssociation bool) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                              = "acctestag-%d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = azurerm_resource_group.test.location
  force_firewall_policy_association = %t

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger, forceFirewallPolicyAssociation)
}

func (r ApplicationGatewayResource) trustedRootCertificate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  auth_cert_name                 = "${azurerm_virtual_network.test.name}-auth"
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_public_ip" "teststd" {
  name                = "acctest-PubIpStd-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "WAF_v2"
    tier     = "WAF_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  waf_configuration {
    enabled                  = true
    firewall_mode            = "Detection"
    rule_set_type            = "OWASP"
    rule_set_version         = "3.0"
    file_upload_limit_mb     = 100
    request_body_check       = true
    max_request_body_size_kb = 100
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.teststd.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 443
    protocol              = "Https"
    request_timeout       = 1

    pick_host_name_from_backend_address = true
    trusted_root_certificate_names      = [local.auth_cert_name]
  }

  trusted_root_certificate {
    name = local.auth_cert_name
    data = file("testdata/application_gateway_test.cer")
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) customFirewallPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_public_ip" "teststd" {
  name                = "acctest-PubIpStd-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_web_application_firewall_policy" "testfwp" {
  name                = "acctest-fwp-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  policy_settings {
    enabled                     = true
    mode                        = "Prevention"
    file_upload_limit_in_mb     = 100
    max_request_body_size_in_kb = 128
    request_body_check          = "true"
  }

  managed_rules {
    managed_rule_set {
      type    = "OWASP"
      version = "3.1"
    }
  }
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "WAF_v2"
    tier     = "WAF_v2"
    capacity = 2
  }

  firewall_policy_id = azurerm_web_application_firewall_policy.testfwp.id

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.teststd.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 443
    protocol              = "Https"
    request_timeout       = 1

    pick_host_name_from_backend_address = true
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) customHttpListenerFirewallPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_public_ip" "teststd" {
  name                = "acctest-PubIpStd-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_web_application_firewall_policy" "testfwp" {
  name                = "acctest-fwp-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  policy_settings {
    enabled = true
    mode    = "Prevention"
  }

  managed_rules {
    managed_rule_set {
      type    = "OWASP"
      version = "3.1"
    }
  }
}

resource "azurerm_web_application_firewall_policy" "testfwp_listener" {
  name                = "acctest-fwp-listener-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  policy_settings {
    enabled = true
    mode    = "Prevention"
  }

  managed_rules {
    managed_rule_set {
      type    = "OWASP"
      version = "3.1"
    }
  }
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "WAF_v2"
    tier     = "WAF_v2"
    capacity = 2
  }

  firewall_policy_id = azurerm_web_application_firewall_policy.testfwp.id

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.teststd.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 443
    protocol              = "Https"
    request_timeout       = 1

    pick_host_name_from_backend_address = true
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
    firewall_policy_id             = azurerm_web_application_firewall_policy.testfwp_listener.id
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) customPathRuleFirewallPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  path_rule_name                 = "${azurerm_virtual_network.test.name}-pathrule1"
  url_path_map_name              = "${azurerm_virtual_network.test.name}-urlpath1"
}

resource "azurerm_public_ip" "teststd" {
  name                = "acctest-PubIpStd-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_web_application_firewall_policy" "testfwp" {
  name                = "acctest-fwp-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  policy_settings {
    enabled = true
    mode    = "Prevention"
  }

  managed_rules {
    managed_rule_set {
      type    = "OWASP"
      version = "3.1"
    }
  }
}

resource "azurerm_web_application_firewall_policy" "testfwp_path_rule" {
  name                = "acctest-fwp-path-rule-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  policy_settings {
    enabled = true
    mode    = "Prevention"
  }

  managed_rules {
    managed_rule_set {
      type    = "OWASP"
      version = "3.1"
    }
  }
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "WAF_v2"
    tier     = "WAF_v2"
    capacity = 2
  }

  firewall_policy_id = azurerm_web_application_firewall_policy.testfwp.id

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.teststd.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 443
    protocol              = "Https"
    request_timeout       = 1

    pick_host_name_from_backend_address = true
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name               = local.request_routing_rule_name
    rule_type          = "PathBasedRouting"
    url_path_map_name  = local.url_path_map_name
    http_listener_name = local.listener_name
    priority           = 10
  }

  url_path_map {
    name                               = local.url_path_map_name
    default_backend_address_pool_name  = local.backend_address_pool_name
    default_backend_http_settings_name = local.http_setting_name

    path_rule {
      name                       = local.path_rule_name
      backend_address_pool_name  = local.backend_address_pool_name
      backend_http_settings_name = local.http_setting_name
      firewall_policy_id         = azurerm_web_application_firewall_policy.testfwp_path_rule.id
      paths = [
        "/test",
      ]
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) authCertificateUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  auth_cert_name                 = "${azurerm_virtual_network.test.name}-auth2"
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 443
    protocol              = "Https"
    request_timeout       = 1

    authentication_certificate {
      name = local.auth_cert_name
    }
  }

  authentication_certificate {
    name = local.auth_cert_name
    data = file("testdata/application_gateway_test_2.crt")
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) trustedRootCertificateUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  auth_cert_name                 = "${azurerm_virtual_network.test.name}-auth2"
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_public_ip" "teststd" {
  name                = "acctest-PubIpStd-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "WAF_v2"
    tier     = "WAF_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.teststd.id
  }

  waf_configuration {
    enabled                  = true
    firewall_mode            = "Detection"
    rule_set_type            = "OWASP"
    rule_set_version         = "3.0"
    file_upload_limit_mb     = 100
    request_body_check       = true
    max_request_body_size_kb = 100
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 443
    protocol              = "Https"
    request_timeout       = 1

    pick_host_name_from_backend_address = true
    trusted_root_certificate_names      = [local.auth_cert_name]
  }

  trusted_root_certificate {
    name = local.auth_cert_name
    data = file("testdata/application_gateway_test_2.crt")
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) pathBasedRouting(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  path_rule_name                 = "${azurerm_virtual_network.test.name}-pathrule1"
  url_path_map_name              = "${azurerm_virtual_network.test.name}-urlpath1"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name               = local.request_routing_rule_name
    rule_type          = "PathBasedRouting"
    url_path_map_name  = local.url_path_map_name
    http_listener_name = local.listener_name
  }

  url_path_map {
    name                               = local.url_path_map_name
    default_backend_address_pool_name  = local.backend_address_pool_name
    default_backend_http_settings_name = local.http_setting_name

    path_rule {
      name                       = local.path_rule_name
      backend_address_pool_name  = local.backend_address_pool_name
      backend_http_settings_name = local.http_setting_name

      paths = [
        "/test",
      ]
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) pathBasedRoutingwithRedirection(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  path_rule_name                 = "${azurerm_virtual_network.test.name}-pathrule1"
  url_path_map_name              = "${azurerm_virtual_network.test.name}-urlpath1"
  redirect_configuration_name    = "${azurerm_virtual_network.test.name}-redirect"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name               = local.request_routing_rule_name
    rule_type          = "PathBasedRouting"
    url_path_map_name  = local.url_path_map_name
    http_listener_name = local.listener_name
  }

  redirect_configuration {
    name                 = local.redirect_configuration_name
    redirect_type        = "Permanent"
    target_url           = "http://example.com"
    include_path         = true
    include_query_string = false
  }

  url_path_map {
    name                                = local.url_path_map_name
    default_redirect_configuration_name = local.redirect_configuration_name

    path_rule {
      name                       = local.path_rule_name
      backend_address_pool_name  = local.backend_address_pool_name
      backend_http_settings_name = local.http_setting_name

      paths = [
        "/test",
      ]
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) routingRedirect_httpListener(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_port_name2            = "${azurerm_virtual_network.test.name}-feport2"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  target_listener_name           = "${azurerm_virtual_network.test.name}-trgthttplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  path_rule_name                 = "${azurerm_virtual_network.test.name}-pathrule1"
  url_path_map_name              = "${azurerm_virtual_network.test.name}-urlpath1"
  redirect_configuration_name    = "${azurerm_virtual_network.test.name}-Port80To8888Redirect"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_port {
    name = local.frontend_port_name2
    port = 8888
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
    host_name                      = "application.test.com"
    require_sni                    = false
  }

  http_listener {
    name                           = local.target_listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name2
    protocol                       = "Http"
  }

  request_routing_rule {
    name                        = local.request_routing_rule_name
    rule_type                   = "Basic"
    http_listener_name          = local.listener_name
    redirect_configuration_name = local.redirect_configuration_name
  }

  redirect_configuration {
    name                 = local.redirect_configuration_name
    redirect_type        = "Temporary"
    target_listener_name = local.target_listener_name
    include_path         = true
    include_query_string = false
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) routingRedirect_httpListenerError(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_port_name2            = "${azurerm_virtual_network.test.name}-feport2"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  target_listener_name           = "${azurerm_virtual_network.test.name}-trgthttplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  path_rule_name                 = "${azurerm_virtual_network.test.name}-pathrule1"
  url_path_map_name              = "${azurerm_virtual_network.test.name}-urlpath1"
  redirect_configuration_name    = "${azurerm_virtual_network.test.name}-Port80To8888Redirect"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_port {
    name = local.frontend_port_name2
    port = 8888
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  http_listener {
    name                           = local.target_listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name2
    protocol                       = "Http"
  }

  request_routing_rule {
    name                        = local.request_routing_rule_name
    rule_type                   = "Basic"
    http_listener_name          = local.listener_name
    backend_address_pool_name   = local.backend_address_pool_name
    redirect_configuration_name = local.redirect_configuration_name
  }

  redirect_configuration {
    name                 = local.redirect_configuration_name
    redirect_type        = "Temporary"
    target_listener_name = local.target_listener_name
    include_path         = true
    include_query_string = false
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) routingRedirect_pathBased(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_port_name2            = "${azurerm_virtual_network.test.name}-feport2"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  target_listener_name           = "${azurerm_virtual_network.test.name}-trgthttplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  path_rule_name                 = "${azurerm_virtual_network.test.name}-pathrule1"
  path_rule_name2                = "${azurerm_virtual_network.test.name}-pathrule2"
  url_path_map_name              = "${azurerm_virtual_network.test.name}-urlpath1"
  redirect_configuration_name    = "${azurerm_virtual_network.test.name}-PathRedirect"
  redirect_configuration_name2   = "${azurerm_virtual_network.test.name}-PathRedirect2"
  target_url                     = "http://www.example.com"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_port {
    name = local.frontend_port_name2
    port = 8888
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  http_listener {
    name                           = local.target_listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name2
    protocol                       = "Http"
  }

  request_routing_rule {
    name               = local.request_routing_rule_name
    rule_type          = "PathBasedRouting"
    url_path_map_name  = local.url_path_map_name
    http_listener_name = local.listener_name
  }

  url_path_map {
    name                               = local.url_path_map_name
    default_backend_address_pool_name  = local.backend_address_pool_name
    default_backend_http_settings_name = local.http_setting_name

    path_rule {
      name                        = local.path_rule_name
      redirect_configuration_name = local.redirect_configuration_name

      paths = [
        "/test",
      ]
    }

    path_rule {
      name                        = local.path_rule_name2
      redirect_configuration_name = local.redirect_configuration_name2

      paths = [
        "/test2",
      ]
    }
  }

  redirect_configuration {
    name                 = local.redirect_configuration_name
    redirect_type        = "Found"
    target_url           = local.target_url
    include_query_string = true
  }

  redirect_configuration {
    name                 = local.redirect_configuration_name2
    redirect_type        = "Permanent"
    target_listener_name = local.target_listener_name
    include_path         = false
    include_query_string = false
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) probes(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  probe1_name                    = "${azurerm_virtual_network.test.name}-probe1"
  probe2_name                    = "${azurerm_virtual_network.test.name}-probe2"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    probe_name            = local.probe1_name
    protocol              = "Http"
    request_timeout       = 1
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}-2"
    cookie_based_affinity = "Disabled"
    port                  = 8080
    probe_name            = local.probe2_name
    protocol              = "Http"
    request_timeout       = 1
  }

  probe {
    name                = local.probe1_name
    protocol            = "Http"
    path                = "/test"
    host                = "azure.com"
    timeout             = 120
    interval            = 300
    unhealthy_threshold = 8
    minimum_servers     = 0
  }

  probe {
    name                = local.probe2_name
    protocol            = "Http"
    path                = "/other"
    host                = "azure.com"
    timeout             = 120
    interval            = 300
    unhealthy_threshold = 8
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) probesEmptyMatch(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  probe1_name                    = "${azurerm_virtual_network.test.name}-probe1"
  probe2_name                    = "${azurerm_virtual_network.test.name}-probe2"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    probe_name            = local.probe1_name
    protocol              = "Http"
    request_timeout       = 1
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}-2"
    cookie_based_affinity = "Disabled"
    port                  = 8080
    probe_name            = local.probe2_name
    protocol              = "Http"
    request_timeout       = 1
  }

  probe {
    name                = local.probe1_name
    protocol            = "Http"
    path                = "/test"
    host                = "azure.com"
    timeout             = 120
    interval            = 300
    unhealthy_threshold = 8
  }

  probe {
    name                = local.probe2_name
    protocol            = "Http"
    path                = "/other"
    host                = "azure.com"
    timeout             = 120
    interval            = 300
    unhealthy_threshold = 8
    match {
      body        = ""
      status_code = ["404-500"]
    }
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) probesPickHostNameFromBackendHTTPSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  probe_name                     = "${azurerm_virtual_network.test.name}-probe"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                                = local.http_setting_name
    cookie_based_affinity               = "Disabled"
    pick_host_name_from_backend_address = true
    port                                = 80
    probe_name                          = local.probe_name
    protocol                            = "Http"
    request_timeout                     = 1
  }

  probe {
    name                                      = local.probe_name
    protocol                                  = "Http"
    path                                      = "/test"
    timeout                                   = 120
    interval                                  = 300
    unhealthy_threshold                       = 8
    pick_host_name_from_backend_http_settings = true
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) probesWithPort(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  probe1_name                    = "${azurerm_virtual_network.test.name}-probe1"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    probe_name            = local.probe1_name
    protocol              = "Http"
    request_timeout       = 1
  }

  probe {
    name                = local.probe1_name
    protocol            = "Http"
    port                = "8082"
    path                = "/test"
    host                = "azure.com"
    timeout             = 120
    interval            = 300
    unhealthy_threshold = 8
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) backendHttpSettingsHostName(data acceptance.TestData, hostName string, pick bool) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                                = local.http_setting_name
    cookie_based_affinity               = "Disabled"
    host_name                           = "%s"
    port                                = 80
    protocol                            = "Http"
    request_timeout                     = 1
    pick_host_name_from_backend_address = %t
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger, hostName, pick)
}

func (r ApplicationGatewayResource) withHttpListenerHostNames(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_public_ip" "test2" {
  name                = "acctest-pubip2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  domain_name_label   = "testdns-123"
}

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test2.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
    host_names                     = ["testdns-123"]
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) settingsPickHostNameFromBackendAddress(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                                = local.http_setting_name
    cookie_based_affinity               = "Disabled"
    pick_host_name_from_backend_address = true
    port                                = 80
    protocol                            = "Http"
    request_timeout                     = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) sslCertificate_keyvault_versionless(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  auth_cert_name                 = "${azurerm_virtual_network.test.name}-auth"
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  ssl_certificate_name           = "${azurerm_virtual_network.test.name}-sslcert"
}

data "azurerm_client_config" "test" {}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "acctest%[2]d"
}

resource "azurerm_public_ip" "testStd" {
  name                = "acctest-PubIpStd-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_key_vault" "test" {
  name                = "acct%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.test.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id               = data.azurerm_client_config.test.tenant_id
    object_id               = data.azurerm_client_config.test.object_id
    secret_permissions      = ["Delete", "Get", "Set"]
    certificate_permissions = ["Create", "Delete", "Get", "Import", "Purge"]
  }

  access_policy {
    tenant_id               = data.azurerm_client_config.test.tenant_id
    object_id               = azurerm_user_assigned_identity.test.principal_id
    secret_permissions      = ["Get"]
    certificate_permissions = ["Get"]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctest%[2]d"
  key_vault_id = azurerm_key_vault.test.id

  certificate {
    contents = filebase64("testdata/app_service_certificate.pfx")
    password = "terraform"
  }

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = false
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }
  }
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "WAF_v2"
    tier     = "WAF_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  waf_configuration {
    enabled                  = true
    firewall_mode            = "Detection"
    rule_set_type            = "OWASP"
    rule_set_version         = "3.0"
    file_upload_limit_mb     = 100
    request_body_check       = true
    max_request_body_size_kb = 100
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  frontend_port {
    name = local.frontend_port_name
    port = 443
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.testStd.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Https"
    ssl_certificate_name           = local.ssl_certificate_name
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }

  ssl_certificate {
    name                = local.ssl_certificate_name
    key_vault_secret_id = "${azurerm_key_vault.test.vault_uri}secrets/${azurerm_key_vault_certificate.test.name}"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) sslCertificate_keyvault_versioned(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  auth_cert_name                 = "${azurerm_virtual_network.test.name}-auth"
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  ssl_certificate_name           = "${azurerm_virtual_network.test.name}-sslcert"
}

data "azurerm_client_config" "test" {}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "acctest%[2]d"
}

resource "azurerm_public_ip" "testStd" {
  name                = "acctest-PubIpStd-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_key_vault" "test" {
  name                = "acct%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.test.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id               = data.azurerm_client_config.test.tenant_id
    object_id               = data.azurerm_client_config.test.object_id
    secret_permissions      = ["Delete", "Get", "Set"]
    certificate_permissions = ["Create", "Delete", "Get", "Import", "Purge"]
  }

  access_policy {
    tenant_id               = data.azurerm_client_config.test.tenant_id
    object_id               = azurerm_user_assigned_identity.test.principal_id
    secret_permissions      = ["Get"]
    certificate_permissions = ["Get"]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctest%[2]d"
  key_vault_id = azurerm_key_vault.test.id

  certificate {
    contents = filebase64("testdata/app_service_certificate.pfx")
    password = "terraform"
  }

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = false
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }
  }
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "WAF_v2"
    tier     = "WAF_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  waf_configuration {
    enabled                  = true
    firewall_mode            = "Detection"
    rule_set_type            = "OWASP"
    rule_set_version         = "3.0"
    file_upload_limit_mb     = 100
    request_body_check       = true
    max_request_body_size_kb = 100
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  frontend_port {
    name = local.frontend_port_name
    port = 443
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.testStd.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Https"
    ssl_certificate_name           = local.ssl_certificate_name
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }

  ssl_certificate {
    name                = local.ssl_certificate_name
    key_vault_secret_id = azurerm_key_vault_certificate.test.secret_id
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) sslCertificate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  ssl_certificate_name           = "${azurerm_virtual_network.test.name}-ssl1"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 443
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Https"
    ssl_certificate_name           = local.ssl_certificate_name
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }

  ssl_certificate {
    name     = local.ssl_certificate_name
    data     = filebase64("testdata/application_gateway_test.pfx")
    password = "terraform"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) sslCertificateUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  ssl_certificate_name           = "${azurerm_virtual_network.test.name}-ssl2"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 443
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Https"
    ssl_certificate_name           = local.ssl_certificate_name
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }

  ssl_certificate {
    name     = local.ssl_certificate_name
    data     = filebase64("testdata/application_gateway_test_2.pfx")
    password = "hello-world"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) sslCertificateEmptyPassword(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  ssl_certificate_name           = "${azurerm_virtual_network.test.name}-ssl3"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 443
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Https"
    ssl_certificate_name           = local.ssl_certificate_name
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }

  ssl_certificate {
    name     = local.ssl_certificate_name
    data     = filebase64("testdata/application_gateway_test_3.pfx")
    password = ""
  }
}
`, r.template(data), data.RandomInteger)
}

func (ApplicationGatewayResource) changeCert(certificateName string) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		// ctx has to refresh timeout value, or it may cause timeout when `WaitForCompletionRef`
		// `time.Minute*90` is the timeout value for ApplicationGatewayResource Create
		ctx, cancel := context.WithTimeout(ctx, time.Minute*90)
		defer cancel()

		gatewayName := state.Attributes["name"]
		resourceGroup := state.Attributes["resource_group_name"]

		agw, err := clients.Network.ApplicationGatewaysClient.Get(ctx, resourceGroup, gatewayName)
		if err != nil {
			return fmt.Errorf("Bad: Get on ApplicationGatewaysClient: %+v", err)
		}

		certPfx, err := os.ReadFile("testdata/application_gateway_test.pfx")
		if err != nil {
			log.Fatal(err)
		}
		certB64 := base64.StdEncoding.EncodeToString(certPfx)

		newSslCertificates := make([]network.ApplicationGatewaySslCertificate, 1)
		newSslCertificates[0] = network.ApplicationGatewaySslCertificate{
			Name: utils.String(certificateName),
			Etag: utils.String("*"),

			ApplicationGatewaySslCertificatePropertiesFormat: &network.ApplicationGatewaySslCertificatePropertiesFormat{
				Data:     utils.String(certB64),
				Password: utils.String("terraform"),
			},
		}

		agw.SslCertificates = &newSslCertificates

		future, err := clients.Network.ApplicationGatewaysClient.CreateOrUpdate(ctx, resourceGroup, gatewayName, agw)
		if err != nil {
			return fmt.Errorf("Bad: updating AGW: %+v", err)
		}

		if err := future.WaitForCompletionRef(ctx, clients.Network.ApplicationGatewaysClient.Client); err != nil {
			return fmt.Errorf("Bad: waiting for update of AGW: %+v", err)
		}

		return nil
	}
}

func (r ApplicationGatewayResource) manualSslCertificateChangeIgnoreChangesConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  ssl_certificate_name           = "acctestcertificate1"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctesttag"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }

  ssl_certificate {
    name     = local.ssl_certificate_name
    data     = filebase64("testdata/application_gateway_test.pfx")
    password = "terraform"
  }

  lifecycle {
    ignore_changes = [
      ssl_certificate,
    ]
  }
}
`, r.template(data))
}

func (r ApplicationGatewayResource) manualSslCertificateChangeIgnoreChangesUpdatedConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  ssl_certificate_name           = "acctestcertificate3"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctesttag"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }

  ssl_certificate {
    name     = local.ssl_certificate_name
    data     = filebase64("testdata/application_gateway_test.pfx")
    password = "terraform"
  }

  lifecycle {
    ignore_changes = [
      ssl_certificate,
    ]
  }
}
`, r.template(data))
}

func (r ApplicationGatewayResource) webApplicationFirewall(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "WAF_Medium"
    tier     = "WAF"
    capacity = 1
  }

  waf_configuration {
    enabled                  = true
    firewall_mode            = "Detection"
    rule_set_type            = "OWASP"
    rule_set_version         = "3.0"
    file_upload_limit_mb     = 100
    request_body_check       = true
    max_request_body_size_kb = 100
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) connectionDraining(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  enable_http2        = true

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1

    connection_draining {
      enabled           = true
      drain_timeout_sec = 1984
    }
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) webApplicationFirewall_disabledRuleGroups(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "WAF_v2"
    tier     = "WAF_v2"
    capacity = 1
  }

  waf_configuration {
    enabled                  = true
    firewall_mode            = "Detection"
    rule_set_type            = "OWASP"
    rule_set_version         = "3.0"
    request_body_check       = true
    max_request_body_size_kb = 128
    file_upload_limit_mb     = 100

    disabled_rule_group {
      rule_group_name = "REQUEST-921-PROTOCOL-ATTACK"
      rules           = [921110, 921151, 921180]
    }

    disabled_rule_group {
      rule_group_name = "REQUEST-930-APPLICATION-ATTACK-LFI"
      rules           = [930120, 930130]
    }

    disabled_rule_group {
      rule_group_name = "REQUEST-942-APPLICATION-ATTACK-SQLI"
    }
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) webApplicationFirewall_disabledRuleGroups_enabled_some_rules(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "WAF_v2"
    tier     = "WAF_v2"
    capacity = 1
  }

  waf_configuration {
    enabled                  = true
    firewall_mode            = "Detection"
    rule_set_type            = "OWASP"
    rule_set_version         = "3.0"
    request_body_check       = true
    max_request_body_size_kb = 128
    file_upload_limit_mb     = 100

    disabled_rule_group {
      rule_group_name = "REQUEST-921-PROTOCOL-ATTACK"
      rules           = [921110, 921151, 921180]
    }

    disabled_rule_group {
      rule_group_name = "REQUEST-942-APPLICATION-ATTACK-SQLI"
    }
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) webApplicationFirewall_exclusions_many(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "WAF_v2"
    tier     = "WAF_v2"
    capacity = 1
  }

  waf_configuration {
    enabled                  = true
    firewall_mode            = "Detection"
    rule_set_type            = "OWASP"
    rule_set_version         = "3.0"
    request_body_check       = true
    max_request_body_size_kb = 128
    file_upload_limit_mb     = 750

    exclusion {
      match_variable          = "RequestArgNames"
      selector_match_operator = "Equals"
      selector                = "displayNameHtml"
    }

    exclusion {
      match_variable          = "RequestCookieNames"
      selector_match_operator = "EndsWith"
      selector                = "username"
    }

    exclusion {
      match_variable          = "RequestHeaderNames"
      selector_match_operator = "StartsWith"
      selector                = "ORIGIN"
    }

    exclusion {
      match_variable          = "RequestHeaderNames"
      selector_match_operator = "Contains"
      selector                = "ORIGIN"
    }

    exclusion {
      match_variable = "RequestHeaderNames"
    }
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) webApplicationFirewall_exclusions_one(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "WAF_v2"
    tier     = "WAF_v2"
    capacity = 1
  }

  waf_configuration {
    enabled                  = true
    firewall_mode            = "Detection"
    rule_set_type            = "OWASP"
    rule_set_version         = "3.0"
    request_body_check       = true
    max_request_body_size_kb = 128
    file_upload_limit_mb     = 750

    exclusion {
      match_variable          = "RequestArgNames"
      selector_match_operator = "Equals"
      selector                = "displayNameHtml"
    }
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) sslPolicy_policyType_predefined(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 1
  }

  ssl_policy {
    policy_name = "AppGwSslPolicy20170401S"
    policy_type = "Predefined"
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) sslPolicy_policyType_custom(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 1
  }

  ssl_policy {
    policy_type          = "Custom"
    min_protocol_version = "TLSv1_1"
    cipher_suites        = ["TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256", "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384", "TLS_RSA_WITH_AES_128_GCM_SHA256"]
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) sslPolicy_policyType_customV2(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 1
  }

  ssl_policy {
    policy_type          = "CustomV2"
    min_protocol_version = "TLSv1_3"
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) sslPolicy_disabledProtocols(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 1
  }

  ssl_policy {
    disabled_protocols = ["TLSv1_0", "TLSv1_1"]
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (ApplicationGatewayResource) template(data acceptance.TestData) string {
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
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
}

resource "azurerm_subnet" "test" {
  name                                          = "subnet-%d"
  resource_group_name                           = azurerm_resource_group.test.name
  virtual_network_name                          = azurerm_virtual_network.test.name
  address_prefixes                              = ["10.0.0.0/24"]
  enforce_private_link_service_network_policies = true
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-pubip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) customErrorConfigurations(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  path_rule_name                 = "${azurerm_virtual_network.test.name}-pathrule1"
  url_path_map_name              = "${azurerm_virtual_network.test.name}-urlpath1"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"

    custom_error_configuration {
      status_code           = "HttpStatus403"
      custom_error_page_url = "http://azure.com/error403_listener.html"
    }

    custom_error_configuration {
      status_code           = "HttpStatus502"
      custom_error_page_url = "http://azure.com/error502_listener.html"
    }
  }

  custom_error_configuration {
    status_code           = "HttpStatus403"
    custom_error_page_url = "http://azure.com/error.html"
  }

  custom_error_configuration {
    status_code           = "HttpStatus502"
    custom_error_page_url = "http://azure.com/error.html"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) rewriteRuleSets_backend(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  rewrite_rule_set_name          = "${azurerm_virtual_network.test.name}-rwset"
  rewrite_rule_name              = "${azurerm_virtual_network.test.name}-rwrule"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name  = local.backend_address_pool_name
    fqdns = ["foo.com", "bar.com"]
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    rewrite_rule_set_name      = local.rewrite_rule_set_name
    priority                   = 10
  }

  rewrite_rule_set {
    name = local.rewrite_rule_set_name

    rewrite_rule {
      name          = local.rewrite_rule_name
      rule_sequence = 1

      condition {
        variable = "var_http_status"
        pattern  = "502"
      }

      request_header_configuration {
        header_name  = "X-custom"
        header_value = "customvalue"
      }
    }
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) rewriteRuleSets_redirect(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_port_name2            = "${azurerm_virtual_network.test.name}-feport2"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  target_listener_name           = "${azurerm_virtual_network.test.name}-trgthttplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  redirect_configuration_name    = "${azurerm_virtual_network.test.name}-Port80To8888Redirect"
  rewrite_rule_set_name          = "${azurerm_virtual_network.test.name}-rwset"
  rewrite_rule_name              = "${azurerm_virtual_network.test.name}-rwrule"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_port {
    name = local.frontend_port_name2
    port = 8888
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  http_listener {
    name                           = local.target_listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name2
    protocol                       = "Http"
  }

  request_routing_rule {
    name                        = local.request_routing_rule_name
    rule_type                   = "Basic"
    http_listener_name          = local.listener_name
    redirect_configuration_name = local.redirect_configuration_name
    rewrite_rule_set_name       = local.rewrite_rule_set_name
    priority                    = 10
  }

  rewrite_rule_set {
    name = local.rewrite_rule_set_name

    rewrite_rule {
      name          = local.rewrite_rule_name
      rule_sequence = 1

      condition {
        variable = "var_http_status"
        pattern  = "502"
      }

      request_header_configuration {
        header_name  = "X-custom"
        header_value = "customvalue"
      }
    }
  }

  redirect_configuration {
    name                 = local.redirect_configuration_name
    redirect_type        = "Temporary"
    target_listener_name = local.target_listener_name
    include_path         = true
    include_query_string = false
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) rewriteRuleSets_rewriteUrl(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_port_name2            = "${azurerm_virtual_network.test.name}-feport2"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  target_listener_name           = "${azurerm_virtual_network.test.name}-trgthttplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  redirect_configuration_name    = "${azurerm_virtual_network.test.name}-Port80To8888Redirect"
  rewrite_rule_set_name          = "${azurerm_virtual_network.test.name}-rwset"
  rewrite_rule_name              = "${azurerm_virtual_network.test.name}-rwrule"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_port {
    name = local.frontend_port_name2
    port = 8888
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  http_listener {
    name                           = local.target_listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name2
    protocol                       = "Http"
  }

  request_routing_rule {
    name                        = local.request_routing_rule_name
    rule_type                   = "Basic"
    http_listener_name          = local.listener_name
    redirect_configuration_name = local.redirect_configuration_name
    rewrite_rule_set_name       = local.rewrite_rule_set_name
    priority                    = 10
  }

  rewrite_rule_set {
    name = local.rewrite_rule_set_name

    rewrite_rule {
      name          = local.rewrite_rule_name
      rule_sequence = 1

      condition {
        variable    = "var_uri_path"
        pattern     = ".*article/(.*)/(.*)"
        ignore_case = false
        negate      = false
      }
      response_header_configuration {
        header_name  = "X-custom"
        header_value = "customvalue"
      }

      url {
        path         = "/article.aspx"
        query_string = "id={var_uri_path_1}&title={var_uri_path_2}"
        reroute      = false
      }
    }
  }

  rewrite_rule_set {
    name = "${local.rewrite_rule_set_name}_1"

    rewrite_rule {
      name          = "${local.rewrite_rule_name}_1"
      rule_sequence = 1

      condition {
        variable = "var_uri_path"
        pattern  = ".*article/(.*)/(.*)"
      }

      url {
        path       = "/article.aspx"
        components = "path_only"
      }
    }

    rewrite_rule {
      name          = "${local.rewrite_rule_name}_2"
      rule_sequence = 2

      condition {
        variable = "var_uri_path"
        pattern  = ".*article2/(.*)/(.*)"
      }

      url {
        query_string = "id={var_uri_path_1}&title={var_uri_path_2}"
        components   = "query_string_only"
      }
    }
  }

  redirect_configuration {
    name                 = local.redirect_configuration_name
    redirect_type        = "Temporary"
    target_listener_name = local.target_listener_name
    include_path         = true
    include_query_string = false
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) rewriteRuleSets_rewriteUrlUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_port_name2            = "${azurerm_virtual_network.test.name}-feport2"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  target_listener_name           = "${azurerm_virtual_network.test.name}-trgthttplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  redirect_configuration_name    = "${azurerm_virtual_network.test.name}-Port80To8888Redirect"
  rewrite_rule_set_name          = "${azurerm_virtual_network.test.name}-rwset"
  rewrite_rule_name              = "${azurerm_virtual_network.test.name}-rwrule"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_port {
    name = local.frontend_port_name2
    port = 8888
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  http_listener {
    name                           = local.target_listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name2
    protocol                       = "Http"
  }

  request_routing_rule {
    name                        = local.request_routing_rule_name
    rule_type                   = "Basic"
    http_listener_name          = local.listener_name
    redirect_configuration_name = local.redirect_configuration_name
    rewrite_rule_set_name       = local.rewrite_rule_set_name
    priority                    = 10
  }

  rewrite_rule_set {
    name = local.rewrite_rule_set_name

    rewrite_rule {
      name          = local.rewrite_rule_name
      rule_sequence = 1

      condition {
        variable    = "var_uri_path"
        pattern     = ".*article/(.*)/(.*)"
        ignore_case = false
        negate      = false
      }
      response_header_configuration {
        header_name  = "X-custom"
        header_value = "customvalue"
      }

      url {
        path         = "/article.aspx"
        query_string = "id={var_uri_path_1}&title={var_uri_path_2}"
        reroute      = false
      }
    }
  }

  rewrite_rule_set {
    name = "${local.rewrite_rule_set_name}_1"

    rewrite_rule {
      name          = "${local.rewrite_rule_name}_1"
      rule_sequence = 1

      condition {
        variable = "var_uri_path"
        pattern  = ".*article/(.*)/(.*)"
      }

      url {
        path = "/article.aspx"
      }
    }

    rewrite_rule {
      name          = "${local.rewrite_rule_name}_2"
      rule_sequence = 2

      condition {
        variable = "var_uri_path"
        pattern  = ".*article2/(.*)/(.*)"
      }

      url {
        query_string = "id={var_uri_path_1}&title={var_uri_path_2}"
        components   = "query_string_only"
      }
    }
  }

  redirect_configuration {
    name                 = local.redirect_configuration_name
    redirect_type        = "Temporary"
    target_listener_name = local.target_listener_name
    include_path         = true
    include_query_string = false
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) cookieAffinity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Enabled"
    affinity_cookie_name  = "testCookieName"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) cookieAffinityUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Enabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) gatewayIPUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test1" {
  name                 = "subnet1-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test1.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) v2SKUCapacity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 124
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) includePathWithTargetURL(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_port_name2            = "${azurerm_virtual_network.test.name}-feport2"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  target_listener_name           = "${azurerm_virtual_network.test.name}-trgthttplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  path_rule_name                 = "${azurerm_virtual_network.test.name}-pathrule1"
  path_rule_name2                = "${azurerm_virtual_network.test.name}-pathrule2"
  url_path_map_name              = "${azurerm_virtual_network.test.name}-urlpath1"
  redirect_configuration_name    = "${azurerm_virtual_network.test.name}-PathRedirect"
  redirect_configuration_name2   = "${azurerm_virtual_network.test.name}-PathRedirect2"
  target_url                     = "http://www.example.com"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_port {
    name = local.frontend_port_name2
    port = 8888
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  http_listener {
    name                           = local.target_listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name2
    protocol                       = "Http"
  }

  request_routing_rule {
    name               = local.request_routing_rule_name
    rule_type          = "PathBasedRouting"
    url_path_map_name  = local.url_path_map_name
    http_listener_name = local.listener_name
  }

  url_path_map {
    name                               = local.url_path_map_name
    default_backend_address_pool_name  = local.backend_address_pool_name
    default_backend_http_settings_name = local.http_setting_name

    path_rule {
      name                        = local.path_rule_name
      redirect_configuration_name = local.redirect_configuration_name

      paths = [
        "/test",
      ]
    }

    path_rule {
      name                        = local.path_rule_name2
      redirect_configuration_name = local.redirect_configuration_name2

      paths = [
        "/test2",
      ]
    }
  }

  redirect_configuration {
    name                 = local.redirect_configuration_name
    redirect_type        = "Found"
    target_url           = local.target_url
    include_query_string = true
    include_path         = true
  }

  redirect_configuration {
    name                 = local.redirect_configuration_name2
    redirect_type        = "Permanent"
    target_listener_name = local.target_listener_name
    include_path         = false
    include_query_string = false
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) backendAddressPoolEmptyIpList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name         = local.backend_address_pool_name
    ip_addresses = []
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) sslProfile(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  ssl_profile_name               = "${azurerm_virtual_network.test.name}-sslprof"
  ssl_certificate_name           = "${azurerm_virtual_network.test.name}-ssl1"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 443
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Https"
    ssl_certificate_name           = local.ssl_certificate_name
    ssl_profile_name               = local.ssl_profile_name
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }

  ssl_profile {
    name = local.ssl_profile_name
    ssl_policy {
      policy_type = "Predefined"
      policy_name = "AppGwSslPolicy20150501"
    }
  }

  ssl_certificate {
    name     = local.ssl_certificate_name
    data     = filebase64("testdata/application_gateway_test.pfx")
    password = "terraform"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) sslProfileUpdateOne(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  ssl_profile_name               = "${azurerm_virtual_network.test.name}-sslprof"
  ssl_certificate_name           = "${azurerm_virtual_network.test.name}-ssl1"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 443
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Https"
    ssl_certificate_name           = local.ssl_certificate_name
    ssl_profile_name               = local.ssl_profile_name
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }

  ssl_profile {
    name                                 = local.ssl_profile_name
    verify_client_cert_issuer_dn         = true
    verify_client_certificate_revocation = "OCSP"
    ssl_policy {
      policy_type = "Predefined"
      policy_name = "AppGwSslPolicy20170401"
    }
  }

  ssl_certificate {
    name     = local.ssl_certificate_name
    data     = filebase64("testdata/application_gateway_test.pfx")
    password = "terraform"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) sslProfileWithClientCertificateVerification(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  ssl_profile_name               = "${azurerm_virtual_network.test.name}-sslprof"
  trusted_client_cert_name       = "${azurerm_virtual_network.test.name}-tcc"
  ssl_certificate_name           = "${azurerm_virtual_network.test.name}-ssl1"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 443
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Https"
    ssl_certificate_name           = local.ssl_certificate_name
    ssl_profile_name               = local.ssl_profile_name
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }

  ssl_profile {
    name                             = local.ssl_profile_name
    trusted_client_certificate_names = [local.trusted_client_cert_name]
    ssl_policy {
      policy_type          = "Custom"
      min_protocol_version = "TLSv1_1"
      cipher_suites        = ["TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256", "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384", "TLS_RSA_WITH_AES_128_GCM_SHA256"]
    }
  }

  trusted_client_certificate {
    name = local.trusted_client_cert_name
    data = filebase64("testdata/application_gateway_mtls_test.cer")
  }

  ssl_certificate {
    name     = local.ssl_certificate_name
    data     = filebase64("testdata/application_gateway_test.pfx")
    password = "terraform"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) sslProfileWithClientCertificateVerificationUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  ssl_profile_name               = "${azurerm_virtual_network.test.name}-sslprof"
  trusted_client_cert_name       = "${azurerm_virtual_network.test.name}-tcc"
  ssl_certificate_name           = "${azurerm_virtual_network.test.name}-ssl1"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 443
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Https"
    ssl_certificate_name           = local.ssl_certificate_name
    ssl_profile_name               = local.ssl_profile_name
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }

  ssl_profile {
    name                             = local.ssl_profile_name
    trusted_client_certificate_names = [local.trusted_client_cert_name]
    verify_client_cert_issuer_dn     = true
    ssl_policy {
      disabled_protocols = ["TLSv1_0", "TLSv1_1"]
    }
  }

  trusted_client_certificate {
    name = local.trusted_client_cert_name
    data = filebase64("testdata/application_gateway_mtls_test_2.cer")
  }

  ssl_certificate {
    name     = local.ssl_certificate_name
    data     = filebase64("testdata/application_gateway_test.pfx")
    password = "terraform"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) sslProfileWithClientCertificateVerificationUpdateTwo(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  ssl_profile_name               = "${azurerm_virtual_network.test.name}-sslprof"
  trusted_client_cert_name       = "${azurerm_virtual_network.test.name}-tcc"
  ssl_certificate_name           = "${azurerm_virtual_network.test.name}-ssl1"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 443
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Https"
    ssl_certificate_name           = local.ssl_certificate_name
    ssl_profile_name               = local.ssl_profile_name
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }

  ssl_profile {
    name                             = local.ssl_profile_name
    trusted_client_certificate_names = [local.trusted_client_cert_name]
  }

  trusted_client_certificate {
    name = local.trusted_client_cert_name
    data = filebase64("testdata/application_gateway_mtls_test_2.cer")
  }

  ssl_certificate {
    name     = local.ssl_certificate_name
    data     = filebase64("testdata/application_gateway_test.pfx")
    password = "terraform"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) requestRoutingRulePrioritySet(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name_1           = "${azurerm_virtual_network.test.name}-feport-1"
  frontend_port_name_2           = "${azurerm_virtual_network.test.name}-feport-2"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name_1                = "${azurerm_virtual_network.test.name}-httplstn-1"
  listener_name_2                = "${azurerm_virtual_network.test.name}-httplstn-2"
  request_routing_rule_name_1    = "${azurerm_virtual_network.test.name}-rqrt-1"
  request_routing_rule_name_2    = "${azurerm_virtual_network.test.name}-rqrt-2"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 1
  }

  ssl_policy {
    policy_type          = "Custom"
    min_protocol_version = "TLSv1_1"
    cipher_suites        = ["TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256", "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384", "TLS_RSA_WITH_AES_128_GCM_SHA256"]
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name_1
    port = 80
  }

  frontend_port {
    name = local.frontend_port_name_2
    port = 8080
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name_1
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name_1
    protocol                       = "Http"
  }

  http_listener {
    name                           = local.listener_name_2
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name_2
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name_1
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name_1
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 1
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name_2
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name_2
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 20000
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) privateLink(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name               = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name                      = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name          = "${azurerm_virtual_network.test.name}-feip"
  frontend_ip_configuration_internal_name = "${azurerm_virtual_network.test.name}-feipint"
  http_setting_name                       = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                           = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name               = "${azurerm_virtual_network.test.name}-rqrt"
  private_link_configuration_name         = "private_link"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-standard-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test_standard.id
  }

  frontend_ip_configuration {
    name                            = local.frontend_ip_configuration_internal_name
    subnet_id                       = azurerm_subnet.test.id
    private_ip_address_allocation   = "Static"
    private_ip_address              = "10.0.0.10"
    private_link_configuration_name = local.private_link_configuration_name
  }

  private_link_configuration {
    name = local.private_link_configuration_name
    ip_configuration {
      name                          = "primary"
      subnet_id                     = azurerm_subnet.test.id
      private_ip_address_allocation = "Dynamic"
      primary                       = true
    }
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
    priority                   = 10
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApplicationGatewayResource) updateEnableFips(data acceptance.TestData, enableFips bool) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  fips_enabled        = %t

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger, enableFips)
}

func (r ApplicationGatewayResource) updateFeipConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name          = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name                 = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name     = "${azurerm_virtual_network.test.name}-feip"
  frontend_ip_configuration_name_new = "${azurerm_virtual_network.test.name}-feip-new"
  http_setting_name                  = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                      = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name          = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                          = local.frontend_ip_configuration_name_new
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Static"
    private_ip_address            = "10.0.0.10"
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name_new
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) deletePublicFeip(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name          = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name                 = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name     = "${azurerm_virtual_network.test.name}-feip"
  frontend_ip_configuration_name_new = "${azurerm_virtual_network.test.name}-feip-new"
  http_setting_name                  = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                      = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name          = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                          = local.frontend_ip_configuration_name_new
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Static"
    private_ip_address            = "10.0.0.10"
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name_new
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationGatewayResource) withoutRequestTimeout(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.test.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.test.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }
}
`, r.template(data), data.RandomInteger)
}
