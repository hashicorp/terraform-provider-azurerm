package azurerm

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/automation"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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

			"location": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: azureRMNormalizeLocation,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

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
				Required: true,
				ForceNew: true,
			},

			"publish_content_link": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uri": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set: resourceAzureRMAutomationRunbookContentLinkHash,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAzureRMAutomationRunbookContentLinkHash(v interface{}) int {
	m := v.(map[string]interface{})

	uri := m["uri"].(string)

	return hashcode.String(uri)
}

func resourceArmAutomationRunbookCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationRunbookClient
	log.Printf("[INFO] preparing arguments for AzureRM Automation Runbook creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
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

		Name:     &name,
		Location: &location,
		Tags:     expandTags(tags),
	}

	_, err := client.CreateOrUpdate(resGroup, accName, name, parameters)
	if err != nil {
		return err
	}

	read, err := client.Get(resGroup, accName, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Automation Runbook '%s' (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAutomationRunbookRead(d, meta)
}

func resourceArmAutomationRunbookRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationRunbookClient
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	accName := id.Path["automationAccounts"]
	name := id.Path["runbooks"]

	resp, err := client.Get(resGroup, accName, name)
	if err != nil {
		if responseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on AzureRM Automation Runbook '%s': %s", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))
	d.Set("resource_group_name", resGroup)

	d.Set("account_name", accName)
	d.Set("log_verbose", resp.LogVerbose)
	d.Set("log_progress", resp.LogProgress)
	d.Set("runbook_type", resp.RunbookType)
	d.Set("description", resp.Description)
	d.Set("publish_content_link", nil) //publish content link is not set during Get()

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmAutomationRunbookDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationRunbookClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	accName := id.Path["automationAccounts"]
	name := id.Path["runbooks"]

	resp, err := client.Delete(resGroup, accName, name)

	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}

		return fmt.Errorf("Error issuing AzureRM delete request for Automation Runbook '%s': %+v", name, err)
	}

	return nil
}

func expandContentLink(d *schema.ResourceData) automation.ContentLink {
	inputs := d.Get("publish_content_link").(*schema.Set).List()
	input := inputs[0].(map[string]interface{})
	uri := input["uri"].(string)

	contentLink := automation.ContentLink{
		URI: &uri,
	}

	return contentLink
}
