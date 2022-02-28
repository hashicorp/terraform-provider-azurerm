package automation

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2020-01-13-preview/automation"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAutomationRunbook() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationRunbookCreateUpdate,
		Read:   resourceAutomationRunbookRead,
		Update: resourceAutomationRunbookCreateUpdate,
		Delete: resourceAutomationRunbookDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.RunbookID(id)
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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"runbook_type": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifferenceV2Only,
				ValidateFunc: validation.StringInSlice([]string{
					string(automation.RunbookTypeEnumGraph),
					string(automation.RunbookTypeEnumGraphPowerShell),
					string(automation.RunbookTypeEnumGraphPowerShellWorkflow),
					string(automation.RunbookTypeEnumPowerShell),
					string(automation.RunbookTypeEnumPowerShellWorkflow),
					string(automation.RunbookTypeEnumScript),
				}, !features.ThreePointOh()),
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
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"content", "publish_content_link"},
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"job_schedule": helper.JobScheduleSchema(),

			"publish_content_link": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				MaxItems:     1,
				AtLeastOneOf: []string{"content", "publish_content_link"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"uri": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"version": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"hash": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"algorithm": {
										Type:     pluginsdk.TypeString,
										Required: true,
									},
									"value": {
										Type:     pluginsdk.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceAutomationRunbookCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.RunbookClient
	jsClient := meta.(*clients.Client).Automation.JobScheduleClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Automation Runbook creation.")

	id := parse.NewRunbookID(client.SubscriptionID, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_automation_runbook", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	runbookType := automation.RunbookTypeEnum(d.Get("runbook_type").(string))
	logProgress := d.Get("log_progress").(bool)
	logVerbose := d.Get("log_verbose").(bool)
	description := d.Get("description").(string)

	parameters := automation.RunbookCreateOrUpdateParameters{
		RunbookCreateOrUpdateProperties: &automation.RunbookCreateOrUpdateProperties{
			LogVerbose:  &logVerbose,
			LogProgress: &logProgress,
			RunbookType: runbookType,
			Description: &description,
		},

		Location: &location,
		Tags:     tags.Expand(t),
	}

	contentLink := expandContentLink(d.Get("publish_content_link").([]interface{}))
	if contentLink != nil {
		parameters.RunbookCreateOrUpdateProperties.PublishContentLink = contentLink
	} else {
		parameters.RunbookCreateOrUpdateProperties.Draft = &automation.RunbookDraft{}
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if v, ok := d.GetOk("content"); ok {
		content := v.(string)
		reader := io.NopCloser(bytes.NewBufferString(content))
		draftClient := meta.(*clients.Client).Automation.RunbookDraftClient

		if _, err := draftClient.ReplaceContent(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name, reader); err != nil {
			return fmt.Errorf("setting the draft for %s: %+v", id, err)
		}

		if _, err := client.Publish(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name); err != nil {
			return fmt.Errorf("publishing the updated %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	for jsIterator, err := jsClient.ListByAutomationAccountComplete(ctx, id.ResourceGroup, id.AutomationAccountName, ""); jsIterator.NotDone(); err = jsIterator.NextWithContext(ctx) {
		if err != nil {
			return fmt.Errorf("loading %s Job Schedule List: %+v", id, err)
		}
		if props := jsIterator.Value().JobScheduleProperties; props != nil {
			if props.Runbook.Name != nil && *props.Runbook.Name == id.Name {
				if jsIterator.Value().JobScheduleID == nil || *jsIterator.Value().JobScheduleID == "" {
					return fmt.Errorf("job schedule Id is nil or empty listed by %s Job Schedule List: %+v", id, err)
				}
				jsId, err := uuid.FromString(*jsIterator.Value().JobScheduleID)
				if err != nil {
					return fmt.Errorf("parsing job schedule Id listed by %s Job Schedule List:%v", id, err)
				}
				if resp, err := jsClient.Delete(ctx, id.ResourceGroup, id.AutomationAccountName, jsId); err != nil {
					if !utils.ResponseWasNotFound(resp) {
						return fmt.Errorf("deleting job schedule Id listed by %s Job Schedule List:%v", id, err)
					}
				}
			}
		}
	}

	if v, ok := d.GetOk("job_schedule"); ok {
		jsMap, err := helper.ExpandAutomationJobSchedule(v.(*pluginsdk.Set).List(), id.Name)
		if err != nil {
			return err
		}
		for jsuuid, js := range *jsMap {
			if _, err := jsClient.Create(ctx, id.ResourceGroup, id.AutomationAccountName, jsuuid, js); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
		}
	}

	return resourceAutomationRunbookRead(d, meta)
}

func resourceAutomationRunbookRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.RunbookClient
	jsClient := meta.(*clients.Client).Automation.JobScheduleClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.RunbookID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on AzureRM Automation Runbook %q (Account %q / Resource Group %q): %+v", id.Name, id.AutomationAccountName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("automation_account_name", id.AutomationAccountName)
	if props := resp.RunbookProperties; props != nil {
		d.Set("log_verbose", props.LogVerbose)
		d.Set("log_progress", props.LogProgress)
		d.Set("runbook_type", props.RunbookType)
		d.Set("description", props.Description)
	}

	response, err := client.GetContent(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(response.Response) {
			d.Set("content", "")
		} else {
			return fmt.Errorf("retrieving content for Automation Runbook %q (Account %q / Resource Group %q): %+v", id.Name, id.AutomationAccountName, id.ResourceGroup, err)
		}
	}

	if v := response.Value; v != nil {
		if contentBytes := *response.Value; contentBytes != nil {
			buf := new(bytes.Buffer)
			if _, err := buf.ReadFrom(contentBytes); err != nil {
				return fmt.Errorf("reading from Automation Runbook buffer %q: %+v", id.Name, err)
			}
			content := buf.String()
			d.Set("content", content)
		}
	}

	jsMap := make(map[uuid.UUID]automation.JobScheduleProperties)
	for jsIterator, err := jsClient.ListByAutomationAccountComplete(ctx, id.ResourceGroup, id.AutomationAccountName, ""); jsIterator.NotDone(); err = jsIterator.NextWithContext(ctx) {
		if err != nil {
			return fmt.Errorf("loading Automation Account %q Job Schedule List: %+v", id.AutomationAccountName, err)
		}
		if props := jsIterator.Value().JobScheduleProperties; props != nil {
			if props.Runbook.Name != nil && *props.Runbook.Name == id.Name {
				if jsIterator.Value().JobScheduleID == nil || *jsIterator.Value().JobScheduleID == "" {
					return fmt.Errorf("job schedule Id is nil or empty listed by Automation Account %q Job Schedule List: %+v", id.AutomationAccountName, err)
				}
				jsId, err := uuid.FromString(*jsIterator.Value().JobScheduleID)
				if err != nil {
					return fmt.Errorf("parsing job schedule Id listed by Automation Account %q Job Schedule List:%v", id.AutomationAccountName, err)
				}
				jsMap[jsId] = *props
			}
		}
	}

	jobSchedule := helper.FlattenAutomationJobSchedule(jsMap)
	if err := d.Set("job_schedule", jobSchedule); err != nil {
		return fmt.Errorf("setting `job_schedule`: %+v", err)
	}

	if t := resp.Tags; t != nil {
		return tags.FlattenAndSet(d, t)
	}

	return nil
}

func resourceAutomationRunbookDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.RunbookClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.RunbookID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("issuing AzureRM delete request for Automation Runbook '%s': %+v", id.Name, err)
	}

	return nil
}

func expandContentLink(inputs []interface{}) *automation.ContentLink {
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

		return &automation.ContentLink{
			URI:     &uri,
			Version: &version,
			ContentHash: &automation.ContentHash{
				Algorithm: &hashAlgorithm,
				Value:     &hashValue,
			},
		}
	}

	return &automation.ContentLink{
		URI:     &uri,
		Version: &version,
	}
}
