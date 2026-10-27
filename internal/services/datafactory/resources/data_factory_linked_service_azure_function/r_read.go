// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_azure_function

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryLinkedServiceAzureFunctionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LinkedServiceID(d.Id())
	if err != nil {
		return err
	}

	dataFactoryId := factories.NewFactoryID(id.SubscriptionId, id.ResourceGroup, id.FactoryName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	azureFunction, ok := resp.Properties.AsAzureFunctionLinkedService()
	if !ok {
		return fmt.Errorf("classifying Data Factory Linked Service Azure Function %s: Expected: %q Received: %q", id, datafactory.TypeBasicLinkedServiceTypeAzureFunction, *resp.Type)
	}

	d.Set("url", azureFunction.FunctionAppURL)

	d.Set("additional_properties", azureFunction.AdditionalProperties)
	d.Set("description", azureFunction.Description)

	if functionKey := azureFunction.FunctionKey; functionKey != nil {
		if keyVaultKey, ok := functionKey.AsAzureKeyVaultSecretReference(); !ok {
			return fmt.Errorf("classifying Data Factory Azure Key Vault Secret %s: Expected: %q Received: %q", id, datafactory.TypeAzureKeyVaultSecret, keyVaultKey.Type)
		} else if err := d.Set("key_vault_key", helpers.FlattenAzureKeyVaultSecretReference(keyVaultKey)); err != nil {
			return fmt.Errorf("setting `key_vault_key`: %+v", err)
		}
	}

	if err := d.Set("annotations", helpers.FlattenDataFactoryAnnotations(azureFunction.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations` for Data Factory Linked Service Azure Function %q (Data Factory %q) / Resource Group %q): %+v", id.Name, id.FactoryName, id.ResourceGroup, err)
	}

	if err := d.Set("parameters", helpers.FlattenLinkedServiceParameters(azureFunction.Parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	if connectVia := azureFunction.ConnectVia; connectVia != nil {
		if connectVia.ReferenceName != nil {
			d.Set("integration_runtime_name", connectVia.ReferenceName)
		}
	}

	return nil
}
