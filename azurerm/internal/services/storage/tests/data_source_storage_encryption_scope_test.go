package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMStorageEncryptionScope_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_encryption_scope", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageEncryptionScopeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMStorageEncryptionScope_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageEncryptionScopeExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "source"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMStorageEncryptionScope_basic(data acceptance.TestData) string {
	basic := testAccAzureRMStorageEncryptionScope_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_encryption_scope" "test" {
  name               = azurerm_storage_encryption_scope.test.name
  storage_account_id = azurerm_storage_encryption_scope.test.storage_account_id
}
`, basic)
}
