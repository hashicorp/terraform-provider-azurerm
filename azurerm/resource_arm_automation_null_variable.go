package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAutomationNullVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationNullVariableCreateUpdate,
		Read:   resourceArmAutomationNullVariableRead,
		Update: resourceArmAutomationNullVariableCreateUpdate,
		Delete: resourceArmAutomationNullVariableDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": resourceGroupNameSchema(),

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
		},
	}
}

func resourceArmAutomationNullVariableCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationVariableClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("automation_account_name").(string)

	if requireResourcesToBeImported {
		resp, err := client.Get(ctx, resourceGroup, accountName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for present of existing Automation Null Variable %q (Automation Account Name %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_automation_null_variable", *resp.ID)
		}
	}

	description := d.Get("description").(string)

	parameters := automation.VariableCreateOrUpdateParameters{
		Name: utils.String(name),
		VariableCreateOrUpdateProperties: &automation.VariableCreateOrUpdateProperties{
			Description: utils.String(description),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, accountName, name, parameters); err != nil {
		return fmt.Errorf("Error creating Automation Null Variable %q (Automation Account Name %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Automation Null Variable %q (Automation Account Name %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Automation Null Variable %q (Automation Account Name %q / Resource Group %q) ID", name, accountName, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmAutomationNullVariableRead(d, meta)
}

func resourceArmAutomationNullVariableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationVariableClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	accountName := id.Path["automationAccounts"]
	name := id.Path["variables"]

	resp, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Automation Null Variable %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Automation Null Variable %q (Automation Account Name %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("automation_account_name", accountName)
	if properties := resp.VariableProperties; properties != nil {
		d.Set("description", properties.Description)
	}

	return nil
}

func resourceArmAutomationNullVariableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationVariableClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	accountName := id.Path["automationAccounts"]
	name := id.Path["variables"]

	if _, err := client.Delete(ctx, resourceGroup, accountName, name); err != nil {
		return fmt.Errorf("Error deleting Automation Null Variable %q (Automation Account Name %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}

	return nil
}
