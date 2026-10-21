// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_azure_file_storage

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

func resourceDataFactoryLinkedServiceAzureFileStorageRead(d *pluginsdk.ResourceData, meta any) error {
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

		return fmt.Errorf("retrieving Data Factory Azure File Storage %s: %+v", *id, err)
	}

	d.Set("name", resp.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	fileStorage, ok := resp.Properties.AsAzureFileStorageLinkedService()
	if !ok {
		return fmt.Errorf("classifying Data Factory Azure File Storage %s: Expected: %q Received: %q", *id, datafactory.TypeBasicLinkedServiceTypeAzureFileStorage, *resp.Type)
	}

	d.Set("additional_properties", fileStorage.AdditionalProperties)
	d.Set("description", fileStorage.Description)

	if password := fileStorage.Password; password != nil {
		if keyVaultPassword, ok := password.AsAzureKeyVaultSecretReference(); ok {
			if err := d.Set("key_vault_password", helpers.FlattenAzureKeyVaultSecretReference(keyVaultPassword)); err != nil {
				return fmt.Errorf("setting `key_vault_password`: %+v", err)
			}
		}
	}

	if fileStorage.UserID != nil {
		d.Set("user_id", fileStorage.UserID.(string))
	}

	if err := d.Set("annotations", helpers.FlattenDataFactoryAnnotations(fileStorage.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations` for Data Factory Azure File Storage %s: %+v", *id, err)
	}

	if err := d.Set("parameters", helpers.FlattenLinkedServiceParameters(fileStorage.Parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	if connectVia := fileStorage.ConnectVia; connectVia != nil {
		if connectVia.ReferenceName != nil {
			d.Set("integration_runtime_name", connectVia.ReferenceName)
		}
	}

	if props := fileStorage.AzureFileStorageLinkedServiceTypeProperties; props != nil {
		if props.FileShare != nil {
			d.Set("file_share", props.FileShare)
		}
	}

	return nil
}
