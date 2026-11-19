// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_custom_certificate

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name web_pubsub_custom_certificate -service-package-name signalr -properties "name" -compare-values "subscription_id:web_pubsub_id,resource_group_name:web_pubsub_id,web_pubsub_name:web_pubsub_id"

type Resource struct{}

var (
	_ sdk.Resource             = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &webpubsub.CustomCertificateId{}
}

func (r Resource) ResourceType() string {
	return webPubsubCustomCertificateResourceType
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return webpubsub.ValidateCustomCertificateID
}
