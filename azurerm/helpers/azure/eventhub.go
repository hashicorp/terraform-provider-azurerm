package azure

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventhub/mgmt/2018-01-01-preview/eventhub"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// validation
func ValidateEventHubNamespaceName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]{4,48}[a-zA-Z0-9]$"),
		"The namespace name can contain only letters, numbers and hyphens. The namespace must start with a letter, and it must end with a letter or number and be between 6 and 50 characters long.",
	)
}

func ValidateEventHubName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-zA-Z0-9]([-._a-zA-Z0-9]{0,48}[a-zA-Z0-9])?$"),
		"The event hub name can contain only letters, numbers, periods (.), hyphens (-),and underscores (_), up to 50 characters, and it must begin and end with a letter or number.",
	)
}

func ValidateEventHubConsumerName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-zA-Z0-9]([-._a-zA-Z0-9]{0,48}[a-zA-Z0-9])?$"),
		"The consumer group name can contain only letters, numbers, periods (.), hyphens (-),and underscores (_), up to 50 characters, and it must begin and end with a letter or number.",
	)
}

func ValidateEventHubAuthorizationRuleName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-zA-Z0-9]([-._a-zA-Z0-9]{0,48}[a-zA-Z0-9])?$"),
		"The authorization rule name can contain only letters, numbers, periods, hyphens and underscores. The name must start and end with a letter or number and be up to 50 characters long.",
	)
}

// schema
func ExpandEventHubAuthorizationRuleRights(d *schema.ResourceData) *[]eventhub.AccessRights {
	rights := make([]eventhub.AccessRights, 0)

	if d.Get("listen").(bool) {
		rights = append(rights, eventhub.Listen)
	}

	if d.Get("send").(bool) {
		rights = append(rights, eventhub.SendEnumValue)
	}

	if d.Get("manage").(bool) {
		rights = append(rights, eventhub.Manage)
	}

	return &rights
}

func FlattenEventHubAuthorizationRuleRights(rights *[]eventhub.AccessRights) (listen, send, manage bool) {
	// zero (initial) value for a bool in go is false

	if rights != nil {
		for _, right := range *rights {
			switch right {
			case eventhub.Listen:
				listen = true
			case eventhub.SendEnumValue:
				send = true
			case eventhub.Manage:
				manage = true
			default:
				log.Printf("[DEBUG] Unknown Authorization Rule Right '%s'", right)
			}
		}
	}

	return listen, send, manage
}

func EventHubAuthorizationRuleSchemaFrom(s map[string]*schema.Schema) map[string]*schema.Schema {
	authSchema := map[string]*schema.Schema{
		"listen": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},

		"manage": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},

		"primary_connection_string": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_connection_string_alias": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_key": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_connection_string": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_connection_string_alias": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_key": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"send": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
	return MergeSchema(s, authSchema)
}

func EventHubAuthorizationRuleCustomizeDiff(d *schema.ResourceDiff, _ interface{}) error {
	listen, hasListen := d.GetOk("listen")
	send, hasSend := d.GetOk("send")
	manage, hasManage := d.GetOk("manage")

	if !hasListen && !hasSend && !hasManage {
		return fmt.Errorf("One of the `listen`, `send` or `manage` properties needs to be set")
	}

	if manage.(bool) && !listen.(bool) && !send.(bool) {
		return fmt.Errorf("if `manage` is set both `listen` and `send` must be set to true too")
	}

	return nil
}
