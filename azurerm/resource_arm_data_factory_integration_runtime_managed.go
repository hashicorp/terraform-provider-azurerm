package azurerm

import (
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func resourceArmDataFactoryIntegrationRuntimeManaged() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDataFactoryIntegrationRuntimeManagedCreate,
		Read:   resourceArmDataFactoryIntegrationRuntimeManagedRead,
		Update: resourceArmDataFactoryIntegrationRuntimeManagedUpdate,
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

			"resource_group": azure.SchemaResourceGroupName(),

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
							Type:     schema.TypeString,
							Required: true,
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
							Type:     schema.TypeString,
							Required: true,
						},
						"sas_token": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
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
							Type:     schema.TypeString,
							Required: true,
						},
						"administrator_login": {
							Type:     schema.TypeString,
							Required: true,
						},
						"administrator_password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
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

func resourceArmDataFactoryIntegrationRuntimeManagedCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceArmDataFactoryIntegrationRuntimeManagedRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceArmDataFactoryIntegrationRuntimeManagedUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceArmDataFactoryIntegrationRuntimeManagedDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
