// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package communication_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/domains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type EmailCommunicationServiceDomainResource struct{}

func TestAccEmailServiceDomain_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_email_communication_service_domain", "test")
	r := EmailCommunicationServiceDomainResource{}

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

func TestAccEmailServiceDomain_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_email_communication_service_domain", "test")
	r := EmailCommunicationServiceDomainResource{}

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

func TestAccEmailServiceDomain_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_email_communication_service_domain", "test")
	r := EmailCommunicationServiceDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "false"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, "true"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r EmailCommunicationServiceDomainResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	domainClient := client.Communication.DomainClient
	id, err := domains.ParseDomainID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := domainClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}

		return nil, fmt.Errorf("retrieving Email Domain Communication Service %q: %+v", state.ID, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r EmailCommunicationServiceDomainResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_email_communication_service_domain" "test" {
  name             = "AzureManagedDomain"
  email_service_id = azurerm_email_communication_service.test.id

  domain_management = "AzureManaged"
}
`, r.template(data))
}

func (r EmailCommunicationServiceDomainResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_email_communication_service_domain" "import" {
  name             = azurerm_email_communication_service_domain.test.name
  email_service_id = azurerm_email_communication_service_domain.test.email_service_id

  domain_management = azurerm_email_communication_service_domain.test.domain_management
}
`, config)
}

func (r EmailCommunicationServiceDomainResource) complete(data acceptance.TestData, userTrackingEnabled string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_email_communication_service_domain" "test" {
  name             = "example.com"
  email_service_id = azurerm_email_communication_service.test.id

  domain_management                = "CustomerManaged"
  user_engagement_tracking_enabled = %s

  tags = {
    env = "Test"
  }
}
`, r.template(data), userTrackingEnabled)
}

func (r EmailCommunicationServiceDomainResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-communicationservice-%[1]d"
  location = "%[2]s"
}

resource "azurerm_email_communication_service" "test" {
  name                = "acctest-CommunicationService-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  data_location       = "United States"
}
`, data.RandomInteger, data.Locations.Primary)
}
