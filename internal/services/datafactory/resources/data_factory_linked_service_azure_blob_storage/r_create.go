// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_azure_blob_storage

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

func resourceDataFactoryLinkedServiceBlobStorageCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
					return fmt.Errorf("checking for presence of existing Data Factory Blob Storage Anonymous %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_data_factory_linked_service_azure_blob_storage", id.ID())
			}
		}
	}

	blobStorageProperties := &datafactory.AzureBlobStorageLinkedServiceTypeProperties{}

	if v, ok := d.GetOk("connection_string"); ok {
		blobStorageProperties.ConnectionString = &datafactory.SecureString{
			Value: pointer.To(v.(string)),
			Type:  datafactory.TypeSecureString,
		}
	}

	if v, ok := d.GetOk("connection_string_insecure"); ok {
		blobStorageProperties.ConnectionString = v.(string)
	}

	if v, ok := d.GetOk("sas_uri"); ok {
		blobStorageProperties.SasURI = &datafactory.SecureString{
			Value: pointer.To(v.(string)),
			Type:  datafactory.TypeSecureString,
		}

		if sasToken, ok := d.GetOk("sas_token_linked_key_vault_key"); ok {
			blobStorageProperties.SasToken = helpers.ExpandAzureKeyVaultSecretReference(sasToken.([]any))
			blobStorageProperties.SasURI = pointer.To(v.(string))
		}
	}

	if d.Get("use_managed_identity").(bool) {
		if v, ok := d.GetOk("service_endpoint"); ok {
			blobStorageProperties.ServiceEndpoint = pointer.To(v.(string))
		}
	} else {
		if v, ok := d.GetOk("service_endpoint"); ok {
			blobStorageProperties.ServiceEndpoint = pointer.To(v.(string))
		}
		if kvsp, ok := d.GetOk("service_principal_linked_key_vault_key"); ok {
			blobStorageProperties.ServicePrincipalKey = helpers.ExpandAzureKeyVaultSecretReference(kvsp.([]any))
		} else {
			secureString := datafactory.SecureString{
				Value: pointer.To(d.Get("service_principal_key").(string)),
				Type:  datafactory.TypeSecureString,
			}
			blobStorageProperties.ServicePrincipalKey = &secureString
		}

		blobStorageProperties.ServicePrincipalID = pointer.To(d.Get("service_principal_id").(string))
		blobStorageProperties.Tenant = pointer.To(d.Get("tenant_id").(string))
	}

	blobStorageLinkedService := &datafactory.AzureBlobStorageLinkedService{
		Description: pointer.To(d.Get("description").(string)),
		AzureBlobStorageLinkedServiceTypeProperties: blobStorageProperties,
		Type: datafactory.TypeBasicLinkedServiceTypeAzureBlobStorage,
	}

	if v, ok := d.GetOk("parameters"); ok {
		blobStorageLinkedService.Parameters = helpers.ExpandLinkedServiceParameters(v.(map[string]any))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		blobStorageLinkedService.ConnectVia = helpers.ExpandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("storage_kind"); ok {
		blobStorageLinkedService.AccountKind = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		blobStorageLinkedService.AdditionalProperties = v.(map[string]any)
	}

	if v, ok := d.GetOk("annotations"); ok {
		blobStorageLinkedService.Annotations = pointer.To(v.([]any))
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: blobStorageLinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, linkedService, ""); err != nil {
		return fmt.Errorf("creating/updating Data Factory Blob Storage %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceDataFactoryLinkedServiceBlobStorageRead(d, meta)
}
