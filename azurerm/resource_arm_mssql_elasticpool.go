package azurerm

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-10-01-preview/sql"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMsSqlElasticPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMsSqlElasticPoolCreateUpdate,
		Read:   resourceArmMsSqlElasticPoolRead,
		Update: resourceArmMsSqlElasticPoolCreateUpdate,
		Delete: resourceArmMsSqlElasticPoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateMsSqlElasticPoolName,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"server_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateMsSqlServerName,
			},

			"sku": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"BasicPool",
								"StandardPool",
								"PremiumPool",
								"GP_Gen4",
								"GP_Gen5",
								"BC_Gen4",
								"BC_Gen5",
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"capacity": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"tier": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Basic",
								"Standard",
								"Premium",
								"GeneralPurpose",
								"BusinessCritical",
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"family": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Gen4",
								"Gen5",
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},
				},
			},

			"per_database_settings": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_capacity": {
							Type:         schema.TypeFloat,
							Required:     true,
							ValidateFunc: validate.FloatAtLeast(0.0),
						},

						"max_capacity": {
							Type:         schema.TypeFloat,
							Required:     true,
							ValidateFunc: validate.FloatAtLeast(0.0),
						},
					},
				},
			},

			"elastic_pool_properties": {
				Type:       schema.TypeList,
				Computed:   true,
				MaxItems:   1,
				Deprecated: "These properties herein have been moved to the top level or removed",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "This property has been removed",
						},

						"creation_date": {
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "This property has been removed",
						},

						"max_size_bytes": {
							Type:       schema.TypeInt,
							Computed:   true,
							Deprecated: "This property has been moved to the top level",
						},

						"zone_redundant": {
							Type:       schema.TypeBool,
							Computed:   true,
							Deprecated: "This property has been moved to the top level",
						},

						"license_type": {
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "This property has been removed",
						},
					},
				},
			},

			"max_size_gb": {
				Type:         schema.TypeFloat,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.FloatAtLeast(0),
			},

			"zone_redundant": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"tags": tagsSchema(),
		},

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {

			name, _ := diff.GetOk("sku.0.name")
			tier, _ := diff.GetOk("sku.0.tier")
			capacity, _ := diff.GetOk("sku.0.capacity")
			family, _ := diff.GetOk("sku.0.family")
			maxSizeGb, _ := diff.GetOk("max_size_gb")
			minCapacity, _ := diff.GetOk("per_database_settings.0.min_capacity")
			maxCapacity, _ := diff.GetOk("per_database_settings.0.max_capacity")

			// Basic Checks
			if strings.ToLower(tier.(string)) == "basic" {
				if !azure.BasicIsCapacityValid(capacity.(int)) {
					return fmt.Errorf("service tier 'Basic' must have a capacity of 50, 100, 200, 300, 400, 800, 1200 or 1600 DTUs")
				}

				// Basic SKU does not let you pick your max_size_GB they are fixed values
				maxAllowedGB := azure.BasicGetMaxSizeGB(capacity.(int))

				if maxSizeGb.(float64) != maxAllowedGB {
					return fmt.Errorf("service tier 'Basic' with a capacity of %d must have a max_size_gb of %.7f GB, got %.7f GB", capacity.(int), maxAllowedGB, maxSizeGb.(float64))
				}
			}

			// Standard Checks
			if strings.ToLower(tier.(string)) == "standard" {
				if !azure.StandardCapacityValid(capacity.(int)) {
					return fmt.Errorf("service tier 'Standard' must have a capacity of 50, 100, 200, 300, 400, 800, 1200, 1600, 2000, 2500 or 3000 DTUs")
				}

				maxAllowedGB := azure.StandardGetMaxSizeGB(capacity.(int))

				if maxSizeGb.(float64) > maxAllowedGB {
					return fmt.Errorf("service tier 'Standard' with a capacity of %d must have a max_size_gb no greater than %d GB, got %d GB", capacity.(int), int(maxAllowedGB), int(maxSizeGb.(float64)))
				}
			}

			// Premium Checks
			if strings.ToLower(tier.(string)) == "premium" {
				if !azure.PremiumCapacityValid(capacity.(int)) {
					return fmt.Errorf("service tier 'Premium' must have a capacity of 125, 250, 500, 1000, 1500, 2000, 2500, 3000, 3500 or 4000 DTUs")
				}

				maxAllowedGB := azure.PremiumGetMaxSizeGB(capacity.(int))

				if maxSizeGb.(float64) > maxAllowedGB {
					return fmt.Errorf("service tier 'Premium' with a capacity of %d must have a max_size_gb no greater than %d GB, got %d GB", capacity.(int), int(maxAllowedGB), int(maxSizeGb.(float64)))
				}
			}

			// GeneralPurpose Checks
			if strings.HasPrefix(strings.ToLower(name.(string)), "gp_") {
				// Gen4 Checks
				if strings.ToLower(family.(string)) == "gen4" {
					if !azure.GeneralPurposeCapacityValid(capacity.(int), "gen4") {
						return fmt.Errorf("service tier 'GeneralPurpose' Gen4 must have a capacity of 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 16 or 24 vCores")
					}

					maxAllowedGB := azure.GeneralPurposeGetMaxSizeGB(capacity.(int), "gen4")

					if maxSizeGb.(float64) > maxAllowedGB {
						return fmt.Errorf("service tier 'GeneralPurpose' Gen4 with a capacity of %d vCores must have a max_size_gb between 5 GB and %d GB, got %d GB", capacity.(int), int(maxAllowedGB), int(maxSizeGb.(float64)))
					}
				}

				// Gen5 Checks
				if strings.ToLower(family.(string)) == "gen5" {
					if !azure.GeneralPurposeCapacityValid(capacity.(int), "gen5") {
						return fmt.Errorf("service tier 'GeneralPurpose' Gen5 must have a capacity of 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 24, 32, 40 or 80 vCores")
					}

					maxAllowedGB := azure.GeneralPurposeGetMaxSizeGB(capacity.(int), "gen5")

					if maxSizeGb.(float64) > maxAllowedGB {
						return fmt.Errorf("service tier 'GeneralPurpose' Gen5 with a capacity of %d vCores must have a max_size_gb between 5 GB and %d GB, got %d GB", capacity.(int), int(maxAllowedGB), int(maxSizeGb.(float64)))
					}
				}
			}

			// BusinessCritical Checks
			if strings.HasPrefix(strings.ToLower(name.(string)), "bc_") {
				// Gen4 Checks
				if strings.ToLower(family.(string)) == "gen4" {
					if !azure.BusinessCriticalCapacityValid(capacity.(int), "gen4") {
						return fmt.Errorf("service tier 'BusinessCritical' Gen4 must have a capacity of 2, 3, 4, 5, 6, 7, 8, 9, 10, 16 or 24 vCores")
					}

					maxAllowedGB := azure.BusinessCriticalGetMaxSizeGB(capacity.(int), "gen4")

					if maxSizeGb.(float64) > maxAllowedGB {
						return fmt.Errorf("service tier 'BusinessCritical' Gen4 with a capacity of %d vCores must have a max_size_gb between 5 GB and %d GB, got %d GB", capacity.(int), int(maxAllowedGB), int(maxSizeGb.(float64)))
					}
				}

				// Gen5 Checks
				if strings.ToLower(family.(string)) == "gen5" {
					if !azure.BusinessCriticalCapacityValid(capacity.(int), "gen5") {
						return fmt.Errorf("service tier 'BusinessCritical' Gen5 must have a capacity of 4, 6, 8, 10, 12, 14, 16, 18, 20, 24, 32, 40 or 80 vCores")
					}

					maxAllowedGB := azure.BusinessCriticalGetMaxSizeGB(capacity.(int), "gen5")

					if maxSizeGb.(float64) > maxAllowedGB {
						return fmt.Errorf("service tier 'BusinessCritical' Gen5 with a capacity of %d vCores must have a max_size_gb between 5 GB and %d GB, got %d GB", capacity.(int), int(maxAllowedGB), int(maxSizeGb.(float64)))
					}
				}
			}

			// General Checks based off SKU type...
			if strings.HasPrefix(strings.ToLower(name.(string)), "gp_") || strings.HasPrefix(strings.ToLower(name.(string)), "bc_") {
				// vCore Based
				if int(maxSizeGb.(float64)) < 5 {
					return fmt.Errorf("service tier 'GeneralPurpose' and 'BusinessCritical' must have a max_size_gb value equal to or greater than 5 GB, got %d GB", int(maxSizeGb.(float64)))
				}

				if err, ok := azure.IsMaxGBValid(1, maxSizeGb.(float64)); !ok {
					return fmt.Errorf(err)
				}

				if !azure.NameFamilyValid(name.(string), family.(string)) {
					return fmt.Errorf("SKU has a name family mismatch, got name = %s, family = %s", name.(string), family.(string))
				}

				if maxCapacity.(float64) > float64(capacity.(int)) {
					return fmt.Errorf("service tiers 'GeneralPurpose' and 'BusinessCritical' perDatabaseSettings maxCapacity must not be higher than the SKUs capacity value")
				}

				if minCapacity.(float64) > maxCapacity.(float64) {
					return fmt.Errorf("perDatabaseSettings maxCapacity must be greater than or equal to the perDatabaseSettings minCapacity value")
				}
			} else {
				// DTU Based
				if strings.ToLower(tier.(string)) != "basic" {
					if int(maxSizeGb.(float64)) < 50 {
						return fmt.Errorf("service tiers 'Standard', and 'Premium' must have a max_size_gb value equal to or greater than 50 GB, got %d GB", int(maxSizeGb.(float64)))
					}

					if err, ok := azure.IsMaxGBValid(50, maxSizeGb.(float64)); !ok {
						return fmt.Errorf(err)
					}
				}

				if maxCapacity.(float64) != math.Trunc(maxCapacity.(float64)) {
					return fmt.Errorf("service tiers 'Basic', 'Standard', and 'Premium' must have whole numbers as their maxCapacity")
				}

				if minCapacity.(float64) != math.Trunc(minCapacity.(float64)) {
					return fmt.Errorf("service tiers 'Basic', 'Standard', and 'Premium' must have whole numbers as their minCapacity")
				}

				if minCapacity.(float64) < 0.0 {
					return fmt.Errorf("service tiers 'Basic', 'Standard', and 'Premium' per_database_settings min_capacity must be equal to or greater than zero")
				}
			}

			return nil
		},
	}
}

func resourceArmMsSqlElasticPoolCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).msSqlElasticPoolsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for MsSQL ElasticPool creation.")

	elasticPoolName := d.Get("name").(string)
	serverName := d.Get("server_name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	sku := expandAzureRmMsSqlElasticPoolSku(d)
	tags := d.Get("tags").(map[string]interface{})

	elasticPool := sql.ElasticPool{
		Name:     &elasticPoolName,
		Location: &location,
		Sku:      sku,
		Tags:     expandTags(tags),
		ElasticPoolProperties: &sql.ElasticPoolProperties{
			PerDatabaseSettings: expandAzureRmMsSqlElasticPoolPerDatabaseSettings(d),
		},
	}

	if v, ok := d.GetOk("max_size_gb"); ok {
		maxSizeBytes := v.(float64) * 1073741824
		elasticPool.MaxSizeBytes = utils.Int64(int64(maxSizeBytes))
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, serverName, elasticPoolName, elasticPool)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, serverName, elasticPoolName)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read MsSQL ElasticPool %q (resource group %q) ID", elasticPoolName, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmMsSqlElasticPoolRead(d, meta)
}

func resourceArmMsSqlElasticPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).msSqlElasticPoolsClient
	ctx := meta.(*ArmClient).StopContext

	resGroup, serverName, name, err := parseArmMsSqlElasticPoolId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, resGroup, serverName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on MsSql Elastic Pool %s: %s", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	d.Set("server_name", serverName)

	if err := d.Set("sku", flattenAzureRmMsSqlElasticPoolSku(resp.Sku)); err != nil {
		return fmt.Errorf("Error setting `sku`: %+v", err)
	}

	if properties := resp.ElasticPoolProperties; properties != nil {
		d.Set("max_size_bytes", properties.MaxSizeBytes)
		d.Set("zone_redundant", properties.ZoneRedundant)

		//todo remove in 2.0
		if err := d.Set("elastic_pool_properties", flattenAzureRmMsSqlElasticPoolProperties(resp.ElasticPoolProperties)); err != nil {
			return fmt.Errorf("Error setting `elastic_pool_properties`: %+v", err)
		}

		if err := d.Set("per_database_settings", flattenAzureRmMsSqlElasticPoolPerDatabaseSettings(properties.PerDatabaseSettings)); err != nil {
			return fmt.Errorf("Error setting `per_database_settings`: %+v", err)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmMsSqlElasticPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).msSqlElasticPoolsClient
	ctx := meta.(*ArmClient).StopContext

	resGroup, serverName, name, err := parseArmSqlElasticPoolId(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, resGroup, serverName, name)
	return err
}

