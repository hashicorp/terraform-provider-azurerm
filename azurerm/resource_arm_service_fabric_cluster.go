package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/servicefabric/mgmt/2018-02-01/servicefabric"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmServiceFabricCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceFabricClusterCreateUpdate,
		Read:   resourceArmServiceFabricClusterRead,
		Update: resourceArmServiceFabricClusterCreateUpdate,
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

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

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
					string(servicefabric.UpgradeModeAutomatic),
					string(servicefabric.UpgradeModeManual),
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

			"azure_active_directory": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.UUID,
						},
						"cluster_application_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.UUID,
						},
						"client_application_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.UUID,
						},
					},
				},
			},

			"certificate": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"certificate_common_names"},
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

			"certificate_common_names": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"certificate"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"common_names": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate_common_name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.NoEmptyStrings,
									},
									"certificate_issuer_thumbprint": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validate.NoEmptyStrings,
									},
								},
							},
						},
						"x509_store_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"reverse_proxy_certificate": {
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
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_account_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"protected_account_key_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"blob_endpoint": {
							Type:     schema.TypeString,
							Required: true,
						},
						"queue_endpoint": {
							Type:     schema.TypeString,
							Required: true,
						},
						"table_endpoint": {
							Type:     schema.TypeString,
							Required: true,
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
						},
						"placement_properties": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"capacities": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"instance_count": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"is_primary": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"client_endpoint_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"http_endpoint_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"reverse_proxy_endpoint_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validate.PortNumber,
						},
						"durability_level": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(servicefabric.Bronze),
							ValidateFunc: validation.StringInSlice([]string{
								string(servicefabric.Bronze),
								string(servicefabric.Gold),
								string(servicefabric.Silver),
							}, false),
						},

						"application_ports": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_port": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"end_port": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},

						"ephemeral_ports": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_port": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"end_port": {
										Type:     schema.TypeInt,
										Required: true,
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

func resourceArmServiceFabricClusterCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceFabric.ClustersClient
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

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Service Fabric Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_service_fabric_cluster", *existing.ID)
		}
	}

	addOnFeaturesRaw := d.Get("add_on_features").(*schema.Set).List()
	addOnFeatures := expandServiceFabricClusterAddOnFeatures(addOnFeaturesRaw)

	azureActiveDirectoryRaw := d.Get("azure_active_directory").([]interface{})
	azureActiveDirectory := expandServiceFabricClusterAzureActiveDirectory(azureActiveDirectoryRaw)

	reverseProxyCertificateRaw := d.Get("reverse_proxy_certificate").([]interface{})
	reverseProxyCertificate := expandServiceFabricClusterReverseProxyCertificate(reverseProxyCertificateRaw)

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
			AzureActiveDirectory:            azureActiveDirectory,
			CertificateCommonNames:          expandServiceFabricClusterCertificateCommonNames(d),
			ReverseProxyCertificate:         reverseProxyCertificate,
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

	if certificateRaw, ok := d.GetOk("certificate"); ok {
		certificate := expandServiceFabricClusterCertificate(certificateRaw.([]interface{}))
		cluster.ClusterProperties.Certificate = certificate
	}

	if clusterCodeVersion != "" {
		cluster.ClusterProperties.ClusterCodeVersion = utils.String(clusterCodeVersion)
	}

	future, err := client.Create(ctx, resourceGroup, name, cluster)
	if err != nil {
		return fmt.Errorf("Error creating Service Fabric Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
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

func resourceArmServiceFabricClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceFabric.ClustersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
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
		d.Set("location", azure.NormalizeLocation(*location))
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

		azureActiveDirectory := flattenServiceFabricClusterAzureActiveDirectory(props.AzureActiveDirectory)
		if err := d.Set("azure_active_directory", azureActiveDirectory); err != nil {
			return fmt.Errorf("Error setting `azure_active_directory`: %+v", err)
		}

		certificate := flattenServiceFabricClusterCertificate(props.Certificate)
		if err := d.Set("certificate", certificate); err != nil {
			return fmt.Errorf("Error setting `certificate`: %+v", err)
		}

		certificateCommonNames := flattenServiceFabricClusterCertificateCommonNames(props.CertificateCommonNames)
		if err := d.Set("certificate_common_names", certificateCommonNames); err != nil {
			return fmt.Errorf("Error setting `certificate_common_names`: %+v", err)
		}

		reverseProxyCertificate := flattenServiceFabricClusterReverseProxyCertificate(props.ReverseProxyCertificate)
		if err := d.Set("reverse_proxy_certificate", reverseProxyCertificate); err != nil {
			return fmt.Errorf("Error setting `reverse_proxy_certificate`: %+v", err)
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
	client := meta.(*ArmClient).serviceFabric.ClustersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
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

func expandServiceFabricClusterAzureActiveDirectory(input []interface{}) *servicefabric.AzureActiveDirectory {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	tenantId := v["tenant_id"].(string)
	clusterApplication := v["cluster_application_id"].(string)
	clientApplication := v["client_application_id"].(string)

	config := servicefabric.AzureActiveDirectory{
		TenantID:           utils.String(tenantId),
		ClusterApplication: utils.String(clusterApplication),
		ClientApplication:  utils.String(clientApplication),
	}
	return &config
}

func flattenServiceFabricClusterAzureActiveDirectory(input *servicefabric.AzureActiveDirectory) []interface{} {
	results := make([]interface{}, 0)

	if v := input; v != nil {
		output := make(map[string]interface{})

		if name := v.TenantID; name != nil {
			output["tenant_id"] = *name
		}

		if name := v.ClusterApplication; name != nil {
			output["cluster_application_id"] = *name
		}

		if endpoint := v.ClientApplication; endpoint != nil {
			output["client_application_id"] = *endpoint
		}

		results = append(results, output)
	}

	return results
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

func expandServiceFabricClusterCertificateCommonNames(d *schema.ResourceData) *servicefabric.ServerCertificateCommonNames {
	i := d.Get("certificate_common_names").([]interface{})
	if len(i) <= 0 || i[0] == nil {
		return nil
	}
	input := i[0].(map[string]interface{})

	commonNamesRaw := input["common_names"].(*schema.Set).List()
	commonNames := make([]servicefabric.ServerCertificateCommonName, 0)

	for _, commonName := range commonNamesRaw {
		commonNameDetails := commonName.(map[string]interface{})
		certificateCommonName := commonNameDetails["certificate_common_name"].(string)
		certificateIssuerThumbprint := commonNameDetails["certificate_issuer_thumbprint"].(string)

		commonName := servicefabric.ServerCertificateCommonName{
			CertificateCommonName:       &certificateCommonName,
			CertificateIssuerThumbprint: &certificateIssuerThumbprint,
		}

		commonNames = append(commonNames, commonName)
	}

	x509StoreName := input["x509_store_name"].(string)

	output := servicefabric.ServerCertificateCommonNames{
		CommonNames:   &commonNames,
		X509StoreName: servicefabric.X509StoreName1(x509StoreName),
	}

	return &output
}

func flattenServiceFabricClusterCertificateCommonNames(in *servicefabric.ServerCertificateCommonNames) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})

	if commonNames := in.CommonNames; commonNames != nil {
		common_names := make([]map[string]interface{}, 0)
		for _, i := range *commonNames {
			commonName := make(map[string]interface{})

			if i.CertificateCommonName != nil {
				commonName["certificate_common_name"] = *i.CertificateCommonName
			}

			if i.CertificateIssuerThumbprint != nil {
				commonName["certificate_issuer_thumbprint"] = *i.CertificateIssuerThumbprint
			}

			common_names = append(common_names, commonName)
		}

		output["common_names"] = common_names
	}

	output["x509_store_name"] = string(in.X509StoreName)

	return []interface{}{output}
}

func expandServiceFabricClusterReverseProxyCertificate(input []interface{}) *servicefabric.CertificateDescription {
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

func flattenServiceFabricClusterReverseProxyCertificate(input *servicefabric.CertificateDescription) []interface{} {
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

		if props, ok := node["placement_properties"]; ok {
			placementProperties := make(map[string]*string)
			for key, value := range props.(map[string]interface{}) {
				placementProperties[key] = utils.String(value.(string))
			}

			result.PlacementProperties = placementProperties
		}

		if caps, ok := node["capacities"]; ok {
			capacities := make(map[string]*string)
			for key, value := range caps.(map[string]interface{}) {
				capacities[key] = utils.String(value.(string))
			}

			result.Capacities = capacities
		}

		if v := int32(node["reverse_proxy_endpoint_port"].(int)); v != 0 {
			result.ReverseProxyEndpointPort = utils.Int32(v)
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

		if placementProperties := v.PlacementProperties; placementProperties != nil {
			output["placement_properties"] = placementProperties
		}

		if capacities := v.Capacities; capacities != nil {
			output["capacities"] = capacities
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

		if port := v.ReverseProxyEndpointPort; port != nil {
			output["reverse_proxy_endpoint_port"] = *port
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

		ephemeralPorts := make([]interface{}, 0)
		if ports := v.EphemeralPorts; ports != nil {
			r := make(map[string]interface{})
			if start := ports.StartPort; start != nil {
				r["start_port"] = int(*start)
			}
			if end := ports.EndPort; end != nil {
				r["end_port"] = int(*end)
			}
			ephemeralPorts = append(ephemeralPorts, r)
		}
		output["ephemeral_ports"] = ephemeralPorts

		results = append(results, output)
	}

	return results
}
