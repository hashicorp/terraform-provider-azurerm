package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMDiskEncryptionSet_basic(t *testing.T) {
	dataSourceName := "data.azurerm_disk_encryption_set.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(6)
	resourceGroup := fmt.Sprintf("acctestRG-%d", ri)
	vaultName := fmt.Sprintf("vault%d", ri)
	keyName := fmt.Sprintf("key-%s", rs)
	desName := fmt.Sprintf("acctestdes-%d", ri)
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDiskEncryptionSetDestroy,
		Steps: []resource.TestStep{
			// These two steps are used for setting up a valid keyVault, which enables soft-delete and purge protection.
			{
				Config:  testAccPrepareKeyvaultAndKey(resourceGroup, location, vaultName, keyName),
				Destroy: false,
				Check:   resource.ComposeTestCheckFunc(),
			},
			// This step is not negligible, without this step, the final step will fail on refresh complaining `Disk Encryption Set does not exist`
			{
				PreConfig: func() { enableSoftDeleteAndPurgeProtectionForKeyvault(resourceGroup, vaultName) },
				Config:    testAccAzureRMDiskEncryptionSet_basic(resourceGroup, location, vaultName, keyName, desName),
			},
			{
				Config: testAccDataSourceDiskEncryptionSet_basic(resourceGroup, location, vaultName, keyName, desName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "identity.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "identity.0.type", "SystemAssigned"),
					resource.TestCheckResourceAttrSet(dataSourceName, "identity.0.principal_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "identity.0.tenant_id"),
					resource.TestCheckResourceAttr(dataSourceName, "active_key.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "active_key.0.source_vault_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "active_key.0.key_url"),
				),
			},
		},
	})
}

func testAccDataSourceDiskEncryptionSet_basic(resourceGroup, location, vaultName, keyName, desName string) string {
	config := testAccAzureRMDiskEncryptionSet_basic(resourceGroup, location, vaultName, keyName, desName)
	return fmt.Sprintf(`
%s

data "azurerm_disk_encryption_set" "test" {
  resource_group_name = "${azurerm_disk_encryption_set.test.resource_group_name}"
  name                = "${azurerm_disk_encryption_set.test.name}"
}
`, config)
}
