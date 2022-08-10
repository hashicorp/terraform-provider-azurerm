package containers

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/containerservice/mgmt/2022-03-02-preview/containerservice"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/kubernetes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	laparse "github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceKubernetesCluster() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Read: dataSourceKubernetesClusterRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"aci_connector_linux": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"subnet_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"agent_pool_profile": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"count": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"max_count": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"min_count": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						// TODO 4.0: change this from enable_* to *_enabled
						"enable_auto_scaling": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"vm_size": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"tags": commonschema.TagsDataSource(),

						"os_disk_size_gb": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"vnet_subnet_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"os_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"orchestrator_version": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"max_pods": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"node_labels": {
							Type:     pluginsdk.TypeMap,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"node_taints": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},

						// TODO 4.0: change this from enable_* to *_enabled
						"enable_node_public_ip": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"node_public_ip_prefix_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"upgrade_settings": upgradeSettingsForDataSourceSchema(),

						"zones": commonschema.ZonesMultipleComputed(),
					},
				},
			},

			"azure_active_directory_role_based_access_control": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"client_app_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"server_app_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"managed": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"azure_rbac_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"admin_group_object_ids": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"azure_policy_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"dns_prefix": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"http_application_routing_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"http_application_routing_zone_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"ingress_application_gateway": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"gateway_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"gateway_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"subnet_cidr": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
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

			"key_vault_secrets_provider": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"secret_rotation_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
						"secret_rotation_interval": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"secret_identity": {
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

			"api_server_authorized_ip_ranges": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"disk_encryption_set_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"microsoft_defender": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"log_analytics_workspace_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"oms_agent": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"log_analytics_workspace_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
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

			"open_service_mesh_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"private_cluster_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"private_fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

			"kubernetes_version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"kube_admin_config": {
				Type:      pluginsdk.TypeList,
				Computed:  true,
				Sensitive: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"host": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"username": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"password": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"client_certificate": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"client_key": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"cluster_ca_certificate": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"kube_admin_config_raw": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"kube_config": {
				Type:      pluginsdk.TypeList,
				Computed:  true,
				Sensitive: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"host": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"username": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"password": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"client_certificate": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"client_key": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"cluster_ca_certificate": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"kube_config_raw": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"kubelet_identity": {
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

			"linux_profile": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"admin_username": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"ssh_key": {
							Type:     pluginsdk.TypeList,
							Computed: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"key_data": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"windows_profile": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"admin_username": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"network_profile": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"network_plugin": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"network_policy": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"service_cidr": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"dns_service_ip": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"docker_bridge_cidr": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"pod_cidr": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"load_balancer_sku": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"node_resource_group": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"oidc_issuer_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"oidc_issuer_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"role_based_access_control_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"service_principal": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"client_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": commonschema.TagsDataSource(),
		},
	}

	return resource
}

func dataSourceKubernetesClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.KubernetesClustersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.ManagedClusterName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	profile, err := client.GetAccessProfile(ctx, id.ResourceGroup, id.ManagedClusterName, "clusterUser")
	if err != nil {
		return fmt.Errorf("retrieving Access Profile for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.ManagedClusterProperties; props != nil {
		d.Set("dns_prefix", props.DNSPrefix)
		d.Set("fqdn", props.Fqdn)
		d.Set("disk_encryption_set_id", props.DiskEncryptionSetID)
		d.Set("private_fqdn", props.PrivateFQDN)
		d.Set("kubernetes_version", props.KubernetesVersion)
		d.Set("node_resource_group", props.NodeResourceGroup)

		if accessProfile := props.APIServerAccessProfile; accessProfile != nil {
			apiServerAuthorizedIPRanges := utils.FlattenStringSlice(accessProfile.AuthorizedIPRanges)
			if err := d.Set("api_server_authorized_ip_ranges", apiServerAuthorizedIPRanges); err != nil {
				return fmt.Errorf("setting `api_server_authorized_ip_ranges`: %+v", err)
			}

			d.Set("private_cluster_enabled", accessProfile.EnablePrivateCluster)
		}

		addOns := flattenKubernetesClusterDataSourceAddOns(props.AddonProfiles)
		d.Set("aci_connector_linux", addOns["aci_connector_linux"])
		d.Set("azure_policy_enabled", addOns["azure_policy_enabled"].(bool))
		d.Set("http_application_routing_enabled", addOns["http_application_routing_enabled"].(bool))
		d.Set("http_application_routing_zone_name", addOns["http_application_routing_zone_name"])
		d.Set("oms_agent", addOns["oms_agent"])
		d.Set("ingress_application_gateway", addOns["ingress_application_gateway"])
		d.Set("open_service_mesh_enabled", addOns["open_service_mesh_enabled"].(bool))
		d.Set("key_vault_secrets_provider", addOns["key_vault_secrets_provider"])

		agentPoolProfiles := flattenKubernetesClusterDataSourceAgentPoolProfiles(props.AgentPoolProfiles)
		if err := d.Set("agent_pool_profile", agentPoolProfiles); err != nil {
			return fmt.Errorf("setting `agent_pool_profile`: %+v", err)
		}

		kubeletIdentity, err := flattenKubernetesClusterDataSourceIdentityProfile(props.IdentityProfile)
		if err != nil {
			return err
		}
		if err := d.Set("kubelet_identity", kubeletIdentity); err != nil {
			return fmt.Errorf("setting `kubelet_identity`: %+v", err)
		}

		linuxProfile := flattenKubernetesClusterDataSourceLinuxProfile(props.LinuxProfile)
		if err := d.Set("linux_profile", linuxProfile); err != nil {
			return fmt.Errorf("setting `linux_profile`: %+v", err)
		}

		windowsProfile := flattenKubernetesClusterDataSourceWindowsProfile(props.WindowsProfile)
		if err := d.Set("windows_profile", windowsProfile); err != nil {
			return fmt.Errorf("setting `windows_profile`: %+v", err)
		}

		networkProfile := flattenKubernetesClusterDataSourceNetworkProfile(props.NetworkProfile)
		if err := d.Set("network_profile", networkProfile); err != nil {
			return fmt.Errorf("setting `network_profile`: %+v", err)
		}

		oidcIssuerEnabled := false
		oidcIssuerUrl := ""
		if props.OidcIssuerProfile != nil {
			if props.OidcIssuerProfile.Enabled != nil {
				oidcIssuerEnabled = *props.OidcIssuerProfile.Enabled
			}
			if props.OidcIssuerProfile.IssuerURL != nil {
				oidcIssuerUrl = *props.OidcIssuerProfile.IssuerURL
			}
		}

		if err := d.Set("oidc_issuer_enabled", oidcIssuerEnabled); err != nil {
			return fmt.Errorf("setting `oidc_issuer_enabled`: %+v", err)
		}
		if err := d.Set("oidc_issuer_url", oidcIssuerUrl); err != nil {
			return fmt.Errorf("setting `oidc_issuer_url`: %+v", err)
		}

		rbacEnabled := true
		if props.EnableRBAC != nil {
			rbacEnabled = *props.EnableRBAC
		}
		d.Set("role_based_access_control_enabled", rbacEnabled)

		microsoftDefender := flattenKubernetesClusterDataSourceMicrosoftDefender(props.SecurityProfile)
		if err := d.Set("microsoft_defender", microsoftDefender); err != nil {
			return fmt.Errorf("setting `microsoft_defender`: %+v", err)
		}

		aadRbac := flattenKubernetesClusterDataSourceAzureActiveDirectoryRoleBasedAccessControl(props)
		if err := d.Set("azure_active_directory_role_based_access_control", aadRbac); err != nil {
			return fmt.Errorf("setting `azure_active_directory_role_based_access_control`: %+v", err)
		}

		servicePrincipal := flattenKubernetesClusterDataSourceServicePrincipalProfile(props.ServicePrincipalProfile)
		if err := d.Set("service_principal", servicePrincipal); err != nil {
			return fmt.Errorf("setting `service_principal`: %+v", err)
		}

		// adminProfile is only available for RBAC enabled clusters with AAD and without local accounts disabled
		if props.AadProfile != nil && (props.DisableLocalAccounts == nil || !*props.DisableLocalAccounts) {
			adminProfile, err := client.GetAccessProfile(ctx, id.ResourceGroup, id.ManagedClusterName, "clusterAdmin")
			if err != nil {
				return fmt.Errorf("retrieving Admin Access Profile for %s: %+v", id, err)
			}

			adminKubeConfigRaw, adminKubeConfig := flattenKubernetesClusterAccessProfile(adminProfile)
			d.Set("kube_admin_config_raw", adminKubeConfigRaw)
			if err := d.Set("kube_admin_config", adminKubeConfig); err != nil {
				return fmt.Errorf("setting `kube_admin_config`: %+v", err)
			}
		} else {
			d.Set("kube_admin_config_raw", "")
			d.Set("kube_admin_config", []interface{}{})
		}
	}

	identity, err := flattenClusterDataSourceIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	kubeConfigRaw, kubeConfig := flattenKubernetesClusterDataSourceAccessProfile(profile)
	d.Set("kube_config_raw", kubeConfigRaw)
	if err := d.Set("kube_config", kubeConfig); err != nil {
		return fmt.Errorf("setting `kube_config`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenKubernetesClusterDataSourceAccessProfile(profile containerservice.ManagedClusterAccessProfile) (*string, []interface{}) {
	if profile.AccessProfile == nil {
		return nil, []interface{}{}
	}

	if kubeConfigRaw := profile.AccessProfile.KubeConfig; kubeConfigRaw != nil {
		rawConfig := string(*kubeConfigRaw)
		var flattenedKubeConfig []interface{}

		if strings.Contains(rawConfig, "apiserver-id:") {
			kubeConfigAAD, err := kubernetes.ParseKubeConfigAAD(rawConfig)
			if err != nil {
				return utils.String(rawConfig), []interface{}{}
			}

			flattenedKubeConfig = flattenKubernetesClusterDataSourceKubeConfigAAD(*kubeConfigAAD)
		} else {
			kubeConfig, err := kubernetes.ParseKubeConfig(rawConfig)
			if err != nil {
				return utils.String(rawConfig), []interface{}{}
			}

			flattenedKubeConfig = flattenKubernetesClusterDataSourceKubeConfig(*kubeConfig)
		}

		return utils.String(rawConfig), flattenedKubeConfig
	}

	return nil, []interface{}{}
}

func flattenKubernetesClusterDataSourceAddOns(profile map[string]*containerservice.ManagedClusterAddonProfile) map[string]interface{} {
	aciConnectors := make([]interface{}, 0)
	if aciConnector := kubernetesAddonProfileLocate(profile, aciConnectorKey); aciConnector != nil {
		if enabled := aciConnector.Enabled; enabled != nil && *enabled {
			subnetName := ""
			if v := aciConnector.Config["SubnetName"]; v != nil {
				subnetName = *v
			}

			aciConnectors = append(aciConnectors, map[string]interface{}{
				"subnet_name": subnetName,
			})
		}

	}

	azurePolicyEnabled := false
	if azurePolicy := kubernetesAddonProfileLocate(profile, azurePolicyKey); azurePolicy != nil {
		if enabledVal := azurePolicy.Enabled; enabledVal != nil {
			azurePolicyEnabled = *enabledVal
		}
	}

	httpApplicationRoutingEnabled := false
	httpApplicationRoutingZone := ""
	if httpApplicationRouting := kubernetesAddonProfileLocate(profile, httpApplicationRoutingKey); httpApplicationRouting != nil {
		if enabledVal := httpApplicationRouting.Enabled; enabledVal != nil {
			httpApplicationRoutingEnabled = *enabledVal
		}

		if v := kubernetesAddonProfilelocateInConfig(httpApplicationRouting.Config, "HTTPApplicationRoutingZoneName"); v != nil {
			httpApplicationRoutingZone = *v
		}
	}

	omsAgents := make([]interface{}, 0)
	if omsAgent := kubernetesAddonProfileLocate(profile, omsAgentKey); omsAgent != nil {
		if enabled := omsAgent.Enabled; enabled != nil && *enabled {
			workspaceID := ""
			if v := kubernetesAddonProfilelocateInConfig(omsAgent.Config, "logAnalyticsWorkspaceResourceID"); v != nil {
				if lawid, err := laparse.LogAnalyticsWorkspaceID(*v); err == nil {
					workspaceID = lawid.ID()
				}
			}

			omsAgentIdentity := flattenKubernetesClusterAddOnIdentityProfile(omsAgent.Identity)

			omsAgents = append(omsAgents, map[string]interface{}{
				"log_analytics_workspace_id": workspaceID,
				"oms_agent_identity":         omsAgentIdentity,
			})
		}
	}

	ingressApplicationGateways := make([]interface{}, 0)
	if ingressApplicationGateway := kubernetesAddonProfileLocate(profile, ingressApplicationGatewayKey); ingressApplicationGateway != nil {
		if enabled := ingressApplicationGateway.Enabled; enabled != nil && *enabled {
			gatewayId := ""
			if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "applicationGatewayId"); v != nil {
				gatewayId = *v
			}

			gatewayName := ""
			if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "applicationGatewayName"); v != nil {
				gatewayName = *v
			}

			effectiveGatewayId := ""
			if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "effectiveApplicationGatewayId"); v != nil {
				effectiveGatewayId = *v
			}

			subnetCIDR := ""
			if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "subnetCIDR"); v != nil {
				subnetCIDR = *v
			}

			subnetId := ""
			if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "subnetId"); v != nil {
				subnetId = *v
			}

			ingressApplicationGatewayIdentity := flattenKubernetesClusterAddOnIdentityProfile(ingressApplicationGateway.Identity)

			ingressApplicationGateways = append(ingressApplicationGateways, map[string]interface{}{
				"gateway_id":                           gatewayId,
				"gateway_name":                         gatewayName,
				"effective_gateway_id":                 effectiveGatewayId,
				"subnet_cidr":                          subnetCIDR,
				"subnet_id":                            subnetId,
				"ingress_application_gateway_identity": ingressApplicationGatewayIdentity,
			})
		}
	}

	openServiceMeshEnabled := false
	if openServiceMesh := kubernetesAddonProfileLocate(profile, openServiceMeshKey); openServiceMesh != nil {
		if enabledVal := openServiceMesh.Enabled; enabledVal != nil {
			openServiceMeshEnabled = *enabledVal
		}
	}

	azureKeyVaultSecretsProviders := make([]interface{}, 0)
	if azureKeyVaultSecretsProvider := kubernetesAddonProfileLocate(profile, azureKeyvaultSecretsProviderKey); azureKeyVaultSecretsProvider != nil {
		if enabled := azureKeyVaultSecretsProvider.Enabled; enabled != nil && *enabled {
			enableSecretRotation := false
			if v := kubernetesAddonProfilelocateInConfig(azureKeyVaultSecretsProvider.Config, "enableSecretRotation"); v != nil && *v != "false" {
				enableSecretRotation = true
			}

			rotationPollInterval := ""
			if v := kubernetesAddonProfilelocateInConfig(azureKeyVaultSecretsProvider.Config, "rotationPollInterval"); v != nil {
				rotationPollInterval = *v
			}

			azureKeyvaultSecretsProviderIdentity := flattenKubernetesClusterAddOnIdentityProfile(azureKeyVaultSecretsProvider.Identity)

			azureKeyVaultSecretsProviders = append(azureKeyVaultSecretsProviders, map[string]interface{}{
				"secret_rotation_enabled":  enableSecretRotation,
				"secret_rotation_interval": rotationPollInterval,
				"secret_identity":          azureKeyvaultSecretsProviderIdentity,
			})
		}
	}

	return map[string]interface{}{
		"aci_connector_linux":                aciConnectors,
		"azure_policy_enabled":               azurePolicyEnabled,
		"http_application_routing_enabled":   httpApplicationRoutingEnabled,
		"http_application_routing_zone_name": httpApplicationRoutingZone,
		"oms_agent":                          omsAgents,
		"ingress_application_gateway":        ingressApplicationGateways,
		"open_service_mesh_enabled":          openServiceMeshEnabled,
		"key_vault_secrets_provider":         azureKeyVaultSecretsProviders,
	}
}

