package azurerm

import (
	"fmt"
	"log"

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
				Type:     schema.TypeString,
				Optional: true,
				Default:  "~3.6",
				ValidateFunc: validation.StringInSlice([]string{
					"~3.6",
					"3.5",
					"3.4",
					"3.3",
					"3.2",
					"3.1",
					"3.0",
				}, false),
			},

			//TODO: Add support for Windows
			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Linux",
				ValidateFunc: validation.StringInSlice([]string{
					"Linux",
				}, false),
			},

			"tier": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Standard",
				ValidateFunc: validation.StringInSlice([]string{
					"Standard",
					"Premium",
				}, false),
			},

			//TODO: add support for configurations
			"cluster_definition": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blueprint": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Required: true,
						},
						"component_version": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{},
							},
						},
						"configurations": {
							Type:     schema.TypeMap,
							Optional: true,
						},
					},
				},
			},

			"security_profile": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"roles": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"min_instance_count": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"target_instance_count": {
										Type:     schema.TypeInt,
										Required: true,
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
												"linux_operation_system_profile": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"username": {
																Type:     schema.TypeString,
																Required: true,
															},
															"password": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"ssh_keys": {
																Type:     schema.TypeList,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"key_data": {
																			Type:     schema.TypeString,
																			Required: true,
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
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"subnet": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"data_disks_group": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"disks_per_node": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"storage_account_type": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"data_size_gb": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
									"script_actions": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"uri": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"parameters": {
													Type:     schema.TypeString,
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

			"storage_profile": {
				Type:     schema.TypeList,
				Required: true,
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
	client := meta.(*ArmClient).hdInsightClustersClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM HDInsight creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resGroup := d.Get("resource_group_name").(string)
	clusterVersion := d.Get("cluster_version").(string)
	osType := d.Get("os_type").(hdinsight.OSType)
	tier := d.Get("tier").(hdinsight.Tier)

	clusterDefinition := expandAzureRmHDInsightClusterDefinition(d)
	securityProfile := expandAzureRmHDInsightSecurityProfile(d)
	computeProfile := expandAzureRmHDInsightComputeProfile(d)
	storageProfile := expandAzureRmHDInsightStorageProfile(d)

	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	clusterCreateProperties := &hdinsight.ClusterCreateProperties{
		ClusterVersion:    utils.String(clusterVersion),
		OsType:            osType,
		Tier:              tier,
		ClusterDefinition: clusterDefinition,
		SecurityProfile:   securityProfile,
		ComputeProfile:    computeProfile,
		StorageProfile:    storageProfile,
	}
	clusterCreateParametersExtended := hdinsight.ClusterCreateParametersExtended{
		Location:   utils.String(location),
		Tags:       expandedTags,
		Properties: clusterCreateProperties,
	}

	future, err := client.Create(ctx, resGroup, name, clusterCreateParametersExtended)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating HDInsight cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating HDInsight cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
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
	name := id.Path["hdinsight"]

	ctx := client.StopContext
	resp, err := hdinsightClustersClient.Get(ctx, resGroup, name)
	if err != nil {
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
	d.Set("cluster_version", resp.Properties.ClusterVersion)
	d.Set("os_type", resp.Properties.OsType)
	d.Set("tier", resp.Properties.Tier)

	clusterDefinition := flattenAzureRmHDinsightClusterDefinition(resp.Properties.ClusterDefinition)
	if err := d.Set("cluster_definition", &clusterDefinition); err != nil {
		return fmt.Errorf("Error setting `cluster_definition`: %+v", err)
	}

	securityProfile := flattenAzureRmHDinsightSecurityProfile(resp.Properties.SecurityProfile)
	if err := d.Set("security_profile", &securityProfile); err != nil {
		return fmt.Errorf("Error setting `security_profile`: %+v", err)
	}

	computeProfile := flattenAzureRmHDinsightComputeProfile(resp.Properties.ComputeProfile)
	if err := d.Set("compute_profile", &computeProfile); err != nil {
		return fmt.Errorf("Error setting `compute_profile`: %+v", err)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmHDInsightClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).loadBalancerClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return errwrap.Wrapf("Error Parsing Azure Resource ID {{err}}", err)
	}
	resGroup := id.ResourceGroup
	name := id.Path["hdinsight"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting HDInsight cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the deleting HDInsight cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func resourceArmHDInsightClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).hdInsightClustersClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for AzureRM HDInsight update.")

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	tags := d.Get("tags").(map[string]interface{})

	parameters := hdinsight.ClusterPatchParameters{
		Tags: expandTags(tags),
	}

	//TODO: Does not return future - do we need to wait for the update?
	_, err := client.Update(ctx, resourceGroup, name, parameters)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read HDInsight cluster %q (resource group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmContainerRegistryRead(d, meta)
}

func flattenList(interfaceList []interface{}) (map[string]interface{}, error) {
	listLength := len(interfaceList)
	if listLength != 1 {
		return nil, fmt.Errorf("All lists should contain 1 item")
	}
	only := interfaceList[0]
	flatmap, ok := only.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Cannot cast list to flatmap")
	}
	return flatmap, nil
}

func expandAzureRmHDInsightClusterDefinition(d *schema.ResourceData) *hdinsight.ClusterDefinition {
	clusterDefinitionInterfaceList := d.Get("cluster_definition").([]interface{})
	clusterDefinition := &hdinsight.ClusterDefinition{}

	clusterDefinitionFlatMap, err := flattenList(clusterDefinitionInterfaceList)
	if err != nil {
		panic(err)
	}

	if v := clusterDefinitionFlatMap["blueprint"].(string); v != "" {
		clusterDefinition.Blueprint = &v
	}
	if v := clusterDefinitionFlatMap["kind"].(string); v != "" {
		clusterDefinition.Kind = &v
	}

	componentVersion := make(map[string]*string)
	componentVersionsFlatMap := clusterDefinitionFlatMap["component_version"].(map[string]string)
	for k, v := range componentVersionsFlatMap {
		componentVersion[k] = &v
	}
	clusterDefinition.ComponentVersion = &componentVersion

	return clusterDefinition
}

func expandAzureRmHDInsightSecurityProfile(d *schema.ResourceData) *hdinsight.SecurityProfile {
	securityProfileInterfaceList := d.Get("security_profile").([]interface{})
	securityProfile := &hdinsight.SecurityProfile{}

	securityProfileFlatMap, err := flattenList(securityProfileInterfaceList)
	if err != nil {
		panic(err)
	}
	if v := securityProfileFlatMap["directory_type"].(hdinsight.DirectoryType); v != "" {
		securityProfile.DirectoryType = v
	}

	if v := securityProfileFlatMap["domain"].(string); v != "" {
		securityProfile.Domain = &v
	}

	if v := securityProfileFlatMap["organizational_unit_dn"].(string); v != "" {
		securityProfile.OrganizationalUnitDN = &v
	}

	ldapsUrls := []string{}
	ldapsUrlsList := securityProfileFlatMap["ldaps_urls"].([]interface{})
	for _, ldapsUrl := range ldapsUrlsList {
		if v := ldapsUrl.(string); v != "" {
			ldapsUrls = append(ldapsUrls, v)
		}
	}
	securityProfile.LdapsUrls = &ldapsUrls

	if v := securityProfileFlatMap["domain_username"].(string); v != "" {
		securityProfile.DomainUsername = &v
	}

	if v := securityProfileFlatMap["domain_password"].(string); v != "" {
		securityProfile.DomainUserPassword = &v
	}

	clusterUsersGroupDNS := []string{}
	clusterUsersGroupDNSList := securityProfileFlatMap["cluster_users_group_dns"].([]interface{})
	for _, clustUsersGroupDNSItem := range clusterUsersGroupDNSList {
		if v := clustUsersGroupDNSItem.(string); v != "" {
			clusterUsersGroupDNS = append(clusterUsersGroupDNS, v)
		}
	}
	securityProfile.ClusterUsersGroupDNS = &clusterUsersGroupDNS

	return securityProfile
}

func expandAzureRmHDInsightStorageProfile(d *schema.ResourceData) *hdinsight.StorageProfile {
	storageProfileInterfaceList := d.Get("storage_profile").([]interface{})
	storageProfile := &hdinsight.StorageProfile{}

	storageProfileFlatMap, err := flattenList(storageProfileInterfaceList)
	if err != nil {
		panic(err)
	}

	storageAccounts := []hdinsight.StorageAccount{}
	storageAccountsInterfaceList := storageProfileFlatMap["storage_accounts"].([]interface{})

	for _, storageAccountInterface := range storageAccountsInterfaceList {
		storageAccountFlatMap := storageAccountInterface.(map[string]interface{})
		var storageAccount hdinsight.StorageAccount
		if v := storageAccountFlatMap["name"].(string); v != "" {
			storageAccount.Name = &v
		}
		if v := storageAccountFlatMap["is_default"].(bool); true {
			storageAccount.IsDefault = &v
		}
		if v := storageAccountFlatMap["container"].(string); v != "" {
			storageAccount.Container = &v
		}
		if v := storageAccountFlatMap["key"].(string); v != "" {
			storageAccount.Key = &v
		}
		storageAccounts = append(storageAccounts, storageAccount)
	}

	storageProfile.Storageaccounts = &storageAccounts
	return storageProfile
}

func expandAzureRmHDInsightComputeProfile(d *schema.ResourceData) *hdinsight.ComputeProfile {
	computeProfileInterfaceList := d.Get("compute_profile").([]interface{})
	computeProfile := &hdinsight.ComputeProfile{}

	computeProfileFlatMap, err := flattenList(computeProfileInterfaceList)
	if err != nil {
		panic(err)
	}

	roles := []hdinsight.Role{}
	rolesInterfaceList := computeProfileFlatMap["roles"].([]interface{})
	for _, roleInterface := range rolesInterfaceList {
		roleFlatMap := roleInterface.(map[string]interface{})
		var role hdinsight.Role
		if v := roleFlatMap["name"].(string); v != "" {
			role.Name = &v
		}
		if v := roleFlatMap["min_instance_count"].(int32); v != 0 {
			role.MinInstanceCount = &v
		}
		if v := roleFlatMap["target_instance_count"].(int32); v != 0 {
			role.TargetInstanceCount = &v
		}
		hardwareProfileInterfaceList := roleFlatMap["hardware_profile"].([]interface{})
		for _, hardwareProfileInterface := range hardwareProfileInterfaceList {
			hardwareProfileFlatMap := hardwareProfileInterface.(map[string]interface{})
			if v := hardwareProfileFlatMap["vm_size"].(string); v != "" {
				role.HardwareProfile = &hdinsight.HardwareProfile{
					VMSize: &v,
				}
			}
		}
		osProfileInterfaceList := roleFlatMap["os_profile"].([]interface{})
		for _, osProfileInterface := range osProfileInterfaceList {
			osProfileFlatMap := osProfileInterface.(map[string]interface{})
			linuxOperatingSystemProfile := &hdinsight.LinuxOperatingSystemProfile{}
			linuxOperatingSystemProfileInterfaceList := osProfileFlatMap["linux_operating_system_profile"].([]interface{})

			for _, linuxOperatingSystemProfileInterface := range linuxOperatingSystemProfileInterfaceList {

				linuxOperatingSystemProfileFlatMap := linuxOperatingSystemProfileInterface.(map[string]interface{})
				if v := linuxOperatingSystemProfileFlatMap["username"].(string); v != "" {
					linuxOperatingSystemProfile.Username = &v
				}
				if v := linuxOperatingSystemProfileFlatMap["password"].(string); v != "" {
					linuxOperatingSystemProfile.Password = &v
				}
				var sshPublicKeys []hdinsight.SSHPublicKey

				sshProfileInterfaceList := linuxOperatingSystemProfileFlatMap["ssh_key"].([]interface{})
				for _, sshProfileInterface := range sshProfileInterfaceList {
					sshProfileFlatMap := sshProfileInterface.(map[string]interface{})
					if v := sshProfileFlatMap["key_data"].(string); v != "" {
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

		virtualNetworkProfileInterfaceList := roleFlatMap["virtual_network_profile"].([]interface{})
		for _, virtualNetworkProfileInterface := range virtualNetworkProfileInterfaceList {
			virtualNetworkProfile := &hdinsight.VirtualNetworkProfile{}

			virtualNetworkProfileFlatMap := virtualNetworkProfileInterface.(map[string]interface{})
			if v := virtualNetworkProfileFlatMap["id"].(string); v != "" {
				virtualNetworkProfile.ID = &v
			}
			if v := virtualNetworkProfileFlatMap["subnet"].(string); v != "" {
				virtualNetworkProfile.Subnet = &v
			}
			role.VirtualNetworkProfile = virtualNetworkProfile
		}

		var dataDisksGroups []hdinsight.DataDisksGroups

		dataDisksGroupsInterfaceList := roleFlatMap["data_disk_group"].([]interface{})
		for _, dataDisksGroupsInterface := range dataDisksGroupsInterfaceList {
			dataDisksGroupsItem := &hdinsight.DataDisksGroups{}

			dataDisksGroupsFlatMap := dataDisksGroupsInterface.(map[string]interface{})
			if v := dataDisksGroupsFlatMap["disks_per_node"].(int32); v != 0 {
				dataDisksGroupsItem.DisksPerNode = &v
			}
			if v := dataDisksGroupsFlatMap["storage_account_type"].(string); v != "" {
				dataDisksGroupsItem.StorageAccountType = &v
			}
			if v := dataDisksGroupsFlatMap["data_size_gb"].(int32); v != 0 {
				dataDisksGroupsItem.DiskSizeGB = &v
			}
			dataDisksGroups = append(dataDisksGroups, *dataDisksGroupsItem)
		}
		role.DataDisksGroups = &dataDisksGroups

		var scriptActions []hdinsight.ScriptAction

		scriptActionsInterfaceList := roleFlatMap["script_actions"].([]interface{})
		for _, scriptActionInterface := range scriptActionsInterfaceList {
			var scriptAction hdinsight.ScriptAction

			scriptActionFlatMap := scriptActionInterface.(map[string]interface{})
			if v := scriptActionFlatMap["name"].(string); v != "" {
				scriptAction.Name = &v
			}
			if v := scriptActionFlatMap["uri"].(string); v != "" {
				scriptAction.URI = &v
			}
			if v := scriptActionFlatMap["parameters"].(string); v != "" {
				scriptAction.Parameters = &v
			}
			scriptActions = append(scriptActions, scriptAction)
		}
		role.ScriptActions = &scriptActions

		roles = append(roles, role)
	}
	computeProfile.Roles = &roles

	return computeProfile
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
	if componentVersion := clusterDefinition.ComponentVersion; componentVersion != nil {
		componentVersionFlat := make(map[string]*string)
		for k, v := range *componentVersion {
			componentVersionFlat[k] = v
		}
		clusterDefinitionFlat["component_version"] = componentVersionFlat
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
				linuxOperatingSystemProfileList := make([]interface{}, 0)
				if linuxOperatinSystemProfile := osProfile.LinuxOperatingSystemProfile; linuxOperatinSystemProfile != nil {
					linuxOperationSystemProfile := make(map[string]interface{})
					if username := osProfile.LinuxOperatingSystemProfile.Username; username != nil {
						linuxOperationSystemProfile["username"] = *username
					}
					if password := osProfile.LinuxOperatingSystemProfile.Password; password != nil {
						linuxOperationSystemProfile["password"] = *password
					}
					if sshProfile := osProfile.LinuxOperatingSystemProfile.SSHProfile; sshProfile != nil {
						sshKeysList := make([]interface{}, 0, len(*sshProfile.PublicKeys))
						for _, ssh := range *sshProfile.PublicKeys {
							key := make(map[string]interface{})
							key["key_data"] = *ssh.CertificateData
							sshKeysList = append(sshKeysList, key)
						}
						linuxOperationSystemProfile["ssh_key"] = sshKeysList
					}
					linuxOperatingSystemProfileList = append(linuxOperatingSystemProfileList, linuxOperationSystemProfile)
				}
				osProfileFlat["linux_operation_system_profile"] = linuxOperatingSystemProfileList
				osProfileList = append(osProfileList, osProfileFlat)
				roleFlat["os_profile"] = osProfileList
			}
			if virtualNetworkProfile := role.VirtualNetworkProfile; virtualNetworkProfile != nil {
				virtualNetworkList := make([]interface{}, 0)
				virtualNetworkFlat := make(map[string]interface{})
				if virtualNetworkID := role.VirtualNetworkProfile.ID; virtualNetworkID != nil {
					virtualNetworkFlat["id"] = *virtualNetworkID
				}
				if virtualNetworkSubnet := role.VirtualNetworkProfile.Subnet; virtualNetworkSubnet != nil {
					virtualNetworkFlat["subnet"] = *virtualNetworkSubnet
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
			if scriptActions := role.ScriptActions; scriptActions != nil {
				scriptActionList := make([]interface{}, 0, len(*scriptActions))
				for _, s := range *scriptActions {
					scriptActionsFlat := make(map[string]interface{})
					if name := s.Name; name != nil {
						scriptActionsFlat["name"] = *name
					}
					if URI := s.URI; URI != nil {
						scriptActionsFlat["uri"] = *URI
					}
					if parameters := s.Parameters; parameters != nil {
						scriptActionsFlat["parameters"] = *parameters
					}
					scriptActionList = append(scriptActionList, scriptActionsFlat)
				}
				roleFlat["script_actions"] = scriptActionList
			}
			rolesList = append(rolesList, roleFlat)
		}
		rolesFlat := make(map[string]interface{})
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
