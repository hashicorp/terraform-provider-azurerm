// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/applicationgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/desktop"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var applicationGroupType = "azurerm_virtual_desktop_application_group"

func resourceVirtualDesktopApplicationGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualDesktopApplicationGroupCreateUpdate,
		Read:   resourceVirtualDesktopApplicationGroupRead,
		Update: resourceVirtualDesktopApplicationGroupCreateUpdate,
		Delete: resourceVirtualDesktopApplicationGroupDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := applicationgroup.ParseApplicationGroupID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ApplicationGroupV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringIsNotEmpty,
					validation.StringMatch(
						regexp.MustCompile("^[-a-zA-Z0-9]{1,260}$"),
						"Virtual desktop application group name must be 1 - 260 characters long, contain only letters, numbers and hyphens.",
					),
				),
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(applicationgroup.ApplicationGroupTypeDesktop),
					string(applicationgroup.ApplicationGroupTypeRemoteApp),
				}, false),
			},

			"host_pool_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: hostpool.ValidateHostPoolID,
			},

			"friendly_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},

			"default_desktop_display_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 512),
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceVirtualDesktopApplicationGroupCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ApplicationGroupsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	log.Printf("[INFO] preparing arguments for Virtual Desktop Application Group creation")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	locks.ByName(name, applicationGroupType)
	defer locks.UnlockByName(name, applicationGroupType)

	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := applicationgroup.NewApplicationGroupID(subscriptionId, resourceGroup, name)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_virtual_desktop_application_group", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	payload := applicationgroup.ApplicationGroup{
		Location: &location,
		Tags:     tags.Expand(t),
		Properties: applicationgroup.ApplicationGroupProperties{
			ApplicationGroupType: applicationgroup.ApplicationGroupType(d.Get("type").(string)),
			FriendlyName:         utils.String(d.Get("friendly_name").(string)),
			Description:          utils.String(d.Get("description").(string)),
			HostPoolArmPath:      d.Get("host_pool_id").(string),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating Virtual Desktop Application Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if applicationgroup.ApplicationGroupType(d.Get("type").(string)) == applicationgroup.ApplicationGroupTypeDesktop {
		if desktopFriendlyName := utils.String(d.Get("default_desktop_display_name").(string)); desktopFriendlyName != nil {
			desktopClient := meta.(*clients.Client).DesktopVirtualization.DesktopsClient
			// default desktop name created for Application Group is 'sessionDesktop'
			desktopId := desktop.NewDesktopID(id.SubscriptionId, id.ResourceGroupName, id.ApplicationGroupName, "sessionDesktop")
			model, err := desktopClient.Get(ctx, desktopId)
			if err != nil {
				if !response.WasNotFound(model.HttpResponse) {
					return fmt.Errorf("retrieving default desktop for %s: %+v", id, err)
				}
			}

			desktopPatch := desktop.DesktopPatch{
				Properties: &desktop.DesktopPatchProperties{
					FriendlyName: desktopFriendlyName,
				},
			}

			if _, err := desktopClient.Update(ctx, desktopId, desktopPatch); err != nil {
				return fmt.Errorf("setting friendly name for default desktop %s: %+v", id, err)
			}
		}
	}

	d.SetId(id.ID())

	return resourceVirtualDesktopApplicationGroupRead(d, meta)
}

func resourceVirtualDesktopApplicationGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ApplicationGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := applicationgroup.ParseApplicationGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ApplicationGroupName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		props := model.Properties

		d.Set("friendly_name", props.FriendlyName)
		d.Set("description", props.Description)
		d.Set("type", string(props.ApplicationGroupType))
		defaultDesktopDisplayName := ""
		if props.ApplicationGroupType == applicationgroup.ApplicationGroupTypeDesktop {
			desktopClient := meta.(*clients.Client).DesktopVirtualization.DesktopsClient
			// default desktop name created for Application Group is 'sessionDesktop'
			desktopId := desktop.NewDesktopID(id.SubscriptionId, id.ResourceGroupName, id.ApplicationGroupName, "sessionDesktop")
			desktopResp, err := desktopClient.Get(ctx, desktopId)
			if err != nil {
				if !response.WasNotFound(desktopResp.HttpResponse) {
					return fmt.Errorf("retrieving default desktop for %s: %+v", *id, err)
				}
			}
			// if the default desktop was found then set the display name attribute
			if desktopModel := desktopResp.Model; desktopModel != nil && desktopModel.Properties != nil && desktopModel.Properties.FriendlyName != nil {
				defaultDesktopDisplayName = *desktopModel.Properties.FriendlyName
			}
		}
		d.Set("default_desktop_display_name", defaultDesktopDisplayName)

		hostPoolId, err := hostpool.ParseHostPoolIDInsensitively(props.HostPoolArmPath)
		if err != nil {
			return fmt.Errorf("parsing Host Pool ID %q: %+v", props.HostPoolArmPath, err)
		}
		d.Set("host_pool_id", hostPoolId.ID())

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceVirtualDesktopApplicationGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ApplicationGroupsClient

	id, err := applicationgroup.ParseApplicationGroupID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ApplicationGroupName, applicationGroupType)
	defer locks.UnlockByName(id.ApplicationGroupName, applicationGroupType)

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	if _, err = client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
