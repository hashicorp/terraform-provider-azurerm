package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"sku": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          string(automation.Basic),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							ValidateFunc: validation.StringInSlice([]string{
								string(automation.Basic),
							}, true),
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmAutomationAccountCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationAccountClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Automation Account creation.")

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	sku := expandAutomationAccountSku(d)

	parameters := automation.AccountCreateOrUpdateParameters{
		AccountCreateOrUpdateProperties: &automation.AccountCreateOrUpdateProperties{
			Sku: &sku,
		},

		Location: &location,
		Tags:     expandTags(tags),
	}

	_, err := client.CreateOrUpdate(ctx, resGroup, name, parameters)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Automation Account '%s' (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAutomationAccountRead(d, meta)
}

func resourceArmAutomationAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationAccountClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["automationAccounts"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on AzureRM Automation Account '%s': %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	flattenAndSetAutomationAccountSku(d, resp.Sku)

	if tags := resp.Tags; tags != nil {
		flattenAndSetTags(d, tags)
	}

	return nil
}

func resourceArmAutomationAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationAccountClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["automationAccounts"]

	resp, err := client.Delete(ctx, resGroup, name)

	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error issuing AzureRM delete request for Automation Account '%s': %+v", name, err)
	}

	return nil
}

func flattenAndSetAutomationAccountSku(d *schema.ResourceData, sku *automation.Sku) {
	results := make([]interface{}, 1)

	result := map[string]interface{}{}
	result["name"] = string(sku.Name)
	results[0] = result

	d.Set("sku", &results)
}

func expandAutomationAccountSku(d *schema.ResourceData) automation.Sku {
	inputs := d.Get("sku").([]interface{})
	input := inputs[0].(map[string]interface{})
	name := automation.SkuNameEnum(input["name"].(string))

	sku := automation.Sku{
		Name: name,
	}

	return sku
}
