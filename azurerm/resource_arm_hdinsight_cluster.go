package azurerm

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/hdinsight/mgmt/2015-03-01-preview/hdinsight"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type hdinsightClusterComputeProfile struct {
	headNode      *hdinsight.Role
	workerNode    *hdinsight.Role
	zookeeperNode *hdinsight.Role
}

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
				// TODO: validation
				// The name must be 59 characters or less and can contain letters, numbers, and hyphens (but the first and last character must be a letter or number).
				// The name can not contain a reserved key word.
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"tier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(hdinsight.Standard),
					string(hdinsight.Premium),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"cluster": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							ValidateFunc: validation.StringInSlice([]string{
								"hadoop",
								"hbase",
								"interactivehive",
								"kafka",
								"rserver",
								"storm",
								"spark",
							}, true),
						},

						"version": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validateHDInsightsClusterVersion,
						},

						"gateway": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
										ForceNew: true,
									},
									"username": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"password": {
										Type:      schema.TypeString,
										Required:  true,
										ForceNew:  true,
										Sensitive: true,
										// TODO: validation
										// The password must be at least 10 characters in length and must contain at least one digit, one uppercase and one lower case letter, one non-alphanumeric character (except characters ' " ` \).
									},
								},
							},
						},
					},
				},
			},

			"storage_profile": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Data Lake Stores aren't supported by the Swagger at this time, hence this is required
						"storage_account": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"storage_account_name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"storage_account_key": {
										Type:      schema.TypeString,
										Required:  true,
										ForceNew:  true,
										Sensitive: true,
									},
									"container_name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"is_default": {
										Type:     schema.TypeBool,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},

			"head_node": hdinsightClusterNodeProfile(2, true),

			"worker_node": hdinsightClusterNodeProfile(2, false),

			"zookeeper_node": hdinsightClusterNodeProfile(3, true),

			"tags": tagsSchema(),

			"connectivity_endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceArmHDInsightClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).hdinsightClustersClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for HDInsight Cluster creation.")

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)
	tier := d.Get("tier").(string)
	tags := d.Get("tags").(map[string]interface{})

	computeProfile := expandHDInsightClusterComputeProfile(d)
	clusterDefinition, clusterVersion := expandHDInsightClusterDetails(d)
	storageProfile := expandHDInsightStorageProfile(d)

	properties := hdinsight.ClusterCreateParametersExtended{
		Location: utils.String(location),
		Properties: &hdinsight.ClusterCreateProperties{
			OsType:            hdinsight.Linux,
			Tier:              hdinsight.Tier(tier),
			ClusterVersion:    utils.String(clusterVersion),
			ClusterDefinition: clusterDefinition,
			ComputeProfile:    computeProfile,
			StorageProfile:    storageProfile,
		},
		Tags: expandTags(tags),
	}

	future, err := client.Create(ctx, resourceGroup, name, properties)
	if err != nil {
		return fmt.Errorf("Error creating HDInsight Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for creation of HDInsight Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("[ERROR] Cannot read HDInsight Cluster %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmHDInsightClusterRead(d, meta)
}

func resourceArmHDInsightClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	clustersClient := meta.(*ArmClient).hdinsightClustersClient
	configurationsClient := meta.(*ArmClient).hdinsightConfigurationsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["clusters"]

	if d.HasChange("tags") {
		tags := d.Get("tags").(map[string]interface{})
		parameters := hdinsight.ClusterPatchParameters{
			Tags: expandTags(tags),
		}

		_, err = clustersClient.Update(ctx, resourceGroup, name, parameters)
		if err != nil {
			return fmt.Errorf("Error updating HDInsight Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	if d.HasChange("worker_node") {
		workerNode := expandHDInsightClusterNodeProfile("workernode", d.Get("worker_node"))
		parameters := hdinsight.ClusterResizeParameters{
			TargetInstanceCount: workerNode.TargetInstanceCount,
		}
		future, err := clustersClient.Resize(ctx, resourceGroup, name, parameters)
		if err != nil {
			return fmt.Errorf("Error resizing the number of worker nodes for HDInsight Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		err = future.WaitForCompletion(ctx, clustersClient.Client)
		if err != nil {
			return fmt.Errorf("Error waiting for resizing of worker nodes for HDInsight Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	if d.HasChange("cluster") {
		credentials := expandHDInsightsClusterGatewayCredentials(d)
		future, err := configurationsClient.UpdateHTTPSettings(ctx, resourceGroup, name, credentials)
		if err != nil {
			return fmt.Errorf("Error updating Gateway Connectivity Settings for HDInsight Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		err = future.WaitForCompletion(ctx, configurationsClient.Client)
		if err != nil {
			return fmt.Errorf("Error waiting for HDInsights Cluster %q (Resource Group %q) to finish updating Gateway Connectivity Settings: %+v", name, resourceGroup, err)
		}
	}

	return resourceArmHDInsightClusterRead(d, meta)
}

func resourceArmHDInsightClusterRead(d *schema.ResourceData, meta interface{}) error {
	clustersClient := meta.(*ArmClient).hdinsightClustersClient
	configurationsClient := meta.(*ArmClient).hdinsightConfigurationsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["clusters"]

	resp, err := clustersClient.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERROR] Error making Read request for HDInsight Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	configuration, err := configurationsClient.Get(ctx, resourceGroup, name, "gateway")
	if err != nil {
		return fmt.Errorf("[ERROR] Error loading configuration for HDInsight Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.Properties; props != nil {
		d.Set("tier", string(props.Tier))

		cluster := flattenHDInsightClusterDefinition(props, configuration)
		if err := d.Set("cluster", cluster); err != nil {
			return fmt.Errorf("Error setting `cluster`: %+v", err)
		}

		computeProfile, err := populateHDInsightClusterComputeProfile(props.ComputeProfile)
		if err != nil {
			return fmt.Errorf("Error parsing Compute Profile for HDInsight Cluster: %+v", err)
		}

		log.Printf("[TOMTOMTOMTOMTOM] Head: %d / Worker %d / Zookeeper %d",
			int(*computeProfile.headNode.TargetInstanceCount),
			int(*computeProfile.workerNode.TargetInstanceCount),
			int(*computeProfile.zookeeperNode.TargetInstanceCount))

		headNode := flattenHDInsightClusterComputeNodeProfile(computeProfile.headNode)
		if err := d.Set("head_node", headNode); err != nil {
			return fmt.Errorf("Error setting `head_node`: %+v", err)
		}

		workerNode := flattenHDInsightClusterComputeNodeProfile(computeProfile.workerNode)
		if err := d.Set("worker_node", workerNode); err != nil {
			return fmt.Errorf("Error setting `worker_node`: %+v", err)
		}

		zookeeperNode := flattenHDInsightClusterComputeNodeProfile(computeProfile.zookeeperNode)
		if err := d.Set("zookeeper_node", zookeeperNode); err != nil {
			return fmt.Errorf("Error setting `zookeeper_node`: %+v", err)
		}

		connectivityEndpoints := flattenHDInsightClusterConnectivityEndpoints(props.ConnectivityEndpoints)
		if err := d.Set("connectivity_endpoints", connectivityEndpoints); err != nil {
			return fmt.Errorf("Error setting `connectivity_endpoints`: %+v", err)
		}

		// this is a hack since storage_profile isn't returned from the API ¯\_(ツ)_/¯
		// TODO: API bug link
		storageProfile := expandHDInsightStorageProfile(d)
		flattenedStorageProfile := flattenHDInsightStorageProfile(storageProfile)
		if err := d.Set("storage_profile", flattenedStorageProfile); err != nil {
			return fmt.Errorf("Error setting `storage_profile`: %+v", err)
		}
	}
	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmHDInsightClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).hdinsightClustersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["clusters"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting HDInsight Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for deletion of HDInsight Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func hdinsightClusterNodeProfile(minTargetInstanceCount int, numberOfNodesForceNew bool) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Required: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"min_instance_count": {
					Type:     schema.TypeInt,
					Optional: true,
					ForceNew: true,
				},

				"target_instance_count": {
					Type:         schema.TypeInt,
					Required:     true,
					ForceNew:     numberOfNodesForceNew,
					ValidateFunc: validation.IntAtLeast(minTargetInstanceCount),
				},

				"hardware_profile": {
					Type:     schema.TypeList,
					MaxItems: 1,
					Required: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"vm_size": {
								Type:     schema.TypeString,
								Required: true,
								ForceNew: true,
							},
						},
					},
				},

				"os_profile": {
					Type:     schema.TypeSet,
					MaxItems: 1,
					Required: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"username": {
								Type:     schema.TypeString,
								Required: true,
								ForceNew: true,
							},

							"password": {
								Type:      schema.TypeString,
								Required:  true,
								ForceNew:  true,
								Sensitive: true,
							},

							// TODO: ssh auth
						},
					},
					Set: resourceAzureRMHDInsightClusterComputeNodeOSProfileHash,
				},

				"virtual_network_profile": {
					Type:     schema.TypeList,
					MaxItems: 1,
					Optional: true,
					ForceNew: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"virtual_network_id": {
								Type:     schema.TypeString,
								Required: true,
								ForceNew: true,
							},
							"subnet_id": {
								Type:     schema.TypeString,
								Required: true,
								ForceNew: true,
							},
						},
					},
				},
			},
		},
	}
}

