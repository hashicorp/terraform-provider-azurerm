package iotcentral

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/iotcentral/mgmt/2018-09-01/iotcentral"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iotcentral/migration"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iotcentral/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iotcentral/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceIotCentralApplication() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotCentralAppCreate,
		Read:   resourceIotCentralAppRead,
		Update: resourceIotCentralAppUpdate,
		Delete: resourceIotCentralAppDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ApplicationID(id)
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
					string(iotcentral.F1),
					string(iotcentral.S1),
					string(iotcentral.ST1),
					string(iotcentral.ST2),
				}, true),
				Default: iotcentral.ST1,
			},
			"template": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ApplicationTemplateName,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceIotCentralAppCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTCentral.AppsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewApplicationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_iotcentral_application", id.ID())
	}

	resp, err := client.CheckNameAvailability(ctx, iotcentral.OperationInputs{
		Name: utils.String(name),
	})
	if err != nil {
		return fmt.Errorf("checking if the name %q was globally available:  %+v", id.IoTAppName, err)
	}
	if !*resp.NameAvailable {
		return fmt.Errorf("the name %q cannot be used. Reason: %q Message: %q", name, *resp.Reason, *resp.Message)
	}

	displayName := d.Get("display_name").(string)
	if displayName == "" {
		displayName = name
	}

	subdomain := d.Get("sub_domain").(string)
	template := d.Get("template").(string)
	location := d.Get("location").(string)
	app := iotcentral.App{
		AppProperties: &iotcentral.AppProperties{
			DisplayName: &displayName,
			Subdomain:   &subdomain,
			Template:    &template,
		},
		Sku: &iotcentral.AppSkuInfo{
			Name: iotcentral.AppSku(d.Get("sku").(string)),
		},
		Location: &location,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IoTAppName, app)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceIotCentralAppRead(d, meta)
}

func resourceIotCentralAppUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTCentral.AppsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApplicationID(d.Id())
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	if displayName == "" {
		displayName = id.IoTAppName
	}

	subdomain := d.Get("sub_domain").(string)
	template := d.Get("template").(string)
	appPatch := iotcentral.AppPatch{
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
		AppProperties: &iotcentral.AppProperties{
			DisplayName: &displayName,
			Subdomain:   &subdomain,
			Template:    &template,
		},
	}
	future, err := client.Update(ctx, id.ResourceGroup, id.IoTAppName, appPatch)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	return resourceIotCentralAppRead(d, meta)
}

func resourceIotCentralAppRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTCentral.AppsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApplicationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.IoTAppName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.IoTAppName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if err := d.Set("sku", resp.Sku.Name); err != nil {
		return fmt.Errorf("setting `sku`: %+v", err)
	}

	if props := resp.AppProperties; props != nil {
		d.Set("sub_domain", props.Subdomain)
		d.Set("display_name", props.DisplayName)
		d.Set("template", props.Template)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceIotCentralAppDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTCentral.AppsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApplicationID(d.Id())
	if err != nil {
		return err
	}
	resp, err := client.Delete(ctx, id.ResourceGroup, id.IoTAppName)
	if err != nil {
		if !response.WasNotFound(resp.Response()) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	if err := resp.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}
