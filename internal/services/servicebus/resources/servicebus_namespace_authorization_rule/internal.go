// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_namespace_authorization_rule

import (
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespacesauthorizationrule"
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
