// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"html"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/apimanagementservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementPolicyCreateUpdate,
		Read:   resourceApiManagementPolicyRead,
		Update: resourceApiManagementPolicyCreateUpdate,
		Delete: resourceApiManagementPolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := policy.ParseServiceID(id)
			return err
		}),

		SchemaVersion: 3,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ApiManagementPolicyV0ToV1{},
			1: migration.ApiManagementPolicyV1ToV2{},
			2: migration.ApiManagementPolicyV2ToV3{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"api_management_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: apimanagementservice.ValidateServiceID,
			},

			"xml_content": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				Computed:         true,
				ConflictsWith:    []string{"xml_link"},
				ExactlyOneOf:     []string{"xml_link", "xml_content"},
				DiffSuppressFunc: XmlWithDotNetInterpolationsDiffSuppress,
			},

			"xml_link": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ConflictsWith: []string{"xml_content"},
				ExactlyOneOf:  []string{"xml_link", "xml_content"},
			},
		},
	}
}

func resourceApiManagementPolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Either `xml_content` or `xml_link` must be set")
	}

	policyServiceId := policy.NewServiceID(apiMgmtId.SubscriptionId, resourceGroup, serviceName)
	_, err = client.CreateOrUpdate(ctx, policyServiceId, parameters, policy.CreateOrUpdateOperationOptions{})
	if err != nil {
		return fmt.Errorf("creating %s: %+v", policyServiceId, err)
	}

	id := policy.NewServiceID(apiMgmtId.SubscriptionId, apiMgmtId.ResourceGroupName, apiMgmtId.ServiceName)
	d.SetId(id.ID())

	return resourceApiManagementPolicyRead(d, meta)
}

func resourceApiManagementPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceApiManagementPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.PolicyClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := policy.ParseServiceID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id, policy.DeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
