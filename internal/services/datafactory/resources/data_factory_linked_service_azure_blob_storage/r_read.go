// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_azure_blob_storage

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

func resourceDataFactoryLinkedServiceBlobStorageRead(d *pluginsdk.ResourceData, meta any) error {
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

		return fmt.Errorf("retrieving Data Factory Blob Storage %s: %+v", *id, err)
	}

	d.Set("name", resp.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	blobStorage, ok := resp.Properties.AsAzureBlobStorageLinkedService()
	if !ok {
		return fmt.Errorf("classifying Data Factory Blob Storage %s: Expected: %q Received: %q", *id, datafactory.TypeBasicLinkedServiceTypeAzureBlobStorage, *resp.Type)
	}

	if blobStorage != nil {
		if blobStorage.Tenant != nil {
			d.Set("tenant_id", blobStorage.Tenant)
		}

		if blobStorage.ServicePrincipalID != nil {
			d.Set("service_principal_id", blobStorage.ServicePrincipalID)
			d.Set("use_managed_identity", false)
		} else {
			d.Set("service_endpoint", blobStorage.ServiceEndpoint)
			d.Set("use_managed_identity", true)
		}

		// blobStorage.ConnectionString is returned as a String when using `connection_string_insecure` and SecureString when using `connection_string`
		if insecureConnectionString, ok := blobStorage.ConnectionString.(string); ok {
			d.Set("connection_string_insecure", insecureConnectionString)
		}
	}

	if properties := blobStorage.AzureBlobStorageLinkedServiceTypeProperties; properties != nil {
		d.Set("storage_kind", properties.AccountKind)
		if sasToken := properties.SasToken; sasToken != nil {
			if keyVaultPassword, ok := sasToken.AsAzureKeyVaultSecretReference(); ok {
				if err := d.Set("sas_token_linked_key_vault_key", helpers.FlattenAzureKeyVaultSecretReference(keyVaultPassword)); err != nil {
					return fmt.Errorf("setting `sas_token_linked_key_vault_key`: %+v", err)
				}
			}
		}

		if spKey := properties.ServicePrincipalKey; spKey != nil {
			if kvSPkey, ok := spKey.AsAzureKeyVaultSecretReference(); ok {
				if err := d.Set("service_principal_linked_key_vault_key", helpers.FlattenAzureKeyVaultSecretReference(kvSPkey)); err != nil {
					return fmt.Errorf("setting `service_principal_linked_key_vault_key`: %+v", err)
				}
			}
		}
	}

	d.Set("additional_properties", blobStorage.AdditionalProperties)
	d.Set("description", blobStorage.Description)

	if err := d.Set("annotations", helpers.FlattenDataFactoryAnnotations(blobStorage.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations` for Data Factory Azure Blob Storage %s: %+v", *id, err)
	}

	if err := d.Set("parameters", helpers.FlattenLinkedServiceParameters(blobStorage.Parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	if connectVia := blobStorage.ConnectVia; connectVia != nil {
		if connectVia.ReferenceName != nil {
			d.Set("integration_runtime_name", connectVia.ReferenceName)
		}
	}

	return nil
}
