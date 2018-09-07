package azurerm

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAutomationRunbook() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationRunbookCreateUpdate,
		Read:   resourceArmAutomationRunbookRead,
		Update: resourceArmAutomationRunbookCreateUpdate,
		Delete: resourceArmAutomationRunbookDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"runbook_type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"publish_content_link": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
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

			"tags": tagsSchema(),
		},
	}
}

func resourceArmAutomationRunbookCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationRunbookClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Automation Runbook creation.")

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	accName := d.Get("account_name").(string)
	runbookType := automation.RunbookTypeEnum(d.Get("runbook_type").(string))
	logProgress := d.Get("log_progress").(bool)
	logVerbose := d.Get("log_verbose").(bool)
	description := d.Get("description").(string)

	contentLink := expandContentLink(d)

	parameters := automation.RunbookCreateOrUpdateParameters{
		RunbookCreateOrUpdateProperties: &automation.RunbookCreateOrUpdateProperties{
			LogVerbose:         &logVerbose,
			LogProgress:        &logProgress,
			RunbookType:        runbookType,
			Description:        &description,
			PublishContentLink: &contentLink,
		},

		Location: &location,
		Tags:     expandTags(tags),
	}

	_, err := client.CreateOrUpdate(ctx, resGroup, accName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating/updating Automation Runbook %q (Account %q / Resource Group %q): %+v", name, accName, resGroup, err)
	}

	if v, ok := d.GetOk("content"); ok {
		content := v.(string)
		reader := ioutil.NopCloser(bytes.NewBufferString(content))
		draftClient := meta.(*ArmClient).automationRunbookDraftClient
		_, err := draftClient.ReplaceContent(ctx, resGroup, accName, name, reader)
		if err != nil {
			return fmt.Errorf("Error setting the draft Automation Runbook %q (Account %q / Resource Group %q): %+v", name, accName, resGroup, err)
		}

		_, err = draftClient.Publish(ctx, resGroup, accName, name)
		if err != nil {
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

	return resourceArmAutomationRunbookRead(d, meta)
}

func resourceArmAutomationRunbookRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationRunbookClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	d.Set("account_name", accName)
	if props := resp.RunbookProperties; props != nil {
		d.Set("log_verbose", props.LogVerbose)
		d.Set("log_progress", props.LogProgress)
		d.Set("runbook_type", props.RunbookType)
		d.Set("description", props.Description)
	}

	if tags := resp.Tags; tags != nil {
		flattenAndSetTags(d, tags)
	}

	response, err := client.GetContent(ctx, resGroup, accName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving content for Automation Runbook %q (Account %q / Resource Group %q): %+v", name, accName, resGroup, err)
	}

	if v := response.Value; v != nil {
		if contentBytes := *response.Value; contentBytes != nil {
			buf := new(bytes.Buffer)
			buf.ReadFrom(contentBytes)
			content := buf.String()
			d.Set("content", content)
		}
	}

	return nil
}

func resourceArmAutomationRunbookDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationRunbookClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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

func expandContentLink(d *schema.ResourceData) automation.ContentLink {
	inputs := d.Get("publish_content_link").([]interface{})
	input := inputs[0].(map[string]interface{})
	uri := input["uri"].(string)
	version := input["version"].(string)

	hashes := input["hash"].([]interface{})

	if len(hashes) > 0 {
		hash := hashes[0].(map[string]interface{})
		hashValue := hash["value"].(string)
		hashAlgorithm := hash["algorithm"].(string)

		return automation.ContentLink{
			URI:     &uri,
			Version: &version,
			ContentHash: &automation.ContentHash{
				Algorithm: &hashAlgorithm,
				Value:     &hashValue,
			},
		}
	}

	return automation.ContentLink{
		URI:     &uri,
		Version: &version,
	}
}
