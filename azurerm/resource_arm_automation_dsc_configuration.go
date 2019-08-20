package azurerm

import (
	"bytes"
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAutomationDscConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationDscConfigurationCreateUpdate,
		Read:   resourceArmAutomationDscConfigurationRead,
		Update: resourceArmAutomationDscConfigurationCreateUpdate,
		Delete: resourceArmAutomationDscConfigurationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[a-zA-Z0-9_]{1,64}$`),
					`The name length must be from 1 to 64 characters. The name can only contain letters, numbers and underscores.`,
				),
			},

			"automation_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"content_embedded": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"log_verbose": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmAutomationDscConfigurationCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automation.DscConfigurationClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Automation Dsc Configuration creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	accName := d.Get("automation_account_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, accName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Automation DSC Configuration %q (Account %q / Resource Group %q): %s", name, accName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_automation_dsc_configuration", *existing.ID)
		}
	}

	contentEmbedded := d.Get("content_embedded").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	logVerbose := d.Get("log_verbose").(bool)
	description := d.Get("description").(string)

	parameters := automation.DscConfigurationCreateOrUpdateParameters{
		DscConfigurationCreateOrUpdateProperties: &automation.DscConfigurationCreateOrUpdateProperties{
			LogVerbose:  utils.Bool(logVerbose),
			Description: utils.String(description),
			Source: &automation.ContentSource{
				Type:  automation.EmbeddedContent,
				Value: utils.String(contentEmbedded),
			},
		},
		Location: utils.String(location),
	}

	if _, err := client.CreateOrUpdate(ctx, resGroup, accName, name, parameters); err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, accName, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Automation Dsc Configuration %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAutomationDscConfigurationRead(d, meta)
}

func resourceArmAutomationDscConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automation.DscConfigurationClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	accName := id.Path["automationAccounts"]
	name := id.Path["configurations"]

	resp, err := client.Get(ctx, resGroup, accName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on AzureRM Automation Dsc Configuration %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("automation_account_name", accName)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.DscConfigurationProperties; props != nil {
		d.Set("log_verbose", props.LogVerbose)
		d.Set("description", props.Description)
		d.Set("state", resp.State)
	}

	contentresp, err := client.GetContent(ctx, resGroup, accName, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on AzureRM Automation Dsc Configuration content %q: %+v", name, err)
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(contentresp.Body); err != nil {
		return fmt.Errorf("Error reading from AzureRM Automation Dsc Configuration buffer %q: %+v", name, err)
	}
	content := buf.String()

	d.Set("content_embedded", content)

	return nil
}

func resourceArmAutomationDscConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automation.DscConfigurationClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	accName := id.Path["automationAccounts"]
	name := id.Path["configurations"]

	resp, err := client.Delete(ctx, resGroup, accName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error issuing AzureRM delete request for Automation Dsc Configuration %q: %+v", name, err)
	}

	return nil
}
