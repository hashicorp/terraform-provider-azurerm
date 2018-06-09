package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMAppService_importBasic(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppService_import32Bit(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppService_importAlwaysOn(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppService_importAppSettings(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppService_importClientAffinityEnabled(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppService_importConnectionStrings(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppService_importIpRestriction(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppService_importDefaultDocuments(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppService_importEnabled(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppService_importLocalMySql(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppService_importManagedPipelineMode(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppService_importRemoteDebugging(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppService_importWindowsDotNet2(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppService_importWindowsDotNet4(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppService_importWindowsJava7Jetty(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppService_importWindowsJava8Jetty(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppService_importWindowsJava7Tomcat(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppService_importWindowsJava8Tomcat(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppService_importWindowsPHP7(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppService_importWindowsPython(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAppService_importWebSockets(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
