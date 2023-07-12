// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationnetworkmappings"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSiteRecoveryNetworkMapping() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSiteRecoveryNetworkMappingCreate,
		Read:   resourceSiteRecoveryNetworkMappingRead,
		Delete: resourceSiteRecoveryNetworkMappingDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := replicationnetworkmappings.ParseReplicationNetworkMappingID(id)
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
			"source_recovery_fabric_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"target_recovery_fabric_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"source_network_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"target_network_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
		},
	}
}

func resourceSiteRecoveryNetworkMappingCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	fabricName := d.Get("source_recovery_fabric_name").(string)
	targetFabricName := d.Get("target_recovery_fabric_name").(string)
	sourceNetworkId := d.Get("source_network_id").(string)
	targetNetworkId := d.Get("target_network_id").(string)
	name := d.Get("name").(string)

	client := meta.(*clients.Client).RecoveryServices.NetworkMappingClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// get network name from id
	parsedSourceNetworkId, err := commonids.ParseVirtualNetworkID(sourceNetworkId)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to parse source_network_id '%s' (network mapping %s): %+v", sourceNetworkId, name, err)
	}

	id := replicationnetworkmappings.NewReplicationNetworkMappingID(subscriptionId, resGroup, vaultName, fabricName, parsedSourceNetworkId.VirtualNetworkName, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) && !wasBadRequestWithNotExist(existing.HttpResponse, err) {
				return fmt.Errorf("checking for presence of existing site recovery network mapping %s (vault %s): %+v", name, vaultName, err)
			}
		}

		if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
			return tf.ImportAsExistsError("azurerm_site_recovery_network_mapping", *existing.Model.Id)
		}
	}

	parameters := replicationnetworkmappings.CreateNetworkMappingInput{
		Properties: replicationnetworkmappings.CreateNetworkMappingInputProperties{
			RecoveryNetworkId:  targetNetworkId,
			RecoveryFabricName: &targetFabricName,
			FabricSpecificDetails: replicationnetworkmappings.AzureToAzureCreateNetworkMappingInput{
				PrimaryNetworkId: sourceNetworkId,
			},
		},
	}

	err = client.CreateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating site recovery network mapping %s (vault %s): %+v", name, vaultName, err)
	}

	d.SetId(id.ID())

	return resourceSiteRecoveryNetworkMappingRead(d, meta)
}

func resourceSiteRecoveryNetworkMappingRead(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := replicationnetworkmappings.ParseReplicationNetworkMappingID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.NetworkMappingClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on site recovery network mapping %q: %+v", id, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("retrieving site recovery network mapping %q: `model` was nil", id)
	}
	model := resp.Model

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.VaultName)
	d.Set("source_recovery_fabric_name", id.ReplicationFabricName)
	d.Set("name", id.ReplicationNetworkMappingName)
	if props := model.Properties; props != nil {
		d.Set("source_network_id", props.PrimaryNetworkId)
		d.Set("target_network_id", props.RecoveryNetworkId)

		targetFabricId, err := parse.ReplicationFabricID(handleAzureSdkForGoBug2824(*props.RecoveryFabricArmId))
		if err != nil {
			return err
		}
		d.Set("target_recovery_fabric_name", targetFabricId.Name)
	}

	return nil
}

func resourceSiteRecoveryNetworkMappingDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := replicationnetworkmappings.ParseReplicationNetworkMappingID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.NetworkMappingClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting site recovery network mapping %q: %+v", id, err)
	}

	return nil
}
