// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/resourcemanagementprivatelink"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = ResourceManagementPrivateLinkResource{}

type ResourceManagementPrivateLinkResource struct{}

func (r ResourceManagementPrivateLinkResource) ModelObject() interface{} {
	return &ResourceManagementPrivateLinkResourceSchema{}
}

type ResourceManagementPrivateLinkResourceSchema struct {
	Location          string `tfschema:"location"`
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`
}

func (r ResourceManagementPrivateLinkResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return resourcemanagementprivatelink.ValidateResourceManagementPrivateLinkID
}

func (r ResourceManagementPrivateLinkResource) ResourceType() string {
	return "azurerm_resource_management_private_link"
}

func (r ResourceManagementPrivateLinkResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			ForceNew:     true,
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"location":            commonschema.Location(),
	}
}

func (r ResourceManagementPrivateLinkResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ResourceManagementPrivateLinkResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.ResourceManagementPrivateLinkClient

			var config ResourceManagementPrivateLinkResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId

			id := resourcemanagementprivatelink.NewResourceManagementPrivateLinkID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := resourcemanagementprivatelink.ResourceManagementPrivateLinkLocation{
				Location: pointer.To(location.Normalize(config.Location)),
			}

			if _, err := client.Put(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ResourceManagementPrivateLinkResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.ResourceManagementPrivateLinkClient

			id, err := resourcemanagementprivatelink.ParseResourceManagementPrivateLinkID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			schema := ResourceManagementPrivateLinkResourceSchema{
				Name:              id.ResourceManagementPrivateLinkName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				schema.Location = location.NormalizeNilable(model.Location)
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r ResourceManagementPrivateLinkResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.ResourceManagementPrivateLinkClient

			id, err := resourcemanagementprivatelink.ParseResourceManagementPrivateLinkID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
