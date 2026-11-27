// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_variables

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/variable"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	automationClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataSource struct{}

type AutomationVariablesDataSourceModel struct {
	AutomationAccountId string                      `tfschema:"automation_account_id"`
	BooleanVariables    []helpers.BooleanVariable   `tfschema:"bool"`
	DateTimeVariables   []helpers.DateTimeVariable  `tfschema:"datetime"`
	EncryptedVariables  []helpers.EncryptedVariable `tfschema:"encrypted"`
	IntegerVariables    []helpers.IntegerVariable   `tfschema:"int"`
	NullVariables       []helpers.NullVariable      `tfschema:"null"`
	ObjectVariables     []helpers.ObjectVariable    `tfschema:"object"`
	StringVariables     []helpers.StringVariable    `tfschema:"string"`
}

var _ sdk.DataSource = DataSource{}

func (v DataSource) ResourceType() string {
	return "azurerm_automation_variables"
}

func (v DataSource) Description() string {
	return "Use this data source to get all variables in an Automation Account."
}

func (v DataSource) ModelObject() any {
	return &AutomationVariablesDataSourceModel{}
}

func (v DataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return variable.ValidateAutomationAccountID
}

func (v DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"automation_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: variable.ValidateAutomationAccountID,
		},
	}
}

func (v DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"bool": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: helpers.DataSourceAutomationVariableCommonSchema(pluginsdk.TypeBool),
			},
		},

		"datetime": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: helpers.DataSourceAutomationVariableCommonSchema(pluginsdk.TypeString),
			},
		},

		"encrypted": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: helpers.DataSourceAutomationVariableCommonSchema(pluginsdk.TypeString),
			},
		},

		"int": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: helpers.DataSourceAutomationVariableCommonSchema(pluginsdk.TypeInt),
			},
		},

		"null": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: helpers.DataSourceAutomationVariableCommonSchema(pluginsdk.TypeString),
			},
		},

		"object": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: helpers.DataSourceAutomationVariableCommonSchema(pluginsdk.TypeString),
			},
		},

		"string": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: helpers.DataSourceAutomationVariableCommonSchema(pluginsdk.TypeString),
			},
		},
	}
}

func (v DataSource) Read() sdk.ResourceFunc {
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

			client := automationClient.From(metadata.Client).Variable

			variableList, err := client.ListByAutomationAccountComplete(ctx, pointer.From(automationAccountId))
			if err != nil {
				return fmt.Errorf("listing variables in %s: %+v", automationAccountId, err)
			}

			var var_bool []helpers.BooleanVariable
			var var_dt []helpers.DateTimeVariable
			var var_encrypt []helpers.EncryptedVariable
			var var_int []helpers.IntegerVariable
			var var_null []helpers.NullVariable
			var var_object []helpers.ObjectVariable
			var var_str []helpers.StringVariable

			for _, v := range variableList.Items {
				if _, err := variable.ParseVariableID(pointer.From(v.Id)); err != nil {
					return err
				}

				datePattern := regexp.MustCompile(`"\\/Date\((-?[0-9]+)\)\\/"`)
				var objVar map[string]any

				if pointer.From(v.Properties.IsEncrypted) {
					var_encrypt = append(var_encrypt, helpers.EncryptedVariable{
						ID:          pointer.From(v.Id),
						Name:        pointer.From(v.Name),
						Description: pointer.From(v.Properties.Description),
						IsEncrypted: pointer.From(v.Properties.IsEncrypted),
					})
				} else if v.Properties.Value == nil {
					var_null = append(var_null, helpers.NullVariable{
						ID:          pointer.From(v.Id),
						Name:        pointer.From(v.Name),
						Description: pointer.From(v.Properties.Description),
						IsEncrypted: pointer.From(v.Properties.IsEncrypted),
					})
				} else if i, err := strconv.ParseInt(pointer.From(v.Properties.Value), 10, 32); err == nil {
					var_int = append(var_int, helpers.IntegerVariable{
						ID:          pointer.From(v.Id),
						Name:        pointer.From(v.Name),
						Description: pointer.From(v.Properties.Description),
						IsEncrypted: pointer.From(v.Properties.IsEncrypted),
						Value:       i,
					})
				} else if b, err := strconv.ParseBool(pointer.From(v.Properties.Value)); err == nil {
					var_bool = append(var_bool, helpers.BooleanVariable{
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
					var_dt = append(var_dt, helpers.DateTimeVariable{
						ID:          pointer.From(v.Id),
						Name:        pointer.From(v.Name),
						Description: pointer.From(v.Properties.Description),
						IsEncrypted: pointer.From(v.Properties.IsEncrypted),
						Value:       value,
					})
				} else if err := json.Unmarshal([]byte(*v.Properties.Value), &objVar); err == nil {
					var_object = append(var_object, helpers.ObjectVariable{
						ID:          pointer.From(v.Id),
						Name:        pointer.From(v.Name),
						Description: pointer.From(v.Properties.Description),
						IsEncrypted: pointer.From(v.Properties.IsEncrypted),
						Value:       pointer.From(v.Properties.Value),
					})
				} else if s, err := strconv.Unquote(pointer.From(v.Properties.Value)); err == nil {
					var_str = append(var_str, helpers.StringVariable{
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
