// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace_schema_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/schemaregistry"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	eventhubClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventHubNamespaceSchemaRegistryCreate(d *pluginsdk.ResourceData, meta any) error {
	client := eventhubClient.From(meta).SchemaRegistryClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	namespaceId, err := namespaces.ParseNamespaceID(d.Get("namespace_id").(string))
	if err != nil {
		return fmt.Errorf("parsing eventhub namespace %s error: %+v", namespaceId.ID(), err)
	}

	id := schemaregistry.NewSchemaGroupID(subscriptionId, namespaceId.ResourceGroupName, namespaceId.NamespaceName, d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_eventhub_namespace_schema_group", id.ID())
		}
	}

	schemaCompatibilityType := schemaregistry.SchemaCompatibility(d.Get("schema_compatibility").(string))
	schemaType := schemaregistry.SchemaType(d.Get("schema_type").(string))

	parameters := schemaregistry.SchemaGroup{
		Properties: &schemaregistry.SchemaGroupProperties{
			SchemaCompatibility: &schemaCompatibilityType,
			SchemaType:          &schemaType,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceEventHubNamespaceSchemaRegistryRead(d, meta)
}
