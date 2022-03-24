package iothub

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/provisioningservices/mgmt/2021-10-15/iothub"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceIotHubDPSCertificate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubDPSCertificateCreateUpdate,
		Read:   resourceIotHubDPSCertificateRead,
		Update: resourceIotHubDPSCertificateCreateUpdate,
		Delete: resourceIotHubDPSCertificateDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DpsCertificateID(id)
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"iot_dps_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"certificate_content": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Sensitive:    true,
			},
		},
	}
}

func resourceIotHubDPSCertificateCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSCertificateClient
	subscriptionId := meta.(*clients.Client).IoTHub.DPSResourceClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewDpsCertificateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iot_dps_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.CertificateName, id.ResourceGroup, id.ProvisioningServiceName, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing IoT Device Provisioning Service Certificate %s: %+v", id.String(), err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_iothub_dps_certificate", id.ID())
		}
	}

	certificate := iothub.CertificateBodyDescription{
		Certificate: utils.String(d.Get("certificate_content").(string)),
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ProvisioningServiceName, id.CertificateName, certificate, ""); err != nil {
		return fmt.Errorf("creating/updating IoT Device Provisioning Service Certificate %s: %+v", id.String(), err)
	}

	d.SetId(id.ID())

	return resourceIotHubDPSCertificateRead(d, meta)
}

func resourceIotHubDPSCertificateRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSCertificateClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DpsCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.CertificateName, id.ResourceGroup, id.ProvisioningServiceName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id.String(), err)
	}

	d.Set("name", id.CertificateName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("iot_dps_name", id.ProvisioningServiceName)
	// We are unable to set `certificate_content` since it is not returned from the API

	return nil
}

func resourceIotHubDPSCertificateDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSCertificateClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DpsCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.CertificateName, id.ResourceGroup, id.ProvisioningServiceName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Etag == nil {
		return fmt.Errorf("deleting  %s because Etag is nil", id)
	}

	// TODO address this delete call if https://github.com/Azure/azure-rest-api-specs/pull/6311 get's merged
	if _, err := client.Delete(ctx, id.ResourceGroup, *resp.Etag, id.ProvisioningServiceName, id.CertificateName, "", nil, nil, iothub.CertificatePurposeServerAuthentication, nil, nil, nil, ""); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}
