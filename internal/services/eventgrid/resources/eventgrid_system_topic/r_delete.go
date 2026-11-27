// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_system_topic

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/systemtopics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	eventgridClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventGridSystemTopicDelete(d *pluginsdk.ResourceData, meta any) error {
	client := eventgridClient.From(meta).SystemTopics
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := systemtopics.ParseSystemTopicID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
