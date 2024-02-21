// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/policyfragment"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementPolicyFragment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementPolicyFragmentCreateUpdate,
		Read:   resourceApiManagementPolicyFragmentRead,
		Update: resourceApiManagementPolicyFragmentCreateUpdate,
		Delete: resourceApiManagementPolicyFragmentDelete,
		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := policyfragment.ParsePolicyFragmentIDInsensitively(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
			client := meta.(*clients.Client).ApiManagement.PolicyFragmentClient

			id, err := policyfragment.ParsePolicyFragmentIDInsensitively(d.Id())
			if err != nil {
				return nil, err
			}

			resp, err := client.Get(ctx, *id, policyfragment.GetOperationOptions{
				Format: pointer.To(policyfragment.PolicyFragmentContentFormatXml),
			})
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return nil, fmt.Errorf("Api Management Policy Fragment %q was not found in Api Management instance %q in Resource Group %q", id.PolicyFragmentName, id.ServiceName, id.ResourceGroupName)
				}

				return nil, fmt.Errorf("retrieving Api Management Policy Fragment %q (Api Management: %q, Resource Group %q): %+v", id.PolicyFragmentName, id.ServiceName, id.ResourceGroupName, err)
			}

			d.Set("format", policyfragment.PolicyFragmentContentFormatXml)

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": schemaz.SchemaApiManagementChildName(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"format": {
				Type:     pluginsdk.TypeString,
				ForceNew: true,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(policyfragment.PolicyFragmentContentFormatRawxml),
					string(policyfragment.PolicyFragmentContentFormatXml),
				}, false),
				Default: policyfragment.PolicyFragmentContentFormatXml,
			},

			"value": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				DiffSuppressFunc: XmlWithDotNetInterpolationsDiffSuppress,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceApiManagementPolicyFragmentCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.PolicyFragmentClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	id := policyfragment.NewPolicyFragmentID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))
	format := policyfragment.PolicyFragmentContentFormat(d.Get("format").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id, policyfragment.GetOperationOptions{
			Format: &format,
		})
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_policy_fragment", id.ID())
		}
	}

	description := d.Get("description").(string)
	value := d.Get("value").(string)

	parameters := policyfragment.PolicyFragmentContract{
		Properties: &policyfragment.PolicyFragmentContractProperties{
			Description: pointer.To(description),
			Format:      pointer.To(format),
			Value:       value,
		},
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters, policyfragment.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementPolicyFragmentRead(d, meta)
}

func resourceApiManagementPolicyFragmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.PolicyFragmentClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := policyfragment.ParsePolicyFragmentIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	format := policyfragment.PolicyFragmentContentFormat(d.Get("format").(string))
	resp, err := client.Get(ctx, *id, policyfragment.GetOperationOptions{
		Format: &format,
	})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s does not exist - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.PolicyFragmentName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("api_management_name", id.ServiceName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", pointer.From(props.Description))
			d.Set("value", props.Value)
			// the api only returns a format field when it's requested in the GET request as param '?format=rawxml' and only does so when it's not 'xml'
			if props.Format == nil {
				d.Set("format", "xml")
			} else {
				d.Set("format", string(pointer.From(props.Format)))
			}
		}
	}

	return nil
}

func resourceApiManagementPolicyFragmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.PolicyFragmentClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := policyfragment.ParsePolicyFragmentIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id, policyfragment.DeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %s", id, err)
		}
	}

	return nil
}
