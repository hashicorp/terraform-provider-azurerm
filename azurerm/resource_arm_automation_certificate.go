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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAutomationCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationCertificateCreateUpdate,
		Read:   resourceArmAutomationCertificateRead,
		Update: resourceArmAutomationCertificateCreateUpdate,
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"automation_account_name": {
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

			"exportable": {
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
	client := meta.(*clients.Client).Automation.CertificateClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Automation Certificate creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("automation_account_name").(string)

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

	description := d.Get("description").(string)

	parameters := automation.CertificateCreateOrUpdateParameters{
		Name: &name,
		CertificateCreateOrUpdateProperties: &automation.CertificateCreateOrUpdateProperties{
			Description: &description,
		},
	}

	if v, ok := d.GetOk("base64"); ok {
		base64 := v.(string)
		parameters.CertificateCreateOrUpdateProperties.Base64Value = &base64
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, accountName, name, parameters); err != nil {
		return fmt.Errorf("Error creating/updating Certificate %q (Automation Account %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Certificate %q (Automation Account %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("ID was nil for Automation Certificate %q (Automation Account %q / Resource Group %q)", name, accountName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAutomationCertificateRead(d, meta)
}

func resourceArmAutomationCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.CertificateClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	accountName := id.Path["automationAccounts"]
	name := id.Path["certificates"]

	resp, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Certificate %q (Automation Account %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("automation_account_name", accountName)

	if props := resp.CertificateProperties; props != nil {
		d.Set("exportable", props.IsExportable)
		d.Set("thumbprint", props.Thumbprint)
		d.Set("description", props.Description)
	}

	return nil
}

func resourceArmAutomationCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.CertificateClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	accountName := id.Path["automationAccounts"]
	name := id.Path["certificates"]

	resp, err := client.Delete(ctx, resourceGroup, accountName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting Certificate %q (Automation Account %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
		}
	}

	return nil
}
