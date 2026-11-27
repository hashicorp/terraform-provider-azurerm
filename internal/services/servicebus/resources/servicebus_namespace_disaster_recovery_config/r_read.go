// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_namespace_disaster_recovery_config

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/armdisasterrecoveries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/disasterrecoveryconfigs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	servicebusClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceBusNamespaceDisasterRecoveryConfigRead(d *pluginsdk.ResourceData, meta any) error {
	client := servicebusClient.From(meta).ArmDisasterRecoveriesClient
	authRuleClient := servicebusClient.From(meta).DisasterRecoveryConfigsClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := armdisasterrecoveries.ParseDisasterRecoveryConfigID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.DisasterRecoveryConfigsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	primaryId := armdisasterrecoveries.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)

	d.Set("name", id.DisasterRecoveryConfigName)
	d.Set("primary_namespace_id", primaryId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("partner_namespace_id", props.PartnerNamespace)
		}
	}

	// the auth rule cannot be retrieved by dr config name, the shared access policy should either be specified by user or using the default one which is `RootManageSharedAccessKey`
	authRuleId := disasterrecoveryconfigs.NewDisasterRecoveryConfigAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.DisasterRecoveryConfigName, serviceBusNamespaceDefaultAuthorizationRule)
	if input := d.Get("alias_authorization_rule_id").(string); input != "" {
		ruleId, err := disasterrecoveryconfigs.ParseDisasterRecoveryConfigAuthorizationRuleID(input)
		if err != nil {
			return fmt.Errorf("parsing primary namespace auth rule id error: %+v", err)
		}
		authRuleId = *ruleId
	}

	keys, err := authRuleClient.ListKeys(ctx, authRuleId)

	if err != nil {
		log.Printf("[WARN] listing default keys for %s: %+v", id, err)
	} else {
		if keysModel := keys.Model; keysModel != nil {
			d.Set("primary_connection_string_alias", keysModel.AliasPrimaryConnectionString)
			d.Set("secondary_connection_string_alias", keysModel.AliasSecondaryConnectionString)
			d.Set("default_primary_key", keysModel.PrimaryKey)
			d.Set("default_secondary_key", keysModel.SecondaryKey)
		}
	}

	return nil
}

func resourceServiceBusNamespaceDisasterRecoveryConfigWaitForState(ctx context.Context, client *armdisasterrecoveries.ArmDisasterRecoveriesClient, id armdisasterrecoveries.DisasterRecoveryConfigId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return errors.New("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(armdisasterrecoveries.ProvisioningStateDRAccepted)},
		Target:     []string{string(armdisasterrecoveries.ProvisioningStateDRSucceeded)},
		MinTimeout: 30 * time.Second,
		Timeout:    time.Until(deadline),
		Refresh: func() (any, string, error) {
			resp, err := client.DisasterRecoveryConfigsGet(ctx, id)
			if err != nil {
				return nil, "error", fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if *props.ProvisioningState == armdisasterrecoveries.ProvisioningStateDRFailed {
						return resp, "failed", fmt.Errorf("replication Failed for %s: %+v", id, err)
					}
					return resp, string(*props.ProvisioningState), nil
				}
			}

			return resp, "nil", fmt.Errorf("waiting on replication of %s: %+v", id, err)
		},
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
