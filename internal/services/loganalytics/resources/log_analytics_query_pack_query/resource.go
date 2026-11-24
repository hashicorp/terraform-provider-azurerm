// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_query_pack_query

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2019-09-01/querypackqueries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_log_analytics_query_pack_query"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return querypackqueries.ValidateQueryID
}
