package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/servicefabric/mgmt/2018-02-01/servicefabric"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmServiceFabricCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceFabricClusterCreate,
		Read:   resourceArmServiceFabricClusterRead,
		Update: resourceArmServiceFabricClusterUpdate,
		Delete: resourceArmServiceFabricClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

			"reliability_level": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(servicefabric.ReliabilityLevelNone),
					string(servicefabric.ReliabilityLevelBronze),
					string(servicefabric.ReliabilityLevelSilver),
					string(servicefabric.ReliabilityLevelGold),
					string(servicefabric.ReliabilityLevelPlatinum),
				}, false),
			},

			"upgrade_mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(servicefabric.Automatic),
					string(servicefabric.Manual),
				}, false),
			},

			"cluster_code_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"management_endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vm_image": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"add_on_features": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"certificate": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"thumbprint": {
							Type:     schema.TypeString,
							Required: true,
						},
						"thumbprint_secondary": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"x509_store_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"client_certificate_thumbprint": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"thumbprint": {
							Type:     schema.TypeString,
							Required: true,
						},
						"is_admin": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},

			"diagnostics_config": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_account_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"protected_account_key_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"blob_endpoint": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"queue_endpoint": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"table_endpoint": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"fabric_settings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"parameters": {
							Type:     schema.TypeMap,
							Optional: true,
						},
					},
				},
			},

			"node_type": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"instance_count": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"is_primary": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},
						"client_endpoint_port": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"http_endpoint_port": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"durability_level": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(servicefabric.Bronze),
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(servicefabric.Bronze),
								string(servicefabric.Gold),
								string(servicefabric.Silver),
							}, false),
						},

						"application_ports": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_port": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: true,
									},
									"end_port": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},

						"ephemeral_ports": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_port": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: true,
									},
									"end_port": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},

			"tags": tagsSchema(),

			"cluster_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmServiceFabricClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceFabricClustersClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Service Fabric Cluster creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := d.Get("location").(string)
	reliabilityLevel := d.Get("reliability_level").(string)
	managementEndpoint := d.Get("management_endpoint").(string)
	upgradeMode := d.Get("upgrade_mode").(string)
	clusterCodeVersion := d.Get("cluster_code_version").(string)
	vmImage := d.Get("vm_image").(string)
	tags := d.Get("tags").(map[string]interface{})

	addOnFeaturesRaw := d.Get("add_on_features").(*schema.Set).List()
	addOnFeatures := expandServiceFabricClusterAddOnFeatures(addOnFeaturesRaw)

	certificateRaw := d.Get("certificate").([]interface{})
	certificate := expandServiceFabricClusterCertificate(certificateRaw)

	clientCertificateThumbprintRaw := d.Get("client_certificate_thumbprint").([]interface{})
	clientCertificateThumbprints := expandServiceFabricClusterClientCertificateThumbprints(clientCertificateThumbprintRaw)

	diagnosticsRaw := d.Get("diagnostics_config").([]interface{})
	diagnostics := expandServiceFabricClusterDiagnosticsConfig(diagnosticsRaw)

	fabricSettingsRaw := d.Get("fabric_settings").([]interface{})
	fabricSettings := expandServiceFabricClusterFabricSettings(fabricSettingsRaw)

	nodeTypesRaw := d.Get("node_type").([]interface{})
	nodeTypes := expandServiceFabricClusterNodeTypes(nodeTypesRaw)

	cluster := servicefabric.Cluster{
		Location: utils.String(location),
		Tags:     expandTags(tags),
		ClusterProperties: &servicefabric.ClusterProperties{
			AddOnFeatures:                   addOnFeatures,
			Certificate:                     certificate,
			ClientCertificateThumbprints:    clientCertificateThumbprints,
			DiagnosticsStorageAccountConfig: diagnostics,
			FabricSettings:                  fabricSettings,
			ManagementEndpoint:              utils.String(managementEndpoint),
			NodeTypes:                       nodeTypes,
			ReliabilityLevel:                servicefabric.ReliabilityLevel(reliabilityLevel),
			UpgradeMode:                     servicefabric.UpgradeMode(upgradeMode),
			VMImage:                         utils.String(vmImage),
		},
	}

	if clusterCodeVersion != "" {
		cluster.ClusterProperties.ClusterCodeVersion = utils.String(clusterCodeVersion)
	}

	future, err := client.Create(ctx, resourceGroup, name, cluster)
	if err != nil {
		return fmt.Errorf("Error creating Service Fabric Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for creation of Service Fabric Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Service Fabric Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read ID of Service Fabric Cluster %q (Resource Group %q)", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmServiceFabricClusterRead(d, meta)
}

func resourceArmServiceFabricClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceFabricClustersClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Service Fabric Cluster update.")

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	reliabilityLevel := d.Get("reliability_level").(string)
	upgradeMode := d.Get("upgrade_mode").(string)
	clusterCodeVersion := d.Get("cluster_code_version").(string)
	tags := d.Get("tags").(map[string]interface{})

	addOnFeaturesRaw := d.Get("add_on_features").(*schema.Set).List()
	addOnFeatures := expandServiceFabricClusterAddOnFeatures(addOnFeaturesRaw)

	certificateRaw := d.Get("certificate").([]interface{})
	certificate := expandServiceFabricClusterCertificate(certificateRaw)

	clientCertificateThumbprintsRaw := d.Get("client_certificate_thumbprint").([]interface{})
	clientCertificateThumbprints := expandServiceFabricClusterClientCertificateThumbprints(clientCertificateThumbprintsRaw)

	fabricSettingsRaw := d.Get("fabric_settings").([]interface{})
	fabricSettings := expandServiceFabricClusterFabricSettings(fabricSettingsRaw)

	nodeTypesRaw := d.Get("node_type").([]interface{})
	nodeTypes := expandServiceFabricClusterNodeTypes(nodeTypesRaw)

	parameters := servicefabric.ClusterUpdateParameters{
		ClusterPropertiesUpdateParameters: &servicefabric.ClusterPropertiesUpdateParameters{
			AddOnFeatures:                addOnFeatures,
			Certificate:                  certificate,
			ClientCertificateThumbprints: clientCertificateThumbprints,
			FabricSettings:               fabricSettings,
			NodeTypes:                    nodeTypes,
			ReliabilityLevel:             servicefabric.ReliabilityLevel1(reliabilityLevel),
			UpgradeMode:                  servicefabric.UpgradeMode1(upgradeMode),
		},
		Tags: expandTags(tags),
	}

	if clusterCodeVersion != "" {
		parameters.ClusterPropertiesUpdateParameters.ClusterCodeVersion = utils.String(clusterCodeVersion)
	}

	future, err := client.Update(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error updating Service Fabric Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for update of Service Fabric Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return resourceArmServiceFabricClusterRead(d, meta)
}

func resourceArmServiceFabricClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceFabricClustersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["clusters"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Service Fabric Cluster %q (Resource Group %q) was not found - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Service Fabric Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.ClusterProperties; props != nil {
		d.Set("cluster_code_version", props.ClusterCodeVersion)
		d.Set("cluster_endpoint", props.ClusterEndpoint)
		d.Set("management_endpoint", props.ManagementEndpoint)
		d.Set("reliability_level", string(props.ReliabilityLevel))
		d.Set("vm_image", props.VMImage)
		d.Set("upgrade_mode", string(props.UpgradeMode))

		addOnFeatures := flattenServiceFabricClusterAddOnFeatures(props.AddOnFeatures)
		if err := d.Set("add_on_features", schema.NewSet(schema.HashString, addOnFeatures)); err != nil {
			return fmt.Errorf("Error setting `add_on_features`: %+v", err)
		}

		certificate := flattenServiceFabricClusterCertificate(props.Certificate)
		if err := d.Set("certificate", certificate); err != nil {
			return fmt.Errorf("Error setting `certificate`: %+v", err)
		}

		clientCertificateThumbprints := flattenServiceFabricClusterClientCertificateThumbprints(props.ClientCertificateThumbprints)
		if err := d.Set("client_certificate_thumbprint", clientCertificateThumbprints); err != nil {
			return fmt.Errorf("Error setting `client_certificate_thumbprint`: %+v", err)
		}

		diagnostics := flattenServiceFabricClusterDiagnosticsConfig(props.DiagnosticsStorageAccountConfig)
		if err := d.Set("diagnostics_config", diagnostics); err != nil {
			return fmt.Errorf("Error setting `diagnostics_config`: %+v", err)
		}

		fabricSettings := flattenServiceFabricClusterFabricSettings(props.FabricSettings)
		if err := d.Set("fabric_settings", fabricSettings); err != nil {
			return fmt.Errorf("Error setting `fabric_settings`: %+v", err)
		}

		nodeTypes := flattenServiceFabricClusterNodeTypes(props.NodeTypes)
		if err := d.Set("node_type", nodeTypes); err != nil {
			return fmt.Errorf("Error setting `node_type`: %+v", err)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmServiceFabricClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceFabricClustersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["clusters"]

	log.Printf("[DEBUG] Deleting Service Fabric Cluster %q (Resource Group %q)", name, resourceGroup)

	resp, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting Service Fabric Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandServiceFabricClusterAddOnFeatures(input []interface{}) *[]string {
	output := make([]string, 0)

	for _, v := range input {
		output = append(output, v.(string))
	}

	return &output
}

func flattenServiceFabricClusterAddOnFeatures(input *[]string) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		for _, v := range *input {
			output = append(output, v)
		}
	}

	return output
}

func expandServiceFabricClusterCertificate(input []interface{}) *servicefabric.CertificateDescription {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	thumbprint := v["thumbprint"].(string)
	x509StoreName := v["x509_store_name"].(string)

	result := servicefabric.CertificateDescription{
		Thumbprint:    utils.String(thumbprint),
		X509StoreName: servicefabric.X509StoreName(x509StoreName),
	}

	if thumb, ok := v["thumbprint_secondary"]; ok {
		result.ThumbprintSecondary = utils.String(thumb.(string))
	}

	return &result
}

func flattenServiceFabricClusterCertificate(input *servicefabric.CertificateDescription) []interface{} {
	results := make([]interface{}, 0)

	if v := input; v != nil {
		output := make(map[string]interface{})

		if thumbprint := input.Thumbprint; thumbprint != nil {
			output["thumbprint"] = *thumbprint
		}

		if thumbprint := input.ThumbprintSecondary; thumbprint != nil {
			output["thumbprint_secondary"] = *thumbprint
		}

		output["x509_store_name"] = string(input.X509StoreName)
		results = append(results, output)
	}

	return results
}

func expandServiceFabricClusterClientCertificateThumbprints(input []interface{}) *[]servicefabric.ClientCertificateThumbprint {
	results := make([]servicefabric.ClientCertificateThumbprint, 0)

	for _, v := range input {
		val := v.(map[string]interface{})

		thumbprint := val["thumbprint"].(string)
		isAdmin := val["is_admin"].(bool)

		result := servicefabric.ClientCertificateThumbprint{
			CertificateThumbprint: utils.String(thumbprint),
			IsAdmin:               utils.Bool(isAdmin),
		}
		results = append(results, result)
	}

	return &results
}

func flattenServiceFabricClusterClientCertificateThumbprints(input *[]servicefabric.ClientCertificateThumbprint) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)

	for _, v := range *input {
		result := make(map[string]interface{})

		if thumbprint := v.CertificateThumbprint; thumbprint != nil {
			result["thumbprint"] = *thumbprint
		}

		if isAdmin := v.IsAdmin; isAdmin != nil {
			result["is_admin"] = *isAdmin
		}

		results = append(results, result)
	}

	return results
}

