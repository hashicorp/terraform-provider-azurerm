// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace_disaster_recovery_config

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/disasterrecoveryconfigs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventHubNamespaceDisasterRecoveryConfigUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Eventhub.DisasterRecoveryConfigsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := disasterrecoveryconfigs.ParseDisasterRecoveryConfigID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.NamespaceName, eventHubNamespaceResourceName)
	defer locks.UnlockByName(id.NamespaceName, eventHubNamespaceResourceName)

	pairingStatus, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("checking the status of eventhub disaster recovery error: %+v", err)
	}

	// need to check if DCR needs pair-breaking first
	breakPairFirst := false
	if model := pairingStatus.Model; model != nil {
		if model.Properties != nil {
			if model.Properties.PartnerNamespace != nil && *model.Properties.PartnerNamespace != "" {
				breakPairFirst = true
			}
		}
	}

	if d.HasChange("partner_namespace_id") && breakPairFirst {
		// break pairing
		if _, err := client.BreakPairing(ctx, *id); err != nil {
			return fmt.Errorf("breaking the pairing for %s: %+v", *id, err)
		}
	}

	if err := resourceEventHubNamespaceDisasterRecoveryConfigWaitForState(ctx, client, *id); err != nil {
		return fmt.Errorf("waiting for the pairing to be broken for %s: %+v", *id, err)
	}

	parameters := disasterrecoveryconfigs.ArmDisasterRecovery{
		Properties: &disasterrecoveryconfigs.ArmDisasterRecoveryProperties{
			PartnerNamespace: pointer.To(d.Get("partner_namespace_id").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if err := resourceEventHubNamespaceDisasterRecoveryConfigWaitForState(ctx, client, *id); err != nil {
		return fmt.Errorf("waiting for replication after update of %s: %+v", *id, err)
	}

	return resourceEventHubNamespaceDisasterRecoveryConfigRead(d, meta)
}
