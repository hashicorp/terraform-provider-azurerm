// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_dsc_configuration

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/dscconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationDscConfigurationCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.DscConfiguration
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := dscconfiguration.NewConfigurationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_automation_dsc_configuration", id.ID())
		}
	}

	parameters := dscconfiguration.DscConfigurationCreateOrUpdateParameters{
		Properties: dscconfiguration.DscConfigurationCreateOrUpdateProperties{
			LogVerbose:  pointer.To(d.Get("log_verbose").(bool)),
			Description: pointer.To(d.Get("description").(string)),
			Source: dscconfiguration.ContentSource{
				Type:  pointer.To(dscconfiguration.ContentSourceTypeEmbeddedContent),
				Value: pointer.To(d.Get("content_embedded").(string)),
			},
		},
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:     pointer.To(helpers.ExpandStringInterfaceMap(d.Get("tags").(map[string]any))),
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceAutomationDscConfigurationRead(d, meta)
}
