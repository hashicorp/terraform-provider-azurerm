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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/helper"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSqlElasticPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceSqlElasticPoolCreateUpdate,
		Read:   resourceSqlElasticPoolRead,
		Update: resourceSqlElasticPoolCreateUpdate,
		Delete: resourceSqlElasticPoolDelete,

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
				ValidateFunc: helper.ValidateMsSqlServerName,
			},

			"edition": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.ElasticPoolEditionBasic),
					string(sql.ElasticPoolEditionStandard),
					string(sql.ElasticPoolEditionPremium),
				}, false),
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

func resourceSqlElasticPoolCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("creating/updating ElasticPool %q (Server %q / Resource Group %q): %+v", name, serverName, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of ElasticPool %q (Server %q / Resource Group %q): %+v", name, serverName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, serverName, name)
	if err != nil {
		return fmt.Errorf("retrieving ElasticPool %q (Server %q / Resource Group %q): %+v", name, serverName, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read SQL ElasticPool %q (resource group %q, server %q) ID", name, serverName, resGroup)
	}

	d.SetId(*read.ID)

	return resourceSqlElasticPoolRead(d, meta)
}

func resourceSqlElasticPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ElasticPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ElasticPoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving ElasticPool %q (Server %q / Resource Group %q): %+v", id.Name, id.ServerName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("server_name", id.ServerName)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.ElasticPoolProperties; props != nil {
		creationDate := ""
		if props.CreationDate != nil {
			creationDate = props.CreationDate.Format(time.RFC3339)
		}
		d.Set("creation_date", creationDate)

		dtu := 0
		if props.Dtu != nil {
			dtu = int(*props.Dtu)
		}
		d.Set("dtu", dtu)

		databaseDtuMin := 0
		if props.DatabaseDtuMin != nil {
			databaseDtuMin = int(*props.DatabaseDtuMin)
		}
		d.Set("db_dtu_min", databaseDtuMin)

		databaseDtuMax := 0
		if props.DatabaseDtuMax != nil {
			databaseDtuMax = int(*props.DatabaseDtuMax)
		}
		d.Set("db_dtu_max", databaseDtuMax)

		d.Set("edition", string(props.Edition))

		storageMb := 0
		if props.StorageMB != nil {
			storageMb = int(*props.StorageMB)
		}
		d.Set("pool_size", storageMb)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceSqlElasticPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ElasticPoolsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ElasticPoolID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, id.ResourceGroup, id.ServerName, id.Name); err != nil {
		return fmt.Errorf("deleting ElasticPool %q (Server %q / Resource Group %q): %+v", id.Name, id.ServerName, id.ResourceGroup, err)
	}

	return nil
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
