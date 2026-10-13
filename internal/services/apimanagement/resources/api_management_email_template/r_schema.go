// Copyright IBM Corp. 2014, 2025
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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
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
	}
}
