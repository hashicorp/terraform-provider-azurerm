// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/armdisasterrecoveries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func AuthorizationRuleSchemaFrom(s map[string]*pluginsdk.Schema) map[string]*pluginsdk.Schema {
	s["listen"] = &pluginsdk.Schema{
		Type:     pluginsdk.TypeBool,
		Optional: true,
		Default:  false,
	}
	s["send"] = &pluginsdk.Schema{
		Type:     pluginsdk.TypeBool,
		Optional: true,
		Default:  false,
	}
	s["manage"] = &pluginsdk.Schema{
		Type:     pluginsdk.TypeBool,
		Optional: true,
		Default:  false,
	}
	s["primary_key"] = &pluginsdk.Schema{
		Type:      pluginsdk.TypeString,
		Computed:  true,
		Sensitive: true,
	}
	s["primary_connection_string"] = &pluginsdk.Schema{
		Type:      pluginsdk.TypeString,
		Computed:  true,
		Sensitive: true,
	}
	s["secondary_key"] = &pluginsdk.Schema{
		Type:      pluginsdk.TypeString,
		Computed:  true,
		Sensitive: true,
	}
	s["secondary_connection_string"] = &pluginsdk.Schema{
		Type:      pluginsdk.TypeString,
		Computed:  true,
		Sensitive: true,
	}
	s["primary_connection_string_alias"] = &pluginsdk.Schema{
		Type:      pluginsdk.TypeString,
		Computed:  true,
		Sensitive: true,
	}
	s["secondary_connection_string_alias"] = &pluginsdk.Schema{
		Type:      pluginsdk.TypeString,
		Computed:  true,
		Sensitive: true,
	}
	return s
}

func AuthorizationRuleCustomizeDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ any) error {
	listen, hasListen := d.GetOk("listen")
	send, hasSend := d.GetOk("send")
	manage, hasManage := d.GetOk("manage")

	if !hasListen && !hasSend && !hasManage {
		return errors.New("one of the `listen`, `send` or `manage` properties needs to be set")
	}

	if manage.(bool) && (!listen.(bool) || !send.(bool)) {
		return errors.New("if `manage` is set both `listen` and `send` must be set to true too")
	}

	return nil
}

func WaitForPairedNamespaceReplication(ctx context.Context, meta any, id namespaces.NamespaceId, timeout time.Duration) error {
	namespaceClient := client.From(meta).NamespacesClient
	resp, err := namespaceClient.Get(ctx, id)

	if model := resp.Model; model != nil {
		if !strings.EqualFold(string(model.Sku.Name), "Premium") {
			return err
		}
	}

	disasterRecoveryClient := client.From(meta).ArmDisasterRecoveriesClient
	disasterRecoveryNamespaceId := armdisasterrecoveries.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
	disasterRecoveryResponse, err := disasterRecoveryClient.DisasterRecoveryConfigsList(ctx, disasterRecoveryNamespaceId)

	if disasterRecoveryResponse.Model == nil {
		return err
	}

	if len(*disasterRecoveryResponse.Model) != 1 {
		return err
	}

	aliasName := (*disasterRecoveryResponse.Model)[0].Name

	disasterRecoveryConfigId := armdisasterrecoveries.NewDisasterRecoveryConfigID(disasterRecoveryNamespaceId.SubscriptionId, disasterRecoveryNamespaceId.ResourceGroupName, disasterRecoveryNamespaceId.NamespaceName, *aliasName)

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(armdisasterrecoveries.ProvisioningStateDRAccepted)},
		Target:     []string{string(armdisasterrecoveries.ProvisioningStateDRSucceeded)},
		MinTimeout: 30 * time.Second,
		Timeout:    timeout,
		Refresh: func() (any, string, error) {
			resp, err := disasterRecoveryClient.DisasterRecoveryConfigsGet(ctx, disasterRecoveryConfigId)
			if err != nil {
				return nil, "error", fmt.Errorf("wait read for %s: %v", disasterRecoveryConfigId, err)
			}

			if model := resp.Model; model != nil {
				if *model.Properties.ProvisioningState == armdisasterrecoveries.ProvisioningStateDRFailed {
					return resp, "failed", fmt.Errorf("replication for %s failed", disasterRecoveryConfigId)
				}
				return resp, string(*model.Properties.ProvisioningState), nil
			}

			return resp, "nil", fmt.Errorf("waiting for replication error for %s: provisioning state is nil", disasterRecoveryConfigId)
		},
	}

	_, waitErr := stateConf.WaitForStateContext(ctx)
	return waitErr
}
