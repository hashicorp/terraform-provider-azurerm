// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-04-02-preview/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/kubernetes"
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

			"current_kubernetes_version": {
				Type:     pluginsdk.TypeString,
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

			"custom_ca_trust_certificates_base64": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
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
						"msi_auth_for_monitoring_enabled": {
							Type:     pluginsdk.TypeBool,
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

			"key_management_service": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_vault_key_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"key_vault_network_access": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

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

			"node_resource_group_id": {
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

			"storage_profile": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"blob_driver_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
						"disk_driver_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
						"disk_driver_version": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"file_driver_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
						"snapshot_controller_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"service_mesh_profile": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"mode": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"internal_ingress_gateway_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
						"external_ingress_gateway_enabled": {
							Type:     pluginsdk.TypeBool,
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

	id := commonids.NewKubernetesClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	userCredentialsResp, err := client.ListClusterUserCredentials(ctx, id, managedclusters.ListClusterUserCredentialsOperationOptions{})
	// only raise the error if it's not a limited permissions error, since this is the Data Source
	if err != nil && !response.WasStatusCode(userCredentialsResp.HttpResponse, http.StatusForbidden) {
		return fmt.Errorf("retrieving User Credentials for %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if model := resp.Model; model != nil {
		d.Set("name", id.ManagedClusterName)
		d.Set("resource_group_name", id.ResourceGroupName)
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("dns_prefix", pointer.From(props.DnsPrefix))
			d.Set("fqdn", pointer.From(props.Fqdn))
			d.Set("disk_encryption_set_id", pointer.From(props.DiskEncryptionSetID))
			d.Set("private_fqdn", pointer.From(props.PrivateFQDN))
			d.Set("kubernetes_version", pointer.From(props.KubernetesVersion))
			d.Set("current_kubernetes_version", pointer.From(props.CurrentKubernetesVersion))

			nodeResourceGroup := ""
			if v := props.NodeResourceGroup; v != nil {
				nodeResourceGroup = *props.NodeResourceGroup
			}
			d.Set("node_resource_group", nodeResourceGroup)

			nodeResourceGroupId := commonids.NewResourceGroupID(id.SubscriptionId, nodeResourceGroup)
			d.Set("node_resource_group_id", nodeResourceGroupId.ID())

			if accessProfile := props.ApiServerAccessProfile; accessProfile != nil {
				apiServerAuthorizedIPRanges := utils.FlattenStringSlice(accessProfile.AuthorizedIPRanges)
				if err := d.Set("api_server_authorized_ip_ranges", apiServerAuthorizedIPRanges); err != nil {
					return fmt.Errorf("setting `api_server_authorized_ip_ranges`: %+v", err)
				}

				d.Set("private_cluster_enabled", accessProfile.EnablePrivateCluster)
			}

			if addonProfiles := props.AddonProfiles; addonProfiles != nil {
				addOns := flattenKubernetesClusterDataSourceAddOns(*addonProfiles)
				d.Set("aci_connector_linux", addOns["aci_connector_linux"])
				d.Set("azure_policy_enabled", addOns["azure_policy_enabled"].(bool))
				d.Set("http_application_routing_enabled", addOns["http_application_routing_enabled"].(bool))
				d.Set("http_application_routing_zone_name", addOns["http_application_routing_zone_name"])
				d.Set("oms_agent", addOns["oms_agent"])
				d.Set("ingress_application_gateway", addOns["ingress_application_gateway"])
				d.Set("open_service_mesh_enabled", addOns["open_service_mesh_enabled"].(bool))
				d.Set("key_vault_secrets_provider", addOns["key_vault_secrets_provider"])
			}

			agentPoolProfiles := flattenKubernetesClusterDataSourceAgentPoolProfiles(props.AgentPoolProfiles)
			if err := d.Set("agent_pool_profile", agentPoolProfiles); err != nil {
				return fmt.Errorf("setting `agent_pool_profile`: %+v", err)
			}

			azureKeyVaultKms := flattenKubernetesClusterDataSourceKeyVaultKms(props.SecurityProfile)
			if err := d.Set("key_management_service", azureKeyVaultKms); err != nil {
				return fmt.Errorf("setting `key_management_service`: %+v", err)
			}

			customCaTrustCertList := flattenCustomCaTrustCerts(props.SecurityProfile.CustomCATrustCertificates)
			if err := d.Set("custom_ca_trust_certificates_base64", customCaTrustCertList); err != nil {
				return fmt.Errorf("setting `custom_ca_trust_certificates_base64`: %+v", err)
			}

			serviceMeshProfile := flattenKubernetesClusterAzureServiceMeshProfile(props.ServiceMeshProfile)
			if err := d.Set("service_mesh_profile", serviceMeshProfile); err != nil {
				return fmt.Errorf("setting `service_mesh_profile`: %+v", err)
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

			storageProfile := flattenKubernetesClusterDataSourceStorageProfile(props.StorageProfile)
			if err := d.Set("storage_profile", storageProfile); err != nil {
				return fmt.Errorf("setting `storage_profile`: %+v", err)
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
			adminKubeConfig := make([]interface{}, 0)
			var adminKubeConfigRaw *string
			if props.AadProfile != nil && (props.DisableLocalAccounts == nil || !*props.DisableLocalAccounts) {
				adminCredentialsResp, err := client.ListClusterAdminCredentials(ctx, id, managedclusters.ListClusterAdminCredentialsOperationOptions{})
				// only raise the error if it's not a limited permissions error, since this is the Data Source
				if err != nil && !response.WasStatusCode(adminCredentialsResp.HttpResponse, http.StatusForbidden) {
					return fmt.Errorf("retrieving Admin Credentials for %s: %+v", id, err)
				}

				adminKubeConfigRaw, adminKubeConfig = flattenKubernetesClusterCredentials(adminCredentialsResp.Model, "clusterAdmin")
			}
			d.Set("kube_admin_config_raw", adminKubeConfigRaw)
			if err := d.Set("kube_admin_config", adminKubeConfig); err != nil {
				return fmt.Errorf("setting `kube_admin_config`: %+v", err)
			}
		}

		identity, err := flattenClusterDataSourceIdentity(model.Identity)
		if err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		kubeConfigRaw, kubeConfig := flattenKubernetesClusterCredentials(userCredentialsResp.Model, "clusterUser")
		d.Set("kube_config_raw", kubeConfigRaw)
		if err := d.Set("kube_config", kubeConfig); err != nil {
			return fmt.Errorf("setting `kube_config`: %+v", err)
		}

		d.Set("tags", tags.Flatten(model.Tags))
	}

	return nil
}

func flattenKubernetesClusterDataSourceKeyVaultKms(input *managedclusters.ManagedClusterSecurityProfile) []interface{} {
	azureKeyVaultKms := make([]interface{}, 0)

	if input != nil && input.AzureKeyVaultKms != nil && input.AzureKeyVaultKms.Enabled != nil && *input.AzureKeyVaultKms.Enabled {
		keyId := ""
		if v := input.AzureKeyVaultKms.KeyId; v != nil {
			keyId = *v
		}

		networkAccess := ""
		if v := input.AzureKeyVaultKms.KeyVaultNetworkAccess; v != nil {
			networkAccess = string(*v)
		}

		azureKeyVaultKms = append(azureKeyVaultKms, map[string]interface{}{
			"key_vault_key_id":         keyId,
			"key_vault_network_access": networkAccess,
		})
	}

	return azureKeyVaultKms
}

func flattenKubernetesClusterDataSourceStorageProfile(input *managedclusters.ManagedClusterStorageProfile) []interface{} {
	storageProfile := make([]interface{}, 0)

	if input != nil {
		blobEnabled := false
		if input.BlobCSIDriver != nil && input.BlobCSIDriver.Enabled != nil {
			blobEnabled = *input.BlobCSIDriver.Enabled
		}

		diskEnabled := true
		if input.DiskCSIDriver != nil && input.DiskCSIDriver.Enabled != nil {
			diskEnabled = *input.DiskCSIDriver.Enabled
		}

		diskVersion := ""
		if input.DiskCSIDriver != nil && input.DiskCSIDriver.Version != nil {
			diskVersion = *input.DiskCSIDriver.Version
		}

		fileEnabled := true
		if input.FileCSIDriver != nil && input.FileCSIDriver.Enabled != nil {
			fileEnabled = *input.FileCSIDriver.Enabled
		}

		snapshotController := true
		if input.SnapshotController != nil && input.SnapshotController.Enabled != nil {
			snapshotController = *input.SnapshotController.Enabled
		}

		storageProfile = append(storageProfile, map[string]interface{}{
			"blob_driver_enabled":         blobEnabled,
			"disk_driver_enabled":         diskEnabled,
			"disk_driver_version":         diskVersion,
			"file_driver_enabled":         fileEnabled,
			"snapshot_controller_enabled": snapshotController,
		})
	}

	return storageProfile
}

func flattenKubernetesClusterCredentials(model *managedclusters.CredentialResults, configName string) (*string, []interface{}) {
	if model == nil || model.Kubeconfigs == nil || len(*model.Kubeconfigs) < 1 {
		return nil, []interface{}{}
	}

	for _, c := range *model.Kubeconfigs {
		if c.Name == nil || *c.Name != configName {
			continue
		}
		if kubeConfigRaw := c.Value; kubeConfigRaw != nil {
			rawConfig := *kubeConfigRaw
			if base64IsEncoded(*kubeConfigRaw) {
				rawConfig = base64Decode(*kubeConfigRaw)
			}

			var flattenedKubeConfig []interface{}

			if strings.Contains(rawConfig, "apiserver-id:") || strings.Contains(rawConfig, "exec") {
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
	}

	return nil, []interface{}{}
}

func flattenKubernetesClusterDataSourceAddOns(profile map[string]managedclusters.ManagedClusterAddonProfile) map[string]interface{} {
	aciConnectors := make([]interface{}, 0)
	aciConnector := kubernetesAddonProfileLocate(profile, aciConnectorKey)
	if enabled := aciConnector.Enabled; enabled {
		subnetName := ""
		if v := aciConnector.Config; v != nil && (*v)["SubnetName"] != "" {
			subnetName = (*v)["SubnetName"]
		}

		aciConnectors = append(aciConnectors, map[string]interface{}{
			"subnet_name": subnetName,
		})
	}

	azurePolicyEnabled := false
	azurePolicy := kubernetesAddonProfileLocate(profile, azurePolicyKey)
	if enabledVal := azurePolicy.Enabled; enabledVal {
		azurePolicyEnabled = enabledVal
	}

	httpApplicationRoutingEnabled := false
	httpApplicationRoutingZone := ""
	httpApplicationRouting := kubernetesAddonProfileLocate(profile, httpApplicationRoutingKey)
	if enabledVal := httpApplicationRouting.Enabled; enabledVal {
		httpApplicationRoutingEnabled = enabledVal
	}

	if v := kubernetesAddonProfilelocateInConfig(httpApplicationRouting.Config, "HTTPApplicationRoutingZoneName"); v != "" {
		httpApplicationRoutingZone = v
	}

	omsAgents := make([]interface{}, 0)
	omsAgent := kubernetesAddonProfileLocate(profile, omsAgentKey)
	if enabled := omsAgent.Enabled; enabled {
		workspaceID := ""
		useAADAuth := false

		if v := kubernetesAddonProfilelocateInConfig(omsAgent.Config, "logAnalyticsWorkspaceResourceID"); v != "" {
			if lawid, err := workspaces.ParseWorkspaceID(v); err == nil {
				workspaceID = lawid.ID()
			}
		}

		if v := kubernetesAddonProfilelocateInConfig(omsAgent.Config, "useAADAuth"); v != "false" && v != "" {
			useAADAuth = true
		}

		omsAgentIdentity := flattenKubernetesClusterAddOnIdentityProfile(omsAgent.Identity)

		omsAgents = append(omsAgents, map[string]interface{}{
			"log_analytics_workspace_id":      workspaceID,
			"msi_auth_for_monitoring_enabled": useAADAuth,
			"oms_agent_identity":              omsAgentIdentity,
		})
	}

	ingressApplicationGateways := make([]interface{}, 0)
	ingressApplicationGateway := kubernetesAddonProfileLocate(profile, ingressApplicationGatewayKey)
	if enabled := ingressApplicationGateway.Enabled; enabled {
		gatewayId := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "applicationGatewayId"); v != "" {
			gatewayId = v
		}

		gatewayName := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "applicationGatewayName"); v != "" {
			gatewayName = v
		}

		effectiveGatewayId := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "effectiveApplicationGatewayId"); v != "" {
			effectiveGatewayId = v
		}

		subnetCIDR := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "subnetCIDR"); v != "" {
			subnetCIDR = v
		}

		subnetId := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "subnetId"); v != "" {
			subnetId = v
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

	openServiceMeshEnabled := false
	openServiceMesh := kubernetesAddonProfileLocate(profile, openServiceMeshKey)
	if enabledVal := openServiceMesh.Enabled; enabledVal {
		openServiceMeshEnabled = enabledVal
	}

	azureKeyVaultSecretsProviders := make([]interface{}, 0)
	azureKeyVaultSecretsProvider := kubernetesAddonProfileLocate(profile, azureKeyvaultSecretsProviderKey)
	if enabled := azureKeyVaultSecretsProvider.Enabled; enabled {
		enableSecretRotation := false
		if v := kubernetesAddonProfilelocateInConfig(azureKeyVaultSecretsProvider.Config, "enableSecretRotation"); v != "false" {
			enableSecretRotation = true
		}

		rotationPollInterval := ""
		if v := kubernetesAddonProfilelocateInConfig(azureKeyVaultSecretsProvider.Config, "rotationPollInterval"); v != "" {
			rotationPollInterval = v
		}

		azureKeyvaultSecretsProviderIdentity := flattenKubernetesClusterAddOnIdentityProfile(azureKeyVaultSecretsProvider.Identity)

		azureKeyVaultSecretsProviders = append(azureKeyVaultSecretsProviders, map[string]interface{}{
			"secret_rotation_enabled":  enableSecretRotation,
			"secret_rotation_interval": rotationPollInterval,
			"secret_identity":          azureKeyvaultSecretsProviderIdentity,
		})
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

func flattenKubernetesClusterDataSourceAgentPoolProfiles(input *[]managedclusters.ManagedClusterAgentPoolProfile) []interface{} {
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

		name := profile.Name

		nodePublicIPPrefixID := profile.NodePublicIPPrefixID

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
			for k, v := range *profile.NodeLabels {
				if v == "" {
					continue
				}

				nodeLabels[k] = v
			}
		}

		nodeTaints := make([]string, 0)
		if profile.NodeTaints != nil {
			nodeTaints = *profile.NodeTaints
		}

		vmSize := profile.VMSize

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
			"os_type":                  string(*profile.OsType),
			"tags":                     tags.Flatten(profile.Tags),
			"type":                     string(*profile.Type),
			"upgrade_settings":         flattenKubernetesClusterDataSourceUpgradeSettings(profile.UpgradeSettings),
			"vm_size":                  vmSize,
			"vnet_subnet_id":           vnetSubnetId,
			"zones":                    zones.FlattenUntyped(profile.AvailabilityZones),
		}
		agentPoolProfiles = append(agentPoolProfiles, out)
	}

	return agentPoolProfiles
}

func flattenKubernetesClusterDataSourceAzureActiveDirectoryRoleBasedAccessControl(input *managedclusters.ManagedClusterProperties) []interface{} {
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

func flattenKubernetesClusterDataSourceIdentityProfile(profile *map[string]managedclusters.UserAssignedIdentity) ([]interface{}, error) {
	if profile == nil || *profile == nil {
		return []interface{}{}, nil
	}

	kubeletIdentity := make([]interface{}, 0)
	kubeletidentity := (*profile)["kubeletidentity"]
	clientId := ""
	if clientid := kubeletidentity.ClientId; clientid != nil {
		clientId = *clientid
	}

	objectId := ""
	if objectid := kubeletidentity.ObjectId; objectid != nil {
		objectId = *objectid
	}

	userAssignedIdentityId := ""
	if resourceid := kubeletidentity.ResourceId; resourceid != nil {
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

	return kubeletIdentity, nil
}

func flattenKubernetesClusterDataSourceLinuxProfile(input *managedclusters.ContainerServiceLinuxProfile) []interface{} {
	values := make(map[string]interface{})
	sshKeys := make([]interface{}, 0)

	if profile := input; profile != nil {
		if username := profile.AdminUsername; username != "" {
			values["admin_username"] = username
		}

		ssh := profile.Ssh
		if keys := ssh.PublicKeys; keys != nil {
			for _, sshKey := range keys {
				if keyData := sshKey.KeyData; keyData != "" {
					outputs := make(map[string]interface{})
					outputs["key_data"] = keyData
					sshKeys = append(sshKeys, outputs)
				}
			}
		}
	}

	values["ssh_key"] = sshKeys

	return []interface{}{values}
}

func flattenKubernetesClusterDataSourceWindowsProfile(input *managedclusters.ManagedClusterWindowsProfile) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	values := make(map[string]interface{})

	if username := input.AdminUsername; username != "" {
		values["admin_username"] = username
	}

	return []interface{}{values}
}

func flattenKubernetesClusterDataSourceNetworkProfile(profile *managedclusters.ContainerServiceNetworkProfile) []interface{} {
	values := make(map[string]interface{})

	values["network_plugin"] = profile.NetworkPlugin

	if profile.NetworkPolicy != nil {
		values["network_policy"] = string(*profile.NetworkPolicy)
	}

	if profile.ServiceCidr != nil {
		values["service_cidr"] = *profile.ServiceCidr
	}

	if profile.DnsServiceIP != nil {
		values["dns_service_ip"] = *profile.DnsServiceIP
	}

	if profile.PodCidr != nil {
		values["pod_cidr"] = *profile.PodCidr
	}

	if profile.LoadBalancerSku != nil {
		values["load_balancer_sku"] = string(*profile.LoadBalancerSku)
	}

	return []interface{}{values}
}

func flattenKubernetesClusterDataSourceServicePrincipalProfile(profile *managedclusters.ManagedClusterServicePrincipalProfile) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	values := make(map[string]interface{})

	if clientID := profile.ClientId; clientID != "" {
		values["client_id"] = clientID
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

func flattenClusterDataSourceIdentity(input *identity.SystemOrUserAssignedMap) (*[]interface{}, error) {
	return identity.FlattenSystemOrUserAssignedMap(input)
}

func flattenKubernetesClusterDataSourceMicrosoftDefender(input *managedclusters.ManagedClusterSecurityProfile) []interface{} {
	if input == nil || input.Defender == nil || input.Defender.SecurityMonitoring == nil || (input.Defender.SecurityMonitoring.Enabled != nil && !*input.Defender.SecurityMonitoring.Enabled) {
		return []interface{}{}
	}

	logAnalyticsWorkspace := ""
	if v := input.Defender.LogAnalyticsWorkspaceResourceId; v != nil {
		logAnalyticsWorkspace = *v
	}

	return []interface{}{
		map[string]interface{}{
			"log_analytics_workspace_id": logAnalyticsWorkspace,
		},
	}
}

func flattenKubernetesClusterDataSourceUpgradeSettings(input *managedclusters.AgentPoolUpgradeSettings) []interface{} {
	maxSurge := ""
	if input != nil && input.MaxSurge != nil {
		maxSurge = *input.MaxSurge
	}

	if maxSurge == "" {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"max_surge": maxSurge,
		},
	}
}

func flattenCustomCaTrustCerts(input *[]string) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	customCaTrustCertInterface := make([]interface{}, len(*input))

	for index, value := range *input {
		customCaTrustCertInterface[index] = value
	}

	return customCaTrustCertInterface
}
