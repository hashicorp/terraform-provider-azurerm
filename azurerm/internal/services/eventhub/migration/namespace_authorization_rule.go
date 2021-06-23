package migration

import (
	"context"
	"log"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = NamespaceAuthorizationRuleV0ToV1{}

type NamespaceAuthorizationRuleV0ToV1 struct{}

func (NamespaceAuthorizationRuleV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return authorizationRuleSchemaForV0AndV1()
}

func (NamespaceAuthorizationRuleV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)

		newId := strings.Replace(rawState["id"].(string), "/authorizationRules/", "/AuthorizationRules/", 1)

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId

		return rawState, nil
	}
}

var _ pluginsdk.StateUpgrade = NamespaceAuthorizationRuleV1ToV2{}

type NamespaceAuthorizationRuleV1ToV2 struct{}

func (NamespaceAuthorizationRuleV1ToV2) Schema() map[string]*pluginsdk.Schema {
	return authorizationRuleSchemaForV0AndV1()
}

func (NamespaceAuthorizationRuleV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)

		newId := strings.Replace(rawState["id"].(string), "/AuthorizationRules/", "/authorizationRules/", 1)

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId

		return rawState, nil
	}
}

func authorizationRuleSchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"namespace_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}
