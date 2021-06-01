package migration

import (
	"context"
	"log"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = MongoCollectionV0ToV1{}

type MongoCollectionV0ToV1 struct{}

func (MongoCollectionV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"account_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"database_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"shard_key": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"default_ttl_seconds": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"throughput": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},
	}
}

func (MongoCollectionV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId := strings.Replace(rawState["id"].(string), "apis/mongodb/databases", "mongodbDatabases", 1)

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId

		return rawState, nil
	}
}
