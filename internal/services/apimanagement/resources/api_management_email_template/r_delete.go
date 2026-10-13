// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_email_template

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/emailtemplates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementEmailTemplateDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.EmailTemplatesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := emailtemplates.ParseTemplateIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	templateName := emailtemplates.TemplateName(titleCase(string(id.TemplateName)))
	newId := emailtemplates.NewTemplateID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, templateName)

	if resp, err := client.EmailTemplateDelete(ctx, newId, emailtemplates.EmailTemplateDeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %s", newId, err)
		}
	}

	return nil
}
