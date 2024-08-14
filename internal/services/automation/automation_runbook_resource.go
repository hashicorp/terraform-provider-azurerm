// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/jobschedule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/runbook"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/runbookdraft"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func contentLinkSchema(isDraft bool) *pluginsdk.Schema {
	ins := &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"uri": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.Any(
						validation.IsURLWithScheme([]string{"http", "https"}),
						validation.StringIsEmpty,
					),
				},

				"version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"hash": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"algorithm": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"value": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},
		},
	}
	if !isDraft {
		ins.AtLeastOneOf = []string{"content", "publish_content_link", "draft"}
	}
	return ins
}

func resourceAutomationRunbook() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationRunbookCreateUpdate,
		Read:   resourceAutomationRunbookRead,
		Update: resourceAutomationRunbookCreateUpdate,
		Delete: resourceAutomationRunbookDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := runbook.ParseRunbookID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.RunbookName(),
			},

			"automation_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AutomationAccount(),
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"runbook_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(runbook.RunbookTypeEnumGraph),
					string(runbook.RunbookTypeEnumGraphPowerShell),
					string(runbook.RunbookTypeEnumGraphPowerShellWorkflow),
					string(runbook.RunbookTypeEnumPowerShell),
					string(runbook.RunbookTypeEnumPowerShellSevenTwo),
					string(runbook.RunbookTypeEnumPythonTwo),
					string(runbook.RunbookTypeEnumPythonThree),
					string(runbook.RunbookTypeEnumPowerShellWorkflow),
					string(runbook.RunbookTypeEnumPowerShellSevenTwo),
					string(runbook.RunbookTypeEnumScript),
				}, false),
			},

			"log_progress": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},

			"log_verbose": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"content": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				// NOTE: O+C the API returns some defaults for this if `publish_content_link` is specified
				Computed:     true,
				AtLeastOneOf: []string{"content", "publish_content_link", "draft"},
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"job_schedule": {
				Type:       pluginsdk.TypeSet,
				Optional:   true,
				Computed:   true,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"schedule_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.ScheduleName(),
						},

						"parameters": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
							ValidateFunc: validate.ParameterNames,
						},

						"run_on": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"job_schedule_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
				Set: helper.ResourceAutomationJobScheduleHash,
			},

			"publish_content_link": contentLinkSchema(false),

			"draft": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*schema.Schema{
						"creation_time": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"content_link": contentLinkSchema(true),

						"edit_mode_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"last_modified_time": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"output_types": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"parameters": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"type": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"mandatory": {
										Type:     pluginsdk.TypeBool,
										Default:  false,
										Optional: true,
									},

									"position": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
									},

									"default_value": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},

			"log_activity_trace_level": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceAutomationRunbookCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	autoCli := meta.(*clients.Client).Automation
	client := autoCli.Runbook
	jsClient := autoCli.JobSchedule
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Automation Runbook creation.")
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId

	id := runbook.NewRunbookID(subscriptionID, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_automation_runbook", id.ID())
		}
	}

	// for existing runbook, if only job_schedule field updated, then skip update runbook
	if d.IsNewResource() || d.HasChangeExcept("job_schedule") {

		location := azure.NormalizeLocation(d.Get("location").(string))
		t := d.Get("tags").(map[string]interface{})

		runbookType := runbook.RunbookTypeEnum(d.Get("runbook_type").(string))
		logProgress := d.Get("log_progress").(bool)
		logVerbose := d.Get("log_verbose").(bool)
		description := d.Get("description").(string)

		parameters := runbook.RunbookCreateOrUpdateParameters{
			Properties: runbook.RunbookCreateOrUpdateProperties{
				LogVerbose:       &logVerbose,
				LogProgress:      &logProgress,
				RunbookType:      runbookType,
				Description:      &description,
				LogActivityTrace: utils.Int64(int64(d.Get("log_activity_trace_level").(int))),
			},

			Location: &location,
		}
		if tagsVal := expandStringInterfaceMap(t); tagsVal != nil {
			parameters.Tags = &tagsVal
		}

		contentLink := expandContentLink(d.Get("publish_content_link").([]interface{}))
		if contentLink != nil {
			parameters.Properties.PublishContentLink = contentLink
		} else {
			parameters.Properties.Draft = &runbook.RunbookDraft{}
			if draft := expandDraft(d.Get("draft").([]interface{})); draft != nil {
				parameters.Properties.Draft = draft
			}
		}

		if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
			return fmt.Errorf("creating/updating %s: %+v", id, err)
		}

		if v, ok := d.GetOk("content"); ok {
			content := v.(string)
			draftRunbookID := runbookdraft.NewRunbookID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.RunbookName)
			if err := autoCli.RunbookDraft.ReplaceContentThenPoll(ctx, draftRunbookID, []byte(content)); err != nil {
				return fmt.Errorf("setting the draft for %s: %+v", id, err)
			}

			if err := autoCli.Runbook.PublishThenPoll(ctx, id); err != nil {
				return fmt.Errorf("publishing the updated %s: %+v", id, err)
			}
		}

		d.SetId(id.ID())
	}

	// **don't need** to list job schedules and delete all of them. update the runbook will recreate these job schedules automatically,
	// but with a different job schedule id
	// crosscheck these existing jobs and jobs from tf, delete the ones not in tf, and create the ones not in api
	// Fix issue: https://github.com/hashicorp/terraform-provider-azurerm/issues/8634
	jsValue, ok := d.GetOk("job_schedule")
	if ok && d.HasChange("job_schedule") {
		jsMap, err := helper.ExpandAutomationJobSchedule(jsValue.(*pluginsdk.Set).List(), id.RunbookName)
		if err != nil {
			return err
		}

		if err := updatedLinkedJobSchedules(ctx, subscriptionID, jsClient, &id, *jsMap); err != nil {
			return fmt.Errorf("update job schedule links: %v", err)
		}
	}

	return resourceAutomationRunbookRead(d, meta)
}

