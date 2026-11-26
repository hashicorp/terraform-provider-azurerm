// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace_schema_group

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/schemaregistry"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func eventhubNamespaceSchemaGroupSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this schema group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateSchemaGroupName(),
		},

		"namespace_id": {
			Description:  "The ID of the EventHub Namespace.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: namespaces.ValidateNamespaceID,
		},

		"schema_compatibility": {
			Description:  "The compatibility of this schema group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(schemaregistry.PossibleValuesForSchemaCompatibility(), false),
		},

		"schema_type": {
			Description: "The Type of this schema group.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(schemaregistry.SchemaTypeUnknown),
				string(schemaregistry.SchemaTypeAvro),
				"Json",
			}, false),
		},
	}
}
