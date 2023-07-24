// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid

import (
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2021-12-01/eventgrid" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/domains"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func IdentitySchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(eventgrid.IdentityTypeNone),
						string(eventgrid.IdentityTypeSystemAssigned),
						string(eventgrid.IdentityTypeUserAssigned),
					}, false),
				},

				"identity_ids": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},

				"principal_id": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"tenant_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func eventSubscriptionPublicNetworkAccessEnabled() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	}
}

func localAuthEnabled() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeBool,
		Optional: true,
		Default:  true,
	}
}

func eventSubscriptionInboundIPRule() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:       pluginsdk.TypeList,
		Optional:   true,
		MaxItems:   128,
		ConfigMode: pluginsdk.SchemaConfigModeAttr,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"ip_mask": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"action": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(eventgrid.IPActionTypeAllow),
					ValidateFunc: validation.StringInSlice([]string{
						string(eventgrid.IPActionTypeAllow),
					}, false),
				},
			},
		},
	}
}

func expandInboundIPRules(input []interface{}) *[]domains.InboundIPRule {
	if len(input) == 0 {
		return nil
	}

	rules := make([]domains.InboundIPRule, 0)
	for _, item := range input {
		rawRule := item.(map[string]interface{})
		rules = append(rules, domains.InboundIPRule{
			Action: pointer.To(domains.IPActionType(rawRule["action"].(string))),
			IPMask: utils.String(rawRule["ip_mask"].(string)),
		})
	}
	return &rules
}

func flattenPublicNetworkAccess(in eventgrid.PublicNetworkAccess) bool {
	return in == eventgrid.PublicNetworkAccessEnabled
}

func flattenInboundIPRules(input *[]domains.InboundIPRule) []interface{} {
	rules := make([]interface{}, 0)
	if input == nil {
		return rules
	}

	for _, r := range *input {
		action := ""
		if r.Action != nil {
			action = string(*r.Action)
		}

		rules = append(rules, map[string]interface{}{
			"action":  action,
			"ip_mask": pointer.From(r.IPMask),
		})
	}
	return rules
}
