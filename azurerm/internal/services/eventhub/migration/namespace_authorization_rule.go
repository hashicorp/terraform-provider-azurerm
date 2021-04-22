package migration

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func NamespaceAuthorizationRuleV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		Type:    namespaceAuthorizationRuleSchemaForV0AndV1().CoreConfigSchema().ImpliedType(),
		Upgrade: namespaceAuthorizationRuleUpgradeV0ToV1,
		Version: 0,
	}
}

func NamespaceAuthorizationRuleV1ToV2() schema.StateUpgrader {
	return schema.StateUpgrader{
		Type:    namespaceAuthorizationRuleSchemaForV0AndV1().CoreConfigSchema().ImpliedType(),
		Upgrade: namespaceAuthorizationRuleUpgradeV1ToV2,
		Version: 1,
	}
}

func namespaceAuthorizationRuleSchemaForV0AndV1() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"namespace_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func namespaceAuthorizationRuleUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	oldId := rawState["id"].(string)

	newId := strings.Replace(rawState["id"].(string), "/authorizationRules/", "/AuthorizationRules/", 1)

	log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

	rawState["id"] = newId

	return rawState, nil
}

func namespaceAuthorizationRuleUpgradeV1ToV2(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	oldId := rawState["id"].(string)

	newId := strings.Replace(rawState["id"].(string), "/AuthorizationRules/", "/authorizationRules/", 1)

	log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

	rawState["id"] = newId

	return rawState, nil
}
