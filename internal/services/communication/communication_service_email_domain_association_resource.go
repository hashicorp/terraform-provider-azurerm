// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package communication

import (
	"context"
	"fmt"
	"log"
	"slices"
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

var _ sdk.Resource = CommunicationServiceEmailDomainAssociationResource{}

type CommunicationServiceEmailDomainAssociationResource struct{}

type CommunicationServiceEmailDomainAssociationResourceModel struct {
	CommunicationServiceId string `tfschema:"communication_service_id"`
	EMailServiceDomainId   string `tfschema:"email_service_domain_id"`
}

func (CommunicationServiceEmailDomainAssociationResource) Arguments() map[string]*pluginsdk.Schema {
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

func (CommunicationServiceEmailDomainAssociationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (CommunicationServiceEmailDomainAssociationResource) ModelObject() interface{} {
	return &CommunicationServiceEmailDomainAssociationResourceModel{}
}

func (CommunicationServiceEmailDomainAssociationResource) ResourceType() string {
	return "azurerm_communication_service_email_domain_association"
}

func (r CommunicationServiceEmailDomainAssociationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.ServiceClient
			domainClient := metadata.Client.Communication.DomainClient

			var model CommunicationServiceEmailDomainAssociationResourceModel

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			communicationServiceId, err := communicationservices.ParseCommunicationServiceID(model.CommunicationServiceId)
			if err != nil {
				return fmt.Errorf("parsing Communication Service ID: %w", err)
			}

			eMailServiceDomainId, err := domains.ParseDomainID(model.EMailServiceDomainId)
			if err != nil {
				return fmt.Errorf("parsing EMail Service Domain ID: %w", err)
			}

			locks.ByName(communicationServiceId.CommunicationServiceName, "azurerm_communication_service")
			defer locks.UnlockByName(communicationServiceId.CommunicationServiceName, "azurerm_communication_service")

			locks.ByName(eMailServiceDomainId.DomainName, "azurerm_email_communication_service_domain")
			defer locks.UnlockByName(eMailServiceDomainId.DomainName, "azurerm_email_communication_service_domain")

			existingEMailServiceDomain, err := domainClient.Get(ctx, *eMailServiceDomainId)
			if err != nil && !response.WasNotFound(existingEMailServiceDomain.HttpResponse) {
				return fmt.Errorf("checking for the presence of existing EMail Service Domain %q: %+v", model.EMailServiceDomainId, err)
			}

			if response.WasNotFound(existingEMailServiceDomain.HttpResponse) {
				return fmt.Errorf("EMail Service Domain %q does not exsits", model.EMailServiceDomainId)
			}

			existingCommunicationService, err := client.Get(ctx, *communicationServiceId)
			if err != nil && !response.WasNotFound(existingCommunicationService.HttpResponse) {
				return fmt.Errorf("checking for the presence of existing Communication Service %q: %+v", model.CommunicationServiceId, err)
			}

			if response.WasNotFound(existingCommunicationService.HttpResponse) {
				return fmt.Errorf("Communication Service %q does not exists", model.CommunicationServiceId)
			}

			if existingCommunicationService.Model == nil || existingCommunicationService.Model.Properties == nil {
				return fmt.Errorf("model/properties for %s was nil", model.CommunicationServiceId)
			}

			domainList := existingCommunicationService.Model.Properties.LinkedDomains
			if domainList == nil {
				domainList = pointer.FromSliceOfStrings(make([]string, 0, 1))
			}

			id := commonids.NewCompositeResourceID(communicationServiceId, eMailServiceDomainId)
			if slices.Contains(*domainList, eMailServiceDomainId.ID()) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			*domainList = append(*domainList, eMailServiceDomainId.ID())
			existingCommunicationService.Model.Properties.LinkedDomains = domainList

			input := communicationservices.CommunicationServiceResourceUpdate{
				Properties: &communicationservices.CommunicationServiceUpdateProperties{
					LinkedDomains: domainList,
				},
			}

			if _, err := client.Update(ctx, *communicationServiceId, input); err != nil {
				return fmt.Errorf("updating %s: %+v", *communicationServiceId, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (CommunicationServiceEmailDomainAssociationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.ServiceClient
			domainClient := metadata.Client.Communication.DomainClient

			id, err := commonids.ParseCompositeResourceID(metadata.ResourceData.Id(), &communicationservices.CommunicationServiceId{}, &domains.DomainId{})
			if err != nil {
				return err
			}

			state := CommunicationServiceEmailDomainAssociationResourceModel{}
			state.CommunicationServiceId = id.First.ID()
			state.EMailServiceDomainId = id.Second.ID()

			communicationServiceId, err := communicationservices.ParseCommunicationServiceID(state.CommunicationServiceId)
			if err != nil {
				return fmt.Errorf("parsing Communication Service ID: %w", err)
			}

			eMailServiceDomainId, err := domains.ParseDomainID(state.EMailServiceDomainId)
			if err != nil {
				return fmt.Errorf("parsing EMail Service Domain ID: %w", err)
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
				return fmt.Errorf("EMail Service Domain %q does not exsits", state.EMailServiceDomainId)
			}

			existingCommunicationService, err := client.Get(ctx, *communicationServiceId)
			if err != nil && !response.WasNotFound(existingCommunicationService.HttpResponse) {
				return fmt.Errorf("checking for the presence of existing Communication Service %q: %+v", state.CommunicationServiceId, err)
			}

			if response.WasNotFound(existingCommunicationService.HttpResponse) {
				return fmt.Errorf("Communication Service %q does not exsits", state.CommunicationServiceId)
			}

			if existingCommunicationService.Model == nil || existingCommunicationService.Model.Properties == nil {
				return fmt.Errorf("model/properties for %s was nil", state.CommunicationServiceId)
			}

			domainList := existingCommunicationService.Model.Properties.LinkedDomains
			if domainList == nil {
				domainList = pointer.FromSliceOfStrings(make([]string, 0, 1))
			}

			if !slices.Contains(*domainList, eMailServiceDomainId.ID()) {
				log.Printf("EMail Service Domain %q does not exsits in %q, removing from state.", eMailServiceDomainId, communicationServiceId)
				err := metadata.MarkAsGone(id)
				if err != nil {
					return err
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (CommunicationServiceEmailDomainAssociationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.ServiceClient
			domainClient := metadata.Client.Communication.DomainClient

			var model CommunicationServiceEmailDomainAssociationResourceModel

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			communicationServiceId, err := communicationservices.ParseCommunicationServiceID(model.CommunicationServiceId)
			if err != nil {
				return fmt.Errorf("parsing Communication Service ID: %w", err)
			}

			eMailServiceDomainId, err := domains.ParseDomainID(model.EMailServiceDomainId)
			if err != nil {
				return fmt.Errorf("parsing EMail Service Domain ID: %w", err)
			}

			locks.ByName(communicationServiceId.CommunicationServiceName, "azurerm_communication_service")
			defer locks.UnlockByName(communicationServiceId.CommunicationServiceName, "azurerm_communication_service")

			locks.ByName(eMailServiceDomainId.DomainName, "azurerm_email_communication_service_domain")
			defer locks.UnlockByName(eMailServiceDomainId.DomainName, "azurerm_email_communication_service_domain")

			existingEMailServiceDomain, err := domainClient.Get(ctx, *eMailServiceDomainId)
			if err != nil && !response.WasNotFound(existingEMailServiceDomain.HttpResponse) {
				return fmt.Errorf("checking for the presence of existing EMail Service Domain %q: %+v", model.EMailServiceDomainId, err)
			}

			if response.WasNotFound(existingEMailServiceDomain.HttpResponse) {
				return fmt.Errorf("EMail Service Domain %q does not exsits", model.EMailServiceDomainId)
			}

			existingCommunicationService, err := client.Get(ctx, *communicationServiceId)
			if err != nil && !response.WasNotFound(existingCommunicationService.HttpResponse) {
				return fmt.Errorf("checking for the presence of existing Communication Service %q: %+v", model.CommunicationServiceId, err)
			}

			if response.WasNotFound(existingCommunicationService.HttpResponse) {
				return fmt.Errorf("Communication Service %q does not exsits", model.CommunicationServiceId)
			}

			if existingCommunicationService.Model == nil || existingCommunicationService.Model.Properties == nil {
				return fmt.Errorf("model/properties for %s was nil", model.CommunicationServiceId)
			}

			domainList := existingCommunicationService.Model.Properties.LinkedDomains
			if domainList == nil {
				domainList = pointer.FromSliceOfStrings(make([]string, 0, 1))
			}

			id := commonids.NewCompositeResourceID(communicationServiceId, eMailServiceDomainId)
			metadata.SetID(id)

			if !slices.Contains(*domainList, eMailServiceDomainId.ID()) {
				return nil
			}

			*domainList = slices.DeleteFunc(*domainList, func(n string) bool {
				return n == eMailServiceDomainId.ID()
			})

			existingCommunicationService.Model.Properties.LinkedDomains = domainList

			input := communicationservices.CommunicationServiceResourceUpdate{
				Properties: &communicationservices.CommunicationServiceUpdateProperties{
					LinkedDomains: domainList,
				},
			}

			if _, err := client.Update(ctx, *communicationServiceId, input); err != nil {
				return fmt.Errorf("updating %s: %+v", *communicationServiceId, err)
			}

			return nil
		},
	}
}

func (CommunicationServiceEmailDomainAssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
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
