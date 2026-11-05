// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_virtual_hard_disk

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/virtualharddisks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.Resource           = Resource{}
	_ sdk.ResourceWithUpdate = Resource{}
)

type Resource struct{}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return virtualharddisks.ValidateVirtualHardDiskID
}

func (Resource) ResourceType() string {
	return "azurerm_stack_hci_virtual_hard_disk"
}
