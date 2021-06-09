package eventgrid

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/eventgrid/mgmt/2020-10-15-preview/eventgrid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func eventSubscriptionPublicNetworkAccessEnabled() *pluginsdk.Schema {
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
					Default:  string(eventgrid.Allow),
					ValidateFunc: validation.StringInSlice([]string{
						string(eventgrid.Allow),
					}, false),
				},
			},
		},
	}
}

func expandPublicNetworkAccess(d *pluginsdk.ResourceData) eventgrid.PublicNetworkAccess {
	if v, ok := d.GetOk("public_network_access_enabled"); ok {
		enabled := eventgrid.Disabled
		if v.(bool) {
			enabled = eventgrid.Enabled
		}
		return enabled
	}
	return eventgrid.Disabled
}

func expandInboundIPRules(d *pluginsdk.ResourceData) *[]eventgrid.InboundIPRule {
	inboundIPRuleList := d.Get("inbound_ip_rule").([]interface{})
	if len(inboundIPRuleList) == 0 {
		return nil
	}

	rules := make([]eventgrid.InboundIPRule, 0)

	for _, r := range inboundIPRuleList {
		rawRule := r.(map[string]interface{})
		rule := &eventgrid.InboundIPRule{
			Action: eventgrid.IPActionType(rawRule["action"].(string)),
			IPMask: utils.String(rawRule["ip_mask"].(string)),
		}

		rules = append(rules, *rule)
	}
	return &rules
}

func flattenPublicNetworkAccess(in eventgrid.PublicNetworkAccess) bool {
	return in == eventgrid.Enabled
}

func flattenInboundIPRules(in *[]eventgrid.InboundIPRule) []interface{} {
	rules := make([]interface{}, 0)
	if in == nil {
		return rules
	}

	for _, r := range *in {
		rawRule := make(map[string]interface{})

		rawRule["action"] = string(r.Action)

		if r.IPMask != nil {
			rawRule["ip_mask"] = *r.IPMask
		}
		rules = append(rules, rawRule)
	}
	return rules
}
