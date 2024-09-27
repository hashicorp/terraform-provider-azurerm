// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/emailtemplates"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementEmailTemplate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementEmailTemplateCreateUpdate,
		Read:   resourceApiManagementEmailTemplateRead,
		Update: resourceApiManagementEmailTemplateCreateUpdate,
		Delete: resourceApiManagementEmailTemplateDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := emailtemplates.ParseTemplateIDInsensitively(id)
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
					azure.TitleCase(string(emailtemplates.TemplateNameAccountClosedDeveloper)),
					azure.TitleCase(string(emailtemplates.TemplateNameApplicationApprovedNotificationMessage)),
					azure.TitleCase(string(emailtemplates.TemplateNameConfirmSignUpIdentityDefault)),
					azure.TitleCase(string(emailtemplates.TemplateNameEmailChangeIdentityDefault)),
					azure.TitleCase(string(emailtemplates.TemplateNameInviteUserNotificationMessage)),
					azure.TitleCase(string(emailtemplates.TemplateNameNewCommentNotificationMessage)),
					azure.TitleCase(string(emailtemplates.TemplateNameNewDeveloperNotificationMessage)),
					azure.TitleCase(string(emailtemplates.TemplateNameNewIssueNotificationMessage)),
					azure.TitleCase(string(emailtemplates.TemplateNamePasswordResetByAdminNotificationMessage)),
					azure.TitleCase(string(emailtemplates.TemplateNamePasswordResetIdentityDefault)),
					azure.TitleCase(string(emailtemplates.TemplateNamePurchaseDeveloperNotificationMessage)),
					azure.TitleCase(string(emailtemplates.TemplateNameQuotaLimitApproachingDeveloperNotificationMessage)),
					azure.TitleCase(string(emailtemplates.TemplateNameRejectDeveloperNotificationMessage)),
					azure.TitleCase(string(emailtemplates.TemplateNameRequestDeveloperNotificationMessage)),
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
	client := meta.(*clients.Client).ApiManagement.EmailTemplatesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	id := emailtemplates.NewTemplateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), emailtemplates.TemplateName(d.Get("template_name").(string)))
	if d.IsNewResource() {
		existing, err := client.EmailTemplateGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		// in case the template has been edited (is not default anymore) this errors and the resource should be imported manually into the state (terraform import).
		if !response.WasNotFound(existing.HttpResponse) {
			if model := existing.Model; model != nil && model.Properties != nil && model.Properties.IsDefault != nil && !*model.Properties.IsDefault {
				return tf.ImportAsExistsError("azurerm_api_management_email_template", id.ID())
			}
		}
	}

	subject := d.Get("subject").(string)
	body := d.Get("body").(string)

	emailTemplateUpdateParameters := emailtemplates.EmailTemplateUpdateParameters{
		Properties: &emailtemplates.EmailTemplateUpdateParameterProperties{
			Subject: pointer.To(subject),
			Body:    pointer.To(body),
		},
	}

	if _, err := client.EmailTemplateCreateOrUpdate(ctx, id, emailTemplateUpdateParameters, emailtemplates.EmailTemplateCreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementEmailTemplateRead(d, meta)
}

func resourceApiManagementEmailTemplateRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.EmailTemplatesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := emailtemplates.ParseTemplateIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	templateName := emailtemplates.TemplateName(azure.TitleCase(string(id.TemplateName)))
	newId := emailtemplates.NewTemplateID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, templateName)
	resp, err := client.EmailTemplateGet(ctx, newId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s does not exist - removing from state!", newId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", newId, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("api_management_name", id.ServiceName)
	d.Set("template_name", templateName)
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("title", pointer.From(props.Title))
			d.Set("description", pointer.From(props.Description))
			d.Set("subject", props.Subject)
			d.Set("body", props.Body)
		}
	}

	return nil
}

func resourceApiManagementEmailTemplateDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.EmailTemplatesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := emailtemplates.ParseTemplateIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	templateName := emailtemplates.TemplateName(azure.TitleCase(string(id.TemplateName)))
	newId := emailtemplates.NewTemplateID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, templateName)

	if resp, err := client.EmailTemplateDelete(ctx, newId, emailtemplates.EmailTemplateDeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %s", newId, err)
		}
	}

	return nil
}
