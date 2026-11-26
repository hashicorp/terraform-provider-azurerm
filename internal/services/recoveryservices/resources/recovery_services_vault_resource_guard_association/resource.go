// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package recovery_services_vault_resource_guard_association

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/resourceguardproxy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = Resource{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return resourceguardproxy.ValidateBackupResourceGuardProxyID
}

func (r Resource) ResourceType() string {
	return "azurerm_recovery_services_vault_resource_guard_association"
}
