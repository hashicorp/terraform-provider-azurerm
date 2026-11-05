// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_hybrid_runbook_worker_group

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/hybridrunbookworkergroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = (*Resource)(nil)

func (m Resource) ResourceType() string {
	return "azurerm_automation_hybrid_runbook_worker_group"
}

func (m Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return hybridrunbookworkergroup.ValidateHybridRunbookWorkerGroupID
}
