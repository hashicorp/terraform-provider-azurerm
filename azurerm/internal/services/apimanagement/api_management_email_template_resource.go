package apimanagement

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementEmailTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiManagementEmailTemplateCreateUpdate,
		Read:   resourceApiManagementEmailTemplateRead,
		Update: resourceApiManagementEmailTemplateCreateUpdate,
		Delete: resourceApiManagementEmailTemplateDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			// There is an open issue for the capitalization of the template names: https://github.com/Azure/azure-rest-api-specs/issues/13341
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					strings.Title(string(apimanagement.AccountClosedDeveloper)),
					strings.Title(string(apimanagement.ApplicationApprovedNotificationMessage)),
					strings.Title(string(apimanagement.ConfirmSignUpIdentityDefault)),
					strings.Title(string(apimanagement.EmailChangeIdentityDefault)),
					strings.Title(string(apimanagement.InviteUserNotificationMessage)),
					strings.Title(string(apimanagement.NewCommentNotificationMessage)),
					strings.Title(string(apimanagement.NewDeveloperNotificationMessage)),
					strings.Title(string(apimanagement.NewIssueNotificationMessage)),
					strings.Title(string(apimanagement.PasswordResetByAdminNotificationMessage)),
					strings.Title(string(apimanagement.PasswordResetIdentityDefault)),
					strings.Title(string(apimanagement.PurchaseDeveloperNotificationMessage)),
					strings.Title(string(apimanagement.QuotaLimitApproachingDeveloperNotificationMessage)),
					strings.Title(string(apimanagement.RejectDeveloperNotificationMessage)),
					strings.Title(string(apimanagement.RequestDeveloperNotificationMessage)),
				}, false),
			},
			"body": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"subject": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Computed:
			"title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceApiManagementEmailTemplateCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.EmailTemplateClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	templateName := apimanagement.TemplateName(d.Get("template_name").(string))

	id := parse.NewEmailTemplateID(subscriptionId, resourceGroup, serviceName, d.Get("template_name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, templateName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		// in case the template has been edited (is not default anymore) this errors and the resource should be imported manually into the state (terraform import).
		if !utils.ResponseWasNotFound(existing.Response) && (existing.IsDefault != nil && !*existing.IsDefault) {
			return tf.ImportAsExistsError("azurerm_api_management_email_template", id.ID())
		}
	}

	subject := d.Get("subject").(string)
	body := d.Get("body").(string)

	emailTemplateUpdateParameters := apimanagement.EmailTemplateUpdateParameters{
		EmailTemplateUpdateParameterProperties: &apimanagement.EmailTemplateUpdateParameterProperties{
			Subject: utils.String(subject),
			Body:    utils.String(body),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, templateName, emailTemplateUpdateParameters, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementEmailTemplateRead(d, meta)
}

func resourceApiManagementEmailTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.EmailTemplateClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EmailTemplateID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName

	templateName := apimanagement.TemplateName(id.TemplateName)

	resp, err := client.Get(ctx, resourceGroup, serviceName, templateName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s does not exist - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)
	d.Set("template_name", templateName)
	if properties := resp.EmailTemplateContractProperties; properties != nil {
		d.Set("title", properties.Title)
		d.Set("description", properties.Description)
		d.Set("subject", properties.Subject)
		d.Set("body", properties.Body)
	}

	return nil
}

func resourceApiManagementEmailTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.EmailTemplateClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EmailTemplateID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	templateName := apimanagement.TemplateName(id.TemplateName)

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, templateName, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting %s: %s", *id, err)
		}
	}

	return nil
}
