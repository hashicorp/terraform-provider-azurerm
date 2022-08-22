package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type TriggerScheduleResource struct{}

func TestAccDataFactoryTriggerSchedule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_schedule", "test")
	r := TriggerScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryTriggerSchedule_pipeline(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_schedule", "test")
	r := TriggerScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.pipeline(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryTriggerSchedule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_schedule", "test")
	r := TriggerScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryTriggerSchedule_scheduleWeekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_schedule", "test")
	r := TriggerScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.scheduleWeekly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryTriggerSchedule_scheduleMonthly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_schedule", "test")
	r := TriggerScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.scheduleMonthly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t TriggerScheduleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.TriggerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataFactory.TriggersClient.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (TriggerScheduleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_pipeline" "test" {
  name            = "acctest%[1]d"
  data_factory_id = azurerm_data_factory.test.id

  parameters = {
    test = "testparameter"
  }
}

resource "azurerm_data_factory_trigger_schedule" "test" {
  name            = "acctestdf%[1]d"
  data_factory_id = azurerm_data_factory.test.id
  pipeline_name   = azurerm_data_factory_pipeline.test.name

  annotations = ["test1", "test2", "test3"]
}
`, data.RandomInteger, data.Locations.Primary)
}

func (TriggerScheduleResource) pipeline(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%[1]d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_pipeline" "test1" {
  name            = "acctest%[1]d"
  data_factory_id = azurerm_data_factory.test.id

  parameters = {
    test = "testparameter1"
  }
}

resource "azurerm_data_factory_pipeline" "test2" {
  name            = "acctests%[1]d"
  data_factory_id = azurerm_data_factory.test.id

  parameters = {
    test = "testparameter2"
  }
}

resource "azurerm_data_factory_trigger_schedule" "test" {
  name            = "acctestdf%[1]d"
  data_factory_id = azurerm_data_factory.test.id

  pipeline {
    name       = azurerm_data_factory_pipeline.test1.name
    parameters = azurerm_data_factory_pipeline.test1.parameters
  }

  pipeline {
    name       = azurerm_data_factory_pipeline.test2.name
    parameters = azurerm_data_factory_pipeline.test2.parameters
  }

  annotations = ["test1", "test2", "test3"]
}
`, data.RandomInteger, data.Locations.Primary)
}

func (TriggerScheduleResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_pipeline" "test" {
  name            = "acctest%[1]d"
  data_factory_id = azurerm_data_factory.test.id

  parameters = {
    test = "testparameter"
  }
}

resource "azurerm_data_factory_trigger_schedule" "test" {
  name                = "acctestdf%[1]d"
  data_factory_id     = azurerm_data_factory.test.id
  pipeline_name       = azurerm_data_factory_pipeline.test.name
  description         = "test"
  pipeline_parameters = azurerm_data_factory_pipeline.test.parameters
  annotations         = ["test5"]
  frequency           = "Day"
  interval            = 5
  activated           = true
  end_time            = "2022-09-22T00:00:00Z"
  start_time          = "2022-09-21T00:00:00Z"
  time_zone           = "GMT Standard Time"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (TriggerScheduleResource) scheduleWeekly(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_pipeline" "test" {
  name            = "acctest%[1]d"
  data_factory_id = azurerm_data_factory.test.id

  parameters = {
    test = "testparameter"
  }
}

resource "azurerm_data_factory_trigger_schedule" "test" {
  name            = "acctestdf%[1]d"
  data_factory_id = azurerm_data_factory.test.id
  pipeline_name   = azurerm_data_factory_pipeline.test.name

  annotations = ["test1", "test2", "test3"]
  activated   = true
  frequency   = "Week"

  schedule {
    minutes      = [0, 30, 59]
    hours        = [0, 12, 23]
    days_of_week = ["Monday", "Tuesday"]
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (TriggerScheduleResource) scheduleMonthly(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_pipeline" "test" {
  name            = "acctest%[1]d"
  data_factory_id = azurerm_data_factory.test.id

  parameters = {
    test = "testparameter"
  }
}

resource "azurerm_data_factory_trigger_schedule" "test" {
  name            = "acctestdf%[1]d"
  data_factory_id = azurerm_data_factory.test.id
  pipeline_name   = azurerm_data_factory_pipeline.test.name

  annotations = ["test1", "test2", "test3"]
  frequency   = "Month"
  interval    = 1
  activated   = true

  schedule {
    hours         = [0, 12, 23]
    minutes       = [0, 30, 59]
    days_of_month = [1, 2, 3]
    monthly {
      weekday = "Monday"
      week    = 1
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
