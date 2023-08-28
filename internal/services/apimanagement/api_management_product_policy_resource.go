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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/productpolicy"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementProductPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementProductPolicyCreateUpdate,
		Read:   resourceApiManagementProductPolicyRead,
		Update: resourceApiManagementProductPolicyCreateUpdate,
		Delete: resourceApiManagementProductPolicyDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := productpolicy.ParseProductID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ApiManagementProductPolicyV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": commonschema.ResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"product_id": schemaz.SchemaApiManagementChildName(),

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

func resourceApiManagementProductPolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := productpolicy.NewProductID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("product_id").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id, productpolicy.GetOperationOptions{Format: pointer.To(productpolicy.PolicyExportFormatXml)})
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_product_policy", id.ID())
		}
	}

	parameters := productpolicy.PolicyContract{}

	xmlContent := d.Get("xml_content").(string)
	xmlLink := d.Get("xml_link").(string)

	if xmlContent != "" {
		parameters.Properties = &productpolicy.PolicyContractProperties{
			Format: pointer.To(productpolicy.PolicyContentFormatRawxml),
			Value:  xmlContent,
		}
	}

	if xmlLink != "" {
		parameters.Properties = &productpolicy.PolicyContractProperties{
			Format: pointer.To(productpolicy.PolicyContentFormatRawxmlNegativelink),
			Value:  xmlLink,
		}
	}

	if parameters.Properties == nil {
		return fmt.Errorf("Either `xml_content` or `xml_link` must be set")
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, productpolicy.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	return resourceApiManagementProductPolicyRead(d, meta)
}

func resourceApiManagementProductPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := productpolicy.ParseProductID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, productpolicy.GetOperationOptions{Format: pointer.To(productpolicy.PolicyExportFormatXml)})
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
	d.Set("product_id", id.ProductId)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			// when you submit an `xml_link` to the API, the API downloads this link and stores it as `xml_content`
			// as such there is no way to set `xml_link` and we'll let Terraform handle it
			d.Set("xml_content", html.UnescapeString(props.Value))
		}
	}

	return nil
}

func resourceApiManagementProductPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := productpolicy.ParseProductID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id, productpolicy.DeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
