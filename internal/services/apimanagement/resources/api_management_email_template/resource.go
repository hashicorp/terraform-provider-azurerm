// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_email_template

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/emailtemplates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
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

		Schema: apiManagementEmailTemplateSchema(),
	}
}
