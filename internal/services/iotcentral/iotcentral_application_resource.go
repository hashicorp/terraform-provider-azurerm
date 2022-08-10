package iotcentral

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iotcentral/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iotcentral/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotCentralApplication() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotCentralAppCreate,
		Read:   resourceIotCentralAppRead,
		Update: resourceIotCentralAppUpdate,
		Delete: resourceIotCentralAppDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := apps.ParseIotAppID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ApplicationV0ToV1{},
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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

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
				Computed:     true,
				ValidateFunc: validate.ApplicationTemplateName,
			},

			"tags": commonschema.Tags(),
		},
	}
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
		Name: id.ResourceName,
	}

	resp, err := client.CheckNameAvailability(ctx, commonids.NewSubscriptionID(id.SubscriptionId), inputs)
	if err != nil {
		return fmt.Errorf("checking if the name %q was globally available:  %+v", id.ResourceName, err)
	}
	if model := resp.Model; model != nil {
		if !*model.NameAvailable {
			return fmt.Errorf("the name %q cannot be used. Reason: %q Message: %q", id.ResourceName, *model.Reason, *model.Message)
		}
	}

	displayName := d.Get("display_name").(string)
	if displayName == "" {
		displayName = id.ResourceGroupName
	}

	subdomain := d.Get("sub_domain").(string)
	template := d.Get("template").(string)
	app := apps.App{
		Properties: &apps.AppProperties{
			DisplayName: &displayName,
			Subdomain:   &subdomain,
			Template:    &template,
		},
		Sku: apps.AppSkuInfo{
			Name: apps.AppSku(d.Get("sku").(string)),
		},
		Location: d.Get("location").(string),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, app); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
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

	displayName := d.Get("display_name").(string)
	if displayName == "" {
		displayName = id.ResourceName
	}

	subdomain := d.Get("sub_domain").(string)
	template := d.Get("template").(string)
	appPatch := apps.AppPatch{
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: &apps.AppProperties{
			DisplayName: &displayName,
			Subdomain:   &subdomain,
			Template:    &template,
		},
	}

	if err := client.UpdateThenPoll(ctx, *id, appPatch); err != nil {
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

	d.Set("name", id.ResourceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		d.Set("sku", model.Sku.Name)

		if props := model.Properties; props != nil {
			d.Set("sub_domain", props.Subdomain)
			d.Set("display_name", props.DisplayName)
			d.Set("template", props.Template)
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
