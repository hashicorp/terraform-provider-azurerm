// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_data_lake_storage_gen2

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

func resourceDataFactoryLinkedServiceDataLakeStorageGen2CreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
					return fmt.Errorf("checking for presence of existing Data Factory Data Lake Storage Gen2 %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_data_factory_linked_service_data_lake_storage_gen2", id.ID())
			}
		}
	}

	var datalakeStorageGen2Properties *datafactory.AzureBlobFSLinkedServiceTypeProperties

	if d.Get("use_managed_identity").(bool) {
		datalakeStorageGen2Properties = &datafactory.AzureBlobFSLinkedServiceTypeProperties{
			URL: pointer.To(d.Get("url").(string)),
		}
	} else if v, ok := d.GetOk("storage_account_key"); ok {
		datalakeStorageGen2Properties = &datafactory.AzureBlobFSLinkedServiceTypeProperties{
			URL: pointer.To(d.Get("url").(string)),
			AccountKey: datafactory.SecureString{
				Value: pointer.To(v.(string)),
				Type:  datafactory.TypeSecureString,
			},
		}
	} else {
		secureString := datafactory.SecureString{
			Value: pointer.To(d.Get("service_principal_key").(string)),
			Type:  datafactory.TypeSecureString,
		}

		datalakeStorageGen2Properties = &datafactory.AzureBlobFSLinkedServiceTypeProperties{
			URL:                 pointer.To(d.Get("url").(string)),
			ServicePrincipalID:  pointer.To(d.Get("service_principal_id").(string)),
			Tenant:              pointer.To(d.Get("tenant").(string)),
			ServicePrincipalKey: &secureString,
		}
	}

	datalakeStorageGen2LinkedService := &datafactory.AzureBlobFSLinkedService{
		Description:                            pointer.To(d.Get("description").(string)),
		AzureBlobFSLinkedServiceTypeProperties: datalakeStorageGen2Properties,
		Type:                                   datafactory.TypeBasicLinkedServiceTypeAzureBlobFS,
	}

	if v, ok := d.GetOk("parameters"); ok {
		datalakeStorageGen2LinkedService.Parameters = helpers.ExpandLinkedServiceParameters(v.(map[string]any))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		datalakeStorageGen2LinkedService.ConnectVia = helpers.ExpandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		datalakeStorageGen2LinkedService.AdditionalProperties = v.(map[string]any)
	}

	if v, ok := d.GetOk("annotations"); ok {
		datalakeStorageGen2LinkedService.Annotations = pointer.To(v.([]any))
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: datalakeStorageGen2LinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, linkedService, ""); err != nil {
		return fmt.Errorf("creating/updating Data Factory Data Lake Storage Gen2 %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryLinkedServiceDataLakeStorageGen2Read(d, meta)
}
