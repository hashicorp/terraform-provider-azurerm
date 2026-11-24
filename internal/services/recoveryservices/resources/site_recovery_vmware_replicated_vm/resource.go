// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_vmware_replicated_vm

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_site_recovery_vmware_replicated_vm"
}

func (r Resource) Description() string {
	return "Manages a VMWare replicated VM using Azure Site Recovery (VMWare to Azure only). A replicated VM keeps a copiously updated image of the VM in Azure in order to be able to start the VM in Azure in case of a disaster."
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return replicationprotecteditems.ValidateReplicationProtectedItemID
}

var (
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
)
