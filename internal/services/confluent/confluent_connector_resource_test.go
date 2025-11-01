// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package confluent_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/connectorresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ConfluentConnectorResource struct{}

func TestAccConfluentConnector_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_confluent_connector", "test")
	r := ConfluentConnectorResource{}

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

func TestAccConfluentConnector_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_confluent_connector", "test")
	r := ConfluentConnectorResource{}

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

func TestAccConfluentConnector_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_confluent_connector", "test")
	r := ConfluentConnectorResource{}

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

func (r ConfluentConnectorResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := connectorresources.ParseConnectorID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Confluent.ConnectorClient.ConnectorGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (r ConfluentConnectorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-confluent-%d"
  location = "%s"
}

resource "azurerm_confluent_organization" "test" {
  name                = "acctest-co-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  offer_detail {
    id           = "confluent-cloud-azure-prod"
    plan_id      = "confluent-cloud-azure-payg-prod"
    plan_name    = "Confluent Cloud - Pay as you Go"
    publisher_id = "confluentinc"
    term_unit    = "P1M"
  }

  user_detail {
    email_address = "test-%d@example.com"
  }
}

resource "azurerm_confluent_environment" "test" {
  environment_id      = "env-%d"
  organization_id     = azurerm_confluent_organization.test.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_confluent_cluster" "test" {
  cluster_id          = "lkc-%d"
  environment_id      = azurerm_confluent_environment.test.environment_id
  organization_id     = azurerm_confluent_organization.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ConfluentConnectorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_confluent_connector" "test" {
  connector_name      = "connector-%d"
  cluster_id          = azurerm_confluent_cluster.test.cluster_id
  environment_id      = azurerm_confluent_environment.test.environment_id
  organization_id     = azurerm_confluent_organization.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, r.template(data), data.RandomInteger)
}

func (r ConfluentConnectorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_confluent_connector" "import" {
  connector_name      = azurerm_confluent_connector.test.connector_name
  cluster_id          = azurerm_confluent_connector.test.cluster_id
  environment_id      = azurerm_confluent_connector.test.environment_id
  organization_id     = azurerm_confluent_connector.test.organization_id
  resource_group_name = azurerm_confluent_connector.test.resource_group_name
}
`, r.basic(data))
}

func (r ConfluentConnectorResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_confluent_connector" "test" {
  connector_name      = "connector-%d"
  cluster_id          = azurerm_confluent_cluster.test.cluster_id
  environment_id      = azurerm_confluent_environment.test.environment_id
  organization_id     = azurerm_confluent_organization.test.name
  resource_group_name = azurerm_resource_group.test.name
  connector_type      = "SINK"
  connector_class     = "AZUREBLOBSINK"
}
`, r.template(data), data.RandomInteger)
}
