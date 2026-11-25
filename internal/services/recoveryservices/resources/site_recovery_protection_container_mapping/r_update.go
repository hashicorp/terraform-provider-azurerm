// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_protection_container_mapping

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainermappings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSiteRecoveryContainerMappingUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).RecoveryServices.ContainerMappingClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := replicationprotectioncontainermappings.ParseReplicationProtectionContainerMappingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on site recovery protection container mapping %s : %+v", id.String(), err)
	}

	if resp.Model == nil {
		return fmt.Errorf("making Read request on site recovery protection container mapping %s : `Model` is nil", id.String())
	}

	if resp.Model.Properties == nil {
		return fmt.Errorf("making Read request on site recovery protection container mapping %s : `Properties` is nil", id.String())
	}

	update := replicationprotectioncontainermappings.UpdateProtectionContainerMappingInput{
		Properties: &replicationprotectioncontainermappings.UpdateProtectionContainerMappingInputProperties{},
	}

	if d.HasChange("automatic_update") {
		autoUpdateEnabledValue, automationAccountArmId, authType := expandAutoUpdateSettings(d.Get("automatic_update").([]any))
		update.Properties.ProviderSpecificInput = replicationprotectioncontainermappings.A2AUpdateContainerMappingInput{
			AgentAutoUpdateStatus:               &autoUpdateEnabledValue,
			AutomationAccountArmId:              automationAccountArmId,
			AutomationAccountAuthenticationType: authType,
		}
	}

	if err = client.UpdateThenPoll(ctx, *id, update); err != nil {
		return fmt.Errorf("update %s: %+v", id, err)
	}

	return resourceSiteRecoveryContainerMappingRead(d, meta)
}
