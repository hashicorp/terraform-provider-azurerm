package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMKeyVaultCertificate_importPFX(t *testing.T) {
	resourceName := "azurerm_key_vault_certificate.test"

	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultCertificate_basicImportPFX(rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate"},
			},
		},
	})
}

func TestAccAzureRMKeyVaultCertificate_importGenerated(t *testing.T) {
	resourceName := "azurerm_key_vault_certificate.test"

	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultCertificate_basicGenerate(rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultCertificateDestroy,
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