func flattenKubernetesClusterDataSourceAgentPoolProfiles(input *[]containerservice.ManagedClusterAgentPoolProfile) []interface{} {
	agentPoolProfiles := make([]interface{}, 0)

	if input == nil {
		return agentPoolProfiles
	}

	for _, profile := range *input {
		count := 0
		if profile.Count != nil {
			count = int(*profile.Count)
		}

		enableNodePublicIP := false
		if profile.EnableNodePublicIP != nil {
			enableNodePublicIP = *profile.EnableNodePublicIP
		}

		minCount := 0
		if profile.MinCount != nil {
			minCount = int(*profile.MinCount)
		}

		maxCount := 0
		if profile.MaxCount != nil {
			maxCount = int(*profile.MaxCount)
		}

		enableAutoScaling := false
		if profile.EnableAutoScaling != nil {
			enableAutoScaling = *profile.EnableAutoScaling
		}

		name := ""
		if profile.Name != nil {
			name = *profile.Name
		}

		nodePublicIPPrefixID := ""
		if profile.NodePublicIPPrefixID != nil {
			nodePublicIPPrefixID = *profile.NodePublicIPPrefixID
		}

		osDiskSizeGb := 0
		if profile.OsDiskSizeGB != nil {
			osDiskSizeGb = int(*profile.OsDiskSizeGB)
		}

		vnetSubnetId := ""
		if profile.VnetSubnetID != nil {
			vnetSubnetId = *profile.VnetSubnetID
		}

		orchestratorVersion := ""
		if profile.OrchestratorVersion != nil && *profile.OrchestratorVersion != "" {
			orchestratorVersion = *profile.OrchestratorVersion
		}

		maxPods := 0
		if profile.MaxPods != nil {
			maxPods = int(*profile.MaxPods)
		}

		nodeLabels := make(map[string]string)
		if profile.NodeLabels != nil {
			for k, v := range profile.NodeLabels {
				if v == nil {
					continue
				}

				nodeLabels[k] = *v
			}
		}

		nodeTaints := make([]string, 0)
		if profile.NodeTaints != nil {
			nodeTaints = *profile.NodeTaints
		}

		vmSize := ""
		if profile.VMSize != nil {
			vmSize = *profile.VMSize
		}

		out := map[string]interface{}{
			"count":                    count,
			"enable_auto_scaling":      enableAutoScaling,
			"enable_node_public_ip":    enableNodePublicIP,
			"max_count":                maxCount,
			"max_pods":                 maxPods,
			"min_count":                minCount,
			"name":                     name,
			"node_labels":              nodeLabels,
			"node_public_ip_prefix_id": nodePublicIPPrefixID,
			"node_taints":              nodeTaints,
			"orchestrator_version":     orchestratorVersion,
			"os_disk_size_gb":          osDiskSizeGb,
			"os_type":                  string(profile.OsType),
			"tags":                     tags.Flatten(profile.Tags),
			"type":                     string(profile.Type),
			"upgrade_settings":         flattenUpgradeSettings(profile.UpgradeSettings),
			"vm_size":                  vmSize,
			"vnet_subnet_id":           vnetSubnetId,
			"zones":                    zones.Flatten(profile.AvailabilityZones),
		}
		agentPoolProfiles = append(agentPoolProfiles, out)
	}

	return agentPoolProfiles
}

