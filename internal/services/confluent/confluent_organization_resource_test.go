// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package confluent_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/organizationresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ConfluentOrganizationResource struct{}

func TestAccConfluentOrganization_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_confluent_organization", "test")
	r := ConfluentOrganizationResource{}

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

func TestAccConfluentOrganization_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_confluent_organization", "test")
	r := ConfluentOrganizationResource{}

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

func TestAccConfluentOrganization_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_confluent_organization", "test")
	r := ConfluentOrganizationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("link_organization.0.token"),
	})
}

func TestAccConfluentOrganization_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_confluent_organization", "test")
	r := ConfluentOrganizationResource{}

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

func (r ConfluentOrganizationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := organizationresources.ParseOrganizationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Confluent.OrganizationResourcesClient.OrganizationGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ConfluentOrganizationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-confluent-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ConfluentOrganizationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

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
    first_name    = "Test"
    last_name     = "User"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ConfluentOrganizationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_confluent_organization" "import" {
  name                = azurerm_confluent_organization.test.name
  resource_group_name = azurerm_confluent_organization.test.resource_group_name
  location            = azurerm_confluent_organization.test.location

  offer_detail {
    id           = "confluent-cloud-azure-prod"
    plan_id      = "confluent-cloud-azure-payg-prod"
    plan_name    = "Confluent Cloud - Pay as you Go"
    publisher_id = "confluentinc"
    term_unit    = "P1M"
  }

  user_detail {
    email_address = "test-%d@example.com"
    first_name    = "Test"
    last_name     = "User"
  }
}
`, r.basic(data), data.RandomInteger)
}

func (r ConfluentOrganizationResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

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
    term_id      = "term-001"
  }

  user_detail {
    email_address = "test-%d@example.com"
    first_name    = "Test"
    last_name     = "User"
  }

  link_organization {
    token = "test-link-token-123"
  }

  tags = {
    Environment = "Test"
    Purpose     = "Acceptance"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ConfluentOrganizationResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

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
    first_name    = "Test"
    last_name     = "User"
  }

  tags = {
    Environment = "Production"
    Updated     = "true"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}
