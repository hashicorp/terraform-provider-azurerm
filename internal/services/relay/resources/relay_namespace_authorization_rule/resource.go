// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package relay_namespace_authorization_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_relay_namespace_authorization_rule"
}

func (r Resource) Description() string {
	return "Manages an Azure Relay Namespace Authorization Rule."
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return namespaces.ValidateAuthorizationRuleID
}