func flattenKubernetesClusterDataSourceAzureActiveDirectoryRoleBasedAccessControl(input *containerservice.ManagedClusterProperties) []interface{} {
	results := make([]interface{}, 0)
	if profile := input.AadProfile; profile != nil {
		adminGroupObjectIds := utils.FlattenStringSlice(profile.AdminGroupObjectIDs)

		clientAppId := ""
		if profile.ClientAppID != nil {
			clientAppId = *profile.ClientAppID
		}

		managed := false
		if profile.Managed != nil {
			managed = *profile.Managed
		}

		azureRbacEnabled := false
		if profile.EnableAzureRBAC != nil {
			azureRbacEnabled = *profile.EnableAzureRBAC
		}

		serverAppId := ""
		if profile.ServerAppID != nil {
			serverAppId = *profile.ServerAppID
		}

		tenantId := ""
		if profile.TenantID != nil {
			tenantId = *profile.TenantID
		}

		results = append(results, map[string]interface{}{
			"admin_group_object_ids": adminGroupObjectIds,
			"client_app_id":          clientAppId,
			"managed":                managed,
			"server_app_id":          serverAppId,
			"tenant_id":              tenantId,
			"azure_rbac_enabled":     azureRbacEnabled,
		})
	}

	return results
}

func flattenKubernetesClusterDataSourceIdentityProfile(profile map[string]*containerservice.UserAssignedIdentity) ([]interface{}, error) {
	if profile == nil {
		return []interface{}{}, nil
	}

	kubeletIdentity := make([]interface{}, 0)
	if kubeletidentity := profile["kubeletidentity"]; kubeletidentity != nil {
		clientId := ""
		if clientid := kubeletidentity.ClientID; clientid != nil {
			clientId = *clientid
		}

		objectId := ""
		if objectid := kubeletidentity.ObjectID; objectid != nil {
			objectId = *objectid
		}

		userAssignedIdentityId := ""
		if resourceid := kubeletidentity.ResourceID; resourceid != nil {
			parsedId, err := commonids.ParseUserAssignedIdentityIDInsensitively(*resourceid)
			if err != nil {
				return nil, err
			}
			userAssignedIdentityId = parsedId.ID()
		}

		kubeletIdentity = append(kubeletIdentity, map[string]interface{}{
			"client_id":                 clientId,
			"object_id":                 objectId,
			"user_assigned_identity_id": userAssignedIdentityId,
		})
	}

	return kubeletIdentity, nil
}

