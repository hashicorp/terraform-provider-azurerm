// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package trusted_signing_account

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13/codesigningaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = (*Resource)(nil)

func (m Resource) ResourceType() string {
	return "azurerm_trusted_signing_account"
}

func (m Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return codesigningaccounts.ValidateCodeSigningAccountID
}
