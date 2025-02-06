// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package communication

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/communicationservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/domains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = EmailDomainAssociationResource{}

type EmailDomainAssociationResource struct{}

type EmailDomainAssociationResourceModel struct {
	CommunicationServiceId string `tfschema:"communication_service_id"`
	EMailServiceDomainId   string `tfschema:"email_service_domain_id"`
}

func (EmailDomainAssociationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"communication_service_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: communicationservices.ValidateCommunicationServiceID,
		},
		"email_service_domain_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: domains.ValidateDomainID,
		},
	}
}

func (EmailDomainAssociationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (EmailDomainAssociationResource) ModelObject() interface{} {
	return &EmailDomainAssociationResourceModel{}
}

func (EmailDomainAssociationResource) ResourceType() string {
	return "azurerm_communication_service_email_domain_association"
}

func (r EmailDomainAssociationResource) Create() sdk.ResourceFunc {
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

			for _, v := range domainList {
				tmpID, tmpErr := domains.ParseDomainIDInsensitively(v)
				if tmpErr != nil {
					return fmt.Errorf("parsing domain ID %q from LinkedDomains for %s: %+v", v, communicationServiceId, err)
				}

				if strings.EqualFold(eMailServiceDomainId.ID(), tmpID.ID()) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			domainList = append(domainList, eMailServiceDomainId.ID())

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

func (EmailDomainAssociationResource) Read() sdk.ResourceFunc {
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

func (EmailDomainAssociationResource) Delete() sdk.ResourceFunc {
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

func (EmailDomainAssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(input interface{}, key string) (warnings []string, errors []error) {
		v, ok := input.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected %q to be a string", key))
			return
		}

		if _, err := commonids.ParseCompositeResourceID(v, &communicationservices.CommunicationServiceId{}, &domains.DomainId{}); err != nil {
			errors = append(errors, err)
		}

		return
	}
}
