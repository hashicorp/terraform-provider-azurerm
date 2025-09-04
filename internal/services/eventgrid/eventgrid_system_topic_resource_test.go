// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/systemtopics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type EventGridSystemTopicResource struct{}

func TestAccEventGridSystemTopic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic", "test")
	r := EventGridSystemTopicResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_arm_resource_id").Exists(),
				check.That(data.ResourceName).Key("topic_type").Exists(),
				check.That(data.ResourceName).Key("metric_arm_resource_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridSystemTopic_policyStates(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic", "test")
	r := EventGridSystemTopicResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.policyStates(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_arm_resource_id").Exists(),
				check.That(data.ResourceName).Key("topic_type").Exists(),
				check.That(data.ResourceName).Key("metric_arm_resource_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridSystemTopic_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic", "test")
	r := EventGridSystemTopicResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_eventgrid_system_topic"),
		},
	})
}

func TestAccEventGridSystemTopic_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic", "test")
	r := EventGridSystemTopicResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Foo").HasValue("Bar"),
				check.That(data.ResourceName).Key("source_arm_resource_id").Exists(),
				check.That(data.ResourceName).Key("topic_type").Exists(),
				check.That(data.ResourceName).Key("metric_arm_resource_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridSystemTopic_basicWithSystemManagedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic", "test")
	r := EventGridSystemTopicResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithSystemManagedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("0"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridSystemTopic_basicWithUserAssignedManagedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_system_topic", "test")
	r := EventGridSystemTopicResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithUserAssignedManagedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsEmpty(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func (EventGridSystemTopicResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := systemtopics.ParseSystemTopicID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.EventGrid.SystemTopics.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (EventGridSystemTopicResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestegst%[3]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctestEGST%[1]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_storage_account.test.id
  topic_type             = "Microsoft.Storage.StorageAccounts"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(12))
}

func (r EventGridSystemTopicResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_eventgrid_system_topic" "import" {
  name                   = azurerm_eventgrid_system_topic.test.name
  location               = azurerm_eventgrid_system_topic.test.location
  resource_group_name    = azurerm_eventgrid_system_topic.test.resource_group_name
  source_arm_resource_id = azurerm_eventgrid_system_topic.test.source_arm_resource_id
  topic_type             = azurerm_eventgrid_system_topic.test.topic_type
}
`, r.basic(data))
}

func (EventGridSystemTopicResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestegst%[3]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctestEGST%[1]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_storage_account.test.id
  topic_type             = "Microsoft.Storage.StorageAccounts"

  tags = {
    "Foo" = "Bar"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(12))
}

func (EventGridSystemTopicResource) policyStates(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctestEGST%[1]d"
  location               = "Global"
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = format("/subscriptions/%%s", data.azurerm_subscription.current.subscription_id)
  topic_type             = "Microsoft.PolicyInsights.PolicyStates"

  tags = {
    "Foo" = "Bar"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(10))
}

func (EventGridSystemTopicResource) basicWithSystemManagedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestegst%[3]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctesteg-%[1]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_storage_account.test.id
  topic_type             = "Microsoft.Storage.StorageAccounts"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(12))
}

func (EventGridSystemTopicResource) basicWithUserAssignedManagedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestegst%[3]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctesteg-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_eventgrid_system_topic" "test" {
  name                   = "acctesteg-%[1]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  source_arm_resource_id = azurerm_storage_account.test.id
  topic_type             = "Microsoft.Storage.StorageAccounts"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(12))
}
