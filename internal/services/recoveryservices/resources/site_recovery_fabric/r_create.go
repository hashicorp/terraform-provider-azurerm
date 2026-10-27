// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_fabric

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationfabrics"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSiteRecoveryFabricCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	location := location.Normalize(d.Get("location").(string))
	name := d.Get("name").(string)

	client := meta.(*clients.Client).RecoveryServices.FabricClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := replicationfabrics.NewReplicationFabricID(subscriptionId, resGroup, vaultName, name)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, replicationfabrics.DefaultGetOperationOptions())
		if err != nil {
			// NOTE: Bad Request due to https://github.com/Azure/azure-rest-api-specs/issues/12759
			if !response.WasNotFound(existing.HttpResponse) && !helpers.WasBadRequestWithNotExist(existing.HttpResponse, err) {
				return fmt.Errorf("checking for presence of existing site recovery fabric %s (vault %s): %+v", name, vaultName, err)
			}
		}

		if model := existing.Model; model != nil && model.Id != nil && *model.Id != "" {
			return tf.ImportAsExistsError("azurerm_site_recovery_fabric", helpers.HandleAzureSdkForGoBug2824(*model.Id))
		}
	}

	parameters := replicationfabrics.FabricCreationInput{
		Properties: &replicationfabrics.FabricCreationInputProperties{
			CustomDetails: replicationfabrics.AzureFabricCreationInput{
				Location: &location,
			},
		},
	}

	if err := client.CreateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating site recovery fabric %s (vault %s): %+v", name, vaultName, err)
	}
	d.SetId(id.ID())

	return resourceSiteRecoveryFabricRead(d, meta)
}
