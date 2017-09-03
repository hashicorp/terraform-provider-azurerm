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
			resource.TestStep{
				Config: config,
			},

			resource.TestStep{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_dns_registration", "skip_custom_domain_verification", "skip_dns_registration", "delete_metrics"},
			},
		},
	})
}

func TestAccAzureRMAppService_importComplete(t *testing.T) {
	resourceName := "azurerm_app_service.test"

	ri := acctest.RandInt()
	config := testAccAzureRMAppService_complete(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
			},

			resource.TestStep{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_dns_registration", "skip_custom_domain_verification", "skip_dns_registration", "delete_metrics"},
			},
		},
	})
}
