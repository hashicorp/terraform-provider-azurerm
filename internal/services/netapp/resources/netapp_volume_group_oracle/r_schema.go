// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume_group_oracle

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VolumeGroupName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountName,
		},

		"group_description": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"application_identifier": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z][\w-]{2,11}$`),
				"`application_identifier` must be between 3 and 12 characters, may contain alphanumerics, hyphens, and underscores, and must begin with a letter.",
			),
		},

		"volume": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 2,
			MaxItems: 12,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validate.VolumeName,
					},

					"capacity_pool_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: azure.ValidateResourceID,
					},

					"proximity_placement_group_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"zone": commonschema.ZoneSingleOptionalForceNew(),

					"volume_spec_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(validate.PossibleValuesForVolumeSpecNameOracle(), false),
					},

					"volume_path": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validate.VolumePath,
					},

					"service_level": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(volumegroups.ServiceLevelPremium),
							string(volumegroups.ServiceLevelStandard),
							string(volumegroups.ServiceLevelUltra),
							string(volumegroups.ServiceLevelFlexible),
						}, false),
					},

					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: azure.ValidateResourceID,
					},

					"protocols": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						MaxItems: 1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(validate.PossibleValuesForProtocolTypeVolumeGroupOracle(), false),
						},
					},

					"security_style": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(volumegroups.PossibleValuesForSecurityStyle(), false),
					},

					"storage_quota_in_gb": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(50, 102400),
					},

					"throughput_in_mibps": {
						Type:         pluginsdk.TypeFloat,
						Required:     true,
						ValidateFunc: validation.FloatAtLeast(0.1),
					},

					"export_policy_rule": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						MaxItems: 5,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"rule_index": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 5),
								},

								"allowed_clients": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"nfsv3_enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},

								"nfsv41_enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},

								"unix_read_only": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"unix_read_write": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  true,
								},

								"root_access_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  true,
								},
							},
						},
					},

					"tags": commonschema.Tags(),

					"snapshot_directory_visible": {
						Type:     pluginsdk.TypeBool,
						Required: true,
						ForceNew: true,
					},

					"network_features": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// Note: O+C because the service team is changing network features on the backend to upgrade everyone from Basic to Standard and there is a feature that allows customers to change network features from portal but not the API. This could cause drift that forces data loss that we want to avoid
						Computed:     true,
						ValidateFunc: validation.StringInSlice(volumegroups.PossibleValuesForNetworkFeatures(), false),
					},

					"mount_ip_addresses": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"data_protection_replication": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						ForceNew: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"endpoint_type": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      string(volumegroups.EndpointTypeDst),
									ForceNew:     true,
									ValidateFunc: validation.StringInSlice(volumegroups.PossibleValuesForEndpointType(), false),
								},

								"remote_volume_location": commonschema.LocationWithoutForceNew(),

								"remote_volume_resource_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: azure.ValidateResourceID,
								},

								"replication_frequency": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringInSlice(models.PossibleValuesForReplicationSchedule(), false),
								},
							},
						},
					},

					"data_protection_snapshot_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						// Note: O+C - Adding this because Terraform is not being able to build proper deletion graph, it is trying to delete the snapshot policy before the volume because this is in a deeper level within the schema inside an array of volumes
						Computed: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"snapshot_policy_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: azure.ValidateResourceID,
								},
							},
						},
					},

					"encryption_key_source": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						// Note: O+C since there is a feature coming up that will allow customers to change the encryption key source from portal but not the API. This could cause drift if configuration is not updated
						Computed:     true,
						ValidateFunc: validation.StringInSlice(volumes.PossibleValuesForEncryptionKeySource(), false),
					},

					"key_vault_private_endpoint_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						// Note: O+C since there is a feature coming up that will allow customers to change the encryption key source from portal but not the API (tied to encryption_key_source). This could cause drift if configuration is not updated
						Computed:     true,
						ValidateFunc: azure.ValidateResourceID,
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			// Validate NFSv3 to NFSv4.1 protocol conversion restrictions for volume groups
			volumes := rd.Get("volume").([]any)
			for i := range volumes {
				protocolsKey := fmt.Sprintf("volume.%d.protocols", i)

				if rd.HasChange(protocolsKey) {
					old, new := rd.GetChange(protocolsKey)
					oldProtocols := old.([]any)
					newProtocols := new.([]any)

					// Convert to string slices for validation
					oldProtocolsStr := make([]string, len(oldProtocols))
					newProtocolsStr := make([]string, len(newProtocols))

					for j, v := range oldProtocols {
						oldProtocolsStr[j] = v.(string)
					}
					for j, v := range newProtocols {
						newProtocolsStr[j] = v.(string)
					}

					// Get the export policy rules configuration for this volume
					exportPolicyRulesKey := fmt.Sprintf("volume.%d.export_policy_rule", i)
					exportPolicyRules := rd.Get(exportPolicyRulesKey).([]any)

					// For volume groups, kerberos and data replication are not directly supported, so we pass empty values
					var kerberosEnabled bool
					var dataReplication []any

					validationErrors := validate.ValidateNetAppVolumeProtocolConversion(oldProtocolsStr, newProtocolsStr, kerberosEnabled, dataReplication, exportPolicyRules)
					for _, err := range validationErrors {
						return err
					}
				}
			}

			return nil
		},
	}
}
