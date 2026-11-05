// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package email_communication_service_domain_sender_username

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/senderusernames"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.Communication.SenderUsernamesClient

			var model EmailCommunicationServiceDomainSenderUsernameModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			eMailServiceDomainID, err := senderusernames.ParseDomainID(model.EMailServiceDomainID)
			if err != nil {
				return err
			}

			id := senderusernames.NewSenderUsernameID(subscriptionId, eMailServiceDomainID.ResourceGroupName, eMailServiceDomainID.EmailServiceName, eMailServiceDomainID.DomainName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := senderusernames.SenderUsernameResource{
				Properties: &senderusernames.SenderUsernameProperties{
					// service API response indicates that the username must be same with the resource name
					Username: model.Name,
				},
			}

			if v := model.DisplayName; v != "" {
				parameters.Properties.DisplayName = pointer.To(v)
			}

			if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
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
