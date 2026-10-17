// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_policy

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/policy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementPolicyCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.PolicyClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apiManagementID := d.Get("api_management_id").(string)
	apiMgmtId, err := apimanagementservice.ParseServiceID(apiManagementID)
	if err != nil {
		return err
	}
	resourceGroup := apiMgmtId.ResourceGroupName
	serviceName := apiMgmtId.ServiceName

	/*
		Other resources would have a check for d.IsNewResource() at this location, and would error out using `tf.ImportAsExistsError` if the resource already existed.
		However, this is a sub-resource, and the API always returns a policy when queried, either a default policy or one configured by the user or by this pluginsdk.
		Instead of the usual check, the resource documentation clearly states that any existing policy will be overwritten if the resource is used.
	*/

	parameters := policy.PolicyContract{}

	xmlContent := d.Get("xml_content").(string)
	xmlLink := d.Get("xml_link").(string)

	if xmlLink != "" {
		parameters.Properties = &policy.PolicyContractProperties{
			Format: pointer.To(policy.PolicyContentFormatRawxmlNegativelink),
			Value:  xmlLink,
		}
	} else if xmlContent != "" {
		// this is intentionally an else-if since `xml_content` is computed

		// clear out any existing value for xml_link
		if !d.IsNewResource() {
			d.Set("xml_link", "")
		}

		parameters.Properties = &policy.PolicyContractProperties{
			Format: pointer.To(policy.PolicyContentFormatRawxml),
			Value:  xmlContent,
		}
	}

	if parameters.Properties == nil {
		return errors.New("either `xml_content` or `xml_link` must be set")
	}

	policyServiceId := policy.NewServiceID(apiMgmtId.SubscriptionId, resourceGroup, serviceName)
	if _, err = client.CreateOrUpdate(ctx, policyServiceId, parameters, policy.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating %s: %+v", policyServiceId, err)
	}

	id := policy.NewServiceID(apiMgmtId.SubscriptionId, apiMgmtId.ResourceGroupName, apiMgmtId.ServiceName)
	d.SetId(id.ID())

	return resourceApiManagementPolicyRead(d, meta)
}
