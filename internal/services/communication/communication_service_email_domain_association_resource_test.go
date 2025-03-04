// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package communication_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/communicationservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/domains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CommunicationServiceEmailDomainAssociationResource struct{}

func TestAccCommunicationServiceEmailDomainAssociationResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_communication_service_email_domain_association", "test")
	r := CommunicationServiceEmailDomainAssociationResource{}
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

func TestAccCommunicationServiceEmailDomainAssociationResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_communication_service_email_domain_association", "test")
	r := CommunicationServiceEmailDomainAssociationResource{}
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

func (r CommunicationServiceEmailDomainAssociationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseCompositeResourceID(state.ID, &communicationservices.CommunicationServiceId{}, &domains.DomainId{})
	if err != nil {
		return pointer.To(false), fmt.Errorf("parsing ID: %w", err)
	}

	serviceClient := client.Communication.ServiceClient
	existingCommunicationService, err := serviceClient.Get(ctx, *id.First)
	if err != nil && !response.WasNotFound(existingCommunicationService.HttpResponse) {
		return pointer.To(false), fmt.Errorf("checking for the presence of existing %s: %+v", id.First, err)
	}

	if response.WasNotFound(existingCommunicationService.HttpResponse) {
		return pointer.To(false), fmt.Errorf("%s does not exist", id.First)
	}

	input := existingCommunicationService
	if input.Model != nil && input.Model.Properties != nil {
		for _, v := range pointer.From(input.Model.Properties.LinkedDomains) {
			tmpID, tmpErr := domains.ParseDomainID(v)
			if tmpErr != nil {
				return pointer.To(false), fmt.Errorf("parsing domain ID %q from LinkedDomains for %s: %+v", v, id.First, err)
			}
			if strings.EqualFold(id.Second.ID(), tmpID.ID()) {
				return pointer.To(true), nil
			}
		}
	}

	return pointer.To(false), nil
}

func (r CommunicationServiceEmailDomainAssociationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_communication_service_email_domain_association" "test" {
  communication_service_id = azurerm_communication_service.test.id
  email_service_domain_id  = azurerm_email_communication_service_domain.test.id
}
`, r.template(data))
}

func (r CommunicationServiceEmailDomainAssociationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-communicationservice-%[1]d"
  location = "%[2]s"
}

resource "azurerm_communication_service" "test" {
  name                = "acctest-CommunicationService-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  data_location       = "United States"

  tags = {
    env = "Test2"
  }
}

resource "azurerm_email_communication_service" "test" {
  name                = "acctest-CommunicationService-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  data_location       = "United States"
}

resource "azurerm_email_communication_service_domain" "test" {
  name             = "AzureManagedDomain"
  email_service_id = azurerm_email_communication_service.test.id

  domain_management = "AzureManaged"
}

`, data.RandomInteger, data.Locations.Primary)
}

func (r CommunicationServiceEmailDomainAssociationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_communication_service_email_domain_association" "import" {
  communication_service_id = azurerm_communication_service_email_domain_association.test.communication_service_id
  email_service_domain_id  = azurerm_communication_service_email_domain_association.test.email_service_domain_id
}
`, r.basic(data))
}
