// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package digitaltwins

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/digitaltwinsinstance"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/digitaltwins/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDigitalTwinsInstance() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDigitalTwinsInstanceCreate,
		Read:   resourceDigitalTwinsInstanceRead,
		Update: resourceDigitalTwinsInstanceUpdate,
		Delete: resourceDigitalTwinsInstanceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := digitaltwinsinstance.ParseDigitalTwinsInstanceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DigitalTwinsInstanceName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"host_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"tags": commonschema.Tags(),
		},
	}
}

func resourceDigitalTwinsInstanceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DigitalTwins.InstanceClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := digitaltwinsinstance.NewDigitalTwinsInstanceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.DigitalTwinsGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for an existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_digital_twins_instance", id.ID())
	}

	expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	payload := digitaltwinsinstance.DigitalTwinsDescription{
		Location: location.Normalize(d.Get("location").(string)),
		Identity: expandedIdentity,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.DigitalTwinsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDigitalTwinsInstanceRead(d, meta)
}

func resourceDigitalTwinsInstanceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DigitalTwins.InstanceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := digitaltwinsinstance.ParseDigitalTwinsInstanceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.DigitalTwinsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.DigitalTwinsInstanceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("host_name", props.HostName)
		}

		flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceDigitalTwinsInstanceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DigitalTwins.InstanceClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := digitaltwinsinstance.ParseDigitalTwinsInstanceID(d.Id())
	if err != nil {
		return err
	}

	props := digitaltwinsinstance.DigitalTwinsPatchDescription{}

	if d.HasChange("identity") {
		expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		props.Identity = expandedIdentity
	}

	if d.HasChange("tags") {
		props.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.DigitalTwinsUpdateThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceDigitalTwinsInstanceRead(d, meta)
}

func resourceDigitalTwinsInstanceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DigitalTwins.InstanceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := digitaltwinsinstance.ParseDigitalTwinsInstanceID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DigitalTwinsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
