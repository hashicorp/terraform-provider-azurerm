// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dashboard_grafana_managed_private_endpoint

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2025-08-01/managedprivateendpointmodels"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name dashboard_grafana_managed_private_endpoint -service-package-name dashboard -properties "name" -compare-values "subscription_id:grafana_id,resource_group_name:grafana_id,grafana_name:grafana_id"

type Resource struct{}

var (
	_ sdk.Resource             = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &managedprivateendpointmodels.ManagedPrivateEndpointId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_dashboard_grafana_managed_private_endpoint"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return managedprivateendpointmodels.ValidateManagedPrivateEndpointID
}
