// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-05-02/databaseprincipalassignments"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKustoDatabasePrincipalAssignment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoDatabasePrincipalAssignmentCreate,
		Read:   resourceKustoDatabasePrincipalAssignmentRead,
		Delete: resourceKustoDatabasePrincipalAssignmentDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KustoDatabasePrincipalAssignmentV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := databaseprincipalassignments.ParseDatabasePrincipalAssignmentID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": commonschema.ResourceGroupName(),

			"cluster_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ClusterName,
			},

			"database_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DatabaseName,
			},

			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DatabasePrincipalAssignmentName,
			},

			"tenant_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"tenant_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"principal_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"principal_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"principal_type": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(databaseprincipalassignments.PossibleValuesForPrincipalType(), false),
			},

			"role": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(databaseprincipalassignments.PossibleValuesForDatabasePrincipalRole(), false),
			},
		},
	}
}

func resourceKustoDatabasePrincipalAssignmentCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DatabasePrincipalAssignmentsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := databaseprincipalassignments.NewDatabasePrincipalAssignmentID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cluster_name").(string), d.Get("database_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_kusto_database_principal_assignment", id.ID())
	}

	principalAssignment := databaseprincipalassignments.DatabasePrincipalAssignment{
		Properties: &databaseprincipalassignments.DatabasePrincipalProperties{
			TenantId:      utils.String(d.Get("tenant_id").(string)),
			PrincipalId:   d.Get("principal_id").(string),
			PrincipalType: databaseprincipalassignments.PrincipalType(d.Get("principal_type").(string)),
			Role:          databaseprincipalassignments.DatabasePrincipalRole(d.Get("role").(string)),
		},
	}

	err = client.CreateOrUpdateThenPoll(ctx, id, principalAssignment)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceKustoDatabasePrincipalAssignmentRead(d, meta)
}

func resourceKustoDatabasePrincipalAssignmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DatabasePrincipalAssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := databaseprincipalassignments.ParseDatabasePrincipalAssignmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("cluster_name", id.ClusterName)
	d.Set("database_name", id.DatabaseName)
	d.Set("name", id.PrincipalAssignmentName)

	model := resp.Model

	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("principal_id", props.PrincipalId)
			d.Set("principal_name", props.PrincipalName)
			d.Set("principal_type", string(props.PrincipalType))
			d.Set("role", string(props.Role))
			d.Set("tenant_id", props.TenantId)
			d.Set("tenant_name", props.TenantName)
		}
	}

	return nil
}

func resourceKustoDatabasePrincipalAssignmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DatabasePrincipalAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := databaseprincipalassignments.ParseDatabasePrincipalAssignmentID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
