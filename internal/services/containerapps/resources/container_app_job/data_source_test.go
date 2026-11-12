// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_app_job_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ContainerAppJobDataSource struct{}

func TestAccContainerAppJobDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_app_job", "test")
	r := ContainerAppJobDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("container_app_environment_id").Exists(),
				check.That(data.ResourceName).Key("replica_timeout_in_seconds").HasValue("10"),
				check.That(data.ResourceName).Key("replica_retry_limit").HasValue("10"),
				check.That(data.ResourceName).Key("manual_trigger_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("manual_trigger_config.0.parallelism").HasValue("4"),
				check.That(data.ResourceName).Key("manual_trigger_config.0.replica_completion_count").HasValue("1"),
				check.That(data.ResourceName).Key("template.#").HasValue("1"),
				check.That(data.ResourceName).Key("template.0.container.#").HasValue("1"),
			),
		},
	})
}

func TestAccContainerAppJobDataSource_readSecretsDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_app_job", "test")
	r := ContainerAppJobDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.readSecretsDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("read_secrets").HasValue("false"),
				check.That(data.ResourceName).Key("secret.#").HasValue("0"),
			),
		},
	})
}

func (d ContainerAppJobDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_app_job" "test" {
  name                = azurerm_container_app_job.test.name
  resource_group_name = azurerm_container_app_job.test.resource_group_name
}
`, ContainerAppJobResource{}.basic(data))
}

func (d ContainerAppJobDataSource) readSecretsDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_app_job" "test" {
  name                = azurerm_container_app_job.test.name
  resource_group_name = azurerm_container_app_job.test.resource_group_name
  read_secrets        = false
}
`, ContainerAppJobResource{}.basic(data))
}
