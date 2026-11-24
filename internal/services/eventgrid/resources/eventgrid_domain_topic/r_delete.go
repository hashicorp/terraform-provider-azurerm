// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_domain_topic

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/domaintopics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventGridDomainTopicDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).EventGrid.DomainTopics
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := domaintopics.ParseDomainTopicID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
