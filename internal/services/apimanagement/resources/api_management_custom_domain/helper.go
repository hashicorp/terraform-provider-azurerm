// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_custom_domain

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func apiManagementRefreshFunc(ctx context.Context, client *apimanagementservice.ApiManagementServiceClient, id apimanagementservice.ServiceId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		log.Printf("[DEBUG] Checking to see if API Management Service %q (Resource Group: %q) is available..", id.ServiceName, id.ResourceGroupName)

		resp, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return nil, "NotFound", nil
			}

			return nil, "", fmt.Errorf("polling for the state of the API Management Service %q (Resource Group: %q): %+v", id.ServiceName, id.ResourceGroupName, err)
		}

		state := ""
		if model := resp.Model; model != nil {
			if provisioningState := model.Properties.ProvisioningState; provisioningState != nil {
				state = pointer.From(provisioningState)
			}
		}

		return resp, state, nil
	}
}
