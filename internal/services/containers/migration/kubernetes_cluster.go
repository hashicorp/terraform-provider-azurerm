// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = KubernetesClusterV0ToV1{}

var _ pluginsdk.StateUpgrade = KubernetesClusterV1ToV2{}

type KubernetesClusterV0ToV1 struct{}

type KubernetesClusterV1ToV2 struct{}

func (k KubernetesClusterV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		log.Printf("[DEBUG] Migrating ID to correct casing for Kubernetes Cluster")
		rawId := rawState["id"].(string)

		id, err := parse.ClusterID(rawId)
		if err != nil {
			return nil, err
		}

		rawState["id"] = id.ID()
		return rawState, nil
	}
}

func (k KubernetesClusterV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": {
			Type:     schema.TypeString,
			Required: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"dns_prefix": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"dns_prefix_private_cluster": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"kubernetes_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"default_node_pool": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
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

					"enable_node_public_ip": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"enable_host_encryption": {
						Type:     pluginsdk.TypeBool,
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

					"min_count": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"node_count": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Computed: true,
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

					"tags": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
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
						Computed: true,
					},

					"ultra_ssd_enabled": {
						Type:     pluginsdk.TypeBool,
						Default:  false,
						Optional: true,
					},

					"vnet_subnet_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"orchestrator_version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"pod_subnet_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"proximity_placement_group_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"only_critical_addons_enabled": {
						Type:     pluginsdk.TypeBool,
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
				},
			},
		},

		"addon_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"aci_connector_linux": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},

								"subnet_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},

					"azure_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
							},
						},
					},

					"kube_dashboard": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
							},
						},
					},

					"http_application_routing": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
								"http_application_routing_zone_name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"oms_agent": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
								"log_analytics_workspace_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"oms_agent_identity": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"client_id": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
											"object_id": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
											"user_assigned_identity_id": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},

					"ingress_application_gateway": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
								"gateway_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"gateway_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"subnet_cidr": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"subnet_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"effective_gateway_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"ingress_application_gateway_identity": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"client_id": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
											"object_id": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
											"user_assigned_identity_id": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},

					"open_service_mesh": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
							},
						},
					},
				},
			},
		},

		"api_server_authorized_ip_ranges": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"auto_scaler_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"balance_similar_node_groups": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"expander": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"max_graceful_termination_sec": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"max_node_provisioning_time": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"max_unready_nodes": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
					"max_unready_percentage": {
						Type:     pluginsdk.TypeFloat,
						Optional: true,
					},
					"new_pod_scale_up_delay": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"scan_interval": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"scale_down_delay_after_add": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"scale_down_delay_after_delete": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"scale_down_delay_after_failure": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"scale_down_unneeded": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"scale_down_unready": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"scale_down_utilization_threshold": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"empty_bulk_delete_max": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"skip_nodes_with_local_storage": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"skip_nodes_with_system_pods": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},

		"disk_encryption_set_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"enable_pod_security_policy": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"identity": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"user_assigned_identity_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"principal_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"kubelet_identity": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"object_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"user_assigned_identity_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},

		"linux_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"admin_username": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"ssh_key": {
						Type:     pluginsdk.TypeList,
						Required: true,

						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"key_data": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},

		"local_account_disabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"maintenance_window": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allowed": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"day": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"hours": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeInt,
									},
								},
							},
						},
					},

					"not_allowed": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"end": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"start": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},

		"network_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"network_plugin": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"network_mode": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"network_policy": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"dns_service_ip": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"docker_bridge_cidr": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"pod_cidr": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"service_cidr": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"load_balancer_sku": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"outbound_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"load_balancer_profile": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"outbound_ports_allocated": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},
								"idle_timeout_in_minutes": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},
								"managed_outbound_ip_count": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									Computed: true,
								},
								"outbound_ip_prefix_ids": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"outbound_ip_address_ids": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"effective_outbound_ips": {
									Type:     pluginsdk.TypeSet,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},

		"node_resource_group": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"private_fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"portal_fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"private_link_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"private_cluster_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"private_cluster_public_fqdn_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"private_dns_zone_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"role_based_access_control": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"azure_active_directory": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_app_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"server_app_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"server_app_secret": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"tenant_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Computed: true,
								},

								"managed": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},

								"azure_rbac_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},

								"admin_group_object_ids": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},

		"service_principal": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"client_secret": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"sku_tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"windows_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"admin_username": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"admin_password": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"license": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"automatic_channel_upgrade": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kube_admin_config": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"host": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"username": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"password": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"client_certificate": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"client_key": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"cluster_ca_certificate": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"kube_admin_config_raw": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kube_config": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"host": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"username": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"password": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"client_certificate": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"client_key": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"cluster_ca_certificate": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"kube_config_raw": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (k KubernetesClusterV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	// since `server_app_secret` isn't returned from the AKS API we need to populate that value in the new `azure_active_directory_role_based_access_control` block
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		secretRaw := ""
		if rbac, ok := rawState["role_based_access_control"]; ok {
			rbacRaw := rbac.([]interface{})[0].(map[string]interface{})
			if aad, ok := rbacRaw["azure_active_directory"]; ok {
				aadRaw := aad.([]interface{})
				if len(aadRaw) == 0 {
					return rawState, nil
				}
				aadMap := aadRaw[0].(map[string]interface{})
				if secret, ok := aadMap["server_app_secret"]; ok {
					log.Printf("[DEBUG] found value for `role_based_access_control.0.azure_active_directory.0.server_app_secret`")
					secretRaw = secret.(string)
				}
			}
		}

		if secretRaw != "" {
			log.Printf("[DEBUG] copying value to `azure_active_directory_role_based_access_control.0.server_app_secret`")
			rawState["azure_active_directory_role_based_access_control"] = []interface{}{
				map[string]interface{}{
					"server_app_secret": secretRaw,
				},
			}
		}

		return rawState, nil
	}
}

