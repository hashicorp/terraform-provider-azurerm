// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApiManagementEmailTemplate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementEmailTemplateCreateUpdate,
		Read:   resourceApiManagementEmailTemplateRead,
		Update: resourceApiManagementEmailTemplateCreateUpdate,
		Delete: resourceApiManagementEmailTemplateDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.EmailTemplateID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": commonschema.ResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			// There is an open issue for the capitalization of the template names: https://github.com/Azure/azure-rest-api-specs/issues/13341
			"template_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					azure.TitleCase(string(apimanagement.TemplateNameAccountClosedDeveloper)),
					azure.TitleCase(string(apimanagement.TemplateNameApplicationApprovedNotificationMessage)),
					azure.TitleCase(string(apimanagement.TemplateNameConfirmSignUpIdentityDefault)),
					azure.TitleCase(string(apimanagement.TemplateNameEmailChangeIdentityDefault)),
					azure.TitleCase(string(apimanagement.TemplateNameInviteUserNotificationMessage)),
					azure.TitleCase(string(apimanagement.TemplateNameNewCommentNotificationMessage)),
					azure.TitleCase(string(apimanagement.TemplateNameNewDeveloperNotificationMessage)),
					azure.TitleCase(string(apimanagement.TemplateNameNewIssueNotificationMessage)),
					azure.TitleCase(string(apimanagement.TemplateNamePasswordResetByAdminNotificationMessage)),
					azure.TitleCase(string(apimanagement.TemplateNamePasswordResetIdentityDefault)),
					azure.TitleCase(string(apimanagement.TemplateNamePurchaseDeveloperNotificationMessage)),
					azure.TitleCase(string(apimanagement.TemplateNameQuotaLimitApproachingDeveloperNotificationMessage)),
					azure.TitleCase(string(apimanagement.TemplateNameRejectDeveloperNotificationMessage)),
					azure.TitleCase(string(apimanagement.TemplateNameRequestDeveloperNotificationMessage)),
				}, false),
			},
			"body": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"subject": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},
			// Computed:
			"title": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"description": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceApiManagementEmailTemplateCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceApiManagementEmailTemplateRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceApiManagementEmailTemplateDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