func expandHDInsightClusterNodeProfile(name string, d interface{}) hdinsight.Role {
	vs := d.([]interface{})
	v := vs[0].(map[string]interface{})

	targetInstanceCount := v["target_instance_count"].(int)

	hardwareProfiles := v["hardware_profile"].([]interface{})
	hardwareProfile := hardwareProfiles[0].(map[string]interface{})

	vmSize := hardwareProfile["vm_size"].(string)

	osProfiles := v["os_profile"].(*schema.Set).List()
	osProfile := osProfiles[0].(map[string]interface{})

	username := osProfile["username"].(string)
	password := osProfile["password"].(string)

	role := hdinsight.Role{
		Name:                utils.String(name),
		TargetInstanceCount: utils.Int32(int32(targetInstanceCount)),
		HardwareProfile: &hdinsight.HardwareProfile{
			VMSize: utils.String(vmSize),
		},
		OsProfile: &hdinsight.OsProfile{
			LinuxOperatingSystemProfile: &hdinsight.LinuxOperatingSystemProfile{
				Username: utils.String(username),
				Password: utils.String(password),
			},
		},
	}

	if v, ok := v["min_instance_count"]; ok {
		val := v.(int)
		if val > 0 {
			role.MinInstanceCount = utils.Int32(int32(val))
		}
	}

	virtualNetworks := v["virtual_network_profile"].([]interface{})
	if len(virtualNetworks) > 0 {
		virtualNetwork := virtualNetworks[0].(map[string]interface{})

		virtualNetworkId := virtualNetwork["virtual_network_id"].(string)
		subnetId := virtualNetwork["subnet_id"].(string)
		role.VirtualNetworkProfile = &hdinsight.VirtualNetworkProfile{
			ID:     utils.String(virtualNetworkId),
			Subnet: utils.String(subnetId),
		}
	}

	return role
}

