package netapp

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2021-10-01/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2021-10-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type NetappVolumeGroupResource struct{}

type NetappVolumeGroupModel struct {
}

var _ sdk.ResourceWithUpdate = NetappVolumeGroupResource{}

func (r NetappVolumeGroupResource) ModelObject() interface{} {
	return &NetappVolumeGroupModel{}
}

func (r NetappVolumeGroupResource) ResourceType() string {
	return "azurerm_netapp_volume_group"
}

func (r NetappVolumeGroupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	panic("Implement me") // TODO - Add Validation func return here
}

func (r NetappVolumeGroupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.VolumeGroupName,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"location": azure.SchemaLocation(),

		"account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.AccountName,
		},

		"group_description": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"application_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(netapp.ApplicationTypeSAPHANA),
			}, false),
		},

		"application_identifier": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 3),
		},

		"deployment_spec_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"volume": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MinItems: 5,
			MaxItems: 5,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: netAppValidate.VolumeName,
					},

					"volume_path": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: netAppValidate.VolumePath,
					},

					"service_level": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(volumes.ServiceLevelPremium),
							string(volumes.ServiceLevelStandard),
							string(volumes.ServiceLevelUltra),
						}, false),
					},

					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: azure.ValidateResourceID,
					},

					"create_from_snapshot_resource_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: snapshots.ValidateSnapshotID,
					},

					"network_features": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(volumes.NetworkFeaturesBasic),
							string(volumes.NetworkFeaturesStandard),
						}, false),
					},

					"protocols": {
						Type:     pluginsdk.TypeSet,
						ForceNew: true,
						Optional: true,
						Computed: true,
						MaxItems: 2,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"NFSv3",
								"NFSv4.1",
								"CIFS",
							}, false),
						},
					},

					"security_style": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Computed: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Unix", // Using hardcoded values instead of SDK enum since no matter what case is passed,
							"Ntfs", // ANF changes casing to Pascal case in the backend. Please refer to https://github.com/Azure/azure-sdk-for-go/issues/14684
						}, false),
					},

					"storage_quota_in_gb": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(100, 102400),
					},

					"throughput_in_mibps": {
						Type:     pluginsdk.TypeFloat,
						Optional: true,
						Computed: true,
					},

					"capacity_pool_id": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						Computed:         true,
						ForceNew:         true,
						DiffSuppressFunc: suppress.CaseDifference,
						ValidateFunc:     azure.ValidateResourceID,
					},

					"proximity_placement_group_id": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						ForceNew:         true,
						DiffSuppressFunc: suppress.CaseDifference,
						ValidateFunc:     azure.ValidateResourceID,
					},

					"volume_spec_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},

					"export_policy_rule": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 5,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"rule_index": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 5),
								},

								"allowed_clients": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validate.CIDR,
									},
								},

								"protocols_enabled": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Computed: true,
									MaxItems: 1,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											"NFSv3",
											"NFSv4.1",
											"CIFS",
										}, false),
									},
								},

								"unix_read_only": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},

								"unix_read_write": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},

								"root_access_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},
							},
						},
					},

					"tags": commonschema.Tags(),

					"mount_ip_addresses": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"snapshot_directory_visible": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
					},

					"data_protection_replication": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						ForceNew: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"endpoint_type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "dst",
									ValidateFunc: validation.StringInSlice([]string{
										"dst",
									}, false),
								},

								"remote_volume_location": azure.SchemaLocation(),

								"remote_volume_resource_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: azure.ValidateResourceID,
								},

								"replication_frequency": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										"10minutes",
										"daily",
										"hourly",
									}, false),
								},
							},
						},
					},

					"data_protection_snapshot_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
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
				},
			},
		},

		// Can't use tags.Schema since there is no patch available
		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r NetappVolumeGroupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		/*
			TODO - This section is for `Computed: true` only items, i.e. useful values that are returned by the
			datasource that can be used as outputs or passed programmatically to other resources or data sources.
		*/
	}
}

func (r NetappVolumeGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Create Func
			// TODO - Don't forget to set the ID! e.g. metadata.SetID(id)
			return nil
		},
	}
}

func (r NetappVolumeGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Read Func
			return nil
		},
	}
}

func (r NetappVolumeGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Delete Func
			return nil
		},
	}
}

func (r NetappVolumeGroupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Update Func
			return nil
		},
	}
}
