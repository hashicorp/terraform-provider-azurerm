// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package labservice_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/schedule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LabServiceScheduleResource struct{}

func TestAccLabServiceSchedule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lab_service_schedule", "test")
	r := LabServiceScheduleResource{}
	stopTime := time.Now().Add(time.Hour * 1).Format(time.RFC3339)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, stopTime),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLabServiceSchedule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lab_service_schedule", "test")
	r := LabServiceScheduleResource{}
	stopTime := time.Now().Add(time.Hour * 1).Format(time.RFC3339)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, stopTime),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, stopTime),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
	})
}

func TestAccLabServiceSchedule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lab_service_schedule", "test")
	r := LabServiceScheduleResource{}

	startTime := time.Now().Format(time.RFC3339)
	stopTime := time.Now().Add(time.Hour * 1).Format(time.RFC3339)
	expirationDate := time.Now().Add(time.Hour * 1).Format(time.RFC3339)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, startTime, stopTime, expirationDate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLabServiceSchedule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lab_service_schedule", "test")
	r := LabServiceScheduleResource{}

	startTime := time.Now().Format(time.RFC3339)
	stopTime := time.Now().Add(time.Hour * 1).Format(time.RFC3339)
	expirationDate := time.Now().Add(time.Hour * 1).Format(time.RFC3339)

	updatedStartTime := time.Now().Add(time.Hour * 1).Format(time.RFC3339)
	updatedStopTime := time.Now().Add(time.Hour * 2).Format(time.RFC3339)
	updatedExpirationDate := time.Now().Add(time.Hour * 2).Format(time.RFC3339)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, startTime, stopTime, expirationDate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, updatedStartTime, updatedStopTime, updatedExpirationDate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LabServiceScheduleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := schedule.ParseScheduleID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.LabService.ScheduleClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r LabServiceScheduleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-labschedule-%d"
  location = "%s"
}

resource "azurerm_lab_service_lab" "test" {
  name                = "acctest-lab-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  title               = "Test Title"

  security {
    open_access_enabled = false
  }

  virtual_machine {
    admin_user {
      username = "testadmin"
      password = "Password1234!"
    }

    image_reference {
      offer     = "0001-com-ubuntu-server-focal"
      publisher = "canonical"
      sku       = "20_04-lts"
      version   = "latest"
    }

    sku {
      name     = "Classic_Fsv2_2_4GB_128_S_SSD"
      capacity = 1
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(17))
}

func (r LabServiceScheduleResource) basic(data acceptance.TestData, stopTime string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lab_service_schedule" "test" {
  name      = "acctest-labschedule-%d"
  lab_id    = azurerm_lab_service_lab.test.id
  stop_time = "%s"
  time_zone = "America/Los_Angeles"
}
`, r.template(data), data.RandomInteger, stopTime)
}

func (r LabServiceScheduleResource) requiresImport(data acceptance.TestData, stopTime string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lab_service_schedule" "import" {
  name      = azurerm_lab_service_schedule.test.name
  lab_id    = azurerm_lab_service_schedule.test.lab_id
  stop_time = azurerm_lab_service_schedule.test.stop_time
  time_zone = azurerm_lab_service_schedule.test.time_zone
}
`, r.basic(data, stopTime))
}

func (r LabServiceScheduleResource) complete(data acceptance.TestData, startTime, stopTime, expirationDate string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lab_service_schedule" "test" {
  name       = "acctest-labschedule-%d"
  lab_id     = azurerm_lab_service_lab.test.id
  notes      = "Testing"
  start_time = "%s"
  stop_time  = "%s"
  time_zone  = "America/Los_Angeles"

  recurrence {
    expiration_date = "%s"
    frequency       = "Weekly"
    interval        = 1
    week_days       = ["Friday", "Thursday"]
  }
}
`, r.template(data), data.RandomInteger, startTime, stopTime, expirationDate)
}

func (r LabServiceScheduleResource) update(data acceptance.TestData, startTime, stopTime, expirationDate string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lab_service_schedule" "test" {
  name       = "acctest-labschedule-%d"
  lab_id     = azurerm_lab_service_lab.test.id
  notes      = "Testing2"
  start_time = "%s"
  stop_time  = "%s"
  time_zone  = "America/Grenada"

  recurrence {
    expiration_date = "%s"
    frequency       = "Daily"
    interval        = 2
  }
}
`, r.template(data), data.RandomInteger, startTime, stopTime, expirationDate)
}
