// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package search

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2024-06-01-preview/services"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2024-06-01-preview/sharedprivatelinkresources"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SharedPrivateLinkServiceResource struct{}

var (
	_ sdk.Resource           = SharedPrivateLinkServiceResource{}
	_ sdk.ResourceWithUpdate = SharedPrivateLinkServiceResource{}
)

type SharedPrivateLinkServiceModel struct {
	Name             string `tfschema:"name"`
	SearchServiceId  string `tfschema:"search_service_id"`
	SubResourceName  string `tfschema:"subresource_name"`
	TargetResourceId string `tfschema:"target_resource_id"`
	RequestMessage   string `tfschema:"request_message"`
	Status           string `tfschema:"status"`
}

func (r SharedPrivateLinkServiceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"search_service_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: services.ValidateSearchServiceID,
		},

		"subresource_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: networkValidate.PrivateLinkSubResourceName,
		},

		"target_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"request_message": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r SharedPrivateLinkServiceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r SharedPrivateLinkServiceResource) ResourceType() string {
	return "azurerm_search_shared_private_link_service"
}

func (r SharedPrivateLinkServiceResource) ModelObject() interface{} {
	return &SharedPrivateLinkServiceModel{}
}

func (r SharedPrivateLinkServiceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sharedprivatelinkresources.ValidateSharedPrivateLinkResourceID
}

func (r SharedPrivateLinkServiceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SharedPrivateLinkServiceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Search.SearchSharedPrivateLinkResourceClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			searchServiceId, err := services.ParseSearchServiceID(model.SearchServiceId)
			if err != nil {
				return err
			}

			locks.ByID(searchServiceId.ID())
			defer locks.UnlockByID(searchServiceId.ID())

			id := sharedprivatelinkresources.NewSharedPrivateLinkResourceID(subscriptionId, searchServiceId.ResourceGroupName, searchServiceId.SearchServiceName, model.Name)

			existing, err := client.Get(ctx, id, sharedprivatelinkresources.GetOperationOptions{})
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing shared private link resource %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := sharedprivatelinkresources.SharedPrivateLinkResource{
				Properties: &sharedprivatelinkresources.SharedPrivateLinkResourceProperties{
					GroupId:               pointer.To(model.SubResourceName),
					PrivateLinkResourceId: pointer.To(model.TargetResourceId),
				},
			}

			if model.RequestMessage != "" {
				parameters.Properties.RequestMessage = pointer.To(model.RequestMessage)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, parameters, sharedprivatelinkresources.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 60 * time.Minute,
	}
}

func (r SharedPrivateLinkServiceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Search.SearchSharedPrivateLinkResourceClient
			id, err := sharedprivatelinkresources.ParseSharedPrivateLinkResourceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id, sharedprivatelinkresources.GetOperationOptions{})
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					metadata.Logger.Infof("%q was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := &SharedPrivateLinkServiceModel{
				Name:            id.SharedPrivateLinkResourceName,
				SearchServiceId: services.NewSearchServiceID(id.SubscriptionId, id.ResourceGroupName, id.SearchServiceName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if props.GroupId != nil {
						state.SubResourceName = *props.GroupId
					}

					if props.PrivateLinkResourceId != nil {
						state.TargetResourceId = *props.PrivateLinkResourceId
					}

					if props.RequestMessage != nil {
						state.RequestMessage = *props.RequestMessage
					}

					if props.Status != nil {
						state.Status = string(*props.Status)
					}
				}
			}

			return metadata.Encode(state)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r SharedPrivateLinkServiceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := sharedprivatelinkresources.ParseSharedPrivateLinkResourceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			searchServiceId := sharedprivatelinkresources.NewSearchServiceID(id.SubscriptionId, id.ResourceGroupName, id.SearchServiceName)
			locks.ByID(searchServiceId.ID())
			defer locks.UnlockByID(searchServiceId.ID())

			var state SharedPrivateLinkServiceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Search.SearchSharedPrivateLinkResourceClient

			if metadata.ResourceData.HasChange("request_message") {
				props := sharedprivatelinkresources.SharedPrivateLinkResource{
					Properties: &sharedprivatelinkresources.SharedPrivateLinkResourceProperties{
						RequestMessage: pointer.To(state.RequestMessage),
					},
				}
				if err := client.CreateOrUpdateThenPoll(ctx, *id, props, sharedprivatelinkresources.CreateOrUpdateOperationOptions{}); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}
			return nil
		},
		Timeout: 60 * time.Minute,
	}
}

func (r SharedPrivateLinkServiceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Search.SearchSharedPrivateLinkResourceClient
			id, err := sharedprivatelinkresources.ParseSharedPrivateLinkResourceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			searchServiceId := sharedprivatelinkresources.NewSearchServiceID(id.SubscriptionId, id.ResourceGroupName, id.SearchServiceName)
			locks.ByID(searchServiceId.ID())
			defer locks.UnlockByID(searchServiceId.ID())

			if err := client.DeleteThenPoll(ctx, *id, sharedprivatelinkresources.DeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}
			return nil
		},
		Timeout: 60 * time.Minute,
	}
}
