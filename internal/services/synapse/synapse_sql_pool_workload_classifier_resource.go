// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"fmt"
	"log"
	"regexp"
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

func resourceSynapseSQLPoolWorkloadClassifier() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseSQLPoolWorkloadClassifierCreateUpdate,
		Read:   resourceSynapseSQLPoolWorkloadClassifierRead,
		Update: resourceSynapseSQLPoolWorkloadClassifierCreateUpdate,
		Delete: resourceSynapseSQLPoolWorkloadClassifierDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SqlPoolWorkloadClassifierID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"workload_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SqlPoolWorkloadGroupID,
			},

			"member_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"context": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"end_time": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^\d{2}:\d{2}$`),
					"The `end_time` is of the `HH:MM` format in UTC time zone",
				),
			},

			"importance": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"low",
					"below_normal",
					"normal",
					"above_normal",
					"high",
				}, false),
			},

			"label": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"start_time": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^\d{2}:\d{2}$`),
					"The `start_time` is of the `HH:MM` format in UTC time zone",
				),
			},
		},
	}
}

func resourceSynapseSQLPoolWorkloadClassifierCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.SQLPoolWorkloadClassifierClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workloadGroupId, err := parse.SqlPoolWorkloadGroupID(d.Get("workload_group_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewSqlPoolWorkloadClassifierID(workloadGroupId.SubscriptionId, workloadGroupId.ResourceGroup, workloadGroupId.WorkspaceName, workloadGroupId.SqlPoolName, workloadGroupId.WorkloadGroupName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, id.WorkloadGroupName, id.WorkloadClassifierName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %q: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_synapse_sql_pool_workload_classifier", id.ID())
		}
	}

	parameters := synapse.WorkloadClassifier{
		WorkloadClassifierProperties: &synapse.WorkloadClassifierProperties{
			Context:    utils.String(d.Get("context").(string)),
			EndTime:    utils.String(d.Get("end_time").(string)),
			Importance: utils.String(d.Get("importance").(string)),
			Label:      utils.String(d.Get("label").(string)),
			MemberName: utils.String(d.Get("member_name").(string)),
			StartTime:  utils.String(d.Get("start_time").(string)),
		},
	}
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, id.WorkloadGroupName, id.WorkloadClassifierName, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %q: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSynapseSQLPoolWorkloadClassifierRead(d, meta)
}

func resourceSynapseSQLPoolWorkloadClassifierRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.SQLPoolWorkloadClassifierClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlPoolWorkloadClassifierID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, id.WorkloadGroupName, id.WorkloadClassifierName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] synapse %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}
	d.Set("name", id.WorkloadClassifierName)
	d.Set("workload_group_id", parse.NewSqlPoolWorkloadGroupID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, id.WorkloadGroupName).ID())
	if props := resp.WorkloadClassifierProperties; props != nil {
		d.Set("context", props.Context)
		d.Set("end_time", props.EndTime)
		d.Set("importance", props.Importance)
		d.Set("label", props.Label)
		d.Set("member_name", props.MemberName)
		d.Set("start_time", props.StartTime)
	}
	return nil
}

func resourceSynapseSQLPoolWorkloadClassifierDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.SQLPoolWorkloadClassifierClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlPoolWorkloadClassifierID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, id.WorkloadGroupName, id.WorkloadClassifierName)
	if err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %q: %+v", id, err)
	}
	return nil
}
