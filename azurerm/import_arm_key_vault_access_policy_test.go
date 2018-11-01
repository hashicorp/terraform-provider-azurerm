package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMKeyVaultAccessPolicy_importBasic(t *testing.T) {
	resourceName := "azurerm_key_vault_access_policy.test"

	rs := acctest.RandString(5)
	config := testAccAzureRMKeyVaultAccessPolicy_basic(rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
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

func TestAccAzureRMKeyVaultAccessPolicy_importMultiple(t *testing.T) {
	rs := acctest.RandString(5)
	config := testAccAzureRMKeyVaultAccessPolicy_multiple(rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_key_vault_access_policy.test_with_application_id",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "azurerm_key_vault_access_policy.test_no_application_id",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
