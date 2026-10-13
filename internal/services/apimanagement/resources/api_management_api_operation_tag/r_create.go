// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_operation_tag

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apioperationtag"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/tag"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiOperationTagCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	tagClient := meta.(*clients.Client).ApiManagement.TagClient
	client := meta.(*clients.Client).ApiManagement.ApiOperationTagClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apiOperationId, err := apioperationtag.ParseOperationID(d.Get("api_operation_id").(string))
	if err != nil {
		return err
	}

	apiName := helpers.GetApiName(apiOperationId.ApiId)

	id := apioperationtag.NewOperationTagID(subscriptionId, apiOperationId.ResourceGroupName, apiOperationId.ServiceName, apiName, apiOperationId.OperationId, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.TagGetByOperation(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing Tag %q: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_api_management_api_operation_tag", id.ID())
			}
		}
	}

	parameters := tag.TagCreateUpdateParameters{
		Properties: &tag.TagContractProperties{
			DisplayName: d.Get("display_name").(string),
		},
	}

	tagId := tag.NewTagID(subscriptionId, apiOperationId.ResourceGroupName, apiOperationId.ServiceName, d.Get("name").(string))
	if _, err := tagClient.CreateOrUpdate(ctx, tagId, parameters, tag.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	if _, err := client.TagAssignToOperation(ctx, id); err != nil {
		return fmt.Errorf("assigning to operation %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementApiOperationTagRead(d, meta)
}
