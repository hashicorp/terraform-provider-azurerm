// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace_schema_group

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2026-01-01/schemagroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceEventHubNamespaceSchemaRegistryCreate,
		Read:   resourceEventHubNamespaceSchemaRegistryRead,
		Delete: resourceEventHubNamespaceSchemaRegistryDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := schemagroups.ParseSchemaGroupID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: eventhubNamespaceSchemaGroupSchema(),
	}
}
