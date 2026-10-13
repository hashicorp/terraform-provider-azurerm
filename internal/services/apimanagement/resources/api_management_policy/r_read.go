// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_policy

import (
	"fmt"
	"html"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/policy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementPolicyRead(d *pluginsdk.ResourceData, meta any) error {
	serviceClient := meta.(*clients.Client).ApiManagement.ServiceClient
	client := meta.(*clients.Client).ApiManagement.PolicyClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := policy.ParseServiceID(d.Id())
	if err != nil {
		return err
	}

	apimServiceId := apimanagementservice.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName)
	serviceResp, err := serviceClient.Get(ctx, apimServiceId)
	if err != nil {
		if response.WasNotFound(serviceResp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", apimServiceId)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", apimServiceId, err)
	}

	if model := serviceResp.Model; model != nil {
		d.Set("api_management_id", pointer.From(model.Id))
	}

	serviceId := policy.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName)
	resp, err := client.Get(ctx, serviceId, policy.GetOperationOptions{Format: pointer.To(policy.PolicyExportFormatXml)})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			// when you submit an `xml_link` to the API, the API downloads this link and stores it as `xml_content`
			// as such there is no way to set `xml_link` and we'll let Terraform handle it
			d.Set("xml_content", html.UnescapeString(props.Value))
		}
	}

	return nil
}
