// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package communication_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/senderusernames"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type EmailCommunicationServiceDomainSenderUsernameResource struct{}

func TestAccEmailServiceDomainSenderUsername_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_email_communication_service_domain_sender_username", "test")
	r := EmailCommunicationServiceDomainSenderUsernameResource{}

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

func TestAccEmailServiceDomainSenderUsername_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_email_communication_service_domain_sender_username", "test")
	r := EmailCommunicationServiceDomainSenderUsernameResource{}

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

func TestAccEmailServiceDomainSenderUsername_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_email_communication_service_domain_sender_username", "test")
	r := EmailCommunicationServiceDomainSenderUsernameResource{}

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

func TestAccEmailServiceDomainSenderUsername_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_email_communication_service_domain_sender_username", "test")
	r := EmailCommunicationServiceDomainSenderUsernameResource{}

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
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r EmailCommunicationServiceDomainSenderUsernameResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := senderusernames.ParseSenderUsernameID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Communication.SenderUsernamesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r EmailCommunicationServiceDomainSenderUsernameResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_email_communication_service_domain_sender_username" "test" {
  name                    = "acctest-su-%d"
  email_service_domain_id = azurerm_email_communication_service_domain.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r EmailCommunicationServiceDomainSenderUsernameResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_email_communication_service_domain_sender_username" "import" {
  name                    = azurerm_email_communication_service_domain_sender_username.test.name
  email_service_domain_id = azurerm_email_communication_service_domain_sender_username.test.email_service_domain_id
}
`, r.basic(data))
}

func (r EmailCommunicationServiceDomainSenderUsernameResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_email_communication_service_domain_sender_username" "test" {
  name                    = "acctest-su-%d"
  email_service_domain_id = azurerm_email_communication_service_domain.test.id
  display_name            = "TFTester"
}
`, r.template(data), data.RandomInteger)
}

func (r EmailCommunicationServiceDomainSenderUsernameResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_email_communication_service_domain_sender_username" "test" {
  name                    = "acctest-su-%d"
  email_service_domain_id = azurerm_email_communication_service_domain.test.id
  display_name            = "TFTester2"
}
`, r.template(data), data.RandomInteger)
}

func (r EmailCommunicationServiceDomainSenderUsernameResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-communicationservice-%d"
  location = "%s"
}

resource "azurerm_email_communication_service" "test" {
  name                = "acctest-CommunicationService-%d"
  resource_group_name = azurerm_resource_group.test.name
  data_location       = "United States"
}

resource "azurerm_email_communication_service_domain" "test" {
  name              = "AzureManagedDomain"
  email_service_id  = azurerm_email_communication_service.test.id
  domain_management = "AzureManaged"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
