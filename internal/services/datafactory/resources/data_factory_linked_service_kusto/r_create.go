// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_kusto

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

func resourceDataFactoryLinkedServiceKustoCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
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
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_data_factory_linked_service_kusto", id.ID())
			}
		}
	}

	kustoLinkedService := &datafactory.AzureDataExplorerLinkedService{
		AzureDataExplorerLinkedServiceTypeProperties: &datafactory.AzureDataExplorerLinkedServiceTypeProperties{
			Endpoint: d.Get("kusto_endpoint").(string),
			Database: d.Get("kusto_database_name").(string),
		},
		Description: pointer.To(d.Get("description").(string)),
		Type:        datafactory.TypeBasicLinkedServiceTypeAzureDataExplorer,
	}

	if d.Get("use_managed_identity").(bool) {
		kustoLinkedService.AzureDataExplorerLinkedServiceTypeProperties = &datafactory.AzureDataExplorerLinkedServiceTypeProperties{
			Endpoint: d.Get("kusto_endpoint").(string),
			Database: d.Get("kusto_database_name").(string),
		}
	} else if v, ok := d.GetOk("service_principal_id"); ok {
		kustoLinkedService.AzureDataExplorerLinkedServiceTypeProperties = &datafactory.AzureDataExplorerLinkedServiceTypeProperties{
			Endpoint:           d.Get("kusto_endpoint").(string),
			Database:           d.Get("kusto_database_name").(string),
			ServicePrincipalID: v.(string),
			ServicePrincipalKey: &datafactory.SecureString{
				Value: pointer.To(d.Get("service_principal_key").(string)),
				Type:  datafactory.TypeSecureString,
			},
			Tenant: pointer.To(d.Get("tenant").(string)),
		}
	} else {
		return fmt.Errorf("one of Managed Identity and service principal authentication must be set")
	}

	if v, ok := d.GetOk("parameters"); ok {
		kustoLinkedService.Parameters = helpers.ExpandLinkedServiceParameters(v.(map[string]any))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		kustoLinkedService.ConnectVia = helpers.ExpandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		kustoLinkedService.AdditionalProperties = v.(map[string]any)
	}

	if v, ok := d.GetOk("annotations"); ok {
		kustoLinkedService.Annotations = pointer.To(v.([]any))
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: kustoLinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, linkedService, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryLinkedServiceKustoRead(d, meta)
}