func parseArmMsSqlElasticPoolId(sqlElasticPoolId string) (string, string, string, error) {
	id, err := parseAzureResourceID(sqlElasticPoolId)
	if err != nil {
		return "", "", "", fmt.Errorf("[ERROR] Unable to parse MsSQL ElasticPool ID %q: %+v", sqlElasticPoolId, err)
	}

	return id.ResourceGroup, id.Path["servers"], id.Path["elasticPools"], nil
}

func expandAzureRmMsSqlElasticPoolPerDatabaseSettings(d *schema.ResourceData) *sql.ElasticPoolPerDatabaseSettings {
	perDatabaseSettings := d.Get("per_database_settings").([]interface{})
	perDatabaseSetting := perDatabaseSettings[0].(map[string]interface{})

	minCapacity := perDatabaseSetting["min_capacity"].(float64)
	maxCapacity := perDatabaseSetting["max_capacity"].(float64)

	return &sql.ElasticPoolPerDatabaseSettings{
		MinCapacity: utils.Float(minCapacity),
		MaxCapacity: utils.Float(maxCapacity),
	}
}

func expandAzureRmMsSqlElasticPoolSku(d *schema.ResourceData) *sql.Sku {
	skus := d.Get("sku").([]interface{})
	sku := skus[0].(map[string]interface{})

	name := sku["name"].(string)
	tier := sku["tier"].(string)
	family := sku["family"].(string)
	capacity := sku["capacity"].(int)

	return &sql.Sku{
		Name:     utils.String(name),
		Tier:     utils.String(tier),
		Family:   utils.String(family),
		Capacity: utils.Int32(int32(capacity)),
	}
}

