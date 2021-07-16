package datafactory

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/validate"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataFactoryIntegrationRuntimeAzureSsis() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryIntegrationRuntimeAzureSsisCreateUpdate,
		Read:   resourceDataFactoryIntegrationRuntimeAzureSsisRead,
		Update: resourceDataFactoryIntegrationRuntimeAzureSsisCreateUpdate,
		Delete: resourceDataFactoryIntegrationRuntimeAzureSsisDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^([a-zA-Z0-9](-|-?[a-zA-Z0-9]+)+[a-zA-Z0-9])$`),
					`Invalid name for Managed Integration Runtime: minimum 3 characters, must start and end with a number or a letter, may only consist of letters, numbers and dashes and no consecutive dashes.`,
				),
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"data_factory_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataFactoryName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"node_size": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Standard_D2_v3",
					"Standard_D4_v3",
					"Standard_D8_v3",
					"Standard_D16_v3",
					"Standard_D32_v3",
					"Standard_D64_v3",
					"Standard_E2_v3",
					"Standard_E4_v3",
					"Standard_E8_v3",
					"Standard_E16_v3",
					"Standard_E32_v3",
					"Standard_E64_v3",
					"Standard_D1_v2",
					"Standard_D2_v2",
					"Standard_D3_v2",
					"Standard_D4_v2",
					"Standard_A4_v2",
					"Standard_A8_v2",
				}, false),
			},

			"number_of_nodes": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 10),
			},

			"max_parallel_executions_per_node": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 16),
			},

			"edition": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(datafactory.IntegrationRuntimeEditionStandard),
				ValidateFunc: validation.StringInSlice([]string{
					string(datafactory.IntegrationRuntimeEditionStandard),
					string(datafactory.IntegrationRuntimeEditionEnterprise),
				}, false),
			},

			"license_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(datafactory.IntegrationRuntimeLicenseTypeLicenseIncluded),
				ValidateFunc: validation.StringInSlice([]string{
					string(datafactory.IntegrationRuntimeLicenseTypeLicenseIncluded),
					string(datafactory.IntegrationRuntimeLicenseTypeBasePrice),
				}, false),
			},

			"vnet_integration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"vnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"subnet_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"public_ips": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MinItems: 2,
							MaxItems: 2,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: networkValidate.PublicIpAddressID,
							},
						},
					},
				},
			},

			"custom_setup_script": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"blob_container_uri": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"sas_token": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"catalog_info": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"server_endpoint": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"administrator_login": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"administrator_password": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"pricing_tier": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(datafactory.IntegrationRuntimeSsisCatalogPricingTierBasic),
							ValidateFunc: validation.StringInSlice([]string{
								string(datafactory.IntegrationRuntimeSsisCatalogPricingTierBasic),
								string(datafactory.IntegrationRuntimeSsisCatalogPricingTierStandard),
								string(datafactory.IntegrationRuntimeSsisCatalogPricingTierPremium),
								string(datafactory.IntegrationRuntimeSsisCatalogPricingTierPremiumRS),
							}, false),
						},
						"dual_standby_pair_name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"express_custom_setup": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"environment": {
							Type:         pluginsdk.TypeMap,
							Optional:     true,
							AtLeastOneOf: []string{"express_custom_setup.0.environment", "express_custom_setup.0.powershell_version", "express_custom_setup.0.component", "express_custom_setup.0.command_key"},
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"powershell_version": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							AtLeastOneOf: []string{"express_custom_setup.0.environment", "express_custom_setup.0.powershell_version", "express_custom_setup.0.component", "express_custom_setup.0.command_key"},
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"command_key": {
							Type:         pluginsdk.TypeList,
							Optional:     true,
							AtLeastOneOf: []string{"express_custom_setup.0.environment", "express_custom_setup.0.powershell_version", "express_custom_setup.0.component", "express_custom_setup.0.command_key"},
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"target_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"user_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"password": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										Sensitive:    true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},

						"component": {
							Type:         pluginsdk.TypeList,
							Optional:     true,
							AtLeastOneOf: []string{"express_custom_setup.0.environment", "express_custom_setup.0.powershell_version", "express_custom_setup.0.component", "express_custom_setup.0.command_key"},
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"license": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Sensitive:    true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},

			"package_store": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"linked_service_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"proxy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"self_hosted_integration_runtime_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"staging_storage_linked_service_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"path": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
	}
}

func resourceDataFactoryIntegrationRuntimeAzureSsisCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	factoryName := d.Get("data_factory_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, factoryName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Data Factory Azure-SSIS Integration Runtime %q (Resource Group %q, Data Factory %q): %s", name, resourceGroup, factoryName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_factory_integration_runtime_azure_ssis", *existing.ID)
		}
	}

	description := d.Get("description").(string)
	managedIntegrationRuntime := datafactory.ManagedIntegrationRuntime{
		Description: &description,
		Type:        datafactory.TypeBasicIntegrationRuntimeTypeManaged,
		ManagedIntegrationRuntimeTypeProperties: &datafactory.ManagedIntegrationRuntimeTypeProperties{
			ComputeProperties: expandDataFactoryIntegrationRuntimeAzureSsisComputeProperties(d),
			SsisProperties:    expandDataFactoryIntegrationRuntimeAzureSsisProperties(d),
		},
	}

	basicIntegrationRuntime, _ := managedIntegrationRuntime.AsBasicIntegrationRuntime()

	integrationRuntime := datafactory.IntegrationRuntimeResource{
		Name:       &name,
		Properties: basicIntegrationRuntime,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, factoryName, name, integrationRuntime, ""); err != nil {
		return fmt.Errorf("Error creating/updating Data Factory Azure-SSIS Integration Runtime %q (Resource Group %q, Data Factory %q): %+v", name, resourceGroup, factoryName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, factoryName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Data Factory Azure-SSIS Integration Runtime %q (Resource Group %q, Data Factory %q): %+v", name, resourceGroup, factoryName, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Data Factory Azure-SSIS Integration Runtime %q (Resource Group %q, Data Factory %q) ID", name, resourceGroup, factoryName)
	}

	d.SetId(*resp.ID)

	return resourceDataFactoryIntegrationRuntimeAzureSsisRead(d, meta)
}

func resourceDataFactoryIntegrationRuntimeAzureSsisRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	factoryName := id.Path["factories"]
	name := id.Path["integrationruntimes"]

	resp, err := client.Get(ctx, resourceGroup, factoryName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Data Factory Azure-SSIS Integration Runtime %q (Resource Group %q, Data Factory %q): %+v", name, resourceGroup, factoryName, err)
	}

	d.Set("name", name)
	d.Set("data_factory_name", factoryName)
	d.Set("resource_group_name", resourceGroup)

	managedIntegrationRuntime, convertSuccess := resp.Properties.AsManagedIntegrationRuntime()
	if !convertSuccess {
		return fmt.Errorf("Error converting integration runtime to Azure-SSIS integration runtime %q (Resource Group %q, Data Factory %q)", name, resourceGroup, factoryName)
	}

	if managedIntegrationRuntime.Description != nil {
		d.Set("description", managedIntegrationRuntime.Description)
	}

	if computeProps := managedIntegrationRuntime.ComputeProperties; computeProps != nil {
		if location := computeProps.Location; location != nil {
			d.Set("location", location)
		}

		if nodeSize := computeProps.NodeSize; nodeSize != nil {
			d.Set("node_size", nodeSize)
		}

		if numberOfNodes := computeProps.NumberOfNodes; numberOfNodes != nil {
			d.Set("number_of_nodes", numberOfNodes)
		}

		if maxParallelExecutionsPerNode := computeProps.MaxParallelExecutionsPerNode; maxParallelExecutionsPerNode != nil {
			d.Set("max_parallel_executions_per_node", maxParallelExecutionsPerNode)
		}

		if err := d.Set("vnet_integration", flattenDataFactoryIntegrationRuntimeAzureSsisVnetIntegration(computeProps.VNetProperties)); err != nil {
			return fmt.Errorf("Error setting `vnet_integration`: %+v", err)
		}
	}

	if ssisProps := managedIntegrationRuntime.SsisProperties; ssisProps != nil {
		d.Set("edition", string(ssisProps.Edition))
		d.Set("license_type", string(ssisProps.LicenseType))

		if err := d.Set("catalog_info", flattenDataFactoryIntegrationRuntimeAzureSsisCatalogInfo(ssisProps.CatalogInfo, d)); err != nil {
			return fmt.Errorf("setting `catalog_info`: %+v", err)
		}

		if err := d.Set("custom_setup_script", flattenDataFactoryIntegrationRuntimeAzureSsisCustomSetupScript(ssisProps.CustomSetupScriptProperties, d)); err != nil {
			return fmt.Errorf("setting `custom_setup_script`: %+v", err)
		}

		if err := d.Set("express_custom_setup", flattenDataFactoryIntegrationRuntimeAzureSsisExpressCustomSetUp(ssisProps.ExpressCustomSetupProperties, d)); err != nil {
			return fmt.Errorf("setting `express_custom_setup`: %+v", err)
		}

		if err := d.Set("package_store", flattenDataFactoryIntegrationRuntimeAzureSsisPackageStore(ssisProps.PackageStores)); err != nil {
			return fmt.Errorf("setting `package_store`: %+v", err)
		}

		if err := d.Set("proxy", flattenDataFactoryIntegrationRuntimeAzureSsisProxy(ssisProps.DataProxyProperties)); err != nil {
			return fmt.Errorf("setting `proxy`: %+v", err)
		}
	}

	return nil
}

func resourceDataFactoryIntegrationRuntimeAzureSsisDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	factoryName := id.Path["factories"]
	name := id.Path["integrationruntimes"]

	response, err := client.Delete(ctx, resourceGroup, factoryName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("Error deleting Data Factory Azure-SSIS Integration Runtime %q (Resource Group %q, Data Factory %q): %+v", name, resourceGroup, factoryName, err)
		}
	}

	return nil
}

func expandDataFactoryIntegrationRuntimeAzureSsisComputeProperties(d *pluginsdk.ResourceData) *datafactory.IntegrationRuntimeComputeProperties {
	location := azure.NormalizeLocation(d.Get("location").(string))
	computeProperties := datafactory.IntegrationRuntimeComputeProperties{
		Location:                     &location,
		NodeSize:                     utils.String(d.Get("node_size").(string)),
		NumberOfNodes:                utils.Int32(int32(d.Get("number_of_nodes").(int))),
		MaxParallelExecutionsPerNode: utils.Int32(int32(d.Get("max_parallel_executions_per_node").(int))),
	}

	if vnetIntegrations, ok := d.GetOk("vnet_integration"); ok && len(vnetIntegrations.([]interface{})) > 0 {
		vnetProps := vnetIntegrations.([]interface{})[0].(map[string]interface{})
		computeProperties.VNetProperties = &datafactory.IntegrationRuntimeVNetProperties{
			VNetID: utils.String(vnetProps["vnet_id"].(string)),
			Subnet: utils.String(vnetProps["subnet_name"].(string)),
		}

		if publicIPs := vnetProps["public_ips"].([]interface{}); len(publicIPs) > 0 {
			computeProperties.VNetProperties.PublicIPs = utils.ExpandStringSlice(publicIPs)
		}
	}

	return &computeProperties
}

func expandDataFactoryIntegrationRuntimeAzureSsisProperties(d *pluginsdk.ResourceData) *datafactory.IntegrationRuntimeSsisProperties {
	ssisProperties := &datafactory.IntegrationRuntimeSsisProperties{
		LicenseType:                  datafactory.IntegrationRuntimeLicenseType(d.Get("license_type").(string)),
		DataProxyProperties:          expandDataFactoryIntegrationRuntimeAzureSsisProxy(d.Get("proxy").([]interface{})),
		Edition:                      datafactory.IntegrationRuntimeEdition(d.Get("edition").(string)),
		ExpressCustomSetupProperties: expandDataFactoryIntegrationRuntimeAzureSsisExpressCustomSetUp(d.Get("express_custom_setup").([]interface{})),
		PackageStores:                expandDataFactoryIntegrationRuntimeAzureSsisPackageStore(d.Get("package_store").([]interface{})),
	}

	if catalogInfos, ok := d.GetOk("catalog_info"); ok && len(catalogInfos.([]interface{})) > 0 {
		catalogInfo := catalogInfos.([]interface{})[0].(map[string]interface{})

		ssisProperties.CatalogInfo = &datafactory.IntegrationRuntimeSsisCatalogInfo{
			CatalogServerEndpoint: utils.String(catalogInfo["server_endpoint"].(string)),
			CatalogPricingTier:    datafactory.IntegrationRuntimeSsisCatalogPricingTier(catalogInfo["pricing_tier"].(string)),
		}

		if adminUserName := catalogInfo["administrator_login"]; adminUserName.(string) != "" {
			ssisProperties.CatalogInfo.CatalogAdminUserName = utils.String(adminUserName.(string))
		}

		if adminPassword := catalogInfo["administrator_password"]; adminPassword.(string) != "" {
			ssisProperties.CatalogInfo.CatalogAdminPassword = &datafactory.SecureString{
				Value: utils.String(adminPassword.(string)),
				Type:  datafactory.TypeSecureString,
			}
		}

		if dualStandbyPairName := catalogInfo["dual_standby_pair_name"].(string); dualStandbyPairName != "" {
			ssisProperties.CatalogInfo.DualStandbyPairName = utils.String(dualStandbyPairName)
		}
	}

	if customSetupScripts, ok := d.GetOk("custom_setup_script"); ok && len(customSetupScripts.([]interface{})) > 0 {
		customSetupScript := customSetupScripts.([]interface{})[0].(map[string]interface{})

		sasToken := &datafactory.SecureString{
			Value: utils.String(customSetupScript["sas_token"].(string)),
			Type:  datafactory.TypeSecureString,
		}

		ssisProperties.CustomSetupScriptProperties = &datafactory.IntegrationRuntimeCustomSetupScriptProperties{
			BlobContainerURI: utils.String(customSetupScript["blob_container_uri"].(string)),
			SasToken:         sasToken,
		}
	}

	return ssisProperties
}

func expandDataFactoryIntegrationRuntimeAzureSsisProxy(input []interface{}) *datafactory.IntegrationRuntimeDataProxyProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	raw := input[0].(map[string]interface{})

	result := &datafactory.IntegrationRuntimeDataProxyProperties{
		ConnectVia: &datafactory.EntityReference{
			Type:          datafactory.IntegrationRuntimeEntityReferenceTypeIntegrationRuntimeReference,
			ReferenceName: utils.String(raw["self_hosted_integration_runtime_name"].(string)),
		},
		StagingLinkedService: &datafactory.EntityReference{
			Type:          datafactory.IntegrationRuntimeEntityReferenceTypeLinkedServiceReference,
			ReferenceName: utils.String(raw["staging_storage_linked_service_name"].(string)),
		},
	}
	if path := raw["path"].(string); len(path) > 0 {
		result.Path = utils.String(path)
	}
	return result
}

func expandDataFactoryIntegrationRuntimeAzureSsisExpressCustomSetUp(input []interface{}) *[]datafactory.BasicCustomSetupBase {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	raw := input[0].(map[string]interface{})

	result := make([]datafactory.BasicCustomSetupBase, 0)
	if env := raw["environment"].(map[string]interface{}); len(env) > 0 {
		for k, v := range env {
			result = append(result, &datafactory.EnvironmentVariableSetup{
				Type: datafactory.TypeBasicCustomSetupBaseTypeEnvironmentVariableSetup,
				EnvironmentVariableSetupTypeProperties: &datafactory.EnvironmentVariableSetupTypeProperties{
					VariableName:  utils.String(k),
					VariableValue: utils.String(v.(string)),
				},
			})
		}
	}
	if powershellVersion := raw["powershell_version"].(string); powershellVersion != "" {
		result = append(result, &datafactory.AzPowerShellSetup{
			Type: datafactory.TypeBasicCustomSetupBaseTypeAzPowerShellSetup,
			AzPowerShellSetupTypeProperties: &datafactory.AzPowerShellSetupTypeProperties{
				Version: utils.String(powershellVersion),
			},
		})
	}
	if components := raw["component"].([]interface{}); len(components) > 0 {
		for _, item := range components {
			raw := item.(map[string]interface{})
			component := &datafactory.ComponentSetup{
				Type: datafactory.TypeBasicCustomSetupBaseTypeComponentSetup,
				LicensedComponentSetupTypeProperties: &datafactory.LicensedComponentSetupTypeProperties{
					ComponentName: utils.String(raw["name"].(string)),
				},
			}
			if license := raw["license"].(string); license != "" {
				component.LicensedComponentSetupTypeProperties.LicenseKey = &datafactory.SecureString{
					Type:  datafactory.TypeSecureString,
					Value: utils.String(license),
				}
			}

			result = append(result, component)
		}
	}
	if cmdKeys := raw["command_key"].([]interface{}); len(cmdKeys) > 0 {
		for _, item := range cmdKeys {
			raw := item.(map[string]interface{})
			result = append(result, &datafactory.CmdkeySetup{
				Type: datafactory.TypeBasicCustomSetupBaseTypeCmdkeySetup,
				CmdkeySetupTypeProperties: &datafactory.CmdkeySetupTypeProperties{
					TargetName: utils.String(raw["target_name"].(string)),
					UserName:   utils.String(raw["user_name"].(string)),
					Password: &datafactory.SecureString{
						Type:  datafactory.TypeSecureString,
						Value: utils.String(raw["password"].(string)),
					},
				},
			})
		}
	}

	return &result
}

func expandDataFactoryIntegrationRuntimeAzureSsisPackageStore(input []interface{}) *[]datafactory.PackageStore {
	if len(input) == 0 {
		return nil
	}

	result := make([]datafactory.PackageStore, 0)
	for _, item := range input {
		raw := item.(map[string]interface{})
		result = append(result, datafactory.PackageStore{
			Name: utils.String(raw["name"].(string)),
			PackageStoreLinkedService: &datafactory.EntityReference{
				Type:          datafactory.IntegrationRuntimeEntityReferenceTypeLinkedServiceReference,
				ReferenceName: utils.String(raw["linked_service_name"].(string)),
			},
		})
	}
	return &result
}

func flattenDataFactoryIntegrationRuntimeAzureSsisVnetIntegration(vnetProperties *datafactory.IntegrationRuntimeVNetProperties) []interface{} {
	if vnetProperties == nil {
		return []interface{}{}
	}

	var vnetId, subnetName string
	if vnetProperties.VNetID != nil {
		vnetId = *vnetProperties.VNetID
	}
	if vnetProperties.Subnet != nil {
		subnetName = *vnetProperties.Subnet
	}

	return []interface{}{
		map[string]interface{}{
			"vnet_id":     vnetId,
			"subnet_name": subnetName,
			"public_ips":  utils.FlattenStringSlice(vnetProperties.PublicIPs),
		},
	}
}

func flattenDataFactoryIntegrationRuntimeAzureSsisCatalogInfo(ssisProperties *datafactory.IntegrationRuntimeSsisCatalogInfo, d *pluginsdk.ResourceData) []interface{} {
	if ssisProperties == nil {
		return []interface{}{}
	}

	var serverEndpoint, catalogAdminUserName, administratorPassword, dualStandbyPairName string
	if ssisProperties.CatalogServerEndpoint != nil {
		serverEndpoint = *ssisProperties.CatalogServerEndpoint
	}
	if ssisProperties.CatalogAdminUserName != nil {
		catalogAdminUserName = *ssisProperties.CatalogAdminUserName
	}
	if ssisProperties.DualStandbyPairName != nil {
		dualStandbyPairName = *ssisProperties.DualStandbyPairName
	}

	// read back
	if adminPassword, ok := d.GetOk("catalog_info.0.administrator_password"); ok {
		administratorPassword = adminPassword.(string)
	}

	return []interface{}{
		map[string]interface{}{
			"server_endpoint":        serverEndpoint,
			"pricing_tier":           string(ssisProperties.CatalogPricingTier),
			"administrator_login":    catalogAdminUserName,
			"administrator_password": administratorPassword,
			"dual_standby_pair_name": dualStandbyPairName,
		},
	}
}

func flattenDataFactoryIntegrationRuntimeAzureSsisProxy(input *datafactory.IntegrationRuntimeDataProxyProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var path, selfHostedIntegrationRuntimeName, stagingStorageLinkedServiceName string
	if input.Path != nil {
		path = *input.Path
	}
	if input.ConnectVia != nil && input.ConnectVia.ReferenceName != nil {
		selfHostedIntegrationRuntimeName = *input.ConnectVia.ReferenceName
	}
	if input.StagingLinkedService != nil && input.StagingLinkedService.ReferenceName != nil {
		stagingStorageLinkedServiceName = *input.StagingLinkedService.ReferenceName
	}
	return []interface{}{
		map[string]interface{}{
			"path":                                 path,
			"self_hosted_integration_runtime_name": selfHostedIntegrationRuntimeName,
			"staging_storage_linked_service_name":  stagingStorageLinkedServiceName,
		},
	}
}

func flattenDataFactoryIntegrationRuntimeAzureSsisCustomSetupScript(customSetupScriptProperties *datafactory.IntegrationRuntimeCustomSetupScriptProperties, d *pluginsdk.ResourceData) []interface{} {
	if customSetupScriptProperties == nil {
		return []interface{}{}
	}

	customSetupScript := map[string]string{
		"blob_container_uri": *customSetupScriptProperties.BlobContainerURI,
	}

	if sasToken, ok := d.GetOk("custom_setup_script.0.sas_token"); ok {
		customSetupScript["sas_token"] = sasToken.(string)
	}

	return []interface{}{customSetupScript}
}

func flattenDataFactoryIntegrationRuntimeAzureSsisPackageStore(input *[]datafactory.PackageStore) []interface{} {
	if input == nil {
		return nil
	}

	result := make([]interface{}, 0)
	for _, item := range *input {
		var name, linkedServiceName string
		if item.Name != nil {
			name = *item.Name
		}
		if item.PackageStoreLinkedService != nil && item.PackageStoreLinkedService.ReferenceName != nil {
			linkedServiceName = *item.PackageStoreLinkedService.ReferenceName
		}

		result = append(result, map[string]interface{}{
			"name":                name,
			"linked_service_name": linkedServiceName,
		})
	}
	return result
}

func flattenDataFactoryIntegrationRuntimeAzureSsisExpressCustomSetUp(input *[]datafactory.BasicCustomSetupBase, d *pluginsdk.ResourceData) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	// retrieve old state
	oldState := make(map[string]interface{})
	if arr := d.Get("express_custom_setup").([]interface{}); len(arr) > 0 {
		oldState = arr[0].(map[string]interface{})
	}
	oldComponents := make([]interface{}, 0)
	if rawComponent, ok := oldState["component"]; ok {
		if v := rawComponent.([]interface{}); len(v) > 0 {
			oldComponents = v
		}
	}
	oldCmdKey := make([]interface{}, 0)
	if rawCmdKey, ok := oldState["command_key"]; ok {
		if v := rawCmdKey.([]interface{}); len(v) > 0 {
			oldCmdKey = v
		}
	}

	env := make(map[string]interface{})
	powershellVersion := ""
	components := make([]interface{}, 0)
	cmdkeys := make([]interface{}, 0)
	for _, item := range *input {
		switch v := item.(type) {
		case datafactory.AzPowerShellSetup:
			if v.Version != nil {
				powershellVersion = *v.Version
			}
		case datafactory.ComponentSetup:
			var name string
			if v.ComponentName != nil {
				name = *v.ComponentName
			}
			components = append(components, map[string]interface{}{
				"name": name,
				"license": readBackSensitiveValue(oldComponents, "license", map[string]string{
					"name": name,
				}),
			})
		case datafactory.EnvironmentVariableSetup:
			if v.VariableName != nil && v.VariableValue != nil {
				env[*v.VariableName] = *v.VariableValue
			}
		case datafactory.CmdkeySetup:
			var name, userName string
			if v.TargetName != nil {
				if v, ok := v.TargetName.(string); ok {
					name = v
				}
			}
			if v.UserName != nil {
				if v, ok := v.UserName.(string); ok {
					userName = v
				}
			}
			cmdkeys = append(cmdkeys, map[string]interface{}{
				"target_name": name,
				"user_name":   userName,
				"password": readBackSensitiveValue(oldCmdKey, "password", map[string]string{
					"target_name": name,
					"user_name":   userName,
				}),
			})
		}
	}

	return []interface{}{
		map[string]interface{}{
			"environment":        env,
			"powershell_version": powershellVersion,
			"component":          components,
			"command_key":        cmdkeys,
		},
	}
}

func readBackSensitiveValue(input []interface{}, propertyName string, filters map[string]string) string {
	if len(input) == 0 {
		return ""
	}
	for _, item := range input {
		raw := item.(map[string]interface{})
		found := true
		for k, v := range filters {
			if raw[k].(string) != v {
				found = false
				break
			}
		}
		if found {
			return raw[propertyName].(string)
		}
	}
	return ""
}
