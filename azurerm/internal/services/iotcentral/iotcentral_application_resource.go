package iotcentral

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/iotcentral/mgmt/2018-09-01/iotcentral"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iotcentral/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iotcentral/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmIotCentralApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmIotCentralAppCreate,
		Read:   resourceArmIotCentralAppRead,
		Update: resourceArmIotCentralAppUpdate,
		Delete: resourceArmIotCentralAppDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IotCentralAppName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sub_domain": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.IotCentralAppSubdomain,
			},

			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.IotCentralAppDisplayName,
			},

			"sku": {
				Type:     schema.TypeString,
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
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.IotCentralAppTemplateName,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmIotCentralAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTCentral.AppsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of existing IoT Central Application  %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_iotcentral_application", *existing.ID)
	}

	resp, err := client.CheckNameAvailability(ctx, iotcentral.OperationInputs{
		Name: utils.String(name),
	})
	if err != nil {
		return fmt.Errorf("Error happened on check name availability. %q (Group Name %q). Error:  %+v", name, resourceGroup, err)
	}
	if !*resp.NameAvailable {
		return fmt.Errorf("Resource name not available. name: %s, Reason:  %q, Message  %q", name, *resp.Reason, *resp.Message)
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

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, app)
	if err != nil {
		return fmt.Errorf("Error creating Iot Central Application.  %v", err)
	}

	if err = future.Future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creating IoT Central Application %q (Resource Group %q):  %+v", name, resourceGroup, err)
	}

	response, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving IoT Central Application %q (Resource Group %q):  %+v", name, resourceGroup, err)
	}

	if response.ID == nil || *response.ID == "" {
		return fmt.Errorf("Error creating IoT Central Application %q (Resource Group %q):  %+v", name, resourceGroup, err)
	}

	d.SetId(*response.ID)
	return resourceArmIotCentralAppRead(d, meta)
}

func resourceArmIotCentralAppUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTCentral.AppsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IoTCentralApplicationID(d.Id())
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	if displayName == "" {
		displayName = id.Name
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
	future, err := client.Update(ctx, id.ResourceGroup, id.Name, appPatch)
	if err != nil {
		return fmt.Errorf("Error update Iot Central Application %q (Resource Group %q).  %+v", id.Name, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the completion of update Iot Central Application %q (Resource Group %q):  %+v", id.Name, id.ResourceGroup, err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error retrieving IoT Central Application %q (Resource Group %q):  %+v", id.Name, id.ResourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read IoT Central Application %q (Resource Group %q):  %+v", id.Name, id.ResourceGroup, err)
	}

	d.SetId(*resp.ID)
	return resourceArmIotCentralAppRead(d, meta)
}

func resourceArmIotCentralAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTCentral.AppsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IoTCentralApplicationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving IoT Central Application %q (Resource Group %q):  %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if err := d.Set("sku", resp.Sku.Name); err != nil {
		return fmt.Errorf("Error setting `sku`:  %+v", err)
	}

	if props := resp.AppProperties; props != nil {
		d.Set("sub_domain", props.Subdomain)
		d.Set("display_name", props.DisplayName)
		d.Set("template", props.Template)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmIotCentralAppDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTCentral.AppsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IoTCentralApplicationID(d.Id())
	if err != nil {
		return err
	}
	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !response.WasNotFound(resp.Response()) {
			return fmt.Errorf("Error delete Iot Central Application %q (Resource Group %q).  %+v", id.Name, id.ResourceGroup, err)
		}
	}

	if err := resp.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error delete Iot Central Application %q Resource Group %q).  %+v", id.Name, id.ResourceGroup, err)
	}
	return nil
}
