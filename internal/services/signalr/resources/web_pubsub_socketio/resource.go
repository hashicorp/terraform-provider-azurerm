// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_socketio

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
	_ sdk.ResourceWithIdentity      = Resource{}
)

func (w Resource) Identity() resourceids.ResourceId {
	return &webpubsub.WebPubSubId{}
}

func (w Resource) ResourceType() string {
	return "azurerm_web_pubsub_socketio"
}

func (w Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return webpubsub.ValidateWebPubSubID
}
