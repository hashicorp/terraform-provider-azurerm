// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_dsc_nodeconfiguration

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/dscnodeconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationDscNodeConfigurationUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.DscNodeConfiguration
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dscnodeconfiguration.ParseNodeConfigurationID(d.Id())
	if err != nil {
		return err
	}

	// there is only one property that can be updated and is not force new: content_embedded
	// if content_embedded is changed, we need to update the dsc node configuration
	// otherwise we can just return as the resource is already in the desired state
	if d.HasChange("content_embedded") {
		parameters := dscnodeconfiguration.DscNodeConfigurationCreateOrUpdateParameters{
			Properties: &dscnodeconfiguration.DscNodeConfigurationCreateOrUpdateParametersProperties{
				Source: dscnodeconfiguration.ContentSource{
					Type: pointer.To(dscnodeconfiguration.ContentSourceTypeEmbeddedContent),
					// content_embedded is not returned by the API, so we get it from d.Get
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

		if err := client.CreateOrUpdateThenPoll(ctx, *id, parameters); err != nil {
			return fmt.Errorf("updating %s: %+v", *id, err)
		}
	}

	return resourceAutomationDscNodeConfigurationRead(d, meta)
}
