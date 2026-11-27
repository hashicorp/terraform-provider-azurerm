// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_automation

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2019-01-01-preview/automations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSecurityCenterAutomationCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SecurityCenter.AutomationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := automations.NewAutomationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_security_center_automation", id.ID())
			}
		}
	}

	location := location.Normalize(d.Get("location").(string))
	enabled := d.Get("enabled").(bool)

	// Build automation struct
	automation := automations.Automation{
		Location: &location,
		Properties: &automations.AutomationProperties{
			Description: pointer.To(d.Get("description").(string)),
			IsEnabled:   &enabled,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	automation.Properties.Scopes = expandSecurityCenterAutomationScopes(d.Get("scopes").([]any))

	var err error
	automation.Properties.Actions, err = expandSecurityCenterAutomationActions(d.Get("action").([]any))
	if err != nil {
		return err
	}

	automation.Properties.Sources, err = expandSecurityCenterAutomationSources(d.Get("source").([]any))
	if err != nil {
		return err
	}

	if _, err := client.CreateOrUpdate(ctx, id, automation); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSecurityCenterAutomationRead(d, meta)
}
