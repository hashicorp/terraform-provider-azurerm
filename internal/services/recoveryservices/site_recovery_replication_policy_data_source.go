// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceSiteRecoveryReplicationPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{

		Read: dataSourceSiteRecoveryReplicationPolicyRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
			"recovery_vault_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"recovery_point_retention_in_minutes": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},
			"application_consistent_snapshot_frequency_in_minutes": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceSiteRecoveryReplicationPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := replicationpolicies.NewReplicationPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("recovery_vault_name").(string), d.Get("name").(string))

	client := meta.(*clients.Client).RecoveryServices.ReplicationPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) || wasBadRequestWithNotExist(resp.HttpResponse, err) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("making Read request on site recovery replication policy %s : %+v", id.String(), err)
	}

	d.SetId(id.ID())
	d.Set("name", id.ReplicationPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.VaultName)

	if model := resp.Model; model != nil {
		a2APolicyDetails, isA2A := expandA2APolicyDetail(model)
		if isA2A {
			d.Set("recovery_point_retention_in_minutes", a2APolicyDetails.RecoveryPointHistory)
			d.Set("application_consistent_snapshot_frequency_in_minutes", a2APolicyDetails.AppConsistentFrequencyInMinutes)
		}
	}

	return nil
}
