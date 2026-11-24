// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_namespace_topic

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/namespacetopics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name eventgrid_namespace_topic -properties "name" -compare-values "subscription_id:eventgrid_namespace_id,resource_group_name:eventgrid_namespace_id,namespace_name:eventgrid_namespace_id"

var (
	_ sdk.ResourceWithUpdate   = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_eventgrid_namespace_topic"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return namespacetopics.ValidateNamespaceTopicID
}

func (r Resource) Identity() resourceids.ResourceId {
	return new(namespacetopics.NamespaceTopicId)
}
