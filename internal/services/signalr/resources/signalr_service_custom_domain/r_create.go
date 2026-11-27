// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_service_custom_domain

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	signalrClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var customDomainSignalrServiceModel CustomDomainSignalrServiceModel
			if err := metadata.Decode(&customDomainSignalrServiceModel); err != nil {
				return err
			}
			client := signalrClient.From(metadata.Client).SignalRClient

			signalRServiceId, err := signalr.ParseSignalRIDInsensitively(metadata.ResourceData.Get("signalr_service_id").(string))
			if err != nil {
				return fmt.Errorf("parsing signalr service id error: %+v", err)
			}

			id := signalr.NewCustomDomainID(signalRServiceId.SubscriptionId, signalRServiceId.ResourceGroupName, signalRServiceId.SignalRName, metadata.ResourceData.Get("name").(string))

			locks.ByID(signalRServiceId.ID())
			defer locks.UnlockByID(signalRServiceId.ID())

			if _, err := signalr.ParseCustomCertificateIDInsensitively(customDomainSignalrServiceModel.SignalrCustomCertificateId); err != nil {
				return fmt.Errorf("parsing custom certificate for %s: %+v", id, err)
			}

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.CustomDomainsGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			customDomainObj := signalr.CustomDomain{
				Properties: signalr.CustomDomainProperties{
					DomainName: customDomainSignalrServiceModel.DomainName,
					CustomCertificate: signalr.ResourceReference{
						Id: pointer.To(customDomainSignalrServiceModel.SignalrCustomCertificateId),
					},
				},
			}
			if err := client.CustomDomainsCreateOrUpdateCallbackThenPoll(ctx, id, customDomainObj, metadata.SetIDAndIdentityCallback(&id)); err != nil {
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
