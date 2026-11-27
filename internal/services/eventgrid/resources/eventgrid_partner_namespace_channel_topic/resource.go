// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_namespace_channel_topic

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/channels"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name eventgrid_partner_namespace_channel_topic -properties name -compare-values subscription_id:partner_namespace_id,resource_group_name:partner_namespace_id,partner_namespace_name:partner_namespace_id

var (
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithIdentity      = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
)

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_eventgrid_partner_namespace_channel_topic"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return channels.ValidateChannelID
}

func (r Resource) Identity() resourceids.ResourceId {
	return &channels.ChannelId{}
}
