// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SpringCloudCustomizedAcceleratorV0ToV1 struct{}

func (s SpringCloudCustomizedAcceleratorV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"spring_cloud_accelerator_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"git_repository": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"url": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"basic_auth": {
						Type:          pluginsdk.TypeList,
						Optional:      true,
						ForceNew:      true,
						MaxItems:      1,
						ConflictsWith: []string{"git_repository.0.ssh_auth"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"username": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"password": {
									Type:      pluginsdk.TypeString,
									Required:  true,
									Sensitive: true,
								},
							},
						},
					},

					"ssh_auth": {
						Type:          pluginsdk.TypeList,
						Optional:      true,
						ForceNew:      true,
						MaxItems:      1,
						ConflictsWith: []string{"git_repository.0.basic_auth"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"private_key": {
									Type:      pluginsdk.TypeString,
									Required:  true,
									Sensitive: true,
								},

								"host_key": {
									Type:      pluginsdk.TypeString,
									Optional:  true,
									Sensitive: true,
								},

								"host_key_algorithm": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},

					"branch": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ExactlyOneOf: []string{"git_repository.0.branch", "git_repository.0.commit", "git_repository.0.git_tag"},
					},

					"ca_certificate_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"commit": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ExactlyOneOf: []string{"git_repository.0.branch", "git_repository.0.commit", "git_repository.0.git_tag"},
					},

					"git_tag": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ExactlyOneOf: []string{"git_repository.0.branch", "git_repository.0.commit", "git_repository.0.git_tag"},
					},

					"interval_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"accelerator_tags": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"icon_url": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (s SpringCloudCustomizedAcceleratorV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.SpringCloudCustomizedAcceleratorIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
