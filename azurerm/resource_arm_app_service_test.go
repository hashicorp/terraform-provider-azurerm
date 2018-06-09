package azurerm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAzureRMAppServiceName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "ab",
			ErrCount: 0,
		},
		{
			Value:    "abc",
			ErrCount: 0,
		},
		{
			Value:    "webapp1",
			ErrCount: 0,
		},
		{
			Value:    "hello-world",
			ErrCount: 0,
		},
		{
			Value:    "hello_world",
			ErrCount: 1,
		},
		{
			Value:    "helloworld21!",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateAppServiceName(tc.Value, "azurerm_app_service")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the App Service Name to trigger a validation error for '%s'", tc.Value)
		}
	}
}

func TestAccAzureRMAppService_basic(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_freeTier(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_freeTier(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_sharedTier(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_sharedTier(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_32Bit(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_32Bit(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.use_32_bit_worker_process", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_http2Enabled(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_http2Enabled(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.http2_enabled", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_alwaysOn(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_alwaysOn(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.always_on", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_httpsOnly(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_httpsOnly(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "https_only", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_appSettings(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_appSettings(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "app_settings.foo", "bar"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_clientAffinityEnabled(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_clientAffinityEnabled(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "client_affinity_enabled", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_clientAffinityDisabled(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_clientAffinityDisabled(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "client_affinity_enabled", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_enableManageServiceIdentity(t *testing.T) {

	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_mangedServiceIdentity(ri, testLocation())

	uuidMatch := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(resourceName, "identity.0.principal_id", uuidMatch),
					resource.TestMatchResourceAttr(resourceName, "identity.0.tenant_id", uuidMatch),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_updateResourceByEnablingManageServiceIdentity(t *testing.T) {

	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()

	basicResourceNoManagedIdentity := testAccAzureRMAppService_basic(ri, testLocation())
	managedIdentityEnabled := testAccAzureRMAppService_mangedServiceIdentity(ri, testLocation())

	uuidMatch := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: basicResourceNoManagedIdentity,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity.#", "0"),
				),
			},
			{
				Config: managedIdentityEnabled,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(resourceName, "identity.0.principal_id", uuidMatch),
					resource.TestMatchResourceAttr(resourceName, "identity.0.tenant_id", uuidMatch),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_clientAffinityUpdate(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_clientAffinity(ri, testLocation(), true)
	updatedConfig := testAccAzureRMAppService_clientAffinity(ri, testLocation(), false)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "client_affinity_enabled", "true"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "client_affinity_enabled", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_connectionStrings(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_connectionStrings(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "connection_string.0.name", "Example"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.0.value", "some-postgresql-connection-string"),
					resource.TestCheckResourceAttr(resourceName, "connection_string.0.type", "PostgreSQL"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_oneIpRestriction(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_oneIpRestriction(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.0.ip_address", "10.10.10.10"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.0.subnet_mask", "255.255.255.255"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_manyIpRestrictions(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_manyIpRestrictions(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.0.ip_address", "10.10.10.10"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.0.subnet_mask", "255.255.255.255"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.1.ip_address", "20.20.20.0"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.1.subnet_mask", "255.255.255.0"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.2.ip_address", "30.30.0.0"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.ip_restriction.2.subnet_mask", "255.255.0.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_defaultDocuments(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_defaultDocuments(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.default_documents.0", "first.html"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.default_documents.1", "second.jsp"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.default_documents.2", "third.aspx"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_enabled(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_enabled(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_localMySql(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_localMySql(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.local_mysql_enabled", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_managedPipelineMode(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_managedPipelineMode(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.managed_pipeline_mode", "Classic"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_tagsUpdate(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_tags(ri, testLocation())
	updatedConfig := testAccAzureRMAppService_tagsUpdated(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Hello", "World"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.Hello", "World"),
					resource.TestCheckResourceAttr(resourceName, "tags.Terraform", "AcceptanceTests"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_remoteDebugging(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_remoteDebugging(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.remote_debugging_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.remote_debugging_version", "VS2015"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_windowsDotNet2(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_windowsDotNet(ri, testLocation(), "v2.0")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.dotnet_framework_version", "v2.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_windowsDotNet4(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_windowsDotNet(ri, testLocation(), "v4.0")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.dotnet_framework_version", "v4.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_windowsDotNetUpdate(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_windowsDotNet(ri, testLocation(), "v2.0")
	updatedConfig := testAccAzureRMAppService_windowsDotNet(ri, testLocation(), "v4.0")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.dotnet_framework_version", "v2.0"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.dotnet_framework_version", "v4.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_windowsJava7Jetty(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_windowsJava(ri, testLocation(), "1.7", "JETTY", "9.3")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_version", "1.7"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container", "JETTY"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container_version", "9.3"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_windowsJava8Jetty(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_windowsJava(ri, testLocation(), "1.8", "JETTY", "9.3")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_version", "1.8"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container", "JETTY"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container_version", "9.3"),
				),
			},
		},
	})
}
func TestAccAzureRMAppService_windowsJava7Tomcat(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_windowsJava(ri, testLocation(), "1.7", "TOMCAT", "9.0")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_version", "1.7"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container", "TOMCAT"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container_version", "9.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_windowsJava8Tomcat(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_windowsJava(ri, testLocation(), "1.8", "TOMCAT", "9.0")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_version", "1.8"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container", "TOMCAT"),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.java_container_version", "9.0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_windowsPHP7(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_windowsPHP(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.php_version", "7.1"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_windowsPython(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_windowsPython(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.python_version", "3.4"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_webSockets(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_webSockets(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.websockets_enabled", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMAppService_scmType(t *testing.T) {
	resourceName := "azurerm_app_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppService_scmType(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_config.0.scm_type", "LocalGit"),
					resource.TestCheckResourceAttr(resourceName, "source_control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "site_credential.#", "1"),
				),
			},
		},
	})
}

func testCheckAzureRMAppServiceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).appServicesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service" {
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

func testCheckAzureRMAppServiceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		appServiceName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Service: %s", appServiceName)
		}

		client := testAccProvider.Meta().(*ArmClient).appServicesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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

func testAccAzureRMAppService_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_freeTier(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Free"
    size = "F1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  site_config {
    use_32_bit_worker_process = true
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_sharedTier(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Free"
    size = "F1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  site_config {
    use_32_bit_worker_process = true
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_alwaysOn(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    always_on = true
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_httpsOnly(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  https_only          = true
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_32Bit(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    use_32_bit_worker_process = true
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_http2Enabled(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    http2_enabled = true
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_appSettings(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  app_settings {
    "foo" = "bar"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_clientAffinityEnabled(rInt int, location string) string {
	return testAccAzureRMAppService_clientAffinity(rInt, location, true)
}

func testAccAzureRMAppService_clientAffinityDisabled(rInt int, location string) string {
	return testAccAzureRMAppService_clientAffinity(rInt, location, false)
}

func testAccAzureRMAppService_clientAffinity(rInt int, location string, clientAffinity bool) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                    = "acctestAS-%d"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  app_service_plan_id     = "${azurerm_app_service_plan.test.id}"
  client_affinity_enabled = %t
}
`, rInt, location, rInt, rInt, clientAffinity)
}

func testAccAzureRMAppService_mangedServiceIdentity(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                    = "acctestAS-%d"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  app_service_plan_id     = "${azurerm_app_service_plan.test.id}"
  identity = {
    type = "SystemAssigned"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_connectionStrings(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  connection_string {
    name  = "Example"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_oneIpRestriction(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10"
    }
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_manyIpRestrictions(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10"
    }

    ip_restriction {
      ip_address  = "20.20.20.0"
      subnet_mask = "255.255.255.0"
    }

    ip_restriction {
      ip_address  = "30.30.0.0"
      subnet_mask = "255.255.0.0"
    }
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_defaultDocuments(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    default_documents = [
      "first.html",
      "second.jsp",
      "third.aspx",
    ]
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_enabled(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  enabled = false
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_localMySql(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    local_mysql_enabled = true
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_managedPipelineMode(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    managed_pipeline_mode = "Classic"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_remoteDebugging(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    remote_debugging_enabled = true
    remote_debugging_version = "VS2015"
  }

  tags {
    "Hello" = "World"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_tags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  tags {
    "Hello" = "World"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_tagsUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  tags {
    "Hello"     = "World"
    "Terraform" = "AcceptanceTests"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_windowsDotNet(rInt int, location, version string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    dotnet_framework_version = "%s"
  }
}
`, rInt, location, rInt, rInt, version)
}

func testAccAzureRMAppService_windowsJava(rInt int, location, javaVersion, container, containerVersion string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    java_version           = "%s"
    java_container         = "%s"
    java_container_version = "%s"
  }
}
`, rInt, location, rInt, rInt, javaVersion, container, containerVersion)
}

func testAccAzureRMAppService_windowsPHP(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    php_version = "7.1"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_windowsPython(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    python_version = "3.4"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_webSockets(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    websockets_enabled = true
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAppService_scmType(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
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

resource "azurerm_app_service" "test" {
  name                = "acctestAS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    scm_type = "LocalGit"
  }
}
`, rInt, location, rInt, rInt)
}
