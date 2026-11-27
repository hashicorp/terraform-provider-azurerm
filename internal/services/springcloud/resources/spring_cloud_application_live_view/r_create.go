// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_application_live_view

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func (s Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudApplicationLiveViewModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := springcloudClient.From(metadata.Client).ApplicationLiveViewsClient
			// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
			// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
			springId, err := commonids.ParseSpringCloudServiceIDInsensitively(model.SpringCloudServiceId)
			if err != nil {
				return fmt.Errorf("parsing spring service ID: %+v", err)
			}
			id := appplatform_rm.NewApplicationLiveViewID(springId.SubscriptionId, springId.ResourceGroupName, springId.ServiceName, model.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.ApplicationLiveViewName)
				if err != nil && !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.Response.Response) {
					return metadata.ResourceRequiresImport(s.ResourceType(), id)
				}
			}

			applicationLiveViewResource := appplatform.ApplicationLiveViewResource{}
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.SpringName, id.ApplicationLiveViewName, applicationLiveViewResource)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
			}

			return nil
		},
	}
}
