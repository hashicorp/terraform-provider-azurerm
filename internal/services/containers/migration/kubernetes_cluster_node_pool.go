// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = KubernetesClusterNodePoolV0ToV1{}

type KubernetesClusterNodePoolV0ToV1 struct{}

func (k KubernetesClusterNodePoolV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		log.Printf("[DEBUG] Migrating ID to correct casing for Kubernetes Cluster")

		originClusterId := rawState["kubernetes_cluster_id"].(string)
		clusterId, err := parse.ClusterID(originClusterId)
		if err != nil {
			return nil, err
		}
		rawState["kubernetes_cluster_id"] = clusterId.ID()

		id := rawState["id"].(string)
		poolId, err := parse.NodePoolID(id)
		if err != nil {
			return nil, err
		}
		rawState["id"] = poolId.ID()

		return rawState, nil
	}
}

func (k KubernetesClusterNodePoolV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"kubernetes_cluster_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"node_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"vm_size": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"availability_zones": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"enable_auto_scaling": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"enable_host_encryption": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"enable_node_public_ip": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"eviction_policy": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"kubelet_config": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cpu_manager_policy": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"cpu_cfs_quota_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"cpu_cfs_quota_period": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"image_gc_high_threshold": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"image_gc_low_threshold": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"topology_manager_policy": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"allowed_unsafe_sysctls": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"container_log_max_size_mb": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					// TODO 4.0: change this to `container_log_max_files`
					"container_log_max_line": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"pod_max_pid": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"linux_os_config": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"sysctl_config": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"fs_aio_max_nr": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"fs_file_max": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"fs_inotify_max_user_watches": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"fs_nr_open": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"kernel_threads_max": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"net_core_netdev_max_backlog": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"net_core_optmem_max": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"net_core_rmem_default": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"net_core_rmem_max": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"net_core_somaxconn": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"net_core_wmem_default": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"net_core_wmem_max": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"net_ipv4_ip_local_port_range_min": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"net_ipv4_ip_local_port_range_max": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"net_ipv4_neigh_default_gc_thresh1": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"net_ipv4_neigh_default_gc_thresh2": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"net_ipv4_neigh_default_gc_thresh3": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"net_ipv4_tcp_fin_timeout": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"net_ipv4_tcp_keepalive_intvl": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"net_ipv4_tcp_keepalive_probes": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"net_ipv4_tcp_keepalive_time": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"net_ipv4_tcp_max_syn_backlog": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"net_ipv4_tcp_max_tw_buckets": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"net_ipv4_tcp_tw_reuse": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},

								"net_netfilter_nf_conntrack_buckets": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"net_netfilter_nf_conntrack_max": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"vm_max_map_count": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"vm_swappiness": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"vm_vfs_cache_pressure": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},
							},
						},
					},

					"transparent_huge_page_enabled": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"transparent_huge_page_defrag": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"swap_file_size_mb": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"fips_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"kubelet_disk_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"max_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"max_pods": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"min_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"node_labels": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"node_public_ip_prefix_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"node_taints": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"orchestrator_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"os_disk_size_gb": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"os_disk_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"os_sku": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // defaults to Ubuntu if using Linux
		},

		"os_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"pod_subnet_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"priority": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"proximity_placement_group_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"spot_max_price": {
			Type:     pluginsdk.TypeFloat,
			Optional: true,
		},

		"ultra_ssd_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"vnet_subnet_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"upgrade_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"max_surge": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},
	}
}
