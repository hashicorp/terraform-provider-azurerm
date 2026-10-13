// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_policy

import (
	"fmt"
	"html"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apipolicy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementAPIPolicyRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ApiPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apipolicy.ParseApiID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, apipolicy.GetOperationOptions{Format: pointer.To(apipolicy.PolicyExportFormatXml)})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for %s: %+v", *id, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("api_management_name", id.ServiceName)
	d.Set("api_name", helpers.GetApiName(id.ApiId))

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			policyContent := ""
			if pc := props.Value; pc != "" {
				policyContent = html.UnescapeString(pc)
			}

			// when you submit an `xml_link` to the API, the API downloads this link and stores it as `xml_content`
			// as such there is no way to set `xml_link` and we'll let Terraform handle it
			d.Set("xml_content", policyContent)
		}
	}
	return nil
}
