// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_azure_advanced_threat_protection

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	sentinelClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSentinelDataConnectorAzureAdvancedThreatProtectionDelete(d *pluginsdk.ResourceData, meta any) error {
	client := sentinelClient.From(meta).DataConnectorsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := dataconnectors.ParseDataConnectorID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
