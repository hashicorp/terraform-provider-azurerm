// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_rai_policy

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/raipolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.ResourceWithUpdate        = &Resource{}
	_ sdk.ResourceWithCustomizeDiff = &Resource{}
)

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_cognitive_account_rai_policy"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return raipolicies.ValidateAccountRaiPolicyID
}