func flattenKubernetesClusterDataSourceLinuxProfile(input *containerservice.LinuxProfile) []interface{} {
	values := make(map[string]interface{})
	sshKeys := make([]interface{}, 0)

	if profile := input; profile != nil {
		if username := profile.AdminUsername; username != nil {
			values["admin_username"] = *username
		}

		if ssh := profile.SSH; ssh != nil {
			if keys := ssh.PublicKeys; keys != nil {
				for _, sshKey := range *keys {
					if keyData := sshKey.KeyData; keyData != nil {
						outputs := make(map[string]interface{})
						outputs["key_data"] = *keyData
						sshKeys = append(sshKeys, outputs)
					}
				}
			}
		}
	}

	values["ssh_key"] = sshKeys

	return []interface{}{values}
}

func flattenKubernetesClusterDataSourceWindowsProfile(input *containerservice.ManagedClusterWindowsProfile) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	values := make(map[string]interface{})

	if username := input.AdminUsername; username != nil {
		values["admin_username"] = *username
	}

	return []interface{}{values}
}

func flattenKubernetesClusterDataSourceNetworkProfile(profile *containerservice.NetworkProfile) []interface{} {
	values := make(map[string]interface{})

	values["network_plugin"] = profile.NetworkPlugin

	if profile.NetworkPolicy != "" {
		values["network_policy"] = string(profile.NetworkPolicy)
	}

	if profile.ServiceCidr != nil {
		values["service_cidr"] = *profile.ServiceCidr
	}

	if profile.DNSServiceIP != nil {
		values["dns_service_ip"] = *profile.DNSServiceIP
	}

	if profile.DockerBridgeCidr != nil {
		values["docker_bridge_cidr"] = *profile.DockerBridgeCidr
	}

	if profile.PodCidr != nil {
		values["pod_cidr"] = *profile.PodCidr
	}

	if profile.LoadBalancerSku != "" {
		values["load_balancer_sku"] = string(profile.LoadBalancerSku)
	}

	return []interface{}{values}
}

