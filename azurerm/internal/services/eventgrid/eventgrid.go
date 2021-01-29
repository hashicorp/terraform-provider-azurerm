package eventgrid

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/eventgrid/mgmt/2020-04-01-preview/eventgrid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func eventSubscriptionPublicNetworkAccessEnabled() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	}
}

func eventSubscriptionInboundIPRule() *schema.Schema {
	return &schema.Schema{
		Type:       schema.TypeList,
		Optional:   true,
		MaxItems:   128,
		ConfigMode: schema.SchemaConfigModeAttr,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ip_mask": {
					Type:     schema.TypeString,
					Required: true,
				},
				"action": {
					Type:     schema.TypeString,
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

func expandPublicNetworkAccess(d *schema.ResourceData) eventgrid.PublicNetworkAccess {
	if v, ok := d.GetOk("public_network_access_enabled"); ok {
		enabled := eventgrid.Disabled
		if v.(bool) {
			enabled = eventgrid.Enabled
		}
		return enabled
	}
	return eventgrid.Disabled
}

func expandInboundIPRules(d *schema.ResourceData) *[]eventgrid.InboundIPRule {
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
