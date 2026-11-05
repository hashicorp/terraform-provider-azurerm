// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_connection_custom_keys

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/accountconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name cognitive_account_connection_custom_keys -properties "name" -compare-values "subscription_id:cognitive_account_id,resource_group_name:cognitive_account_id,account_name:cognitive_account_id" -test-name "basic" -test-expect-non-empty

var (
	_ sdk.ResourceWithUpdate         = Resource{}
	_ sdk.ResourceWithIdentity       = Resource{}
	_ sdk.ResourceWithCustomImporter = Resource{}
)

type Resource struct{}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return helpers.CognitiveAccountConnectionImporter(accountconnectionresource.ConnectionAuthTypeCustomKeys, r.ResourceType())
}

func (r Resource) Identity() resourceids.ResourceId {
	return new(accountconnectionresource.ConnectionId)
}

func (r Resource) ResourceType() string {
	return "azurerm_cognitive_account_connection_custom_keys"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return accountconnectionresource.ValidateConnectionID
}
