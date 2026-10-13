// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_custom_domain

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func apiManagementCustomDomainDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CustomDomainID(d.Id())
	if err != nil {
		return err
	}

	apiMgmtId := apimanagementservice.NewServiceID(id.SubscriptionId, id.ResourceGroup, id.ServiceName)

	resp, err := client.Get(ctx, apiMgmtId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("%s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	// Wait for the ProvisioningState to become "Succeeded" before attempting to update
	log.Printf("[DEBUG] Waiting for %s to become ready", *id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Updating", "Unknown"},
		Target:                    []string{"Succeeded", "Ready"},
		Refresh:                   apiManagementRefreshFunc(ctx, client, apiMgmtId),
		MinTimeout:                1 * time.Minute,
		Timeout:                   d.Timeout(pluginsdk.TimeoutDelete),
		ContinuousTargetOccurence: 6,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become ready: %+v", *id, err)
	}

	if resp.Model != nil {
		resp.Model.Properties.HostnameConfigurations = nil
	}

	if err := client.CreateOrUpdateThenPoll(ctx, apiMgmtId, *resp.Model); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	// Wait for the ProvisioningState to become "Succeeded" before attempting to update
	log.Printf("[DEBUG] Waiting for %s to become ready", *id)
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become ready: %+v", *id, err)
	}

	return nil
}
