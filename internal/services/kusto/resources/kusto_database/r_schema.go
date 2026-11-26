// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kusto_database

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func kustoDatabaseSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Kusto Database to create.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DatabaseName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"cluster_name": {
			Description:  "The name of the Kusto Cluster this database will be added to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ClusterName,
		},

		"soft_delete_period": {
			Description:  "The time the data should be kept before it stops being accessible to queries as ISO 8601 timespan. For more information see: [ISO 8601 Timespan](https://en.wikipedia.org/wiki/ISO_8601#Durations).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.ISO8601Duration,
		},

		"hot_cache_period": {
			Description:  "The time the data that should be kept in cache for fast queries as ISO 8601 timespan. For more information see: [ISO 8601 Timespan](https://en.wikipedia.org/wiki/ISO_8601#Durations).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.ISO8601Duration,
		},

		"size": {
			Description: "The size of the database in bytes.",
			Type:        pluginsdk.TypeFloat,
			Computed:    true,
		},
	}
}