func flattenKubernetesClusterDataSourceServicePrincipalProfile(profile *containerservice.ManagedClusterServicePrincipalProfile) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	values := make(map[string]interface{})

	if clientID := profile.ClientID; clientID != nil {
		values["client_id"] = *clientID
	}

	return []interface{}{values}
}

func flattenKubernetesClusterDataSourceKubeConfig(config kubernetes.KubeConfig) []interface{} {
	values := make(map[string]interface{})

	cluster := config.Clusters[0].Cluster
	user := config.Users[0].User
	name := config.Users[0].Name

	values["host"] = cluster.Server
	values["username"] = name
	values["password"] = user.Token
	values["client_certificate"] = user.ClientCertificteData
	values["client_key"] = user.ClientKeyData
	values["cluster_ca_certificate"] = cluster.ClusterAuthorityData

	return []interface{}{values}
}

func flattenKubernetesClusterDataSourceKubeConfigAAD(config kubernetes.KubeConfigAAD) []interface{} {
	values := make(map[string]interface{})

	cluster := config.Clusters[0].Cluster
	name := config.Users[0].Name

	values["host"] = cluster.Server
	values["username"] = name

	values["password"] = ""
	values["client_certificate"] = ""
	values["client_key"] = ""

	values["cluster_ca_certificate"] = cluster.ClusterAuthorityData

	return []interface{}{values}
}

func flattenClusterDataSourceIdentity(input *containerservice.ManagedClusterIdentity) (*[]interface{}, error) {
	var transform *identity.SystemOrUserAssignedMap

	if input != nil {
		transform = &identity.SystemOrUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		if input.PrincipalID != nil {
			transform.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transform.TenantId = *input.TenantID
		}
		for k, v := range input.UserAssignedIdentities {
			transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				ClientId:    v.ClientID,
				PrincipalId: v.PrincipalID,
			}
		}
	}

	return identity.FlattenSystemOrUserAssignedMap(transform)
}

func flattenKubernetesClusterDataSourceMicrosoftDefender(input *containerservice.ManagedClusterSecurityProfile) []interface{} {
	if input == nil || input.AzureDefender == nil || (input.AzureDefender.Enabled != nil && !*input.AzureDefender.Enabled) {
		return []interface{}{}
	}

	logAnalyticsWorkspace := ""
	if v := input.AzureDefender.LogAnalyticsWorkspaceResourceID; v != nil {
		logAnalyticsWorkspace = *v
	}

	return []interface{}{
		map[string]interface{}{
			"log_analytics_workspace_id": logAnalyticsWorkspace,
		},
	}
}
