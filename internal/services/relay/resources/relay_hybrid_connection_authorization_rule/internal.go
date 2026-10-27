// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package relay_hybrid_connection_authorization_rule

import (
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandHybridConnectionAuthorizationRuleRights(d *pluginsdk.ResourceData) []hybridconnections.AccessRights {
	rights := make([]hybridconnections.AccessRights, 0)

	if d.Get("listen").(bool) {
		rights = append(rights, hybridconnections.AccessRightsListen)
	}

	if d.Get("send").(bool) {
		rights = append(rights, hybridconnections.AccessRightsSend)
	}

	if d.Get("manage").(bool) {
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
