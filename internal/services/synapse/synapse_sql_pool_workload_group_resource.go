// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSynapseSQLPoolWorkloadGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseSQLPoolWorkloadGroupCreateUpdate,
		Read:   resourceSynapseSQLPoolWorkloadGroupRead,
		Update: resourceSynapseSQLPoolWorkloadGroupCreateUpdate,
		Delete: resourceSynapseSQLPoolWorkloadGroupDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SqlPoolWorkloadGroupID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"sql_pool_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SqlPoolID,
			},

			"max_resource_percent": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 100),
			},

			"min_resource_percent": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 100),
			},

			"importance": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      "normal",
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"max_resource_percent_per_request": {
				Type:         pluginsdk.TypeFloat,
				Optional:     true,
				Default:      3,
				ValidateFunc: validation.FloatBetween(0, 100),
			},

			"min_resource_percent_per_request": {
				Type:         pluginsdk.TypeFloat,
				Optional:     true,
				ValidateFunc: validation.FloatBetween(0, 100),
			},

			"query_execution_timeout_in_seconds": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
		},
	}
}

func resourceSynapseSQLPoolWorkloadGroupCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.SQLPoolWorkloadGroupClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	sqlPoolId, err := parse.SqlPoolID(d.Get("sql_pool_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewSqlPoolWorkloadGroupID(sqlPoolId.SubscriptionId, sqlPoolId.ResourceGroup, sqlPoolId.WorkspaceName, sqlPoolId.Name, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, id.WorkloadGroupName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %q: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_synapse_sql_pool_workload_group", id.ID())
		}
	}

	parameters := synapse.WorkloadGroup{
		WorkloadGroupProperties: &synapse.WorkloadGroupProperties{
			Importance:                   utils.String(d.Get("importance").(string)),
			MaxResourcePercent:           utils.Int32(int32(d.Get("max_resource_percent").(int))),
			MaxResourcePercentPerRequest: utils.Float(d.Get("max_resource_percent_per_request").(float64)),
			MinResourcePercent:           utils.Int32(int32(d.Get("min_resource_percent").(int))),
			MinResourcePercentPerRequest: utils.Float(d.Get("min_resource_percent_per_request").(float64)),
		},
	}

	if timeout, ok := d.GetOk("query_execution_timeout_in_seconds"); ok {
		parameters.WorkloadGroupProperties.QueryExecutionTimeout = utils.Int32(int32(timeout.(int)))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, id.WorkloadGroupName, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %q: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSynapseSQLPoolWorkloadGroupRead(d, meta)
}

func resourceSynapseSQLPoolWorkloadGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.SQLPoolWorkloadGroupClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlPoolWorkloadGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, id.WorkloadGroupName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] synapse %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}
	d.Set("name", id.WorkloadGroupName)
	d.Set("sql_pool_id", parse.NewSqlPoolID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName).ID())
	if props := resp.WorkloadGroupProperties; props != nil {
		d.Set("importance", props.Importance)
		d.Set("max_resource_percent", props.MaxResourcePercent)
		d.Set("max_resource_percent_per_request", props.MaxResourcePercentPerRequest)
		d.Set("min_resource_percent", props.MinResourcePercent)
		d.Set("min_resource_percent_per_request", props.MinResourcePercentPerRequest)
		d.Set("query_execution_timeout_in_seconds", props.QueryExecutionTimeout)
	}
	return nil
}

func resourceSynapseSQLPoolWorkloadGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.SQLPoolWorkloadGroupClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlPoolWorkloadGroupID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, id.WorkloadGroupName)
	if err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %q: %+v", id, err)
	}
	return nil
}
