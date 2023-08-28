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
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSiteRecoveryFabric() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSiteRecoveryFabricCreate,
		Read:   resourceSiteRecoveryFabricRead,
		Delete: resourceSiteRecoveryFabricDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ReplicationFabricID(id)
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"resource_group_name": commonschema.ResourceGroupName(),

			"recovery_vault_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.RecoveryServicesVaultName,
			},
			"location": commonschema.Location(),
		},
	}
}

func resourceSiteRecoveryFabricCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	name := d.Get("name").(string)

	client := meta.(*clients.Client).RecoveryServices.FabricClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := replicationfabrics.NewReplicationFabricID(subscriptionId, resGroup, vaultName, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id, replicationfabrics.DefaultGetOperationOptions())
		if err != nil {
			// NOTE: Bad Request due to https://github.com/Azure/azure-rest-api-specs/issues/12759
			if !response.WasNotFound(existing.HttpResponse) && !wasBadRequestWithNotExist(existing.HttpResponse, err) {
				return fmt.Errorf("checking for presence of existing site recovery fabric %s (vault %s): %+v", name, vaultName, err)
			}
		}

		if model := existing.Model; model != nil && model.Id != nil && *model.Id != "" {
			return tf.ImportAsExistsError("azurerm_site_recovery_fabric", handleAzureSdkForGoBug2824(*model.Id))
		}
	}

	parameters := replicationfabrics.FabricCreationInput{
		Properties: &replicationfabrics.FabricCreationInputProperties{
			CustomDetails: replicationfabrics.AzureFabricCreationInput{
				Location: &location,
			},
		},
	}

	err := client.CreateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating site recovery fabric %s (vault %s): %+v", name, vaultName, err)
	}

	d.SetId(id.ID())

	return resourceSiteRecoveryFabricRead(d, meta)
}

func resourceSiteRecoveryFabricRead(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := replicationfabrics.ParseReplicationFabricID(d.Id())
	if err != nil {
		return err
	}

	fabricClient := meta.(*clients.Client).RecoveryServices.FabricClient
	client := fabricClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, *id, replicationfabrics.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making read request on site recovery fabric %s: %+v", id.String(), err)
	}

	d.Set("name", id.ReplicationFabricName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.VaultName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if details := props.CustomDetails; details != nil {
				fabric, ok := details.(replicationfabrics.AzureFabricSpecificDetails)
				if !ok {
					return fmt.Errorf("expected `details` to be an AzureFabricSpecificDetails but it wasn't: %+v", details)
				}
				d.Set("location", location.NormalizeNilable(fabric.Location))
			}
		}
	}

	return nil
}

func resourceSiteRecoveryFabricDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := replicationfabrics.ParseReplicationFabricID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.FabricClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
