// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_network_anchor

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/networkanchors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

var (
	_ sdk.ResourceWithUpdate   = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

type Resource struct{}

func (Resource) ResourceType() string {
	return "azurerm_oracle_network_anchor"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return networkanchors.ValidateNetworkAnchorID
}

func (Resource) Identity() resourceids.ResourceId {
	return &networkanchors.NetworkAnchorId{}
}
