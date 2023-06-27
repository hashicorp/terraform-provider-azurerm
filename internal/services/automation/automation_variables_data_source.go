package automation

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/variable"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceAutomationVariables() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceAutomationVariablesRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"automation_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: variable.ValidateAutomationAccountID,
			},

			"bool": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
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
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
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
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
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
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
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
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
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
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
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
		},
	}
}

func dataSourceAutomationVariablesRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.VariableClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	automationId := variable.NewAutomationAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string))

	variableList, err := client.ListByAutomationAccountComplete(ctx, automationId)
	if err != nil {
		return fmt.Errorf("listing variables in %s: %+v", automationId, err)
	}

	d.SetId(automationId.ID())
	d.Set("resource_group_name", automationId.ResourceGroupName)
	d.Set("automation_account_name", automationId.AutomationAccountName)

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

		if *v.Properties.IsEncrypted {
			var_encrypt = append(var_encrypt, res)
		} else if v.Properties.Value == nil {
			res["value"] = nil
			var_null = append(var_null, res)
		} else if i, err := strconv.ParseInt(*v.Properties.Value, 10, 32); err == nil {
			res["value"] = i
			var_int = append(var_int, res)
		} else if b, err := strconv.ParseBool(*v.Properties.Value); err == nil {
			res["value"] = b
			var_bool = append(var_bool, res)
		} else if matches := datePattern.FindStringSubmatch(*v.Properties.Value); len(matches) == 2 && matches[0] == *v.Properties.Value {
			if t, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
				res["value"] = time.UnixMilli(t).In(time.UTC).Format("2006-01-02T15:04:05.999Z")
			}
			var_dt = append(var_dt, res)
		} else if s, err := strconv.Unquote(*v.Properties.Value); err == nil {
			res["value"] = s
			var_str = append(var_str, res)
		} else {
			return fmt.Errorf("cannot determine type of variable %q", *v.Name)
		}
	}

	d.Set("bool", var_bool)
	d.Set("datetime", var_dt)
	d.Set("encrypted", var_encrypt)
	d.Set("int", var_int)
	d.Set("null", var_null)
	d.Set("string", var_str)

	return nil
}
