package azurerm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApplicationGateway_basic(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_Small"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "2"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_autoscaleConfiguration(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_autoscaleConfiguration(ri, testLocation(), 2, 10),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_v2"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "Standard_v2"),
					resource.TestCheckResourceAttr(resourceName, "autoscale_configuration.0.min_capacity", "2"),
					resource.TestCheckResourceAttr(resourceName, "autoscale_configuration.0.max_capacity", "10"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.#", "0"),
				),
			},
			{
				Config: testAccAzureRMApplicationGateway_autoscaleConfiguration(ri, testLocation(), 4, 12),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_v2"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "Standard_v2"),
					resource.TestCheckResourceAttr(resourceName, "autoscale_configuration.0.min_capacity", "4"),
					resource.TestCheckResourceAttr(resourceName, "autoscale_configuration.0.max_capacity", "12"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_autoscaleConfigurationNoMaxCapacity(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_autoscaleConfigurationNoMaxCapacity(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_v2"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "Standard_v2"),
					resource.TestCheckResourceAttr(resourceName, "autoscale_configuration.0.min_capacity", "2"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_zones(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_zones(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_v2"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "Standard_v2"),
					resource.TestCheckResourceAttr(resourceName, "zones.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "2"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_overridePath(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_overridePath(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "backend_http_settings.0.path", "/path1/"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_http2(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_http2(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_http2", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMApplicationGateway_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_application_gateway"),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_authCertificate(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_authCertificate(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "authentication_certificate.0.name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// since these are read from the existing state
					"authentication_certificate.0.data",
				},
			},
			{
				Config: testAccAzureRMApplicationGateway_authCertificateUpdated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "authentication_certificate.0.name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// since these are read from the existing state
					"authentication_certificate.0.data",
				},
			},
		},
	})
}

// TODO required soft delete on the keyvault
func TestAccAzureRMApplicationGateway_trustedRootCertificate_keyvault(t *testing.T) {
	t.Skip()

	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_trustedRootCertificate_keyvault(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "trusted_root_certificate.0.name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_trustedRootCertificate(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_trustedRootCertificate(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "trusted_root_certificate.0.name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// since these are read from the existing state
					"trusted_root_certificate.0.data",
				},
			},
			{
				Config: testAccAzureRMApplicationGateway_trustedRootCertificateUpdated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "trusted_root_certificate.0.name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// since these are read from the existing state
					"trusted_root_certificate.0.data",
				},
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_pathBasedRouting(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_pathBasedRouting(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_routingRedirect_httpListener(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_routingRedirect_httpListener(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "redirect_configuration.0.name"),
					resource.TestCheckResourceAttr(resourceName, "redirect_configuration.0.redirect_type", "Temporary"),
					resource.TestCheckResourceAttrSet(resourceName, "redirect_configuration.0.target_listener_name"),
					resource.TestCheckResourceAttr(resourceName, "redirect_configuration.0.include_path", "true"),
					resource.TestCheckResourceAttr(resourceName, "redirect_configuration.0.include_query_string", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_routingRedirect_httpListenerError(t *testing.T) {
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMApplicationGateway_routingRedirect_httpListenerError(ri, testLocation()),
				ExpectError: regexp.MustCompile("Conflict between `backend_address_pool_name` and `redirect_configuration_name`"),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_routingRedirect_pathBased(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_routingRedirect_pathBased(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "redirect_configuration.0.name"),
					resource.TestCheckResourceAttr(resourceName, "redirect_configuration.0.redirect_type", "Found"),
					resource.TestCheckResourceAttr(resourceName, "redirect_configuration.0.target_url", "http://www.example.com"),
					resource.TestCheckResourceAttr(resourceName, "redirect_configuration.0.include_query_string", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "redirect_configuration.1.name"),
					resource.TestCheckResourceAttr(resourceName, "redirect_configuration.1.redirect_type", "Permanent"),
					resource.TestCheckResourceAttrSet(resourceName, "redirect_configuration.1.target_listener_name"),
					resource.TestCheckResourceAttr(resourceName, "redirect_configuration.1.include_path", "false"),
					resource.TestCheckResourceAttr(resourceName, "redirect_configuration.1.include_query_string", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_customErrorConfigurations(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_customErrorConfigurations(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_rewriteRuleSets_backend(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_rewriteRuleSets_backend(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "rewrite_rule_set.0.name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_rewriteRuleSets_redirect(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_rewriteRuleSets_redirect(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "rewrite_rule_set.0.name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_probes(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_probes(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_probesPickHostNameFromBackendHTTPSettings(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_probesPickHostNameFromBackendHTTPSettings(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "probe.0.pick_host_name_from_backend_http_settings", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_backendHttpSettingsHostName(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()
	hostName := "example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_backendHttpSettingsHostName(ri, testLocation(), hostName, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "backend_http_settings.0.host_name", hostName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_backendHttpSettingsHostNameAndPick(t *testing.T) {
	ri := tf.AccRandTimeInt()
	hostName := "example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMApplicationGateway_backendHttpSettingsHostName(ri, testLocation(), hostName, true),
				ExpectError: regexp.MustCompile("Only one of `host_name` or `pick_host_name_from_backend_address` can be set"),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_settingsPickHostNameFromBackendAddress(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_settingsPickHostNameFromBackendAddress(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "backend_http_settings.0.pick_host_name_from_backend_address", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_sslCertificate(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_sslCertificate(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// since these are read from the existing state
					"ssl_certificate.0.data",
					"ssl_certificate.0.password",
				},
			},
			{
				Config: testAccAzureRMApplicationGateway_sslCertificateUpdated(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// since these are read from the existing state
					"ssl_certificate.0.data",
					"ssl_certificate.0.password",
				},
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_webApplicationFirewall(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_webApplicationFirewall(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "WAF_Medium"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "WAF"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.firewall_mode", "Detection"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_type", "OWASP"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_version", "3.0"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.file_upload_limit_mb", "100"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.request_body_check", "true"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.max_request_body_size_kb", "100"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_connectionDraining(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_connectionDraining(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "backend_http_settings.0.connection_draining.0.enabled", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMApplicationGateway_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_Small"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "2"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.#", "0"),
					resource.TestCheckNoResourceAttr(resourceName, "backend_http_settings.0.connection_draining.0.enabled"),
					resource.TestCheckNoResourceAttr(resourceName, "backend_http_settings.0.connection_draining.0.drain_timeout_sec"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_webApplicationFirewall_disabledRuleGroups(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_webApplicationFirewall_disabledRuleGroups(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "WAF_v2"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "WAF_v2"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.firewall_mode", "Detection"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_type", "OWASP"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_version", "3.0"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.request_body_check", "true"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.max_request_body_size_kb", "128"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.file_upload_limit_mb", "100"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.disabled_rule_group.0.rule_group_name", "REQUEST-921-PROTOCOL-ATTACK"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.disabled_rule_group.0.rules.0", "921110"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.disabled_rule_group.0.rules.1", "921151"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.disabled_rule_group.0.rules.2", "921180"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.disabled_rule_group.1.rule_group_name", "REQUEST-930-APPLICATION-ATTACK-LFI"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.disabled_rule_group.1.rules.0", "930120"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.disabled_rule_group.1.rules.1", "930130"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.disabled_rule_group.2.rule_group_name", "REQUEST-942-APPLICATION-ATTACK-SQLI"),
				),
			},
			{
				Config: testAccAzureRMApplicationGateway_webApplicationFirewall_disabledRuleGroups_enabled_some_rules(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "WAF_v2"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "WAF_v2"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.firewall_mode", "Detection"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_type", "OWASP"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_version", "3.0"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.request_body_check", "true"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.max_request_body_size_kb", "128"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.file_upload_limit_mb", "100"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.disabled_rule_group.0.rule_group_name", "REQUEST-921-PROTOCOL-ATTACK"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.disabled_rule_group.0.rules.0", "921110"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.disabled_rule_group.0.rules.1", "921151"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.disabled_rule_group.0.rules.2", "921180"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.disabled_rule_group.1.rule_group_name", "REQUEST-942-APPLICATION-ATTACK-SQLI"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_webApplicationFirewall_exclusions(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_webApplicationFirewall_exclusions_many(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "WAF_v2"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "WAF_v2"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.firewall_mode", "Detection"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_type", "OWASP"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_version", "3.0"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.request_body_check", "true"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.max_request_body_size_kb", "128"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.file_upload_limit_mb", "100"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.exclusion.0.match_variable", "RequestArgNames"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.exclusion.0.selector_match_operator", "Equals"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.exclusion.0.selector", "displayNameHtml"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.exclusion.1.match_variable", "RequestCookieNames"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.exclusion.1.selector_match_operator", "EndsWith"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.exclusion.1.selector", "username"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.exclusion.2.match_variable", "RequestHeaderNames"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.exclusion.2.selector_match_operator", "StartsWith"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.exclusion.2.selector", "ORIGIN"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.exclusion.3.match_variable", "RequestHeaderNames"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.exclusion.3.selector_match_operator", "Contains"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.exclusion.3.selector", "ORIGIN"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.exclusion.4.match_variable", "RequestHeaderNames"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.exclusion.4.selector_match_operator", ""),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.exclusion.4.selector", ""),
				),
			},
			{
				Config: testAccAzureRMApplicationGateway_webApplicationFirewall_exclusions_one(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "WAF_v2"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "WAF_v2"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.firewall_mode", "Detection"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_type", "OWASP"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.rule_set_version", "3.0"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.request_body_check", "true"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.max_request_body_size_kb", "128"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.file_upload_limit_mb", "100"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.exclusion.0.match_variable", "RequestArgNames"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.exclusion.0.selector_match_operator", "Equals"),
					resource.TestCheckResourceAttr(resourceName, "waf_configuration.0.exclusion.0.selector", "displayNameHtml"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_sslPolicy_policyType_predefined(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_sslPolicy_policyType_predefined(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "ssl_policy.0.policy_type", "Predefined"),
					resource.TestCheckResourceAttr(resourceName, "ssl_policy.0.policy_name", "AppGwSslPolicy20170401S"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_sslPolicy_policyType_custom(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_sslPolicy_policyType_custom(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "ssl_policy.0.policy_type", "Custom"),
					resource.TestCheckResourceAttr(resourceName, "ssl_policy.0.min_protocol_version", "TLSv1_1"),
					resource.TestCheckResourceAttr(resourceName, "ssl_policy.0.cipher_suites.0", "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256"),
					resource.TestCheckResourceAttr(resourceName, "ssl_policy.0.cipher_suites.1", "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384"),
					resource.TestCheckResourceAttr(resourceName, "ssl_policy.0.cipher_suites.2", "TLS_RSA_WITH_AES_128_GCM_SHA256"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_sslPolicy_disabledProtocols(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_sslPolicy_disabledProtocols(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "ssl_policy.0.disabled_protocols.0", "TLSv1_0"),
					resource.TestCheckResourceAttr(resourceName, "ssl_policy.0.disabled_protocols.1", "TLSv1_1"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_disabledSslProtocols(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_disabledSslProtocols(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "disabled_ssl_protocols.0", "TLSv1_0"),
					resource.TestCheckResourceAttr(resourceName, "disabled_ssl_protocols.1", "TLSv1_1"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_cookieAffinity(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_cookieAffinity(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "backend_http_settings.0.affinity_cookie_name", "testCookieName"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMApplicationGateway_cookieAffinityUpdated(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "backend_http_settings.0.affinity_cookie_name", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_gatewayIP(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMApplicationGateway_gatewayIPUpdated(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApplicationGateway_UserAssignedIdentity(t *testing.T) {
	resourceName := "azurerm_application_gateway.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(14)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationGateway_UserDefinedIdentity(ri, testLocation(), rs),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity.0.type", "UserAssigned"),
					resource.TestCheckResourceAttr(resourceName, "identity.0.identity_ids.#", "1"),
				),
			},
		},
	})
}

func testCheckAzureRMApplicationGatewayExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		gatewayName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Application Gateway: %q", gatewayName)
		}

		client := testAccProvider.Meta().(*ArmClient).network.ApplicationGatewaysClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
	client := testAccProvider.Meta().(*ArmClient).network.ApplicationGatewaysClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAccAzureRMApplicationGateway_basic(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
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
`, template, rInt)
}

func testAccAzureRMApplicationGateway_UserDefinedIdentity(rInt int, location string, rString string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  name = "acctest%s"
}

resource "azurerm_public_ip" "test_standard" {
  name                = "acctest-pubip-%d-standard"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 1
  }

  identity {
    identity_ids = ["${azurerm_user_assigned_identity.test.id}"]
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test_standard.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
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
`, template, rString, rInt, rInt)
}

func testAccAzureRMApplicationGateway_zones(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  zones               = ["1", "2"]

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test_standard.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
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
`, template, rInt, rInt)
}

func testAccAzureRMApplicationGateway_autoscaleConfiguration(rInt int, location string, minCapacity int, maxCapacity int) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
  }

  autoscale_configuration {
    min_capacity = %d
	max_capacity = %d
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test_standard.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
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
`, template, rInt, rInt, minCapacity, maxCapacity)
}

func testAccAzureRMApplicationGateway_autoscaleConfigurationNoMaxCapacity(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name = "Standard_v2"
    tier = "Standard_v2"
  }

  autoscale_configuration {
    min_capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test_standard.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
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
`, template, rInt, rInt)
}

func testAccAzureRMApplicationGateway_overridePath(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    path                  = "/path1/"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
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
`, template, rInt)
}

func testAccAzureRMApplicationGateway_http2(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  enable_http2        = true

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
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
`, template, rInt)
}

func testAccAzureRMApplicationGateway_requiresImport(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_application_gateway" "import" {
  name                = "${azurerm_application_gateway.test.name}"
  resource_group_name = "${azurerm_application_gateway.test.resource_group_name}"
  location            = "${azurerm_application_gateway.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
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
`, template)
}

func testAccAzureRMApplicationGateway_authCertificate(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
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

    authentication_certificate {
      name = "${local.auth_cert_name}"
    }
  }

  authentication_certificate {
    name = "${local.auth_cert_name}"
    data = "${file("testdata/application_gateway_test.cer")}"
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
`, template, rInt)
}

func testAccAzureRMApplicationGateway_trustedRootCertificate_keyvault(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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

data "azuread_service_principal" "test" {
  display_name = "Microsoft Azure App Service"
}

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
    object_id               = "${data.azurerm_client_config.test.service_principal_object_id }"
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
`, template, rInt)
}

func testAccAzureRMApplicationGateway_trustedRootCertificate(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  sku                 = "Standard" 
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

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.teststd.id}"
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
    name = "${local.auth_cert_name}"
    data = "${file("testdata/application_gateway_test.cer")}"
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
`, template, rInt)
}

func testAccAzureRMApplicationGateway_authCertificateUpdated(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
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

    authentication_certificate {
      name = "${local.auth_cert_name}"
    }
  }

  authentication_certificate {
    name = "${local.auth_cert_name}"
    data = "${file("testdata/application_gateway_test_2.crt")}"
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
`, template, rInt)
}

func testAccAzureRMApplicationGateway_trustedRootCertificateUpdated(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  sku                 = "Standard" 
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

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.teststd.id}"
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
    name = "${local.auth_cert_name}"
    data = "${file("testdata/application_gateway_test_2.crt")}"
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
`, template, rInt)
}

func testAccAzureRMApplicationGateway_pathBasedRouting(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name               = "${local.request_routing_rule_name}"
    rule_type          = "PathBasedRouting"
    url_path_map_name  = "${local.url_path_map_name}"
    http_listener_name = "${local.listener_name}"
  }

  url_path_map {
    name                               = "${local.url_path_map_name}"
    default_backend_address_pool_name  = "${local.backend_address_pool_name}"
    default_backend_http_settings_name = "${local.http_setting_name}"

    path_rule {
      name                       = "${local.path_rule_name}"
      backend_address_pool_name  = "${local.backend_address_pool_name}"
      backend_http_settings_name = "${local.http_setting_name}"

      paths = [
        "/test",
      ]
    }
  }
}
`, template, rInt)
}

func testAccAzureRMApplicationGateway_routingRedirect_httpListener(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_port {
    name = "${local.frontend_port_name2}"
    port = 8888
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Http"
  }

  http_listener {
    name                           = "${local.target_listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name2}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                        = "${local.request_routing_rule_name}"
    rule_type                   = "Basic"
    http_listener_name          = "${local.listener_name}"
    redirect_configuration_name = "${local.redirect_configuration_name}"
  }

  redirect_configuration {
    name                 = "${local.redirect_configuration_name}"
    redirect_type        = "Temporary"
    target_listener_name = "${local.target_listener_name}"
    include_path         = true
    include_query_string = false
  }
}
`, template, rInt)
}

func testAccAzureRMApplicationGateway_routingRedirect_httpListenerError(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_port {
    name = "${local.frontend_port_name2}"
    port = 8888
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Http"
  }

  http_listener {
    name                           = "${local.target_listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name2}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                        = "${local.request_routing_rule_name}"
    rule_type                   = "Basic"
    http_listener_name          = "${local.listener_name}"
    backend_address_pool_name   = "${local.backend_address_pool_name}"
    redirect_configuration_name = "${local.redirect_configuration_name}"
  }

  redirect_configuration {
    name                 = "${local.redirect_configuration_name}"
    redirect_type        = "Temporary"
    target_listener_name = "${local.target_listener_name}"
    include_path         = true
    include_query_string = false
  }
}
`, template, rInt)
}

func testAccAzureRMApplicationGateway_routingRedirect_pathBased(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_port {
    name = "${local.frontend_port_name2}"
    port = 8888
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Http"
  }

  http_listener {
    name                           = "${local.target_listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name2}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name               = "${local.request_routing_rule_name}"
    rule_type          = "PathBasedRouting"
    url_path_map_name  = "${local.url_path_map_name}"
    http_listener_name = "${local.listener_name}"
  }

  url_path_map {
    name                               = "${local.url_path_map_name}"
    default_backend_address_pool_name  = "${local.backend_address_pool_name}"
    default_backend_http_settings_name = "${local.http_setting_name}"

    path_rule {
      name                        = "${local.path_rule_name}"
      redirect_configuration_name = "${local.redirect_configuration_name}"

      paths = [
        "/test",
      ]
    }

    path_rule {
      name                        = "${local.path_rule_name2}"
      redirect_configuration_name = "${local.redirect_configuration_name2}"

      paths = [
        "/test2",
      ]
    }
  }

  redirect_configuration {
    name                 = "${local.redirect_configuration_name}"
    redirect_type        = "Found"
    target_url           = "${local.target_url}"
    include_query_string = true
  }

  redirect_configuration {
    name                 = "${local.redirect_configuration_name2}"
    redirect_type        = "Permanent"
    target_listener_name = "${local.target_listener_name}"
    include_path         = false
    include_query_string = false
  }
}
`, template, rInt)
}

func testAccAzureRMApplicationGateway_probes(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    probe_name            = "${local.probe1_name}"
    protocol              = "Http"
    request_timeout       = 1
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}-2"
    cookie_based_affinity = "Disabled"
    port                  = 8080
    probe_name            = "${local.probe2_name}"
    protocol              = "Http"
    request_timeout       = 1
  }

  probe {
    name                = "${local.probe1_name}"
    protocol            = "Http"
    path                = "/test"
    host                = "azure.com"
    timeout             = 120
    interval            = 300
    unhealthy_threshold = 8
  }

  probe {
    name                = "${local.probe2_name}"
    protocol            = "Http"
    path                = "/other"
    host                = "azure.com"
    timeout             = 120
    interval            = 300
    unhealthy_threshold = 8
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
`, template, rInt)
}

func testAccAzureRMApplicationGateway_probesPickHostNameFromBackendHTTPSettings(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                                = "${local.http_setting_name}"
    cookie_based_affinity               = "Disabled"
    pick_host_name_from_backend_address = true
    port                                = 80
    probe_name                          = "${local.probe_name}"
    protocol                            = "Http"
    request_timeout                     = 1
  }

  probe {
    name                                      = "${local.probe_name}"
    protocol                                  = "Http"
    path                                      = "/test"
    timeout                                   = 120
    interval                                  = 300
    unhealthy_threshold                       = 8
    pick_host_name_from_backend_http_settings = true
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
`, template, rInt)
}

func testAccAzureRMApplicationGateway_backendHttpSettingsHostName(rInt int, location string, hostName string, pick bool) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    host_name             = "%s"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
    pick_host_name_from_backend_address = %t
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
`, template, rInt, hostName, pick)
}

func testAccAzureRMApplicationGateway_settingsPickHostNameFromBackendAddress(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                                = "${local.http_setting_name}"
    cookie_based_affinity               = "Disabled"
    pick_host_name_from_backend_address = true
    port                                = 80
    protocol                            = "Http"
    request_timeout                     = 1
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
`, template, rInt)
}

func testAccAzureRMApplicationGateway_sslCertificate(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 443
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Https"
    ssl_certificate_name           = "${local.ssl_certificate_name}"
  }

  request_routing_rule {
    name                       = "${local.request_routing_rule_name}"
    rule_type                  = "Basic"
    http_listener_name         = "${local.listener_name}"
    backend_address_pool_name  = "${local.backend_address_pool_name}"
    backend_http_settings_name = "${local.http_setting_name}"
  }

  ssl_certificate {
    name     = "${local.ssl_certificate_name}"
    data     = "${filebase64("testdata/application_gateway_test.pfx")}"
    password = "terraform"
  }
}
`, template, rInt)
}

func testAccAzureRMApplicationGateway_sslCertificateUpdated(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 443
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Https"
    ssl_certificate_name           = "${local.ssl_certificate_name}"
  }

  request_routing_rule {
    name                       = "${local.request_routing_rule_name}"
    rule_type                  = "Basic"
    http_listener_name         = "${local.listener_name}"
    backend_address_pool_name  = "${local.backend_address_pool_name}"
    backend_http_settings_name = "${local.http_setting_name}"
  }

  ssl_certificate {
    name     = "${local.ssl_certificate_name}"
    data     = "${filebase64("testdata/application_gateway_test_2.pfx")}"
    password = "hello-world"
  }
}
`, template, rInt)
}

func testAccAzureRMApplicationGateway_webApplicationFirewall(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "WAF_Medium"
    tier     = "WAF"
    capacity = 1
  }

  disabled_ssl_protocols = [
    "TLSv1_0",
  ]

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
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
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
`, template, rInt)
}

func testAccAzureRMApplicationGateway_connectionDraining(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  enable_http2        = true

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
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
`, template, rInt)
}
func testAccAzureRMApplicationGateway_webApplicationFirewall_disabledRuleGroups(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

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
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test_standard.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
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
`, template, rInt, rInt)
}

func testAccAzureRMApplicationGateway_webApplicationFirewall_disabledRuleGroups_enabled_some_rules(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

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
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test_standard.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
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
`, template, rInt, rInt)
}

func testAccAzureRMApplicationGateway_webApplicationFirewall_exclusions_many(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

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
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test_standard.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
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
`, template, rInt, rInt)
}
func testAccAzureRMApplicationGateway_webApplicationFirewall_exclusions_one(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

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

    exclusion {
      match_variable          = "RequestArgNames"
      selector_match_operator = "Equals"
      selector                = "displayNameHtml"
    }
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test_standard.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
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
`, template, rInt, rInt)
}

func testAccAzureRMApplicationGateway_sslPolicy_policyType_predefined(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

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
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test_standard.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
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
`, template, rInt, rInt)
}

func testAccAzureRMApplicationGateway_sslPolicy_policyType_custom(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

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
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test_standard.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
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
`, template, rInt, rInt)
}

func testAccAzureRMApplicationGateway_sslPolicy_disabledProtocols(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

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
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test_standard.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
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
`, template, rInt, rInt)
}

func testAccAzureRMApplicationGateway_disabledSslProtocols(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 1
  }

  disabled_ssl_protocols = ["TLSv1_0", "TLSv1_1"]

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test_standard.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
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
`, template, rInt, rInt)
}

func testAccAzureRMApplicationGateway_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "subnet-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.0.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-pubip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Dynamic"
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMApplicationGateway_customErrorConfigurations(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
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
    name                       = "${local.request_routing_rule_name}"
    rule_type                  = "Basic"
    http_listener_name         = "${local.listener_name}"
    backend_address_pool_name  = "${local.backend_address_pool_name}"
    backend_http_settings_name = "${local.http_setting_name}"
  }
}
`, template, rInt)
}

func testAccAzureRMApplicationGateway_rewriteRuleSets_backend(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test_standard.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
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
    rewrite_rule_set_name      = "${local.rewrite_rule_set_name}"
  }

  rewrite_rule_set {
    name = "${local.rewrite_rule_set_name}"

    rewrite_rule {
      name          = "${local.rewrite_rule_name}"
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
`, template, rInt, rInt)
}

func testAccAzureRMApplicationGateway_rewriteRuleSets_redirect(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestag-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_port {
    name = "${local.frontend_port_name2}"
    port = 8888
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test_standard.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Http"
  }

  http_listener {
    name                           = "${local.target_listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name2}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                        = "${local.request_routing_rule_name}"
    rule_type                   = "Basic"
    http_listener_name          = "${local.listener_name}"
    redirect_configuration_name = "${local.redirect_configuration_name}"
    rewrite_rule_set_name       = "${local.rewrite_rule_set_name}"
  }

  rewrite_rule_set {
    name = "${local.rewrite_rule_set_name}"

    rewrite_rule {
      name          = "${local.rewrite_rule_name}"
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
    name                 = "${local.redirect_configuration_name}"
    redirect_type        = "Temporary"
    target_listener_name = "${local.target_listener_name}"
    include_path         = true
    include_query_string = false
  }
}
`, template, rInt, rInt)
}

func testAccAzureRMApplicationGateway_cookieAffinity(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Enabled"
    affinity_cookie_name  = "testCookieName"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
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
`, template, rInt)
}

func testAccAzureRMApplicationGateway_cookieAffinityUpdated(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Enabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
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
`, template, rInt)
}

func testAccAzureRMApplicationGateway_gatewayIPUpdated(rInt int, location string) string {
	template := testAccAzureRMApplicationGateway_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test1" {
  name                 = "subnet1-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
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
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.test1.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
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
`, template, rInt, rInt)
}
