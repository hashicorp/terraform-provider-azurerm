// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func eventHubAuthorizationRuleSchemaFrom(s map[string]*pluginsdk.Schema) map[string]*pluginsdk.Schema {
	s["listen"] = &pluginsdk.Schema{
		Type:     pluginsdk.TypeBool,
		Optional: true,
		Default:  false,
	}
	s["manage"] = &pluginsdk.Schema{
		Type:     pluginsdk.TypeBool,
		Optional: true,
		Default:  false,
	}
	s["primary_connection_string"] = &pluginsdk.Schema{
		Type:      pluginsdk.TypeString,
		Computed:  true,
		Sensitive: true,
	}
	s["primary_connection_string_alias"] = &pluginsdk.Schema{
		Type:      pluginsdk.TypeString,
		Computed:  true,
		Sensitive: true,
	}
	s["primary_key"] = &pluginsdk.Schema{
		Type:      pluginsdk.TypeString,
		Computed:  true,
		Sensitive: true,
	}
	s["secondary_connection_string"] = &pluginsdk.Schema{
		Type:      pluginsdk.TypeString,
		Computed:  true,
		Sensitive: true,
	}
	s["secondary_connection_string_alias"] = &pluginsdk.Schema{
		Type:      pluginsdk.TypeString,
		Computed:  true,
		Sensitive: true,
	}

	s["secondary_key"] = &pluginsdk.Schema{
		Type:      pluginsdk.TypeString,
		Computed:  true,
		Sensitive: true,
	}

	s["send"] = &pluginsdk.Schema{
		Type:     pluginsdk.TypeBool,
		Optional: true,
		Default:  false,
	}
	return s
}

func eventHubAuthorizationRuleCustomizeDiff(_ context.Context, d *pluginsdk.ResourceDiff, _ interface{}) error {
	listen, hasListen := d.GetOk("listen")
	send, hasSend := d.GetOk("send")
	manage, hasManage := d.GetOk("manage")

	if !hasListen && !hasSend && !hasManage {
		return errors.New("one of the `listen`, `send` or `manage` properties needs to be set")
	}

	if manage.(bool) && !listen.(bool) && !send.(bool) {
		return errors.New("if `manage` is set both `listen` and `send` must be set to true too")
	}

	return nil
}
