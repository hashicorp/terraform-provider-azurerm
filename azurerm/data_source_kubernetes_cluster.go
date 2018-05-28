package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2017-09-30/containerservice"
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

			"dns_prefix": {
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

						"dns_prefix": {
							Type:     schema.TypeString,
							Computed: true,
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
					},
				},
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

	profile, err := kubernetesClustersClient.GetAccessProfiles(ctx, resourceGroup, name, "clusterUser")
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

		linuxProfile := flattenKubernetesClusterDataSourceLinuxProfile(props.LinuxProfile)
		if err := d.Set("linux_profile", linuxProfile); err != nil {
			return fmt.Errorf("Error setting `linux_profile`: %+v", err)
		}

		agentPoolProfiles := flattenKubernetesClusterDataSourceAgentPoolProfiles(props.AgentPoolProfiles)
		if err := d.Set("agent_pool_profile", agentPoolProfiles); err != nil {
			return fmt.Errorf("Error setting `agent_pool_profile`: %+v", err)
		}

		servicePrincipal := flattenKubernetesClusterDataSourceServicePrincipalProfile(resp.ManagedClusterProperties.ServicePrincipalProfile)
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
						outputs := make(map[string]interface{}, 0)
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

func flattenKubernetesClusterDataSourceAgentPoolProfiles(input *[]containerservice.AgentPoolProfile) []interface{} {
	agentPoolProfiles := make([]interface{}, 0)

	if input == nil {
		return agentPoolProfiles
	}

	for _, profile := range *input {
		agentPoolProfile := make(map[string]interface{})

		if profile.Count != nil {
			agentPoolProfile["count"] = int(*profile.Count)
		}

		if profile.DNSPrefix != nil {
			agentPoolProfile["dns_prefix"] = *profile.DNSPrefix
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

		agentPoolProfiles = append(agentPoolProfiles, agentPoolProfile)
	}

	return agentPoolProfiles
}

func flattenKubernetesClusterDataSourceServicePrincipalProfile(profile *containerservice.ServicePrincipalProfile) []interface{} {
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

		kubeConfig, err := kubernetes.ParseKubeConfig(rawConfig)
		if err != nil {
			return utils.String(rawConfig), []interface{}{}
		}

		flattenedKubeConfig := flattenKubernetesClusterDataSourceKubeConfig(*kubeConfig)
		return utils.String(rawConfig), flattenedKubeConfig
	}

	return nil, []interface{}{}
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
