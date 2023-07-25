// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/iotconnectors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type HealthCareMedTechServiceFhirDestinationV0ToV1 struct{}

func (s HealthCareMedTechServiceFhirDestinationV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"medtech_service_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"destination_fhir_service_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"destination_identity_resolution_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"destination_fhir_mapping_json": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (s HealthCareMedTechServiceFhirDestinationV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := iotconnectors.ParseFhirDestinationIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
