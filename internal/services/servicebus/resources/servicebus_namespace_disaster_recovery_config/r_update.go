// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_namespace_disaster_recovery_config

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/armdisasterrecoveries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceBusNamespaceDisasterRecoveryConfigUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ServiceBus.ArmDisasterRecoveriesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := armdisasterrecoveries.ParseDisasterRecoveryConfigID(d.State().ID)
	if err != nil {
		return err
	}

	locks.ByName(id.NamespaceName, serviceBusNamespaceResourceName)
	defer locks.UnlockByName(id.NamespaceName, serviceBusNamespaceResourceName)

	if d.HasChange("partner_namespace_id") {
		if _, err := client.DisasterRecoveryConfigsBreakPairing(ctx, *id); err != nil {
			return fmt.Errorf("breaking the pairing for %s: %+v", *id, err)
		}
		if err := resourceServiceBusNamespaceDisasterRecoveryConfigWaitForState(ctx, client, *id); err != nil {
			return fmt.Errorf("waiting for the pairing to break for %s: %+v", *id, err)
		}
	}

	parameters := armdisasterrecoveries.ArmDisasterRecovery{
		Properties: &armdisasterrecoveries.ArmDisasterRecoveryProperties{
			PartnerNamespace: pointer.To(d.Get("partner_namespace_id").(string)),
		},
	}

	if _, err := client.DisasterRecoveryConfigsCreateOrUpdate(ctx, *id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", *id, err)
	}
	if err := resourceServiceBusNamespaceDisasterRecoveryConfigWaitForState(ctx, client, *id); err != nil {
		return fmt.Errorf("waiting for %s to finish replicating: %+v", *id, err)
	}

	return resourceServiceBusNamespaceDisasterRecoveryConfigRead(d, meta)
}
