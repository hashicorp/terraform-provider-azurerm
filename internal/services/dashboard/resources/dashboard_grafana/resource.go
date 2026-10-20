// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dashboard_grafana

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2025-08-01/managedgrafanas"
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
	return &managedgrafanas.GrafanaId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_dashboard_grafana"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return managedgrafanas.ValidateGrafanaID
}
