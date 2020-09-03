package mssql

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
							ValidateFunc: validation.FloatAtLeast(0.0),
						},

						"max_capacity": {
							Type:         schema.TypeFloat,
							Required:     true,
							ValidateFunc: validation.FloatAtLeast(0.0),
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
				ValidateFunc:  validation.FloatAtLeast(0),
			},

			"zone_redundant": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"license_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.BasePrice),
					string(sql.LicenseIncluded),
				}, false),
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
	client := meta.(*clients.Client).MSSQL.ElasticPoolsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for MSSQL ElasticPool creation.")

	elasticPoolName := d.Get("name").(string)
	serverName := d.Get("server_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
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
			LicenseType:         sql.ElasticPoolLicenseType(d.Get("license_type").(string)),
			PerDatabaseSettings: expandAzureRmMsSqlElasticPoolPerDatabaseSettings(d),
			ZoneRedundant:       utils.Bool(d.Get("zone_redundant").(bool)),
		},
	}

	if d.HasChange("max_size_gb") {
		if v, ok := d.GetOk("max_size_gb"); ok {
			maxSizeBytes := v.(float64) * 1073741824
			elasticPool.MaxSizeBytes = utils.Int64(int64(maxSizeBytes))
		}
	} else if v, ok := d.GetOk("max_size_bytes"); ok {
		elasticPool.MaxSizeBytes = utils.Int64(int64(v.(int)))
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
	client := meta.(*clients.Client).MSSQL.ElasticPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	elasticPool, err := parse.MSSqlElasticPoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, elasticPool.ResourceGroup, elasticPool.MsSqlServer, elasticPool.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on MsSql Elastic Pool %s: %s", elasticPool.Name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", elasticPool.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("server_name", elasticPool.MsSqlServer)

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
		d.Set("license_type", string(properties.LicenseType))

		if err := d.Set("per_database_settings", flattenAzureRmMsSqlElasticPoolPerDatabaseSettings(properties.PerDatabaseSettings)); err != nil {
			return fmt.Errorf("Error setting `per_database_settings`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMsSqlElasticPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ElasticPoolsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	elasticPool, err := parse.MSSqlElasticPoolID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, elasticPool.ResourceGroup, elasticPool.MsSqlServer, elasticPool.Name)
	if err != nil {
		return fmt.Errorf("deleting ElasticPool %q (Server %q / Resource Group %q): %+v", elasticPool.Name, elasticPool.MsSqlServer, elasticPool.ResourceGroup, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of ElasticPool %q (Server %q / Resource Group %q): %+v", elasticPool.Name, elasticPool.MsSqlServer, elasticPool.ResourceGroup, err)
	}

	return nil
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

func flattenAzureRmMsSqlElasticPoolSku(input *sql.Sku) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	values := map[string]interface{}{}

	if name := input.Name; name != nil {
		values["name"] = *name
	}

	if tier := input.Tier; tier != nil {
		values["tier"] = *tier
	}

	if family := input.Family; family != nil {
		values["family"] = *family
	}

	if capacity := input.Capacity; capacity != nil {
		values["capacity"] = *capacity
	}

	return []interface{}{values}
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
