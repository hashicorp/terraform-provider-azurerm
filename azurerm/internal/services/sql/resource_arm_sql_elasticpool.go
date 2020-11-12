package sql

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSqlElasticPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSqlElasticPoolCreateUpdate,
		Read:   resourceArmSqlElasticPoolRead,
		Update: resourceArmSqlElasticPoolCreateUpdate,
		Delete: resourceArmSqlElasticPoolDelete,

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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"server_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateMsSqlServerName,
			},

			"edition": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSqlElasticPoolEdition(),
			},

			"dtu": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"db_dtu_min": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"db_dtu_max": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"pool_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmSqlElasticPoolCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ElasticPoolsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for SQL ElasticPool creation.")

	name := d.Get("name").(string)
	serverName := d.Get("server_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, serverName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing SQL ElasticPool %q (resource group %q, server %q) ID", name, serverName, resGroup)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_sql_elasticpool", *existing.ID)
		}
	}

	elasticPool := sql.ElasticPool{
		Name:                  &name,
		Location:              &location,
		ElasticPoolProperties: getArmSqlElasticPoolProperties(d),
		Tags:                  tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, serverName, name, elasticPool)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, serverName, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read SQL ElasticPool %q (resource group %q, server %q) ID", name, serverName, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmSqlElasticPoolRead(d, meta)
}

func resourceArmSqlElasticPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ElasticPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

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
		return fmt.Errorf("Error making Read request on SQL ElasticPool %q (resource group %q, server %q) ID", name, serverName, resGroup)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
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

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmSqlElasticPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ElasticPoolsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup, serverName, name, err := parseArmSqlElasticPoolId(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, resGroup, serverName, name)

	return err
}

func getArmSqlElasticPoolProperties(d *schema.ResourceData) *sql.ElasticPoolProperties {
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

func parseArmSqlElasticPoolId(sqlElasticPoolId string) (string, string, string, error) {
	id, err := azure.ParseAzureResourceID(sqlElasticPoolId)
	if err != nil {
		return "", "", "", fmt.Errorf("[ERROR] Unable to parse SQL ElasticPool ID %q: %+v", sqlElasticPoolId, err)
	}

	return id.ResourceGroup, id.Path["servers"], id.Path["elasticPools"], nil
}

func validateSqlElasticPoolEdition() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(sql.ElasticPoolEditionBasic),
		string(sql.ElasticPoolEditionStandard),
		string(sql.ElasticPoolEditionPremium),
	}, false)
}
