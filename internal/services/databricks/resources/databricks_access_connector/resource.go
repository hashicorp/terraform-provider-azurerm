// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package databricks_access_connector

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/accessconnector"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate   = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &accessconnector.AccessConnectorId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_databricks_access_connector"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return accessconnector.ValidateAccessConnectorID
}
