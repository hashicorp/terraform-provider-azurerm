package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-10-01-preview/sql"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSql2017ElasticPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSql2017ElasticPoolCreate,
		Read:   resourceArmSql2017ElasticPoolRead,
		Update: resourceArmSql2017ElasticPoolCreate,
		Delete: resourceArmSql2017ElasticPoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"server_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
								"B_Gen4_1",
								"B_Gen4_2",
								"B_Gen5_1",
								"B_Gen5_2",
								"GP_Gen4_2",
								"GP_Gen4_4",
								"GP_Gen4_8",
								"GP_Gen4_16",
								"GP_Gen4_32",
								"GP_Gen5_2",
								"GP_Gen5_4",
								"GP_Gen5_8",
								"GP_Gen5_16",
								"GP_Gen5_32",
								"MO_Gen5_2",
								"MO_Gen5_4",
								"MO_Gen5_8",
								"MO_Gen5_16",
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"capacity": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: validateIntInSlice([]int{
								1,
								2,
								4,
								8,
								16,
								32,
							}),
						},

						"size": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"tier": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Basic",
								"GeneralPurpose",
								"MemoryOptimized",
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

			"elastic_pool_properties": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:     schema.TypeInt,
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
	}
}

func resourceArmSql2017ElasticPoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sql2017ElasticPoolsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for SQL2017 ElasticPool creation.")

	elasticPoolName := d.Get("name").(string)
	serverName := azureRMNormalizeLocation(d.Get("server_name").(string))
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	sku := expandAzureRmSql2017ElasticPoolSku(d)
	properties := expandAzureRmSql2017ElasticPoolProperties(d)
	resourceType := d.Get("type").(string)
	tags := d.Get("tags").(map[string]interface{})

	elasticPool := sql.ElasticPool{
		Sku: sku,
		ElasticPoolProperties: properties,
		Location:              &location,
		Tags:                  expandTags(tags),
		Name:                  &elasticPoolName,
		Type:                  &resourceType,
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
		return fmt.Errorf("Cannot read SQL2017 ElasticPool %q (resource group %q) ID", elasticPoolName, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmSql2017ElasticPoolRead(d, meta)
}

func resourceArmSql2017ElasticPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sql2017ElasticPoolsClient
	ctx := meta.(*ArmClient).StopContext

	resGroup, serverName, name, err := parseArmSql2017ElasticPoolId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, resGroup, serverName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Sql2017 Elastic Pool %s: %s", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	d.Set("server_name", serverName)

	if elasticPoolProperties := resp.ElasticPoolProperties; elasticPoolProperties != nil {
		if err := d.Set("per_database_settings", flattenAzureRmSql2017ElasticPoolPerDatabaseSettings(elasticPoolProperties.PerDatabaseSettings)); err != nil {
			return fmt.Errorf("Error flattening `per_database_settings`: %+v", err)
		}

		if maxSizeBytes := elasticPoolProperties.MaxSizeBytes; maxSizeBytes != nil {
			d.Set("max_size_bytes", elasticPoolProperties.MaxSizeBytes)
		}

		d.Set("state", string(elasticPoolProperties.State))

		if date := elasticPoolProperties.CreationDate; date != nil {
			d.Set("creation_date", date.String())
		}

		if zoneRedundant := elasticPoolProperties.ZoneRedundant; zoneRedundant != nil {
			d.Set("zone_redundant", bool(*elasticPoolProperties.ZoneRedundant))
		}

		d.Set("license_type", string(elasticPoolProperties.LicenseType))
	}

	return nil
}

func resourceArmSql2017ElasticPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sql2017ElasticPoolsClient
	ctx := meta.(*ArmClient).StopContext

	resGroup, serverName, name, err := parseArmSqlElasticPoolId(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, resGroup, serverName, name)

	return err
}

func expandAzureRmSql2017ElasticPoolProperties(d *schema.ResourceData) *sql.ElasticPoolProperties {
	elasticPoolProperties := d.Get("elastic_pool_properties").([]interface{})
	elasticPoolProperty := elasticPoolProperties[0].(map[string]interface{})

	state := sql.ElasticPoolState(elasticPoolProperty["State"].(string))
	creationDate, _ := elasticPoolProperty["CreationDate"].(date.Time)
	maxSizeBytes := elasticPoolProperty["MaxSizeBytes"].(int64)
	perDatabaseSettings := expandAzureRmSql2017ElasticPoolPerDatabaseSettings(elasticPoolProperty)
	zoneRedundant := elasticPoolProperty["ZoneRedundant"].(bool)
	licenseType := sql.ElasticPoolLicenseType(elasticPoolProperty["LicenseType"].(string))

	props := &sql.ElasticPoolProperties{
		State:               state,
		CreationDate:        &creationDate,
		MaxSizeBytes:        &maxSizeBytes,
		PerDatabaseSettings: perDatabaseSettings,
		ZoneRedundant:       utils.Bool(zoneRedundant),
		LicenseType:         licenseType,
	}

	return props
}

func parseArmSql2017ElasticPoolId(sqlElasticPoolId string) (string, string, string, error) {
	id, err := parseAzureResourceID(sqlElasticPoolId)
	if err != nil {
		return "", "", "", fmt.Errorf("[ERROR] Unable to parse SQL2017 ElasticPool ID %q: %+v", sqlElasticPoolId, err)
	}

	return id.ResourceGroup, id.Path["servers"], id.Path["elasticPools"], nil
}

func expandAzureRmSql2017ElasticPoolPerDatabaseSettings(d map[string]interface{}) *sql.ElasticPoolPerDatabaseSettings {
	perDatabasSettings := d["per_database_settings"].([]interface{})
	perDatabaseSetting := perDatabasSettings[0].(map[string]interface{})

	minCapacity := perDatabaseSetting["MinCapacity"].(float64)
	maxCapacity := perDatabaseSetting["MaxCapacity"].(float64)

	elasticPoolPerDatabaseSettings := &sql.ElasticPoolPerDatabaseSettings{
		MinCapacity: utils.Float(minCapacity),
		MaxCapacity: utils.Float(maxCapacity),
	}

	return elasticPoolPerDatabaseSettings
}

func expandAzureRmSql2017ElasticPoolSku(d *schema.ResourceData) *sql.Sku {
	skus := d.Get("sku").([]interface{})
	sku := skus[0].(map[string]interface{})

	name := sku["name"].(string)
	size := sku["size"].(string)
	tier := sku["tier"].(string)
	family := sku["family"].(string)
	capacity := sku["capacity"].(int)

	return &sql.Sku{
		Name:     utils.String(name),
		Size:     utils.String(size),
		Tier:     utils.String(tier),
		Family:   utils.String(family),
		Capacity: utils.Int32(int32(capacity)),
	}
}

func flattenAzureRmSql2017ElasticPoolPerDatabaseSettings(resp *sql.ElasticPoolPerDatabaseSettings) []interface{} {
	perDatabaseSettings := map[string]interface{}{}

	if minCapacity := resp.MinCapacity; minCapacity != nil {
		perDatabaseSettings["min_capacity"] = *minCapacity
	}

	if maxCapacity := resp.MinCapacity; maxCapacity != nil {
		perDatabaseSettings["max_capacity"] = *maxCapacity
	}

	return []interface{}{perDatabaseSettings}
}
