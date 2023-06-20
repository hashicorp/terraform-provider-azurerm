package relay

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2017-04-01/hybridconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2017-04-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

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
	return s
}

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
