package azurerm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppServiceCustomHostnameBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceCustomHostnameBindingCreate,
		Read:   resourceArmAppServiceCustomHostnameBindingRead,
		Delete: resourceArmAppServiceCustomHostnameBindingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Minute * 10),
			Delete: schema.DefaultTimeout(time.Minute * 10),
		},

		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"app_service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceArmAppServiceCustomHostnameBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for App Service Hostname Binding creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	appServiceName := d.Get("app_service_name").(string)
	hostname := d.Get("hostname").(string)

	// first check if there's one in this subscription requiring import
	resp, err := client.GetHostNameBinding(ctx, resourceGroup, appServiceName, hostname)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error checking for the existence of Hostname Binding %q (App Service Name %q / Resource Group %q): %+v", hostname, appServiceName, resourceGroup, err)
		}
	}

	if resp.ID != nil {
		return tf.ImportAsExistsError("azurerm_app_service_custom_hostname_binding", *resp.ID)
	}

	properties := web.HostNameBinding{
		HostNameBindingProperties: &web.HostNameBindingProperties{
			SiteName: utils.String(appServiceName),
		},
	}
	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()
	_, err = client.CreateOrUpdateHostNameBinding(waitCtx, resourceGroup, appServiceName, hostname, properties)
	if err != nil {
		return err
	}

	read, err := client.GetHostNameBinding(ctx, resourceGroup, appServiceName, hostname)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Hostname Binding %q (App Service %q / Resource Group %q) ID", hostname, appServiceName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAppServiceCustomHostnameBindingRead(d, meta)
}

func resourceArmAppServiceCustomHostnameBindingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	appServiceName := id.Path["sites"]
	hostname := id.Path["hostNameBindings"]

	ctx := meta.(*ArmClient).StopContext
	resp, err := client.GetHostNameBinding(ctx, resourceGroup, appServiceName, hostname)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Hostname Binding %q (App Service %q / Resource Group %q) was not found - removing from state", hostname, appServiceName, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on App Service Hostname Binding %q (App Service %q / Resource Group %q): %+v", hostname, appServiceName, resourceGroup, err)
	}

	d.Set("hostname", hostname)
	d.Set("app_service_name", appServiceName)
	d.Set("resource_group_name", resourceGroup)

	return nil
}

func resourceArmAppServiceCustomHostnameBindingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appServicesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	appServiceName := id.Path["sites"]
	hostname := id.Path["hostNameBindings"]

	log.Printf("[DEBUG] Deleting App Service Hostname Binding %q (App Service %q / Resource Group %q)", hostname, appServiceName, resGroup)

	ctx := meta.(*ArmClient).StopContext
	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutDelete))
	defer cancel()
	resp, err := client.DeleteHostNameBinding(waitCtx, resGroup, appServiceName, hostname)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return err
		}
	}

	return nil
}
