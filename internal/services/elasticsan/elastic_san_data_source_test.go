// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package elasticsan_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ElasticSANDataSource struct{}

func TestAccElasticSANDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_elastic_san", "test")
	d := ElasticSANDataSource{}

	data.DataSourceTestInSequence(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku.#").HasValue("1"),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Premium_LRS"),
				check.That(data.ResourceName).Key("sku.0.tier").HasValue("Premium"),
				check.That(data.ResourceName).Key("location").IsNotEmpty(),
				check.That(data.ResourceName).Key("base_size_in_tib").HasValue("2"),
				check.That(data.ResourceName).Key("extended_size_in_tib").HasValue("4"),
				check.That(data.ResourceName).Key("zones.#").HasValue("2"),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("total_iops").Exists(),
				check.That(data.ResourceName).Key("total_mbps").Exists(),
				check.That(data.ResourceName).Key("total_size_in_tib").Exists(),
				check.That(data.ResourceName).Key("total_volume_size_in_gib").Exists(),
				check.That(data.ResourceName).Key("volume_group_count").Exists(),
			),
		},
	})
}

func (d ElasticSANDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_elastic_san" "test" {
  name                = azurerm_elastic_san.test.name
  resource_group_name = azurerm_elastic_san.test.resource_group_name
}
`, ElasticSANTestResource{}.complete(data))
}
