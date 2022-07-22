package search

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2020-08-01/SharedPrivateLinkResources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2020-08-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"time"
)

type SearchSharedPrivateLinkServiceResource struct{}

var (
	_ sdk.Resource           = SearchSharedPrivateLinkServiceResource{}
	_ sdk.ResourceWithUpdate = SearchSharedPrivateLinkServiceResource{}
)

type SearchSharedPrivateLinkServiceModel struct {
	Name             string `tfschema:"name"`
	SearchServiceId  string `tfschema:"search_service_id"`
	SubResourceName  string `tfschema:"subresource_name"`
	TargetResourceId string `tfschema:"target_resource_id"`
	RequestMessage   string `tfschema:"request_message"`
	Status           string `tfschema:"status"`
}

func (r SearchSharedPrivateLinkServiceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsEmpty,
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

		"status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r SearchSharedPrivateLinkServiceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r SearchSharedPrivateLinkServiceResource) ResourceType() string {
	return "azurerm_search_shared_private_link_resource"
}

func (r SearchSharedPrivateLinkServiceResource) ModelObject() interface{} {
	return &SearchSharedPrivateLinkServiceModel{}
}

func (r SearchSharedPrivateLinkServiceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sharedprivatelinkresources.ValidateSharedPrivateLinkResourceID
}

func (r SearchSharedPrivateLinkServiceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SearchSharedPrivateLinkServiceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Search.SearchSharedPrivateLinkResourceClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			searchServiceId, err := services.ParseSearchServiceID(model.SearchServiceId)
			if err != nil {
				return err
			}

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
					GroupId:               utils.String(model.SubResourceName),
					PrivateLinkResourceId: utils.String(model.TargetResourceId),
				},
			}

			if model.RequestMessage != "" {
				parameters.Properties.RequestMessage = utils.String(model.RequestMessage)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, parameters, sharedprivatelinkresources.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("creating/ updating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r SearchSharedPrivateLinkServiceResource) Read() sdk.ResourceFunc {
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

			state := &SearchSharedPrivateLinkServiceModel{
				Name:            id.SharedPrivateLinkResourceName,
				SearchServiceId: id.SearchServiceName,
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

func (r SearchSharedPrivateLinkServiceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Search.SearchSharedPrivateLinkResourceClient
			id, err := sharedprivatelinkresources.ParseSharedPrivateLinkResourceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, sharedprivatelinkresources.DeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
