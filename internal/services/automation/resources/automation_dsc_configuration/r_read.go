// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_dsc_configuration

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/dscconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationDscConfigurationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.DscConfiguration
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dscconfiguration.ParseConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.ConfigurationName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("automation_account_name", id.AutomationAccountName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("log_verbose", props.LogVerbose)
			d.Set("description", props.Description)
			d.Set("state", pointer.FromEnum(props.State))
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	contentResp, err := client.GetContent(ctx, *id)
	if err != nil {
		return fmt.Errorf("making Read request on AzureRM Automation Dsc Configuration content %q: %+v", id.ConfigurationName, err)
	}

	if contentHttpResponse := contentResp.HttpResponse; contentHttpResponse != nil {
		if contentHttpResponse.Body != nil {
			buf := new(bytes.Buffer)
			if _, err := buf.ReadFrom(contentResp.HttpResponse.Body); err != nil {
				return fmt.Errorf("reading from AzureRM Automation Dsc Configuration buffer %q: %+v", id.ConfigurationName, err)
			}
			d.Set("content_embedded", buf.String())
		}
	}

	return nil
}
