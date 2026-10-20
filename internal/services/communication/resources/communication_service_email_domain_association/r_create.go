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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.ServiceClient
			domainClient := metadata.Client.Communication.DomainClient

			var model EmailDomainAssociationResourceModel

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			communicationServiceId, err := communicationservices.ParseCommunicationServiceID(model.CommunicationServiceId)
			if err != nil {
				return err
			}

			eMailServiceDomainId, err := domains.ParseDomainID(model.EMailServiceDomainId)
			if err != nil {
				return err
			}

			locks.ByName(communicationServiceId.CommunicationServiceName, "azurerm_communication_service")
			defer locks.UnlockByName(communicationServiceId.CommunicationServiceName, "azurerm_communication_service")

			locks.ByName(eMailServiceDomainId.DomainName, "azurerm_email_communication_service_domain")
			defer locks.UnlockByName(eMailServiceDomainId.DomainName, "azurerm_email_communication_service_domain")

			existingEMailServiceDomain, err := domainClient.Get(ctx, *eMailServiceDomainId)
			if err != nil && !response.WasNotFound(existingEMailServiceDomain.HttpResponse) {
				return fmt.Errorf("checking for the presence of existing %s: %+v", *eMailServiceDomainId, err)
			}

			if response.WasNotFound(existingEMailServiceDomain.HttpResponse) {
				return fmt.Errorf("%s was not found", eMailServiceDomainId)
			}

			existingCommunicationService, err := client.Get(ctx, *communicationServiceId)
			if err != nil && !response.WasNotFound(existingCommunicationService.HttpResponse) {
				return fmt.Errorf("checking for the presence of existing  %s: %+v", communicationServiceId, err)
			}

			if response.WasNotFound(existingCommunicationService.HttpResponse) {
				return fmt.Errorf("%s was not found", communicationServiceId)
			}

			if existingCommunicationService.Model == nil {
				return fmt.Errorf("model for %s was nil", communicationServiceId)
			}

			if existingCommunicationService.Model.Properties == nil {
				return fmt.Errorf("properties for %s was nil", communicationServiceId)
			}

			domainList := make([]string, 0)
			if existingDomainList := existingCommunicationService.Model.Properties.LinkedDomains; existingDomainList != nil {
				domainList = pointer.From(existingDomainList)
			}

			id := commonids.NewCompositeResourceID(communicationServiceId, eMailServiceDomainId)

			exists := false
			for _, v := range domainList {
				tmpID, tmpErr := domains.ParseDomainIDInsensitively(v)
				if tmpErr != nil {
					return fmt.Errorf("parsing domain ID %q from LinkedDomains for %s: %+v", v, communicationServiceId, err)
				}

				if strings.EqualFold(eMailServiceDomainId.ID(), tmpID.ID()) {
					exists = true
					if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
						return metadata.ResourceRequiresImport(r.ResourceType(), id)
					}
				}
			}

			if !exists {
				domainList = append(domainList, eMailServiceDomainId.ID())
			}

			input := communicationservices.CommunicationServiceResourceUpdate{
				Properties: &communicationservices.CommunicationServiceUpdateProperties{
					LinkedDomains: pointer.To(domainList),
				},
			}

			if _, err = client.Update(ctx, *communicationServiceId, input); err != nil {
				return fmt.Errorf("updating %s: %+v", *communicationServiceId, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
