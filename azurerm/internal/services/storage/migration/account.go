package migration

import (
	"context"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = AccountV0ToV1{}

type AccountV0ToV1 struct{}

func (AccountV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return accountSchemaForV0AndV1()
}

func (AccountV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// this should have been applied from pre-0.12 migration system; backporting just in-case
		accountType := rawState["account_type"].(string)
		split := strings.Split(accountType, "_")
		rawState["account_tier"] = split[0]
		rawState["account_replication_type"] = split[1]
		return rawState, nil
	}
}

var _ pluginsdk.StateUpgrade = AccountV1ToV2{}

type AccountV1ToV2 struct{}

func (AccountV1ToV2) Schema() map[string]*pluginsdk.Schema {
	return accountSchemaForV0AndV1()
}

func (AccountV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	// this should have been applied from pre-0.12 migration system; backporting just in-case
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		rawState["account_encryption_source"] = "Microsoft.Storage"
		return rawState, nil
	}
}

func accountSchemaForV0AndV1() map[string]*pluginsdk.Schema {
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

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"account_kind": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "Storage",
		},

		"account_type": {
			Type:       pluginsdk.TypeString,
			Optional:   true,
			Computed:   true,
			Deprecated: "This field has been split into `account_tier` and `account_replication_type`",
		},

		"account_tier": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"account_replication_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		// Only valid for BlobStorage accounts, defaults to "Hot" in create function
		"access_tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"custom_domain": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"use_subdomain": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"enable_blob_encryption": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"enable_file_encryption": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"enable_https_traffic_only": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"primary_location": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_location": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_blob_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_blob_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_queue_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_queue_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_table_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_table_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		// NOTE: The API does not appear to expose a secondary file endpoint
		"primary_file_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_access_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_access_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_blob_connection_string": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_blob_connection_string": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}
