// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_operation_policy

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apioperationpolicy"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementAPIOperationPolicyCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ApiOperationPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := apioperationpolicy.NewOperationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("api_name").(string), d.Get("operation_id").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id, apioperationpolicy.GetOperationOptions{Format: pointer.To(apioperationpolicy.PolicyExportFormatXml)})
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_api_management_api_operation_policy", id.ID())
			}
		}
	}

	parameters := apioperationpolicy.PolicyContract{}

	xmlContent := d.Get("xml_content").(string)
	xmlLink := d.Get("xml_link").(string)

	if xmlContent != "" {
		parameters.Properties = &apioperationpolicy.PolicyContractProperties{
			Format: pointer.To(apioperationpolicy.PolicyContentFormatRawxml),
			Value:  xmlContent,
		}
	}

	if xmlLink != "" {
		parameters.Properties = &apioperationpolicy.PolicyContractProperties{
			Format: pointer.To(apioperationpolicy.PolicyContentFormatRawxmlNegativelink),
			Value:  xmlLink,
		}
	}

	if parameters.Properties == nil {
		return errors.New("either `xml_content` or `xml_link` must be set")
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, apioperationpolicy.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementAPIOperationPolicyRead(d, meta)
}
