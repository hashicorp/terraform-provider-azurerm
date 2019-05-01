package azurerm

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func parseAzureRmAutomationVariableValue(resource string, input *string) (interface{}, error) {
	if input == nil {
		if resource != "azurerm_automation_null_variable" {
			return nil, fmt.Errorf("Expected value \"nil\" to be %q, actual type is \"azurerm_automation_null_variable\"", resource)
		}
		return nil, nil
	}

	var value interface{}
	var err error
	actualResource := "Unknown"
	datePattern := regexp.MustCompile(`"\\/Date\((-?[0-9]+)\)\\/"`)
	matches := datePattern.FindStringSubmatch(*input)

	if len(matches) == 2 && matches[0] == *input {
		if ticks, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
			value = time.Unix(ticks/1000, ticks%1000*1000000).In(time.UTC)
			actualResource = "azurerm_automation_datetime_variable"
		}
	} else if value, err = strconv.Unquote(*input); err == nil {
		actualResource = "azurerm_automation_string_variable"
	} else if value, err = strconv.ParseBool(*input); err == nil {
		actualResource = "azurerm_automation_bool_variable"
	} else if value, err = strconv.ParseInt(*input, 10, 32); err == nil {
		value = int32(value.(int64))
		actualResource = "azurerm_automation_int_variable"
	}

	if actualResource != resource {
		return nil, fmt.Errorf("Expected value %q to be %q, actual type is %q", *input, resource, actualResource)
	}
	return value, nil
}

func AutomationVariableCommonSchemaFrom(s map[string]*schema.Schema) map[string]*schema.Schema {

	varSchema := map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.NoEmptyStrings,
		},

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
	}
	return azure.MergeSchema(s, varSchema)
}

func resourceArmAutomationVariableCreateUpdate(d *schema.ResourceData, meta interface{}, varType string) error {
	client := meta.(*ArmClient).automationVariableClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("automation_account_name").(string)
	varTypeLower := strings.ToLower(varType)

	if requireResourcesToBeImported {
		resp, err := client.Get(ctx, resourceGroup, accountName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for present of existing Automation %s Variable %q (Automation Account Name %q / Resource Group %q): %+v", varType, name, accountName, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError(fmt.Sprintf("azurerm_automation_%s_variable", varTypeLower), *resp.ID)
		}
	}

	description := d.Get("description").(string)
	encrypted := d.Get("encrypted").(bool)
	value := ""

	if varTypeLower == "datetime" {
		vTime, parseErr := time.Parse(time.RFC3339, d.Get("value").(string))
		if parseErr != nil {
			return fmt.Errorf("Error invalid time format: %+v", parseErr)
		}
		value = fmt.Sprintf("\"\\/Date(%d)\\/\"", vTime.UnixNano()/1000000)
	} else if varTypeLower == "bool" {
		value = strconv.FormatBool(d.Get("value").(bool))
	} else if varTypeLower == "int" {
		value = strconv.Itoa(d.Get("value").(int))
	} else if varTypeLower == "string" {
		value = strconv.Quote(d.Get("value").(string))
	}

	parameters := automation.VariableCreateOrUpdateParameters{}

	if varTypeLower == "null" {
		parameters = automation.VariableCreateOrUpdateParameters{
			Name: utils.String(name),
			VariableCreateOrUpdateProperties: &automation.VariableCreateOrUpdateProperties{
				Description: utils.String(description),
				IsEncrypted: utils.Bool(encrypted),
			},
		}
	} else {
		parameters = automation.VariableCreateOrUpdateParameters{
			Name: utils.String(name),
			VariableCreateOrUpdateProperties: &automation.VariableCreateOrUpdateProperties{
				Description: utils.String(description),
				IsEncrypted: utils.Bool(encrypted),
				Value:       utils.String(value),
			},
		}
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, accountName, name, parameters); err != nil {
		return fmt.Errorf("Error creating Automation %s Variable %q (Automation Account Name %q / Resource Group %q): %+v", varType, name, accountName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Automation %s Variable %q (Automation Account Name %q / Resource Group %q): %+v", varType, name, accountName, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Automation %s Variable %q (Automation Account Name %q / Resource Group %q) ID", varType, name, accountName, resourceGroup)
	}
	d.SetId(*resp.ID)

	return nil
}

func resourceArmAutomationVariableRead(d *schema.ResourceData, meta interface{}, varType string) error {
	client := meta.(*ArmClient).automationVariableClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	accountName := id.Path["automationAccounts"]
	name := id.Path["variables"]
	varTypeLower := strings.ToLower(varType)

	resp, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Automation %s Variable %q does not exist - removing from state", varType, d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Automation %s Variable %q (Automation Account Name %q / Resource Group %q): %+v", varType, name, accountName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("automation_account_name", accountName)
	if properties := resp.VariableProperties; properties != nil {
		d.Set("description", properties.Description)
		d.Set("encrypted", properties.IsEncrypted)
		if !d.Get("encrypted").(bool) {
			value, err := parseAzureRmAutomationVariableValue(fmt.Sprintf("azurerm_automation_%s_variable", varTypeLower), properties.Value)
			if err != nil {
				return err
			}

			if varTypeLower == "datetime" {
				d.Set("value", value.(time.Time).Format("2006-01-02T15:04:05.999Z"))
			} else if varTypeLower != "null" {
				d.Set("value", value)
			}
		}
	}

	return nil
}

func resourceArmAutomationVariableDelete(d *schema.ResourceData, meta interface{}, varType string) error {
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
		return fmt.Errorf("Error deleting Automation %s Variable %q (Automation Account Name %q / Resource Group %q): %+v", varType, name, accountName, resourceGroup, err)
	}

	return nil
}
