// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = SiteRecoveryReplicatedVMV0ToV1{}

type SiteRecoveryReplicatedVMV0ToV1 struct{}

func (SiteRecoveryReplicatedVMV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		return rawState, nil
	}
}

func (SiteRecoveryReplicatedVMV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"recovery_vault_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"source_recovery_fabric_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"source_vm_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"target_recovery_fabric_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"recovery_replication_policy_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"source_recovery_protection_container_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"target_recovery_protection_container_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"target_resource_group_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"target_availability_set_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"target_zone": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"target_network_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"test_network_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"target_edge_zone": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"unmanaged_disk": {
			Type:       pluginsdk.TypeSet,
			Optional:   true,
			Computed:   true,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"disk_uri": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"staging_storage_account_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"target_storage_account_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"multi_vm_group_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"managed_disk": {
			Type:       pluginsdk.TypeSet,
			Optional:   true,
			Computed:   true,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"disk_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"staging_storage_account_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"target_resource_group_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"target_disk_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"target_replica_disk_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"target_disk_encryption_set_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"target_disk_encryption": {
						Type:       pluginsdk.TypeList,
						Optional:   true,
						ConfigMode: pluginsdk.SchemaConfigModeAttr,
						MaxItems:   1,
						Elem:       diskEncryptionResourceV0(),
					},
				},
			},
		},

		"target_proximity_placement_group_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"target_boot_diagnostic_storage_account_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"target_capacity_reservation_group_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"target_virtual_machine_scale_set_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"target_virtual_machine_size": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"network_interface": {
			Type:       pluginsdk.TypeSet,
			Optional:   true,
			Computed:   true,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem:       networkInterfaceResourceV0(),
		},
	}
}

func networkInterfaceResourceV0() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"source_network_interface_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
			},
			"failover_test_static_ip": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
			},
			"target_static_ip": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
			"failover_test_subnet_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
			},
			"target_subnet_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
			"failover_test_public_ip_address_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
			},
			"recovery_load_balancer_backend_address_pool_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
			"recovery_public_ip_address_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
		},
	}
}

func diskEncryptionResourceV0() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"disk_encryption_key": {
				Type:       pluginsdk.TypeList,
				Required:   true,
				MaxItems:   1,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"secret_url": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"vault_id": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					},
				},
			},
			"key_encryption_key": {
				Type:       pluginsdk.TypeList,
				Optional:   true,
				MaxItems:   1,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_url": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"vault_id": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}
