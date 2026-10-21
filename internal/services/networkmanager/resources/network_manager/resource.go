// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkmanagers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -test-sequential

var (
	_ sdk.ResourceWithUpdate   = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

type Resource struct{}

func (r Resource) Identity() resourceids.ResourceId {
	return &networkmanagers.NetworkManagerId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_network_manager"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return networkmanagers.ValidateNetworkManagerID
}
