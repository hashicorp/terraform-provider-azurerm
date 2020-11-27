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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/automation/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/automation/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceAutomationConnectionClassicCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutomationConnectionClassicCertificateCreateUpdate,
		Read:   resourceAutomationConnectionClassicCertificateRead,
		Update: resourceAutomationConnectionClassicCertificateCreateUpdate,
		Delete: resourceAutomationConnectionClassicCertificateDelete,

		Importer: azSchema.ValidateResourceIDPriorToImportThen(func(id string) error {
			_, err := parse.ConnectionID(id)
			return err
		}, importAutomationConnection("AzureClassicCertificate")),

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

			"subscription_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"subscription_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"certificate_asset_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAutomationConnectionClassicCertificateCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.ConnectionClient
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
			return tf.ImportAsExistsError("azurerm_automation_connection_classic_certificate", *existing.ID)
		}
	}

	parameters := automation.ConnectionCreateOrUpdateParameters{
		Name: &name,
		ConnectionCreateOrUpdateProperties: &automation.ConnectionCreateOrUpdateProperties{
			Description: utils.String(d.Get("description").(string)),
			ConnectionType: &automation.ConnectionTypeAssociationProperty{
				Name: utils.String("AzureClassicCertificate"),
			},
			FieldDefinitionValues: map[string]*string{
				"SubscriptionName":     utils.String(d.Get("subscription_name").(string)),
				"SubscriptionId":       utils.String(d.Get("subscription_id").(string)),
				"CertificateAssetName": utils.String(d.Get("certificate_asset_name").(string)),
			},
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

	return resourceAutomationConnectionClassicCertificateRead(d, meta)
}

func resourceAutomationConnectionClassicCertificateRead(d *schema.ResourceData, meta interface{}) error {
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
	d.Set("description", resp.Description)

	if props := resp.ConnectionProperties; props != nil {
		if v, ok := props.FieldDefinitionValues["CertificateAssetName"]; ok {
			d.Set("certificate_asset_name", v)
		}
		if v, ok := props.FieldDefinitionValues["SubscriptionId"]; ok {
			d.Set("subscription_id", v)
		}
		if v, ok := props.FieldDefinitionValues["SubscriptionName"]; ok {
			d.Set("subscription_name", v)
		}
	}

	return nil
}

func resourceAutomationConnectionClassicCertificateDelete(d *schema.ResourceData, meta interface{}) error {
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
