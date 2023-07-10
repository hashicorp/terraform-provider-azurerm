// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/disasterrecoveryconfigs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/namespacesauthorizationrule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2022-01-01-preview/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandAuthorizationRuleRights(d *pluginsdk.ResourceData) *[]namespacesauthorizationrule.AccessRights {
	rights := make([]namespacesauthorizationrule.AccessRights, 0)

	if d.Get("listen").(bool) {
		rights = append(rights, namespacesauthorizationrule.AccessRightsListen)
	}

	if d.Get("send").(bool) {
		rights = append(rights, namespacesauthorizationrule.AccessRightsSend)
	}

	if d.Get("manage").(bool) {
		rights = append(rights, namespacesauthorizationrule.AccessRightsManage)
	}

	return &rights
}

func flattenAuthorizationRuleRights(rights *[]namespacesauthorizationrule.AccessRights) (listen, send, manage bool) {
	// zero (initial) value for a bool in go is false

	if rights != nil {
		for _, right := range *rights {
			switch right {
			case namespacesauthorizationrule.AccessRightsListen:
				listen = true
			case namespacesauthorizationrule.AccessRightsSend:
				send = true
			case namespacesauthorizationrule.AccessRightsManage:
				manage = true
			default:
				log.Printf("[DEBUG] Unknown Authorization Rule Right '%s'", right)
			}
		}
	}

	return listen, send, manage
}

func authorizationRuleSchemaFrom(s map[string]*pluginsdk.Schema) map[string]*pluginsdk.Schema {
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

func authorizationRuleCustomizeDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ interface{}) error {
	listen, hasListen := d.GetOk("listen")
	send, hasSend := d.GetOk("send")
	manage, hasManage := d.GetOk("manage")

	if !hasListen && !hasSend && !hasManage {
		return fmt.Errorf("One of the `listen`, `send` or `manage` properties needs to be set")
	}

	if manage.(bool) && (!listen.(bool) || !send.(bool)) {
		return fmt.Errorf("if `manage` is set both `listen` and `send` must be set to true too")
	}

	return nil
}

func waitForPairedNamespaceReplication(ctx context.Context, meta interface{}, id namespaces.NamespaceId, timeout time.Duration) error {
	namespaceClient := meta.(*clients.Client).ServiceBus.NamespacesClient
	resp, err := namespaceClient.Get(ctx, id)

	if model := resp.Model; model != nil {
		if !strings.EqualFold(string(model.Sku.Name), "Premium") {
			return err
		}
	}

	disasterRecoveryClient := meta.(*clients.Client).ServiceBus.DisasterRecoveryConfigsClient
	disasterRecoveryNamespaceId := disasterrecoveryconfigs.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
	disasterRecoveryResponse, err := disasterRecoveryClient.List(ctx, disasterRecoveryNamespaceId)

	if disasterRecoveryResponse.Model == nil {
		return err
	}

	if len(*disasterRecoveryResponse.Model) != 1 {
		return err
	}

	aliasName := (*disasterRecoveryResponse.Model)[0].Name

	disasterRecoveryConfigId := disasterrecoveryconfigs.NewDisasterRecoveryConfigID(disasterRecoveryNamespaceId.SubscriptionId, disasterRecoveryNamespaceId.ResourceGroupName, disasterRecoveryNamespaceId.NamespaceName, *aliasName)

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(disasterrecoveryconfigs.ProvisioningStateDRAccepted)},
		Target:     []string{string(disasterrecoveryconfigs.ProvisioningStateDRSucceeded)},
		MinTimeout: 30 * time.Second,
		Timeout:    timeout,
		Refresh: func() (interface{}, string, error) {
			resp, err := disasterRecoveryClient.Get(ctx, disasterRecoveryConfigId)
			if err != nil {
				return nil, "error", fmt.Errorf("wait read for %s: %v", disasterRecoveryConfigId, err)
			}

			if model := resp.Model; model != nil {
				if *model.Properties.ProvisioningState == disasterrecoveryconfigs.ProvisioningStateDRFailed {
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

func waitForNamespaceStatusToBeReady(ctx context.Context, meta interface{}, id namespaces.NamespaceId, timeout time.Duration) error {
	namespaceClient := meta.(*clients.Client).ServiceBus.NamespacesClient
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{
			string(namespaces.EndPointProvisioningStateUpdating),
			string(namespaces.EndPointProvisioningStateCreating),
			string(namespaces.EndPointProvisioningStateDeleting),
		},
		Target:                    []string{string(namespaces.EndPointProvisioningStateSucceeded)},
		Refresh:                   serviceBusNamespaceProvisioningStateRefreshFunc(ctx, namespaceClient, id),
		Timeout:                   timeout,
		PollInterval:              10 * time.Second,
		ContinuousTargetOccurence: 5,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return err
	}
	return nil
}

func serviceBusNamespaceProvisioningStateRefreshFunc(ctx context.Context, client *namespaces.NamespacesClient, id namespaces.NamespaceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving servicebus namespace error: %+v", err)
		}
		if res.Model == nil || res.Model.Properties == nil || res.Model.Properties.ProvisioningState == nil {
			return nil, "", fmt.Errorf("retrieving %s: model/provisioningState was nil", id)
		}

		return res, *res.Model.Properties.ProvisioningState, nil
	}
}
