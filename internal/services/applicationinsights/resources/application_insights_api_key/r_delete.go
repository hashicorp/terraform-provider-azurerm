// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_insights_api_key

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentapikeysapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	applicationinsightsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationInsightsAPIKeyDelete(d *pluginsdk.ResourceData, meta any) error {
	client := applicationinsightsClient.From(meta).APIKeysClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := componentapikeysapis.ParseApiKeyID(d.Id())
	if err != nil {
		return err
	}

	result, err := client.APIKeysDelete(ctx, *id)
	if err != nil {
		if response.WasNotFound(result.HttpResponse) {
			return nil
		}
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
