// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_san_volume_group

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
	elasticsanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/elasticsan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elasticsan/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataSource struct{}

var _ sdk.DataSource = DataSource{}

type ElasticSANVolumeGroupDataSourceModel struct {
	SanId          string                                          `tfschema:"elastic_san_id"`
	EncryptionType string                                          `tfschema:"encryption_type"`
	Encryption     []ElasticSANVolumeGroupResourceEncryptionModel  `tfschema:"encryption"`
	Identity       []identity.ModelSystemAssignedUserAssigned      `tfschema:"identity"`
	Name           string                                          `tfschema:"name"`
	NetworkRule    []ElasticSANVolumeGroupResourceNetworkRuleModel `tfschema:"network_rule"`
	ProtocolType   string                                          `tfschema:"protocol_type"`
}

func (r DataSource) ResourceType() string {
	return "azurerm_elastic_san_volume_group"
}

func (r DataSource) Description() string {
	return "Use this data source to access information about an existing Elastic SAN Volume Group."
}

func (r DataSource) ModelObject() any {
	return &ElasticSANVolumeGroupDataSourceModel{}
}

func (r DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Elastic SAN Volume Group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ElasticSanVolumeGroupName,
		},

		"elastic_san_id": commonschema.ResourceIDReferenceRequired(&volumegroups.ElasticSanId{}, "Elastic SAN"),
	}
}

func (r DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"encryption_type": {
			Description: "The type of the key used to encrypt the data of the disk.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"encryption": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Description: "The Key Vault Key URI for Customer Managed Key encryption, which can be either a full URI or a versionless URI.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"user_assigned_identity_id": {
						Description: "The ID of the User Assigned Identity used by this Elastic SAN Volume Group.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"current_versioned_key_expiration_timestamp": {
						Description: "The timestamp of the expiration time for the current version of the Customer Managed Key.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"current_versioned_key_id": {
						Description: "The ID of the current versioned Key Vault Key in use.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"last_key_rotation_timestamp": {
						Description: "The timestamp of the last rotation of the Key Vault Key.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
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
						Description: "The ID of the Subnet from which access to this Elastic SAN Volume Group is allowed.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"action": {
						Description: "The action to take when an access attempt to this Elastic SAN Volume Group from this Subnet is made.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"protocol_type": {
			Description: "The type of the storage target.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"identity": commonschema.SystemOrUserAssignedIdentityComputed(),
	}
}

func (r DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := elasticsanClient.From(metadata.Client).VolumeGroups

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
					state.EncryptionType = pointer.FromEnum(model.Properties.Encryption)
					state.NetworkRule = FlattenVolumeGroupNetworkRules(model.Properties.NetworkAcls)

					if model.Properties.ProtocolType != nil {
						state.ProtocolType = pointer.FromEnum(model.Properties.ProtocolType)
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
