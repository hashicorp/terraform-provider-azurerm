// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DataProtectionBackupPolicyBlobStorageDataSource struct{}

func TestAccDataProtectionBackupPolicyBlobStorageDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_protection_backup_policy_blob_storage", "test")
	r := DataProtectionBackupPolicyBlobStorageDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(DataProtectionBackupPolicyBlobStorageResource{}),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("vault_id").Exists(),
				check.That(data.ResourceName).Key("operational_default_retention_duration").HasValue("P30D"),
			),
		},
	})
}

func (r DataProtectionBackupPolicyBlobStorageDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_data_protection_backup_policy_blob_storage" "test" {
  name     = azurerm_data_protection_backup_policy_blob_storage.test.name
  vault_id = azurerm_data_protection_backup_vault.test.id
}
`, DataProtectionBackupPolicyBlobStorageResource{}.basic(data))
}