func flattenHDInsightClusterComputeNodeProfile(input *hdinsight.Role) []interface{} {
	output := make(map[string]interface{}, 0)

	if input != nil {
		hardwareOutput := make(map[string]interface{}, 0)
		if profile := input.HardwareProfile; profile != nil {
			if size := profile.VMSize; size != nil {
				hardwareOutput["vm_size"] = *size
			}
		}
		output["hardware_profile"] = []interface{}{hardwareOutput}

		if count := input.MinInstanceCount; count != nil {
			output["min_instance_count"] = int(*count)
		}
		if count := input.TargetInstanceCount; count != nil {
			output["target_instance_count"] = int(*count)
		}

		osProfile := make(map[string]interface{}, 0)

		if profile := input.OsProfile; profile != nil {
			if linux := profile.LinuxOperatingSystemProfile; linux != nil {
				if username := linux.Username; username != nil {
					osProfile["username"] = *username
				}
			}
		}

		output["os_profile"] = schema.NewSet(resourceAzureRMHDInsightClusterComputeNodeOSProfileHash, []interface{}{osProfile})
	}

	return []interface{}{output}
}

func expandHDInsightStorageProfile(d *schema.ResourceData) *hdinsight.StorageProfile {
	storageProfiles := d.Get("storage_profile").([]interface{})

	// when importing there's nothing in the state, so..
	if len(storageProfiles) == 0 {
		return nil
	}

	storageProfile := storageProfiles[0].(map[string]interface{})
	storageAccounts := storageProfile["storage_account"].([]interface{})

	accounts := make([]hdinsight.StorageAccount, 0)

	for _, v := range storageAccounts {
		storageAccount := v.(map[string]interface{})

		storageAccountName := storageAccount["storage_account_name"].(string)
		storageAccountKey := storageAccount["storage_account_key"].(string)
		storageContainerName := storageAccount["container_name"].(string)
		isDefault := storageAccount["is_default"].(bool)

		account := hdinsight.StorageAccount{
			Name:      utils.String(storageAccountName),
			Key:       utils.String(storageAccountKey),
			Container: utils.String(storageContainerName),
			IsDefault: utils.Bool(isDefault),
		}

		accounts = append(accounts, account)
	}

	return &hdinsight.StorageProfile{
		Storageaccounts: &accounts,
	}
}

