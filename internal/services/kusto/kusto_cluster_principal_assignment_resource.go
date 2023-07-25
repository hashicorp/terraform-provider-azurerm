// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-05-02/clusterprincipalassignments"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKustoClusterPrincipalAssignment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoClusterPrincipalAssignmentCreateUpdate,
		Read:   resourceKustoClusterPrincipalAssignmentRead,
		Delete: resourceKustoClusterPrincipalAssignmentDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KustoAttachedClusterPrincipalAssignmentV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := clusterprincipalassignments.ParsePrincipalAssignmentID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
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

			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ClusterPrincipalAssignmentName,
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
				ValidateFunc: validation.StringInSlice(clusterprincipalassignments.PossibleValuesForPrincipalType(), false),
			},

			"role": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(clusterprincipalassignments.PossibleValuesForClusterPrincipalRole(), false),
			},
		},
	}
}

func resourceKustoClusterPrincipalAssignmentCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClusterPrincipalAssignmentsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := clusterprincipalassignments.NewPrincipalAssignmentID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cluster_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		principalAssignment, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(principalAssignment.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(principalAssignment.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_kusto_cluster_principal_assignment", id.ID())
		}
	}

	tenantID := d.Get("tenant_id").(string)
	principalID := d.Get("principal_id").(string)
	principalType := d.Get("principal_type").(string)
	role := d.Get("role").(string)

	principalAssignment := clusterprincipalassignments.ClusterPrincipalAssignment{
		Properties: &clusterprincipalassignments.ClusterPrincipalProperties{
			TenantId:      utils.String(tenantID),
			PrincipalId:   principalID,
			PrincipalType: clusterprincipalassignments.PrincipalType(principalType),
			Role:          clusterprincipalassignments.ClusterPrincipalRole(role),
		},
	}

	err := client.CreateOrUpdateThenPoll(ctx, id, principalAssignment)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceKustoClusterPrincipalAssignmentRead(d, meta)
}

func resourceKustoClusterPrincipalAssignmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClusterPrincipalAssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusterprincipalassignments.ParsePrincipalAssignmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.PrincipalAssignmentName)
	d.Set("cluster_name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if resp.Model != nil {
		props := resp.Model.Properties

		if props != nil {
			d.Set("principal_id", props.PrincipalId)

			principalName := ""
			if props.PrincipalName != nil {
				principalName = *props.PrincipalName
			}
			d.Set("principal_name", principalName)

			d.Set("principal_type", string(props.PrincipalType))
			d.Set("role", string(props.Role))

			tenantID := ""
			if props.TenantId != nil {
				tenantID = *props.TenantId
			}
			d.Set("tenant_id", tenantID)

			tenantName := ""
			if props.TenantName != nil {
				tenantName = *props.TenantName
			}
			d.Set("tenant_name", tenantName)
		}
	}

	return nil
}

func resourceKustoClusterPrincipalAssignmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClusterPrincipalAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusterprincipalassignments.ParsePrincipalAssignmentID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
