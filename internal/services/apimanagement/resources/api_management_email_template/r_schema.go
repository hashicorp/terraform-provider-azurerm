// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_email_template

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/emailtemplates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementEmailTemplateSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"api_management_name": schemaz.SchemaApiManagementName(),

		// There is an open issue for the capitalization of the template names: https://github.com/Azure/azure-rest-api-specs/issues/13341
		"template_name": {
			Description: "The name of the Email Template.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				titleCase(string(emailtemplates.TemplateNameAccountClosedDeveloper)),
				titleCase(string(emailtemplates.TemplateNameApplicationApprovedNotificationMessage)),
				titleCase(string(emailtemplates.TemplateNameConfirmSignUpIdentityDefault)),
				titleCase(string(emailtemplates.TemplateNameEmailChangeIdentityDefault)),
				titleCase(string(emailtemplates.TemplateNameInviteUserNotificationMessage)),
				titleCase(string(emailtemplates.TemplateNameNewCommentNotificationMessage)),
				titleCase(string(emailtemplates.TemplateNameNewDeveloperNotificationMessage)),
				titleCase(string(emailtemplates.TemplateNameNewIssueNotificationMessage)),
				titleCase(string(emailtemplates.TemplateNamePasswordResetByAdminNotificationMessage)),
				titleCase(string(emailtemplates.TemplateNamePasswordResetIdentityDefault)),
				titleCase(string(emailtemplates.TemplateNamePurchaseDeveloperNotificationMessage)),
				titleCase(string(emailtemplates.TemplateNameQuotaLimitApproachingDeveloperNotificationMessage)),
				titleCase(string(emailtemplates.TemplateNameRejectDeveloperNotificationMessage)),
				titleCase(string(emailtemplates.TemplateNameRequestDeveloperNotificationMessage)),
			}, false),
		},
		"body": {
			Description:  "The body of the Email. Its format has to be a well-formed HTML document.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"subject": {
			Description: "The subject of the Email.",
			Type:        pluginsdk.TypeString,
			Required:    true,
		},
		// Computed:
		"title": {
			Description: "The title of the Email Template.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"description": {
			Description: "The description of the Email Template.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
