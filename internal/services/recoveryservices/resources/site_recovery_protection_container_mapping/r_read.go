// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_protection_container_mapping

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainermappings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSiteRecoveryContainerMappingRead(d *pluginsdk.ResourceData, meta any) error {
	id, err := replicationprotectioncontainermappings.ParseReplicationProtectionContainerMappingID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.ContainerMappingClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on site recovery protection container mapping %s : %+v", id.String(), err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("retrieving site recovery protection container mapping %s (vault %s): model is nil", id.ReplicationProtectionContainerMappingName, id.VaultName)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.VaultName)
	d.Set("recovery_fabric_name", id.ReplicationFabricName)
	d.Set("name", model.Name)
	if prop := model.Properties; prop != nil {
		d.Set("recovery_source_protection_container_name", model.Properties.SourceProtectionContainerFriendlyName)
		d.Set("recovery_replication_policy_id", model.Properties.PolicyId)
		d.Set("recovery_target_protection_container_id", model.Properties.TargetProtectionContainerId)

		if detail, ok := prop.ProviderSpecificDetails.(replicationprotectioncontainermappings.A2AProtectionContainerMappingDetails); ok {
			d.Set("automatic_update", flattenAutoUpdateSettings(&detail))
		} else {
			d.Set("automatic_update", flattenAutoUpdateSettings(nil))
		}
	}

	return nil
}

func flattenAutoUpdateSettings(input *replicationprotectioncontainermappings.A2AProtectionContainerMappingDetails) []any {
	if input == nil || pointer.From(input.AgentAutoUpdateStatus) == replicationprotectioncontainermappings.AgentAutoUpdateStatusDisabled {
		return []any{}
	}

	output := map[string]any{
		"automation_account_id": input.AutomationAccountArmId,
		"authentication_type":   pointer.FromEnum(input.AutomationAccountAuthenticationType),
	}

	return []any{output}
}
