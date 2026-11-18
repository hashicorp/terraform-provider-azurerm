// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
)

var (
	NewSubnetAndVnetPoller     = helpers.NewSubnetAndVnetPoller
	SubnetResourceName         = helpers.SubnetResourceName
	VirtualNetworkResourceName = helpers.VirtualNetworkResourceName
)