func resourceAutomationRunbookRead(d *pluginsdk.ResourceData, meta interface{}) error {
	autoCli := meta.(*clients.Client).Automation
	client := autoCli.Runbook
	jsClient := autoCli.JobSchedule
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := runbook.ParseRunbookID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on AzureRM Automation Runbook %q (Account %q / Resource Group %q): %+v", id.RunbookName, id.AutomationAccountName, id.ResourceGroupName, err)
	}

	d.Set("name", id.RunbookName)
	d.Set("resource_group_name", id.ResourceGroupName)
	model := resp.Model
	if location := model.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("automation_account_name", id.AutomationAccountName)
	if props := model.Properties; props != nil {
		d.Set("log_verbose", props.LogVerbose)
		d.Set("log_progress", props.LogProgress)
		d.Set("runbook_type", string(pointer.From(props.RunbookType)))
		d.Set("description", props.Description)
		d.Set("log_activity_trace_level", props.LogActivityTrace)
	}

	// GetContent need to use preview version client RunbookClientHack
	// move to stable Runbook once this issue fixed: https://github.com/Azure/azure-sdk-for-go/issues/17591#issuecomment-1233676539
	contentResp, err := autoCli.Runbook.GetContent(ctx, *id)
	if err != nil {
		if response.WasNotFound(contentResp.HttpResponse) {
			d.Set("content", "")
		} else {
			return fmt.Errorf("retrieving content for Automation Runbook %s: %+v", id, err)
		}
	}

	if v := contentResp.Model; v != nil && *v != nil {
		d.Set("content", string(*v))
	}

	jsMap := make(map[uuid.UUID]jobschedule.JobScheduleProperties)
	automationAccountId := jobschedule.NewAutomationAccountID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName)

	filter := fmt.Sprintf("properties/runbook/name eq '%s'", id.RunbookName)
	jsIterator, err := jsClient.ListByAutomationAccount(ctx, automationAccountId, jobschedule.ListByAutomationAccountOperationOptions{Filter: &filter})
	if err != nil {
		return fmt.Errorf("loading Automation Account %q Job Schedule List: %+v", id.AutomationAccountName, err)
	}
	for _, item := range pointer.From(jsIterator.Model) {
		if itemProps := item.Properties; itemProps != nil {
			if itemProps.JobScheduleId == nil || *itemProps.JobScheduleId == "" {
				return fmt.Errorf("job schedule Id is nil or empty listed by Automation Account %q Job Schedule List: %+v", id.AutomationAccountName, err)
			}
			jsId, err := uuid.FromString(*itemProps.JobScheduleId)
			if err != nil {
				return fmt.Errorf("parsing job schedule Id listed by Automation Account %q Job Schedule List: %v", id.AutomationAccountName, err)
			}
			// get job schedule from GET API, `ListByAutomationAccountComplete` lost parameters
			jobscheduleID, err := jobschedule.ParseJobScheduleID(pointer.From(item.Id))
			if err != nil {
				return fmt.Errorf("parsing job schedule Id listed by Automation Account %q Job Schedule List: %v", id.AutomationAccountName, err)
			}
			jsResult, err := jsClient.Get(ctx, *jobscheduleID)
			if err != nil {
				return fmt.Errorf("retrieving job schedule by %s: %v", *jobscheduleID, err)
			}
			if jsResult.Model != nil && jsResult.Model.Properties != nil {
				jsMap[jsId] = *jsResult.Model.Properties
			}
		}
	}

	jobSchedule := helper.FlattenAutomationJobSchedule(jsMap)
	if err := d.Set("job_schedule", jobSchedule); err != nil {
		return fmt.Errorf("setting `job_schedule`: %+v", err)
	}

	if err := tags.FlattenAndSet(d, model.Tags); err != nil {
		return err
	}

	return nil
}

func resourceAutomationRunbookDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.Runbook
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := runbook.ParseRunbookID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("issuing AzureRM delete request for Automation Runbook '%s': %+v", id.RunbookName, err)
	}

	return nil
}

func expandContentLink(inputs []interface{}) *runbook.ContentLink {
	if len(inputs) == 0 || inputs[0] == nil {
		return nil
	}

	input := inputs[0].(map[string]interface{})
	uri := input["uri"].(string)
	version := input["version"].(string)
	hashes := input["hash"].([]interface{})

	if len(hashes) > 0 {
		hash := hashes[0].(map[string]interface{})
		hashValue := hash["value"].(string)
		hashAlgorithm := hash["algorithm"].(string)

		return &runbook.ContentLink{
			Uri:     &uri,
			Version: &version,
			ContentHash: &runbook.ContentHash{
				Algorithm: hashAlgorithm,
				Value:     hashValue,
			},
		}
	}

	return &runbook.ContentLink{
		Uri:     &uri,
		Version: &version,
	}
}

func expandDraft(inputs []interface{}) *runbook.RunbookDraft {
	if len(inputs) == 0 || inputs[0] == nil {
		return nil
	}

	input := inputs[0].(map[string]interface{})
	var res runbook.RunbookDraft

	res.DraftContentLink = expandContentLink(input["content_link"].([]interface{}))
	res.InEdit = utils.Bool(input["edit_mode_enabled"].(bool))
	parameter := map[string]runbook.RunbookParameter{}

	for _, iparam := range input["parameters"].([]interface{}) {
		param := iparam.(map[string]interface{})
		key := param["key"].(string)
		parameter[key] = runbook.RunbookParameter{
			Type:         utils.String(param["type"].(string)),
			IsMandatory:  utils.Bool(param["mandatory"].(bool)),
			Position:     utils.Int64(int64(param["position"].(int))),
			DefaultValue: utils.String(param["default_value"].(string)),
		}
	}
	res.Parameters = &parameter

	var types []string
	for _, v := range input["output_types"].([]interface{}) {
		types = append(types, v.(string))
	}

	if len(types) > 0 {
		res.OutputTypes = &types
	}

	return &res
}

// if job in jsIterator but not in jsMap, then delete it
// if job in both jsIterator and jsMap, remove the entry in jsMap
// at last, create jobs still in jsMap
func updatedLinkedJobSchedules(ctx context.Context, subscriptionID string, client *jobschedule.JobScheduleClient, id *runbook.RunbookId, jsMap map[string]jobschedule.JobScheduleCreateParameters) error {
	automationAccountId := jobschedule.NewAutomationAccountID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName)
	filter := fmt.Sprintf("properties/runbook/name eq '%s'", id.RunbookName)
	jsIterator, err := client.ListByAutomationAccount(ctx, automationAccountId, jobschedule.ListByAutomationAccountOperationOptions{Filter: &filter})
	if err != nil {
		return fmt.Errorf("loading Automation Account %q Job Schedule List: %+v", id.AutomationAccountName, err)
	}

	for _, item := range pointer.From(jsIterator.Model) {
		prop := item.Properties
		jobDigest := helper.ResourceAutomationJobScheduleDigest(prop)

		if _, ok := jsMap[jobDigest]; ok {
			delete(jsMap, jobDigest)
		} else {
			if prop == nil || prop.JobScheduleId == nil || *prop.JobScheduleId == "" {
				return fmt.Errorf("job schedule Id is nil or empty listed by %s Job Schedule List: %+v", id, err)
			}
			parsedId := jobschedule.NewJobScheduleID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, pointer.From(item.Properties.JobScheduleId))
			if resp, err := client.Delete(ctx, parsedId); err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("deleting job schedule Id listed by %s Job Schedule List:%v", id, err)
				}
			}
		}
	}

	// create jobs still in jsMap
	for _, js := range jsMap {
		// skip if the schedule name is empty
		if pointer.From(js.Properties.Schedule.Name) == "" {
			continue
		}
		jsuuid, err := uuid.NewV4()
		if err != nil {
			return fmt.Errorf("creating job schedule Id(UUID) for %s: %+v", id, err)
		}

		jsId := jobschedule.NewJobScheduleID(subscriptionID, id.ResourceGroupName, id.AutomationAccountName, jsuuid.String())
		if _, err := client.Create(ctx, jsId, js); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
	}

	return nil
}
