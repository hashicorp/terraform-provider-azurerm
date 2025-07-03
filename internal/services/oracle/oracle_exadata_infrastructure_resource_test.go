// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudexadatainfrastructures"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ExadataInfraResource struct{}

func (a ExadataInfraResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := cloudexadatainfrastructures.ParseCloudExadataInfrastructureID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Oracle.OracleClient.CloudExadataInfrastructures.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestExaInfra_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ExadataInfraResource{}.ResourceType(), "test")
	r := ExadataInfraResource{}
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

func TestExaInfra_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ExadataInfraResource{}.ResourceType(), "test")
	r := ExadataInfraResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestExaInfra_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ExadataInfraResource{}.ResourceType(), "test")
	r := ExadataInfraResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestExaInfra_update(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ExadataInfraResource{}.ResourceType(), "test")
	r := ExadataInfraResource{}
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

func (a ExadataInfraResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_oracle_exadata_infrastructure" "test" {
  name                = "OFakeacctest%[2]d"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  compute_count       = "2"
  display_name        = "OFakeacctest%[2]d"
  shape               = "Exadata.X9M"
  storage_count       = "3"
  zones               = ["3"]
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a ExadataInfraResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_oracle_exadata_infrastructure" "test" {
  name                = "OFakeacctest%[2]d"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  compute_count       = "2"
  display_name        = "OFakeacctest%[2]d"
  shape               = "Exadata.X9M"
  storage_count       = "3"
  zones               = ["3"]
  customer_contacts   = ["test@test.com"]

  maintenance_window {
    days_of_week       = ["Monday"]
    hours_of_day       = [4]
    months             = ["January"]
    weeks_of_month     = [2]
    lead_time_in_weeks = 1
    patching_mode      = "Rolling"
    preference         = "NoPreference"
  }

  tags = {
    test = "testTag1"
  }
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a ExadataInfraResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_oracle_exadata_infrastructure" "test" {
  name                = "OFakeacctest%[2]d"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  compute_count       = "2"
  display_name        = "OFakeacctest%[2]d"
  shape               = "Exadata.X9M"
  storage_count       = "3"
  zones               = ["3"]
  tags = {
    test = "testTag1"
  }
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a ExadataInfraResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_exadata_infrastructure" "import" {
  name                = azurerm_oracle_exadata_infrastructure.test.name
  location            = azurerm_oracle_exadata_infrastructure.test.location
  resource_group_name = azurerm_oracle_exadata_infrastructure.test.resource_group_name
  compute_count       = azurerm_oracle_exadata_infrastructure.test.compute_count
  display_name        = azurerm_oracle_exadata_infrastructure.test.display_name
  shape               = azurerm_oracle_exadata_infrastructure.test.shape
  storage_count       = azurerm_oracle_exadata_infrastructure.test.storage_count
  zones               = azurerm_oracle_exadata_infrastructure.test.zones
}
`, a.basic(data))
}

func (a ExadataInfraResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
