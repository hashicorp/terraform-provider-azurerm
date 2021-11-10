package eventgrid

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventgrid/mgmt/2020-10-15-preview/eventgrid"
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

func IdentitySchemaForDataSource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
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

func expandIdentity(input []interface{}) (*eventgrid.IdentityInfo, error) {
	if len(input) == 0 || input[0] == nil {
		return &eventgrid.IdentityInfo{
			Type: eventgrid.IdentityTypeNone,
		}, nil
	}

	raw := input[0].(map[string]interface{})

	identity := eventgrid.IdentityInfo{
		Type: eventgrid.IdentityType(raw["type"].(string)),
	}

	identityIdsRaw := raw["identity_ids"].(*schema.Set).List()
	identityIds := make(map[string]*eventgrid.UserIdentityProperties)
	for _, v := range identityIdsRaw {
		identityIds[v.(string)] = &eventgrid.UserIdentityProperties{}
	}

	if len(identityIds) > 0 {
		if identity.Type != eventgrid.IdentityTypeUserAssigned && identity.Type != eventgrid.IdentityTypeSystemAssignedUserAssigned {
			return nil, fmt.Errorf("`identity_ids` can only be specified when `type` includes `UserAssigned`")
		}

		identity.UserAssignedIdentities = identityIds
	}

	return &identity, nil
}

func flattenIdentity(input *eventgrid.IdentityInfo) []interface{} {
	if input == nil || input.Type == eventgrid.IdentityTypeNone {
		return []interface{}{}
	}

	identityIds := make([]string, 0)
	if input.UserAssignedIdentities != nil {
		for k := range input.UserAssignedIdentities {
			identityIds = append(identityIds, k)
		}
	}

	principalID := ""
	if input.PrincipalID != nil {
		principalID = *input.PrincipalID
	}

	tenantID := ""
	if input.TenantID != nil {
		tenantID = *input.TenantID
	}

	return []interface{}{
		map[string]interface{}{
			"type":         string(input.Type),
			"identity_ids": identityIds,
			"principal_id": principalID,
			"tenant_id":    tenantID,
		},
	}
}
