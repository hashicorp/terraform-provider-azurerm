// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_web

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

func resourceDataFactoryLinkedServiceWebCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
					return fmt.Errorf("checking for presence of existing Data Factory Web Anonymous %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_data_factory_linked_service_web", id.ID())
			}
		}
	}

	description := d.Get("description").(string)

	webLinkedService := &datafactory.WebLinkedService{
		Description: &description,
		Type:        datafactory.TypeBasicLinkedServiceTypeWeb,
	}

	url := d.Get("url").(string)
	authenticationType := d.Get("authentication_type").(string)

	if authenticationType == "Anonymous" {
		webLinkedService.TypeProperties = &datafactory.WebAnonymousAuthentication{
			AuthenticationType: datafactory.AuthenticationType(authenticationType),
			URL:                pointer.To(url),
		}
	}

	if authenticationType == "Basic" {
		username := d.Get("username").(string)
		password := d.Get("password").(string)
		passwordSecureString := datafactory.SecureString{
			Value: &password,
			Type:  datafactory.TypeSecureString,
		}
		webLinkedService.TypeProperties = &datafactory.WebBasicAuthentication{
			AuthenticationType: datafactory.AuthenticationType(authenticationType),
			URL:                pointer.To(url),
			Username:           username,
			Password:           passwordSecureString,
		}
	}

	if v, ok := d.GetOk("parameters"); ok {
		webLinkedService.Parameters = helpers.ExpandLinkedServiceParameters(v.(map[string]any))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		webLinkedService.ConnectVia = helpers.ExpandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		webLinkedService.AdditionalProperties = v.(map[string]any)
	}

	if v, ok := d.GetOk("annotations"); ok {
		webLinkedService.Annotations = pointer.To(v.([]any))
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: webLinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, linkedService, ""); err != nil {
		return fmt.Errorf("creating/updating Data Factory Web Anonymous %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryLinkedServiceWebRead(d, meta)
}
