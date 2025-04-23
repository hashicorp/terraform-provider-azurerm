// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

type SpringCloudApplicationLiveViewModel struct {
	Name                 string `tfschema:"name"`
	SpringCloudServiceId string `tfschema:"spring_cloud_service_id"`
}

type SpringCloudApplicationLiveViewResource struct{}

func (s SpringCloudApplicationLiveViewResource) DeprecationMessage() string {
	return features.DeprecatedInFivePointOh("Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_application_live_view` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.")
}

var (
	_ sdk.Resource                                = SpringCloudApplicationLiveViewResource{}
	_ sdk.ResourceWithDeprecationAndNoReplacement = SpringCloudApplicationLiveViewResource{}
)

func (s SpringCloudApplicationLiveViewResource) ResourceType() string {
	return "azurerm_spring_cloud_application_live_view"
}

func (s SpringCloudApplicationLiveViewResource) ModelObject() interface{} {
	return &SpringCloudApplicationLiveViewModel{}
}

func (s SpringCloudApplicationLiveViewResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.SpringCloudApplicationLiveViewID
}

func (s SpringCloudApplicationLiveViewResource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"default"}, false),
		},

		"spring_cloud_service_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SpringCloudServiceID,
		},
	}
}

func (s SpringCloudApplicationLiveViewResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (s SpringCloudApplicationLiveViewResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudApplicationLiveViewModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AppPlatform.ApplicationLiveViewsClient
			springId, err := parse.SpringCloudServiceID(model.SpringCloudServiceId)
			if err != nil {
				return fmt.Errorf("parsing spring service ID: %+v", err)
			}
			id := parse.NewSpringCloudApplicationLiveViewID(springId.SubscriptionId, springId.ResourceGroup, springId.SpringName, model.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.ApplicationLiveViewName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(s.ResourceType(), id)
			}

			applicationLiveViewResource := appplatform.ApplicationLiveViewResource{}
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.ApplicationLiveViewName, applicationLiveViewResource)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (s SpringCloudApplicationLiveViewResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.ApplicationLiveViewsClient

			id, err := parse.SpringCloudApplicationLiveViewID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.ApplicationLiveViewName)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			state := SpringCloudApplicationLiveViewModel{
				Name:                 id.ApplicationLiveViewName,
				SpringCloudServiceId: parse.NewSpringCloudServiceID(id.SubscriptionId, id.ResourceGroup, id.SpringName).ID(),
			}
			return metadata.Encode(&state)
		},
	}
}

func (s SpringCloudApplicationLiveViewResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.ApplicationLiveViewsClient

			id, err := parse.SpringCloudApplicationLiveViewID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.SpringName, id.ApplicationLiveViewName)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				if !response.WasNotFound(future.Response()) {
					return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}
