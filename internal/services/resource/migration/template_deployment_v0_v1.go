// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var _ pluginsdk.StateUpgrade = TemplateDeploymentV0ToV1{}

type TemplateDeploymentV0ToV1 struct{}

func (t TemplateDeploymentV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"template_body": {
			Type:      pluginsdk.TypeString,
			Optional:  true,
			Computed:  true,
			StateFunc: utils.NormalizeJson,
		},

		"parameters": {
			Type:          pluginsdk.TypeMap,
			Optional:      true,
			ConflictsWith: []string{"parameters_body"},
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"parameters_body": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			StateFunc:     utils.NormalizeJson,
			ConflictsWith: []string{"parameters"},
		},

		"deployment_mode": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"outputs": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (t TemplateDeploymentV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)

		id, err := parse.ResourceGroupTemplateDeploymentIDInsensitively(oldId)
		if err != nil {
			return rawState, fmt.Errorf("parsing existing Resource ID %q: %+v", oldId, err)
		}

		newId := id.ID()
		log.Printf("[DEBUG] Updating Resource ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
