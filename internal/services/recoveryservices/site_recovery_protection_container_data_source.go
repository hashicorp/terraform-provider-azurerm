// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotectioncontainers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceSiteRecoveryProtectionContainer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceSiteRecoveryProtectionContainerRead,

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
			"recovery_fabric_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func dataSourceSiteRecoveryProtectionContainerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := replicationprotectioncontainers.NewReplicationProtectionContainerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("recovery_vault_name").(string), d.Get("recovery_fabric_name").(string), d.Get("name").(string))

	client := meta.(*clients.Client).RecoveryServices.ProtectionContainerClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("making Read request on site recovery protection container %s : %+v", id.String(), err)
	}

	d.SetId(handleAzureSdkForGoBug2824(id.ID()))
	d.Set("name", id.ReplicationProtectionContainerName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.VaultName)
	d.Set("recovery_fabric_name", id.ReplicationFabricName)

	return nil
}