func flattenHDInsightStorageProfile(input *hdinsight.StorageProfile) []interface{} {
	accounts := make([]interface{}, 0)

	if input != nil {
		for _, inputAccount := range *input.Storageaccounts {
			account := make(map[string]interface{}, 0)

			if name := inputAccount.Name; name != nil {
				account["storage_account_name"] = *name
			}

			if key := inputAccount.Key; key != nil {
				account["storage_account_key"] = *key
			}

			if container := inputAccount.Container; container != nil {
				account["container_name"] = *container
			}

			if isDefault := inputAccount.IsDefault; isDefault != nil {
				account["is_default"] = *isDefault
			}

			accounts = append(accounts, account)
		}
	}

	profile := make(map[string]interface{}, 0)
	profile["storage_account"] = accounts
	return []interface{}{profile}
}

func expandHDInsightsClusterGatewayCredentials(d *schema.ResourceData) hdinsight.HTTPConnectivitySettings {
	clusters := d.Get("cluster").([]interface{})
	cluster := clusters[0].(map[string]interface{})

	gateways := cluster["gateway"].([]interface{})
	gateway := gateways[0].(map[string]interface{})

	enabled := gateway["enabled"].(bool)
	username := gateway["username"].(string)
	password := gateway["password"].(string)

	if enabled {
		return hdinsight.HTTPConnectivitySettings{
			EnabledCredential: hdinsight.True,
			Username:          utils.String(username),
			Password:          utils.String(password),
		}
	}

	return hdinsight.HTTPConnectivitySettings{
		EnabledCredential: hdinsight.False,
	}
}

func expandHDInsightClusterDetails(d *schema.ResourceData) (*hdinsight.ClusterDefinition, string) {
	clusters := d.Get("cluster").([]interface{})
	cluster := clusters[0].(map[string]interface{})

	clusterKind := cluster["kind"].(string)
	clusterVersion := cluster["version"].(string)

	gatewayCredentials := expandHDInsightsClusterGatewayCredentials(d)
	gatewayConfig := map[string]interface{}{
		"restAuthCredential.isEnabled": gatewayCredentials.EnabledCredential == hdinsight.True,
	}

	if username := gatewayCredentials.Username; username != nil {
		gatewayConfig["restAuthCredential.username"] = *username
	}

	if password := gatewayCredentials.Password; password != nil {
		gatewayConfig["restAuthCredential.password"] = *password
	}

	definition := hdinsight.ClusterDefinition{
		Configurations: map[string]interface{}{
			"gateway": gatewayConfig,
		},
		Kind: utils.String(clusterKind),
	}

	return &definition, clusterVersion
}

