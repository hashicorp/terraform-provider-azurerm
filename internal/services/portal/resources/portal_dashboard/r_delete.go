// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package portal_dashboard

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/dashboard"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	portalClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePortalDashboardDelete(d *pluginsdk.ResourceData, meta any) error {
	client := portalClient.From(meta).DashboardsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := dashboard.ParseDashboardID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
