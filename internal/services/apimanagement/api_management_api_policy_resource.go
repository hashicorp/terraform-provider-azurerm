// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"errors"
	"fmt"
	"html"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apipolicy"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementAPIPolicyCreateUpdate,
		Read:   resourceApiManagementAPIPolicyRead,
		Update: resourceApiManagementAPIPolicyCreateUpdate,
		Delete: resourceApiManagementAPIPolicyDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := apipolicy.ParseApiID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ApiManagementApiPolicyV0ToV1{},
			1: migration.ApiManagementApiPolicyV1ToV2{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": commonschema.ResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"api_name": schemaz.SchemaApiManagementApiName(),

			"xml_content": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				Computed:         true,
				ConflictsWith:    []string{"xml_link"},
				DiffSuppressFunc: XmlWithDotNetInterpolationsDiffSuppress,
			},

			"xml_link": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ConflictsWith: []string{"xml_content"},
			},
		},
	}
}

func resourceApiManagementAPIPolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := apipolicy.NewApiID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("api_name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id, apipolicy.GetOperationOptions{Format: pointer.To(apipolicy.PolicyExportFormatXml)})
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_api_policy", id.ID())
		}
	}

	parameters := apipolicy.PolicyContract{}

	xmlContent := d.Get("xml_content").(string)
	xmlLink := d.Get("xml_link").(string)

	if xmlLink != "" {
		parameters.Properties = &apipolicy.PolicyContractProperties{
			Format: pointer.To(apipolicy.PolicyContentFormatRawxmlNegativelink),
			Value:  xmlLink,
		}
	} else if xmlContent != "" {
		// this is intentionally an else-if since `xml_content` is computed

		// clear out any existing value for xml_link
		if !d.IsNewResource() {
			d.Set("xml_link", "")
		}

		parameters.Properties = &apipolicy.PolicyContractProperties{
			Format: pointer.To(apipolicy.PolicyContentFormatRawxml),
			Value:  xmlContent,
		}
	}

	if parameters.Properties == nil {
		return errors.New("Either `xml_content` or `xml_link` must be set")
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, apipolicy.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementAPIPolicyRead(d, meta)
}

func resourceApiManagementAPIPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
	apiName := getApiName(id.ApiId)
	d.Set("api_name", apiName)

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

func resourceApiManagementAPIPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apipolicy.ParseApiID(d.Id())
	if err != nil {
		return err
	}

	apiName := getApiName(id.ApiId)
	newId := apipolicy.NewApiID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, apiName)
	if resp, err := client.Delete(ctx, newId, apipolicy.DeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
