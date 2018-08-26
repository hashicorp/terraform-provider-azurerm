package azurerm

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMAppServiceCustomHostnameBinding(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests due to
	// the app service name being shared (so the tests don't conflict with each other)
	testCases := map[string]map[string]func(t *testing.T){
		"basic": {
			"basic":  testAccAzureRMAppServiceCustomHostnameBinding_basic,
			"import": testAccAzureRMAppServiceCustomHostnameBinding_import,
		},
	}

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func testAccAzureRMAppServiceCustomHostnameBinding_import(t *testing.T) {
	appServiceEnvVariable := "ARM_TEST_APP_SERVICE"
	appServiceEnv := os.Getenv(appServiceEnvVariable)
	if appServiceEnv == "" {
		t.Skipf("Skipping as %q is not specified", appServiceEnvVariable)
	}

	domainEnvVariable := "ARM_TEST_DOMAIN"
	domainEnv := os.Getenv(domainEnvVariable)
	if domainEnv == "" {
		t.Skipf("Skipping as %q is not specified", domainEnvVariable)
	}

	resourceName := "azurerm_app_service_custom_hostname_binding.test"

	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMAppServiceCustomHostnameBinding_basicConfig(ri, location, appServiceEnv, domainEnv)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceCustomHostnameBindingDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
