// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_provider_feature_registration

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2021-07-01/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name resource_provider_feature_registration -properties "name,provider_name" -known-values "subscription_id:data.Subscriptions.Secondary" -test-name "basicForResourceIdentity"

const (
	Pending       = "Pending"
	Registering   = "Registering"
	Unregistering = "Unregistering"
	Registered    = "Registered"
	NotRegistered = "NotRegistered"
	Unregistered  = "Unregistered"
)

var _ sdk.ResourceWithIdentity = Resource{}

type Resource struct{}

func (r Resource) Identity() resourceids.ResourceId {
	return new(features.FeatureId)
}

func (r Resource) ResourceType() string {
	return "azurerm_resource_provider_feature_registration"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return features.ValidateFeatureID
}
