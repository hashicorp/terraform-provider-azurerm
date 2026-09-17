// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package access_policy_assignment

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_managed_redis_access_policy_assignment"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return databases.ValidateAccessPolicyAssignmentID
}
