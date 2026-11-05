// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_rai_blocklist

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/raiblocklists"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Resource struct{}

var _ sdk.Resource = Resource{}

func (c Resource) IDValidationFunc() func(any, string) ([]string, []error) {
	return raiblocklists.ValidateRaiBlocklistID
}

func (c Resource) ResourceType() string {
	return "azurerm_cognitive_account_rai_blocklist"
}
