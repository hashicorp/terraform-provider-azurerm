// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_policy_fragment

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/policyfragment"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementPolicyFragmentRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.PolicyFragmentClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := policyfragment.ParsePolicyFragmentID(d.Id())
	if err != nil {
		return err
	}

	format := policyfragment.PolicyFragmentContentFormat(d.Get("format").(string))
	opts := policyfragment.DefaultGetOperationOptions()
	opts.Format = &format
	resp, err := client.Get(ctx, *id, opts)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s does not exist - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.PolicyFragmentName)
	apiManagementId := apimanagementservice.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName)
	d.Set("api_management_id", apiManagementId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", pointer.From(props.Description))
			d.Set("value", props.Value)
			// the api only returns a format field when it's requested in the GET request as param '?format=rawxml' and only does so when it's not 'xml'
			format := policyfragment.PolicyFragmentContentFormatXml
			if props.Format != nil {
				format = *props.Format
			}
			d.Set("format", string(format))
		}
	}

	return nil
}
