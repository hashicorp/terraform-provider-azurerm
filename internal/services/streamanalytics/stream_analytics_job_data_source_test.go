// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type StreamAnalyticsJobDataSource struct{}

func TestAccDataSourceStreamAnalyticsJob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_stream_analytics_job", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StreamAnalyticsJobDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("job_id").Exists(),
				check.That(data.ResourceName).Key("streaming_units").Exists(),
				check.That(data.ResourceName).Key("transformation_query").Exists(),
			),
		},
	})
}

func TestAccDataSourceStreamAnalyticsJob_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_stream_analytics_job", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StreamAnalyticsJobDataSource{}.identity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("job_id").Exists(),
				check.That(data.ResourceName).Key("streaming_units").Exists(),
				check.That(data.ResourceName).Key("transformation_query").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
			),
		},
	})
}

func TestAccDataSourceStreamAnalyticsJob_jobSchedule(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_stream_analytics_job", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StreamAnalyticsJobDataSource{}.jobSchedule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("start_mode").Exists(),
				check.That(data.ResourceName).Key("start_time").Exists(),
			),
		},
	})
}

func (d StreamAnalyticsJobDataSource) basic(data acceptance.TestData) string {
	config := StreamAnalyticsJobResource{}.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_stream_analytics_job" "test" {
  name                = azurerm_stream_analytics_job.test.name
  resource_group_name = azurerm_stream_analytics_job.test.resource_group_name
}
`, config)
}

func (d StreamAnalyticsJobDataSource) identity(data acceptance.TestData) string {
	config := StreamAnalyticsJobResource{}.identity(data)
	return fmt.Sprintf(`
%s

data "azurerm_stream_analytics_job" "test" {
  name                = azurerm_stream_analytics_job.test.name
  resource_group_name = azurerm_stream_analytics_job.test.resource_group_name
}
`, config)
}

func (d StreamAnalyticsJobDataSource) jobSchedule(data acceptance.TestData) string {
	config := StreamAnalyticsJobScheduleResource{}.customTime(data)
	return fmt.Sprintf(`
%s

data "azurerm_stream_analytics_job" "test" {
  name                = azurerm_stream_analytics_job.test.name
  resource_group_name = azurerm_stream_analytics_job.test.resource_group_name

  depends_on = [azurerm_stream_analytics_job_schedule.test]
}
`, config)
}
