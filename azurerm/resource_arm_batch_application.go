package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2018-12-01/batch"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmBatchApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmBatchApplicationCreate,
		Read:   resourceArmBatchApplicationRead,
		Update: resourceArmBatchApplicationUpdate,
		Delete: resourceArmBatchApplicationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMBatchApplicationName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMBatchAccountName,
			},

			"allow_updates": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"default_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAzureRMBatchApplicationVersion,
			},

			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAzureRMBatchApplicationDisplayName,
			},
		},
	}
}

func resourceArmBatchApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).batch.ApplicationClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)

	if features.ShouldResourcesBeImported() {
		resp, err := client.Get(ctx, resourceGroup, accountName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for present of existing Batch Application %q (Account Name %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_batch_application", *resp.ID)
		}
	}

	allowUpdates := d.Get("allow_updates").(bool)
	defaultVersion := d.Get("default_version").(string)
	displayName := d.Get("display_name").(string)

	parameters := batch.Application{
		ApplicationProperties: &batch.ApplicationProperties{
			AllowUpdates:   utils.Bool(allowUpdates),
			DefaultVersion: utils.String(defaultVersion),
			DisplayName:    utils.String(displayName),
		},
	}

	if _, err := client.Create(ctx, resourceGroup, accountName, name, &parameters); err != nil {
		return fmt.Errorf("Error creating Batch Application %q (Account Name %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Batch Application %q (Account Name %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Batch Application %q (Account Name %q / Resource Group %q) ID", name, accountName, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmBatchApplicationRead(d, meta)
}

func resourceArmBatchApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).batch.ApplicationClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	accountName := id.Path["batchAccounts"]
	name := id.Path["applications"]

	resp, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Batch Application %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Batch Application %q (Account Name %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("account_name", accountName)
	if applicationProperties := resp.ApplicationProperties; applicationProperties != nil {
		d.Set("allow_updates", applicationProperties.AllowUpdates)
		d.Set("default_version", applicationProperties.DefaultVersion)
		d.Set("display_name", applicationProperties.DisplayName)
	}

	return nil
}

func resourceArmBatchApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).batch.ApplicationClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)
	allowUpdates := d.Get("allow_updates").(bool)
	defaultVersion := d.Get("default_version").(string)
	displayName := d.Get("display_name").(string)

	parameters := batch.Application{
		ApplicationProperties: &batch.ApplicationProperties{
			AllowUpdates:   utils.Bool(allowUpdates),
			DefaultVersion: utils.String(defaultVersion),
			DisplayName:    utils.String(displayName),
		},
	}

	if _, err := client.Update(ctx, resourceGroup, accountName, name, parameters); err != nil {
		return fmt.Errorf("Error updating Batch Application %q (Account Name %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}

	return resourceArmBatchApplicationRead(d, meta)
}

func resourceArmBatchApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).batch.ApplicationClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	accountName := id.Path["batchAccounts"]
	name := id.Path["applications"]

	if _, err := client.Delete(ctx, resourceGroup, accountName, name); err != nil {
		return fmt.Errorf("Error deleting Batch Application %q (Account Name %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}

	return nil
}

func validateAzureRMBatchApplicationName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[-_\da-zA-Z]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q can contain any combination of alphanumeric characters, hyphens, and underscores: %q", k, value))
	}

	if 1 > len(value) {
		errors = append(errors, fmt.Errorf("%q cannot be less than 1 character: %q", k, value))
	}

	if len(value) > 64 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 64 characters: %q %d", k, value, len(value)))
	}

	return warnings, errors
}

func validateAzureRMBatchApplicationVersion(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[-._\da-zA-Z]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q can contain any combination of alphanumeric characters, hyphens, underscores, and periods: %q", k, value))
	}

	if 1 > len(value) {
		errors = append(errors, fmt.Errorf("%q cannot be less than 1 character: %q", k, value))
	}

	if len(value) > 64 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 64 characters: %q %d", k, value, len(value)))
	}

	return warnings, errors
}

func validateAzureRMBatchApplicationDisplayName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if 1 > len(value) {
		errors = append(errors, fmt.Errorf("%q cannot be less than 1 character: %q", k, value))
	}

	if len(value) > 1024 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 1024 characters: %q %d", k, value, len(value)))
	}

	return warnings, errors
}
