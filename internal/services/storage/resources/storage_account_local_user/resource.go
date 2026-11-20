// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_local_user

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/localuseroperationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name storage_account_local_user -service-package-name storage -properties "name" -compare-values "subscription_id:storage_account_id,resource_group_name:storage_account_id,storage_account_name:storage_account_id" -test-name "passwordOnly"

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate   = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_storage_account_local_user"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return localuseroperationgroup.ValidateLocalUserID
}

func (r Resource) Identity() resourceids.ResourceId {
	return &localuseroperationgroup.LocalUserId{}
}