func (k KubernetesClusterV1ToV2) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": {
			Type:     schema.TypeString,
			Required: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"dns_prefix": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"dns_prefix_private_cluster": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"kubernetes_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"default_node_pool": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
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

					"enable_node_public_ip": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"enable_host_encryption": {
						Type:     pluginsdk.TypeBool,
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

					"min_count": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"node_count": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Computed: true,
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

					"tags": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
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
						Computed: true,
					},

					"ultra_ssd_enabled": {
						Type:     pluginsdk.TypeBool,
						Default:  false,
						Optional: true,
					},

					"vnet_subnet_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"orchestrator_version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"pod_subnet_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"proximity_placement_group_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"only_critical_addons_enabled": {
						Type:     pluginsdk.TypeBool,
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
				},
			},
		},

		"addon_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"aci_connector_linux": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},

								"subnet_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},

					"azure_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
							},
						},
					},

					"kube_dashboard": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
							},
						},
					},

					"http_application_routing": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
								"http_application_routing_zone_name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"oms_agent": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
								"log_analytics_workspace_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"oms_agent_identity": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"client_id": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
											"object_id": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
											"user_assigned_identity_id": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},

					"ingress_application_gateway": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
								"gateway_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"gateway_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"subnet_cidr": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"subnet_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"effective_gateway_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"ingress_application_gateway_identity": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"client_id": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
											"object_id": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
											"user_assigned_identity_id": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},

					"open_service_mesh": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
							},
						},
					},
				},
			},
		},

		"api_server_authorized_ip_ranges": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"auto_scaler_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"balance_similar_node_groups": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"expander": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"max_graceful_termination_sec": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"max_node_provisioning_time": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"max_unready_nodes": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
					"max_unready_percentage": {
						Type:     pluginsdk.TypeFloat,
						Optional: true,
					},
					"new_pod_scale_up_delay": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"scan_interval": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"scale_down_delay_after_add": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"scale_down_delay_after_delete": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"scale_down_delay_after_failure": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"scale_down_unneeded": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"scale_down_unready": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"scale_down_utilization_threshold": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"empty_bulk_delete_max": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"skip_nodes_with_local_storage": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"skip_nodes_with_system_pods": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},

		"disk_encryption_set_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"enable_pod_security_policy": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"identity": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"user_assigned_identity_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"principal_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"kubelet_identity": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"object_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"user_assigned_identity_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},

		"linux_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"admin_username": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"ssh_key": {
						Type:     pluginsdk.TypeList,
						Required: true,

						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"key_data": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},

		"local_account_disabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"maintenance_window": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allowed": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"day": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"hours": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeInt,
									},
								},
							},
						},
					},

					"not_allowed": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"end": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"start": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},

		"network_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"network_plugin": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"network_mode": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"network_policy": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"dns_service_ip": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"docker_bridge_cidr": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"pod_cidr": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"service_cidr": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"load_balancer_sku": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"outbound_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"load_balancer_profile": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"outbound_ports_allocated": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},
								"idle_timeout_in_minutes": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},
								"managed_outbound_ip_count": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									Computed: true,
								},
								"outbound_ip_prefix_ids": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"outbound_ip_address_ids": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"effective_outbound_ips": {
									Type:     pluginsdk.TypeSet,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},

		"node_resource_group": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"private_fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"portal_fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"private_link_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"private_cluster_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"private_cluster_public_fqdn_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"private_dns_zone_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"role_based_access_control": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"azure_active_directory": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_app_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"server_app_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"server_app_secret": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"tenant_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Computed: true,
								},

								"managed": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},

								"azure_rbac_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},

								"admin_group_object_ids": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},

		"role_based_access_control_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"azure_active_directory_role_based_access_control": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_app_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"server_app_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"server_app_secret": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"managed": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"azure_rbac_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"admin_group_object_ids": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"service_principal": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"client_secret": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"sku_tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"windows_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"admin_username": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"admin_password": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"license": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"automatic_channel_upgrade": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kube_admin_config": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"host": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"username": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"password": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"client_certificate": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"client_key": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"cluster_ca_certificate": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"kube_admin_config_raw": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kube_config": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"host": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"username": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"password": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"client_certificate": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"client_key": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"cluster_ca_certificate": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"kube_config_raw": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"http_proxy_config": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"http_proxy": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"https_proxy": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"no_proxy": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"trusted_ca": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},
	}
}
