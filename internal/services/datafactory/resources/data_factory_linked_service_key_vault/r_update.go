// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_key_vault

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryLinkedServiceKeyVaultCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	subscriptionId := meta.(*clients.Client).DataFactory.LinkedServiceClient.SubscriptionID
	keyVaultsClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewLinkedServiceID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	keyVaultIdRaw := d.Get("key_vault_id").(string)
	keyVaultBaseUri := pointer.To(d.Get("key_vault_base_url_dynamic_expression").(string))

	if keyVaultIdRaw != "" {
		keyVaultId, err := commonids.ParseKeyVaultID(keyVaultIdRaw)
		if err != nil {
			return fmt.Errorf("parsing Key Vault ID %q: %+v", keyVaultIdRaw, err)
		}

		keyVaultBaseUri, err = keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
		if err != nil {
			return fmt.Errorf("retrieving Key Vault Base URI for %s: %+v", *keyVaultId, err)
		}
	}

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing Data Factory Key Vault %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_data_factory_linked_service_key_vault", id.ID())
			}
		}
	}

	azureKeyVaultProperties := &datafactory.AzureKeyVaultLinkedServiceTypeProperties{
		BaseURL: keyVaultBaseUri,
	}

	azureKeyVaultLinkedService := &datafactory.AzureKeyVaultLinkedService{
		Description:                              pointer.To(d.Get("description").(string)),
		AzureKeyVaultLinkedServiceTypeProperties: azureKeyVaultProperties,
		Type:                                     datafactory.TypeBasicLinkedServiceTypeAzureKeyVault,
	}

	if v, ok := d.GetOk("parameters"); ok {
		azureKeyVaultLinkedService.Parameters = helpers.ExpandLinkedServiceParameters(v.(map[string]any))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		azureKeyVaultLinkedService.ConnectVia = helpers.ExpandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		azureKeyVaultLinkedService.AdditionalProperties = v.(map[string]any)
	}

	if v, ok := d.GetOk("annotations"); ok {
		azureKeyVaultLinkedService.Annotations = pointer.To(v.([]any))
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: azureKeyVaultLinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, linkedService, ""); err != nil {
		return fmt.Errorf("creating/updating Data Factory Key Vault %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryLinkedServiceKeyVaultRead(d, meta)
}
