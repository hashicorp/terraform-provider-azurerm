package azurerm

// Code based on the terraform.provider plugin by Microsoft (R) AutoRest Code Generator.

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/hdinsight/mgmt/2015-03-01-preview/hdinsight"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmHDInsightClusters() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmHDInsightClustersCreate,
		Read:   resourceArmHDInsightClustersRead,
		Update: resourceArmHDInsightClustersUpdate,
		Delete: resourceArmHDInsightClustersDelete,
		Schema: map[string]*schema.Schema{
			"location":            locationSchema(),
			"resource_group_name": resourceGroupNameSchema(),
			"name": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"cluster_type": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"hadoop",
					"hbase",
					"storm",
					"spark",
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},
			"cluster_version": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"component_version": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"restauth_username": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"restauth_password": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"storage_account": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blob_endpoint": {
							Required: true,
							ForceNew: true,
							Type:     schema.TypeString,
						},
						"container": {
							Required: true,
							ForceNew: true,
							Type:     schema.TypeString,
						},
						"access_key": {
							Required: true,
							ForceNew: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"head_node": {
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Type:     schema.TypeList,
				Elem:     hdInsightClustersNodeSchema(),
			},
			"worker_node": {
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Type:     schema.TypeList,
				Elem:     hdInsightClustersNodeSchema(),
			},
			"zookeeper_node": {
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Type:     schema.TypeList,
				Elem:     hdInsightClustersNodeSchema(),
			},
			"tags": tagsSchema(),
			"tier": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
				Default:  string(hdinsight.Standard),
				ValidateFunc: validation.StringInSlice([]string{
					string(hdinsight.Standard),
					string(hdinsight.Premium),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},
			"connectivity_endpoints": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"port": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"protocol": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func hdInsightClustersNodeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"linux_os_profile": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Optional: true,
							ForceNew: true,
							Type:     schema.TypeString,
						},
						"ssh_keys": {
							Required: true,
							ForceNew: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key_data": {
										Required: true,
										ForceNew: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"username": {
							Required: true,
							ForceNew: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"min_instance_count": {
				Optional:     true,
				Default:      1,
				ForceNew:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(1, 32),
			},
			"script_actions": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Required: true,
							ForceNew: true,
							Type:     schema.TypeString,
						},
						"parameters": {
							Required: true,
							ForceNew: true,
							Type:     schema.TypeString,
						},
						"uri": {
							Required: true,
							ForceNew: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"target_instance_count": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(1, 32),
			},
			"vm_size": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"vnet_profile": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Required: true,
							ForceNew: true,
							Type:     schema.TypeString,
						},
						"subnet": {
							Required: true,
							ForceNew: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func resourceArmHDInsightClustersCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).hdInsightClustersClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroupName := d.Get("resource_group_name").(string)
	clusterName := d.Get("name").(string)
	parameters := hdinsight.ClusterCreateParametersExtended{}
	parameters.Location = utils.String(azureRMNormalizeLocation(d.Get("location").(string)))
	tags := d.Get("tags").(map[string]interface{})
	tmpParamOfTags := expandTags(tags)
	parameters.Tags = tmpParamOfTags
	parameters.Properties = &hdinsight.ClusterCreateProperties{}
	if paramValue, paramExists := d.GetOk("cluster_version"); paramExists {
		parameters.Properties.ClusterVersion = utils.String(paramValue.(string))
	}
	parameters.Properties.OsType = hdinsight.Linux
	if paramValue, paramExists := d.GetOk("tier"); paramExists {
		parameters.Properties.Tier = hdinsight.Tier(paramValue.(string))
	}
	parameters.Properties.ClusterDefinition = &hdinsight.ClusterDefinition{}
	parameters.Properties.ClusterDefinition.Kind = utils.String(d.Get("cluster_type").(string))
	if paramValue, paramExists := d.GetOk("component_version"); paramExists {
		tmpParamOfComponentVersion := make(map[string]*string)
		for tmpParamKeyOfComponentVersion, tmpParamItemOfComponentVersion := range paramValue.(map[string]interface{}) {
			parametersPropertiesClusterDefinitionComponentVersion := utils.String(tmpParamItemOfComponentVersion.(string))
			tmpParamOfComponentVersion[tmpParamKeyOfComponentVersion] = parametersPropertiesClusterDefinitionComponentVersion
		}
		parameters.Properties.ClusterDefinition.ComponentVersion = tmpParamOfComponentVersion
	}
	tmpParamOfConfigurations := make(map[string]interface{})
	tmpParamOfGatewayConfigurations := make(map[string]interface{})
	tmpParamOfGatewayConfigurations["restAuthCredential.isEnabled"] = true
	tmpParamOfGatewayConfigurations["restAuthCredential.username"] = d.Get("restauth_username")
	tmpParamOfGatewayConfigurations["restAuthCredential.password"] = d.Get("restauth_password")
	tmpParamOfCoreSiteConfigurations := make(map[string]interface{})
	for _, paramValue := range d.Get("storage_account").([]interface{}) {
		tmpParamOfStorageAccount := paramValue.(map[string]interface{})
		tmpParamOfStorageAccountName := tmpParamOfStorageAccount["blob_endpoint"].(string)
		tmpParamOfStorageAccountName = tmpParamOfStorageAccountName[len("https://") : len(tmpParamOfStorageAccountName)-1]
		tmpParamOfCoreSiteConfigurations["fs.defaultFS"] = "wasb://" + tmpParamOfStorageAccount["container"].(string) + "@" + tmpParamOfStorageAccountName
		tmpParamOfCoreSiteConfigurations["fs.azure.account.key."+tmpParamOfStorageAccountName] = tmpParamOfStorageAccount["access_key"].(string)
	}
	tmpParamOfConfigurations["gateway"] = tmpParamOfGatewayConfigurations
	tmpParamOfConfigurations["core-site"] = tmpParamOfCoreSiteConfigurations
	parameters.Properties.ClusterDefinition.Configurations = &tmpParamOfConfigurations
	parameters.Properties.ComputeProfile = &hdinsight.ComputeProfile{}
	tmpParamOfRoles := make([]hdinsight.Role, 0)
	appendParamRole := func(name string, tmpParamItemOfRoles interface{}) {
		tmpParamValueOfRoles := tmpParamItemOfRoles.(map[string]interface{})
		parametersPropertiesComputeProfileRoles := &hdinsight.Role{}
		parametersPropertiesComputeProfileRoles.Name = utils.String(name)
		if paramValue, paramExists := tmpParamValueOfRoles["min_instance_count"]; paramExists {
			parametersPropertiesComputeProfileRoles.MinInstanceCount = utils.Int32(int32(paramValue.(int)))
		}
		if paramValue, paramExists := tmpParamValueOfRoles["target_instance_count"]; paramExists {
			parametersPropertiesComputeProfileRoles.TargetInstanceCount = utils.Int32(int32(paramValue.(int)))
		}
		parametersPropertiesComputeProfileRoles.HardwareProfile = &hdinsight.HardwareProfile{}
		if paramValue, paramExists := tmpParamValueOfRoles["vm_size"]; paramExists {
			parametersPropertiesComputeProfileRoles.HardwareProfile.VMSize = utils.String(paramValue.(string))
		}
		parametersPropertiesComputeProfileRoles.OsProfile = &hdinsight.OsProfile{}
		for _, paramValue := range tmpParamValueOfRoles["linux_os_profile"].([]interface{}) {
			parametersPropertiesComputeProfileRoles.OsProfile.LinuxOperatingSystemProfile = &hdinsight.LinuxOperatingSystemProfile{}
			tmpParamOfRoleslinuxOsProfile := paramValue.(map[string]interface{})
			if paramValue, paramExists := tmpParamOfRoleslinuxOsProfile["username"]; paramExists {
				parametersPropertiesComputeProfileRoles.OsProfile.LinuxOperatingSystemProfile.Username = utils.String(paramValue.(string))
			}
			if paramValue, paramExists := tmpParamOfRoleslinuxOsProfile["password"]; paramExists {
				parametersPropertiesComputeProfileRoles.OsProfile.LinuxOperatingSystemProfile.Password = utils.String(paramValue.(string))
			}
			if paramValue, paramExists := tmpParamOfRoleslinuxOsProfile["ssh_keys"]; paramExists {
				parametersPropertiesComputeProfileRoles.OsProfile.LinuxOperatingSystemProfile.SSHProfile = &hdinsight.SSHProfile{}
				tmpParamOfRoleslinuxOsProfilesshKeys := make([]hdinsight.SSHPublicKey, 0)
				for _, tmpParamItemOfRoleslinuxOsProfilesshKeys := range paramValue.([]interface{}) {
					tmpParamValueOfRoleslinuxOsProfilesshKeys := tmpParamItemOfRoleslinuxOsProfilesshKeys.(map[string]interface{})
					parametersPropertiesComputeProfileRolesOsProfileLinuxOperatingSystemProfileSSHProfilePublicKeys := &hdinsight.SSHPublicKey{}
					if paramValue, paramExists := tmpParamValueOfRoleslinuxOsProfilesshKeys["key_data"]; paramExists {
						parametersPropertiesComputeProfileRolesOsProfileLinuxOperatingSystemProfileSSHProfilePublicKeys.CertificateData = utils.String(paramValue.(string))
					}
					tmpParamOfRoleslinuxOsProfilesshKeys = append(tmpParamOfRoleslinuxOsProfilesshKeys, *parametersPropertiesComputeProfileRolesOsProfileLinuxOperatingSystemProfileSSHProfilePublicKeys)
				}
				parametersPropertiesComputeProfileRoles.OsProfile.LinuxOperatingSystemProfile.SSHProfile.PublicKeys = &tmpParamOfRoleslinuxOsProfilesshKeys
			}
		}
		for _, paramValue := range tmpParamValueOfRoles["vnet_profile"].([]interface{}) {
			parametersPropertiesComputeProfileRoles.VirtualNetworkProfile = &hdinsight.VirtualNetworkProfile{}
			tmpParamOfRolesvnetProfile := paramValue.(map[string]interface{})
			if paramValue, paramExists := tmpParamOfRolesvnetProfile["id"]; paramExists {
				parametersPropertiesComputeProfileRoles.VirtualNetworkProfile.ID = utils.String(paramValue.(string))
			}
			if paramValue, paramExists := tmpParamOfRolesvnetProfile["subnet"]; paramExists {
				parametersPropertiesComputeProfileRoles.VirtualNetworkProfile.Subnet = utils.String(paramValue.(string))
			}
		}
		if paramValue, paramExists := tmpParamValueOfRoles["script_actions"]; paramExists {
			tmpParamOfRolesscriptActions := make([]hdinsight.ScriptAction, 0)
			for _, tmpParamItemOfRolesscriptActions := range paramValue.([]interface{}) {
				tmpParamValueOfRolesscriptActions := tmpParamItemOfRolesscriptActions.(map[string]interface{})
				parametersPropertiesComputeProfileRolesScriptActions := &hdinsight.ScriptAction{}
				parametersPropertiesComputeProfileRolesScriptActions.Name = utils.String(tmpParamValueOfRolesscriptActions["name"].(string))
				parametersPropertiesComputeProfileRolesScriptActions.URI = utils.String(tmpParamValueOfRolesscriptActions["uri"].(string))
				parametersPropertiesComputeProfileRolesScriptActions.Parameters = utils.String(tmpParamValueOfRolesscriptActions["parameters"].(string))
				tmpParamOfRolesscriptActions = append(tmpParamOfRolesscriptActions, *parametersPropertiesComputeProfileRolesScriptActions)
			}
			parametersPropertiesComputeProfileRoles.ScriptActions = &tmpParamOfRolesscriptActions
		}
		tmpParamOfRoles = append(tmpParamOfRoles, *parametersPropertiesComputeProfileRoles)
	}
	for _, tmpParamItemOfRoles := range d.Get("head_node").([]interface{}) {
		appendParamRole("headnode", tmpParamItemOfRoles)
	}
	for _, tmpParamItemOfRoles := range d.Get("worker_node").([]interface{}) {
		appendParamRole("workernode", tmpParamItemOfRoles)
	}
	for _, tmpParamItemOfRoles := range d.Get("zookeeper_node").([]interface{}) {
		appendParamRole("zookeepernode", tmpParamItemOfRoles)
	}
	parameters.Properties.ComputeProfile.Roles = &tmpParamOfRoles

	future, err := client.Create(ctx, resourceGroupName, clusterName, parameters)
	if err != nil {
		return fmt.Errorf("HD Insight Clusters creation error: %+v", err)
	}
	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("HD Insight Clusters creation future wait for completion error: %+v", err)
	}
	response, err := future.Result(client)
	if err != nil {
		return fmt.Errorf("HD Insight Clusters creation future result error: %+v", err)
	}

	if response.ID == nil {
		return fmt.Errorf("Cannot get the ID of HD Insight Clusters %q (Resource Group %q) ID", clusterName, resourceGroupName)
	}
	d.SetId(*response.ID)

	return resourceArmHDInsightClustersRead(d, meta)
}

func resourceArmHDInsightClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).hdInsightClustersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroupName := id.ResourceGroup
	clusterName := id.Path["clusters"]

	response, err := client.Get(ctx, resourceGroupName, clusterName)
	if err != nil {
		return fmt.Errorf("HD Insight Clusters read error: %+v", err)
	}

	d.Set("name", *response.Name)
	d.Set("resource_group_name", resourceGroupName)
	d.Set("location", azureRMNormalizeLocation(*response.Location))
	if response.Properties != nil {
		clusterVersions := strings.Split(*response.Properties.ClusterVersion, ".")
		d.Set("cluster_version", strings.Join(clusterVersions[0:2], "."))
		d.Set("tier", response.Properties.Tier)
		if response.Properties.ClusterDefinition.Kind != nil {
			d.Set("cluster_type", *response.Properties.ClusterDefinition.Kind)
		}
		if response.Properties.ClusterDefinition.ComponentVersion != nil {
			tmpRespOfComponentVersion := make(map[string]interface{})
			for tmpRespKeyOfComponentVersion, tmpRespItemOfComponentVersion := range response.Properties.ClusterDefinition.ComponentVersion {
				tmpRespValueOfComponentVersion := *tmpRespItemOfComponentVersion
				tmpRespOfComponentVersion[tmpRespKeyOfComponentVersion] = tmpRespValueOfComponentVersion
			}
			d.Set("component_version", tmpRespOfComponentVersion)
		}
		if response.Properties.ClusterDefinition.Configurations != nil {
			tmpRespOfConfigurations := *response.Properties.ClusterDefinition.Configurations.(*map[string]interface{})
			if paramValue, paramExists := tmpRespOfConfigurations["gateway"]; paramExists {
				tmpRespOfGatewayConfigurations := paramValue.(map[string]interface{})
				d.Set("restauth_username", tmpRespOfGatewayConfigurations["restAuthCredential.username"].(string))
				d.Set("restauth_password", tmpRespOfGatewayConfigurations["restAuthCredential.password"].(string))
			}
			if paramValue, paramExists := tmpRespOfConfigurations["core-site"]; paramExists {
				tmpRespOfCoreSiteConfigurations := paramValue.(map[string]interface{})
				tmpRespOfStorageAccount := make(map[string]interface{})
				tmpRespOfStorageAccountURL := tmpRespOfCoreSiteConfigurations["fs.defaultFS"].(string)
				tmpRespOfStorageAccountName := tmpRespOfStorageAccountURL[strings.Index(tmpRespOfStorageAccountURL, "@"):]
				tmpRespOfStorageAccount["blob_endpoint"] = tmpRespOfStorageAccountName
				tmpRespOfStorageAccount["container"] = tmpRespOfStorageAccountURL[len("wasb://"):strings.Index(tmpRespOfStorageAccountURL, "@")]
				tmpRespOfStorageAccount["access_key"] = tmpRespOfCoreSiteConfigurations["fs.azure.account.key."+tmpRespOfStorageAccountName]
				d.Set("storage_account", tmpRespOfStorageAccount)
			}
		}
		if response.Properties.ComputeProfile != nil {
			if response.Properties.ComputeProfile.Roles != nil && len(*response.Properties.ComputeProfile.Roles) > 0 {
				for _, tmpRespItemOfRoles := range *response.Properties.ComputeProfile.Roles {
					tmpRespValueOfRoles := make(map[string]interface{})
					if tmpRespItemOfRoles.MinInstanceCount != nil {
						tmpRespValueOfRoles["min_instance_count"] = *tmpRespItemOfRoles.MinInstanceCount
					}
					if tmpRespItemOfRoles.TargetInstanceCount != nil {
						tmpRespValueOfRoles["target_instance_count"] = *tmpRespItemOfRoles.TargetInstanceCount
					}
					if tmpRespItemOfRoles.HardwareProfile != nil {
						if tmpRespItemOfRoles.HardwareProfile.VMSize != nil {
							tmpRespValueOfRoles["vm_size"] = *tmpRespItemOfRoles.HardwareProfile.VMSize
						}
					}
					if tmpRespItemOfRoles.OsProfile != nil {
						if tmpRespItemOfRoles.OsProfile.LinuxOperatingSystemProfile != nil {
							tmpRespOfRoleslinuxOsProfile := make(map[string]interface{})
							if tmpRespItemOfRoles.OsProfile.LinuxOperatingSystemProfile.Username != nil {
								tmpRespOfRoleslinuxOsProfile["username"] = *tmpRespItemOfRoles.OsProfile.LinuxOperatingSystemProfile.Username
							}
							if tmpRespItemOfRoles.OsProfile.LinuxOperatingSystemProfile.Password != nil {
								tmpRespOfRoleslinuxOsProfile["password"] = *tmpRespItemOfRoles.OsProfile.LinuxOperatingSystemProfile.Password
							}
							if tmpRespItemOfRoles.OsProfile.LinuxOperatingSystemProfile.SSHProfile != nil {
								if tmpRespItemOfRoles.OsProfile.LinuxOperatingSystemProfile.SSHProfile.PublicKeys != nil && len(*tmpRespItemOfRoles.OsProfile.LinuxOperatingSystemProfile.SSHProfile.PublicKeys) > 0 {
									tmpRespOfRoleslinuxOsProfilesshKeys := make([]interface{}, 0)
									for _, tmpRespItemOfRoleslinuxOsProfilesshKeys := range *tmpRespItemOfRoles.OsProfile.LinuxOperatingSystemProfile.SSHProfile.PublicKeys {
										tmpRespValueOfRoleslinuxOsProfilesshKeys := make(map[string]interface{})
										if tmpRespItemOfRoleslinuxOsProfilesshKeys.CertificateData != nil {
											tmpRespValueOfRoleslinuxOsProfilesshKeys["key_data"] = *tmpRespItemOfRoleslinuxOsProfilesshKeys.CertificateData
										}
										tmpRespOfRoleslinuxOsProfilesshKeys = append(tmpRespOfRoleslinuxOsProfilesshKeys, tmpRespValueOfRoleslinuxOsProfilesshKeys)
									}
									tmpRespOfRoleslinuxOsProfile["ssh_keys"] = tmpRespOfRoleslinuxOsProfilesshKeys
								}
							}
							tmpRespValueOfRoles["linux_os_profile"] = tmpRespOfRoleslinuxOsProfile
						}
					}
					if tmpRespItemOfRoles.VirtualNetworkProfile != nil {
						tmpRespOfRolesvnetProfile := make(map[string]interface{})
						if tmpRespItemOfRoles.VirtualNetworkProfile.ID != nil {
							tmpRespOfRolesvnetProfile["id"] = *tmpRespItemOfRoles.VirtualNetworkProfile.ID
						}
						if tmpRespItemOfRoles.VirtualNetworkProfile.Subnet != nil {
							tmpRespOfRolesvnetProfile["subnet"] = *tmpRespItemOfRoles.VirtualNetworkProfile.Subnet
						}
						tmpRespValueOfRoles["vnet_profile"] = tmpRespOfRolesvnetProfile
					}
					if tmpRespItemOfRoles.ScriptActions != nil && len(*tmpRespItemOfRoles.ScriptActions) > 0 {
						tmpRespOfRolesscriptActions := make([]interface{}, 0)
						for _, tmpRespItemOfRolesscriptActions := range *tmpRespItemOfRoles.ScriptActions {
							tmpRespValueOfRolesscriptActions := make(map[string]interface{})
							tmpRespValueOfRolesscriptActions["name"] = *tmpRespItemOfRolesscriptActions.Name
							tmpRespValueOfRolesscriptActions["uri"] = *tmpRespItemOfRolesscriptActions.URI
							tmpRespValueOfRolesscriptActions["parameters"] = *tmpRespItemOfRolesscriptActions.Parameters
							tmpRespOfRolesscriptActions = append(tmpRespOfRolesscriptActions, tmpRespValueOfRolesscriptActions)
						}
						tmpRespValueOfRoles["script_actions"] = tmpRespOfRolesscriptActions
					}
					switch *tmpRespItemOfRoles.Name {
					case "headnode":
						d.Set("head_node", tmpRespValueOfRoles)
					case "workernode":
						d.Set("worker_node", tmpRespValueOfRoles)
					case "zookeepernode":
						d.Set("zookeeper_node", tmpRespValueOfRoles)
					}
				}
			}
		}
		if response.Properties.ConnectivityEndpoints != nil && len(*response.Properties.ConnectivityEndpoints) > 0 {
			tmpRespOfConnectivityEndpoints := make([]interface{}, 0)
			for _, tmpRespItemOfConnectivityEndpoints := range *response.Properties.ConnectivityEndpoints {
				tmpRespValueOfConnectivityEndpoints := make(map[string]interface{})
				if tmpRespItemOfConnectivityEndpoints.Name != nil {
					tmpRespValueOfConnectivityEndpoints["name"] = *tmpRespItemOfConnectivityEndpoints.Name
				}
				if tmpRespItemOfConnectivityEndpoints.Protocol != nil {
					tmpRespValueOfConnectivityEndpoints["protocol"] = *tmpRespItemOfConnectivityEndpoints.Protocol
				}
				if tmpRespItemOfConnectivityEndpoints.Location != nil {
					tmpRespValueOfConnectivityEndpoints["location"] = *tmpRespItemOfConnectivityEndpoints.Location
				}
				if tmpRespItemOfConnectivityEndpoints.Port != nil {
					tmpRespValueOfConnectivityEndpoints["port"] = *tmpRespItemOfConnectivityEndpoints.Port
				}
				tmpRespOfConnectivityEndpoints = append(tmpRespOfConnectivityEndpoints, tmpRespValueOfConnectivityEndpoints)
			}
			d.Set("connectivity_endpoints", tmpRespOfConnectivityEndpoints)
		}
	}

	flattenAndSetTags(d, response.Tags)

	return nil
}

func resourceArmHDInsightClustersUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).hdInsightClustersClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroupName := d.Get("resource_group_name").(string)
	clusterName := d.Get("name").(string)
	parameters := hdinsight.ClusterPatchParameters{}
	tags := d.Get("tags").(map[string]interface{})
	tmpParamOfTags := expandTags(tags)
	parameters.Tags = tmpParamOfTags

	response, err := client.Update(ctx, resourceGroupName, clusterName, parameters)
	if err != nil {
		return fmt.Errorf("HD Insight Clusters update error: %+v", err)
	}

	if response.ID == nil {
		return fmt.Errorf("Cannot get the ID of HD Insight Clusters %q (Resource Group %q) ID", clusterName, resourceGroupName)
	}
	d.SetId(*response.ID)

	return resourceArmHDInsightClustersRead(d, meta)
}

func resourceArmHDInsightClustersDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).hdInsightClustersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroupName := id.ResourceGroup
	clusterName := id.Path["clusters"]

	future, err := client.Delete(ctx, resourceGroupName, clusterName)
	if err != nil {
		return fmt.Errorf("HD Insight Clusters deletion error: %+v", err)
	}
	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("HD Insight Clusters deletion future wait for completion error: %+v", err)
	}

	return nil
}
