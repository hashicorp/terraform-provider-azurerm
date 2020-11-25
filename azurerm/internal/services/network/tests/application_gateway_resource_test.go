package tests

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApplicationGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Standard_Small"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.#", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_autoscaleConfiguration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_autoscaleConfiguration(data, 0, 10),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Standard_v2"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "Standard_v2"),
					resource.TestCheckResourceAttr(data.ResourceName, "autoscale_configuration.0.min_capacity", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "autoscale_configuration.0.max_capacity", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.#", "0"),
				),
			},
			{
				Config: testAccAzureRMApplicationGateway_autoscaleConfiguration(data, 4, 12),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Standard_v2"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "Standard_v2"),
					resource.TestCheckResourceAttr(data.ResourceName, "autoscale_configuration.0.min_capacity", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "autoscale_configuration.0.max_capacity", "12"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.#", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_autoscaleConfigurationNoMaxCapacity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_autoscaleConfigurationNoMaxCapacity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Standard_v2"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "Standard_v2"),
					resource.TestCheckResourceAttr(data.ResourceName, "autoscale_configuration.0.min_capacity", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.#", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_zones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_zones(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Standard_v2"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "Standard_v2"),
					resource.TestCheckResourceAttr(data.ResourceName, "zones.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.#", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_overridePath(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_overridePath(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "backend_http_settings.0.path", "/path1/"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_http2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_http2(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_http2", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMApplicationGateway_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_application_gateway"),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_authCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_authCertificate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "authentication_certificate.0.name"),
				),
			},
			// since these are read from the existing state
			data.ImportStep(

				"authentication_certificate.0.data",
			),
			{
				Config: testAccAzureRMApplicationGateway_authCertificateUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "authentication_certificate.0.name"),
				),
			},
			// since these are read from the existing state
			data.ImportStep(

				"authentication_certificate.0.data",
			),
		},
	})
}

func TestAccAzureRMApplicationGateway_customFirewallPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_customFirewallPolicy(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "firewall_policy_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_customHttpListenerFirewallPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_customHttpListenerFirewallPolicy(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "http_listener.0.firewall_policy_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

// TODO required soft delete on the keyvault
func TestAccAzureRMApplicationGateway_trustedRootCertificate_keyvault(t *testing.T) {
	t.Skip()

	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_trustedRootCertificate_keyvault(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "trusted_root_certificate.0.name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_trustedRootCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_trustedRootCertificate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "trusted_root_certificate.0.name"),
				),
			},
			// since these are read from the existing state
			data.ImportStep(
				"trusted_root_certificate.0.data",
			),
			{
				Config: testAccAzureRMApplicationGateway_trustedRootCertificateUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "trusted_root_certificate.0.name"),
				),
			},
			// since these are read from the existing state
			data.ImportStep(
				"trusted_root_certificate.0.data",
			),
		},
	})
}

