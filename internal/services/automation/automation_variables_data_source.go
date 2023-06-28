package automation

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/variable"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutomationVariablesDataSource struct{}

type AutomationVariablesDataSourceModel struct {
	AutomationAccountId string                   `tfschema:"automation_account_id"`
	BooleanVariables    []map[string]interface{} `tfschema:"bool"`
	DateTimeVariables   []map[string]interface{} `tfschema:"datetime"`
	EncryptedVariables  []map[string]interface{} `tfschema:"encrypted"`
	IntegerVariables    []map[string]interface{} `tfschema:"int"`
	NullVariables       []map[string]interface{} `tfschema:"null"`
	StringVariables     []map[string]interface{} `tfschema:"string"`
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
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
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
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
				},
			},
		},

		"datetime": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
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
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"encrypted": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
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
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"int": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
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
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"null": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
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
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"string": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
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
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
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

			client := metadata.Client.Automation.VariableClient

			variableList, err := client.ListByAutomationAccountComplete(ctx, *automationAccountId)
			if err != nil {
				return fmt.Errorf("listing variables in %s: %+v", automationAccountId, err)
			}

			var var_bool []map[string]interface{}
			var var_dt []map[string]interface{}
			var var_encrypt []map[string]interface{}
			var var_int []map[string]interface{}
			var var_null []map[string]interface{}
			var var_str []map[string]interface{}

			for _, v := range variableList.Items {
				_, err := variable.ParseVariableID(*v.Id)
				if err != nil {
					return err
				}

				res := map[string]interface{}{
					"name":        pointer.From(v.Name),
					"description": pointer.From(v.Properties.Description),
					"encrypted":   pointer.From(v.Properties.IsEncrypted),
				}

				datePattern := regexp.MustCompile(`"\\/Date\((-?[0-9]+)\)\\/"`)

				if pointer.From(v.Properties.IsEncrypted) {
					var_encrypt = append(var_encrypt, res)
				} else if v.Properties.Value == nil {
					res["value"] = nil
					var_null = append(var_null, res)
				} else if i, err := strconv.ParseInt(pointer.From(v.Properties.Value), 10, 32); err == nil {
					res["value"] = i
					var_int = append(var_int, res)
				} else if b, err := strconv.ParseBool(pointer.From(v.Properties.Value)); err == nil {
					res["value"] = b
					var_bool = append(var_bool, res)
				} else if matches := datePattern.FindStringSubmatch(pointer.From(v.Properties.Value)); len(matches) == 2 && matches[0] == pointer.From(v.Properties.Value) {
					if t, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
						res["value"] = time.UnixMilli(t).In(time.UTC).Format("2006-01-02T15:04:05.999Z")
					}
					var_dt = append(var_dt, res)
				} else if s, err := strconv.Unquote(pointer.From(v.Properties.Value)); err == nil {
					res["value"] = s
					var_str = append(var_str, res)
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
			model.StringVariables = var_str
			return metadata.Encode(&model)
		},
	}
}
