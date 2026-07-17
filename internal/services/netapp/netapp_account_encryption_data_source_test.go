// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type NetAppAccountEncryptionDataSource struct{}

func TestAccNetAppAccountEncryptionDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_netapp_account_encryption", "test")
	d := NetAppAccountEncryptionDataSource{}

	tenantID := os.Getenv("ARM_TENANT_ID")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("netapp_account_id").Exists(),
				check.That(data.ResourceName).Key("system_assigned_identity_principal_id").IsSet(),
				check.That(data.ResourceName).Key("encryption_key").IsSet(),
				check.That(data.ResourceName).Key("cross_tenant_key_vault_resource_id").IsEmpty(),
			),
		},
	})
}

func (d NetAppAccountEncryptionDataSource) basic(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
%s

data "azurerm_netapp_account_encryption" "test" {
  netapp_account_id = azurerm_netapp_account_encryption.test.netapp_account_id
}
`, NetAppAccountEncryptionResource{}.cmkSystemAssigned(data, tenantID))
}
