package automation

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2018-06-30-preview/automation"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/automation/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceAutomationAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutomationAccountCreateUpdate,
		Read:   resourceAutomationAccountRead,
		Update: resourceAutomationAccountCreateUpdate,
		Delete: resourceAutomationAccountDelete,
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
				ValidateFunc: validate.AutomationAccount(),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(automation.Basic),
					string(automation.Free),
				}, false),
			},

			"tags": tags.Schema(),

			"dsc_server_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dsc_primary_access_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dsc_secondary_access_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAutomationAccountCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.AccountClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	sku := automation.Sku{
		Name: automation.SkuNameEnum(d.Get("sku_name").(string)),
	}

	log.Printf("[INFO] preparing arguments for Automation Account create/update.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Automation Account %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_automation_account", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	parameters := automation.AccountCreateOrUpdateParameters{
		AccountCreateOrUpdateProperties: &automation.AccountCreateOrUpdateProperties{
			Sku: &sku,
		},
		Location: utils.String(location),
		Tags:     tags.Expand(t),
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters); err != nil {
		return fmt.Errorf("Error creating/updating Automation Account %q (Resource Group %q) %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Automation Account %q (Resource Group %q) %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Automation Account %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceAutomationAccountRead(d, meta)
}

func resourceAutomationAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.AccountClient
	registrationClient := meta.(*clients.Client).Automation.AgentRegistrationInfoClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["automationAccounts"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Automation Account %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Automation Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	keysResp, err := registrationClient.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Agent Registration Info for Automation Account %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request for Agent Registration Info for Automation Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		if err := d.Set("sku_name", string(sku.Name)); err != nil {
			return fmt.Errorf("Error setting 'sku_name': %+v", err)
		}
	} else {
		return fmt.Errorf("Error making Read request on Automation Account %q (Resource Group %q): Unable to retrieve 'sku' value", name, resourceGroup)
	}

	d.Set("dsc_server_endpoint", keysResp.Endpoint)
	if keys := keysResp.Keys; keys != nil {
		d.Set("dsc_primary_access_key", keys.Primary)
		d.Set("dsc_secondary_access_key", keys.Secondary)
	}

	if t := resp.Tags; t != nil {
		return tags.FlattenAndSet(d, t)
	}

	return nil
}

func resourceAutomationAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.AccountClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["automationAccounts"]

	resp, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error issuing AzureRM delete request for Automation Account '%s': %+v", name, err)
	}

	return nil
}
