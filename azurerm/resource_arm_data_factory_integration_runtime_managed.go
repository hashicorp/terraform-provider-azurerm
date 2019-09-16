package azurerm

import (
	"fmt"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataFactoryIntegrationRuntimeManaged() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDataFactoryIntegrationRuntimeManagedCreateUpdate,
		Read:   resourceArmDataFactoryIntegrationRuntimeManagedRead,
		Update: resourceArmDataFactoryIntegrationRuntimeManagedCreateUpdate,
		Delete: resourceArmDataFactoryIntegrationRuntimeManagedDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^([a-zA-Z0-9](-|-?[a-zA-Z0-9]+)+[a-zA-Z0-9])$`),
					`Invalid name for Managed Integration Runtime: minimum 3 characters, must start and end with a number or a letter, may only consist of letters, numbers and dashes and no consecutive dashes.`,
				),
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"data_factory_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`),
					`Invalid name for Data Factory, see https://docs.microsoft.com/en-us/azure/data-factory/naming-rules`,
				),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"node_size": {
				Type:     schema.TypeString,
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
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 10),
			},

			"max_parallel_executions_per_node": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 8),
			},

			"edition": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(datafactory.Standard),
				ValidateFunc: validation.StringInSlice([]string{
					string(datafactory.Standard),
					string(datafactory.Enterprise),
				}, false),
			},

			"license_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(datafactory.LicenseIncluded),
				ValidateFunc: validation.StringInSlice([]string{
					string(datafactory.LicenseIncluded),
					string(datafactory.BasePrice),
				}, false),
			},

			"vnet_integration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vnet_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"subnet_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"custom_setup_script": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blob_container_uri": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"sas_token": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"catalog_info": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_endpoint": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"administrator_login": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"administrator_password": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"pricing_tier": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(datafactory.IntegrationRuntimeSsisCatalogPricingTierBasic),
							ValidateFunc: validation.StringInSlice([]string{
								string(datafactory.IntegrationRuntimeSsisCatalogPricingTierBasic),
								string(datafactory.IntegrationRuntimeSsisCatalogPricingTierStandard),
								string(datafactory.IntegrationRuntimeSsisCatalogPricingTierPremium),
								string(datafactory.IntegrationRuntimeSsisCatalogPricingTierPremiumRS),
							}, false),
						},
					},
				},
			},
		},
	}
}

func resourceArmDataFactoryIntegrationRuntimeManagedCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataFactory.IntegrationRuntimesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	factoryName := d.Get("data_factory_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, factoryName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Data Factory Managed Integration Runtime %q (Resource Group %q, Data Factory %q): %s", name, resourceGroup, factoryName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_factory_integration_runtime_managed", *existing.ID)
		}
	}

	description := d.Get("description").(string)
	managedIntegrationRuntime := datafactory.ManagedIntegrationRuntime{
		Description: &description,
		Type:        datafactory.TypeManaged,
	}

	managedIntegrationRuntime.ComputeProperties = expandArmDataFactoryIntegrationRuntimeManagedComputeProperties(d)
	managedIntegrationRuntime.SsisProperties = expandArmDataFactoryIntegrationRuntimeManagedSsisProperties(d)

	basicIntegrationRuntime, _ := managedIntegrationRuntime.AsBasicIntegrationRuntime()

	integrationRuntime := datafactory.IntegrationRuntimeResource{
		Name:       &name,
		Properties: basicIntegrationRuntime,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, factoryName, name, integrationRuntime, ""); err != nil {
		return fmt.Errorf("Error creating/updating Data Factory Managed Integration Runtime %q (Resource Group %q, Data Factory %q): %+v", name, resourceGroup, factoryName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, factoryName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Data Factory Managed Integration Runtime %q (Resource Group %q, Data Factory %q): %+v", name, resourceGroup, factoryName, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Data Factory Managed Integration Runtime %q (Resource Group %q, Data Factory %q) ID", name, resourceGroup, factoryName)
	}

	d.SetId(*resp.ID)

	return resourceArmDataFactoryRead(d, meta)
}

func resourceArmDataFactoryIntegrationRuntimeManagedRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataFactory.IntegrationRuntimesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	factoryName := id.Path["factories"]
	name := id.Path["integrationRuntimes"]

	resp, err := client.Get(ctx, resourceGroup, factoryName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Data Factory Managed Integration Runtime %q (Resource Group %q, Data Factory %q): %+v", name, resourceGroup, factoryName, err)
	}

	d.Set("name", resp.Name)
	d.Set("data_factory_name", factoryName)
	d.Set("resource_group_name", resourceGroup)

	managedIntegrationRuntime, convertErr := resp.Properties.AsManagedIntegrationRuntime()
	if convertErr {
		return fmt.Errorf("Error converting Data Factory Managed Integration Runtime basic properties to managed properties %q (Resource Group %q, Data Factory %q): %+v", name, resourceGroup, factoryName, err)
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

		if vnetProps := computeProps.VNetProperties; vnetProps != nil {
			d.Set("vnet_integration", flattenArmDataFactoryIntegrationRuntimeManagedVnetIntegration(vnetProps))
		}
	}

	if ssisProps := managedIntegrationRuntime.SsisProperties; ssisProps != nil {
		d.Set("edition", string(ssisProps.Edition))
		d.Set("license_type", string(ssisProps.LicenseType))

		if catalogInfoProps := ssisProps.CatalogInfo; catalogInfoProps != nil {
			d.Set("catalog_info", flattenArmDataFactoryIntegrationRuntimeManagedSsisCatalogInfo(catalogInfoProps))
		}

		if customSetupScriptProps := ssisProps.CustomSetupScriptProperties; customSetupScriptProps != nil {
			d.Set("custom_setup_script", flattenArmDataFactoryIntegrationRuntimeManagedSsisCustomSetupScript(customSetupScriptProps))
		}
	}

	return nil
}

