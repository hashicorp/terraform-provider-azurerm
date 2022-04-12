package migration

import (
	"context"

	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ServiceBusSubscriptionV0ToV1{}

type ServiceBusSubscriptionV0ToV1 struct{}

func (ServiceBusSubscriptionV0ToV1) Schema() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"topic_id": {
			Type:     pluginsdk.TypeString,
			Required: features.ThreePointOhBeta(),
			ForceNew: true,
		},

		"auto_delete_on_idle": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"default_message_ttl": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"lock_duration": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"dead_lettering_on_message_expiration": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"dead_lettering_on_filter_evaluation_error": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"enable_batched_operations": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"max_delivery_count": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"requires_session": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"forward_to": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"forward_dead_lettered_messages_to": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"status": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
	return s
}

func (ServiceBusSubscriptionV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {

		oldId := rawState["id"].(string)
		id, err := parse.SubscriptionID(oldId)
		if err != nil {
			return nil, err
		}

		rawState["id"] = id.ID()

		return rawState, nil
	}
}
