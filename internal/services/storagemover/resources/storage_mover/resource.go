// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_mover

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/storagemovers"
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
	return &storagemovers.StorageMoverId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_storage_mover"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return storagemovers.ValidateStorageMoverID
}
