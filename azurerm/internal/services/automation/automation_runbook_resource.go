package automation

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2018-06-30-preview/automation"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/automation/helper"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/automation/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceAutomationRunbook() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutomationRunbookCreateUpdate,
		Read:   resourceAutomationRunbookRead,
		Update: resourceAutomationRunbookCreateUpdate,
		Delete: resourceAutomationRunbookDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.RunbookName(),
			},

			"automation_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AutomationAccount(),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"runbook_type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(automation.Graph),
					string(automation.GraphPowerShell),
					string(automation.GraphPowerShellWorkflow),
					string(automation.PowerShell),
					string(automation.PowerShellWorkflow),
					string(automation.Script),
				}, true),
			},

			"log_progress": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"log_verbose": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"content", "publish_content_link"},
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"job_schedule": helper.JobScheduleSchema(),

			"publish_content_link": {
				Type:         schema.TypeList,
				Optional:     true,
				MaxItems:     1,
				AtLeastOneOf: []string{"content", "publish_content_link"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uri": {
							Type:     schema.TypeString,
							Required: true,
						},

						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"hash": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"algorithm": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
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

func resourceAutomationRunbookCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.RunbookClient
	jsClient := meta.(*clients.Client).Automation.JobScheduleClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Automation Runbook creation.")

	name := d.Get("name").(string)
	accName := d.Get("automation_account_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, accName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Automation Runbook %q (Account %q / Resource Group %q): %s", name, accName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_automation_runbook", *existing.ID)
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

	if _, err := client.CreateOrUpdate(ctx, resGroup, accName, name, parameters); err != nil {
		return fmt.Errorf("Error creating/updating Automation Runbook %q (Account %q / Resource Group %q): %+v", name, accName, resGroup, err)
	}

	if v, ok := d.GetOk("content"); ok {
		content := v.(string)
		reader := io.NopCloser(bytes.NewBufferString(content))
		draftClient := meta.(*clients.Client).Automation.RunbookDraftClient

		if _, err := draftClient.ReplaceContent(ctx, resGroup, accName, name, reader); err != nil {
			return fmt.Errorf("Error setting the draft Automation Runbook %q (Account %q / Resource Group %q): %+v", name, accName, resGroup, err)
		}

		if _, err := client.Publish(ctx, resGroup, accName, name); err != nil {
			return fmt.Errorf("Error publishing the updated Automation Runbook %q (Account %q / Resource Group %q): %+v", name, accName, resGroup, err)
		}
	}

	read, err := client.Get(ctx, resGroup, accName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Automation Runbook %q (Account %q / Resource Group %q): %+v", name, accName, resGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Automation Runbook %q (Account %q / Resource Group %q) ID", name, accName, resGroup)
	}

	d.SetId(*read.ID)

	for jsIterator, err := jsClient.ListByAutomationAccountComplete(ctx, resGroup, accName, ""); jsIterator.NotDone(); err = jsIterator.NextWithContext(ctx) {
		if err != nil {
			return fmt.Errorf("loading Automation Account %q Job Schedule List: %+v", accName, err)
		}
		if props := jsIterator.Value().JobScheduleProperties; props != nil {
			if props.Runbook.Name != nil && *props.Runbook.Name == name {
				if jsIterator.Value().JobScheduleID == nil || *jsIterator.Value().JobScheduleID == "" {
					return fmt.Errorf("job schedule Id is nil or empty listed by Automation Account %q Job Schedule List: %+v", accName, err)
				}
				jsId, err := uuid.FromString(*jsIterator.Value().JobScheduleID)
				if err != nil {
					return fmt.Errorf("parsing job schedule Id listed by Automation Account %q Job Schedule List:%v", accName, err)
				}
				if _, err := jsClient.Delete(ctx, resGroup, accName, jsId); err != nil {
					return fmt.Errorf("deleting job schedule Id listed by Automation Account %q Job Schedule List:%v", accName, err)
				}
			}
		}
	}

	if v, ok := d.GetOk("job_schedule"); ok {
		jsMap, err := helper.ExpandAutomationJobSchedule(v.(*schema.Set).List(), name)
		if err != nil {
			return err
		}
		for jsuuid, js := range *jsMap {
			if _, err := jsClient.Create(ctx, resGroup, accName, jsuuid, js); err != nil {
				return fmt.Errorf("creating Automation Runbook %q Job Schedules (Account %q / Resource Group %q): %+v", name, accName, resGroup, err)
			}
		}
	}

	return resourceAutomationRunbookRead(d, meta)
}

func resourceAutomationRunbookRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.RunbookClient
	jsClient := meta.(*clients.Client).Automation.JobScheduleClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	accName := id.Path["automationAccounts"]
	name := id.Path["runbooks"]

	resp, err := client.Get(ctx, resGroup, accName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on AzureRM Automation Runbook %q (Account %q / Resource Group %q): %+v", name, accName, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("automation_account_name", accName)
	if props := resp.RunbookProperties; props != nil {
		d.Set("log_verbose", props.LogVerbose)
		d.Set("log_progress", props.LogProgress)
		d.Set("runbook_type", props.RunbookType)
		d.Set("description", props.Description)
	}

	response, err := client.GetContent(ctx, resGroup, accName, name)
	if err != nil {
		if utils.ResponseWasNotFound(response.Response) {
			d.Set("content", "")
		} else {
			return fmt.Errorf("retrieving content for Automation Runbook %q (Account %q / Resource Group %q): %+v", name, accName, resGroup, err)
		}
	}

	if v := response.Value; v != nil {
		if contentBytes := *response.Value; contentBytes != nil {
			buf := new(bytes.Buffer)
			if _, err := buf.ReadFrom(contentBytes); err != nil {
				return fmt.Errorf("Error reading from Automation Runbook buffer %q: %+v", name, err)
			}
			content := buf.String()
			d.Set("content", content)
		}
	}

	jsMap := make(map[uuid.UUID]automation.JobScheduleProperties)
	for jsIterator, err := jsClient.ListByAutomationAccountComplete(ctx, resGroup, accName, ""); jsIterator.NotDone(); err = jsIterator.NextWithContext(ctx) {
		if err != nil {
			return fmt.Errorf("loading Automation Account %q Job Schedule List: %+v", accName, err)
		}
		if props := jsIterator.Value().JobScheduleProperties; props != nil {
			if props.Runbook.Name != nil && *props.Runbook.Name == name {
				if jsIterator.Value().JobScheduleID == nil || *jsIterator.Value().JobScheduleID == "" {
					return fmt.Errorf("job schedule Id is nil or empty listed by Automation Account %q Job Schedule List: %+v", accName, err)
				}
				jsId, err := uuid.FromString(*jsIterator.Value().JobScheduleID)
				if err != nil {
					return fmt.Errorf("parsing job schedule Id listed by Automation Account %q Job Schedule List:%v", accName, err)
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

func resourceAutomationRunbookDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.RunbookClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	accName := id.Path["automationAccounts"]
	name := id.Path["runbooks"]

	resp, err := client.Delete(ctx, resGroup, accName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error issuing AzureRM delete request for Automation Runbook '%s': %+v", name, err)
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
