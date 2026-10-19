// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package email_communication_service_domain

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/domains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/emailservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.Communication.DomainClient

			var model EmailCommunicationServiceDomainResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			eMailServiceID, err := emailservices.ParseEmailServiceID(model.EMailServiceID)
			if err != nil {
				return fmt.Errorf("parsing parent email_service_id: %+v", err)
			}

			id := domains.NewDomainID(subscriptionId, eMailServiceID.ResourceGroupName, eMailServiceID.EmailServiceName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := &domains.DomainProperties{
				DomainManagement: domains.DomainManagement(model.DomainManagement),
			}

			properties.UserEngagementTracking = pointer.To(domains.UserEngagementTrackingDisabled)
			if model.UserEngagementTrackingEnabled {
				properties.UserEngagementTracking = pointer.To(domains.UserEngagementTrackingEnabled)
			}

			param := domains.DomainResource{
				// The location is always `global` from the Azure Portal
				Location:   location.Normalize("global"),
				Properties: properties,
				Tags:       pointer.To(model.Tags),
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, param, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}
