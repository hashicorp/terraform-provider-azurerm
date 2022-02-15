package eventgrid

import (
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2021-12-01/eventgrid"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
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

func expandPublicNetworkAccess(d *pluginsdk.ResourceData) eventgrid.PublicNetworkAccess {
	if v, ok := d.GetOk("public_network_access_enabled"); ok {
		enabled := eventgrid.PublicNetworkAccessDisabled
		if v.(bool) {
			enabled = eventgrid.PublicNetworkAccessEnabled
		}
		return enabled
	}
	return eventgrid.PublicNetworkAccessDisabled
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
	return in == eventgrid.PublicNetworkAccessEnabled
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

func expandIdentity(input []interface{}) (*eventgrid.IdentityInfo, error) {
	expanded, err := identity.ExpandSystemOrUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := eventgrid.IdentityInfo{
		Type: eventgrid.IdentityType(string(expanded.Type)),
	}

	if expanded.Type == identity.TypeUserAssigned {
		out.UserAssignedIdentities = make(map[string]*eventgrid.UserIdentityProperties)
		for k := range expanded.IdentityIds {
			out.UserAssignedIdentities[k] = &eventgrid.UserIdentityProperties{
				// intentionally empty
			}
		}
	}

	return &out, nil
}

func flattenIdentity(input *eventgrid.IdentityInfo) (*[]interface{}, error) {
	var transform *identity.SystemOrUserAssignedMap

	if input != nil {
		transform = &identity.SystemOrUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}

		for k, v := range input.UserAssignedIdentities {
			transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				ClientId:    v.ClientID,
				PrincipalId: v.PrincipalID,
			}
		}
		if input.PrincipalID != nil {
			transform.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transform.TenantId = *input.TenantID
		}
	}

	return identity.FlattenSystemOrUserAssignedMap(transform)
}
