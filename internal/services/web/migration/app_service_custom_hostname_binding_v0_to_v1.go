package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AppServiceCustomHostnameBindingV0ToV1 struct{}

func (AppServiceCustomHostnameBindingV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"hostname": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": {
			Type:     schema.TypeString,
			Required: true,
		},

		"app_service_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"ssl_state": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"thumbprint": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"virtual_ip": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (AppServiceCustomHostnameBindingV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		parsed, err := webapps.ParseHostNameBindingIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		newId := parsed.ID()
		log.Printf("[DEBUG] Upgrading `%s` to `%s`", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
