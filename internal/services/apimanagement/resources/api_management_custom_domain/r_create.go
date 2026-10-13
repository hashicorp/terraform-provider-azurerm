// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_custom_domain

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func apiManagementCustomDomainCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ServiceClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apiMgmtId, err := apimanagementservice.ParseServiceID(d.Get("api_management_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewCustomDomainID(apiMgmtId.SubscriptionId, apiMgmtId.ResourceGroupName, apiMgmtId.ServiceName, "default")

	existing, err := client.Get(ctx, *apiMgmtId)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *apiMgmtId, err)
	}

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			if existing.Model != nil && existing.Model.Properties.HostnameConfigurations != nil && len(*existing.Model.Properties.HostnameConfigurations) > 1 {
				return tf.ImportAsExistsError(apiManagementCustomDomainResourceName, *existing.Model.Id)
			}
		}
	}

	existing.Model.Properties.HostnameConfigurations = expandApiManagementCustomDomains(d)

	// Wait for the ProvisioningState to become "Succeeded" before attempting to update
	log.Printf("[DEBUG] Waiting for %s to become ready", *apiMgmtId)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Updating", "Unknown"},
		Target:                    []string{"Succeeded", "Ready"},
		Refresh:                   apiManagementRefreshFunc(ctx, client, *apiMgmtId),
		MinTimeout:                1 * time.Minute,
		ContinuousTargetOccurence: 6,
	}
	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become ready: %+v", id, err)
	}

	// The API expects user assigned identities to be submitted with nil values
	if existing.Model != nil && existing.Model.Identity != nil {
		for k := range existing.Model.Identity.IdentityIds {
			UserAssignedIdentityDetails := existing.Model.Identity.IdentityIds[k]
			UserAssignedIdentityDetails.ClientId = nil
			UserAssignedIdentityDetails.PrincipalId = nil
		}
	}

	// TODO: implement callback, requires migrating to an ID implementing `resourceids.ResourceId`
	if err := client.CreateOrUpdateThenPoll(ctx, *apiMgmtId, *existing.Model); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	// Wait for the ProvisioningState to become "Succeeded" before attempting to update
	log.Printf("[DEBUG] Waiting for %s to become ready", id)
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become ready: %+v", id, err)
	}
	d.SetId(id.ID())

	return apiManagementCustomDomainRead(d, meta)
}

func expandApiManagementCustomDomains(input *pluginsdk.ResourceData) *[]apimanagementservice.HostnameConfiguration {
	results := make([]apimanagementservice.HostnameConfiguration, 0)

	if managementRawVal, ok := input.GetOk("management"); ok {
		vs := managementRawVal.([]any)
		for _, rawVal := range vs {
			v := rawVal.(map[string]any)
			output := helpers.ExpandApiManagementCommonHostnameConfiguration(v, apimanagementservice.HostnameTypeManagement)
			results = append(results, output)
		}
	}

	if portalRawVal, ok := input.GetOk("portal"); ok {
		vs := portalRawVal.([]any)
		for _, rawVal := range vs {
			v := rawVal.(map[string]any)
			output := helpers.ExpandApiManagementCommonHostnameConfiguration(v, apimanagementservice.HostnameTypePortal)
			results = append(results, output)
		}
	}

	if developerPortalRawVal, ok := input.GetOk("developer_portal"); ok {
		vs := developerPortalRawVal.([]any)
		for _, rawVal := range vs {
			v := rawVal.(map[string]any)
			output := helpers.ExpandApiManagementCommonHostnameConfiguration(v, apimanagementservice.HostnameTypeDeveloperPortal)
			results = append(results, output)
		}
	}

	if gatewayRawVal, ok := input.GetOk("gateway"); ok {
		vs := gatewayRawVal.([]any)
		for _, rawVal := range vs {
			v := rawVal.(map[string]any)
			output := helpers.ExpandApiManagementCommonHostnameConfiguration(v, apimanagementservice.HostnameTypeProxy)

			if value, ok := v["default_ssl_binding"]; ok {
				output.DefaultSslBinding = pointer.To(value.(bool))
			}

			results = append(results, output)
		}
	}

	if scmRawVal, ok := input.GetOk("scm"); ok {
		vs := scmRawVal.([]any)
		for _, rawVal := range vs {
			v := rawVal.(map[string]any)
			output := helpers.ExpandApiManagementCommonHostnameConfiguration(v, apimanagementservice.HostnameTypeScm)
			results = append(results, output)
		}
	}
	return &results
}
