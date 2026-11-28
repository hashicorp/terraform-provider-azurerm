// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package video_indexer_account

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/videoindexer/2025-04-01/accounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate   = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &accounts.AccountId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_video_indexer_account"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return accounts.ValidateAccountID
}
