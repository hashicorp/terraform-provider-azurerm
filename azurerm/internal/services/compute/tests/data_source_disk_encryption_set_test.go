package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMDiskEncryptionSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_disk_encryption_set", "test")
	resourceGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	vaultName := fmt.Sprintf("vault%d", data.RandomInteger)
	keyName := fmt.Sprintf("key-%s", data.RandomString)
	desName := fmt.Sprintf("acctestdes-%d", data.RandomInteger)
	location := data.Locations.Primary

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDiskEncryptionSetDestroy,
		Steps: []resource.TestStep{
			// These two steps are used for setting up a valid keyVault, which enables soft-delete and purge protection.
			// TODO: After applying soft-delete and purge-protection in keyVault, this extra step can be removed.
			{
				Config: testAccPrepareKeyvaultAndKey(resourceGroup, location, vaultName, keyName),
				Check: resource.ComposeTestCheckFunc(
					enableSoftDeleteAndPurgeProtectionForKeyvault(resourceGroup, vaultName),
				),
			},
			// This step is not negligible, without this step, the final step will fail on refresh complaining `Disk Encryption Set does not exist`
			{
				Config: testAccAzureRMDiskEncryptionSet_basic(resourceGroup, location, vaultName, keyName, desName),
			},
			{
				Config: testAccDataSourceDiskEncryptionSet_basic(resourceGroup, location, vaultName, keyName, desName),
				Check:  resource.ComposeTestCheckFunc(),
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
