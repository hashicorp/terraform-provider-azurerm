// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package communication_service_email_domain_association

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/communicationservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/domains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.ServiceClient
			domainClient := metadata.Client.Communication.DomainClient

			var model EmailDomainAssociationResourceModel

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id, err := commonids.ParseCompositeResourceID(metadata.ResourceData.Id(), &communicationservices.CommunicationServiceId{}, &domains.DomainId{})
			if err != nil {
				return err
			}

			communicationServiceId := id.First
			eMailServiceDomainId := id.Second

			locks.ByName(communicationServiceId.CommunicationServiceName, "azurerm_communication_service")
			defer locks.UnlockByName(communicationServiceId.CommunicationServiceName, "azurerm_communication_service")

			locks.ByName(eMailServiceDomainId.DomainName, "azurerm_email_communication_service_domain")
			defer locks.UnlockByName(eMailServiceDomainId.DomainName, "azurerm_email_communication_service_domain")

			existingEMailServiceDomain, err := domainClient.Get(ctx, *eMailServiceDomainId)
			if err != nil && !response.WasNotFound(existingEMailServiceDomain.HttpResponse) {
				return fmt.Errorf("checking for the presence of existing %s: %+v", *eMailServiceDomainId, err)
			}

			if response.WasNotFound(existingEMailServiceDomain.HttpResponse) {
				return metadata.MarkAsGone(id)
			}

			existingCommunicationService, err := client.Get(ctx, *communicationServiceId)
			if err != nil && !response.WasNotFound(existingCommunicationService.HttpResponse) {
				return fmt.Errorf("checking for the presence of existing %s: %+v", *communicationServiceId, err)
			}

			if response.WasNotFound(existingCommunicationService.HttpResponse) {
				return metadata.MarkAsGone(id)
			}

			if existingCommunicationService.Model == nil {
				return fmt.Errorf("model for %s was nil", model.CommunicationServiceId)
			}

			if existingCommunicationService.Model.Properties == nil {
				return fmt.Errorf("properties for %s was nil", model.CommunicationServiceId)
			}

			domainList := existingCommunicationService.Model.Properties.LinkedDomains
			if domainList == nil {
				return metadata.MarkAsGone(id)
			}

			if !slices.Contains(*domainList, eMailServiceDomainId.ID()) {
				return metadata.MarkAsGone(id)
			}

			*domainList = slices.DeleteFunc(*domainList, func(domainID string) bool {
				parsedDomainID, err := domains.ParseDomainIDInsensitively(domainID)
				if err != nil {
					return false
				}

				return strings.EqualFold(parsedDomainID.ID(), eMailServiceDomainId.ID())
			})

			input := communicationservices.CommunicationServiceResourceUpdate{
				Properties: &communicationservices.CommunicationServiceUpdateProperties{
					LinkedDomains: domainList,
				},
			}

			if _, err := client.Update(ctx, *communicationServiceId, input); err != nil {
				return fmt.Errorf("deleting Email Domain Association for %s from %s: %+v", *eMailServiceDomainId, *communicationServiceId, err)
			}

			return nil
		},
	}
}
