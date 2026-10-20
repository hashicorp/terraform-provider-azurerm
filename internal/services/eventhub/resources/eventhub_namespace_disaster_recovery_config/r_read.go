// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace_disaster_recovery_config

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/disasterrecoveryconfigs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventHubNamespaceDisasterRecoveryConfigRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Eventhub.DisasterRecoveryConfigsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := disasterrecoveryconfigs.ParseDisasterRecoveryConfigID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.DisasterRecoveryConfigName)
	d.Set("namespace_name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil && model.Properties != nil {
		d.Set("partner_namespace_id", model.Properties.PartnerNamespace)
	}

	return nil
}

func resourceEventHubNamespaceDisasterRecoveryConfigWaitForState(ctx context.Context, client *disasterrecoveryconfigs.DisasterRecoveryConfigsClient, id disasterrecoveryconfigs.DisasterRecoveryConfigId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(disasterrecoveryconfigs.ProvisioningStateDRAccepted)},
		Target:     []string{string(disasterrecoveryconfigs.ProvisioningStateDRSucceeded)},
		MinTimeout: 30 * time.Second,
		Timeout:    time.Until(deadline),
		Refresh: func() (any, string, error) {
			read, err := client.Get(ctx, id)
			if err != nil {
				return nil, "error", fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := read.Model; model != nil {
				if props := model.Properties; props != nil {
					if props.ProvisioningState == nil {
						return read, "failed", fmt.Errorf("provisioningState was empty")
					}

					if *props.ProvisioningState == disasterrecoveryconfigs.ProvisioningStateDRFailed {
						return read, "failed", fmt.Errorf("replication failed for %s: %+v", id, err)
					}
					return read, string(*props.ProvisioningState), nil
				}
			}

			return read, "nil", fmt.Errorf("waiting for replication of %s: %+v", id, err)
		},
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
