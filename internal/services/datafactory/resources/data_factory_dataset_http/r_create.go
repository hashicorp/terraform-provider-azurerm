// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_dataset_http

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

func resourceDataFactoryDatasetHTTPCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.DatasetClient
	subscriptionId := meta.(*clients.Client).DataFactory.DatasetClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewDataSetID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_data_factory_dataset_http", id.ID())
			}
		}
	}

	httpDatasetProperties := datafactory.HTTPDatasetTypeProperties{
		RelativeURL:   d.Get("relative_url").(string),
		RequestBody:   d.Get("request_body").(string),
		RequestMethod: d.Get("request_method").(string),
	}

	linkedServiceName := d.Get("linked_service_name").(string)
	linkedService := &datafactory.LinkedServiceReference{
		ReferenceName: &linkedServiceName,
		Type:          pointer.To("LinkedServiceReference"),
	}

	description := d.Get("description").(string)
	httpTableset := datafactory.HTTPDataset{
		HTTPDatasetTypeProperties: &httpDatasetProperties,
		LinkedServiceName:         linkedService,
		Description:               &description,
	}

	if v, ok := d.GetOk("folder"); ok {
		httpTableset.Folder = &datafactory.DatasetFolder{
			Name: pointer.To(v.(string)),
		}
	}

	if v, ok := d.GetOk("parameters"); ok {
		httpTableset.Parameters = helpers.ExpandDataSetParameters(v.(map[string]any))
	}

	if v, ok := d.GetOk("annotations"); ok {
		httpTableset.Annotations = pointer.To(v.([]any))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		httpTableset.AdditionalProperties = v.(map[string]any)
	}

	if v, ok := d.GetOk("schema_column"); ok {
		httpTableset.Structure = helpers.ExpandDataFactoryDatasetStructure(v.([]any))
	}

	dataset := datafactory.DatasetResource{
		Properties: &httpTableset,
		Type:       pointer.To(string(datafactory.TypeBasicDatasetTypeHTTPFile)),
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, dataset, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceDataFactoryDatasetHTTPRead(d, meta)
}
