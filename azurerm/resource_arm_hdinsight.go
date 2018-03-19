package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/hdinsight/mgmt/2015-03-01-preview/hdinsight"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmHDInsightCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmHDInsightClusterCreate,
		Read:   resourceArmHDInsightClusterRead,
		Update: resourceArmHDInsightClusterUpdate,
		Delete: resourceArmHDInsightClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"cluster_version": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: azureRmSuppressClusterVersionDiff,
			},

			"tier": {
				Type:     schema.TypeString,
				Default:  "Standard",
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Standard",
				}, false),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			//TODO: add support for more configurations i.e. core-site
			"cluster_definition": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blueprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"hadoop",
								"spark",
								"storm",
								"hbase",
								"rserver",
								"kafka",
								"interactivequery",
								"interactivehive",
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},
						"component_version": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"configurations": {
							Type:             schema.TypeList,
							Optional:         true,
							DiffSuppressFunc: azureRmSuppressSensitiveDiff,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"gateway": {
										Type:             schema.TypeList,
										Optional:         true,
										DiffSuppressFunc: azureRmSuppressSensitiveDiff,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rest_auth_credential_is_enabled": {
													Type:             schema.TypeBool,
													Optional:         true,
													DiffSuppressFunc: azureRmSuppressSensitiveDiff,
												},
												"rest_auth_credential_username": {
													Type:             schema.TypeString,
													Optional:         true,
													DiffSuppressFunc: azureRmSuppressSensitiveDiff,
												},
												"rest_auth_credential_password": {
													Type:             schema.TypeString,
													Optional:         true,
													DiffSuppressFunc: azureRmSuppressSensitiveDiff,
												},
											},
										},
									},
									"rserver": {
										Type:             schema.TypeList,
										Optional:         true,
										DiffSuppressFunc: azureRmSuppressSensitiveDiff,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rstudio": {
													Type:     schema.TypeBool,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"security_profile": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"directory_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"organizational_unit_dn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ldaps_urls": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"domain_username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"domain_password": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cluster_user_group_dns": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"compute_profile": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"roles": {
							Type:             schema.TypeList,
							Required:         true,
							DiffSuppressFunc: azureRmSuppressRolesDiff,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"min_instance_count": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      2,
										ValidateFunc: validation.IntBetween(1, 32),
									},
									"target_instance_count": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      3,
										ValidateFunc: validation.IntBetween(1, 32),
									},
									"hardware_profile": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vm_size": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
									"os_profile": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"linux_operating_system_profile": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"username": {
																Type:     schema.TypeString,
																Required: true,
															},
															"password": {
																Type:             schema.TypeString,
																Optional:         true,
																DiffSuppressFunc: azureRmSuppressSensitiveDiff,
															},
															"ssh_key": {
																Type:     schema.TypeList,
																Required: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"key_data": {
																			Type:             schema.TypeString,
																			Required:         true,
																			DiffSuppressFunc: azureRmSuppressSensitiveDiff,
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									"virtual_network_profile": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"virtual_network_id": {
													Type:     schema.TypeString,
													Required: true,
												},
												"subnet_name": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
									"data_disks_group": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"disks_per_node": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntBetween(0, 10),
												},
												"storage_account_type": {
													Type:     schema.TypeString,
													Optional: true,
													ValidateFunc: validation.StringInSlice([]string{
														"Standard_LRS",
														"Standard_GRS",
														"Standard_RAGRS",
														"Standard_ZRS",
													}, true),
												},
												"data_size_gb": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validateDiskSizeGB,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"storage_profile": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_accounts": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"is_default": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"container": {
										Type:     schema.TypeString,
										Required: true,
									},
									"key": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmHDInsightClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	hdInsightClustersClient := client.hdInsightClustersClient

	log.Printf("[INFO] preparing arguments for Azure ARM HDInsight creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resGroup := d.Get("resource_group_name").(string)
	clusterVersion := d.Get("cluster_version").(string)
	tierStr := d.Get("tier").(string)
	osType := hdinsight.OSType("Linux")
	tier := hdinsight.Tier(tierStr)

	clusterDefinition, err := expandAzureRmHDInsightClusterDefinition(d)
	if err != nil {
		return fmt.Errorf("Error expanding HDInsight cluster %q template: %+v", name, err)
	}

	securityProfile, err := expandAzureRmHDInsightSecurityProfile(d)
	if err != nil {
		return fmt.Errorf("Error expanding HDInsight cluster %q template: %+v", name, err)
	}

	computeProfile, err := expandAzureRmHDInsightComputeProfile(d)
	if err != nil {
		return fmt.Errorf("Error expanding HDInsight cluster %q template: %+v", name, err)
	}

	storageProfile, err := expandAzureRmHDInsightStorageProfile(d)
	if err != nil {
		return fmt.Errorf("Error expanding HDInsight cluster %q template: %+v", name, err)
	}

	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	clusterCreateParametersExtended := hdinsight.ClusterCreateParametersExtended{
		Location: utils.String(location),
		Tags:     expandedTags,
		Properties: &hdinsight.ClusterCreateProperties{
			ClusterVersion:    utils.String(clusterVersion),
			OsType:            osType,
			Tier:              tier,
			ClusterDefinition: clusterDefinition,
			SecurityProfile:   securityProfile,
			ComputeProfile:    computeProfile,
			StorageProfile:    storageProfile,
		},
	}

	ctx := client.StopContext
	future, err := hdInsightClustersClient.Create(ctx, resGroup, name, clusterCreateParametersExtended)
	if err != nil {
		return fmt.Errorf("Error Creating HDInsight cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletion(ctx, hdInsightClustersClient.Client)
	if err != nil {
		return fmt.Errorf("Error Creating HDInsight cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := hdInsightClustersClient.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error Retrieving HDInsight cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read HDInsight cluster %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmHDInsightClusterRead(d, meta)
}

func resourceArmHDInsightClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	hdinsightClustersClient := meta.(*ArmClient).hdInsightClustersClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["clusters"]

	ctx := client.StopContext
	resp, err := hdinsightClustersClient.Get(ctx, resGroup, name)
	if err != nil && resp.Response.StatusCode == http.StatusOK {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on HDInsight Cluster %q (resource group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	if location := resp.Location; location != nil {
		d.Set("location", *location)
	}
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.Properties; props != nil {
		d.Set("cluster_version", props.ClusterVersion)
		d.Set("tier", props.Tier)

		clusterDefinition := flattenAzureRmHDinsightClusterDefinition(props.ClusterDefinition)
		if err := d.Set("cluster_definition", &clusterDefinition); err != nil {
			return fmt.Errorf("Error setting `cluster_definition`: %+v", err)
		}

		securityProfile := flattenAzureRmHDinsightSecurityProfile(props.SecurityProfile)
		if err := d.Set("security_profile", &securityProfile); err != nil {
			return fmt.Errorf("Error setting `security_profile`: %+v", err)
		}

		computeProfile := flattenAzureRmHDinsightComputeProfile(props.ComputeProfile)
		if err := d.Set("compute_profile", &computeProfile); err != nil {
			return fmt.Errorf("Error setting `compute_profile`: %+v", err)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmHDInsightClusterDelete(d *schema.ResourceData, meta interface{}) error {
	hdInsightClient := meta.(*ArmClient).hdInsightClustersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return errwrap.Wrapf("Error Parsing Azure Resource ID {{err}}", err)
	}
	resGroup := id.ResourceGroup
	name := id.Path["clusters"]

	future, err := hdInsightClient.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting HDInsight cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletion(ctx, hdInsightClient.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the deleting HDInsight cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func resourceArmHDInsightClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	hdInsightClient := meta.(*ArmClient).hdInsightClustersClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for AzureRM HDInsight update.")

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	tags := d.Get("tags").(map[string]interface{})

	parameters := hdinsight.ClusterPatchParameters{
		Tags: expandTags(tags),
	}

	_, err := hdInsightClient.Update(ctx, resourceGroup, name, parameters)
	if err != nil {
		return err
	}

	read, err := hdInsightClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read HDInsight cluster %q (resource group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmContainerRegistryRead(d, meta)
}

func expandAzureRmHDInsightClusterDefinition(d *schema.ResourceData) (*hdinsight.ClusterDefinition, error) {
	clusterDefinitionInterfaceList := d.Get("cluster_definition").([]interface{})
	if clusterDefinitionInterfaceList == nil {
		return nil, nil
	}

	var clusterDefinition *hdinsight.ClusterDefinition
	for _, clusterDefinitionInterface := range clusterDefinitionInterfaceList {
		clusterDefinitionFlat := clusterDefinitionInterface.(map[string]interface{})

		blueprint := clusterDefinitionFlat["blueprint"].(string)
		kind := clusterDefinitionFlat["kind"].(string)
		clusterDefinition = &hdinsight.ClusterDefinition{
			Blueprint: utils.String(blueprint),
			Kind:      utils.String(kind),
		}

		if clusterDefinitionFlat["configurations"] != nil {
			configurations := make(map[string]interface{})

			configurationsInterfaceList := clusterDefinitionFlat["configurations"].([]interface{})
			for _, configurationsInterface := range configurationsInterfaceList {
				configurationsFlat := configurationsInterface.(map[string]interface{})

				gatewayInterfaceList := configurationsFlat["gateway"].([]interface{})
				for _, gatewayInterface := range gatewayInterfaceList {
					gatewayFlat := gatewayInterface.(map[string]interface{})

					gateway := make(map[string]interface{})
					gateway["restAuthCredential.isEnabled"] = gatewayFlat["rest_auth_credential_is_enabled"].(bool)
					gateway["restAuthCredential.username"] = gatewayFlat["rest_auth_credential_username"].(string)
					gateway["restAuthCredential.password"] = gatewayFlat["rest_auth_credential_password"].(string)
					configurations["gateway"] = &gateway
				}

				rserverInterfaceList := configurationsFlat["rserver"].([]interface{})
				for _, rserverInterface := range rserverInterfaceList {
					rserverFlat := rserverInterface.(map[string]interface{})

					rserver := make(map[string]interface{})
					rserver["rstudio"] = rserverFlat["rstudio"].(bool)
					configurations["rstudio"] = &rserver
				}
				clusterDefinition.Configurations = &configurations
			}
		}
	}
	return clusterDefinition, nil
}

func expandAzureRmHDInsightSecurityProfile(d *schema.ResourceData) (*hdinsight.SecurityProfile, error) {
	securityProfileInterfaceList := d.Get("security_profile").([]interface{})
	if securityProfileInterfaceList == nil {
		return nil, nil
	}

	var securityProfile *hdinsight.SecurityProfile
	for _, securityProfileInterface := range securityProfileInterfaceList {
		securityProfileFlat := securityProfileInterface.(map[string]interface{})

		securityProfile.DirectoryType = securityProfileFlat["directory_type"].(hdinsight.DirectoryType)

		domain := securityProfileFlat["domain"].(string)
		organizationalUnitDN := securityProfileFlat["organizational_unit_dn"].(string)
		domainUsername := securityProfileFlat["domain_username"].(string)
		domainUserPassword := securityProfileFlat["domain_password"].(string)

		securityProfile = &hdinsight.SecurityProfile{
			Domain:               &domain,
			OrganizationalUnitDN: &organizationalUnitDN,
			DomainUsername:       &domainUsername,
			DomainUserPassword:   &domainUserPassword,
		}

		ldapsURLs := []string{}
		ldapsURLsList := securityProfileFlat["ldaps_urls"].([]interface{})
		for _, ldapsURL := range ldapsURLsList {
			if v := ldapsURL.(string); v != "" {
				ldapsURLs = append(ldapsURLs, v)
			}
		}
		securityProfile.LdapsUrls = &ldapsURLs

		clusterUsersGroupDNS := []string{}
		clusterUsersGroupDNSList := securityProfileFlat["cluster_users_group_dns"].([]interface{})
		for _, clustUsersGroupDNSItem := range clusterUsersGroupDNSList {
			if v := clustUsersGroupDNSItem.(string); v != "" {
				clusterUsersGroupDNS = append(clusterUsersGroupDNS, v)
			}
		}
		securityProfile.ClusterUsersGroupDNS = &clusterUsersGroupDNS
	}
	return securityProfile, nil
}

func expandAzureRmHDInsightStorageProfile(d *schema.ResourceData) (*hdinsight.StorageProfile, error) {
	storageProfileInterfaceList := d.Get("storage_profile").([]interface{})
	if storageProfileInterfaceList == nil {
		return nil, nil
	}

	var storageProfile *hdinsight.StorageProfile
	for _, storageProfileInterface := range storageProfileInterfaceList {
		storageProfileFlat := storageProfileInterface.(map[string]interface{})

		storageAccounts := []hdinsight.StorageAccount{}
		storageAccountsInterfaceList := storageProfileFlat["storage_accounts"].([]interface{})

		for _, storageAccountInterface := range storageAccountsInterfaceList {
			storageAccountFlat := storageAccountInterface.(map[string]interface{})
			var storageAccount hdinsight.StorageAccount

			endpointWithPrefix := storageAccountFlat["name"].(string)
			endpoint := strings.Replace(strings.Replace(endpointWithPrefix, "https://", "", 1), "/", "", 1)
			isDefault := storageAccountFlat["is_default"].(bool)
			container := storageAccountFlat["container"].(string)
			key := storageAccountFlat["key"].(string)

			storageAccount = hdinsight.StorageAccount{
				Name:      &endpoint,
				IsDefault: &isDefault,
				Container: &container,
				Key:       &key,
			}

			storageAccounts = append(storageAccounts, storageAccount)

		}
		storageProfile = &hdinsight.StorageProfile{
			Storageaccounts: &storageAccounts,
		}
	}
	return storageProfile, nil
}

func expandAzureRmHDInsightComputeProfile(d *schema.ResourceData) (*hdinsight.ComputeProfile, error) {
	computeProfileInterfaceList := d.Get("compute_profile").([]interface{})
	if computeProfileInterfaceList == nil {
		return nil, nil
	}

	var computeProfile *hdinsight.ComputeProfile
	for _, computeInterface := range computeProfileInterfaceList {
		computeProfileFlat := computeInterface.(map[string]interface{})

		roles := []hdinsight.Role{}
		rolesInterfaceList := computeProfileFlat["roles"].([]interface{})

		for _, roleInterface := range rolesInterfaceList {
			roleFlat := roleInterface.(map[string]interface{})

			var role hdinsight.Role

			name := roleFlat["name"].(string)
			minInstanceCount := roleFlat["min_instance_count"].(int)
			targetInstanceCount := roleFlat["target_instance_count"].(int)

			role.Name = &name
			if minInstanceCount != 0 {
				mic := int32(minInstanceCount)
				role.MinInstanceCount = &mic
			}
			if targetInstanceCount != 0 {
				tic := int32(targetInstanceCount)
				role.TargetInstanceCount = &tic
			}

			hardwareProfileInterfaceList := roleFlat["hardware_profile"].([]interface{})
			for _, hardwareProfileInterface := range hardwareProfileInterfaceList {
				hardwareProfileFlat := hardwareProfileInterface.(map[string]interface{})

				hardwareProfile := hardwareProfileFlat["vm_size"].(string)
				role.HardwareProfile = &hdinsight.HardwareProfile{
					VMSize: &hardwareProfile,
				}
			}

			osProfileInterfaceList := roleFlat["os_profile"].([]interface{})
			for _, osProfileInterface := range osProfileInterfaceList {
				osProfileFlat := osProfileInterface.(map[string]interface{})

				var linuxOperatingSystemProfile *hdinsight.LinuxOperatingSystemProfile
				linuxOperatingSystemProfileInterfaceList := osProfileFlat["linux_operating_system_profile"].([]interface{})

				for _, linuxOperatingSystemProfileInterface := range linuxOperatingSystemProfileInterfaceList {
					linuxOperatingSystemProfileFlat := linuxOperatingSystemProfileInterface.(map[string]interface{})

					username := linuxOperatingSystemProfileFlat["username"].(string)
					password := linuxOperatingSystemProfileFlat["password"].(string)

					linuxOperatingSystemProfile = &hdinsight.LinuxOperatingSystemProfile{
						Username: &username,
						Password: &password,
					}

					var sshPublicKeys []hdinsight.SSHPublicKey
					sshProfileInterfaceList := linuxOperatingSystemProfileFlat["ssh_key"].([]interface{})

					for _, sshProfileInterface := range sshProfileInterfaceList {
						sshProfileFlat := sshProfileInterface.(map[string]interface{})

						if v := sshProfileFlat["key_data"].(string); v != "" {
							key := &hdinsight.SSHPublicKey{
								CertificateData: &v,
							}
							sshPublicKeys = append(sshPublicKeys, *key)
						}
					}

					sshProfile := &hdinsight.SSHProfile{
						PublicKeys: &sshPublicKeys,
					}
					linuxOperatingSystemProfile.SSHProfile = sshProfile
				}
				osProfile := &hdinsight.OsProfile{
					LinuxOperatingSystemProfile: linuxOperatingSystemProfile,
				}
				role.OsProfile = osProfile
			}

			virtualNetworkProfileInterfaceList := roleFlat["virtual_network_profile"].([]interface{})

			for _, virtualNetworkProfileInterface := range virtualNetworkProfileInterfaceList {
				virtualNetworkProfileFlat := virtualNetworkProfileInterface.(map[string]interface{})

				virtualNetworkProfile := &hdinsight.VirtualNetworkProfile{}
				vnetID := virtualNetworkProfileFlat["virtual_network_id"].(string)
				subnet := virtualNetworkProfileFlat["subnet_name"].(string)

				virtualNetworkProfile.ID = &vnetID
				virtualNetworkProfile.Subnet = &subnet

				role.VirtualNetworkProfile = virtualNetworkProfile
			}

			var dataDisksGroups []hdinsight.DataDisksGroups
			dataDisksGroupsInterfaceList := roleFlat["data_disks_group"].([]interface{})

			for _, dataDisksGroupsInterface := range dataDisksGroupsInterfaceList {
				var dataDisksGroupsItem hdinsight.DataDisksGroups

				dataDisksGroupsFlat := dataDisksGroupsInterface.(map[string]interface{})
				if v := dataDisksGroupsFlat["disks_per_node"].(int); v != 0 {
					dpn := int32(v)
					dataDisksGroupsItem.DisksPerNode = &dpn
				}
				if v := dataDisksGroupsFlat["storage_account_type"].(string); v != "" {
					dataDisksGroupsItem.StorageAccountType = &v
				}
				if v := dataDisksGroupsFlat["data_size_gb"].(int); v != 0 {
					dsgb := int32(v)
					dataDisksGroupsItem.DiskSizeGB = &dsgb
				}
				dataDisksGroups = append(dataDisksGroups, dataDisksGroupsItem)
			}
			role.DataDisksGroups = &dataDisksGroups
			roles = append(roles, role)
		}
		computeProfile = &hdinsight.ComputeProfile{
			Roles: &roles,
		}
	}
	return computeProfile, nil
}

func flattenAzureRmHDinsightClusterDefinition(clusterDefinition *hdinsight.ClusterDefinition) []interface{} {
	if clusterDefinition == nil {
		return nil
	}

	clusterDefinitionList := make([]interface{}, 0)
	clusterDefinitionFlat := make(map[string]interface{})
	if blueprint := clusterDefinition.Blueprint; blueprint != nil {
		clusterDefinitionFlat["blueprint"] = *blueprint
	}
	if kind := clusterDefinition.Kind; kind != nil {
		clusterDefinitionFlat["kind"] = *kind
	}
	if clusterDefinition.ComponentVersion != nil {
		componentVersion := *clusterDefinition.ComponentVersion
		componentVersionFlat := make(map[string]string)
		for k, v := range componentVersion {
			componentVersionFlat[k] = *v
		}

		clusterDefinitionFlat["component_version"] = componentVersionFlat
	}
	if clusterDefinition.Configurations != nil {
		configurations := *clusterDefinition.Configurations
		configurationsFlat := make(map[string]interface{})

		gatewayFlat := make(map[string]interface{})
		gateway := configurations["gateway"].(map[string]interface{})

		gatewayFlat["rest_auth_credentials_is_enabled"] = gateway["restAuthCredentials.isEnabled"].(string)
		gatewayFlat["rest_auth_credentials_username"] = gateway["restAuthCredentials.username"].(string)
		gatewayFlat["rest_auth_credentials_password"] = gateway["restAuthCredentials.password"].(string)
		configurationsFlat["gateway"] = &gatewayFlat

		rserverFlat := make(map[string]interface{})
		rserver := configurations["rserver"].(map[string]interface{})

		rserverFlat["rstudio"] = rserver["rstudio"].(string)
		configurationsFlat["rserver"] = &rserverFlat

		//TODO: support other configuration values i.e. core-site?
	}

	clusterDefinitionList = append(clusterDefinitionList, clusterDefinitionFlat)
	return clusterDefinitionList
}

func flattenAzureRmHDinsightSecurityProfile(securityProfile *hdinsight.SecurityProfile) []interface{} {
	if securityProfile == nil {
		return nil
	}

	securityProfileList := make([]interface{}, 0)
	securityProfileFlat := make(map[string]interface{})
	if directoryType := securityProfile.DirectoryType; directoryType != "" {
		securityProfileFlat["directory_type"] = directoryType
	}
	if domain := securityProfile.Domain; domain != nil {
		securityProfileFlat["domain"] = domain
	}
	if organizationalUnitDN := securityProfile.OrganizationalUnitDN; organizationalUnitDN != nil {
		securityProfileFlat["organizational_unit_dn"] = *organizationalUnitDN
	}
	if ldapsUrls := securityProfile.LdapsUrls; ldapsUrls != nil {
		ldapsUrlsList := make([]string, 0, len(*ldapsUrls))
		for _, ldapsUrl := range *ldapsUrls {
			ldapsUrlsList = append(ldapsUrlsList, ldapsUrl)
		}
		securityProfileFlat["ldaps_urls"] = ldapsUrlsList
	}
	if domainUsername := securityProfile.DomainUsername; domainUsername != nil {
		securityProfileFlat["domain_username"] = *domainUsername
	}
	if domainPassword := securityProfile.DomainUserPassword; domainPassword != nil {
		securityProfileFlat["domain_password"] = *domainPassword
	}
	if clusterUsersGroupDNS := securityProfile.ClusterUsersGroupDNS; clusterUsersGroupDNS != nil {
		clusterUsersGroupDNSList := make([]string, 0, len(*clusterUsersGroupDNS))
		for _, clusterUsersGroupDNSItem := range *clusterUsersGroupDNS {
			clusterUsersGroupDNSList = append(clusterUsersGroupDNSList, clusterUsersGroupDNSItem)
		}
		securityProfileFlat["cluster_users_group_dns"] = clusterUsersGroupDNSList
	}

	securityProfileList = append(securityProfileList, securityProfileFlat)
	return securityProfileList
}

func flattenAzureRmHDinsightComputeProfile(computeProfile *hdinsight.ComputeProfile) []interface{} {
	if computeProfile == nil {
		return nil
	}

	computeProfileList := make([]interface{}, 0)
	if roles := computeProfile.Roles; roles != nil {
		rolesList := make([]interface{}, 0, len(*computeProfile.Roles))
		for _, role := range *roles {
			roleFlat := make(map[string]interface{})
			if name := role.Name; name != nil {
				roleFlat["name"] = *name
			}
			if minInstanceCount := role.MinInstanceCount; minInstanceCount != nil {
				roleFlat["min_instance_count"] = *minInstanceCount
			}
			if targetInstanceCount := role.TargetInstanceCount; targetInstanceCount != nil {
				roleFlat["target_instance_count"] = *targetInstanceCount
			}
			if hardwareProfile := role.HardwareProfile; hardwareProfile != nil {
				hardwareProfileList := make([]interface{}, 0)
				hardwareProfileFlat := make(map[string]interface{})
				if vmSize := hardwareProfile.VMSize; vmSize != nil {
					hardwareProfileFlat["vm_size"] = *vmSize
				}
				hardwareProfileList = append(hardwareProfileList, hardwareProfileFlat)
				roleFlat["hardware_profile"] = hardwareProfileList
			}
			if osProfile := role.OsProfile; osProfile != nil {
				osProfileList := make([]interface{}, 0)
				osProfileFlat := make(map[string]interface{})
				if osProfile.LinuxOperatingSystemProfile != nil {
					linuxOperatingSystemProfileList := make([]interface{}, 0)
					linuxOperatingSystemProfileFlat := make(map[string]interface{})

					if username := osProfile.LinuxOperatingSystemProfile.Username; username != nil {
						linuxOperatingSystemProfileFlat["username"] = *username
					}
					if password := osProfile.LinuxOperatingSystemProfile.Password; password != nil {
						linuxOperatingSystemProfileFlat["password"] = *password
					}
					if sshProfile := osProfile.LinuxOperatingSystemProfile.SSHProfile; sshProfile != nil {
						sshKeysList := make([]interface{}, 0, len(*sshProfile.PublicKeys))
						for _, ssh := range *sshProfile.PublicKeys {
							key := make(map[string]interface{})
							key["key_data"] = *ssh.CertificateData
							sshKeysList = append(sshKeysList, key)
						}
						linuxOperatingSystemProfileFlat["ssh_key"] = sshKeysList
					}
					linuxOperatingSystemProfileList = append(linuxOperatingSystemProfileList, linuxOperatingSystemProfileFlat)
					osProfileFlat["linux_operating_system_profile"] = linuxOperatingSystemProfileList
				}
				osProfileList = append(osProfileList, osProfileFlat)
				roleFlat["os_profile"] = osProfileList
			}
			if virtualNetworkProfile := role.VirtualNetworkProfile; virtualNetworkProfile != nil {
				virtualNetworkList := make([]interface{}, 0)
				virtualNetworkFlat := make(map[string]interface{})

				if virtualNetworkID := role.VirtualNetworkProfile.ID; virtualNetworkID != nil {
					virtualNetworkFlat["virtual_network_id"] = *virtualNetworkID
				}
				if virtualNetworkSubnet := role.VirtualNetworkProfile.Subnet; virtualNetworkSubnet != nil {
					virtualNetworkFlat["subnet_name"] = *virtualNetworkSubnet
				}
				virtualNetworkList = append(virtualNetworkList, virtualNetworkFlat)
				roleFlat["virtual_network_profile"] = virtualNetworkList
			}
			if dataDisksGroups := role.DataDisksGroups; dataDisksGroups != nil {
				dataDiskGroupList := make([]interface{}, 0, len(*dataDisksGroups))
				for _, dataDiskGroup := range *dataDisksGroups {
					dataDiskGroupFlat := make(map[string]interface{})
					if disksPerNode := dataDiskGroup.DisksPerNode; disksPerNode != nil {
						dataDiskGroupFlat["disks_per_node"] = *disksPerNode
					}
					if storageAccountType := dataDiskGroup.StorageAccountType; storageAccountType != nil {
						dataDiskGroupFlat["storage_account_type"] = *storageAccountType
					}
					if diskSizeGB := dataDiskGroup.DiskSizeGB; diskSizeGB != nil {
						dataDiskGroupFlat["disk_size_gb"] = *diskSizeGB
					}
					dataDiskGroupList = append(dataDiskGroupList, dataDiskGroupFlat)
				}
				roleFlat["data_disks_group"] = dataDiskGroupList
			}
			rolesList = append(rolesList, roleFlat)
		}
		rolesFlat := make(map[string][]interface{})
		rolesFlat["roles"] = rolesList
		computeProfileList = append(computeProfileList, rolesFlat)
	}
	return computeProfileList
}

func flattenAzureRmHDinsightStorageProfile(storageProfile *hdinsight.StorageProfile) []interface{} {
	if storageProfile == nil {
		return nil
	}

	storageProfileList := make([]interface{}, 0)
	if storageAccounts := *storageProfile.Storageaccounts; storageAccounts != nil {
		storageAccountsList := make([]interface{}, 0, len(storageAccounts))
		for _, storageAccount := range storageAccounts {
			storageAccountFlat := make(map[string]interface{})
			if name := storageAccount.Name; name != nil {
				storageAccountFlat["name"] = *name
			}
			if isDefault := storageAccount.IsDefault; isDefault != nil {
				storageAccountFlat["is_default"] = *isDefault
			}
			if container := storageAccount.Container; container != nil {
				storageAccountFlat["container"] = *container
			}
			if key := storageAccount.Key; key != nil {
				storageAccountFlat["key"] = *key
			}
			storageAccountsList = append(storageAccountsList, storageAccountFlat)
		}
		storageProfileFlat := make(map[string]interface{})
		storageProfileFlat["storage_accounts"] = storageAccountsList
		storageProfileList = append(storageProfileList, storageProfileFlat)
	}
	return storageProfileList
}

func azureRmSuppressClusterVersionDiff(k, old, new string, d *schema.ResourceData) bool {
	segs := strings.Split(new, ".")
	prefix := strings.Join(segs[0:2], ".")
	return prefix == old
}

func azureRmSuppressSensitiveDiff(k, old, new string, d *schema.ResourceData) bool {
	if new == "" && old != "" {
		return true
	}
	return false
}

func azureRmSuppressRolesDiff(k, old, new string, d *schema.ResourceData) bool {
	if len(old) > 0 {
		if len(new) > len(old) {
			return true
		}
	}
	return false
}
