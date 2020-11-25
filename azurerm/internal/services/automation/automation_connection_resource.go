package automation

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2018-06-30-preview/automation"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/automation/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/automation/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAutomationConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationConnectionCreateUpdate,
		Read:   resourceArmAutomationConnectionRead,
		Update: resourceArmAutomationConnectionCreateUpdate,
		Delete: resourceArmAutomationConnectionDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ConnectionID(id)
			return err
		}),

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
				ValidateFunc: validate.AutomationConnectionName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"automation_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AutomationAccountName(),
			},

			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"values": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceArmAutomationConnectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.ConnectionClient
	connectionTypeClient := meta.(*clients.Client).Automation.ConnectionTypeClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Automation Connection creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("automation_account_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, accountName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Automation Connection %q (Account %q / Resource Group %q): %s", name, accountName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_automation_connection", *existing.ID)
		}
	}

	connectionTypeName := d.Get("type").(string)
	values := utils.ExpandMapStringPtrString(d.Get("values").(map[string]interface{}))

	// check `type` exists and required fields are passed by users
	connectionType, err := connectionTypeClient.Get(ctx, resGroup, accountName, connectionTypeName)
	if err != nil {
		return fmt.Errorf("retrieving Automation Connection type %q (Account %q / Resource Group %q): %s", connectionTypeName, accountName, resGroup, err)
	}
	if connectionType.ConnectionTypeProperties != nil && connectionType.ConnectionTypeProperties.FieldDefinitions != nil {
		var missingFields []string
		for key := range connectionType.ConnectionTypeProperties.FieldDefinitions {
			if _, ok := values[key]; !ok {
				missingFields = append(missingFields, key)
			}
		}
		if len(missingFields) > 0 {
			return fmt.Errorf("%q should be specified in `values` when type is %q for `azurerm_automation_connection`", strings.Join(missingFields, ", "), connectionTypeName)
		}
	}

	parameters := automation.ConnectionCreateOrUpdateParameters{
		Name: &name,
		ConnectionCreateOrUpdateProperties: &automation.ConnectionCreateOrUpdateProperties{
			Description: utils.String(d.Get("description").(string)),
			ConnectionType: &automation.ConnectionTypeAssociationProperty{
				Name: utils.String(connectionTypeName),
			},
			FieldDefinitionValues: values,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resGroup, accountName, name, parameters); err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, accountName, name)
	if err != nil {
		return err
	}

	if read.ID == nil || *read.ID == "" {
		return fmt.Errorf("empty or nil ID for Automation Connection '%s' (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAutomationConnectionRead(d, meta)
}

func resourceArmAutomationConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.ConnectionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Read request on AzureRM Automation Connection '%s': %+v", id.Name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("automation_account_name", id.AutomationAccountName)
	d.Set("values", resp.FieldDefinitionValues)
	d.Set("description", resp.Description)

	if props := resp.ConnectionProperties; props != nil {
		if props.ConnectionType != nil {
			d.Set("type", resp.ConnectionType.Name)
		}

		if err := d.Set("values", utils.FlattenMapStringPtrString(props.FieldDefinitionValues)); err != nil {
			return fmt.Errorf("setting `values`: %+v", err)
		}
	}

	return nil
}

func resourceArmAutomationConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.ConnectionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}

		return fmt.Errorf("deleting Automation Connection '%s': %+v", id.Name, err)
	}

	return nil
}
