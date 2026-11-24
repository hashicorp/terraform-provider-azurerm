// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package netapp_volume

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing NetApp Volume.",
		Read:        dataSourceNetAppVolumeRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the NetApp Volume.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.VolumeName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"zone": commonschema.ZoneSingleComputed(),

			"account_name": {
				Description:  "The name of the NetApp account where the NetApp pool exists.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.AccountName,
			},

			"pool_name": {
				Description:  "The name of the NetApp pool where the NetApp volume exists.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.PoolName,
			},

			"mount_target": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"ip_address": {
							Description: "The IP address of the mount target.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"smb_server_fqdn": {
							Description: "The SMB server's Fully Qualified Domain Name (FQDN). This value is populated when the volume's `protocols` include `CIFS`; otherwise, it is empty.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"volume_path": {
				Description: "The unique file path of the volume.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"service_level": {
				Description: "The service level of the file system.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"subnet_id": {
				Description: "The ID of a Subnet in which the NetApp Volume resides.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"network_features": {
				Description: "Network features in use `Basic` or `Standard`.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"storage_quota_in_gb": {
				Description: "The maximum Storage Quota in Gigabytes allowed for a file system.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"protocols": {
				Description: "A list of protocol types enabled on volume.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},

			"security_style": {
				Description: "Volume security style.",
				Type:        pluginsdk.TypeString,
				Optional:    true,
			},

			"data_protection_replication": {
				Description: "Volume data protection replication block.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"endpoint_type": {
							Description: "The endpoint type.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"remote_volume_location": commonschema.LocationComputed(),

						"remote_volume_resource_id": {
							Description: "Resource ID of the primary volume.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"replication_frequency": {
							Description: "Frequency of replication.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"data_protection_backup_policy": {
				Description: "A data protection backup policy block.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"backup_policy_id": {
							Description: "The Resource ID of the backup policy.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"policy_enabled": {
							Description: "Backup policy is enabled or not.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},

						"backup_vault_id": {
							Description: "The Resource ID of the backup Backup Vault.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"data_protection_advanced_ransomware": {
				Description: "An Advanced Ransomware Protection (ARP) data protection block.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"protection_enabled": {
							Description: "Whether the Advanced Ransomware Protection feature is enabled.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},
					},
				},
			},

			"encryption_key_source": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"key_vault_private_endpoint_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"smb_non_browsable_enabled": {
				Description: "Limits clients from browsing for an SMB share.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"smb_access_based_enumeration_enabled": {
				Description: "Limits enumeration of files and folders (that is, listing the contents) in SMB only to users with allowed access on the share.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"large_volume_enabled": {
				Description: "Whether the volume is a large volume.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"breakthrough_mode_enabled": {
				Description: "Whether the large volume runs in Breakthrough Mode, placing it on dedicated capacity that provides higher throughput and greater capacity.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"accept_grow_capacity_pool_for_short_term_clone_split": {
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Description: "The accept grow capacity pool for short term clone split property.",
			},
		},
	}
}

func dataSourceNetAppVolumeRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetApp.VolumeClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := volumes.NewVolumeID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("pool_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.VolumeName)
	d.Set("pool_name", id.CapacityPoolName)
	d.Set("account_name", id.NetAppAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		zone := ""
		if model.Zones != nil {
			if zones := *model.Zones; len(zones) > 0 {
				zone = zones[0]
			}
		}
		d.Set("zone", zone)

		props := model.Properties
		d.Set("volume_path", props.CreationToken)
		d.Set("service_level", pointer.FromEnum(props.ServiceLevel))
		d.Set("subnet_id", props.SubnetId)
		d.Set("network_features", pointer.FromEnum(props.NetworkFeatures))
		d.Set("encryption_key_source", pointer.FromEnum(props.EncryptionKeySource))
		d.Set("key_vault_private_endpoint_id", props.KeyVaultPrivateEndpointResourceId)
		d.Set("large_volume_enabled", props.IsLargeVolume)
		d.Set("breakthrough_mode_enabled", pointer.From(props.BreakthroughMode) == volumes.BreakthroughModeEnabled)
		d.Set("accept_grow_capacity_pool_for_short_term_clone_split", pointer.FromEnum(props.AcceptGrowCapacityPoolForShortTermCloneSplit))

		smbNonBrowsable := false
		if props.SmbNonBrowsable != nil {
			smbNonBrowsable = strings.EqualFold(string(*props.SmbNonBrowsable), string(volumes.SmbNonBrowsableEnabled))
		}
		d.Set("smb_non_browsable_enabled", smbNonBrowsable)

		smbAccessBasedEnumeration := false
		if props.SmbAccessBasedEnumeration != nil {
			smbAccessBasedEnumeration = strings.EqualFold(string(*props.SmbAccessBasedEnumeration), string(volumes.SmbAccessBasedEnumerationEnabled))
		}
		d.Set("smb_access_based_enumeration_enabled", smbAccessBasedEnumeration)

		protocolTypes := make([]string, 0)
		if prtclTypes := props.ProtocolTypes; prtclTypes != nil {
			protocolTypes = append(protocolTypes, *prtclTypes...)
		}
		d.Set("protocols", protocolTypes)

		d.Set("security_style", pointer.FromEnum(props.SecurityStyle))

		d.Set("storage_quota_in_gb", props.UsageThreshold/1073741824)
		if err := d.Set("mount_target", flattenNetAppVolumeMountTargets(props.MountTargets)); err != nil {
			return fmt.Errorf("setting `mount_target`: %+v", err)
		}
		if err := d.Set("data_protection_replication", flattenNetAppVolumeDataProtectionReplication(props.DataProtection)); err != nil {
			return fmt.Errorf("setting `data_protection_replication`: %+v", err)
		}
		if err := d.Set("data_protection_backup_policy", flattenNetAppVolumeDataProtectionBackupPolicy(props.DataProtection)); err != nil {
			return fmt.Errorf("setting `data_protection_backup_policy`: %+v", err)
		}
		if err := d.Set("data_protection_advanced_ransomware", flattenNetAppVolumeDataProtectionAdvancedRansomwareProtection(props.DataProtection)); err != nil {
			return fmt.Errorf("setting `data_protection_advanced_ransomware`: %+v", err)
		}
	}

	return nil
}
