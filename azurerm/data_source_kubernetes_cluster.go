package azurerm

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2018-03-31/containerservice"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/kubernetes"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmKubernetesClusterRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"location": locationForDataSourceSchema(),

			"addon_profile": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http_application_routing": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"http_application_routing_zone_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},

						"oms_agent": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"log_analytics_workspace_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"agent_pool_profile": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"count": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						// TODO: remove this in a future version
						"dns_prefix": {
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "This field is no longer returned from the Azure API",
						},

						"vm_size": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"os_disk_size_gb": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"vnet_subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"max_pods": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"aad_profile": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_app_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_app_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"dns_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"enable_rbac": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"kubernetes_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"kube_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"client_certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_key": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"cluster_ca_certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"kube_config_raw": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"linux_profile": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admin_username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssh_key": {
							Type:     schema.TypeList,
							Computed: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key_data": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"network_profile": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_plugin": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"service_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"dns_service_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"docker_bridge_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"pod_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"node_resource_group": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"service_principal": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmKubernetesClusterRead(d *schema.ResourceData, meta interface{}) error {
	kubernetesClustersClient := meta.(*ArmClient).kubernetesClustersClient
	client := meta.(*ArmClient)

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	ctx := client.StopContext
	resp, err := kubernetesClustersClient.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: AKS Managed Cluster %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on AKS Managed Cluster %q (resource group %q): %+v", name, resourceGroup, err)
	}

	profile, err := kubernetesClustersClient.GetAccessProfile(ctx, resourceGroup, name, "clusterUser")
	if err != nil {
		return fmt.Errorf("Error getting access profile while making Read request on AKS Managed Cluster %q (resource group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.ManagedClusterProperties; props != nil {
		d.Set("dns_prefix", props.DNSPrefix)
		d.Set("fqdn", props.Fqdn)
		d.Set("kubernetes_version", props.KubernetesVersion)
		d.Set("enable_rbac", props.EnableRBAC)
		d.Set("node_resource_group", props.NodeResourceGroup)

		addonProfiles := flattenKubernetesClusterDataSourceAddonProfiles(props.AddonProfiles)
		if err := d.Set("addon_profile", addonProfiles); err != nil {
			return fmt.Errorf("Error setting `addon_profile`: %+v", err)
		}

		agentPoolProfiles := flattenKubernetesClusterDataSourceAgentPoolProfiles(props.AgentPoolProfiles)
		if err := d.Set("agent_pool_profile", agentPoolProfiles); err != nil {
			return fmt.Errorf("Error setting `agent_pool_profile`: %+v", err)
		}

		aadProfile := flattenKubernetesClusterDataSourceAadProfile(resp.ManagedClusterProperties.AadProfile)
		if err := d.Set("aad_profile", aadProfile); err != nil {
			return fmt.Errorf("Error setting `aad_profile`: %+v", err)
		}

		linuxProfile := flattenKubernetesClusterDataSourceLinuxProfile(props.LinuxProfile)
		if err := d.Set("linux_profile", linuxProfile); err != nil {
			return fmt.Errorf("Error setting `linux_profile`: %+v", err)
		}

		networkProfile := flattenKubernetesClusterDataSourceNetworkProfile(props.NetworkProfile)
		if err := d.Set("network_profile", networkProfile); err != nil {
			return fmt.Errorf("Error setting `network_profile`: %+v", err)
		}

		servicePrincipal := flattenKubernetesClusterDataSourceServicePrincipalProfile(props.ServicePrincipalProfile)
		if err := d.Set("service_principal", servicePrincipal); err != nil {
			return fmt.Errorf("Error setting `service_principal`: %+v", err)
		}
	}

	kubeConfigRaw, kubeConfig := flattenKubernetesClusterDataSourceAccessProfile(&profile)
	d.Set("kube_config_raw", kubeConfigRaw)

	if err := d.Set("kube_config", kubeConfig); err != nil {
		return fmt.Errorf("Error setting `kube_config`: %+v", err)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
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

func flattenKubernetesClusterDataSourceAgentPoolProfiles(input *[]containerservice.ManagedClusterAgentPoolProfile) []interface{} {
	agentPoolProfiles := make([]interface{}, 0)

	if input == nil {
		return agentPoolProfiles
	}

	for _, profile := range *input {
		agentPoolProfile := make(map[string]interface{})

		if profile.Count != nil {
			agentPoolProfile["count"] = int(*profile.Count)
		}

		if profile.Name != nil {
			agentPoolProfile["name"] = *profile.Name
		}

		if profile.VMSize != "" {
			agentPoolProfile["vm_size"] = string(profile.VMSize)
		}

		if profile.OsDiskSizeGB != nil {
			agentPoolProfile["os_disk_size_gb"] = int(*profile.OsDiskSizeGB)
		}

		if profile.VnetSubnetID != nil {
			agentPoolProfile["vnet_subnet_id"] = *profile.VnetSubnetID
		}

		if profile.OsType != "" {
			agentPoolProfile["os_type"] = string(profile.OsType)
		}

		if profile.MaxPods != nil {
			agentPoolProfile["max_pods"] = int(*profile.MaxPods)
		}

		agentPoolProfiles = append(agentPoolProfiles, agentPoolProfile)
	}

	return agentPoolProfiles
}

func flattenKubernetesClusterDataSourceServicePrincipalProfile(profile *containerservice.ManagedClusterServicePrincipalProfile) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	values := make(map[string]interface{})

	if clientId := profile.ClientID; clientId != nil {
		values["client_id"] = *clientId
	}

	return []interface{}{values}
}

func flattenKubernetesClusterDataSourceAccessProfile(profile *containerservice.ManagedClusterAccessProfile) (*string, []interface{}) {
	if profile == nil || profile.AccessProfile == nil {
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

func flattenKubernetesClusterDataSourceAadProfile(profile *containerservice.ManagedClusterAADProfile) []interface{} {
	if profile == nil {
		return nil
	}

	values := make(map[string]interface{})

	if serverAppId := profile.ServerAppID; serverAppId != nil {
		values["server_app_id"] = *serverAppId
	}

	if clientAppId := profile.ClientAppID; clientAppId != nil {
		values["client_app_id"] = *clientAppId
	}

	if tenantId := profile.TenantID; tenantId != nil {
		values["tenant_id"] = *tenantId
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

func flattenKubernetesClusterDataSourceNetworkProfile(profile *containerservice.NetworkProfile) []interface{} {
	values := make(map[string]interface{})

	values["network_plugin"] = profile.NetworkPlugin

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

	return []interface{}{values}
}

func flattenKubernetesClusterDataSourceAddonProfiles(profile map[string]*containerservice.ManagedClusterAddonProfile) interface{} {
	values := make(map[string]interface{})

	routes := make([]interface{}, 0)
	if httpApplicationRouting := profile["httpApplicationRouting"]; httpApplicationRouting != nil {
		enabled := false
		if enabledVal := httpApplicationRouting.Enabled; enabledVal != nil {
			enabled = *enabledVal
		}

		zoneName := ""
		if v := httpApplicationRouting.Config["HTTPApplicationRoutingZoneName"]; v != nil {
			zoneName = *v
		}

		output := map[string]interface{}{
			"enabled":                            enabled,
			"http_application_routing_zone_name": zoneName,
		}
		routes = append(routes, output)
	}
	values["http_application_routing"] = routes

	agents := make([]interface{}, 0)
	if omsAgent := profile["omsAgent"]; omsAgent != nil {
		enabled := false
		if enabledVal := omsAgent.Enabled; enabledVal != nil {
			enabled = *enabledVal
		}

		workspaceId := ""
		if workspaceResourceID := omsAgent.Config["logAnalyticsWorkspaceResourceID"]; workspaceResourceID != nil {
			workspaceId = *workspaceResourceID
		}

		output := map[string]interface{}{
			"enabled":                    enabled,
			"log_analytics_workspace_id": workspaceId,
		}
		agents = append(agents, output)
	}
	values["oms_agent"] = agents

	return []interface{}{values}
}
