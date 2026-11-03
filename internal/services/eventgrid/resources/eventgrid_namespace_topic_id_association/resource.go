// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_namespace_topic_id_association

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.ResourceWithIdentity = Resource{}

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_eventgrid_namespace_topic_id_association"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return namespaces.ValidateNamespaceID
}

func (r Resource) Identity() resourceids.ResourceId {
	return &namespaces.NamespaceId{}
}
