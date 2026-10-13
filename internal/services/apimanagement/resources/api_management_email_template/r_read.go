// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_email_template

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/emailtemplates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func resourceApiManagementEmailTemplateRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.EmailTemplatesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := emailtemplates.ParseTemplateIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	templateName := emailtemplates.TemplateName(titleCase(string(id.TemplateName)))
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

// titleCase upper-cases the first letter and leaves the rest alone, which is
// how the API spells the template names (see the issue linked in the schema).
func titleCase(input string) string {
	return cases.Title(language.English, cases.NoLower).String(input)
}
