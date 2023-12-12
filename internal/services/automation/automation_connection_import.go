// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func importAutomationConnection(connectionType string) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := connection.ParseConnectionID(d.Id())
		if err != nil {
			return []*pluginsdk.ResourceData{}, err
		}

		ctx, cancel := timeouts.ForRead(ctx, d)
		defer cancel()

		client := meta.(*clients.Client).Automation.Connection
		resp, err := client.Get(ctx, *id)
		if err != nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if model := resp.Model; model != nil {
			if resp.Model.Properties == nil || resp.Model.Properties.ConnectionType == nil || resp.Model.Properties.ConnectionType.Name == nil {
				return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: `properties`, `properties.connectionType` or `properties.connectionType.name` was nil", *id)
			}
			if *model.Properties.ConnectionType.Name != connectionType {
				return nil, fmt.Errorf(`automation connection "type" mismatch, expected "%s", got "%s"`, connectionType, *model.Properties.ConnectionType.Name)
			}
		} else {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s", *id)
		}
		return []*pluginsdk.ResourceData{d}, nil
	}
}
