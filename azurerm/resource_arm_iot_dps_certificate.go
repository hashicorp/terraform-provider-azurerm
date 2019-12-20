package azurerm

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/provisioningservices/mgmt/2018-01-22/iothub"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmIotDPSCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmIotDPSCertificateCreateUpdate,
		Read:   resourceArmIotDPSCertificateRead,
		Update: resourceArmIotDPSCertificateCreateUpdate,
		Delete: resourceArmIotDPSCertificateDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: `The 'azurerm_iot_dps_certificate' resource is deprecated in favour of the renamed version 'azurerm_iothub_dps_certificate'.

Information on migrating to the renamed resource can be found here: https://terraform.io/docs/providers/azurerm/guides/migrating-between-renamed-resources.html

As such the existing 'azurerm_iot_dps_certificate' resource is deprecated and will be removed in the next major version of the AzureRM Provider (2.0).
`,

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
				ValidateFunc: validate.IoTHubName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"iot_dps_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"certificate_content": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
				Sensitive:    true,
			},
		},
	}
}

func resourceArmIotDPSCertificateCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSCertificateClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	iotDPSName := d.Get("iot_dps_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, name, resourceGroup, iotDPSName, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing IoT Device Provisioning Service Certificate %q (Device Provisioning Service %q / Resource Group %q): %+v", name, iotDPSName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_iot_dps_certificate", *existing.ID)
		}
	}

	certificate := iothub.CertificateBodyDescription{
		Certificate: utils.String(d.Get("certificate_content").(string)),
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, iotDPSName, name, certificate, ""); err != nil {
		return fmt.Errorf("Error creating/updating IoT Device Provisioning Service Certificate %q (Device Provisioning Service %q / Resource Group %q): %+v", name, iotDPSName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, name, resourceGroup, iotDPSName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving IoT Device Provisioning Service Certificate %q (Device Provisioning Service %q / Resource Group %q): %+v", name, iotDPSName, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read IoT Device Provisioning Service Certificate %q (Device Provisioning Service %q / Resource Group %q): %+v", name, iotDPSName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmIotDPSCertificateRead(d, meta)
}

func resourceArmIotDPSCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSCertificateClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	iotDPSName := id.Path["provisioningServices"]
	name := id.Path["certificates"]

	resp, err := client.Get(ctx, name, resourceGroup, iotDPSName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving IoT Device Provisioning Service Certificate %q (Device Provisioning Service %q / Resource Group %q): %+v", name, iotDPSName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("iot_dps_name", iotDPSName)
	// We are unable to set `certificate_content` since it is not returned from the API

	return nil
}

func resourceArmIotDPSCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSCertificateClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	iotDPSName := id.Path["provisioningServices"]
	name := id.Path["certificates"]

	resp, err := client.Get(ctx, name, resourceGroup, iotDPSName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("Error retrieving IoT Device Provisioning Service Certificate %q (Device Provisioning Service %q / Resource Group %q): %+v", name, iotDPSName, resourceGroup, err)
	}

	if resp.Etag == nil {
		return fmt.Errorf("Error deleting IoT Device Provisioning Service Certificate %q (Device Provisioning Service %q / Resource Group %q) because Etag is nil", name, iotDPSName, resourceGroup)
	}

	// TODO address this delete call if https://github.com/Azure/azure-rest-api-specs/pull/6311 get's merged
	if _, err := client.Delete(ctx, resourceGroup, *resp.Etag, iotDPSName, name, "", nil, nil, iothub.ServerAuthentication, nil, nil, nil, ""); err != nil {
		return fmt.Errorf("Error deleting IoT Device Provisioning Service Certificate %q (Device Provisioning Service %q / Resource Group %q): %+v", name, iotDPSName, resourceGroup, err)
	}
	return nil
}
