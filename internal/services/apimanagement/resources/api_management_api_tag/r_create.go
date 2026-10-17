// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_tag

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apitag"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/tag"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiTagCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	tagClient := meta.(*clients.Client).ApiManagement.TagClient
	client := meta.(*clients.Client).ApiManagement.ApiTagClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apiId, err := api.ParseApiID(d.Get("api_id").(string))
	if err != nil {
		return err
	}

	tagId := tag.NewTagID(subscriptionId, apiId.ResourceGroupName, apiId.ServiceName, d.Get("name").(string))

	id := apitag.NewApiTagID(subscriptionId, apiId.ResourceGroupName, apiId.ServiceName, apiId.ApiId, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		tagExists, err := tagClient.Get(ctx, tagId)
		if err != nil {
			if !response.WasNotFound(tagExists.HttpResponse) {
				return fmt.Errorf("checking for presence of Tag %q: %s", id, err)
			}
		}

		tagAssignmentExist, err := client.TagGetByApi(ctx, id)
		if err != nil {
			if !response.WasNotFound(tagAssignmentExist.HttpResponse) {
				return fmt.Errorf("checking for presence of Tag Assignment %q: %s", id, err)
			}
		}

		if !response.WasNotFound(tagAssignmentExist.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_api_tag", id.ID())
		}
	}

	if _, err := client.TagAssignToApi(ctx, id); err != nil {
		return fmt.Errorf("assigning to Api %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementApiTagRead(d, meta)
}
