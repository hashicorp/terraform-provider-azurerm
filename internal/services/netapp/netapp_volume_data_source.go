// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceNetAppVolume() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceNetAppVolumeRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.VolumeName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"zone": commonschema.ZoneSingleComputed(),

			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.AccountName,
			},

			"pool_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.PoolName,
			},

			"mount_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"volume_path": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"service_level": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"subnet_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"network_features": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"storage_quota_in_gb": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"protocols": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},

			"security_style": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"data_protection_replication": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"endpoint_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"remote_volume_location": commonschema.LocationComputed(),

						"remote_volume_resource_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"replication_frequency": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"data_protection_backup_policy": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"backup_policy_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"policy_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"backup_vault_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
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
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"smb_access_based_enumeration_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceNetAppVolumeRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		d.Set("service_level", string(pointer.From(props.ServiceLevel)))
		d.Set("subnet_id", props.SubnetId)
		d.Set("network_features", string(pointer.From(props.NetworkFeatures)))
		d.Set("encryption_key_source", string(pointer.From(props.EncryptionKeySource)))
		d.Set("key_vault_private_endpoint_id", props.KeyVaultPrivateEndpointResourceId)

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

		d.Set("security_style", string(pointer.From(props.SecurityStyle)))

		d.Set("storage_quota_in_gb", props.UsageThreshold/1073741824)
		if err := d.Set("mount_ip_addresses", flattenNetAppVolumeMountIPAddresses(props.MountTargets)); err != nil {
			return fmt.Errorf("setting `mount_ip_addresses`: %+v", err)
		}
		if err := d.Set("data_protection_replication", flattenNetAppVolumeDataProtectionReplication(props.DataProtection)); err != nil {
			return fmt.Errorf("setting `data_protection_replication`: %+v", err)
		}
		if err := d.Set("data_protection_backup_policy", flattenNetAppVolumeDataProtectionBackupPolicy(props.DataProtection)); err != nil {
			return fmt.Errorf("setting `data_protection_backup_policy`: %+v", err)
		}
	}

	return nil
}
