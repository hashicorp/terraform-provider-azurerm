// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package relay

import (
	"context"
	"errors"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func authorizationRuleArgumentsFrom(s map[string]*pluginsdk.Schema) map[string]*pluginsdk.Schema {
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

	return s
}

func authorizationRuleAttributesFrom(s map[string]*pluginsdk.Schema) map[string]*pluginsdk.Schema {
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

	return s
}

func expandAuthorizationRuleRights(config RelayNamespaceAuthorizationRuleResourceModel) []namespaces.AccessRights {
	rights := make([]namespaces.AccessRights, 0)

	if config.Listen {
		rights = append(rights, namespaces.AccessRightsListen)
	}

	if config.Send {
		rights = append(rights, namespaces.AccessRightsSend)
	}

	if config.Manage {
		rights = append(rights, namespaces.AccessRightsManage)
	}

	return rights
}

func flattenAuthorizationRuleRights(rights []namespaces.AccessRights) (listen, send, manage bool) {
	// zero (initial) value for a bool in go is false

	for _, right := range rights {
		switch right {
		case namespaces.AccessRightsListen:
			listen = true
		case namespaces.AccessRightsSend:
			send = true
		case namespaces.AccessRightsManage:
			manage = true
		default:
			log.Printf("[DEBUG] Unknown Authorization Rule Right '%s'", right)
		}
	}

	return listen, send, manage
}

func expandHybridConnectionAuthorizationRuleRights(config RelayHybridConnectionAuthorizationRuleResourceModel) []hybridconnections.AccessRights {
	rights := make([]hybridconnections.AccessRights, 0)

	if config.Listen {
		rights = append(rights, hybridconnections.AccessRightsListen)
	}

	if config.Send {
		rights = append(rights, hybridconnections.AccessRightsSend)
	}

	if config.Manage {
		rights = append(rights, hybridconnections.AccessRightsManage)
	}

	return rights
}

func flattenHybridConnectionAuthorizationRuleRights(rights []hybridconnections.AccessRights) (listen, send, manage bool) {
	// zero (initial) value for a bool in go is false

	for _, right := range rights {
		switch right {
		case hybridconnections.AccessRightsListen:
			listen = true
		case hybridconnections.AccessRightsSend:
			send = true
		case hybridconnections.AccessRightsManage:
			manage = true
		default:
			log.Printf("[DEBUG] Unknown Authorization Rule Right '%s'", right)
		}
	}

	return listen, send, manage
}

func authorizationRuleCustomizeDiff(ctx context.Context, metadata sdk.ResourceMetaData) error {
	listen, hasListen := metadata.ResourceDiff.GetOk("listen")
	send, hasSend := metadata.ResourceDiff.GetOk("send")
	manage, hasManage := metadata.ResourceDiff.GetOk("manage")

	if !hasListen && !hasSend && !hasManage {
		return errors.New("one of the `listen`, `send` or `manage` properties needs to be set")
	}

	if manage.(bool) && (!listen.(bool) || !send.(bool)) {
		return errors.New("if `manage` is set both `listen` and `send` must be set to true too")
	}

	return nil
}
