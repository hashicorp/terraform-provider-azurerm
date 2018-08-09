package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/hashicorp/terraform/helper/schema"
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
			},

			"automation_account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceArmAutomationDscConfigurationCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationDscConfigurationClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Automation Dsc Configuration creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	accName := d.Get("automation_account_name").(string)
	content := d.Get("content").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))

	parameters := automation.DscConfigurationCreateOrUpdateParameters{
		DscConfigurationCreateOrUpdateProperties: &automation.DscConfigurationCreateOrUpdateProperties{
			Source: &automation.ContentSource{
				Type:  automation.EmbeddedContent,
				Value: &content,
			},
		},
		Location: &location,
		Name:     &name,
	}

	_, err := client.CreateOrUpdate(ctx, resGroup, accName, name, parameters)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, accName, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Automation Dsc Configuration '%s' (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAutomationDscConfigurationRead(d, meta)
}

func resourceArmAutomationDscConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationDscConfigurationClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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

		return fmt.Errorf("Error making Read request on AzureRM Automation Dsc Configuration '%s': %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("automation_account_name", accName)

	return nil
}

func resourceArmAutomationDscConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationDscConfigurationClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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

		return fmt.Errorf("Error issuing AzureRM delete request for Automation Dsc Configuration '%s': %+v", name, err)
	}

	return nil
}
