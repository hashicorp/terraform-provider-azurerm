// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_table_properties

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -parent-id "storage_account_id" -test-name "corsOnly"

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate               = Resource{}
	_ sdk.ResourceWithIdentityTypeOverride = Resource{}
)

func (s Resource) ResourceType() string {
	return "azurerm_storage_account_table_properties"
}

func (s Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateStorageAccountID
}

func (s Resource) Identity() resourceids.ResourceId {
	return &commonids.StorageAccountId{}
}

func (s Resource) IdentityType() pluginsdk.ResourceTypeForIdentity {
	return pluginsdk.ResourceTypeForIdentityVirtual
}
