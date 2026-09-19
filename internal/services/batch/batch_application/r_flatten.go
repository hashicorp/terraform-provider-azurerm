// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch_application

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/application"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceBatchApplicationFlatten(d *pluginsdk.ResourceData, id *application.ApplicationId, model *application.Application) error {
	d.Set("name", id.ApplicationName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("account_name", id.BatchAccountName)

	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("allow_updates", props.AllowUpdates)
			d.Set("default_version", props.DefaultVersion)
			d.Set("display_name", props.DisplayName)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
