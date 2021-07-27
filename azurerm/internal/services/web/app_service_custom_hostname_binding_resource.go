package web

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-01-15/web"
	"github.com/Azure/go-autorest/autorest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var appServiceCustomHostnameBindingResourceName = "azurerm_app_service_custom_hostname_binding"

func resourceAppServiceCustomHostnameBinding() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceCustomHostnameBindingCreate,
		Read:   resourceAppServiceCustomHostnameBindingRead,
		Delete: resourceAppServiceCustomHostnameBindingDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AppServiceCustomHostnameBindingID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"hostname": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"app_service_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"app_service_slot": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"ssl_state": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(web.SslStateIPBasedEnabled),
					string(web.SslStateSniEnabled),
				}, false),
			},

			"thumbprint": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"virtual_ip": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAppServiceCustomHostnameBindingCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for App Service Hostname Binding creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	appServiceName := d.Get("app_service_name").(string)
	appServiceSlot := d.Get("app_service_slot").(string)
	hostname := d.Get("hostname").(string)
	sslState := d.Get("ssl_state").(string)
	thumbprint := d.Get("thumbprint").(string)

	locks.ByName(appServiceName, appServiceCustomHostnameBindingResourceName)
	defer locks.UnlockByName(appServiceName, appServiceCustomHostnameBindingResourceName)

	slotMessage := GetSlotMessage(appServiceSlot)

	if d.IsNewResource() {
		existing, err := GetHostNameBinding(*client, ctx, resourceGroup, appServiceName, appServiceSlot, hostname)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Custom Hostname Binding %q (App Service %q%s / Resource Group %q): %s", hostname, appServiceName, slotMessage, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_app_service_custom_hostname_binding", *existing.ID)
		}
	}

	properties := web.HostNameBinding{
		HostNameBindingProperties: &web.HostNameBindingProperties{
			SiteName: utils.String(appServiceName),
		},
	}

	if sslState != "" {
		if thumbprint == "" {
			return fmt.Errorf("`thumbprint` must be specified when `ssl_state` is set")
		}

		properties.HostNameBindingProperties.SslState = web.SslState(sslState)
	}

	if thumbprint != "" {
		if sslState == "" {
			return fmt.Errorf("`ssl_state` must be specified when `thumbprint` is set")
		}

		properties.HostNameBindingProperties.Thumbprint = utils.String(thumbprint)
	}

	if _, err := CreateOrUpdateHostNameBinding(*client, ctx, resourceGroup, appServiceName, appServiceSlot, hostname, properties); err != nil {
		return fmt.Errorf("Error creating/updating Custom Hostname Binding %q (App Service %q%s / Resource Group %q): %+v", hostname, appServiceName, slotMessage, resourceGroup, err)
	}

	read, err := GetHostNameBinding(*client, ctx, resourceGroup, appServiceName, appServiceSlot, hostname)
	if err != nil {
		return fmt.Errorf("Error retrieving Custom Hostname Binding %q (App Service %q%s / Resource Group %q): %+v", hostname, appServiceName, slotMessage, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Hostname Binding %q (App Service %q%s / Resource Group %q) ID", hostname, appServiceName, slotMessage, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceAppServiceCustomHostnameBindingRead(d, meta)
}

func resourceAppServiceCustomHostnameBindingRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceCustomHostnameBindingID(d.Id())
	if err != nil {
		return err
	}

	slotMessage := GetSlotMessage(id.AppServiceSlot)

	resp, err := GetHostNameBinding(*client, ctx, id.ResourceGroup, id.AppServiceName, id.AppServiceSlot, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Hostname Binding %q (App Service %q%s / Resource Group %q) was not found - removing from state", id.Name, id.AppServiceName, slotMessage, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Custom Hostname Binding %q (App Service %q%s / Resource Group %q): %+v", id.Name, id.AppServiceName, slotMessage, id.ResourceGroup, err)
	}

	d.Set("hostname", id.Name)
	d.Set("app_service_name", id.AppServiceName)
	d.Set("app_service_slot", id.AppServiceSlot)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.HostNameBindingProperties; props != nil {
		d.Set("ssl_state", props.SslState)
		d.Set("thumbprint", props.Thumbprint)
		d.Set("virtual_ip", props.VirtualIP)
	}

	return nil
}

func resourceAppServiceCustomHostnameBindingDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceCustomHostnameBindingID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.AppServiceName, appServiceCustomHostnameBindingResourceName)
	defer locks.UnlockByName(id.AppServiceName, appServiceCustomHostnameBindingResourceName)

	slotMessage := GetSlotMessage(id.AppServiceSlot)
	log.Printf("[DEBUG] Deleting App Service Hostname Binding %q (App Service %q%s / Resource Group %q)", id.Name, id.AppServiceName, slotMessage, id.ResourceGroup)

	resp, err := DeleteHostNameBinding(*client, ctx, id.ResourceGroup, id.AppServiceName, id.AppServiceSlot, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting Custom Hostname Binding %q (App Service %q%s / Resource Group %q): %+v", id.Name, id.AppServiceName, slotMessage, id.ResourceGroup, err)
		}
	}

	return nil
}

func GetSlotMessage(slot string) (result string) {
	if slot != "" {
		return fmt.Sprintf(" / Slot %q", slot)
	}
	return ""
}

func GetHostNameBinding(client web.AppsClient, ctx context.Context, resourceGroupName string, name string, slot string, hostName string) (result web.HostNameBinding, err error) {
	if slot != "" {
		return client.GetHostNameBindingSlot(ctx, resourceGroupName, name, slot, hostName)
	}
	return client.GetHostNameBinding(ctx, resourceGroupName, name, hostName)
}

func CreateOrUpdateHostNameBinding(client web.AppsClient, ctx context.Context, resourceGroupName string, name string, slot string, hostName string, hostNameBinding web.HostNameBinding) (result web.HostNameBinding, err error) {
	if slot != "" {
		return client.CreateOrUpdateHostNameBindingSlot(ctx, resourceGroupName, name, hostName, hostNameBinding, slot)
	}
	return client.CreateOrUpdateHostNameBinding(ctx, resourceGroupName, name, hostName, hostNameBinding)
}

func DeleteHostNameBinding(client web.AppsClient, ctx context.Context, resourceGroupName string, name string, slot string, hostName string) (result autorest.Response, err error) {
	if slot != "" {
		return client.DeleteHostNameBindingSlot(ctx, resourceGroupName, name, slot, hostName)
	}
	return client.DeleteHostNameBinding(ctx, resourceGroupName, name, hostName)
}
