// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerinstance/2023-05-01/containerinstance"
	"github.com/hashicorp/go-azure-sdk/resource-manager/standbypool/2025-03-01/standbycontainergrouppools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = ContainerGroupStandbyPoolResource{}

type ContainerGroupStandbyPoolResource struct{}

type ContainerGroupStandbyPoolResourceModel struct {
	Name                    string            `tfschema:"name"`
	ResourceGroupName       string            `tfschema:"resource_group_name"`
	ContainerGroupId        string            `tfschema:"container_gorup_id"`
	ContainerGroupRevision  int64             `tfschema:"container_group_revision"`
	ContainerGroupSubnetIds []string          `tfschema:"container_group_subnet_ids"`
	MaxReadyCapacity        int64             `tfschema:"max_ready_capacity"`
	RefillPolicy            string            `tfschema:"refill_policy"`
	Zones                   []string          `tfschema:"zones"`
	Tags                    map[string]string `tfschema:"tags"`
}

func (ContainerGroupStandbyPoolResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"container_gorup_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: containerinstance.ValidateContainerGroupID,
		},

		"container_group_revision": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"subnet_ids": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: commonids.ValidateSubnetID,
			},
			Set: pluginsdk.HashString,
		},

		"max_ready_capacity": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"refill_policy": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(standbycontainergrouppools.PossibleValuesForRefillPolicy(), false),
		},

		"zone": commonschema.ZonesMultipleOptional(),

		"tags": tags.Schema(),
	}

}

func (ContainerGroupStandbyPoolResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (ContainerGroupStandbyPoolResource) ModelObject() interface{} {
	return &ContainerGroupStandbyPoolResourceModel{}
}

func (ContainerGroupStandbyPoolResource) ResourceType() string {
	return "azurerm_container_group_standby_pool"
}

func (r ContainerGroupStandbyPoolResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.StandbyContainerGroupPoolsClient

			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ContainerGroupStandbyPoolResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := standbycontainergrouppools.NewStandbyContainerGroupPoolID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := standbycontainergrouppools.StandbyContainerGroupPoolResource{
				Properties: &standbycontainergrouppools.StandbyContainerGroupPoolResourceProperties{
					ContainerGroupProperties: standbycontainergrouppools.ContainerGroupProperties{
						ContainerGroupProfile: standbycontainergrouppools.ContainerGroupProfile{
							Id:       config.ContainerGroupId,
							Revision: pointer.To(config.ContainerGroupRevision),
						},
						SubnetIds: pointer.To(expandContainerGroupStandbyPoolSubnetIds(config.ContainerGroupSubnetIds)),
					},
					ElasticityProfile: standbycontainergrouppools.StandbyContainerGroupPoolElasticityProfile{
						MaxReadyCapacity: config.MaxReadyCapacity,
						RefillPolicy:     pointer.To(standbycontainergrouppools.RefillPolicy(config.RefillPolicy)),
					},
					Zones: pointer.To(zones.Expand(config.Zones)),
				},
				Tags: pointer.To(config.Tags),
			}

			client.CreateOrUpdateThenPoll(ctx, id, payload)

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ContainerGroupStandbyPoolResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.StandbyContainerGroupPoolsClient

			id, err := standbycontainergrouppools.ParseStandbyContainerGroupPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config ContainerGroupStandbyPoolResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload := standbycontainergrouppools.StandbyContainerGroupPoolResourceUpdate{
				Properties: &standbycontainergrouppools.StandbyContainerGroupPoolResourceUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("container_gorup_id") {
				if payload.Properties.ContainerGroupProperties == nil {
					payload.Properties.ContainerGroupProperties = &standbycontainergrouppools.ContainerGroupProperties{}
				}
				payload.Properties.ContainerGroupProperties.ContainerGroupProfile.Id = config.ContainerGroupId
			}

			if metadata.ResourceData.HasChange("container_group_revision") {
				if payload.Properties.ContainerGroupProperties == nil {
					payload.Properties.ContainerGroupProperties = &standbycontainergrouppools.ContainerGroupProperties{}
				}
				payload.Properties.ContainerGroupProperties.ContainerGroupProfile.Revision = &config.ContainerGroupRevision
			}

			if metadata.ResourceData.HasChange("subnet_ids") {
				if payload.Properties.ContainerGroupProperties == nil {
					payload.Properties.ContainerGroupProperties = &standbycontainergrouppools.ContainerGroupProperties{}
				}
				payload.Properties.ContainerGroupProperties.SubnetIds = pointer.To(expandContainerGroupStandbyPoolSubnetIds(config.ContainerGroupSubnetIds))
			}

			if metadata.ResourceData.HasChange("refill_policy") {
				if payload.Properties.ElasticityProfile == nil {
					payload.Properties.ElasticityProfile = &standbycontainergrouppools.StandbyContainerGroupPoolElasticityProfile{}
				}
				payload.Properties.ElasticityProfile.RefillPolicy = pointer.To(standbycontainergrouppools.RefillPolicy(config.RefillPolicy))
			}

			if metadata.ResourceData.HasChange("max_ready_capacity") {
				if payload.Properties.ElasticityProfile == nil {
					payload.Properties.ElasticityProfile = &standbycontainergrouppools.StandbyContainerGroupPoolElasticityProfile{}
				}
				payload.Properties.ElasticityProfile.MaxReadyCapacity = config.MaxReadyCapacity
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(config.Tags)
			}

			if _, err := client.Update(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}

}

func (r ContainerGroupStandbyPoolResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.StandbyContainerGroupPoolsClient

			id, err := standbycontainergrouppools.ParseStandbyContainerGroupPoolID(metadata.ResourceData.Id())
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

			state := ContainerGroupStandbyPoolResourceModel{
				Name:              id.StandbyContainerGroupPoolName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				if prop := model.Properties; prop != nil {
					state.ContainerGroupId = prop.ContainerGroupProperties.ContainerGroupProfile.Id
					state.ContainerGroupRevision = pointer.From(prop.ContainerGroupProperties.ContainerGroupProfile.Revision)
					state.MaxReadyCapacity = prop.ElasticityProfile.MaxReadyCapacity
					state.RefillPolicy = string(pointer.From(prop.ElasticityProfile.RefillPolicy))
					state.Zones = zones.Flatten(prop.Zones)
					state.ContainerGroupSubnetIds = flattenContainerGroupStandbyPoolSubnetIds(prop.ContainerGroupProperties.SubnetIds)
				}
			}

			return nil
		},
	}

}

func (r ContainerGroupStandbyPoolResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.StandbyContainerGroupPoolsClient

			id, err := standbycontainergrouppools.ParseStandbyContainerGroupPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}

}

func (r ContainerGroupStandbyPoolResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return standbycontainergrouppools.ValidateStandbyContainerGroupPoolID
}

func expandContainerGroupStandbyPoolSubnetIds(input []string) []standbycontainergrouppools.Subnet {
	result := make([]standbycontainergrouppools.Subnet, 0)

	if len(input) != 0 {
		for _, id := range input {
			result = append(result, standbycontainergrouppools.Subnet{
				Id: id,
			})
		}
	}

	return result
}
func flattenContainerGroupStandbyPoolSubnetIds(input *[]standbycontainergrouppools.Subnet) []string {
	result := make([]string, 0)

	if input != nil {
		for _, snet := range *input {
			result = append(result, snet.Id)
		}
	}

	return result
}