func TestAccAzureRMApplicationGateway_pathBasedRouting(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_pathBasedRouting(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_routingRedirect_httpListener(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_routingRedirect_httpListener(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_routingRedirect_httpListenerError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMApplicationGateway_routingRedirect_httpListenerError(data),
				ExpectError: regexp.MustCompile("Conflict between `backend_address_pool_name` and `redirect_configuration_name`"),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_routingRedirect_pathBased(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_routingRedirect_pathBased(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_customErrorConfigurations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_customErrorConfigurations(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_rewriteRuleSets_backend(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_rewriteRuleSets_backend(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "rewrite_rule_set.0.name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_rewriteRuleSets_redirect(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_rewriteRuleSets_redirect(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "rewrite_rule_set.0.name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_probes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_probes(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_probesEmptyMatch(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_probesEmptyMatch(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_probesPickHostNameFromBackendHTTPSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_probesPickHostNameFromBackendHTTPSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "probe.0.pick_host_name_from_backend_http_settings", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_probesWithPort(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_probesWithPort(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "probe.0.port", "8082"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_backendHttpSettingsHostName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	hostName := "example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_backendHttpSettingsHostName(data, hostName, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "backend_http_settings.0.host_name", hostName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_withHttpListenerHostNames(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_withHttpListenerHostNames(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_backendHttpSettingsHostNameAndPick(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")
	hostName := "example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMApplicationGateway_backendHttpSettingsHostName(data, hostName, true),
				ExpectError: regexp.MustCompile("Only one of `host_name` or `pick_host_name_from_backend_address` can be set"),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_settingsPickHostNameFromBackendAddress(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_settingsPickHostNameFromBackendAddress(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "backend_http_settings.0.pick_host_name_from_backend_address", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_sslCertificate_keyvault_versionless(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_sslCertificate_keyvault_versionless(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ssl_certificate.0.key_vault_secret_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_sslCertificate_keyvault_versioned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_sslCertificate_keyvault_versioned(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ssl_certificate.0.key_vault_secret_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_sslCertificate_EmptyPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_sslCertificateEmptyPassword(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
				),
			},
			// since these are read from the existing state
			data.ImportStep(
				"ssl_certificate.0.data",
				"ssl_certificate.0.password",
			),
		},
	})
}

func TestAccAzureRMApplicationGateway_manualSslCertificateChangeIgnoreChanges(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_manualSslCertificateChangeIgnoreChangesConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_certificate.0.name", "acctestcertificate1"),
					testCheckAzureRMApplicationGatewayChangeCert(data.ResourceName, "acctestcertificate2"),
				),
			},
			{
				Config: testAccAzureRMApplicationGateway_manualSslCertificateChangeIgnoreChangesUpdatedConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_certificate.0.name", "acctestcertificate2"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_sslCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_sslCertificate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
				),
			},
			// since these are read from the existing state
			data.ImportStep(
				"ssl_certificate.0.data",
				"ssl_certificate.0.password",
			),
			{
				Config: testAccAzureRMApplicationGateway_sslCertificateUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
				),
			},
			// since these are read from the existing state
			data.ImportStep(
				"ssl_certificate.0.data",
				"ssl_certificate.0.password",
			),
		},
	})
}

func TestAccAzureRMApplicationGateway_webApplicationFirewall(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_webApplicationFirewall(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "WAF_Medium"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "WAF"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.firewall_mode", "Detection"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.rule_set_type", "OWASP"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.rule_set_version", "3.0"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.file_upload_limit_mb", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.request_body_check", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.max_request_body_size_kb", "100"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_connectionDraining(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_connectionDraining(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "backend_http_settings.0.connection_draining.0.enabled", "true"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApplicationGateway_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Standard_Small"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.#", "0"),
					resource.TestCheckNoResourceAttr(data.ResourceName, "backend_http_settings.0.connection_draining.0.enabled"),
					resource.TestCheckNoResourceAttr(data.ResourceName, "backend_http_settings.0.connection_draining.0.drain_timeout_sec"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_webApplicationFirewall_disabledRuleGroups(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_webApplicationFirewall_disabledRuleGroups(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "WAF_v2"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "WAF_v2"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.firewall_mode", "Detection"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.rule_set_type", "OWASP"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.rule_set_version", "3.0"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.request_body_check", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.max_request_body_size_kb", "128"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.file_upload_limit_mb", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.disabled_rule_group.0.rule_group_name", "REQUEST-921-PROTOCOL-ATTACK"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.disabled_rule_group.0.rules.0", "921110"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.disabled_rule_group.0.rules.1", "921151"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.disabled_rule_group.0.rules.2", "921180"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.disabled_rule_group.1.rule_group_name", "REQUEST-930-APPLICATION-ATTACK-LFI"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.disabled_rule_group.1.rules.0", "930120"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.disabled_rule_group.1.rules.1", "930130"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.disabled_rule_group.2.rule_group_name", "REQUEST-942-APPLICATION-ATTACK-SQLI"),
				),
			},
			{
				Config: testAccAzureRMApplicationGateway_webApplicationFirewall_disabledRuleGroups_enabled_some_rules(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "WAF_v2"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "WAF_v2"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.firewall_mode", "Detection"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.rule_set_type", "OWASP"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.rule_set_version", "3.0"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.request_body_check", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.max_request_body_size_kb", "128"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.file_upload_limit_mb", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.disabled_rule_group.0.rule_group_name", "REQUEST-921-PROTOCOL-ATTACK"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.disabled_rule_group.0.rules.0", "921110"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.disabled_rule_group.0.rules.1", "921151"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.disabled_rule_group.0.rules.2", "921180"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.disabled_rule_group.1.rule_group_name", "REQUEST-942-APPLICATION-ATTACK-SQLI"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_webApplicationFirewall_exclusions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_webApplicationFirewall_exclusions_many(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "WAF_v2"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "WAF_v2"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.firewall_mode", "Detection"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.rule_set_type", "OWASP"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.rule_set_version", "3.0"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.request_body_check", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.max_request_body_size_kb", "128"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.file_upload_limit_mb", "750"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.exclusion.0.match_variable", "RequestArgNames"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.exclusion.0.selector_match_operator", "Equals"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.exclusion.0.selector", "displayNameHtml"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.exclusion.1.match_variable", "RequestCookieNames"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.exclusion.1.selector_match_operator", "EndsWith"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.exclusion.1.selector", "username"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.exclusion.2.match_variable", "RequestHeaderNames"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.exclusion.2.selector_match_operator", "StartsWith"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.exclusion.2.selector", "ORIGIN"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.exclusion.3.match_variable", "RequestHeaderNames"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.exclusion.3.selector_match_operator", "Contains"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.exclusion.3.selector", "ORIGIN"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.exclusion.4.match_variable", "RequestHeaderNames"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.exclusion.4.selector_match_operator", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.exclusion.4.selector", ""),
				),
			},
			{
				Config: testAccAzureRMApplicationGateway_webApplicationFirewall_exclusions_one(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "WAF_v2"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "WAF_v2"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.firewall_mode", "Detection"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.rule_set_type", "OWASP"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.rule_set_version", "3.0"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.request_body_check", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.max_request_body_size_kb", "128"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.file_upload_limit_mb", "750"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.exclusion.0.match_variable", "RequestArgNames"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.exclusion.0.selector_match_operator", "Equals"),
					resource.TestCheckResourceAttr(data.ResourceName, "waf_configuration.0.exclusion.0.selector", "displayNameHtml"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_sslPolicy_policyType_predefined(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_sslPolicy_policyType_predefined(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_policy.0.policy_type", "Predefined"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_policy.0.policy_name", "AppGwSslPolicy20170401S"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_sslPolicy_policyType_custom(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_sslPolicy_policyType_custom(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_policy.0.policy_type", "Custom"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_policy.0.min_protocol_version", "TLSv1_1"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_policy.0.cipher_suites.0", "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_policy.0.cipher_suites.1", "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_policy.0.cipher_suites.2", "TLS_RSA_WITH_AES_128_GCM_SHA256"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_sslPolicy_disabledProtocols(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_sslPolicy_disabledProtocols(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_policy.0.disabled_protocols.0", "TLSv1_0"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_policy.0.disabled_protocols.1", "TLSv1_1"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_cookieAffinity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_cookieAffinity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "backend_http_settings.0.affinity_cookie_name", "testCookieName"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApplicationGateway_cookieAffinityUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "backend_http_settings.0.affinity_cookie_name", ""),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_gatewayIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApplicationGateway_gatewayIPUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_UserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_UserDefinedIdentity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "UserAssigned"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.identity_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_V2SKUCapacity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_v2SKUCapacity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Standard_v2"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "Standard_v2"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "124"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_IncludePathWithTargetURL(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_includePathWithTargetURL(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationGateway_backendAddressPoolEmptyIpList(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_backendAddressPoolEmptyIpList(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMApplicationGatewayExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ApplicationGatewaysClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		gatewayName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Application Gateway: %q", gatewayName)
		}

		resp, err := client.Get(ctx, resourceGroup, gatewayName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Application Gateway %q (resource group: %q) does not exist", gatewayName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on applicationGatewayClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMApplicationGatewayDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ApplicationGatewaysClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_application_gateway" {
			continue
		}

		gatewayName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, gatewayName)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Application Gateway still exists:\n%#v", resp.ApplicationGatewayPropertiesFormat)
	}

	return nil
}

func testAccAzureRMApplicationGateway_basic(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_UserDefinedIdentity(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
  }
}
`, template, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_zones(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_autoscaleConfiguration(data acceptance.TestData, minCapacity int, maxCapacity int) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
  }
}
`, template, data.RandomInteger, data.RandomInteger, minCapacity, maxCapacity)
}

func testAccAzureRMApplicationGateway_autoscaleConfigurationNoMaxCapacity(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_overridePath(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_http2(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_basic(data)
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
`, template)
}

func testAccAzureRMApplicationGateway_authCertificate(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}

// nolint unused - mistakenly marked as unused
func testAccAzureRMApplicationGateway_trustedRootCertificate_keyvault(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%[1]s

# since these variables are re-used - a locals block makes this more maintainable
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
    secret_permissions      = ["delete", "get", "set"]
    certificate_permissions = ["create", "delete", "get", "import"]
  }

  access_policy {
    tenant_id               = "${data.azurerm_client_config.test.tenant_id}"
    object_id               = "${azurerm_user_assigned_identity.test.principal_id}"
    secret_permissions      = ["get"]
    certificate_permissions = ["get"]
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

  identity {
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
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "${local.request_routing_rule_name}"
    rule_type                  = "Basic"
    http_listener_name         = "${local.listener_name}"
    backend_address_pool_name  = "${local.backend_address_pool_name}"
    backend_http_settings_name = "${local.http_setting_name}"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_trustedRootCertificate(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%[1]s

# since these variables are re-used - a locals block makes this more maintainable
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
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_customFirewallPolicy(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%[1]s

# since these variables are re-used - a locals block makes this more maintainable
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
    max_request_body_size_in_kb = 100
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
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_customHttpListenerFirewallPolicy(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%[1]s

# since these variables are re-used - a locals block makes this more maintainable
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
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_authCertificateUpdated(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%[1]s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_trustedRootCertificateUpdated(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%[1]s

# since these variables are re-used - a locals block makes this more maintainable
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
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_pathBasedRouting(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_routingRedirect_httpListener(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_routingRedirect_httpListenerError(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_routingRedirect_pathBased(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_probes(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_probesEmptyMatch(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
      body = ""
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
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_probesPickHostNameFromBackendHTTPSettings(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_probesWithPort(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_backendHttpSettingsHostName(data acceptance.TestData, hostName string, pick bool) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger, hostName, pick)
}

func testAccAzureRMApplicationGateway_withHttpListenerHostNames(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
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

# since these variables are re-used - a locals block makes this more maintainable
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
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_settingsPickHostNameFromBackendAddress(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_sslCertificate_keyvault_versionless(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
    secret_permissions      = ["delete", "get", "set"]
    certificate_permissions = ["create", "delete", "get", "import"]
  }

  access_policy {
    tenant_id               = data.azurerm_client_config.test.tenant_id
    object_id               = azurerm_user_assigned_identity.test.principal_id
    secret_permissions      = ["get"]
    certificate_permissions = ["get"]
  }

  soft_delete_enabled = true
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

  identity {
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
  }

  ssl_certificate {
    name                = local.ssl_certificate_name
    key_vault_secret_id = "${azurerm_key_vault.test.vault_uri}secrets/${azurerm_key_vault_certificate.test.name}"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_sslCertificate_keyvault_versioned(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
    secret_permissions      = ["delete", "get", "set"]
    certificate_permissions = ["create", "delete", "get", "import"]
  }

  access_policy {
    tenant_id               = data.azurerm_client_config.test.tenant_id
    object_id               = azurerm_user_assigned_identity.test.principal_id
    secret_permissions      = ["get"]
    certificate_permissions = ["get"]
  }

  soft_delete_enabled = true
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

  identity {
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
  }

  ssl_certificate {
    name                = local.ssl_certificate_name
    key_vault_secret_id = azurerm_key_vault_certificate.test.secret_id
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_sslCertificate(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_sslCertificateUpdated(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_sslCertificateEmptyPassword(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}

func testCheckAzureRMApplicationGatewayChangeCert(resourceName, certName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ApplicationGatewaysClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		gatewayName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		agw, err := client.Get(ctx, resourceGroup, gatewayName)
		if err != nil {
			return fmt.Errorf("Bad: Get on ApplicationGatewaysClient: %+v", err)
		}

		certPfx, err := ioutil.ReadFile("testdata/application_gateway_test.pfx")
		if err != nil {
			log.Fatal(err)
		}
		certB64 := base64.StdEncoding.EncodeToString(certPfx)

		newSslCertificates := make([]network.ApplicationGatewaySslCertificate, 1)
		newSslCertificates[0] = network.ApplicationGatewaySslCertificate{
			Name: utils.String(certName),
			Etag: utils.String("*"),

			ApplicationGatewaySslCertificatePropertiesFormat: &network.ApplicationGatewaySslCertificatePropertiesFormat{
				Data:     utils.String(certB64),
				Password: utils.String("terraform"),
			},
		}

		agw.SslCertificates = &newSslCertificates

		future, err := client.CreateOrUpdate(ctx, resourceGroup, gatewayName, agw)
		if err != nil {
			return fmt.Errorf("Bad: updating AGW: %+v", err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Bad: waiting for update of AGW: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApplicationGateway_manualSslCertificateChangeIgnoreChangesConfig(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template)
}

func testAccAzureRMApplicationGateway_manualSslCertificateChangeIgnoreChangesUpdatedConfig(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template)
}

func testAccAzureRMApplicationGateway_webApplicationFirewall(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_connectionDraining(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}
func testAccAzureRMApplicationGateway_webApplicationFirewall_disabledRuleGroups(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_webApplicationFirewall_disabledRuleGroups_enabled_some_rules(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_webApplicationFirewall_exclusions_many(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}
func testAccAzureRMApplicationGateway_webApplicationFirewall_exclusions_one(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_sslPolicy_policyType_predefined(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
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
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_sslPolicy_policyType_custom(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
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
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_sslPolicy_disabledProtocols(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
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
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_template(data acceptance.TestData) string {
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
  name                 = "subnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.0.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-pubip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_customErrorConfigurations(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_rewriteRuleSets_backend(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
    rewrite_rule_set_name      = local.rewrite_rule_set_name
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
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_rewriteRuleSets_redirect(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_cookieAffinity(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_cookieAffinityUpdated(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_gatewayIPUpdated(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test1" {
  name                 = "subnet1-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_v2SKUCapacity(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_includePathWithTargetURL(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}

func testAccAzureRMApplicationGateway_backendAddressPoolEmptyIpList(data acceptance.TestData) string {
	template := testAccAzureRMApplicationGateway_template(data)
	return fmt.Sprintf(`
%s

# since these variables are re-used - a locals block makes this more maintainable
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
`, template, data.RandomInteger)
}