func flattenHDInsightClusterDefinition(input *hdinsight.ClusterGetProperties, gatewayConfig hdinsight.SetString) []interface{} {
	output := make(map[string]interface{}, 0)

	if version := input.ClusterVersion; version != nil {
		segments := strings.Split(*version, ".")
		topLevelVersion := fmt.Sprintf("%s.%s", segments[0], segments[1])
		output["version"] = topLevelVersion
	}

	if definition := input.ClusterDefinition; definition != nil {
		if kind := definition.Kind; kind != nil {
			output["kind"] = *kind
		}
	}

	configuration := make(map[string]interface{}, 0)
	enabled := false
	if v := gatewayConfig.Value["restAuthCredential.isEnabled"]; v != nil {
		enabled, _ = strconv.ParseBool(*v)
		configuration["enabled"] = enabled
	}

	if enabled {
		if username := gatewayConfig.Value["restAuthCredential.username"]; username != nil {
			configuration["username"] = *username
		}

		if password := gatewayConfig.Value["restAuthCredential.password"]; password != nil {
			configuration["password"] = *password
		}
	}

	output["gateway"] = []interface{}{configuration}

	return []interface{}{output}
}

func resourceAzureRMHDInsightClusterComputeNodeOSProfileHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["username"].(string)))
	}

	return hashcode.String(buf.String())
}

func resourceAzureRMHDInsightClusterStorageAccountHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["storage_account_name"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["container_name"].(string)))
		buf.WriteString(fmt.Sprintf("%t-", m["is_default"].(bool)))
	}

	return hashcode.String(buf.String())
}

func flattenHDInsightClusterConnectivityEndpoints(input *[]hdinsight.ConnectivityEndpoint) []interface{} {
	outputs := make([]interface{}, 0)

	if input != nil {
		for _, v := range *input {
			output := make(map[string]interface{}, 0)

			if name := v.Name; name != nil {
				output["name"] = *name
			}
			if location := v.Location; location != nil {
				output["location"] = *location
			}
			if protocol := v.Protocol; protocol != nil {
				output["protocol"] = *protocol
			}
			if port := v.Port; port != nil {
				output["port"] = int(*port)
			}

			outputs = append(outputs, output)
		}
	}

	return outputs
}

func expandHDInsightClusterComputeProfile(d *schema.ResourceData) *hdinsight.ComputeProfile {
	headNode := expandHDInsightClusterNodeProfile("headnode", d.Get("head_node"))
	workerNode := expandHDInsightClusterNodeProfile("workernode", d.Get("worker_node"))
	zookeeperNode := expandHDInsightClusterNodeProfile("zookeepernode", d.Get("zookeeper_node"))

	return &hdinsight.ComputeProfile{
		Roles: &[]hdinsight.Role{
			headNode,
			workerNode,
			zookeeperNode,
		},
	}
}

func populateHDInsightClusterComputeProfile(input *hdinsight.ComputeProfile) (*hdinsightClusterComputeProfile, error) {
	output := hdinsightClusterComputeProfile{}

	if input != nil {
		if roles := input.Roles; roles != nil {
			for _, role := range *roles {
				roleName := *role.Name
				if strings.EqualFold(roleName, "headnode") {
					output.headNode = &role
				} else if strings.EqualFold(roleName, "workernode") {
					output.workerNode = &role
				} else if strings.EqualFold(roleName, "zookeepernode") {
					output.zookeeperNode = &role
				} else {
					return nil, fmt.Errorf("Unsupported node role: %q", roleName)
				}
			}
		}
	}

	return &output, nil
}

func validateHDInsightsClusterVersion(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	segments := strings.Split(v, ".")
	if len(segments) != 2 {
		errors = append(errors, fmt.Errorf("%q should contain 2 segments (e.g. `5.6`): '%q`", k, i))
	} else {
		if segments[0] == "" || segments[1] == "" {
			errors = append(errors, fmt.Errorf("%q: Expected a string in the format `X.Y` - got %q", k, v))
		}
	}

	return
}
