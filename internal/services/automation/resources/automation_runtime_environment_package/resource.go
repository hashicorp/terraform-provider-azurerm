// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_runtime_environment_package

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/packageresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "name" -compare-values "subscription_id:automation_runtime_environment_id,resource_group_name:automation_runtime_environment_id,automation_account_name:automation_runtime_environment_id,runtime_environment_name:automation_runtime_environment_id" -test-expect-non-empty=true

var (
	_ sdk.ResourceWithIdentity = Resource{}
	_ sdk.Resource             = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &packageresource.PackageId{}
}

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_automation_runtime_environment_package"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return packageresource.ValidatePackageID
}
