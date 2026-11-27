// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"context"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

// ClientProvider decouples internal/sdk and resource subpackages from importing
// the monolithic internal/clients package (which transitively imports all 135
// service clients and >1,500 Azure SDK packages).
type ClientProvider interface {
	ClientAccount() *account.ResourceManagerAccount
	ClientStopContext() context.Context
	UserFeatures() features.UserFeatures
}
