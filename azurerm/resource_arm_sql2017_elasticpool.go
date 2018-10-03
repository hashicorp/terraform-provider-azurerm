package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-10-01-preview/sql"
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

			"id": {
				Type:     schema.TypeString,
				Required: false,
				Computed: true,
			},

			"kind": {
				Type:     schema.TypeString,
				Required: false,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: false,
				Computed: true,
			},

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
							Required:         false,
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
							Required: false,
							ValidateFunc: validation.StringInSlice([]string{
								"Gen4",
								"Gen5",
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},
					},
				},
			},

			"properties": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:     schema.TypeInt,
							Required: false,
						},

						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"max_size_bytes": {
							Type:     schema.TypeInt,
							Required: false,
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
							Required: false,
						},

						"license_type": {
							Type:     schema.TypeString,
							Required: false,
							ValidateFunc: validation.StringInSlice([]string{
								"BasicPrice",
								"LicenseIncluded",
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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
	serverName := d.Get("server_name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	sku := expandAzureRmSql2017ElasticPoolSku(d)
	properties := expandAzureRmSql2017ElasticPoolProperties(d)

	elasticPool := sql.ElasticPool{
		Name:                  &elasticPoolName,
		Location:              &location,
		ElasticPoolProperties: expandAzureRmSql2017ElasticPoolProperties(d),
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
	client := meta.(*ArmClient).sqlElasticPoolsClient
	ctx := meta.(*ArmClient).StopContext

	resGroup, serverName, name, err := parseArmSqlElasticPoolId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, resGroup, serverName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Sql Elastic Pool %s: %s", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}
	d.Set("server_name", serverName)

	if elasticPool := resp.ElasticPoolProperties; elasticPool != nil {
		d.Set("edition", string(elasticPool.Edition))
		d.Set("dtu", int(*elasticPool.Dtu))
		d.Set("db_dtu_min", int(*elasticPool.DatabaseDtuMin))
		d.Set("db_dtu_max", int(*elasticPool.DatabaseDtuMax))
		d.Set("pool_size", int(*elasticPool.StorageMB))

		if date := elasticPool.CreationDate; date != nil {
			d.Set("creation_date", date.Format(time.RFC3339))
		}
	}

	return nil
}

func resourceArmSql2017ElasticPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlElasticPoolsClient
	ctx := meta.(*ArmClient).StopContext

	resGroup, serverName, name, err := parseArmSqlElasticPoolId(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, resGroup, serverName, name)

	return err
}

func expandAzureRmSql2017ElasticPoolProperties(d *schema.ResourceData) *sql.ElasticPoolProperties {
	edition := sql.ElasticPoolEdition(d.Get("edition").(string))
	dtu := int32(d.Get("dtu").(int))

	props := &sql.ElasticPoolProperties{
		Edition: edition,
		Dtu:     &dtu,
	}

	if databaseDtuMin, ok := d.GetOk("db_dtu_min"); ok {
		databaseDtuMin := int32(databaseDtuMin.(int))
		props.DatabaseDtuMin = &databaseDtuMin
	}

	if databaseDtuMax, ok := d.GetOk("db_dtu_max"); ok {
		databaseDtuMax := int32(databaseDtuMax.(int))
		props.DatabaseDtuMax = &databaseDtuMax
	}

	if poolSize, ok := d.GetOk("pool_size"); ok {
		poolSize := int32(poolSize.(int))
		props.StorageMB = &poolSize
	}

	return props
}

func parseArmSql2017ElasticPoolId(sqlElasticPoolId string) (string, string, string, error) {
	id, err := parseAzureResourceID(sqlElasticPoolId)
	if err != nil {
		return "", "", "", fmt.Errorf("[ERROR] Unable to parse SQL ElasticPool ID %q: %+v", sqlElasticPoolId, err)
	}

	return id.ResourceGroup, id.Path["servers"], id.Path["elasticPools"], nil
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
		size:     utils.String(size),
		Tier:     utils.String(tier),
		Family:   utils.String(family),
		Capacity: utils.Int32(int32(capacity)),
	}
}
