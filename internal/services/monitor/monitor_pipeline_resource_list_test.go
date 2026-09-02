// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccMonitorPipeline_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_pipeline", "test")
	r := MonitorPipelineResource{}

	listResourceAddress := "azurerm_monitor_pipeline.list"

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 2),
				},
			},
			{
				Query:  true,
				Config: r.basicListQueryByResourceGroupName(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 2),
					querycheck.ExpectIdentity(listResourceAddress, map[string]knownvalue.Check{
						"name":                knownvalue.StringExact(fmt.Sprintf("acctest-pg-0-%d", data.RandomInteger)),
						"resource_group_name": knownvalue.StringExact(fmt.Sprintf("acctestRG-monitor-pipeline-%d", data.RandomInteger)),
						"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
					}),
				},
			},
		},
	})
}

func (r MonitorPipelineResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_monitor_pipeline" "test" {
  count               = 2
  name                = "acctest-pg-${count.index}-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = azurerm_extended_location_custom_location.test.id

  exporter {
    name = "acctest-exporter"

    azure_monitor_workspace_logs {
      api {
        data_collection_endpoint_url      = azurerm_monitor_data_collection_endpoint.test.logs_ingestion_endpoint
        data_collection_rule_immutable_id = azurerm_monitor_data_collection_rule.test.immutable_id
        stream                            = "Custom-${azurerm_log_analytics_workspace_table_custom_log.test.name}"

        schema {
          record_map {
            from = "body"
            to   = "Message"
          }
          record_map {
            from = "time_unix_nano"
            to   = "TimeGenerated"
          }
        }
      }
    }
  }

  receiver {
    name = "acctest-receiver"
    type = "Syslog"

    syslog {
      endpoint = "0.0.0.0:514"
    }
  }

  service {
    pipeline {
      name      = "acctest-pipeline"
      exporters = ["acctest-exporter"]
      receivers = ["acctest-receiver"]
    }
  }

  depends_on = [
    terraform_data.cluster_prereqs,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r MonitorPipelineResource) basicListQuery() string {
	return `
list "azurerm_monitor_pipeline" "list" {
  provider = azurerm
  config {}
}
`
}

func (r MonitorPipelineResource) basicListQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_monitor_pipeline" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctestRG-monitor-pipeline-%[1]d"
  }
}
`, data.RandomInteger)
}
