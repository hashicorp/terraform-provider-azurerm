// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_automation

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2019-01-01-preview/automations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	securitycenterClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSecurityCenterAutomationRead(d *pluginsdk.ResourceData, meta any) error {
	client := securitycenterClient.From(meta).AutomationsClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := automations.ParseAutomationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AutomationName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("location", location.NormalizeNilable(resp.Model.Location))

	if properties := resp.Model.Properties; properties != nil {
		d.Set("description", properties.Description)
		d.Set("enabled", properties.IsEnabled)

		flatScopes, err := flattenSecurityCenterAutomationScopes(properties.Scopes)
		if err != nil {
			return err
		}
		if err := d.Set("scopes", flatScopes); err != nil {
			return fmt.Errorf("reading Security Center automation scopes: %+v", err)
		}

		flatActions, err := flattenSecurityCenterAutomationActions(properties.Actions, d)
		if err != nil {
			return err
		}
		if err = d.Set("action", flatActions); err != nil {
			return fmt.Errorf("reading Security Center automation actions: %+v", err)
		}

		flatSources, err := flattenSecurityCenterAutomationSources(properties.Sources)
		if err != nil {
			return err
		}
		if err = d.Set("source", flatSources); err != nil {
			return fmt.Errorf("reading Security Center automation sources: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Model.Tags)
}
