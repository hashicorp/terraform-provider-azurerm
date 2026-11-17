// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package databox_edge_device

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databoxedge/2022-03-01/devices"
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
	return &devices.DataBoxEdgeDeviceId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_databox_edge_device"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return devices.ValidateDataBoxEdgeDeviceID
}
