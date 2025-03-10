// Copyright (c) HashiCorp, Inc.
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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/ipampools"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = ManagerIpamPoolResource{}

type ManagerIpamPoolResource struct{}

func (ManagerIpamPoolResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return ipampools.ValidateIPamPoolID
}

func (ManagerIpamPoolResource) ResourceType() string {
	return "azurerm_network_manager_ipam_pool"
}

func (ManagerIpamPoolResource) ModelObject() interface{} {
	return &ManagerIpamPoolResourceModel{}
}

type ManagerIpamPoolResourceModel struct {
	AddressPrefixes  []string          `tfschema:"address_prefixes"`
	Description      string            `tfschema:"description"`
	DisplayName      string            `tfschema:"display_name"`
	Location         string            `tfschema:"location"`
	Name             string            `tfschema:"name"`
	NetworkManagerId string            `tfschema:"network_manager_id"`
	ParentPoolName   string            `tfschema:"parent_pool_name"`
	Tags             map[string]string `tfschema:"tags"`
}

func (ManagerIpamPoolResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9\_\.\-]{1,64}$`),
				"`name` must be between 1 and 64 characters long and can only contain letters, numbers, underscores(_), periods(.), and hyphens(-).",
			),
		},

		"network_manager_id": commonschema.ResourceIDReferenceRequiredForceNew(&ipampools.NetworkManagerId{}),

		"location": commonschema.Location(),

		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9\_\.\-]{1,64}$`),
				"`display_name` must be between 1 and 64 characters long and can only contain letters, numbers, underscores(_), periods(.), and hyphens(-).",
			),
		},

		"address_prefixes": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsCIDR,
			},
		},

		"parent_pool_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9\_\.\-]{1,64}$`),
				"`parent_pool_name` must be between 1 and 64 characters long and can only contain letters, numbers, underscores(_), periods(.), and hyphens(-).",
			),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (ManagerIpamPoolResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagerIpamPoolResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.IPamPools
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ManagerIpamPoolResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			networkManagerId, err := ipampools.ParseNetworkManagerID(config.NetworkManagerId)
			if err != nil {
				return err
			}

			id := ipampools.NewIPamPoolID(subscriptionId, networkManagerId.ResourceGroupName, networkManagerId.NetworkManagerName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := ipampools.IPamPool{
				Name:     pointer.To(config.Name),
				Location: location.Normalize(config.Location),
				Tags:     pointer.To(config.Tags),
				Properties: ipampools.IPamPoolProperties{
					AddressPrefixes: config.AddressPrefixes,
					Description:     pointer.To(config.Description),
					DisplayName:     pointer.To(config.DisplayName),
					ParentPoolName:  pointer.To(config.ParentPoolName),
				},
			}

			if err := client.CreateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ManagerIpamPoolResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.IPamPools

			id, err := ipampools.ParseIPamPoolID(metadata.ResourceData.Id())
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

			networkManagerId := ipampools.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName).ID()
			schema := ManagerIpamPoolResourceModel{
				Name:             id.IpamPoolName,
				NetworkManagerId: networkManagerId,
			}

			if model := resp.Model; model != nil {
				schema.Location = location.Normalize(model.Location)
				schema.Tags = pointer.From(model.Tags)

				props := model.Properties
				schema.AddressPrefixes = props.AddressPrefixes
				schema.Description = pointer.From(props.Description)
				schema.DisplayName = pointer.From(props.DisplayName)
				schema.ParentPoolName = pointer.From(props.ParentPoolName)
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r ManagerIpamPoolResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.IPamPools

			id, err := ipampools.ParseIPamPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerIpamPoolResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := ipampools.IPamPoolUpdate{
				Properties: &ipampools.IPamPoolUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = pointer.To(model.Tags)
			}

			if metadata.ResourceData.HasChange("description") {
				parameters.Properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("display_name") {
				parameters.Properties.DisplayName = pointer.To(model.DisplayName)
			}

			if _, err := client.Update(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r ManagerIpamPoolResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.IPamPools

			id, err := ipampools.ParseIPamPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			// https://github.com/Azure/azure-rest-api-specs/issues/31688
			pollerType := custompollers.NewNetworkManagerIPAMPoolDeletePoller(client, *id)
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return err
			}

			return nil
		},
	}
}
