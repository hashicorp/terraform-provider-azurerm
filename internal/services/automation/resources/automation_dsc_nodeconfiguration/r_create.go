// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_dsc_nodeconfiguration

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/dscnodeconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationDscNodeConfigurationCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.DscNodeConfiguration
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := dscnodeconfiguration.NewNodeConfigurationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_automation_dsc_nodeconfiguration", id.ID())
		}
	}

	parameters := dscnodeconfiguration.DscNodeConfigurationCreateOrUpdateParameters{
		Properties: &dscnodeconfiguration.DscNodeConfigurationCreateOrUpdateParametersProperties{
			Source: dscnodeconfiguration.ContentSource{
				Type:  pointer.To(dscnodeconfiguration.ContentSourceTypeEmbeddedContent),
				Value: pointer.To(d.Get("content_embedded").(string)),
			},
			// configuration name is always the first part of the dsc node configuration
			// e.g. webserver.prod or webserver.local will be associated to the dsc configuration webserver
			Configuration: dscnodeconfiguration.DscConfigurationAssociationProperty{
				Name: pointer.To(strings.Split(id.NodeConfigurationName, ".")[0]),
			},
		},
		Name: pointer.To(id.NodeConfigurationName),
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceAutomationDscNodeConfigurationRead(d, meta)
}
