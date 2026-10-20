// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_workspace_network_outbound_rule_private_endpoint

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/managednetwork"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagedNetwork struct {
	IsolationMode string `tfschema:"isolation_mode"`
}

type Resource struct{}

var _ sdk.Resource = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return managednetwork.ValidateOutboundRuleID
}
