// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_tag_description

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apitag"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apitagdescription"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiTagDescriptionCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ApiTagDescriptionClient

	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apiTagId, err := apitag.ParseApiTagID(d.Get("api_tag_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `api_id`: %v", err)
	}

	apiName := helpers.GetApiName(apiTagId.ApiId)

	id := apitagdescription.NewTagDescriptionID(apiTagId.SubscriptionId, apiTagId.ResourceGroupName, apiTagId.ServiceName, apiName, apiTagId.TagId)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_api_management_api_tag_description", id.ID())
			}
		}
	}

	tagDescParameter := apitagdescription.TagDescriptionCreateParameters{Properties: &apitagdescription.TagDescriptionBaseProperties{}}
	if v, ok := d.GetOk("description"); ok {
		tagDescParameter.Properties.Description = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("external_documentation_url"); ok {
		tagDescParameter.Properties.ExternalDocsURL = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("external_documentation_description"); ok {
		tagDescParameter.Properties.ExternalDocsDescription = pointer.To(v.(string))
	}

	if _, err := client.CreateOrUpdate(ctx, id, tagDescParameter, apitagdescription.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementApiTagDescriptionRead(d, meta)
}
