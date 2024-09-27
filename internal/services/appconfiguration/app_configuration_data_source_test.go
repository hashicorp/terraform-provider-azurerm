// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appconfiguration_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AppConfigurationDataSource struct{}

func TestAccAppConfigurationDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_configuration", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: AppConfigurationResource{}.standard(data),
		},
		{
			Config: AppConfigurationDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(AppConfigurationResource{}),
				check.That(data.ResourceName).Key("endpoint").Exists(),
				check.That(data.ResourceName).Key("encryption.#").HasValue("1"),
				check.That(data.ResourceName).Key("encryption.0.key_vault_key_identifier").Exists(),
				check.That(data.ResourceName).Key("encryption.0.identity_client_id").Exists(),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").Exists(),
				check.That(data.ResourceName).Key("identity.0.identity_ids.0").IsSet(),
				check.That(data.ResourceName).Key("local_auth_enabled").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("primary_read_key.0.connection_string").Exists(),
				check.That(data.ResourceName).Key("primary_read_key.0.id").Exists(),
				check.That(data.ResourceName).Key("primary_read_key.0.secret").Exists(),
				check.That(data.ResourceName).Key("primary_write_key.0.connection_string").Exists(),
				check.That(data.ResourceName).Key("primary_write_key.0.id").Exists(),
				check.That(data.ResourceName).Key("primary_write_key.0.secret").Exists(),
				check.That(data.ResourceName).Key("public_network_access").Exists(),
				check.That(data.ResourceName).Key("purge_protection_enabled").Exists(),
				check.That(data.ResourceName).Key("secondary_read_key.0.connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_read_key.0.id").Exists(),
				check.That(data.ResourceName).Key("secondary_read_key.0.secret").Exists(),
				check.That(data.ResourceName).Key("secondary_write_key.0.connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_write_key.0.id").Exists(),
				check.That(data.ResourceName).Key("secondary_write_key.0.secret").Exists(),
				check.That(data.ResourceName).Key("sku").Exists(),
				check.That(data.ResourceName).Key("soft_delete_retention_days").Exists(),
				check.That(data.ResourceName).Key("replica.#").HasValue("2"),
				check.That(data.ResourceName).Key("replica.0.name").Exists(),
				check.That(data.ResourceName).Key("replica.1.id").Exists(),
				check.That(data.ResourceName).Key("replica.1.endpoint").Exists(),
			),
		},
	})
}

func (AppConfigurationDataSource) basic(data acceptance.TestData) string {
	template := AppConfigurationResource{}.complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_configuration" "test" {
  name                = azurerm_app_configuration.test.name
  resource_group_name = azurerm_app_configuration.test.resource_group_name
}
`, template)
}