func resourceArmDataFactoryIntegrationRuntimeManagedDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataFactory.IntegrationRuntimesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	factoryName := id.Path["factories"]
	name := id.Path["integrationRuntimes"]

	response, err := client.Delete(ctx, resourceGroup, factoryName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("Error deleting Data Factory Managed Integration Runtime %q (Resource Group %q, Data Factory %q): %+v", name, resourceGroup, factoryName, err)
		}
	}

	return nil
}

func expandArmDataFactoryIntegrationRuntimeManagedComputeProperties(d *schema.ResourceData) *datafactory.IntegrationRuntimeComputeProperties {
	location := azure.NormalizeLocation(d.Get("location").(string))
	computeProperties := datafactory.IntegrationRuntimeComputeProperties{
		Location:                     &location,
		NodeSize:                     utils.String(d.Get("node_size").(string)),
		NumberOfNodes:                utils.Int32(int32(d.Get("number_of_nodes").(int))),
		MaxParallelExecutionsPerNode: utils.Int32(int32(d.Get("max_parallel_executions_per_node").(int))),
	}

	if _, ok := d.GetOk("vnet_integration"); ok {
		vnetProps := d.Get("vnet_integration").([]map[string]interface{})[0]
		computeProperties.VNetProperties = &datafactory.IntegrationRuntimeVNetProperties{
			VNetID: utils.String(vnetProps["vnet_id"].(string)),
			Subnet: utils.String(vnetProps["subnet_name"].(string)),
		}
	}

	return &computeProperties
}

func expandArmDataFactoryIntegrationRuntimeManagedSsisProperties(d *schema.ResourceData) *datafactory.IntegrationRuntimeSsisProperties {
	ssisProperties := &datafactory.IntegrationRuntimeSsisProperties{
		Edition:     datafactory.IntegrationRuntimeEdition(d.Get("edition").(string)),
		LicenseType: datafactory.IntegrationRuntimeLicenseType(d.Get("license_type").(string)),
	}

	if _, ok := d.GetOk("catalog_info"); ok {
		catalogInfo := d.Get("catalog_info").([]map[string]interface{})[0]

		adminPassword := &datafactory.SecureString{
			Value: utils.String(catalogInfo["administrator_password"].(string)),
			Type:  datafactory.TypeSecureString,
		}

		ssisProperties.CatalogInfo = &datafactory.IntegrationRuntimeSsisCatalogInfo{
			CatalogServerEndpoint: utils.String(catalogInfo["server_endpoint"].(string)),
			CatalogAdminUserName:  utils.String(catalogInfo["administrator_login"].(string)),
			CatalogAdminPassword:  adminPassword,
			CatalogPricingTier:    datafactory.IntegrationRuntimeSsisCatalogPricingTier(catalogInfo["pricing_tier"].(string)),
		}
	}

	if _, ok := d.GetOk("custom_setup_script"); ok {
		customSetupScript := d.Get("custom_setup_script").([]map[string]interface{})[0]

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

func flattenArmDataFactoryIntegrationRuntimeManagedVnetIntegration(vnetProperties *datafactory.IntegrationRuntimeVNetProperties) []interface{} {
	return []interface{}{
		map[string]string{
			"vnet_id":     *vnetProperties.VNetID,
			"subnet_name": *vnetProperties.Subnet,
		},
	}
}

func flattenArmDataFactoryIntegrationRuntimeManagedSsisCatalogInfo(ssisProperties *datafactory.IntegrationRuntimeSsisCatalogInfo) []interface{} {
	return []interface{}{
		map[string]string{
			"server_endpoint":     *ssisProperties.CatalogServerEndpoint,
			"administrator_login": *ssisProperties.CatalogAdminUserName,
			"pricing_tier":        string(ssisProperties.CatalogPricingTier),
		},
	}
}

func flattenArmDataFactoryIntegrationRuntimeManagedSsisCustomSetupScript(customSetupScriptProperties *datafactory.IntegrationRuntimeCustomSetupScriptProperties) []interface{} {
	return []interface{}{
		map[string]string{
			"blob_container_uri": *customSetupScriptProperties.BlobContainerURI,
		},
	}
}
