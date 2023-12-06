// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/variable"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutomationVariablesDataSource struct{}

type AutomationVariablesDataSourceModel struct {
	AutomationAccountId string                     `tfschema:"automation_account_id"`
	BooleanVariables    []helper.BooleanVariable   `tfschema:"bool"`
	DateTimeVariables   []helper.DateTimeVariable  `tfschema:"datetime"`
	EncryptedVariables  []helper.EncryptedVariable `tfschema:"encrypted"`
	IntegerVariables    []helper.IntegerVariable   `tfschema:"int"`
	NullVariables       []helper.NullVariable      `tfschema:"null"`
	ObjectVariables     []helper.ObjectVariable    `tfschema:"object"`
	StringVariables     []helper.StringVariable    `tfschema:"string"`
}

var _ sdk.DataSource = AutomationVariablesDataSource{}

func (v AutomationVariablesDataSource) ResourceType() string {
	return "azurerm_automation_variables"
}

func (v AutomationVariablesDataSource) ModelObject() interface{} {
	return &AutomationVariablesDataSourceModel{}
}

func (v AutomationVariablesDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return variable.ValidateAutomationAccountID
}

func (v AutomationVariablesDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"automation_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: variable.ValidateAutomationAccountID,
		},
	}
}

func (v AutomationVariablesDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"bool": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: helper.DataSourceAutomationVariableCommonSchema(pluginsdk.TypeBool),
			},
		},

		"datetime": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: helper.DataSourceAutomationVariableCommonSchema(pluginsdk.TypeString),
			},
		},

		"encrypted": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: helper.DataSourceAutomationVariableCommonSchema(pluginsdk.TypeString),
			},
		},

		"int": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: helper.DataSourceAutomationVariableCommonSchema(pluginsdk.TypeInt),
			},
		},

		"null": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: helper.DataSourceAutomationVariableCommonSchema(pluginsdk.TypeString),
			},
		},

		"object": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: helper.DataSourceAutomationVariableCommonSchema(pluginsdk.TypeString),
			},
		},

		"string": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: helper.DataSourceAutomationVariableCommonSchema(pluginsdk.TypeString),
			},
		},
	}
}

func (v AutomationVariablesDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model AutomationVariablesDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			automationAccountId, err := variable.ParseAutomationAccountID(model.AutomationAccountId)
			if err != nil {
				return err
			}

			client := metadata.Client.Automation.Variable

			variableList, err := client.ListByAutomationAccountComplete(ctx, pointer.From(automationAccountId))
			if err != nil {
				return fmt.Errorf("listing variables in %s: %+v", automationAccountId, err)
			}

			var var_bool []helper.BooleanVariable
			var var_dt []helper.DateTimeVariable
			var var_encrypt []helper.EncryptedVariable
			var var_int []helper.IntegerVariable
			var var_null []helper.NullVariable
			var var_object []helper.ObjectVariable
			var var_str []helper.StringVariable

			for _, v := range variableList.Items {
				_, err := variable.ParseVariableID(pointer.From(v.Id))
				if err != nil {
					return err
				}

				datePattern := regexp.MustCompile(`"\\/Date\((-?[0-9]+)\)\\/"`)
				var objVar map[string]interface{}

				if pointer.From(v.Properties.IsEncrypted) {
					var_encrypt = append(var_encrypt, helper.EncryptedVariable{
						ID:          pointer.From(v.Id),
						Name:        pointer.From(v.Name),
						Description: pointer.From(v.Properties.Description),
						IsEncrypted: pointer.From(v.Properties.IsEncrypted),
					})
				} else if v.Properties.Value == nil {
					var_null = append(var_null, helper.NullVariable{
						ID:          pointer.From(v.Id),
						Name:        pointer.From(v.Name),
						Description: pointer.From(v.Properties.Description),
						IsEncrypted: pointer.From(v.Properties.IsEncrypted),
					})
				} else if i, err := strconv.ParseInt(pointer.From(v.Properties.Value), 10, 32); err == nil {
					var_int = append(var_int, helper.IntegerVariable{
						ID:          pointer.From(v.Id),
						Name:        pointer.From(v.Name),
						Description: pointer.From(v.Properties.Description),
						IsEncrypted: pointer.From(v.Properties.IsEncrypted),
						Value:       i,
					})
				} else if b, err := strconv.ParseBool(pointer.From(v.Properties.Value)); err == nil {
					var_bool = append(var_bool, helper.BooleanVariable{
						ID:          pointer.From(v.Id),
						Name:        pointer.From(v.Name),
						Description: pointer.From(v.Properties.Description),
						IsEncrypted: pointer.From(v.Properties.IsEncrypted),
						Value:       b,
					})
				} else if matches := datePattern.FindStringSubmatch(pointer.From(v.Properties.Value)); len(matches) == 2 && matches[0] == pointer.From(v.Properties.Value) {
					var value string
					if t, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
						value = time.UnixMilli(t).In(time.UTC).Format("2006-01-02T15:04:05.999Z")
					}
					var_dt = append(var_dt, helper.DateTimeVariable{
						ID:          pointer.From(v.Id),
						Name:        pointer.From(v.Name),
						Description: pointer.From(v.Properties.Description),
						IsEncrypted: pointer.From(v.Properties.IsEncrypted),
						Value:       value,
					})
				} else if err := json.Unmarshal([]byte(*v.Properties.Value), &objVar); err == nil {
					var_object = append(var_object, helper.ObjectVariable{
						ID:          pointer.From(v.Id),
						Name:        pointer.From(v.Name),
						Description: pointer.From(v.Properties.Description),
						IsEncrypted: pointer.From(v.Properties.IsEncrypted),
						Value:       pointer.From(v.Properties.Value),
					})
				} else if s, err := strconv.Unquote(pointer.From(v.Properties.Value)); err == nil {
					var_str = append(var_str, helper.StringVariable{
						ID:          pointer.From(v.Id),
						Name:        pointer.From(v.Name),
						Description: pointer.From(v.Properties.Description),
						IsEncrypted: pointer.From(v.Properties.IsEncrypted),
						Value:       s,
					})
				} else {
					return fmt.Errorf("cannot determine type of variable %q", pointer.From(v.Name))
				}
			}

			metadata.SetID(automationAccountId)
			model.AutomationAccountId = automationAccountId.ID()
			model.BooleanVariables = var_bool
			model.DateTimeVariables = var_dt
			model.EncryptedVariables = var_encrypt
			model.IntegerVariables = var_int
			model.NullVariables = var_null
			model.ObjectVariables = var_object
			model.StringVariables = var_str
			return metadata.Encode(&model)
		},
	}
}
