// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sql

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSqlElasticPool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSqlElasticPoolCreateUpdate,
		Read:   resourceSqlElasticPoolRead,
		Update: resourceSqlElasticPoolCreateUpdate,
		Delete: resourceSqlElasticPoolDelete,

		DeprecationMessage: "The `azurerm_sql_elasticpool_resource` resource is deprecated and will be removed in version 4.0 of the AzureRM provider. Please use the `azurerm_mssql_elasticpool` resource instead.",

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ElasticPoolID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"server_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateMsSqlServerName,
			},

			"edition": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.ElasticPoolEditionBasic),
					string(sql.ElasticPoolEditionStandard),
					string(sql.ElasticPoolEditionPremium),
				}, false),
			},

			"dtu": {
				Type:     pluginsdk.TypeInt,
				Required: true,
			},

			"db_dtu_min": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Computed: true,
			},

			"db_dtu_max": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Computed: true,
			},

			"pool_size": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Computed: true,
			},

			"creation_date": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceSqlElasticPoolCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ElasticPoolsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewElasticPoolID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_sql_elasticpool", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	elasticPool := sql.ElasticPool{
		Name:                  utils.String(id.Name),
		Location:              &location,
		ElasticPoolProperties: getArmSqlElasticPoolProperties(d),
		Tags:                  tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.Name, elasticPool)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for create/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSqlElasticPoolRead(d, meta)
}

func resourceSqlElasticPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceSqlElasticPoolDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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

func getArmSqlElasticPoolProperties(d *pluginsdk.ResourceData) *sql.ElasticPoolProperties {
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
