// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_policy_fragment

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/policyfragment"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementPolicyFragmentCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.PolicyFragmentClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apiManagementId, err := apimanagementservice.ParseServiceID(d.Get("api_management_id").(string))
	if err != nil {
		return err
	}
	id := policyfragment.NewPolicyFragmentID(apiManagementId.SubscriptionId, apiManagementId.ResourceGroupName, apiManagementId.ServiceName, d.Get("name").(string))
	format := policyfragment.PolicyFragmentContentFormat(d.Get("format").(string))

	opts := policyfragment.DefaultGetOperationOptions()
	opts.Format = &format

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, opts)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_policy_fragment", id.ID())
		}
	}

	description := d.Get("description").(string)
	value := d.Get("value").(string)

	parameters := policyfragment.PolicyFragmentContract{
		Properties: &policyfragment.PolicyFragmentContractProperties{
			Description: pointer.To(description),
			Format:      pointer.To(format),
			Value:       value,
		},
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, policyfragment.CreateOrUpdateOperationOptions{}, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementPolicyFragmentRead(d, meta)
}
