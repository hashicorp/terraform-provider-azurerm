// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package communication_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/emailservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type EmailServiceTestResource struct{}

func TestAccEmailService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_email_communication_service", "test")
	r := EmailServiceTestResource{}

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

func TestAccEmailService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_email_communication_service", "test")
	r := EmailServiceTestResource{}

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

func TestAccEmailService_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_email_communication_service", "test")
	r := EmailServiceTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
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

func (r EmailServiceTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	clusterClient := client.Communication.EmailServicesClient
	id, err := emailservices.ParseEmailServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clusterClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving Email Communication Service %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r EmailServiceTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_email_communication_service" "test" {
  name                = "acctest-CommunicationService-%d"
  resource_group_name = azurerm_resource_group.test.name
  data_location       = "United States"
}
`, r.template(data), data.RandomInteger)
}

func (r EmailServiceTestResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_email_communication_service" "import" {
  name                = azurerm_email_communication_service.test.name
  resource_group_name = azurerm_email_communication_service.test.resource_group_name
  data_location       = azurerm_email_communication_service.test.data_location
}
`, config)
}

func (r EmailServiceTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_email_communication_service" "test" {
  name                = "acctest-CommunicationService-%d"
  resource_group_name = azurerm_resource_group.test.name
  data_location       = "United States"

  tags = {
    env = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r EmailServiceTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_email_communication_service" "test" {
  name                = "acctest-CommunicationService-%d"
  resource_group_name = azurerm_resource_group.test.name
  data_location       = "United States"

  tags = {
    env = "Test2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r EmailServiceTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-communicationservice-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
