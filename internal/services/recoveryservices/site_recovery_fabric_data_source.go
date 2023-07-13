// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationfabrics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceSiteRecoveryFabric() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceSiteRecoveryFabricRead,

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
				ValidateFunc: validate.RecoveryServicesVaultName,
			},
			"location": commonschema.LocationComputed(),
		},
	}
}

func dataSourceSiteRecoveryFabricRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := replicationfabrics.NewReplicationFabricID(subscriptionId, d.Get("resource_group_name").(string), d.Get("recovery_vault_name").(string), d.Get("name").(string))

	client := meta.(*clients.Client).RecoveryServices.FabricClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, id, replicationfabrics.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("making read request on site recovery fabric %s: %+v", id.String(), err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("making read request on site recovery fabric %s: model is nil", id.String())
	}

	d.SetId(id.ID())
	d.Set("name", id.ReplicationFabricName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.VaultName)

	normalizedLocation := ""
	if l := model.Location; l != nil {
		normalizedLocation = location.NormalizeNilable(l)
	}
	d.Set("location", normalizedLocation)

	return nil
}
