// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package confluent_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/topicrecords"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

type ConfluentTopicResource struct{}

func TestAccConfluentTopic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_confluent_topic", "test")
	r := ConfluentTopicResource{}

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

func TestAccConfluentTopic_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_confluent_topic", "test")
	r := ConfluentTopicResource{}

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

func TestAccConfluentTopic_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_confluent_topic", "test")
	r := ConfluentTopicResource{}

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

func (r ConfluentTopicResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := topicrecords.ParseTopicID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Confluent.TopicClient.TopicsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(true), nil
}

func (r ConfluentTopicResource) template(data acceptance.TestData) string {
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

func (r ConfluentTopicResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_confluent_topic" "test" {
  topic_name          = "topic-%d"
  cluster_id          = azurerm_confluent_cluster.test.cluster_id
  environment_id      = azurerm_confluent_environment.test.environment_id
  organization_id     = azurerm_confluent_organization.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, r.template(data), data.RandomInteger)
}

func (r ConfluentTopicResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_confluent_topic" "import" {
  topic_name          = azurerm_confluent_topic.test.topic_name
  cluster_id          = azurerm_confluent_topic.test.cluster_id
  environment_id      = azurerm_confluent_topic.test.environment_id
  organization_id     = azurerm_confluent_topic.test.organization_id
  resource_group_name = azurerm_confluent_topic.test.resource_group_name
}
`, r.basic(data))
}

func (r ConfluentTopicResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_confluent_topic" "test" {
  topic_name          = "topic-%d"
  cluster_id          = azurerm_confluent_cluster.test.cluster_id
  environment_id      = azurerm_confluent_environment.test.environment_id
  organization_id     = azurerm_confluent_organization.test.name
  resource_group_name = azurerm_resource_group.test.name
  partitions_count    = "6"
  replication_factor  = "3"

  configs {
    name  = "cleanup.policy"
    value = "compact"
  }

  configs {
    name  = "retention.ms"
    value = "604800000"
  }
}
`, r.template(data), data.RandomInteger)
}
