// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package relay_namespace_authorization_rule

import (
	"log"

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
