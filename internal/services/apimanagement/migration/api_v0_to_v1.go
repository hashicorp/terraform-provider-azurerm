// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ApiV0ToV1{}

type ApiV0ToV1 struct{}

func (ApiV0ToV1) Schema() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementApiName(),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"display_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"path": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"protocols": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"revision": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"revision_description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		// Optional
		"api_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"contact": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"url": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"import": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"content_value": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"content_format": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"wsdl_selector": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"service_name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"endpoint_name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},

		"license": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"url": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"service_url": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"subscription_key_parameter_names": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"header": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"query": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"subscription_required": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"terms_of_service_url": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"source_api_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"oauth2_authorization": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"authorization_server_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"scope": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"openid_authentication": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"openid_provider_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"bearer_token_sending_methods": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		// Computed
		"is_current": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"is_online": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"soap_pass_through": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
			Optional: true,
		},

		"version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
			Optional: true,
		},

		"version_description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"version_set_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
			Optional: true,
		},
	}

	return schema
}

func (ApiV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		apiId := fmt.Sprintf("%s;rev=%s", rawState["name"].(string), rawState["revision"].(string))
		oldId, err := parse.ApiID(rawState["id"].(string))
		if err != nil {
			return rawState, err
		}
		newId := parse.NewApiID(oldId.SubscriptionId, oldId.ResourceGroup, oldId.ServiceName, apiId).ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId
		return rawState, nil
	}
}
