// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package playwright_workspace

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2025-09-01/playwrightworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

type Resource struct{}

var _ sdk.ResourceWithIdentity = Resource{}

func (Resource) ResourceType() string {
	return "azurerm_playwright_workspace"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return playwrightworkspaces.ValidatePlaywrightWorkspaceID
}

func (Resource) Identity() resourceids.ResourceId {
	return &playwrightworkspaces.PlaywrightWorkspaceId{}
}
