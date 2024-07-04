// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package elasticsan_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ElasticSANVolumeGroupDataSource struct{}

func TestAccElasticSANVolumeGroupDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_elastic_san_volume_group", "test")
	d := ElasticSANVolumeGroupDataSource{}

	data.DataSourceTestInSequence(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("encryption_type").HasValue("EncryptionAtRestWithCustomerManagedKey"),
				check.That(data.ResourceName).Key("protocol_type").HasValue("Iscsi"),
				check.That(data.ResourceName).Key("encryption.#").HasValue("1"),
				check.That(data.ResourceName).Key("encryption.0.key_vault_key_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("encryption.0.user_assigned_identity_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("encryption.0.current_versioned_key_expiration_timestamp").IsNotEmpty(),
				check.That(data.ResourceName).Key("encryption.0.current_versioned_key_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("encryption.0.last_key_rotation_timestamp").IsNotEmpty(),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("network_rule.#").HasValue("2"),
				check.That(data.ResourceName).Key("network_rule.0.subnet_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("network_rule.0.action").HasValue("Allow"),
			),
		},
	})
}

func (d ElasticSANVolumeGroupDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_elastic_san_volume_group" "test" {
  name           = azurerm_elastic_san_volume_group.test.name
  elastic_san_id = azurerm_elastic_san_volume_group.test.elastic_san_id
}
`, ElasticSANVolumeGroupTestResource{}.complete(data))
}
