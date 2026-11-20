// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_namespace_topic_id_association

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/namespacetopics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"eventgrid_namespace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: namespaces.ValidateNamespaceID,
		},

		"eventgrid_namespace_topic_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			// `ForceNew` as `eventgrid_namespace_topic_id` can only be assigned to `eventgrid_namespace_id` of same namespace and `eventgrid_namespace_topic_id` is unique for each namespace, changing `eventgrid_namespace_topic_id` means that `eventgrid_namespace_id` is changed too
			ForceNew:     true,
			ValidateFunc: namespacetopics.ValidateNamespaceTopicID,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
