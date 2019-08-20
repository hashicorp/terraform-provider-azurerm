package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApiManagementCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementCertificateCreateUpdate,
		Read:   resourceArmApiManagementCertificateRead,
		Update: resourceArmApiManagementCertificateCreateUpdate,
		Delete: resourceArmApiManagementCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": azure.SchemaApiManagementChildName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": azure.SchemaApiManagementName(),

			"data": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validate.Base64String(),
			},

			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"expiration": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subject": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"thumbprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmApiManagementCertificateCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.CertificatesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	data := d.Get("data").(string)
	password := d.Get("password").(string)

	if requireResourcesToBeImported {
		existing, err := client.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Certificate %q (API Management Service %q / Resource Group %q): %s", name, serviceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_certificate", *existing.ID)
		}
	}

	parameters := apimanagement.CertificateCreateOrUpdateParameters{
		CertificateCreateOrUpdateProperties: &apimanagement.CertificateCreateOrUpdateProperties{
			Data:     utils.String(data),
			Password: utils.String(password),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, name, parameters, ""); err != nil {
		return fmt.Errorf("Error creating or updating Certificate %q (Resource Group %q / API Management Service %q): %+v", name, resourceGroup, serviceName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Certificate %q (Resource Group %q / API Management Service %q): %+v", name, resourceGroup, serviceName, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read ID for Certificate %q (Resource Group %q / API Management Service %q)", name, resourceGroup, serviceName)
	}
	d.SetId(*resp.ID)

	return resourceArmApiManagementCertificateRead(d, meta)
}

func resourceArmApiManagementCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.CertificatesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	name := id.Path["certificates"]

	resp, err := client.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Certificate %q (Resource Group %q / API Management Service %q) was not found - removing from state!", name, resourceGroup, serviceName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request for Certificate %q (Resource Group %q / API Management Service %q): %+v", name, resourceGroup, serviceName, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)

	if props := resp.CertificateContractProperties; props != nil {

		if expiration := props.ExpirationDate; expiration != nil {
			formatted := expiration.Format(time.RFC3339)
			d.Set("expiration", formatted)
		}

		d.Set("subject", props.Thumbprint)
		d.Set("thumbprint", props.Thumbprint)
	}

	return nil
}

func resourceArmApiManagementCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.CertificatesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	name := id.Path["certificates"]

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, name, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting Certificate %q (Resource Group %q / API Management Service %q): %+v", name, resourceGroup, serviceName, err)
		}
	}

	return nil
}
