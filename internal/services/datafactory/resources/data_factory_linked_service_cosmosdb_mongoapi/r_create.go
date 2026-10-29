// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_cosmosdb_mongoapi

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryLinkedServiceCosmosDbMongoAPICreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	subscriptionId := meta.(*clients.Client).DataFactory.LinkedServiceClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewLinkedServiceID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing Data Factory Linked Service CosmosDb %q (Data Factory %q / Resource Group %q): %+v", id.Name, id.FactoryName, id.ResourceGroup, err)
				}
			}

			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_data_factory_linked_service_cosmosdb", id.ID())
			}
		}
	}

	cosmosdbProperties := &datafactory.CosmosDbMongoDbAPILinkedServiceTypeProperties{}

	cosmosdbProperties.ConnectionString = datafactory.SecureString{
		Value: pointer.To(d.Get("connection_string").(string)),
		Type:  datafactory.TypeSecureString,
	}
	cosmosdbProperties.Database = d.Get("database").(string)
	cosmosdbProperties.IsServerVersionAbove32 = d.Get("server_version_is_32_or_higher").(bool)

	cosmosdbLinkedService := &datafactory.CosmosDbMongoDbAPILinkedService{
		Description: pointer.To(d.Get("description").(string)),
		CosmosDbMongoDbAPILinkedServiceTypeProperties: cosmosdbProperties,
		Type: datafactory.TypeBasicLinkedServiceTypeCosmosDbMongoDbAPI,
	}

	if v, ok := d.GetOk("parameters"); ok {
		cosmosdbLinkedService.Parameters = helpers.ExpandLinkedServiceParameters(v.(map[string]any))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		cosmosdbLinkedService.ConnectVia = helpers.ExpandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		cosmosdbLinkedService.AdditionalProperties = v.(map[string]any)
	}

	if v, ok := d.GetOk("annotations"); ok {
		cosmosdbLinkedService.Annotations = pointer.To(v.([]any))
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: cosmosdbLinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, linkedService, ""); err != nil {
		return fmt.Errorf("creating/updating Data Factory Linked Service CosmosDb %q (Data Factory %q / Resource Group %q): %+v", id.Name, id.FactoryName, id.ResourceGroup, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryLinkedServiceCosmosDbMongoAPIRead(d, meta)
}
