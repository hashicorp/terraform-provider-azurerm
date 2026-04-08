// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobagents"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name mssql_job_agent -service-package-name mssql -properties "name" -compare-values "resource_group_name:database_id,server_name:database_id"

func resourceMsSqlJobAgent() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlJobAgentCreate,
		Read:   resourceMsSqlJobAgentRead,
		Update: resourceMsSqlJobAgentUpdate,
		Delete: resourceMsSqlJobAgentDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&jobagents.JobAgentId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&jobagents.JobAgentId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateMsSqlJobAgentName,
			},

			"database_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateSqlDatabaseID,
			},

			"location": commonschema.Location(),

			"identity": commonschema.UserAssignedIdentityOptional(),

			// This is a top level argument rather than a block because while Azure accepts input for both sku name and capacity fields,
			// the capacity must always be equal to the number included in the sku name.
			"sku": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      helper.SqlJobAgentSkuJA100,
				ValidateFunc: validation.StringInSlice(helper.PossibleValuesForJobAgentSku(), false),
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceMsSqlJobAgentCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.JobAgentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Job Agent creation.")

	databaseId := d.Get("database_id").(string)
	dbId, err := commonids.ParseSqlDatabaseID(databaseId)
	if err != nil {
		return err
	}
	id := jobagents.NewJobAgentID(dbId.SubscriptionId, dbId.ResourceGroupName, dbId.ServerName, d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_mssql_job_agent", id.ID())
	}

	params := jobagents.JobAgent{
		Name:     &id.JobAgentName,
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &jobagents.JobAgentProperties{
			DatabaseId: databaseId,
		},
		Sku: &jobagents.Sku{
			Name: d.Get("sku").(string),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	expandedIdentity, err := identity.ExpandUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	params.Identity = expandedIdentity

	err = client.CreateOrUpdateThenPoll(ctx, id, params)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceMsSqlJobAgentRead(d, meta)
}

func resourceMsSqlJobAgentUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.JobAgentsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Job Agent update.")

	databaseId := d.Get("database_id").(string)
	dbId, err := commonids.ParseSqlDatabaseID(databaseId)
	if err != nil {
		return err
	}
	id := jobagents.NewJobAgentID(dbId.SubscriptionId, dbId.ResourceGroupName, dbId.ServerName, d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving existing %s: `model` was nil", id)
	}
	params := existing.Model

	if d.HasChanges("identity") {
		expandedIdentity, err := identity.ExpandUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		params.Identity = expandedIdentity
	}

	if d.HasChanges("sku") {
		params.Sku = &jobagents.Sku{
			Name: d.Get("sku").(string),
		}
	}

	if d.HasChanges("tags") {
		params.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	err = client.CreateOrUpdateThenPoll(ctx, id, *params)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceMsSqlJobAgentRead(d, meta)
}

func resourceMsSqlJobAgentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.JobAgentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := jobagents.ParseJobAgentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %s", *id, err)
	}
	return resourceMssqlJobAgentSetFlatten(d, id, resp.Model)
}

func resourceMssqlJobAgentSetFlatten(d *pluginsdk.ResourceData, id *jobagents.JobAgentId, model *jobagents.JobAgent) error {
	d.Set("name", id.JobAgentName)

	if model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("database_id", props.DatabaseId)
		}

		flattenedIdentity, err := identity.FlattenUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		d.Set("identity", flattenedIdentity)

		if sku := model.Sku; sku != nil {
			d.Set("sku", sku.Name)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}

func resourceMsSqlJobAgentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.JobAgentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := jobagents.ParseJobAgentID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
