// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package communication_service_email_domain_association

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/communicationservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/domains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.ServiceClient
			domainClient := metadata.Client.Communication.DomainClient

			id, err := commonids.ParseCompositeResourceID(metadata.ResourceData.Id(), &communicationservices.CommunicationServiceId{}, &domains.DomainId{})
			if err != nil {
				return err
			}

			state := EmailDomainAssociationResourceModel{}
			state.CommunicationServiceId = id.First.ID()
			state.EMailServiceDomainId = id.Second.ID()

			communicationServiceId, err := communicationservices.ParseCommunicationServiceID(state.CommunicationServiceId)
			if err != nil {
				return err
			}

			eMailServiceDomainId, err := domains.ParseDomainID(state.EMailServiceDomainId)
			if err != nil {
				return err
			}

			locks.ByName(communicationServiceId.CommunicationServiceName, "azurerm_communication_service")
			defer locks.UnlockByName(communicationServiceId.CommunicationServiceName, "azurerm_communication_service")

			locks.ByName(eMailServiceDomainId.DomainName, "azurerm_email_communication_service_domain")
			defer locks.UnlockByName(eMailServiceDomainId.DomainName, "azurerm_email_communication_service_domain")

			existingEMailServiceDomain, err := domainClient.Get(ctx, *eMailServiceDomainId)
			if err != nil && !response.WasNotFound(existingEMailServiceDomain.HttpResponse) {
				return fmt.Errorf("checking for the presence of existing EMail Service Domain %q: %+v", state.EMailServiceDomainId, err)
			}

			if response.WasNotFound(existingEMailServiceDomain.HttpResponse) {
				return metadata.MarkAsGone(id)
			}

			existingCommunicationService, err := client.Get(ctx, *communicationServiceId)
			if err != nil && !response.WasNotFound(existingCommunicationService.HttpResponse) {
				return fmt.Errorf("checking for the presence of existing Communication Service %q: %+v", state.CommunicationServiceId, err)
			}

			if response.WasNotFound(existingCommunicationService.HttpResponse) {
				return metadata.MarkAsGone(id)
			}

			if existingCommunicationService.Model == nil {
				return fmt.Errorf("model for %s was nil", state.CommunicationServiceId)
			}

			if existingCommunicationService.Model.Properties == nil {
				return fmt.Errorf("properties for %s was nil", state.CommunicationServiceId)
			}

			domainList := existingCommunicationService.Model.Properties.LinkedDomains
			if domainList == nil {
				return fmt.Errorf("checking for Domain Association %s for %s", *eMailServiceDomainId, *communicationServiceId)
			}

			var found bool

			for _, v := range pointer.From(domainList) {
				tmpID, tmpErr := domains.ParseDomainIDInsensitively(v)
				if tmpErr != nil {
					return fmt.Errorf("parsing domain ID %q from LinkedDomains for %s: %+v", v, communicationServiceId, err)
				}

				if strings.EqualFold(eMailServiceDomainId.ID(), tmpID.ID()) {
					found = true

					break
				}
			}

			if !found {
				return metadata.MarkAsGone(id)
			}

			return metadata.Encode(&state)
		},
	}
}
