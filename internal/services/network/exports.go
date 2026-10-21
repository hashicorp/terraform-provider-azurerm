// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/subnet"
)

var (
	NewSubnetAndVnetPoller     = helpers.NewSubnetAndVnetPoller
	SubnetResourceName         = subnet.SubnetResourceName
	VirtualNetworkResourceName = subnet.VirtualNetworkResourceName
)
