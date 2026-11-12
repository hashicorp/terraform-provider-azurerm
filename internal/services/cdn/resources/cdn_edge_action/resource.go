// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_edge_action

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

type Resource struct{}

var (
	_ sdk.ResourceWithIdentity = Resource{}
	_ sdk.ResourceWithUpdate   = Resource{}
)

func (Resource) ResourceType() string {
	return "azurerm_cdn_edge_action"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return edgeactions.ValidateEdgeActionID
}

func (Resource) Identity() resourceids.ResourceId {
	return &edgeactions.EdgeActionId{}
}
