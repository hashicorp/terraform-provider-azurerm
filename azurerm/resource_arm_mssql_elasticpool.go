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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMsSqlElasticPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMsSqlElasticPoolCreate,
		Read:   resourceArmMsSqlElasticPoolRead,
		Update: resourceArmMsSqlElasticPoolCreate,
		Delete: resourceArmMsSqlElasticPoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateServiceName,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"server_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateServiceName,
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
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"capacity": {
							Type:     schema.TypeInt,
							Required: true,
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
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"family": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Gen4",
								"Gen5",
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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
							Type:     schema.TypeFloat,
							Required: true,
						},

						"max_capacity": {
							Type:     schema.TypeFloat,
							Required: true,
						},
					},
				},
			},

			"elastic_pool_properties": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"max_size_bytes": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"zone_redundant": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"license_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {

			name, _ := diff.GetOk("sku.0.name")
			capacity, _ := diff.GetOk("sku.0.capacity")
			minCapacity, _ := diff.GetOk("per_database_settings.0.min_capacity")
			maxCapacity, _ := diff.GetOk("per_database_settings.0.max_capacity")

			if strings.HasPrefix(strings.ToLower(name.(string)), "gp_") {

				if capacity.(int) > 24 {
					return fmt.Errorf("GeneralPurpose pricing tier only supports upto 24 vCores")
				}

				if capacity.(int) < 1 {
					return fmt.Errorf("GeneralPurpose pricing tier must have a minimum of 1 vCores")
				}

				switch {
				case capacity.(int) == 1:
				case capacity.(int) == 2:
				case capacity.(int) == 4:
				case capacity.(int) == 8:
				case capacity.(int) == 16:
				case capacity.(int) == 24:
				default:
					return fmt.Errorf("GeneralPurpose pricing tier must have a capacity of 1, 2, 4, 8, 16, or 24 vCores")
				}

			}

			if strings.HasPrefix(strings.ToLower(name.(string)), "bc_") {
				if capacity.(int) > 80 {
					return fmt.Errorf("BusinessCritical pricing tier only supports upto 80 vCores")
				}

				if capacity.(int) < 2 {
					return fmt.Errorf("BusinessCritical pricing tier must have a minimum of 2 vCores")
				}

				switch {
				case capacity.(int) == 1:
				case capacity.(int) == 2:
				case capacity.(int) == 4:
				case capacity.(int) == 8:
				case capacity.(int) == 16:
				case capacity.(int) == 24:
				case capacity.(int) == 32:
				case capacity.(int) == 40:
				case capacity.(int) == 80:
				default:
					return fmt.Errorf("BusinessCritical pricing tier must have a capacity of 2, 4, 8, 16, 24, 32, 40, or 80 vCores")
				}
			}

			// Addutional checks based of SKU type...
			if strings.HasPrefix(strings.ToLower(name.(string)), "gp_") || strings.HasPrefix(strings.ToLower(name.(string)), "bc_") {
				// vCore based
				capacity, _ := diff.GetOk("sku.0.capacity")
				minCapacity, _ := diff.GetOk("per_database_settings.0.min_capacity")
				maxCapacity, _ := diff.GetOk("per_database_settings.0.max_capacity")

				if maxCapacity.(float64) > float64(capacity.(int)) {
					return fmt.Errorf("BusinessCritical pricing tier must have a capacity of 2, 4, 8, 16, 24, 32, 40, or 80 vCores")
				}

				if maxCapacity.(float64) > float64(capacity.(int)) {
					return fmt.Errorf("BusinessCritical and GeneralPurpose pricing tiers perDatabaseSettings maxCapacity must not be higher than the SKUs capacity value")
				}

				if minCapacity.(float64) > maxCapacity.(float64) {
					return fmt.Errorf("perDatabaseSettings maxCapacity must be greater than or equal to the perDatabaseSettings minCapacity value")
				}
			} else {
				// DTU based
				if maxCapacity.(float64) != math.Trunc(maxCapacity.(float64)) {
					return fmt.Errorf("BasicPool, StandardPool, and PremiumPool SKUs must have whole numbers as their maxCapacity")
				}

				if minCapacity.(float64) != math.Trunc(minCapacity.(float64)) {
					return fmt.Errorf("BasicPool, StandardPool, and PremiumPool SKUs must have whole numbers as their minCapacity")
				}

				if minCapacity.(float64) < 0.0 {
					return fmt.Errorf("BasicPool, StandardPool, and PremiumPool SKUs per_database_settings min_capacity must be equal to or greater than zero")
				}
			}

			return nil
		},
	}
}

func resourceArmMsSqlElasticPoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).MsSqlElasticPoolsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for MsSQL ElasticPool creation.")

	elasticPoolName := d.Get("name").(string)
	serverName := d.Get("server_name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	sku := expandAzureRmMsSqlElasticPoolSku(d)
	properties := expandAzureRmMsSqlElasticPoolProperties(d)
	tags := d.Get("tags").(map[string]interface{})

	elasticPool := sql.ElasticPool{
		Sku: sku,
		ElasticPoolProperties: properties,
		Location:              &location,
		Tags:                  expandTags(tags),
		Name:                  &elasticPoolName,
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, serverName, elasticPoolName, elasticPool)
	if err != nil {
		return err
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
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
	client := meta.(*ArmClient).MsSqlElasticPoolsClient
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
		return fmt.Errorf("Error flattening `sku`: %+v", err)
	}

	if err := d.Set("elastic_pool_properties", flattenAzureRmMsSqlElasticPoolProperties(resp.ElasticPoolProperties)); err != nil {
		return fmt.Errorf("Error flattening `elastic_pool_properties`: %+v", err)
	}

	if err := d.Set("per_database_settings", flattenAzureRmMsSqlElasticPoolPerDatabaseSettings(resp.ElasticPoolProperties.PerDatabaseSettings)); err != nil {
		return fmt.Errorf("Error flattening `per_database_settings`: %+v", err)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmMsSqlElasticPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).MsSqlElasticPoolsClient
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

func expandAzureRmMsSqlElasticPoolProperties(d *schema.ResourceData) *sql.ElasticPoolProperties {
	perDatabaseSettings := d.Get("per_database_settings").([]interface{})
	perDatabaseSetting := perDatabaseSettings[0].(map[string]interface{})

	minCapacity := perDatabaseSetting["min_capacity"].(float64)
	maxCapacity := perDatabaseSetting["max_capacity"].(float64)

	elasticPoolPerDatabaseSettings := &sql.ElasticPoolPerDatabaseSettings{
		MinCapacity: utils.Float(minCapacity),
		MaxCapacity: utils.Float(maxCapacity),
	}

	props := &sql.ElasticPoolProperties{
		PerDatabaseSettings: elasticPoolPerDatabaseSettings,
	}

	return props
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

	values["tier"] = *resp.Tier

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

	if maxSizeBytes := resp.MaxSizeBytes; maxSizeBytes != nil {
		elasticPoolProperty["max_size_bytes"] = *maxSizeBytes
	}

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
