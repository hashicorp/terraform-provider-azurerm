package automation

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2018-06-30-preview/automation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/automation/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceAutomationModule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationModuleCreateUpdate,
		Read:   resourceAutomationModuleRead,
		Update: resourceAutomationModuleCreateUpdate,
		Delete: resourceAutomationModuleDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"automation_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AutomationAccount(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"module_link": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"uri": {
							Type:     pluginsdk.TypeString,
							Required: true,
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
		},
	}
}

func resourceAutomationModuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.ModuleClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Automation Module creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	accName := d.Get("automation_account_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, accName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Automation Module %q (Account %q / Resource Group %q): %s", name, accName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_automation_module", *existing.ID)
		}
	}

	contentLink := expandModuleLink(d)

	parameters := automation.ModuleCreateOrUpdateParameters{
		ModuleCreateOrUpdateProperties: &automation.ModuleCreateOrUpdateProperties{
			ContentLink: &contentLink,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resGroup, accName, name, parameters); err != nil {
		return err
	}

	// the API returns 'done' but it's not actually finished provisioning yet
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{
			string(automation.ModuleProvisioningStateActivitiesStored),
			string(automation.ModuleProvisioningStateConnectionTypeImported),
			string(automation.ModuleProvisioningStateContentDownloaded),
			string(automation.ModuleProvisioningStateContentRetrieved),
			string(automation.ModuleProvisioningStateContentStored),
			string(automation.ModuleProvisioningStateContentValidated),
			string(automation.ModuleProvisioningStateCreated),
			string(automation.ModuleProvisioningStateCreating),
			string(automation.ModuleProvisioningStateModuleDataStored),
			string(automation.ModuleProvisioningStateModuleImportRunbookComplete),
			string(automation.ModuleProvisioningStateRunningImportModuleRunbook),
			string(automation.ModuleProvisioningStateStartingImportModuleRunbook),
			string(automation.ModuleProvisioningStateUpdating),
		},
		Target: []string{
			string(automation.ModuleProvisioningStateSucceeded),
		},
		MinTimeout: 30 * time.Second,
		Refresh: func() (interface{}, string, error) {
			resp, err2 := client.Get(ctx, resGroup, accName, name)
			if err2 != nil {
				return resp, "Error", fmt.Errorf("Error retrieving Module %q (Automation Account %q / Resource Group %q): %+v", name, accName, resGroup, err2)
			}

			if properties := resp.ModuleProperties; properties != nil {
				if properties.Error != nil && properties.Error.Message != nil && *properties.Error.Message != "" {
					return resp, string(properties.ProvisioningState), fmt.Errorf(*properties.Error.Message)
				}
				return resp, string(properties.ProvisioningState), nil
			}

			return resp, "Unknown", nil
		},
	}
	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("Error waiting for Module %q (Automation Account %q / Resource Group %q) to finish provisioning: %+v", name, accName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, accName, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Automation Module %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceAutomationModuleRead(d, meta)
}

func resourceAutomationModuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.ModuleClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	accName := id.Path["automationAccounts"]
	name := id.Path["modules"]

	resp, err := client.Get(ctx, resGroup, accName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on AzureRM Automation Module %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("automation_account_name", accName)

	return nil
}

func resourceAutomationModuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.ModuleClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	accName := id.Path["automationAccounts"]
	name := id.Path["modules"]

	resp, err := client.Delete(ctx, resGroup, accName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error issuing AzureRM delete request for Automation Module %q: %+v", name, err)
	}

	return nil
}

func expandModuleLink(d *pluginsdk.ResourceData) automation.ContentLink {
	inputs := d.Get("module_link").([]interface{})
	input := inputs[0].(map[string]interface{})
	uri := input["uri"].(string)

	hashes := input["hash"].([]interface{})

	if len(hashes) > 0 {
		hash := hashes[0].(map[string]interface{})
		hashValue := hash["value"].(string)
		hashAlgorithm := hash["algorithm"].(string)

		return automation.ContentLink{
			URI: &uri,
			ContentHash: &automation.ContentHash{
				Algorithm: &hashAlgorithm,
				Value:     &hashValue,
			},
		}
	}

	return automation.ContentLink{
		URI: &uri,
	}
}
