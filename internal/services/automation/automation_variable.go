// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/variable"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func ParseAzureAutomationVariableValue(resource string, input *string) (interface{}, error) {
	if input == nil {
		if resource != "azurerm_automation_variable_null" {
			return nil, fmt.Errorf("expected value \"nil\" to be %q, actual type is \"azurerm_automation_variable_null\"", resource)
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
			actualResource = "azurerm_automation_variable_datetime"
		}
	} else if value, err = strconv.Unquote(*input); err == nil {
		actualResource = "azurerm_automation_variable_string"
	} else if value, err = strconv.ParseInt(*input, 10, 32); err == nil {
		value = int32(value.(int64))
		actualResource = "azurerm_automation_variable_int"
	} else if value, err = strconv.ParseBool(*input); err == nil {
		actualResource = "azurerm_automation_variable_bool"
	} else if err := json.Unmarshal([]byte(*input), &value); err == nil {
		value = *input
		actualResource = "azurerm_automation_variable_object"
	}

	if actualResource != resource {
		return nil, fmt.Errorf("expected value %q to be %q, actual type is %q", *input, resource, actualResource)
	}
	return value, nil
}

func resourceAutomationVariableCommonSchema(attType pluginsdk.ValueType, validateFunc pluginsdk.SchemaValidateFunc) map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"automation_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RunbookName(),
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"encrypted": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"value": {
			Type:         attType,
			Optional:     true,
			ValidateFunc: validateFunc,
		},
	}
}

func datasourceAutomationVariableCommonSchema(attType pluginsdk.ValueType) map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"automation_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.RunbookName(),
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"encrypted": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"value": {
			Type:     attType,
			Computed: true,
		},
	}
}

func resourceAutomationVariableCreateUpdate(d *pluginsdk.ResourceData, meta interface{}, varType string) error {
	client := meta.(*clients.Client).Automation.Variable
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	varTypeLower := strings.ToLower(varType)

	id := variable.NewVariableID(subscriptionId, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Automation %s Variable %s: %+v", varType, id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError(fmt.Sprintf("azurerm_automation_variable_%s", varTypeLower), id.ID())
		}
	}

	description := d.Get("description").(string)
	encrypted := d.Get("encrypted").(bool)
	value := ""

	switch varTypeLower {
	case "datetime":
		vTime, parseErr := time.Parse(time.RFC3339, d.Get("value").(string))
		if parseErr != nil {
			return fmt.Errorf("invalid time format: %+v", parseErr)
		}
		value = fmt.Sprintf("\"\\/Date(%d)\\/\"", vTime.UnixNano()/1000000)
	case "bool":
		value = strconv.FormatBool(d.Get("value").(bool))
	case "int":
		value = strconv.Itoa(d.Get("value").(int))
	case "object":
		// We don't quote the object so it gets saved as a JSON object
		value = d.Get("value").(string)
	case "string":
		value = strconv.Quote(d.Get("value").(string))
	}

	parameters := variable.VariableCreateOrUpdateParameters{
		Name: id.VariableName,
		Properties: variable.VariableCreateOrUpdateProperties{
			Description: utils.String(description),
			IsEncrypted: utils.Bool(encrypted),
		},
	}

	if varTypeLower != "null" {
		parameters.Properties.Value = utils.String(value)
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating Automation %s Variable %s: %+v", varType, id, err)
	}

	d.SetId(id.ID())

	return resourceAutomationVariableRead(d, meta, varType)
}

func resourceAutomationVariableRead(d *pluginsdk.ResourceData, meta interface{}, varType string) error {
	client := meta.(*clients.Client).Automation.Variable
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := variable.ParseVariableID(d.Id())
	if err != nil {
		return err
	}

	varTypeLower := strings.ToLower(varType)

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Automation %s Variable %q does not exist - removing from state", varType, d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Automation %s Variable %q (Automation Account Name %q / Resource Group %q): %+v", varType, id.VariableName, id.AutomationAccountName, id.ResourceGroupName, err)
	}

	d.Set("name", id.VariableName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("automation_account_name", id.AutomationAccountName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", props.Description)
			d.Set("encrypted", props.IsEncrypted)
			if !d.Get("encrypted").(bool) {
				value, err := ParseAzureAutomationVariableValue(fmt.Sprintf("azurerm_automation_variable_%s", varTypeLower), props.Value)
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
	}

	return nil
}

func dataSourceAutomationVariableRead(d *pluginsdk.ResourceData, meta interface{}, varType string) error {
	client := meta.(*clients.Client).Automation.Variable
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := variable.NewVariableID(subscriptionId, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))
	varTypeLower := strings.ToLower(varType)

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Automation %s Variable %q does not exist - removing from state", varType, d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Automation %s Variable %s: %+v", varType, id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.VariableName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("automation_account_name", id.AutomationAccountName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", props.Description)
			d.Set("encrypted", props.IsEncrypted)
			if !d.Get("encrypted").(bool) {
				value, err := ParseAzureAutomationVariableValue(fmt.Sprintf("azurerm_automation_variable_%s", varTypeLower), props.Value)
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
	}

	return nil
}

func resourceAutomationVariableDelete(d *pluginsdk.ResourceData, meta interface{}, varType string) error {
	client := meta.(*clients.Client).Automation.Variable
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := variable.ParseVariableID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting Automation %s Variable %q (Automation Account Name %q / Resource Group %q): %+v", varType, id.VariableName, id.AutomationAccountName, id.ResourceGroupName, err)
	}

	return nil
}
