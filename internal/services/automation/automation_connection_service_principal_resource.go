package automation

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2020-01-13-preview/automation"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAutomationConnectionServicePrincipal() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationConnectionServicePrincipalCreateUpdate,
		Read:   resourceAutomationConnectionServicePrincipalRead,
		Update: resourceAutomationConnectionServicePrincipalCreateUpdate,
		Delete: resourceAutomationConnectionServicePrincipalDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.ConnectionID(id)
			return err
		}, importAutomationConnection("AzureServicePrincipal")),

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
				ValidateFunc: validate.ConnectionName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"automation_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AutomationAccount(),
			},

			"application_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"certificate_thumbprint": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"subscription_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"tenant_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAutomationConnectionServicePrincipalCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.ConnectionClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Automation Connection creation.")

	id := parse.NewConnectionID(client.SubscriptionID, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_automation_connection_service_principal", id.ID())
		}
	}

	parameters := automation.ConnectionCreateOrUpdateParameters{
		Name: &id.Name,
		ConnectionCreateOrUpdateProperties: &automation.ConnectionCreateOrUpdateProperties{
			Description: utils.String(d.Get("description").(string)),
			ConnectionType: &automation.ConnectionTypeAssociationProperty{
				Name: utils.String("AzureServicePrincipal"),
			},
			FieldDefinitionValues: map[string]*string{
				"ApplicationId":         utils.String(d.Get("application_id").(string)),
				"CertificateThumbprint": utils.String(d.Get("certificate_thumbprint").(string)),
				"SubscriptionId":        utils.String(d.Get("subscription_id").(string)),
				"TenantId":              utils.String(d.Get("tenant_id").(string)),
			},
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name, parameters); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceAutomationConnectionServicePrincipalRead(d, meta)
}

func resourceAutomationConnectionServicePrincipalRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		if v, ok := props.FieldDefinitionValues["ApplicationId"]; ok {
			d.Set("application_id", v)
		}
		if v, ok := props.FieldDefinitionValues["CertificateThumbprint"]; ok {
			d.Set("certificate_thumbprint", v)
		}
		if v, ok := props.FieldDefinitionValues["SubscriptionId"]; ok {
			d.Set("subscription_id", v)
		}
		if v, ok := props.FieldDefinitionValues["TenantId"]; ok {
			d.Set("tenant_id", v)
		}
	}

	return nil
}

func resourceAutomationConnectionServicePrincipalDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.ConnectionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("deleting Automation Connection '%s': %+v", id.Name, err)
	}

	return nil
}
