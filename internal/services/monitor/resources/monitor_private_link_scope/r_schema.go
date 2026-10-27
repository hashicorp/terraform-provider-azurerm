// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_private_link_scope

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-07-01-preview/privatelinkscopesapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func monitorPrivateLinkScopeSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.PrivateLinkScopeName,
		},

		"ingestion_access_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(privatelinkscopesapis.AccessModeOpen),
			ValidateFunc: validation.StringInSlice(privatelinkscopesapis.PossibleValuesForAccessMode(), false),
		},

		"query_access_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(privatelinkscopesapis.AccessModeOpen),
			ValidateFunc: validation.StringInSlice(privatelinkscopesapis.PossibleValuesForAccessMode(), false),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"tags": commonschema.Tags(),
	}
}
