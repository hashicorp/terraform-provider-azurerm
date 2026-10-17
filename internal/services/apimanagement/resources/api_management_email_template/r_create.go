// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_email_template

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/emailtemplates"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementEmailTemplateCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.EmailTemplatesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	id := emailtemplates.NewTemplateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), emailtemplates.TemplateName(d.Get("template_name").(string)))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
