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

func TestAccContainerAppJobDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_app_job", "test")
	r := ContainerAppJobDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("secret.#").HasValue("1"),
				check.That(data.ResourceName).Key("registry.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("template.0.init_container.#").HasValue("1"),
				check.That(data.ResourceName).Key("template.0.volume.#").HasValue("1"),
				check.That(data.ResourceName).Key("tags.ENV").HasValue("test"),
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

func TestAccContainerAppJobDataSource_eventTrigger(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_app_job", "test")
	r := ContainerAppJobDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.eventTrigger(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("event_trigger_config.0.parallelism").HasValue("4"),
				check.That(data.ResourceName).Key("event_trigger_config.0.scale.0.max_executions").HasValue("10"),
				check.That(data.ResourceName).Key("event_trigger_config.0.scale.0.rules.0.name").HasValue("servicebuscalingrule"),
				check.That(data.ResourceName).Key("event_trigger_config.0.scale.0.rules.0.custom_rule_type").HasValue("azure-servicebus"),
				check.That(data.ResourceName).Key("event_trigger_config.0.scale.0.rules.0.metadata.topic_name").HasValue("my-topic"),
				check.That(data.ResourceName).Key("event_trigger_config.0.scale.0.rules.0.identity_id").Exists(),
			),
		},
	})
}

func TestAccContainerAppJobDataSource_scheduleTrigger(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_app_job", "test")
	r := ContainerAppJobDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.scheduleTrigger(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("schedule_trigger_config.0.cron_expression").HasValue("*/1 * * * *"),
				check.That(data.ResourceName).Key("schedule_trigger_config.0.parallelism").HasValue("1"),
				check.That(data.ResourceName).Key("schedule_trigger_config.0.replica_completion_count").HasValue("1"),
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

func (d ContainerAppJobDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_app_job" "test" {
  name                = azurerm_container_app_job.test.name
  resource_group_name = azurerm_container_app_job.test.resource_group_name
}
`, ContainerAppJobResource{}.complete(data))
}

func (d ContainerAppJobDataSource) readSecretsDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_app_job" "test" {
  name                = azurerm_container_app_job.test.name
  resource_group_name = azurerm_container_app_job.test.resource_group_name
  read_secrets        = false
}
`, ContainerAppJobResource{}.complete(data))
}

func (d ContainerAppJobDataSource) eventTrigger(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_app_job" "test" {
  name                = azurerm_container_app_job.test.name
  resource_group_name = azurerm_container_app_job.test.resource_group_name
}
`, ContainerAppJobResource{}.eventTriggerScaleRuleSystemAndUserIdentity(data))
}

func (d ContainerAppJobDataSource) scheduleTrigger(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_app_job" "test" {
  name                = azurerm_container_app_job.test.name
  resource_group_name = azurerm_container_app_job.test.resource_group_name
}
`, ContainerAppJobResource{}.scheduleTrigger(data))
}
