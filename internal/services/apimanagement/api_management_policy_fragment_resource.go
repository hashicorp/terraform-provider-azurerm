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
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apimanagementservice"
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
		Create: resourceApiManagementPolicyFragmentCreate,
		Read:   resourceApiManagementPolicyFragmentRead,
		Update: resourceApiManagementPolicyFragmentUpdate,
		Delete: resourceApiManagementPolicyFragmentDelete,
		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := policyfragment.ParsePolicyFragmentID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
			client := meta.(*clients.Client).ApiManagement.PolicyFragmentClient

			id, err := policyfragment.ParsePolicyFragmentID(d.Id())
			if err != nil {
				return nil, err
			}

			resp, err := client.Get(ctx, *id, policyfragment.GetOperationOptions{
				Format: pointer.To(policyfragment.PolicyFragmentContentFormatXml),
			})
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return nil, fmt.Errorf("%s was not found in Api Management instance %q in Resource Group %q", id, id.ServiceName, id.ResourceGroupName)
				}

				return nil, fmt.Errorf("retrieving %s: %+v", id, err)
			}

			d.Set("format", string(policyfragment.PolicyFragmentContentFormatXml))

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

			"api_management_id": commonschema.ResourceIDReferenceRequiredForceNew(&apimanagementservice.ServiceId{}),

			"format": {
				Type:     pluginsdk.TypeString,
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
				DiffSuppressFunc: XmlWhitespaceDiffSuppress,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceApiManagementPolicyFragmentCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.PolicyFragmentClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apiManagementId, err := apimanagementservice.ParseServiceID(d.Get("api_management_id").(string))
	if err != nil {
		return err
	}
	id := policyfragment.NewPolicyFragmentID(apiManagementId.SubscriptionId, apiManagementId.ResourceGroupName, apiManagementId.ServiceName, d.Get("name").(string))
	format := policyfragment.PolicyFragmentContentFormat(d.Get("format").(string))

	opts := policyfragment.DefaultGetOperationOptions()
	opts.Format = &format
	existing, err := client.Get(ctx, id, opts)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_api_management_policy_fragment", id.ID())
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
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementPolicyFragmentRead(d, meta)
}

func resourceApiManagementPolicyFragmentUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.PolicyFragmentClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := policyfragment.ParsePolicyFragmentID(d.Id())
	if err != nil {
		return err
	}

	format := policyfragment.PolicyFragmentContentFormat(d.Get("format").(string))

	opts := policyfragment.DefaultGetOperationOptions()
	opts.Format = pointer.To(format)
	existing, err := client.Get(ctx, *id, opts)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("description") {
		payload.Properties.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("value") {
		payload.Properties.Value = d.Get("value").(string)
	}

	if d.HasChange("format") {
		payload.Properties.Format = pointer.To(format)
		// on format change we send the value also because it might be interpreted differently because of the changed format
		payload.Properties.Value = d.Get("value").(string)
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload, policyfragment.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementPolicyFragmentRead(d, meta)
}

func resourceApiManagementPolicyFragmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.PolicyFragmentClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := policyfragment.ParsePolicyFragmentID(d.Id())
	if err != nil {
		return err
	}

	format := policyfragment.PolicyFragmentContentFormat(d.Get("format").(string))
	opts := policyfragment.DefaultGetOperationOptions()
	opts.Format = &format
	resp, err := client.Get(ctx, *id, opts)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s does not exist - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.PolicyFragmentName)
	apiManagementId := apimanagementservice.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName)
	d.Set("api_management_id", apiManagementId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", pointer.From(props.Description))
			d.Set("value", props.Value)
			// the api only returns a format field when it's requested in the GET request as param '?format=rawxml' and only does so when it's not 'xml'
			format := policyfragment.PolicyFragmentContentFormatXml
			if props.Format != nil {
				format = *props.Format
			}
			d.Set("format", string(format))
		}
	}

	return nil
}

func resourceApiManagementPolicyFragmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.PolicyFragmentClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := policyfragment.ParsePolicyFragmentID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, *id, policyfragment.DefaultDeleteOperationOptions()); err != nil {
		return fmt.Errorf("deleting %s: %s", *id, err)
	}

	return nil
}
