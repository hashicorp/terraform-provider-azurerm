// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iotcentral

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iotcentral/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iotcentral/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceIotCentralApplication() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceIotCentralAppCreate,
		Read:   resourceIotCentralAppRead,
		Update: resourceIotCentralAppUpdate,
		Delete: resourceIotCentralAppDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := apps.ParseIotAppID(id)
			return err
		}),

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ApplicationV0ToV1{},
			1: migration.ApplicationV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApplicationName,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"sub_domain": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ApplicationSubdomain,
			},

			"display_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ApplicationDisplayName,
			},

			"identity": commonschema.SystemAssignedIdentityOptional(),

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"sku": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(apps.AppSkuSTOne),
					string(apps.AppSkuSTTwo),
					string(apps.AppSkuSTZero),
				}, false),
				Default: string(apps.AppSkuSTOne),
			},
			"template": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "iotc-pnp-preview@1.0.0",
				ValidateFunc: validate.ApplicationTemplateName,
			},

			"tags": commonschema.Tags(),
		},
	}

	if !features.FourPointOhBeta() {
		resource.Schema["template"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: validate.ApplicationTemplateName,
		}
	}

	return resource
}

func resourceIotCentralAppCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTCentral.AppsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := apps.NewIotAppID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_iotcentral_application", id.ID())
	}

	inputs := apps.OperationInputs{
		Name: id.IotAppName,
	}

	resp, err := client.CheckNameAvailability(ctx, commonids.NewSubscriptionID(id.SubscriptionId), inputs)
	if err != nil {
		return fmt.Errorf("checking if the name %q was globally available: %+v", id.IotAppName, err)
	}
	if model := resp.Model; model != nil {
		if !*model.NameAvailable {
			return fmt.Errorf("the name %q cannot be used. Reason: %q Message: %q", id.IotAppName, *model.Reason, *model.Message)
		}
	}

	displayName := d.Get("display_name").(string)
	if displayName == "" {
		displayName = id.IotAppName
		if !features.FourPointOhBeta() {
			displayName = id.ResourceGroupName
		}
	}

	identity, err := identity.ExpandSystemAssigned(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	subdomain := d.Get("sub_domain").(string)
	template := d.Get("template").(string)
	publicNetworkAccess := apps.PublicNetworkAccessEnabled
	app := apps.App{
		Properties: &apps.AppProperties{
			DisplayName:         &displayName,
			PublicNetworkAccess: &publicNetworkAccess,
			Subdomain:           &subdomain,
			Template:            &template,
		},
		Sku: apps.AppSkuInfo{
			Name: apps.AppSku(d.Get("sku").(string)),
		},
		Identity: identity,
		Location: d.Get("location").(string),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, app); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// Public Network Access can only be disabled after creation
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess := apps.PublicNetworkAccessDisabled
		app.Properties.PublicNetworkAccess = &publicNetworkAccess
		if err := client.CreateOrUpdateThenPoll(ctx, id, app); err != nil {
			return fmt.Errorf("updating `public_network_access_enabled` to false for %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())
	return resourceIotCentralAppRead(d, meta)
}

func resourceIotCentralAppUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTCentral.AppsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apps.ParseIotAppID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil || existing.Model == nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model.Properties == nil {
		existing.Model.Properties = &apps.AppProperties{}
	}

	if d.HasChange("sub_domain") {
		existing.Model.Properties.Subdomain = utils.String(d.Get("sub_domain").(string))
	}

	if d.HasChange("display_name") {
		existing.Model.Properties.DisplayName = utils.String(d.Get("display_name").(string))
	}

	if d.HasChange("sku") {
		existing.Model.Sku = apps.AppSkuInfo{
			Name: apps.AppSku(d.Get("sku").(string)),
		}
	}

	if d.HasChange("template") {
		existing.Model.Properties.Template = utils.String(d.Get("template").(string))
	}

	if d.HasChange("tags") {
		existing.Model.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("identity") {
		identity, err := identity.ExpandSystemAssigned(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		existing.Model.Identity = identity
	}

	if d.HasChange("public_network_access_enabled") {
		publicNetworkAccess := apps.PublicNetworkAccessDisabled
		if d.Get("public_network_access_enabled").(bool) {
			publicNetworkAccess = apps.PublicNetworkAccessEnabled
		}
		existing.Model.Properties.PublicNetworkAccess = &publicNetworkAccess
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceIotCentralAppRead(d, meta)
}

func resourceIotCentralAppRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTCentral.AppsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apps.ParseIotAppID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.IotAppName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		d.Set("sku", model.Sku.Name)

		if props := model.Properties; props != nil {
			d.Set("sub_domain", props.Subdomain)
			d.Set("display_name", props.DisplayName)
			d.Set("template", props.Template)

			publicNetworkAccess := true
			if props.PublicNetworkAccess != nil && *props.PublicNetworkAccess == apps.PublicNetworkAccessDisabled {
				publicNetworkAccess = false
			}
			d.Set("public_network_access_enabled", publicNetworkAccess)
		}

		if err := d.Set("identity", identity.FlattenSystemAssigned(model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err = tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceIotCentralAppDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTCentral.AppsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apps.ParseIotAppID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
