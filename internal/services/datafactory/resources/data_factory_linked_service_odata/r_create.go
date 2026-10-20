// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_odata

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

func resourceArmDataFactoryLinkedServiceODataCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
					return fmt.Errorf("checking for presence of existing Data Factory OData Anonymous %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_data_factory_linked_service_odata", id.ID())
			}
		}
	}

	odataLinkedService := &datafactory.ODataLinkedService{
		Description: pointer.To(d.Get("description").(string)),
		Type:        datafactory.TypeBasicLinkedServiceTypeOData,
		ODataLinkedServiceTypeProperties: &datafactory.ODataLinkedServiceTypeProperties{
			AuthenticationType: datafactory.ODataAuthenticationTypeAnonymous,
			URL:                d.Get("url").(string),
		},
	}

	// There are multiple authentication paths. If support for those get added, we can easily add them in
	// a similar format to the below while not messing up the other attributes in ODataLinkedServiceTypeProperties
	if v, ok := d.GetOk("basic_authentication"); ok {
		attrs := v.([]any)
		if len(attrs) != 0 && attrs[0] != nil {
			raw := attrs[0].(map[string]any)
			odataLinkedService.AuthenticationType = datafactory.ODataAuthenticationTypeBasic
			odataLinkedService.UserName = raw["username"].(string)
			odataLinkedService.Password = datafactory.SecureString{
				Value: pointer.To(raw["password"].(string)),
				Type:  datafactory.TypeSecureString,
			}
		}
	}

	if v, ok := d.GetOk("parameters"); ok {
		odataLinkedService.Parameters = helpers.ExpandLinkedServiceParameters(v.(map[string]any))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		odataLinkedService.ConnectVia = helpers.ExpandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		odataLinkedService.AdditionalProperties = v.(map[string]any)
	}

	if v, ok := d.GetOk("annotations"); ok {
		odataLinkedService.Annotations = pointer.To(v.([]any))
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: odataLinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, linkedService, ""); err != nil {
		return fmt.Errorf("creating/updating Data Factory OData Anonymous %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmDataFactoryLinkedServiceODataRead(d, meta)
}
