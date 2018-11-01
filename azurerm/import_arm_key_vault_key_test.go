package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMKeyVaultKey_importBasicEC(t *testing.T) {
	resourceName := "azurerm_key_vault_key.test"

	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultKey_basicEC(rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key_size"},
			},
		},
	})
}

func TestAccAzureRMKeyVaultKey_importBasicRSA(t *testing.T) {
	resourceName := "azurerm_key_vault_key.test"

	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultKey_basicRSA(rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key_size"},
			},
		},
	})
}

func TestAccAzureRMKeyVaultKey_importBasicRSAHSM(t *testing.T) {
	resourceName := "azurerm_key_vault_key.test"

	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultKey_basicRSAHSM(rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key_size"},
			},
		},
	})
}

func TestAccAzureRMKeyVaultKey_importComplete(t *testing.T) {
	resourceName := "azurerm_key_vault_key.test"

	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultKey_complete(rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key_size"},
			},
		},
	})
}
