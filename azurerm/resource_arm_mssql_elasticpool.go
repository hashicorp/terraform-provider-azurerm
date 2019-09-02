package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-10-01-preview/sql"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

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

			"max_size_bytes": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"max_size_gb"},
				ValidateFunc:  validation.IntAtLeast(0),
			},

			"max_size_gb": {
				Type:          schema.TypeFloat,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"max_size_bytes"},
				ValidateFunc:  validate.FloatAtLeast(0),
			},

			"zone_redundant": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"tags": tags.Schema(),
		},

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {

			if err := azure.MSSQLElasticPoolValidateSKU(diff); err != nil {
				return err
			}

			return nil
		},
	}
}

func resourceArmMsSqlElasticPoolCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).mssql.ElasticPoolsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for MSSQL ElasticPool creation.")

	elasticPoolName := d.Get("name").(string)
	serverName := d.Get("server_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, serverName, elasticPoolName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Elastic Pool %q (MSSQL Server %q / Resource Group %q): %s", elasticPoolName, serverName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_mssql_elasticpool", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	sku := expandAzureRmMsSqlElasticPoolSku(d)
	t := d.Get("tags").(map[string]interface{})

	elasticPool := sql.ElasticPool{
		Name:     &elasticPoolName,
		Location: &location,
		Sku:      sku,
		Tags:     tags.Expand(t),
		ElasticPoolProperties: &sql.ElasticPoolProperties{
			PerDatabaseSettings: expandAzureRmMsSqlElasticPoolPerDatabaseSettings(d),
		},
	}

	if v, ok := d.GetOkExists("zone_redundant"); ok {
		elasticPool.ElasticPoolProperties.ZoneRedundant = utils.Bool(v.(bool))
	}

	if d.HasChange("max_size_gb") {
		if v, ok := d.GetOk("max_size_gb"); ok {
			maxSizeBytes := v.(float64) * 1073741824
			elasticPool.MaxSizeBytes = utils.Int64(int64(maxSizeBytes))
		}
	} else {
		if v, ok := d.GetOk("max_size_bytes"); ok {
			elasticPool.MaxSizeBytes = utils.Int64(int64(v.(int)))
		}
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
	client := meta.(*ArmClient).mssql.ElasticPoolsClient
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
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("server_name", serverName)

	if err := d.Set("sku", flattenAzureRmMsSqlElasticPoolSku(resp.Sku)); err != nil {
		return fmt.Errorf("Error setting `sku`: %+v", err)
	}

	if properties := resp.ElasticPoolProperties; properties != nil {
		// Basic tier does not return max_size_bytes, so we need to skip setting this
		// value if the pricing tier is equal to Basic
		if tier, ok := d.GetOk("sku.0.tier"); ok {
			if !strings.EqualFold(tier.(string), "Basic") {
				d.Set("max_size_gb", float64(*properties.MaxSizeBytes/int64(1073741824)))
				d.Set("max_size_bytes", properties.MaxSizeBytes)
			}
		}
		d.Set("zone_redundant", properties.ZoneRedundant)

		//todo remove in 2.0
		if err := d.Set("elastic_pool_properties", flattenAzureRmMsSqlElasticPoolProperties(resp.ElasticPoolProperties)); err != nil {
			return fmt.Errorf("Error setting `elastic_pool_properties`: %+v", err)
		}

		if err := d.Set("per_database_settings", flattenAzureRmMsSqlElasticPoolPerDatabaseSettings(properties.PerDatabaseSettings)); err != nil {
			return fmt.Errorf("Error setting `per_database_settings`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMsSqlElasticPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).mssql.ElasticPoolsClient
	ctx := meta.(*ArmClient).StopContext

	resGroup, serverName, name, err := parseArmSqlElasticPoolId(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, resGroup, serverName, name)
	return err
}

func parseArmMsSqlElasticPoolId(sqlElasticPoolId string) (string, string, string, error) {
	id, err := azure.ParseAzureResourceID(sqlElasticPoolId)
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
