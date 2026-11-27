// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_static_website

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -parent-id "storage_account_id" -test-name "complete"

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate               = Resource{}
	_ sdk.ResourceWithIdentityTypeOverride = Resource{}
)

func (a Resource) ResourceType() string {
	return "azurerm_storage_account_static_website"
}

func (a Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateStorageAccountID
}

func (a Resource) Identity() resourceids.ResourceId {
	return &commonids.StorageAccountId{}
}

func (a Resource) IdentityType() pluginsdk.ResourceTypeForIdentity {
	return pluginsdk.ResourceTypeForIdentityVirtual
}
