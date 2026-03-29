// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package relay

import (
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandAuthorizationRuleRights(d *pluginsdk.ResourceData) []namespaces.AccessRights {
	rights := make([]namespaces.AccessRights, 0)

	if d.Get("listen").(bool) {
		rights = append(rights, namespaces.AccessRightsListen)
	}

	if d.Get("send").(bool) {
		rights = append(rights, namespaces.AccessRightsSend)
	}

	if d.Get("manage").(bool) {
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

func expandHybridConnectionAuthorizationRuleRights(config relayHybridConnectionAuthorizationRuleModel) []hybridconnections.AccessRights {
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
