// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package elastic_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ElasticsearchDataSource struct{}

func TestAccElasticsearchDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_cloud_elasticsearch", "test")
	r := ElasticsearchDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("elastic_cloud_email_address").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("sku_name").Exists(),
				check.That(data.ResourceName).Key("monitoring_enabled").Exists(),
				check.That(data.ResourceName).Key("elastic_cloud_deployment_id").Exists(),
				check.That(data.ResourceName).Key("elastic_cloud_sso_default_url").Exists(),
				check.That(data.ResourceName).Key("elastic_cloud_user_id").Exists(),
				check.That(data.ResourceName).Key("elasticsearch_service_url").Exists(),
				check.That(data.ResourceName).Key("kibana_service_url").Exists(),
				check.That(data.ResourceName).Key("kibana_sso_uri").Exists(),
			),
		},
	})
}

func (ElasticsearchDataSource) basic(data acceptance.TestData) string {
	template := ElasticsearchResource{}.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_elastic_cloud_elasticsearch" "test" {
  name                = azurerm_elastic_cloud_elasticsearch.test.name
  resource_group_name = azurerm_elastic_cloud_elasticsearch.test.resource_group_name
}
`, template)
}
