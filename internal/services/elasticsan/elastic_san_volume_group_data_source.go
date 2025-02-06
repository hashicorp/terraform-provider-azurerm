// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package elasticsan

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elasticsan/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ElasticSANVolumeGroupDataSource struct{}

var _ sdk.DataSource = ElasticSANVolumeGroupDataSource{}

type ElasticSANVolumeGroupDataSourceModel struct {
	SanId          string                                          `tfschema:"elastic_san_id"`
	EncryptionType string                                          `tfschema:"encryption_type"`
	Encryption     []ElasticSANVolumeGroupResourceEncryptionModel  `tfschema:"encryption"`
	Identity       []identity.ModelSystemAssignedUserAssigned      `tfschema:"identity"`
	Name           string                                          `tfschema:"name"`
	NetworkRule    []ElasticSANVolumeGroupResourceNetworkRuleModel `tfschema:"network_rule"`
	ProtocolType   string                                          `tfschema:"protocol_type"`
}

func (r ElasticSANVolumeGroupDataSource) ResourceType() string {
	return "azurerm_elastic_san_volume_group"
}

func (r ElasticSANVolumeGroupDataSource) ModelObject() interface{} {
	return &ElasticSANVolumeGroupDataSourceModel{}
}

func (r ElasticSANVolumeGroupDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ElasticSanVolumeGroupName,
		},

		"elastic_san_id": commonschema.ResourceIDReferenceRequired(&volumegroups.ElasticSanId{}),
	}
}

func (r ElasticSANVolumeGroupDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"encryption_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"encryption": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},
					"user_assigned_identity_id": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},
					"current_versioned_key_expiration_timestamp": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},
					"current_versioned_key_id": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},
					"last_key_rotation_timestamp": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},
				},
			},
		},

		"network_rule": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},
					"action": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},
				},
			},
		},

		"protocol_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"identity": commonschema.SystemOrUserAssignedIdentityComputed(),
	}
}

func (r ElasticSANVolumeGroupDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ElasticSan.VolumeGroups

			var state ElasticSANVolumeGroupDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			elasticSanId, err := volumegroups.ParseElasticSanID(state.SanId)
			if err != nil {
				return err
			}

			id := volumegroups.NewVolumeGroupID(elasticSanId.SubscriptionId, elasticSanId.ResourceGroupName, elasticSanId.ElasticSanName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s does not exist", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state.SanId = elasticSanId.ID()
			state.Name = id.VolumeGroupName

			if model := resp.Model; model != nil {
				flattenedIdentity, err := identity.FlattenSystemOrUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state.Identity = *flattenedIdentity

				if model.Properties != nil {
					state.EncryptionType = string(pointer.From(model.Properties.Encryption))
					state.NetworkRule = FlattenVolumeGroupNetworkRules(model.Properties.NetworkAcls)

					if model.Properties.ProtocolType != nil {
						state.ProtocolType = string(pointer.From(model.Properties.ProtocolType))
					}

					state.Encryption, err = FlattenVolumeGroupEncryption(model.Properties.EncryptionProperties)
					if err != nil {
						return fmt.Errorf("flattening `encryption`: %+v", err)
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
