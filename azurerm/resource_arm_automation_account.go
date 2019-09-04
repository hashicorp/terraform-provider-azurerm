package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAutomationAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationAccountCreateUpdate,
		Read:   resourceArmAutomationAccountRead,
		Update: resourceArmAutomationAccountCreateUpdate,
		Delete: resourceArmAutomationAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					//todo this will not allow single character names, even thou they are valid
					regexp.MustCompile(`^[0-9a-zA-Z]([-0-9a-zA-Z]{0,48}[0-9a-zA-Z])?$`),
					`The account name must not be empty, and must not exceed 50 characters in length.  The account name must start with a letter or number.  The account name can contain letters, numbers, and dashes. The final character must be a letter or a number.`,
				),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			// Remove in 2.0
			"sku": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				Deprecated:    "This property has been deprecated in favour of the 'sku_name' property and will be removed in version 2.0 of the provider",
				ConflictsWith: []string{"sku_name"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(automation.Basic),
								string(automation.Free),
							}, true),
						},
					},
				},
			},

			"sku_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"sku"},
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

func resourceArmAutomationAccountCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automation.AccountClient
	ctx := meta.(*ArmClient).StopContext

	// Remove in 2.0
	var sku automation.Sku

	if inputs := d.Get("sku").([]interface{}); len(inputs) != 0 {
		input := inputs[0].(map[string]interface{})
		v := input["name"].(string)

		sku = automation.Sku{
			Name: automation.SkuNameEnum(v),
		}
	} else {
		// Keep in 2.0
		sku = automation.Sku{
			Name: automation.SkuNameEnum(d.Get("sku_name").(string)),
		}
	}

	if sku.Name == "" {
		return fmt.Errorf("either 'sku_name' or 'sku' must be defined in the configuration file")
	}

	log.Printf("[INFO] preparing arguments for Automation Account create/update.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
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

	return resourceArmAutomationAccountRead(d, meta)
}

func resourceArmAutomationAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automation.AccountClient
	registrationClient := meta.(*ArmClient).automation.AgentRegistrationInfoClient
	ctx := meta.(*ArmClient).StopContext

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
		// Remove in 2.0
		if err := d.Set("sku", flattenAutomationAccountSku(sku)); err != nil {
			return fmt.Errorf("Error setting 'sku': %+v", err)
		}

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

func resourceArmAutomationAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automation.AccountClient
	ctx := meta.(*ArmClient).StopContext

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

// Remove in 2.0
func flattenAutomationAccountSku(sku *automation.Sku) []interface{} {
	if sku == nil {
		return []interface{}{}
	}

	result := map[string]interface{}{}
	result["name"] = string(sku.Name)
	return []interface{}{result}
}