func flattenAzureRmMsSqlElasticPoolSku(resp *sql.Sku) []interface{} {
	values := map[string]interface{}{}

	if name := resp.Name; name != nil {
		values["name"] = *name
	}

	if tier := resp.Tier; tier != nil {
		values["tier"] = *tier
	}

	if family := resp.Family; family != nil {
		values["family"] = *family
	}

	if capacity := resp.Capacity; capacity != nil {
		values["capacity"] = *capacity
	}

	return []interface{}{values}
}

func flattenAzureRmMsSqlElasticPoolProperties(resp *sql.ElasticPoolProperties) []interface{} {
	elasticPoolProperty := map[string]interface{}{}
	elasticPoolProperty["state"] = string(resp.State)

	if date := resp.CreationDate; date != nil {
		elasticPoolProperty["creation_date"] = date.String()
	}

	if zoneRedundant := resp.ZoneRedundant; zoneRedundant != nil {
		elasticPoolProperty["zone_redundant"] = *zoneRedundant
	}

	elasticPoolProperty["license_type"] = string(resp.LicenseType)

	return []interface{}{elasticPoolProperty}
}

func flattenAzureRmMsSqlElasticPoolPerDatabaseSettings(resp *sql.ElasticPoolPerDatabaseSettings) []interface{} {
	perDatabaseSettings := map[string]interface{}{}

	if minCapacity := resp.MinCapacity; minCapacity != nil {
		perDatabaseSettings["min_capacity"] = *minCapacity
	}

	if maxCapacity := resp.MaxCapacity; maxCapacity != nil {
		perDatabaseSettings["max_capacity"] = *maxCapacity
	}

	return []interface{}{perDatabaseSettings}
}
