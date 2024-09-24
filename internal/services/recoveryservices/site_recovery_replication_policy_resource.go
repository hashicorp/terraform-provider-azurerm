// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSiteRecoveryReplicationPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSiteRecoveryReplicationPolicyCreate,
		Read:   resourceSiteRecoveryReplicationPolicyRead,
		Update: resourceSiteRecoveryReplicationPolicyUpdate,
		Delete: resourceSiteRecoveryReplicationPolicyDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ReplicationPolicyID(id)
			return err
		}),
		CustomizeDiff: resourceSiteRecoveryReplicationPolicyCustomDiff,

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
			"recovery_point_retention_in_minutes": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.IntBetween(0, 365*24*60),
			},
			"application_consistent_snapshot_frequency_in_minutes": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.IntBetween(0, 365*24*60),
			},
		},
	}
}

func resourceSiteRecoveryReplicationPolicyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	name := d.Get("name").(string)

	client := meta.(*clients.Client).RecoveryServices.ReplicationPoliciesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := replicationpolicies.NewReplicationPolicyID(subscriptionId, resGroup, vaultName, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			// NOTE: Bad Request due to https://github.com/Azure/azure-rest-api-specs/issues/12759
			if !response.WasNotFound(existing.HttpResponse) && !wasBadRequestWithNotExist(existing.HttpResponse, err) {
				return fmt.Errorf("checking for presence of existing site recovery replication policy %s: %+v", name, err)
			}
		}

		if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
			return tf.ImportAsExistsError("azurerm_site_recovery_replication_policy", *existing.Model.Id)
		}
	}

	recoveryPoint := int64(d.Get("recovery_point_retention_in_minutes").(int))
	appConsistency := int64(d.Get("application_consistent_snapshot_frequency_in_minutes").(int))
	if appConsistency > recoveryPoint {
		return fmt.Errorf("the value of `application_consistent_snapshot_frequency_in_minutes` must be less than or equal to the value of `recovery_point_retention_in_minutes`")
	}
	parameters := replicationpolicies.CreatePolicyInput{
		Properties: &replicationpolicies.CreatePolicyInputProperties{
			ProviderSpecificInput: &replicationpolicies.A2APolicyCreationInput{
				RecoveryPointHistory:            &recoveryPoint,
				AppConsistentFrequencyInMinutes: &appConsistency,
				MultiVMSyncStatus:               replicationpolicies.SetMultiVMSyncStatusEnable,
			},
		},
	}
	err := client.CreateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating site recovery replication policy %s (vault %s): %+v", name, vaultName, err)
	}

	d.SetId(id.ID())

	return resourceSiteRecoveryReplicationPolicyRead(d, meta)
}

func resourceSiteRecoveryReplicationPolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	name := d.Get("name").(string)

	client := meta.(*clients.Client).RecoveryServices.ReplicationPoliciesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := replicationpolicies.NewReplicationPolicyID(subscriptionId, resGroup, vaultName, name)

	recoveryPoint := int64(d.Get("recovery_point_retention_in_minutes").(int))
	appConsistency := int64(d.Get("application_consistent_snapshot_frequency_in_minutes").(int))
	if appConsistency > recoveryPoint {
		return fmt.Errorf("the value of `application_consistent_snapshot_frequency_in_minutes` must be less than or equal to the value of `recovery_point_retention_in_minutes`")
	}

	parameters := replicationpolicies.UpdatePolicyInput{
		Properties: &replicationpolicies.UpdatePolicyInputProperties{
			ReplicationProviderSettings: &replicationpolicies.A2APolicyCreationInput{
				RecoveryPointHistory:            &recoveryPoint,
				AppConsistentFrequencyInMinutes: &appConsistency,
				MultiVMSyncStatus:               replicationpolicies.SetMultiVMSyncStatusEnable,
			},
		},
	}
	err := client.UpdateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("updating site recovery replication policy %s (vault %s): %+v", name, vaultName, err)
	}

	return resourceSiteRecoveryReplicationPolicyRead(d, meta)
}

func resourceSiteRecoveryReplicationPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := replicationpolicies.ParseReplicationPolicyID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.ReplicationPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on site recovery replication policy %s : %+v", id.String(), err)
	}

	d.Set("name", id.ReplicationPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.VaultName)

	if model := resp.Model; model != nil {
		if a2APolicyDetails, isA2A := expandA2APolicyDetail(resp.Model); isA2A {
			d.Set("recovery_point_retention_in_minutes", a2APolicyDetails.RecoveryPointHistory)
			d.Set("application_consistent_snapshot_frequency_in_minutes", a2APolicyDetails.AppConsistentFrequencyInMinutes)
		}
	}
	return nil
}

func resourceSiteRecoveryReplicationPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := replicationpolicies.ParseReplicationPolicyID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.ReplicationPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting site recovery replication policy %s : %+v", id.String(), err)
	}

	return nil
}

func resourceSiteRecoveryReplicationPolicyCustomDiff(ctx context.Context, d *pluginsdk.ResourceDiff, i interface{}) error {
	retention := d.Get("recovery_point_retention_in_minutes").(int)
	frequency := d.Get("application_consistent_snapshot_frequency_in_minutes").(int)

	if retention == 0 && frequency > 0 {
		return fmt.Errorf("application_consistent_snapshot_frequency_in_minutes cannot be greater than zero when recovery_point_retention_in_minutes is set to zero")
	}

	return nil
}

func expandA2APolicyDetail(input *replicationpolicies.Policy) (out *replicationpolicies.A2APolicyDetails, isA2A bool) {
	if input.Properties == nil {
		return nil, false
	}
	if input.Properties.ProviderSpecificDetails == nil {
		return nil, false
	}
	detail, isA2A := input.Properties.ProviderSpecificDetails.(replicationpolicies.A2APolicyDetails)
	if isA2A {
		out = &detail
	}
	return
}
