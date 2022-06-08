package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApiManagementGatewayHostnameConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementGatewayHostnameConfigurationCreateUpdate,
		Read:   resourceApiManagementGatewayHostnameConfigurationRead,
		Update: resourceApiManagementGatewayHostnameConfigurationCreateUpdate,
		Delete: resourceApiManagementGatewayHostnameConfigurationDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.GatewayHostnameConfigurationID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},
			"api_management_gateway_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.GatewayID,
			},
			"hostname": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},
			"certificate_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},
		},
	}
}

func resourceApiManagementGatewayHostnameConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayHostnameConfigurationClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	gwID, err := parse.GatewayID(d.Get("api_management_gateway_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `api_management_gateway_id`: %v", err)
	}

	id := parse.NewGatewayHostnameConfigurationID(gwID.SubscriptionId, gwID.ResourceGroup, gwID.ServiceName, gwID.Name, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.GatewayName, id.HostnameConfigurationName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("making read request %s: %+v", id, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_gateway_hostname_configuration", id.ID())
		}
	}

	hostname := d.Get("hostname").(string)
	certificateID := d.Get("certificate_id").(string)

	parameters := apimanagement.GatewayHostnameConfigurationContract{
		GatewayHostnameConfigurationContractProperties: &apimanagement.GatewayHostnameConfigurationContractProperties{
			Hostname:      &hostname,
			CertificateID: &certificateID,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, id.GatewayName, id.HostnameConfigurationName, parameters, ""); err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementGatewayHostnameConfigurationRead(d, meta)
}

func resourceApiManagementGatewayHostnameConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayHostnameConfigurationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.GatewayHostnameConfigurationID(d.Id())
	if err != nil {
		return err
	}

	gatewayId := parse.NewGatewayID(id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.GatewayName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.GatewayName, id.HostnameConfigurationName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] GatewayHostnameConfiguration %q (Gateway %q / Resource Group %q / API Management Service %q) was not found - removing from state!", id.HostnameConfigurationName, id.GatewayName, id.ResourceGroup, id.ServiceName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making read request for %s: %+v", id, err)
	}

	d.Set("name", resp.Name)
	d.Set("api_management_gateway_id", gatewayId.ID())

	if properties := resp.GatewayHostnameConfigurationContractProperties; properties != nil {
		// TODO: do we have to check for nil values here
		d.Set("certificate_id", resp.CertificateID)
		d.Set("hostname", resp.Hostname)
	}

	return nil
}

func resourceApiManagementGatewayHostnameConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayHostnameConfigurationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.GatewayHostnameConfigurationID(d.Id())
	if err != nil {
		return err
	}
	if resp, err := client.Delete(ctx, id.ResourceGroup, id.ServiceName, id.GatewayName, id.HostnameConfigurationName, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	return nil
}
