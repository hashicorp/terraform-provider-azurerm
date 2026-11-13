// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_workspace

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/monitor/2023-04-03/azuremonitorworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

type Resource struct{}

var (
	_ sdk.ResourceWithIdentity = Resource{}
	_ sdk.ResourceWithUpdate   = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &azuremonitorworkspaces.AccountId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_monitor_workspace"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return azuremonitorworkspaces.ValidateAccountID
}
