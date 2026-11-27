// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_insights_web_test

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	applicationinsightsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationInsightsWebTestsDelete(d *pluginsdk.ResourceData, meta any) error {
	client := applicationinsightsClient.From(meta).WebTestsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := webtestsapis.ParseWebTestID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.WebTestsDelete(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return err
}
