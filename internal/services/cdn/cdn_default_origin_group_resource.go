// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnDefaultOriginGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnDefaultOriginGroupCreate,
		Read:   resourceCdnDefaultOriginGroupRead,
		Update: resourceCdnDefaultOriginGroupUpdate,
		Delete: resourceCdnDefaultOriginGroupDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.EndpointID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"endpoint_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.EndpointID,
			},

			"default_origin_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.OriginGroupID,
			},
		},
	}
}

func resourceCdnDefaultOriginGroupCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	endpointsClient := meta.(*clients.Client).Cdn.EndpointsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM CDN Default Origin Group creation...")

	id, err := parse.EndpointID(d.Get("endpoint_id").(string))
	if err != nil {
		return err
	}

	_, err = endpointsClient.Get(ctx, id.ResourceGroup, id.ProfileName, id.Name)
	if err != nil {
		return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
	}

	// No import error for this resource as it does not make sense...

	endpoint := cdn.EndpointUpdateParameters{
		EndpointPropertiesUpdateParameters: &cdn.EndpointPropertiesUpdateParameters{
			DefaultOriginGroup: &cdn.ResourceReference{
				ID: pointer.To(d.Get("default_origin_group_id").(string)),
			},
		},
	}

	future, err := endpointsClient.Update(ctx, id.ResourceGroup, id.ProfileName, id.Name, endpoint)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, endpointsClient.Client); err != nil {
		return fmt.Errorf("waiting for the creation of `cdn_default_origin_group` for %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnDefaultOriginGroupRead(d, meta)
}

func resourceCdnDefaultOriginGroupUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	endpointsClient := meta.(*clients.Client).Cdn.EndpointsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM CDN Default Origin Group update...")

	id, err := parse.EndpointID(d.Get("endpoint_id").(string))
	if err != nil {
		return err
	}

	endpoint := cdn.EndpointUpdateParameters{
		EndpointPropertiesUpdateParameters: &cdn.EndpointPropertiesUpdateParameters{
			DefaultOriginGroup: &cdn.ResourceReference{
				ID: pointer.To(d.Get("default_origin_group_id").(string)),
			},
		},
	}

	future, err := endpointsClient.Update(ctx, id.ResourceGroup, id.ProfileName, id.Name, endpoint)
	if err != nil {
		return fmt.Errorf("updating `cdn_default_origin_group` for %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, endpointsClient.Client); err != nil {
		return fmt.Errorf("waiting for update of `cdn_default_origin_group` for %s: %+v", *id, err)
	}

	return resourceCdnDefaultOriginGroupRead(d, meta)
}

func resourceCdnDefaultOriginGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.EndpointsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("endpoint_id", id.ID())

	if props := resp.EndpointProperties; props != nil {
		d.Set("default_origin_group_id", *props.DefaultOriginGroup.ID)
	}

	return nil
}

func resourceCdnDefaultOriginGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.EndpointsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EndpointID(d.Id())
	if err != nil {
		return err
	}

	endpoint := cdn.EndpointUpdateParameters{
		EndpointPropertiesUpdateParameters: &cdn.EndpointPropertiesUpdateParameters{
			DefaultOriginGroup: &cdn.ResourceReference{
				ID: nil,
			},
		},
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.ProfileName, id.Name, endpoint)
	if err != nil {
		return fmt.Errorf("deleting `cdn_default_origin_group` for %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of `cdn_default_origin_group` for %s: %+v", *id, err)
	}

	// Remove the resource from state...
	d.SetId("")

	return nil
}
