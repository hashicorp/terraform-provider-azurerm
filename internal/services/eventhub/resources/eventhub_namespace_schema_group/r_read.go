// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace_schema_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/schemaregistry"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventHubNamespaceSchemaRegistryRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Eventhub.SchemaRegistryClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := schemaregistry.ParseSchemaGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.Set("name", id.SchemaGroupName)

	namespaceId := namespaces.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
	d.Set("namespace_id", namespaceId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.SchemaCompatibility != nil {
				d.Set("schema_compatibility", string(*props.SchemaCompatibility))
			}
			if props.SchemaType != nil {
				d.Set("schema_type", string(*props.SchemaType))
			}
		}
	}

	return nil
}
