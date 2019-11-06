package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAutomationCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationCertificateCreateUpdate,
		Read:   resourceArmAutomationCertificateRead,
		Update: resourceArmAutomationCertificateUpdate,
		Delete: resourceArmAutomationCertificateDelete,

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
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"base64": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Sensitive:    true,
				ValidateFunc: validate.Base64String(),
			},

			"is_exportable": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"thumbprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmAutomationCertificateCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Automation.CertificateClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Automation Certificate creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, accountName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Automation Certificate %q (Account %q / Resource Group %q): %s", name, accountName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_automation_certificate", *existing.ID)
		}
	}

	base64 := d.Get("base64").(string)
	description := d.Get("description").(string)

	parameters := automation.CertificateCreateOrUpdateParameters{
		CertificateCreateOrUpdateProperties: &automation.CertificateCreateOrUpdateProperties{
			Base64Value: &base64,
			Description: &description,
		},
		Name: &name,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, accountName, name, parameters); err != nil {
		return err
	}

	read, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Automation Certificate '%s' (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAutomationCertificateRead(d, meta)
}

func resourceArmAutomationCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Automation.CertificateClient
	ctx, cancel := timeouts.ForUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Automation Certificate update.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	accName := d.Get("account_name").(string)

	description := d.Get("description").(string)

	parameters := automation.CertificateUpdateParameters{
		CertificateUpdateProperties: &automation.CertificateUpdateProperties{
			Description: &description,
		},
		Name: &name,
	}

	read, err := client.Update(ctx, resGroup, accName, name, parameters)
	if err != nil {
		return err
	}

	d.SetId(*read.ID)

	return resourceArmAutomationCertificateRead(d, meta)
}

func resourceArmAutomationCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Automation.CertificateClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	accName := id.Path["automationAccounts"]
	name := id.Path["certificates"]

	resp, err := client.Get(ctx, resGroup, accName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on AzureRM Automation Certificate '%s': %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("account_name", accName)

	if props := resp.CertificateProperties; props != nil {
		d.Set("is_exportable", props.IsExportable)
		d.Set("thumbprint", props.Thumbprint)
		d.Set("description", props.Description)
	}

	return nil
}

func resourceArmAutomationCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Automation.CertificateClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	accName := id.Path["automationAccounts"]
	name := id.Path["certificates"]

	resp, err := client.Delete(ctx, resGroup, accName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error issuing AzureRM delete request for Automation Certificate '%s': %+v", name, err)
	}

	return nil
}
