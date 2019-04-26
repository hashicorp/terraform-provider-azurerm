package azurerm

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAutomationBoolVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationBoolVariableCreateUpdate,
		Read:   resourceArmAutomationBoolVariableRead,
		Update: resourceArmAutomationBoolVariableCreateUpdate,
		Delete: resourceArmAutomationBoolVariableDelete,

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

			"encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"value": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceArmAutomationBoolVariableCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationVariableClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("automation_account_name").(string)

	if requireResourcesToBeImported {
		resp, err := client.Get(ctx, resourceGroup, accountName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for present of existing Automation Bool Variable %q (Automation Account Name %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_automation_bool_variable", *resp.ID)
		}
	}

	description := d.Get("description").(string)
	encrypted := d.Get("encrypted").(bool)
	value := strconv.FormatBool(d.Get("value").(bool))

	parameters := automation.VariableCreateOrUpdateParameters{
		Name: utils.String(name),
		VariableCreateOrUpdateProperties: &automation.VariableCreateOrUpdateProperties{
			Description: utils.String(description),
			IsEncrypted: utils.Bool(encrypted),
			Value:       utils.String(value),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, accountName, name, parameters); err != nil {
		return fmt.Errorf("Error creating Automation Bool Variable %q (Automation Account Name %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Automation Bool Variable %q (Automation Account Name %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Automation Bool Variable %q (Automation Account Name %q / Resource Group %q) ID", name, accountName, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmAutomationBoolVariableRead(d, meta)
}

func resourceArmAutomationBoolVariableRead(d *schema.ResourceData, meta interface{}) error {
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
			log.Printf("[INFO] Automation Bool Variable %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Automation Bool Variable %q (Automation Account Name %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("automation_account_name", accountName)
	if properties := resp.VariableProperties; properties != nil {
		d.Set("description", properties.Description)
		d.Set("encrypted", properties.IsEncrypted)
		if !d.Get("encrypted").(bool) {
			value, err := azure.ParseAzureRmAutomationVariableValue("azurerm_automation_bool_variable", properties.Value)
			if err != nil {
				return err
			}
			d.Set("value", value)
		}
	}

	return nil
}

func resourceArmAutomationBoolVariableDelete(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error deleting Automation Bool Variable %q (Automation Account Name %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}

	return nil
}