func expandServiceFabricClusterDiagnosticsConfig(input []interface{}) *servicefabric.DiagnosticsStorageAccountConfig {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	storageAccountName := v["storage_account_name"].(string)
	protectedAccountKeyName := v["protected_account_key_name"].(string)
	blobEndpoint := v["blob_endpoint"].(string)
	queueEndpoint := v["queue_endpoint"].(string)
	tableEndpoint := v["table_endpoint"].(string)

	config := servicefabric.DiagnosticsStorageAccountConfig{
		StorageAccountName:      utils.String(storageAccountName),
		ProtectedAccountKeyName: utils.String(protectedAccountKeyName),
		BlobEndpoint:            utils.String(blobEndpoint),
		QueueEndpoint:           utils.String(queueEndpoint),
		TableEndpoint:           utils.String(tableEndpoint),
	}
	return &config
}

func flattenServiceFabricClusterDiagnosticsConfig(input *servicefabric.DiagnosticsStorageAccountConfig) []interface{} {
	results := make([]interface{}, 0)

	if v := input; v != nil {
		output := make(map[string]interface{})

		if name := v.StorageAccountName; name != nil {
			output["storage_account_name"] = *name
		}

		if name := v.ProtectedAccountKeyName; name != nil {
			output["protected_account_key_name"] = *name
		}

		if endpoint := v.BlobEndpoint; endpoint != nil {
			output["blob_endpoint"] = *endpoint
		}

		if endpoint := v.QueueEndpoint; endpoint != nil {
			output["queue_endpoint"] = *endpoint
		}

		if endpoint := v.TableEndpoint; endpoint != nil {
			output["table_endpoint"] = *endpoint
		}

		results = append(results, output)
	}

	return results
}

