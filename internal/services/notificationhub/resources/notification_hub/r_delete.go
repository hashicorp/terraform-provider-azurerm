// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package notification_hub

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2023-09-01/hubs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	notificationhubClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/notificationhub/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNotificationHubDelete(d *pluginsdk.ResourceData, meta any) error {
	client := notificationhubClient.From(meta).HubsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := hubs.ParseNotificationHubID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.NotificationHubsDelete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
