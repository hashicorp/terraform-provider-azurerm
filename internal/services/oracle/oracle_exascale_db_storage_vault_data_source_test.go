// Copyright © 2025, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
)

type ExascaleDbStorageVaultDataSource struct{}

func TestDbStorageVaultDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ExascaleDbStorageVaultDataSource{}.ResourceType(), "test")
	r := ExascaleDbStorageVaultDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("additional_flash_cache_in_percent").Exists(),
				check.That(data.ResourceName).Key("description").Exists(),
			),
		},
	})
}

func (d ExascaleDbStorageVaultDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_exascale_db_storage_vault" "test" {
  name                = azurerm_oracle_exascale_db_storage_vault.test.name
  resource_group_name = azurerm_oracle_exascale_db_storage_vault.test.resource_group_name
}
`, ExascaleDbStorageVaultResource{}.basic(data))
}
