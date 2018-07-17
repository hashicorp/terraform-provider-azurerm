package azurerm

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/hdinsight/mgmt/2015-03-01-preview/hdinsight"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
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
		CustomizeDiff: resourceArmHDInsightClusterCustomizeDiff,

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
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							// TODO: is `latest` supported here?
							//ValidateFunc: validateHDInsightsClusterVersion,
						},

						"component_versions": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
						},

						"gateway": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Required: true,
										// TODO: can this be updated?
										//ForceNew: true,
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

			"head_node": hdinsightClusterNodeProfile("head_node", true, 2, true),

			"worker_node": hdinsightClusterNodeProfile("worker_node", true, 2, false),

			"zookeeper_node": hdinsightClusterNodeProfile("zookeeper_node", false, 3, true),

			"tags": tagsSchema(),

			"ssh_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"https_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
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

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error checking for the existence of HDInsights Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	if resp.StatusCode == http.StatusOK {
		return fmt.Errorf("An HDInsights Cluster already exists with the name %q in the Resource Group %q - please import it into the state", name, resourceGroup)
	}

	computeProfile, err := expandHDInsightClusterComputeProfile(d)
	if err != nil {
		return err
	}

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
		if response.WasOK(future.Response()) {
			err = resourceArmHDInsightClusterReadError(client, ctx, resourceGroup, name)
		}

		return fmt.Errorf("Error creating HDInsight Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		if response.WasOK(future.Response()) {
			err = resourceArmHDInsightClusterReadError(client, ctx, resourceGroup, name)
		}

		return fmt.Errorf("Error waiting for creation of HDInsight Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving HDInsights Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("[ERROR] Cannot read HDInsight Cluster %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmHDInsightClusterUpdate(d, meta)
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
			err = resourceArmHDInsightClusterReadError(clustersClient, ctx, resourceGroup, name)
			return fmt.Errorf("Error updating HDInsight Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	if d.HasChange("worker_node") {
		workerNode, err := expandHDInsightClusterNodeProfile("workernode", d.Get("worker_node"))
		if err != nil {
			return fmt.Errorf("Error expanding `worker_node`: %+v", err)
		}

		parameters := hdinsight.ClusterResizeParameters{
			TargetInstanceCount: workerNode.TargetInstanceCount,
		}
		future, err := clustersClient.Resize(ctx, resourceGroup, name, parameters)
		if err != nil {
			return fmt.Errorf("Error resizing the number of worker nodes for HDInsight Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		err = future.WaitForCompletionRef(ctx, clustersClient.Client)
		if err != nil {
			return fmt.Errorf("Error waiting for resizing of worker nodes for HDInsight Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	if d.HasChange("cluster") {
		credentials := expandHDInsightsClusterGatewayCredentials(d)
		future, err := configurationsClient.UpdateHTTPSettings(ctx, resourceGroup, name, "gateway", credentials)
		if err != nil {
			return fmt.Errorf("Error updating Gateway Connectivity Settings for HDInsight Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		err = future.WaitForCompletionRef(ctx, configurationsClient.Client)
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
		// Swagger states that this'll only return a 200/201, but it actually returns a 204 no content in some cases ¯\_(ツ)_/¯
		if utils.ResponseWasNotFound(resp.Response) || utils.ResponseWasNoContent(resp.Response) {
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

		if endpoints := props.ConnectivityEndpoints; endpoints != nil {
			for _, endpoint := range *endpoints {
				if v := endpoint.Name; v != nil {
					if strings.EqualFold(*v, "HTTPS") {
						d.Set("https_endpoint", endpoint.Location)
					}

					if strings.EqualFold(*v, "ssh_endpoint") {
						d.Set("ssh_endpoint", endpoint.Location)
					}
				}
			}
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
		// Swagger states that this'll only return a 200/201, but it actually returns a 204 no content in some cases ¯\_(ツ)_/¯
		if response.WasNotFound(future.Response()) || response.WasNoContent(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting HDInsight Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		// Swagger states that this'll only return a 200/201, but it actually returns a 204 no content in some cases ¯\_(ツ)_/¯
		if response.WasNotFound(future.Response()) || response.WasNoContent(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error waiting for deletion of HDInsight Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func hdinsightClusterNodeProfile(blockName string, required bool, minTargetInstanceCount int, numberOfNodesForceNew bool) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: !required,
		Required: required,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
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
								Type:          schema.TypeString,
								Optional:      true,
								ForceNew:      true,
								Sensitive:     true,
								ConflictsWith: []string{fmt.Sprintf("%s.0.os_profile.0.ssh_key", blockName)},
							},

							"ssh_key": {
								Type:          schema.TypeString,
								Optional:      true,
								ForceNew:      true,
								ConflictsWith: []string{fmt.Sprintf("%s.0.os_profile.0.password", blockName)},
							},
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

func expandHDInsightClusterNodeProfile(name string, d interface{}) (*hdinsight.Role, error) {
	vs := d.([]interface{})
	if len(vs) == 0 {
		// optional
		return nil, nil
	}

	v := vs[0].(map[string]interface{})

	targetInstanceCount := v["target_instance_count"].(int)

	hardwareProfiles := v["hardware_profile"].([]interface{})
	hardwareProfile := hardwareProfiles[0].(map[string]interface{})

	vmSize := hardwareProfile["vm_size"].(string)

	osProfiles := v["os_profile"].(*schema.Set).List()
	osProfile := osProfiles[0].(map[string]interface{})

	username := osProfile["username"].(string)

	role := hdinsight.Role{
		Name:                utils.String(name),
		TargetInstanceCount: utils.Int32(int32(targetInstanceCount)),
		HardwareProfile: &hdinsight.HardwareProfile{
			VMSize: utils.String(vmSize),
		},
		OsProfile: &hdinsight.OsProfile{
			LinuxOperatingSystemProfile: &hdinsight.LinuxOperatingSystemProfile{
				Username: utils.String(username),
			},
		},
	}

	// validation handled in CustomizeDiff
	password := osProfile["password"].(string)
	sshKey := osProfile["ssh_key"].(string)
	if password != "" {
		role.OsProfile.LinuxOperatingSystemProfile.Password = utils.String(password)
	} else {
		role.OsProfile.LinuxOperatingSystemProfile.SSHProfile = &hdinsight.SSHProfile{
			PublicKeys: &[]hdinsight.SSHPublicKey{
				{
					CertificateData: utils.String(sshKey),
				},
			},
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

	return &role, nil
}

func flattenHDInsightClusterComputeNodeProfile(input *hdinsight.Role) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{}, 0)

	if input != nil {
		hardwareOutput := make(map[string]interface{}, 0)
		if profile := input.HardwareProfile; profile != nil {
			if size := profile.VMSize; size != nil {
				hardwareOutput["vm_size"] = *size
			}
		}
		output["hardware_profile"] = []interface{}{hardwareOutput}

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

func expandHDInsightsClusterGatewayCredentials(d *schema.ResourceData) map[string]*string {
	clusters := d.Get("cluster").([]interface{})
	cluster := clusters[0].(map[string]interface{})

	gateways := cluster["gateway"].([]interface{})
	gateway := gateways[0].(map[string]interface{})

	username := gateway["username"].(string)
	password := gateway["password"].(string)

	return map[string]*string{
		// hard-coded to true because: "Linux clusters do not support revoking HTTP credentials."
		"restAuthCredential.isEnabled": utils.String("true"),

		// these have to be specified, even if it's Disabled otherwise we get the /totally helpful/ response:
		// {"code":"BadRequest","message":"User input validation failed. Errors: The request payload is invalid."}
		"restAuthCredential.username": utils.String(username),
		"restAuthCredential.password": utils.String(password),
	}
}

func expandHDInsightClusterDetails(d *schema.ResourceData) (*hdinsight.ClusterDefinition, string) {
	clusters := d.Get("cluster").([]interface{})
	cluster := clusters[0].(map[string]interface{})

	clusterKind := cluster["kind"].(string)
	clusterVersion := cluster["version"].(string)

	componentVersions := make(map[string]*string, 0)
	components := cluster["component_versions"].(map[string]interface{})
	for key, value := range components {
		componentVersions[key] = utils.String(value.(string))
	}

	gatewayCredentials := expandHDInsightsClusterGatewayCredentials(d)
	definition := hdinsight.ClusterDefinition{
		Configurations: map[string]interface{}{
			"gateway": gatewayCredentials,
		},
		Kind:             utils.String(clusterKind),
		ComponentVersion: componentVersions,
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

		// TODO: verify that these are returned
		componentVersions := make(map[string]interface{}, 0)
		if versions := definition.ComponentVersion; versions != nil {
			for k, v := range versions {
				componentVersions[k] = *v
			}
		}
		output["component_versions"] = componentVersions
	}

	configuration := make(map[string]interface{}, 0)
	if username := gatewayConfig.Value["restAuthCredential.username"]; username != nil {
		configuration["username"] = *username
	}

	if password := gatewayConfig.Value["restAuthCredential.password"]; password != nil {
		configuration["password"] = *password
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

func expandHDInsightClusterComputeProfile(d *schema.ResourceData) (*hdinsight.ComputeProfile, error) {
	headNode, err := expandHDInsightClusterNodeProfile("headnode", d.Get("head_node"))
	if err != nil {
		return nil, fmt.Errorf("Error expanding `head_node`: %+v", err)
	}

	workerNode, err := expandHDInsightClusterNodeProfile("workernode", d.Get("worker_node"))
	if err != nil {
		return nil, fmt.Errorf("Error expanding `worker_node`: %+v", err)
	}

	zookeeperNode, err := expandHDInsightClusterNodeProfile("zookeepernode", d.Get("zookeeper_node"))
	if err != nil {
		return nil, fmt.Errorf("Error expanding `zookeeper_node`: %+v", err)
	}

	roles := []hdinsight.Role{
		*headNode,
		*workerNode,
	}

	if zookeeperNode != nil {
		roles = append(roles, *zookeeperNode)
	}

	return &hdinsight.ComputeProfile{
		Roles: &roles,
	}, nil
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

func resourceArmHDInsightClusterCustomizeDiff(d *schema.ResourceDiff, v interface{}) error {
	err := resourceArmHDInsightClusterCustomizeDiffUsernameOrSshKey(d, "head_node")
	if err != nil {
		return err
	}

	err = resourceArmHDInsightClusterCustomizeDiffUsernameOrSshKey(d, "worker_node")
	if err != nil {
		return err
	}

	err = resourceArmHDInsightClusterCustomizeDiffUsernameOrSshKey(d, "zookeeper_node")
	if err != nil {
		return err
	}

	/*
		if v, ok := d.GetOk("cluster"); ok {
			clusters := v.([]interface{})
			if len(clusters) > 0 {
				cluster := clusters[0].(map[string]interface{})
				kind, ok := cluster["kind"].(string)
				if !ok {
					return nil
				}

				componentVersions, ok := cluster["component_versions"]
				if !ok {
					return nil
				}

				found := false
				for key, value := range componentVersions.(map[string]interface{}) {

				}


				if !false {
					return fmt.Errorf("The `` block in `cluster` must contain")
				}
				// component_versions
			}
		}

		// TODO: check the Cluster Type is present in the `cluster_versions` map?
	*/

	return nil
}

func resourceArmHDInsightClusterCustomizeDiffUsernameOrSshKey(d *schema.ResourceDiff, field string) error {
	if nodes, ok := d.GetOk(field); ok {
		node := nodes.([]interface{})[0].(map[string]interface{})
		if profiles, ok := node["os_profile"].([]interface{}); ok {
			profile := profiles[0].(map[string]interface{})
			hasPassword := profile["password"].(string) != ""
			hasSshKey := profile["ssh_key"].(string) != ""

			if !hasPassword && !hasSshKey {
				return fmt.Errorf("Either a `password` or a `ssh_key` must be specified for the `os_profile` block in `%s`", field)
			}
		}
	}

	return nil
}

func resourceArmHDInsightClusterReadError(client hdinsight.ClustersClient, ctx context.Context, resourceGroup string, name string) error {
	// the HDInsights errors returns errors as a 200 to the SDK
	// meaning we need to retrieve the cluster to get the error
	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving HDInsights Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if props := resp.Properties; props != nil {
		if errors := props.Errors; errors != nil {
			var err error

			for _, e := range *errors {
				if message := e.Message; message != nil {
					err = multierror.Append(err, fmt.Errorf(*message))
				}
			}

			return fmt.Errorf("Error updating HDInsight Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}