func expandServiceFabricClusterFabricSettings(input []interface{}) *[]servicefabric.SettingsSectionDescription {
	results := make([]servicefabric.SettingsSectionDescription, 0)

	for _, v := range input {
		val := v.(map[string]interface{})

		name := val["name"].(string)
		params := make([]servicefabric.SettingsParameterDescription, 0)
		paramsRaw := val["parameters"].(map[string]interface{})
		for k, v := range paramsRaw {
			param := servicefabric.SettingsParameterDescription{
				Name:  utils.String(k),
				Value: utils.String(v.(string)),
			}
			params = append(params, param)
		}

		result := servicefabric.SettingsSectionDescription{
			Name:       utils.String(name),
			Parameters: &params,
		}
		results = append(results, result)
	}

	return &results
}

func flattenServiceFabricClusterFabricSettings(input *[]servicefabric.SettingsSectionDescription) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)

	for _, v := range *input {
		result := make(map[string]interface{})

		if name := v.Name; name != nil {
			result["name"] = *name
		}

		parameters := make(map[string]interface{})
		if paramsRaw := v.Parameters; paramsRaw != nil {
			for _, p := range *paramsRaw {
				if p.Name == nil || p.Value == nil {
					continue
				}

				parameters[*p.Name] = *p.Value
			}
		}
		result["parameters"] = parameters
		results = append(results, result)
	}

	return results
}

