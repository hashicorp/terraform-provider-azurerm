// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/namespacetopics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type EventgridNamespaceTopicResource struct{}

func TestAccEventgridNamespaceTopic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_namespace_topic", "test")
	r := EventgridNamespaceTopicResource{}

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

func TestAccEventgridNamespaceTopic_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_namespace_topic", "test")
	r := EventgridNamespaceTopicResource{}

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

func TestAccEventgridNamespaceTopic_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_namespace_topic", "test")
	r := EventgridNamespaceTopicResource{}

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

func TestAccEventgridNamespaceTopic_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_namespace_topic", "test")
	r := EventgridNamespaceTopicResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r EventgridNamespaceTopicResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := namespacetopics.ParseNamespaceTopicID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.EventGrid.NamespaceTopicsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r EventgridNamespaceTopicResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_eventgrid_namespace_topic" "test" {
  name                   = "acctest-egnt-%d"
  eventgrid_namespace_id = azurerm_eventgrid_namespace.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r EventgridNamespaceTopicResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_eventgrid_namespace_topic" "import" {
  name                   = azurerm_eventgrid_namespace_topic.test.name
  eventgrid_namespace_id = azurerm_eventgrid_namespace_topic.test.eventgrid_namespace_id
}
`, r.basic(data))
}

func (r EventgridNamespaceTopicResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_eventgrid_namespace_topic" "test" {
  name                    = "acctest-egnt-%d"
  eventgrid_namespace_id  = azurerm_eventgrid_namespace.test.id
  event_retention_in_days = 1
}
`, r.template(data), data.RandomInteger)
}

func (r EventgridNamespaceTopicResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventgrid_namespace" "test" {
  name                = "acctest-egn-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
