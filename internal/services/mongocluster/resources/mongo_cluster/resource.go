// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mongo_cluster

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2026-06-01/mongoclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

var _ sdk.ResourceWithCustomizeDiff = Resource{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return mongoclusters.ValidateMongoClusterID
}

func (r Resource) ResourceType() string {
	return "azurerm_mongo_cluster"
}