func expandServiceFabricClusterNodeTypes(input []interface{}) *[]servicefabric.NodeTypeDescription {
	results := make([]servicefabric.NodeTypeDescription, 0)

	for _, v := range input {
		node := v.(map[string]interface{})

		name := node["name"].(string)
		instanceCount := node["instance_count"].(int)
		clientEndpointPort := node["client_endpoint_port"].(int)
		httpEndpointPort := node["http_endpoint_port"].(int)
		isPrimary := node["is_primary"].(bool)
		durabilityLevel := node["durability_level"].(string)

		result := servicefabric.NodeTypeDescription{
			Name:                         utils.String(name),
			VMInstanceCount:              utils.Int32(int32(instanceCount)),
			IsPrimary:                    utils.Bool(isPrimary),
			ClientConnectionEndpointPort: utils.Int32(int32(clientEndpointPort)),
			HTTPGatewayEndpointPort:      utils.Int32(int32(httpEndpointPort)),
			DurabilityLevel:              servicefabric.DurabilityLevel(durabilityLevel),
		}

		applicationPortsRaw := node["application_ports"].([]interface{})
		if len(applicationPortsRaw) > 0 {
			portsRaw := applicationPortsRaw[0].(map[string]interface{})

			startPort := portsRaw["start_port"].(int)
			endPort := portsRaw["end_port"].(int)

			result.ApplicationPorts = &servicefabric.EndpointRangeDescription{
				StartPort: utils.Int32(int32(startPort)),
				EndPort:   utils.Int32(int32(endPort)),
			}
		}

		ephemeralPortsRaw := node["ephemeral_ports"].([]interface{})
		if len(ephemeralPortsRaw) > 0 {
			portsRaw := ephemeralPortsRaw[0].(map[string]interface{})

			startPort := portsRaw["start_port"].(int)
			endPort := portsRaw["end_port"].(int)

			result.EphemeralPorts = &servicefabric.EndpointRangeDescription{
				StartPort: utils.Int32(int32(startPort)),
				EndPort:   utils.Int32(int32(endPort)),
			}
		}

		results = append(results, result)
	}

	return &results
}

func flattenServiceFabricClusterNodeTypes(input *[]servicefabric.NodeTypeDescription) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)

	for _, v := range *input {
		output := make(map[string]interface{})

		if name := v.Name; name != nil {
			output["name"] = *name
		}

		if count := v.VMInstanceCount; count != nil {
			output["instance_count"] = int(*count)
		}

		if primary := v.IsPrimary; primary != nil {
			output["is_primary"] = *primary
		}

		if port := v.ClientConnectionEndpointPort; port != nil {
			output["client_endpoint_port"] = *port
		}

		if port := v.HTTPGatewayEndpointPort; port != nil {
			output["http_endpoint_port"] = *port
		}

		output["durability_level"] = string(v.DurabilityLevel)

		applicationPorts := make([]interface{}, 0)
		if ports := v.ApplicationPorts; ports != nil {
			r := make(map[string]interface{})
			if start := ports.StartPort; start != nil {
				r["start_port"] = int(*start)
			}
			if end := ports.EndPort; end != nil {
				r["end_port"] = int(*end)
			}
			applicationPorts = append(applicationPorts, r)
		}
		output["application_ports"] = applicationPorts

		ephermeralPorts := make([]interface{}, 0)
		if ports := v.EphemeralPorts; ports != nil {
			r := make(map[string]interface{})
			if start := ports.StartPort; start != nil {
				r["start_port"] = int(*start)
			}
			if end := ports.EndPort; end != nil {
				r["end_port"] = int(*end)
			}
			ephermeralPorts = append(ephermeralPorts, r)
		}
		output["ephemeral_ports"] = ephermeralPorts

		results = append(results, output)
	}

	return results
}
