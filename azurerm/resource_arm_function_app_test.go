package azurerm

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMFunctionApp_basic(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMFunctionApp_basic(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					testCheckAzureRMFunctionAppHasNoContentShare(resourceName),
					resource.TestCheckResourceAttr(resourceName, "version", "~1"),
					resource.TestCheckResourceAttrSet(resourceName, "outbound_ip_addresses"),
					resource.TestCheckResourceAttrSet(resourceName, "possible_outbound_ip_addresses"),
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

func TestAccAzureRMFunctionApp_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_function_app.test"
	ri := tf.AccRandTimeInt()

	location := testLocation()
	rs := strings.ToLower(acctest.RandString(11))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					testCheckAzureRMFunctionAppHasNoContentShare(resourceName),
				),
			},
			{
				Config:      testAccAzureRMFunctionApp_requiresImport(ri, rs, location),
				ExpectError: testRequiresImportError("azurerm_function_app"),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_tags(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMFunctionApp_tags(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "production"),
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

func TestAccAzureRMFunctionApp_tagsUpdate(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMFunctionApp_tags(ri, rs, testLocation())
	updatedConfig := testAccAzureRMFunctionApp_tagsUpdated(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "production"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "production"),
					resource.TestCheckResourceAttr(resourceName, "tags.hello", "Berlin"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_appSettings(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMFunctionApp_basic(ri, rs, testLocation())
	updatedConfig := testAccAzureRMFunctionApp_appSettings(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "app_settings.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "site_credential.#", "1"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "app_settings.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_settings.hello", "world"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "app_settings.%", "0"),
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

func TestAccAzureRMFunctionApp_siteConfig(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMFunctionApp_alwaysOn(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.always_on", "true"),
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

func TestAccAzureRMFunctionApp_linuxFxVersion(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMFunctionApp_linuxFxVersion(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "kind", "functionapp,linux,container"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.linux_fx_version", "DOCKER|(golang:latest)"),
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

func TestAccAzureRMFunctionApp_connectionStrings(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMFunctionApp_connectionStrings(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "connection_string.0.name", "Example"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.0.value", "some-postgresql-connection-string"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.0.type", "PostgreSQL"),
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

func TestAccAzureRMFunctionApp_siteConfigMulti(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	configBase := testAccAzureRMFunctionApp_basic(ri, rs, testLocation())
	configUpdate1 := testAccAzureRMFunctionApp_appSettings(ri, rs, testLocation())
	configUpdate2 := testAccAzureRMFunctionApp_appSettingsAlwaysOn(ri, rs, testLocation())
	configUpdate3 := testAccAzureRMFunctionApp_appSettingsAlwaysOnLinuxFxVersion(ri, rs, testLocation())
	configUpdate4 := testAccAzureRMFunctionApp_appSettingsAlwaysOnLinuxFxVersionConnectionStrings(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBase,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "app_settings.%", "0"),
				),
			},
			{
				Config: configUpdate1,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "app_settings.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_settings.hello", "world"),
				),
			},
			{
				Config: configUpdate2,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "app_settings.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_settings.hello", "world"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.always_on", "true"),
				),
			},
			{
				Config: configUpdate3,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "kind", "functionapp,linux,container"),
					resource.TestCheckResourceAttr(resourceName, "app_settings.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_settings.hello", "world"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.always_on", "true"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.linux_fx_version", "DOCKER|(golang:latest)"),
				),
			},
			{
				Config: configUpdate4,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "kind", "functionapp,linux,container"),
					resource.TestCheckResourceAttr(resourceName, "app_settings.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_settings.hello", "world"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.always_on", "true"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.linux_fx_version", "DOCKER|(golang:latest)"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.0.name", "Example"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.0.value", "some-postgresql-connection-string"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.0.type", "PostgreSQL"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_updateVersion(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	preConfig := testAccAzureRMFunctionApp_version(ri, rs, testLocation(), "~1")
	postConfig := testAccAzureRMFunctionApp_version(ri, rs, testLocation(), "~2")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "version", "~1"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "version", "~2"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_3264bit(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()
	config := testAccAzureRMFunctionApp_basic(ri, rs, location)
	updatedConfig := testAccAzureRMFunctionApp_64bit(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.use_32_bit_worker_process", "true"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.use_32_bit_worker_process", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_httpsOnly(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()
	config := testAccAzureRMFunctionApp_httpsOnly(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "https_only", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_consumptionPlan(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()
	config := testAccAzureRMFunctionApp_consumptionPlan(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					testCheckAzureRMFunctionAppHasContentShare(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.use_32_bit_worker_process", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_consumptionPlanUppercaseName(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()
	config := testAccAzureRMFunctionApp_consumptionPlanUppercaseName(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					testCheckAzureRMFunctionAppHasContentShare(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.use_32_bit_worker_process", "true"),
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

func TestAccAzureRMFunctionApp_createIdentity(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMFunctionApp_basicIdentity(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(resourceName, "identity.0.principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(resourceName, "identity.0.tenant_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_updateIdentity(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))

	preConfig := testAccAzureRMFunctionApp_basic(ri, rs, testLocation())
	postConfig := testAccAzureRMFunctionApp_basicIdentity(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity.#", "0"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(resourceName, "identity.0.principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(resourceName, "identity.0.tenant_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_loggingDisabled(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMFunctionApp_loggingDisabled(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					testCheckAzureRMFunctionAppHasNoContentShare(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_builtin_logging", "false"),
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

func TestAccAzureRMFunctionApp_updateLogging(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	enabledConfig := testAccAzureRMFunctionApp_basic(ri, rs, location)
	disabledConfig := testAccAzureRMFunctionApp_loggingDisabled(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: enabledConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_builtin_logging", "true"),
				),
			},
			{
				Config: disabledConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_builtin_logging", "false"),
				),
			},
			{
				Config: enabledConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_builtin_logging", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_authSettings(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	tenantID := os.Getenv("ARM_TENANT_ID")
	config := testAccAzureRMFunctionApp_authSettings(ri, rs, testLocation(), tenantID)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.runtime_version", "1.0"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.unauthenticated_client_action", "RedirectToLoginPage"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.token_refresh_extension_hours", "75"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.token_store_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.additional_login_params.test_key", "test_value"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.allowed_external_redirect_urls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.allowed_external_redirect_urls.0", "https://terra.form"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(resourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
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

func TestAccAzureRMFunctionApp_corsSettings(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMFunctionApp_corsSettings(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.cors.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.cors.0.support_credentials", "true"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.cors.0.allowed_origins.#", "4"),
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

func TestAccAzureRMFunctionApp_vnetName(t *testing.T) {
	resourceName := "azurerm_function_app.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	vnetName := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMFunctionApp_vnetName(ri, rs, testLocation(), vnetName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.virtual_network_name", vnetName),
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

func testCheckAzureRMFunctionAppDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).web.AppServicesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_function_app" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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

func testCheckAzureRMFunctionAppExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		functionAppName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Function App: %s", functionAppName)
		}

		client := testAccProvider.Meta().(*ArmClient).web.AppServicesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, functionAppName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Function App %q (resource group: %q) does not exist", functionAppName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on appServicesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMFunctionAppHasContentShare(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		functionAppName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Function App: %s", functionAppName)
		}

		client := testAccProvider.Meta().(*ArmClient).web.AppServicesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		appSettingsResp, err := client.ListApplicationSettings(ctx, resourceGroup, functionAppName)
		if err != nil {
			return fmt.Errorf("Error making Read request on AzureRM Function App AppSettings %q: %+v", functionAppName, err)
		}

		for k := range appSettingsResp.Properties {
			if strings.EqualFold("WEBSITE_CONTENTSHARE", k) {
				return nil
			}
		}

		return fmt.Errorf("Function App %q does not contain the Website Content Share!", functionAppName)
	}
}

func testCheckAzureRMFunctionAppHasNoContentShare(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		functionAppName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Function App: %s", functionAppName)
		}

		client := testAccProvider.Meta().(*ArmClient).web.AppServicesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		appSettingsResp, err := client.ListApplicationSettings(ctx, resourceGroup, functionAppName)
		if err != nil {
			return fmt.Errorf("Error making Read request on AzureRM Function App AppSettings %q: %+v", functionAppName, err)
		}

		for k, v := range appSettingsResp.Properties {
			if strings.EqualFold("WEBSITE_CONTENTSHARE", k) && v != nil && *v != "" {
				return fmt.Errorf("Function App %q contains the Website Content Share!", functionAppName)
			}
		}

		return nil
	}
}

func testAccAzureRMFunctionApp_basic(rInt int, storage string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
}
`, rInt, location, storage)
}

func testAccAzureRMFunctionApp_requiresImport(rInt int, storage, location string) string {
	template := testAccAzureRMFunctionApp_basic(rInt, storage, location)
	return fmt.Sprintf(`
%s

resource "azurerm_function_app" "import" {
  name                      = "${azurerm_function_app.test.name}"
  location                  = "${azurerm_function_app.test.location}"
  resource_group_name       = "${azurerm_function_app.test.resource_group_name}"
  app_service_plan_id       = "${azurerm_function_app.test.app_service_plan_id}"
  storage_connection_string = "${azurerm_function_app.test.storage_connection_string}"
}
`, template)
}

func testAccAzureRMFunctionApp_tags(rInt int, storage string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  tags = {
    environment = "production"
  }
}
`, rInt, location, storage)
}

func testAccAzureRMFunctionApp_tagsUpdated(rInt int, storage string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  tags = {
    environment = "production"
    hello       = "Berlin"
  }
}
`, rInt, location, storage)
}

func testAccAzureRMFunctionApp_version(rInt int, storage string, location string, version string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  version                   = "%[4]s"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
}
`, rInt, location, storage, version)
}

func testAccAzureRMFunctionApp_appSettings(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  app_settings = {
    "hello" = "world"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMFunctionApp_alwaysOn(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  site_config {
    always_on = true
  }
}
`, rInt, location, rString)
}

func testAccAzureRMFunctionApp_linuxFxVersion(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  kind                = "Linux"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }

  reserved = true
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  version                   = "~2"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  site_config {
    linux_fx_version = "DOCKER|(golang:latest)"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMFunctionApp_connectionStrings(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  connection_string {
    name  = "Example"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMFunctionApp_appSettingsAlwaysOn(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  app_settings = {
    "hello" = "world"
  }

  site_config {
    always_on = true
  }
}
`, rInt, location, rString)
}

func testAccAzureRMFunctionApp_appSettingsAlwaysOnLinuxFxVersion(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  kind                = "Linux"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }

  reserved = true
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  version                   = "~2"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  app_settings = {
    "hello" = "world"
  }

  site_config {
    always_on        = true
    linux_fx_version = "DOCKER|(golang:latest)"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMFunctionApp_appSettingsAlwaysOnLinuxFxVersionConnectionStrings(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  kind                = "Linux"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }

  reserved = true
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  version                   = "~2"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

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
`, rInt, location, rString)
}

func testAccAzureRMFunctionApp_64bit(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  site_config {
    use_32_bit_worker_process = false
  }
}
`, rInt, location, rString, rInt, rInt)
}

func testAccAzureRMFunctionApp_httpsOnly(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
  https_only                = true
}
`, rInt, location, rString, rInt, rInt)
}

func testAccAzureRMFunctionApp_consumptionPlan(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  kind                = "FunctionApp"

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
}
`, rInt, location, rString, rInt, rInt)
}

func testAccAzureRMFunctionApp_consumptionPlanUppercaseName(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  kind                = "FunctionApp"

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%d-FuncWithUppercase"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
}
`, rInt, location, rString, rInt, rInt)
}

func testAccAzureRMFunctionApp_basicIdentity(rInt int, storage string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  identity {
    type = "SystemAssigned"
  }
}
`, rInt, location, storage)
}

func testAccAzureRMFunctionApp_loggingDisabled(rInt int, storage string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
  enable_builtin_logging    = false
}
`, rInt, location, storage)
}

func testAccAzureRMFunctionApp_authSettings(rInt int, storage string, location string, tenantID string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

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
`, rInt, location, storage, tenantID)
}

func testAccAzureRMFunctionApp_corsSettings(rInt int, storage string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

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
`, rInt, location, storage)
}

func testAccAzureRMFunctionApp_vnetName(rInt int, storage, location, vnetName string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  site_config {
    virtual_network_name = "%[4]s"
  }
}
`, rInt, location, storage, vnetName)
}
