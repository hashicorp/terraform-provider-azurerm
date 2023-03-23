package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApiManagementGatewayHostNameConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementGatewayHostNameConfigurationCreateUpdate,
		Read:   resourceApiManagementGatewayHostNameConfigurationRead,
		Update: resourceApiManagementGatewayHostNameConfigurationCreateUpdate,
		Delete: resourceApiManagementGatewayHostNameConfigurationDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.GatewayHostNameConfigurationID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": schemaz.SchemaApiManagementChildName(),

			"api_management_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApiManagementID,
			},

			"gateway_name": schemaz.SchemaApiManagementChildName(),

			"certificate_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.CertificateID,
			},

			"host_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"request_client_certificate_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"http2_enabled": {
				Type:     pluginsdk.TypeBool,
				Default:  true,
				Optional: true,
			},

			"tls10_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"tls11_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceApiManagementGatewayHostNameConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayHostNameConfigurationClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apimId, err := parse.ApiManagementID(d.Get("api_management_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `api_management_id`: %v", err)
	}

	id := parse.NewGatewayHostNameConfigurationID(apimId.SubscriptionId, apimId.ResourceGroup, apimId.ServiceName, d.Get("gateway_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.GatewayName, d.Get("name").(string))
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_api_management_gateway_host_name_configuration", id.ID())
		}
	}

	parameters := apimanagement.GatewayHostnameConfigurationContract{
		GatewayHostnameConfigurationContractProperties: &apimanagement.GatewayHostnameConfigurationContractProperties{
			Hostname:                   utils.String(d.Get("host_name").(string)),
			CertificateID:              utils.String(d.Get("certificate_id").(string)),
			NegotiateClientCertificate: utils.Bool(d.Get("request_client_certificate_enabled").(bool)),
			TLS10Enabled:               utils.Bool(d.Get("tls10_enabled").(bool)),
			TLS11Enabled:               utils.Bool(d.Get("tls11_enabled").(bool)),
			HTTP2Enabled:               utils.Bool(d.Get("http2_enabled").(bool)),
		},
	}

	_, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, d.Get("gateway_name").(string), d.Get("name").(string), parameters, "")
	if err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementGatewayHostNameConfigurationRead(d, meta)
}

func resourceApiManagementGatewayHostNameConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayHostNameConfigurationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.GatewayHostNameConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.GatewayName, id.HostnameConfigurationName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("%s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making read request for %s: %+v", *id, err)
	}

	apimId := parse.NewApiManagementID(id.SubscriptionId, id.ResourceGroup, id.ServiceName)

	d.Set("name", resp.Name)
	d.Set("api_management_id", apimId.ID())
	d.Set("gateway_name", id.GatewayName)

	if properties := resp.GatewayHostnameConfigurationContractProperties; properties != nil {
		d.Set("host_name", properties.Hostname)
		d.Set("certificate_id", properties.CertificateID)
		d.Set("request_client_certificate_enabled", properties.NegotiateClientCertificate)
		d.Set("tls10_enabled", properties.TLS10Enabled)
		d.Set("tls11_enabled", properties.TLS11Enabled)
		d.Set("http2_enabled", properties.HTTP2Enabled)
	}

	return nil
}

func resourceApiManagementGatewayHostNameConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayHostNameConfigurationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.GatewayHostNameConfigurationID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting %s", *id)
	if resp, err := client.Delete(ctx, id.ResourceGroup, id.ServiceName, id.GatewayName, id.HostnameConfigurationName, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
