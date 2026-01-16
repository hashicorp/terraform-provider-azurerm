// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = NetworkSecurityPerimeterResource{}

type NetworkSecurityPerimeterResource struct{}

type NetworkSecurityPerimeterResourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`
}

func (NetworkSecurityPerimeterResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`(^[a-zA-Z0-9]+[a-zA-Z0-9_.-]{0,78}[a-zA-Z0-9_]+$)|(^[a-zA-Z0-9]$)`),
				"`name` must be between 1 and 80 characters long, start with a letter or number, end with a letter, number, or underscore, and may contain only letters, numbers, underscores (_), periods (.), or hyphens (-).",
			),
			ForceNew:     true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (NetworkSecurityPerimeterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (NetworkSecurityPerimeterResource) ModelObject() interface{} {
	return &NetworkSecurityPerimeterResourceModel{}
}

func (NetworkSecurityPerimeterResource) ResourceType() string {
	return "azurerm_network_security_perimeter"
}

func (r NetworkSecurityPerimeterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{

		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkSecurityPerimetersClient

			subscriptionId := metadata.Client.Account.SubscriptionId

			var config NetworkSecurityPerimeterResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			id := networksecurityperimeters.NewNetworkSecurityPerimeterID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := networksecurityperimeters.NetworkSecurityPerimeter{
				Location: location.Normalize(config.Location),
				Tags:     pointer.To(config.Tags),
			}
			if _, err := client.CreateOrUpdate(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r NetworkSecurityPerimeterResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkSecurityPerimetersClient

			id, err := networksecurityperimeters.ParseNetworkSecurityPerimeterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config NetworkSecurityPerimeterResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("tags") {
				tagReq := networksecurityperimeters.UpdateTagsRequest{
					Id:   pointer.To(id.ID()),
					Tags: pointer.To(config.Tags),
				}

				if _, err := client.Patch(ctx, *id, tagReq); err != nil {
					return fmt.Errorf("updating tags for %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func (NetworkSecurityPerimeterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkSecurityPerimetersClient

			id, err := networksecurityperimeters.ParseNetworkSecurityPerimeterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := NetworkSecurityPerimeterResourceModel{
				Name:              id.NetworkSecurityPerimeterName,
				ResourceGroupName: id.ResourceGroupName,
			}
			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
			}
			return metadata.Encode(&state)
		},
	}
}

func (NetworkSecurityPerimeterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkSecurityPerimetersClient

			id, err := networksecurityperimeters.ParseNetworkSecurityPerimeterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, networksecurityperimeters.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (NetworkSecurityPerimeterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return networksecurityperimeters.ValidateNetworkSecurityPerimeterID
}
