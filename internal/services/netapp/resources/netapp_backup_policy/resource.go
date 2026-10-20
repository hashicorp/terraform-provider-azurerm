// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_backup_policy

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/backuppolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_netapp_backup_policy"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backuppolicies.ValidateBackupPolicyID
}
